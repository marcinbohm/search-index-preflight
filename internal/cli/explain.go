package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

func runExplain(args []string, stdout, stderr io.Writer) int {
	ruleID, format, help, message := parseExplainArgs(args)
	if help {
		writeExplainHelp(stdout)
		return exitSuccess
	}
	if message != "" {
		fmt.Fprintln(stderr, message)
		return exitUsage
	}
	if format != "console" && format != "json" {
		fmt.Fprintf(stderr, "invalid --format %q; expected console or json\n", format)
		return exitUsage
	}

	item, ok, err := findRuleExplainItem(ruleID)
	if err != nil {
		fmt.Fprintf(stderr, "explain rule: %v\n", err)
		return exitInternal
	}
	if !ok {
		fmt.Fprintf(stderr, "unknown rule ID %q\n", ruleID)
		return exitUsage
	}

	if format == "json" {
		output := struct {
			Rule ruleExplainItem `json:"rule"`
		}{Rule: item}
		encoder := json.NewEncoder(stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(output); err != nil {
			fmt.Fprintf(stderr, "write explain JSON: %v\n", err)
			return exitInternal
		}
		return exitSuccess
	}

	writeExplainConsole(stdout, item)
	return exitSuccess
}

type ruleExplainItem struct {
	ID          string   `json:"id"`
	Family      string   `json:"family"`
	Name        string   `json:"name"`
	Category    string   `json:"category"`
	Severity    string   `json:"severity"`
	Confidence  string   `json:"confidence"`
	Determinism string   `json:"determinism"`
	Description string   `json:"description"`
	Notes       []string `json:"notes,omitempty"`
}

func parseExplainArgs(args []string) (ruleID string, format string, help bool, message string) {
	// parseExplainArgs is intentionally hand-rolled instead of using flag.FlagSet.
	// Go's flag package stops parsing flags after the first positional argument,
	// but explain supports both `explain DIF003 --format json` and
	// `explain --format json DIF003`.
	format = "console"
	var positionals []string
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if isHelp(arg) {
			return "", "", true, ""
		}
		if arg == "--format" {
			if i+1 >= len(args) {
				return "", "", false, "--format requires a value"
			}
			i++
			format = args[i]
			continue
		}
		if strings.HasPrefix(arg, "--format=") {
			format = strings.TrimPrefix(arg, "--format=")
			continue
		}
		if strings.HasPrefix(arg, "-") {
			return "", "", false, fmt.Sprintf("unknown explain flag %q", arg)
		}
		positionals = append(positionals, arg)
	}
	if len(positionals) == 0 {
		return "", "", false, "explain requires a rule ID"
	}
	if len(positionals) > 1 {
		return "", "", false, "explain accepts exactly one rule ID"
	}
	return positionals[0], format, false, ""
}

func findRuleExplainItem(ruleID string) (ruleExplainItem, bool, error) {
	items, err := collectRuleListItems("all")
	if err != nil {
		return ruleExplainItem{}, false, err
	}
	for _, item := range items {
		if strings.EqualFold(item.ID, ruleID) {
			return ruleExplainItem{
				ID:          item.ID,
				Family:      item.Family,
				Name:        item.Name,
				Category:    item.Category,
				Severity:    item.Severity,
				Confidence:  item.Confidence,
				Determinism: item.Determinism,
				Description: item.Description,
				Notes:       explainNotes(item.ID),
			}, true, nil
		}
	}
	return ruleExplainItem{}, false, nil
}

func explainNotes(id string) []string {
	switch id {
	case "SIL001":
		return []string{"SIL001 can emit warning findings near the field-count limit and error findings when the limit is exceeded."}
	default:
		return nil
	}
}

func writeExplainConsole(w io.Writer, item ruleExplainItem) {
	fmt.Fprintf(w, "ID:          %s\n", item.ID)
	fmt.Fprintf(w, "Family:      %s\n", item.Family)
	fmt.Fprintf(w, "Name:        %s\n", item.Name)
	fmt.Fprintf(w, "Category:    %s\n", item.Category)
	fmt.Fprintf(w, "Severity:    %s\n", item.Severity)
	fmt.Fprintf(w, "Confidence:  %s\n", item.Confidence)
	fmt.Fprintf(w, "Determinism: %s\n\n", item.Determinism)
	fmt.Fprintln(w, "Description:")
	fmt.Fprintf(w, "  %s\n", item.Description)
	if len(item.Notes) > 0 {
		fmt.Fprintln(w)
		fmt.Fprintln(w, "Notes:")
		for _, note := range item.Notes {
			fmt.Fprintf(w, "  - %s\n", note)
		}
	}
}

func writeExplainHelp(w io.Writer) {
	fmt.Fprint(w, `Usage:
  search-index-preflight explain <RULE_ID> [flags]

Explain a public lint or diff rule by stable rule ID, for example SIL001.

Flags:
  --format <format>   Output format: console or json
`)
}
