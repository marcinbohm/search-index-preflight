package report

import (
	"fmt"
	"io"

	"github.com/marcinbohm/search-index-preflight/internal/model"
)

func WriteConsole(w io.Writer, result model.RunResult) error {
	if len(result.Diagnostics) == 0 && len(result.Findings) == 0 {
		_, err := fmt.Fprintln(w, "SearchIndexPreflight: no diagnostics or findings")
		return err
	}

	for _, diagnostic := range result.Diagnostics {
		location := diagnostic.File
		if diagnostic.Line > 0 {
			location = fmt.Sprintf("%s:%d", location, diagnostic.Line)
		}
		if location == "" {
			location = "input"
		}
		if _, err := fmt.Fprintf(w, "%s: %s: %s\n", diagnostic.Severity, location, diagnostic.Message); err != nil {
			return err
		}
	}
	for _, finding := range result.Findings {
		location := finding.File
		if finding.JSONPointer != "" {
			location = fmt.Sprintf("%s#%s", location, finding.JSONPointer)
		}
		if location == "" {
			location = "input"
		}
		if _, err := fmt.Fprintf(w, "%s %s: %s: %s\n", finding.Severity, finding.ID, location, finding.Message); err != nil {
			return err
		}
		if finding.Remediation != "" {
			if _, err := fmt.Fprintf(w, "  Remediation: %s\n", finding.Remediation); err != nil {
				return err
			}
		}
	}
	return nil
}
