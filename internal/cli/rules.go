package cli

import (
	"fmt"
	"io"
)

func runRules(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 || isHelp(args[0]) {
		writeRulesHelp(stdout)
		return exitSuccess
	}

	switch args[0] {
	case "list":
		return runRulesList(args[1:], stdout, stderr)
	default:
		fmt.Fprintf(stderr, "unknown rules command %q\n\n", args[0])
		writeRulesHelp(stderr)
		return exitUsage
	}
}

func runRulesList(args []string, stdout, stderr io.Writer) int {
	if len(args) > 0 && isHelp(args[0]) {
		writeRulesListHelp(stdout)
		return exitSuccess
	}
	if len(args) > 0 {
		fmt.Fprintln(stderr, "rules list does not accept arguments")
		return exitUsage
	}
	fmt.Fprintln(stdout, "rules list implementation is in progress.")
	return exitSuccess
}

func writeRulesHelp(w io.Writer) {
	fmt.Fprint(w, `Usage:
  search-index-preflight rules <command>

Available Commands:
  list        List available rules
`)
}

func writeRulesListHelp(w io.Writer) {
	fmt.Fprint(w, `Usage:
  search-index-preflight rules list

List available rules.
Implementation is in progress.
`)
}
