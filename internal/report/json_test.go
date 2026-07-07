package report

import (
	"bytes"
	"testing"
)

func TestWriteJSONDeterministicEmptyResult(t *testing.T) {
	var first bytes.Buffer
	if err := WriteJSON(&first, EmptyRunResult()); err != nil {
		t.Fatalf("WriteJSON first run returned error: %v", err)
	}

	var second bytes.Buffer
	if err := WriteJSON(&second, EmptyRunResult()); err != nil {
		t.Fatalf("WriteJSON second run returned error: %v", err)
	}

	if first.String() != second.String() {
		t.Fatalf("WriteJSON output differs between runs:\nfirst:\n%s\nsecond:\n%s", first.String(), second.String())
	}

	const want = `{
  "schema_version": "0.1",
  "tool": {
    "name": "SearchIndexPreflight",
    "version": "0.0.0-dev"
  },
  "summary": {
    "files_scanned": 0,
    "findings_total": 0,
    "critical": 0,
    "error": 0,
    "warning": 0,
    "info": 0,
    "exit_code": 0
  },
  "findings": [],
  "diagnostics": []
}
`
	if first.String() != want {
		t.Fatalf("WriteJSON output =\n%s\nwant:\n%s", first.String(), want)
	}
}
