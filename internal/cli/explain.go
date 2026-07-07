package cli

import (
	"fmt"
	"io"
)

func runExplain(args []string, stdout, stderr io.Writer) int {
	if len(args) > 0 && isHelp(args[0]) {
		writeExplainHelp(stdout)
		return exitSuccess
	}
	if len(args) > 1 {
		fmt.Fprintln(stderr, "explain accepts at most one rule ID")
		return exitUsage
	}
	fmt.Fprintln(stdout, "explain implementation is in progress.")
	return exitSuccess
}

func writeExplainHelp(w io.Writer) {
	fmt.Fprint(w, `Usage:
  search-index-preflight explain [rule-id]

Explain a rule by stable rule ID, for example SIL001.
Implementation is in progress.
`)
}
