package parser

import (
	"testing"

	"github.com/marcinbohm/search-index-preflight/internal/model"
)

func TestParseJSONValid(t *testing.T) {
	document := ParseJSON(testSource("mapping.json"), model.DocumentKindMapping, []byte(`{"properties":{}}`))
	if len(document.Diagnostics) != 0 {
		t.Fatalf("ParseJSON returned diagnostics: %#v", document.Diagnostics)
	}
	if document.Content == nil {
		t.Fatal("ParseJSON returned nil content")
	}
}

func TestParseJSONInvalidReturnsDiagnostic(t *testing.T) {
	document := ParseJSON(testSource("mapping.json"), model.DocumentKindMapping, []byte(`{"properties":`))
	if len(document.Diagnostics) != 1 {
		t.Fatalf("ParseJSON returned %d diagnostics, want 1", len(document.Diagnostics))
	}
	if document.Diagnostics[0].File != "mapping.json" {
		t.Fatalf("diagnostic file = %q, want mapping.json", document.Diagnostics[0].File)
	}
}

func TestParseJSONLValid(t *testing.T) {
	document := ParseJSONL(testSource("samples.jsonl"), []byte("{\"a\":1}\n{\"b\":2}\n"))
	if len(document.Diagnostics) != 0 {
		t.Fatalf("ParseJSONL returned diagnostics: %#v", document.Diagnostics)
	}
	if document.Content == nil {
		t.Fatal("ParseJSONL returned nil content")
	}
}

func TestParseJSONLInvalidReturnsLineDiagnostic(t *testing.T) {
	document := ParseJSONL(testSource("samples.jsonl"), []byte("{\"a\":1}\n{\"b\":\n"))
	if len(document.Diagnostics) != 1 {
		t.Fatalf("ParseJSONL returned %d diagnostics, want 1", len(document.Diagnostics))
	}
	if document.Diagnostics[0].Line != 2 {
		t.Fatalf("diagnostic line = %d, want 2", document.Diagnostics[0].Line)
	}
}

func testSource(path string) model.Source {
	return model.Source{
		Path:         path,
		RelativePath: path,
	}
}
