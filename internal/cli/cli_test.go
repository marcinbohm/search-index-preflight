package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/marcinbohm/search-index-lint/internal/model"
)

func TestRootHelp(t *testing.T) {
	code, stdout, stderr := executeForTest("--help")
	if code != exitSuccess {
		t.Fatalf("Execute returned %d, want %d; stderr=%s", code, exitSuccess, stderr)
	}
	if !strings.Contains(stdout, "SearchIndexLint") {
		t.Fatalf("stdout %q does not contain SearchIndexLint", stdout)
	}
}

func TestVersion(t *testing.T) {
	code, stdout, stderr := executeForTest("version")
	if code != exitSuccess {
		t.Fatalf("Execute returned %d, want %d; stderr=%s", code, exitSuccess, stderr)
	}
	if strings.TrimSpace(stdout) != "SearchIndexLint version 0.0.0-dev" {
		t.Fatalf("version output = %q", stdout)
	}
}

func TestLintNoInputReturnsUsageError(t *testing.T) {
	code, _, stderr := executeForTest("lint")
	if code != exitUsage {
		t.Fatalf("Execute returned %d, want %d", code, exitUsage)
	}
	if !strings.Contains(stderr, "requires at least one input") {
		t.Fatalf("stderr %q does not explain missing input", stderr)
	}
}

func TestLintValidMappingJSONReturnsSuccess(t *testing.T) {
	path := writeTempFile(t, "mapping.json", `{"properties":{}}`)
	code, _, stderr := executeForTest("lint", "--mapping", path)
	if code != exitSuccess {
		t.Fatalf("Execute returned %d, want %d; stderr=%s", code, exitSuccess, stderr)
	}
}

func TestLintInvalidJSONReturnsInputError(t *testing.T) {
	path := writeTempFile(t, "mapping.json", `{"properties":`)
	code, stdout, stderr := executeForTest("lint", "--mapping", path)
	if code != exitInput {
		t.Fatalf("Execute returned %d, want %d; stderr=%s", code, exitInput, stderr)
	}
	if !strings.Contains(stdout, "invalid JSON") {
		t.Fatalf("stdout %q does not contain invalid JSON diagnostic", stdout)
	}
}

func TestLintInvalidJSONLReturnsInputErrorWithLineNumber(t *testing.T) {
	path := writeTempFile(t, "samples.jsonl", "{\"ok\":true}\n{\"bad\":\n")
	code, stdout, stderr := executeForTest("lint", "--sample-docs", path)
	if code != exitInput {
		t.Fatalf("Execute returned %d, want %d; stderr=%s", code, exitInput, stderr)
	}
	if !strings.Contains(stdout, ":2: invalid JSONL") {
		t.Fatalf("stdout %q does not contain JSONL line diagnostic", stdout)
	}
}

func TestLintFormatJSONReturnsValidJSON(t *testing.T) {
	path := writeTempFile(t, "mapping.json", `{"properties":{}}`)
	code, stdout, stderr := executeForTest("lint", "--mapping", path, "--format", "json")
	if code != exitSuccess {
		t.Fatalf("Execute returned %d, want %d; stderr=%s", code, exitSuccess, stderr)
	}

	var result model.RunResult
	if err := json.Unmarshal([]byte(stdout), &result); err != nil {
		t.Fatalf("stdout is not valid JSON: %v\n%s", err, stdout)
	}
	if result.SchemaVersion != "0.1" {
		t.Fatalf("schema_version = %q, want 0.1", result.SchemaVersion)
	}
	if result.Diagnostics == nil {
		t.Fatal("diagnostics is nil, want empty slice")
	}
}

func TestLintDirectoryDiscoveryIgnoresLocal(t *testing.T) {
	root := t.TempDir()
	writeFileAt(t, root, "mapping.json", `{"properties":{}}`)
	writeFileAt(t, root, filepath.Join(".local", "bad.json"), `{"properties":`)

	code, stdout, stderr := executeForTest("lint", root)
	if code != exitSuccess {
		t.Fatalf("Execute returned %d, want %d; stdout=%s stderr=%s", code, exitSuccess, stdout, stderr)
	}
}

func TestLintDirectorySkipsUnsupportedInvalidFiles(t *testing.T) {
	root := t.TempDir()
	writeFileAt(t, root, "mapping.json", `{"properties":{}}`)
	writeFileAt(t, root, "README.md", `{"properties":`)
	writeFileAt(t, root, "main.go", `{"properties":`)
	writeFileAt(t, root, "data.bin", `{"properties":`)

	code, stdout, stderr := executeForTest("lint", root)
	if code != exitSuccess {
		t.Fatalf("Execute returned %d, want %d; stdout=%s stderr=%s", code, exitSuccess, stdout, stderr)
	}
	if strings.Contains(stdout, "README.md") || strings.Contains(stdout, "main.go") || strings.Contains(stdout, "data.bin") {
		t.Fatalf("stdout contains diagnostics for unsupported files: %s", stdout)
	}
}

func TestLintDirectoryUnknownJSONDocumentKindReturnsInputError(t *testing.T) {
	root := t.TempDir()
	writeFileAt(t, root, "unknown.json", `{"name":"not-a-search-schema"}`)

	code, stdout, stderr := executeForTest("lint", root)
	if code != exitInput {
		t.Fatalf("Execute returned %d, want %d; stdout=%s stderr=%s", code, exitInput, stdout, stderr)
	}
	if !strings.Contains(stdout, "unknown JSON document kind") {
		t.Fatalf("stdout %q does not contain unknown kind diagnostic", stdout)
	}
}

func TestLintDirectoryValidMappingJSONReturnsSuccess(t *testing.T) {
	root := t.TempDir()
	writeFileAt(t, root, "mapping.json", `{"properties":{"status":{"type":"keyword"}}}`)

	code, stdout, stderr := executeForTest("lint", root)
	if code != exitSuccess {
		t.Fatalf("Execute returned %d, want %d; stdout=%s stderr=%s", code, exitSuccess, stdout, stderr)
	}
}

func TestLintDirectoryValidIndexTemplateJSONReturnsSuccess(t *testing.T) {
	root := t.TempDir()
	writeFileAt(t, root, "index-template.json", `{
  "index_patterns": ["logs-*"],
  "template": {
    "mappings": {
      "properties": {
        "@timestamp": {
          "type": "date"
        }
      }
    }
  }
}`)

	code, stdout, stderr := executeForTest("lint", root)
	if code != exitSuccess {
		t.Fatalf("Execute returned %d, want %d; stdout=%s stderr=%s", code, exitSuccess, stdout, stderr)
	}
}

func TestLintDirectoryValidComponentTemplateJSONReturnsSuccess(t *testing.T) {
	root := t.TempDir()
	writeFileAt(t, root, "component-template.json", `{
  "template": {
    "mappings": {
      "properties": {
        "service.name": {
          "type": "keyword"
        }
      }
    }
  }
}`)

	code, stdout, stderr := executeForTest("lint", root)
	if code != exitSuccess {
		t.Fatalf("Execute returned %d, want %d; stdout=%s stderr=%s", code, exitSuccess, stdout, stderr)
	}
}

func TestLintDirectoryValidSampleJSONLReturnsSuccess(t *testing.T) {
	root := t.TempDir()
	writeFileAt(t, root, "samples.jsonl", "{\"status\":\"ok\"}\n{\"status\":\"error\"}\n")

	code, stdout, stderr := executeForTest("lint", root)
	if code != exitSuccess {
		t.Fatalf("Execute returned %d, want %d; stdout=%s stderr=%s", code, exitSuccess, stdout, stderr)
	}
}

func executeForTest(args ...string) (int, string, string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Execute(args, &stdout, &stderr)
	return code, stdout.String(), stderr.String()
}

func writeTempFile(t *testing.T, name, content string) string {
	t.Helper()
	root := t.TempDir()
	writeFileAt(t, root, name, content)
	return filepath.Join(root, name)
}

func writeFileAt(t *testing.T, root, name, content string) {
	t.Helper()
	path := filepath.Join(root, name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}
}
