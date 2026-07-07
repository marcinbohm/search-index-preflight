package cli

import (
	"fmt"
	"io"

	"github.com/marcinbohm/search-index-preflight/internal/version"
)

func runVersion(args []string, stdout, stderr io.Writer) int {
	if len(args) > 0 && isHelp(args[0]) {
		writeVersionHelp(stdout)
		return exitSuccess
	}
	if len(args) > 0 {
		fmt.Fprintf(stderr, "version does not accept arguments\n")
		return exitUsage
	}
	fmt.Fprintf(stdout, "%s version %s\n", version.Name, version.Version)
	return exitSuccess
}

func writeVersionHelp(w io.Writer) {
	fmt.Fprint(w, `Usage:
  search-index-preflight version

Print version information.
`)
}
