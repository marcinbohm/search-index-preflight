package cli

import (
	"fmt"
	"io"
	"strings"
)

const (
	exitSuccess  = 0
	exitFindings = 1
	exitUsage    = 2
	exitInput    = 3
	exitInternal = 6
)

func Execute(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 || isHelp(args[0]) {
		writeRootHelp(stdout)
		return exitSuccess
	}

	switch args[0] {
	case "version":
		return runVersion(args[1:], stdout, stderr)
	case "lint":
		return runLint(args[1:], stdout, stderr)
	case "rules":
		return runRules(args[1:], stdout, stderr)
	case "explain":
		return runExplain(args[1:], stdout, stderr)
	default:
		fmt.Fprintf(stderr, "unknown command %q\n\n", args[0])
		writeRootHelp(stderr)
		return exitUsage
	}
}

func isHelp(arg string) bool {
	return arg == "--help" || arg == "-h" || arg == "help"
}

func writeRootHelp(w io.Writer) {
	fmt.Fprint(w, strings.TrimSpace(`SearchIndexPreflight lints Elasticsearch/OpenSearch schemas offline.

Usage:
  search-index-preflight [command]

Available Commands:
  lint        Lint mappings, templates, and sample documents
  rules       Inspect rule metadata
  explain     Explain a rule by ID
  version     Print version information

Use "search-index-preflight <command> --help" for more information about a command.
`)+"\n")
}
