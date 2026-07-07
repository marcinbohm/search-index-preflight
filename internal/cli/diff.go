package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/marcinbohm/search-index-preflight/internal/diff"
	"github.com/marcinbohm/search-index-preflight/internal/diffrules"
	"github.com/marcinbohm/search-index-preflight/internal/input"
	"github.com/marcinbohm/search-index-preflight/internal/model"
	"github.com/marcinbohm/search-index-preflight/internal/normalizer"
	"github.com/marcinbohm/search-index-preflight/internal/report"
)

func runDiff(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("diff", flag.ContinueOnError)
	flags.SetOutput(stderr)

	var basePath string
	var currentPath string
	var format string
	var output string
	var failOn string

	flags.StringVar(&basePath, "base", "", "Base schema file or directory")
	flags.StringVar(&currentPath, "current", "", "Current schema file or directory")
	flags.StringVar(&format, "format", "console", "Output format: console or json")
	flags.StringVar(&output, "output", "", "Output file path; default stdout")
	flags.StringVar(&failOn, "fail-on", "error", "Minimum severity that returns exit code 1")
	flags.Usage = func() { writeDiffHelp(flags.Output()) }

	if err := flags.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return exitSuccess
		}
		return exitUsage
	}

	if flags.NArg() != 0 {
		fmt.Fprintln(stderr, "diff does not accept positional arguments; use --base and --current")
		return exitUsage
	}
	if basePath == "" {
		fmt.Fprintln(stderr, "diff requires --base <path>")
		return exitUsage
	}
	if currentPath == "" {
		fmt.Fprintln(stderr, "diff requires --current <path>")
		return exitUsage
	}
	if format != "console" && format != "json" {
		fmt.Fprintf(stderr, "invalid --format %q; expected console or json\n", format)
		return exitUsage
	}
	failOnSeverity, err := model.ParseSeverity(failOn)
	if err != nil {
		fmt.Fprintf(stderr, "invalid --fail-on %q; expected info, warning, error, or critical\n", failOn)
		return exitUsage
	}

	baseIsFile := isRegularFile(basePath)
	currentIsFile := isRegularFile(currentPath)
	baseSources, baseDiagnostics := collectDiffSources(basePath)
	currentSources, currentDiagnostics := collectDiffSources(currentPath)

	baseInputs := normalizeDiffSources(baseSources, baseDiagnostics)
	currentInputs := normalizeDiffSources(currentSources, currentDiagnostics)

	diagnostics := append([]model.Diagnostic{}, baseInputs.Diagnostics...)
	diagnostics = append(diagnostics, currentInputs.Diagnostics...)
	diagnostics = append(diagnostics, baseInputs.Corpus.Diagnostics...)
	diagnostics = append(diagnostics, currentInputs.Corpus.Diagnostics...)

	exitCode := exitSuccess
	var findings []model.Finding
	if len(diagnostics) > 0 {
		exitCode = exitInput
	} else {
		if baseIsFile && currentIsFile {
			alignSingleFileCorpusResourceIdentity(baseInputs.Sources, currentInputs.Sources, &baseInputs.Corpus)
		}
		diffResult, err := diff.Compare(baseInputs.Corpus, currentInputs.Corpus)
		if err != nil {
			fmt.Fprintf(stderr, "compare schema inputs: %v\n", err)
			return exitInternal
		}
		registry, err := diffrules.BuiltinRegistry()
		if err != nil {
			fmt.Fprintf(stderr, "initialize built-in diff rules: %v\n", err)
			return exitInternal
		}
		ruleResult, err := diffrules.Run(diffrules.Context{}, registry, diffrules.RunRequest{Result: diffResult})
		if err != nil {
			fmt.Fprintf(stderr, "run diff rules: %v\n", err)
			return exitInternal
		}
		findings = ruleResult.Findings
		if hasFindingAtOrAbove(findings, failOnSeverity) {
			exitCode = exitFindings
		}
	}

	result := report.EmptyRunResult()
	result.Summary = summarizeRun(len(baseInputs.Sources)+len(currentInputs.Sources), findings, exitCode)
	if diagnostics != nil {
		result.Diagnostics = diagnostics
	}
	if findings != nil {
		result.Findings = findings
	}

	return writeReport(format, output, result, stdout, stderr, exitCode)
}

func writeDiffHelp(w io.Writer) {
	fmt.Fprint(w, `Usage:
  search-index-preflight diff --base <path> --current <path> [flags]

Compare two schema inputs and report preflight diff findings.
This minimal experimental command currently emits DIF001 field type changes, DIF002 field removals, and DIF003 field additions.

Flags:
  --base <path>       Base schema file or directory
  --current <path>    Current schema file or directory
  --format <format>   Output format: console or json
  --output <path>     Output file path; default stdout
  --fail-on <severity> info, warning, error, or critical
`)
}

type diffInput struct {
	Sources     []input.Source
	Corpus      model.Corpus
	Diagnostics []model.Diagnostic
}

func normalizeDiffSources(sources []input.Source, diagnostics []model.Diagnostic) diffInput {
	documents := parseLintSources(sources)
	for _, document := range documents {
		diagnostics = append(diagnostics, document.Diagnostics...)
	}
	corpus := normalizer.Normalize(documents)
	return diffInput{
		Sources:     sources,
		Corpus:      corpus,
		Diagnostics: diagnostics,
	}
}

func collectDiffSources(path string) ([]input.Source, []model.Diagnostic) {
	sources, err := input.Discover(path)
	if err != nil {
		return nil, []model.Diagnostic{
			{
				Severity: model.SeverityError,
				File:     path,
				Message:  err.Error(),
			},
		}
	}
	rebaseSources(path, sources)
	return sources, nil
}

func isRegularFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode().IsRegular()
}

func rebaseSources(root string, sources []input.Source) {
	info, err := os.Stat(root)
	if err != nil {
		return
	}
	if !info.IsDir() {
		for i := range sources {
			sources[i].RelativePath = filepath.Base(sources[i].Path)
		}
		return
	}
	for i := range sources {
		rel, err := filepath.Rel(root, sources[i].Path)
		if err != nil {
			continue
		}
		sources[i].RelativePath = filepath.Clean(rel)
	}
}

func alignSingleFileCorpusResourceIdentity(baseSources []input.Source, currentSources []input.Source, baseCorpus *model.Corpus) {
	if len(baseSources) != 1 || len(currentSources) != 1 {
		return
	}
	currentRelativePath := currentSources[0].RelativePath
	for i := range baseCorpus.Mappings {
		baseCorpus.Mappings[i].Source.RelativePath = currentRelativePath
	}
	for i := range baseCorpus.IndexTemplates {
		baseCorpus.IndexTemplates[i].Source.RelativePath = currentRelativePath
		if baseCorpus.IndexTemplates[i].Template.Mappings != nil {
			baseCorpus.IndexTemplates[i].Template.Mappings.Source.RelativePath = currentRelativePath
		}
	}
	for i := range baseCorpus.ComponentTemplates {
		baseCorpus.ComponentTemplates[i].Source.RelativePath = currentRelativePath
		if baseCorpus.ComponentTemplates[i].Template.Mappings != nil {
			baseCorpus.ComponentTemplates[i].Template.Mappings.Source.RelativePath = currentRelativePath
		}
	}
}

func writeReport(format string, output string, result model.RunResult, stdout, stderr io.Writer, exitCode int) int {
	w := stdout
	var file *os.File
	if output != "" {
		var err error
		file, err = os.Create(output)
		if err != nil {
			fmt.Fprintf(stderr, "write output %q: %v\n", output, err)
			return exitInput
		}
		defer file.Close()
		w = file
	}

	if format == "json" {
		if err := report.WriteJSON(w, result); err != nil {
			fmt.Fprintf(stderr, "write JSON report: %v\n", err)
			return exitInput
		}
	} else {
		if err := report.WriteConsole(w, result); err != nil {
			fmt.Fprintf(stderr, "write console report: %v\n", err)
			return exitInput
		}
	}
	return exitCode
}
