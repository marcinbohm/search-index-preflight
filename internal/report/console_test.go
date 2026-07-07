package report

import (
	"bytes"
	"strings"
	"testing"

	"github.com/marcinbohm/search-index-preflight/internal/model"
)

func TestWriteConsoleFormatsFindingJSONPointerWithFragmentSeparator(t *testing.T) {
	result := model.RunResult{
		Findings: []model.Finding{
			{
				ID:          "SIL001",
				Severity:    model.SeverityError,
				File:        "bigtemplate.json",
				JSONPointer: "/template/mappings",
				Message:     "test message",
			},
			{
				ID:          "SIL001",
				Severity:    model.SeverityError,
				File:        "mapping.json",
				JSONPointer: "",
				Message:     "root pointer",
			},
		},
	}

	var buf bytes.Buffer
	if err := WriteConsole(&buf, result); err != nil {
		t.Fatalf("WriteConsole returned error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "bigtemplate.json#/template/mappings") {
		t.Fatalf("expected output to contain JSON pointer fragment location, got:\n%s", output)
	}
	if strings.Contains(output, "bigtemplate.json/template/mappings") {
		t.Fatalf("expected output not to contain concatenated path/pointer location, got:\n%s", output)
	}
	if !strings.Contains(output, "mapping.json") {
		t.Fatalf("expected output to contain root finding file location, got:\n%s", output)
	}
	if strings.Contains(output, "mapping.json#/") {
		t.Fatalf("expected output not to contain root JSON pointer fragment location, got:\n%s", output)
	}
}
