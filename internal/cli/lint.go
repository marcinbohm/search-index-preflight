package cli

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/marcinbohm/search-index-preflight/internal/input"
	"github.com/marcinbohm/search-index-preflight/internal/model"
	"github.com/marcinbohm/search-index-preflight/internal/normalizer"
	"github.com/marcinbohm/search-index-preflight/internal/parser"
	"github.com/marcinbohm/search-index-preflight/internal/report"
	"github.com/marcinbohm/search-index-preflight/internal/rules"
)

func runLint(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("lint", flag.ContinueOnError)
	flags.SetOutput(stderr)

	var mapping string
	var template string
	var componentTemplate string
	var sampleDocs string
	var format string
	var output string
	var failOn string

	flags.StringVar(&mapping, "mapping", "", "Mapping JSON file")
	flags.StringVar(&template, "template", "", "Index template JSON file")
	flags.StringVar(&componentTemplate, "component-template", "", "Component template JSON file")
	flags.StringVar(&sampleDocs, "sample-docs", "", "JSONL/NDJSON sample documents")
	flags.StringVar(&format, "format", "console", "Output format: console or json")
	flags.StringVar(&output, "output", "", "Output file path; default stdout")
	flags.StringVar(&failOn, "fail-on", "error", "Minimum severity that returns exit code 1")
	flags.Usage = func() { writeLintHelp(flags.Output()) }

	if err := flags.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return exitSuccess
		}
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

	if flags.NArg() > 1 {
		fmt.Fprintln(stderr, "lint accepts at most one positional directory path")
		return exitUsage
	}

	sources, diagnostics := collectLintSources(mapping, template, componentTemplate, sampleDocs, flags.Args())
	if len(sources) == 0 && len(diagnostics) == 0 {
		fmt.Fprintln(stderr, "lint requires at least one input: --mapping, --template, --component-template, --sample-docs, or a directory path")
		return exitUsage
	}

	documents := parseLintSources(sources)
	for _, document := range documents {
		diagnostics = append(diagnostics, document.Diagnostics...)
	}
	corpus := normalizer.Normalize(documents)
	diagnostics = append(diagnostics, corpus.Diagnostics...)

	exitCode := exitSuccess
	if len(diagnostics) > 0 {
		exitCode = exitInput
	}
	var findings []model.Finding
	if len(diagnostics) == 0 {
		registry, err := rules.BuiltinRegistry()
		if err != nil {
			fmt.Fprintf(stderr, "initialize built-in rules: %v\n", err)
			return exitInternal
		}
		ruleResult, err := rules.Run(rules.Context{}, registry, rules.RunRequest{Corpus: corpus})
		if err != nil {
			fmt.Fprintf(stderr, "run rules: %v\n", err)
			return exitInternal
		}
		findings = ruleResult.Findings
		if hasFindingAtOrAbove(findings, failOnSeverity) {
			exitCode = exitFindings
		}
	}

	result := report.EmptyRunResult()
	result.Summary = summarizeRun(len(sources), findings, exitCode)
	if diagnostics != nil {
		result.Diagnostics = diagnostics
	}
	if findings != nil {
		result.Findings = findings
	}

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

func writeLintHelp(w io.Writer) {
	fmt.Fprint(w, `Usage:
  search-index-preflight lint [path] [flags]

Lint mappings, templates, component templates, and sample documents.
Runs parsing, normalization, and currently implemented built-in rules.

Flags:
  --mapping <path>              Mapping JSON file
  --template <path>             Index template JSON file
  --component-template <path>   Component template JSON file
  --sample-docs <path>          JSONL/NDJSON sample documents
  --format <format>             Output format: console or json
  --output <path>               Output file path; default stdout
  --fail-on <severity>          info, warning, error, or critical
`)
}

func collectLintSources(mapping, template, componentTemplate, sampleDocs string, positional []string) ([]input.Source, []model.Diagnostic) {
	var sources []input.Source
	var diagnostics []model.Diagnostic

	load := func(path string, kind model.DocumentKind) {
		if path == "" {
			return
		}
		source, err := input.LoadFile(path, kind)
		if err != nil {
			diagnostics = append(diagnostics, model.Diagnostic{
				Severity: model.SeverityError,
				File:     path,
				Message:  err.Error(),
			})
			return
		}
		sources = append(sources, source)
	}

	load(mapping, model.DocumentKindMapping)
	load(template, model.DocumentKindIndexTemplate)
	load(componentTemplate, model.DocumentKindComponentTemplate)
	load(sampleDocs, model.DocumentKindSampleDocs)

	if len(positional) == 1 {
		discovered, err := input.Discover(positional[0])
		if err != nil {
			diagnostics = append(diagnostics, model.Diagnostic{
				Severity: model.SeverityError,
				File:     positional[0],
				Message:  err.Error(),
			})
		} else {
			sources = append(sources, discovered...)
		}
	}

	return sources, diagnostics
}

func parseLintSources(sources []input.Source) []model.RawDocument {
	documents := make([]model.RawDocument, 0, len(sources))
	for _, source := range sources {
		modelSource := model.Source{
			Path:         source.Path,
			RelativePath: source.RelativePath,
		}
		documents = append(documents, parser.Parse(modelSource, source.Kind, source.Content))
	}
	return documents
}

func hasFindingAtOrAbove(findings []model.Finding, threshold model.Severity) bool {
	for _, finding := range findings {
		if finding.Severity.AtLeast(threshold) {
			return true
		}
	}
	return false
}

func summarizeRun(filesScanned int, findings []model.Finding, exitCode int) model.Summary {
	summary := model.Summary{
		FilesScanned:  filesScanned,
		FindingsTotal: len(findings),
		ExitCode:      exitCode,
	}
	for _, finding := range findings {
		switch finding.Severity {
		case model.SeverityCritical:
			summary.Critical++
		case model.SeverityError:
			summary.Error++
		case model.SeverityWarning:
			summary.Warning++
		case model.SeverityInfo:
			summary.Info++
		}
	}
	return summary
}
