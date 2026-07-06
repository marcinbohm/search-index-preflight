package cli

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/marcinbohm/search-index-lint/internal/input"
	"github.com/marcinbohm/search-index-lint/internal/model"
	"github.com/marcinbohm/search-index-lint/internal/normalizer"
	"github.com/marcinbohm/search-index-lint/internal/parser"
	"github.com/marcinbohm/search-index-lint/internal/report"
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
	flags.StringVar(&sampleDocs, "sample-docs", "", "JSONL sample documents")
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

	result := report.EmptyRunResult()
	result.Summary = model.Summary{
		FilesScanned: len(sources),
		ExitCode:     exitCode,
	}
	if diagnostics != nil {
		result.Diagnostics = diagnostics
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
  search-index-lint lint [path] [flags]

Lint mappings, templates, component templates, and sample documents.
Implementation is in progress.

Flags:
  --mapping <path>              Mapping JSON file
  --template <path>             Index template JSON file
  --component-template <path>   Component template JSON file
  --sample-docs <path>          JSONL sample documents
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
