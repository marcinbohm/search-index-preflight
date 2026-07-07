package report

import (
	"encoding/json"
	"io"

	"github.com/marcinbohm/search-index-preflight/internal/model"
	"github.com/marcinbohm/search-index-preflight/internal/version"
)

const SchemaVersion = "0.1"

func EmptyRunResult() model.RunResult {
	return model.RunResult{
		SchemaVersion: SchemaVersion,
		Tool: model.Tool{
			Name:    version.Name,
			Version: version.Version,
		},
		Findings:    []model.Finding{},
		Diagnostics: []model.Diagnostic{},
	}
}

func WriteJSON(w io.Writer, result model.RunResult) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}
