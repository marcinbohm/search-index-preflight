package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/marcinbohm/search-index-preflight/internal/model"
)

func TestRootHelp(t *testing.T) {
	code, stdout, stderr := executeForTest("--help")
	if code != exitSuccess {
		t.Fatalf("Execute returned %d, want %d; stderr=%s", code, exitSuccess, stderr)
	}
	if !strings.Contains(stdout, "SearchIndexPreflight") {
		t.Fatalf("stdout %q does not contain SearchIndexPreflight", stdout)
	}
	if !strings.Contains(stdout, "diff") {
		t.Fatalf("stdout %q does not list diff command", stdout)
	}
}

func TestVersion(t *testing.T) {
	code, stdout, stderr := executeForTest("version")
	if code != exitSuccess {
		t.Fatalf("Execute returned %d, want %d; stderr=%s", code, exitSuccess, stderr)
	}
	if strings.TrimSpace(stdout) != "SearchIndexPreflight version 0.0.0-dev" {
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

func TestLintFormatJSONWithSIL001Finding(t *testing.T) {
	path := writeTempFile(t, "mapping.json", mappingJSONWithFields(1000))
	code, stdout, stderr := executeForTest("lint", "--mapping", path, "--format", "json")
	if code != exitFindings {
		t.Fatalf("Execute returned %d, want %d; stderr=%s", code, exitFindings, stderr)
	}

	var result model.RunResult
	if err := json.Unmarshal([]byte(stdout), &result); err != nil {
		t.Fatalf("stdout is not valid JSON: %v\n%s", err, stdout)
	}
	if result.Summary.FindingsTotal != 1 {
		t.Fatalf("findings_total = %d, want 1", result.Summary.FindingsTotal)
	}
	if result.Summary.Error != 1 {
		t.Fatalf("summary.error = %d, want 1", result.Summary.Error)
	}
	if result.Summary.ExitCode != exitFindings {
		t.Fatalf("summary.exit_code = %d, want %d", result.Summary.ExitCode, exitFindings)
	}
	if len(result.Findings) != 1 {
		t.Fatalf("findings length = %d, want 1", len(result.Findings))
	}
	if result.Findings[0].ID != "SIL001" {
		t.Fatalf("finding ID = %q, want SIL001", result.Findings[0].ID)
	}
	if result.Findings[0].Severity != model.SeverityError {
		t.Fatalf("finding severity = %q, want %q", result.Findings[0].Severity, model.SeverityError)
	}
	if result.Diagnostics == nil {
		t.Fatal("diagnostics is nil, want empty slice")
	}
	if len(result.Diagnostics) != 0 {
		t.Fatalf("diagnostics length = %d, want 0", len(result.Diagnostics))
	}
}

func TestLintFormatJSONWithWrappedSIL001Finding(t *testing.T) {
	path := writeTempFile(t, "wrapped.json", wrappedMappingJSONWithFields(1000))
	code, stdout, stderr := executeForTest("lint", "--mapping", path, "--format", "json")
	if code != exitFindings {
		t.Fatalf("Execute returned %d, want %d; stderr=%s", code, exitFindings, stderr)
	}

	var result model.RunResult
	if err := json.Unmarshal([]byte(stdout), &result); err != nil {
		t.Fatalf("stdout is not valid JSON: %v\n%s", err, stdout)
	}
	if len(result.Findings) != 1 {
		t.Fatalf("findings length = %d, want 1", len(result.Findings))
	}
	finding := result.Findings[0]
	if finding.ID != "SIL001" {
		t.Fatalf("finding ID = %q, want SIL001", finding.ID)
	}
	if finding.JSONPointer != "/mappings" {
		t.Fatalf("finding JSON pointer = %q, want /mappings", finding.JSONPointer)
	}
}

func TestLintSIL001FixtureExpectedJSONReports(t *testing.T) {
	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd returned error: %v", err)
	}
	repoRoot := filepath.Clean(filepath.Join("..", ".."))
	if err := os.Chdir(repoRoot); err != nil {
		t.Fatalf("Chdir returned error: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(originalWD); err != nil {
			t.Fatalf("restore working directory: %v", err)
		}
	})

	tests := []struct {
		name     string
		mapping  string
		expected string
		code     int
	}{
		{
			name:     "near limit",
			mapping:  "fixtures/mapping-limits/sil001-total-fields-limit/mapping-near-limit.json",
			expected: "fixtures/mapping-limits/sil001-total-fields-limit/expected-near-limit.json",
			code:     exitSuccess,
		},
		{
			name:     "over limit",
			mapping:  "fixtures/mapping-limits/sil001-total-fields-limit/mapping-over-limit.json",
			expected: "fixtures/mapping-limits/sil001-total-fields-limit/expected-over-limit.json",
			code:     exitFindings,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, stdout, stderr := executeForTest("lint", "--mapping", tt.mapping, "--format", "json")
			if code != tt.code {
				t.Fatalf("Execute returned %d, want %d; stdout=%s stderr=%s", code, tt.code, stdout, stderr)
			}

			expected, err := os.ReadFile(tt.expected)
			if err != nil {
				t.Fatalf("ReadFile returned error: %v", err)
			}
			if stdout != string(expected) {
				t.Fatalf("JSON report mismatch\nactual:\n%s\nexpected:\n%s", stdout, string(expected))
			}
		})
	}
}

func TestLintSIL002FixtureExpectedJSONReport(t *testing.T) {
	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd returned error: %v", err)
	}
	repoRoot := filepath.Clean(filepath.Join("..", ".."))
	if err := os.Chdir(repoRoot); err != nil {
		t.Fatalf("Chdir returned error: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(originalWD); err != nil {
			t.Fatalf("restore working directory: %v", err)
		}
	})

	mapping := "fixtures/dynamic-mapping/sil002-root-dynamic-enabled/mapping-root-dynamic-true.json"
	expectedPath := "fixtures/dynamic-mapping/sil002-root-dynamic-enabled/expected-root-dynamic-true.json"

	code, stdout, stderr := executeForTest("lint", "--mapping", mapping, "--format", "json")
	if code != exitSuccess {
		t.Fatalf("Execute returned %d, want %d; stdout=%s stderr=%s", code, exitSuccess, stdout, stderr)
	}

	expected, err := os.ReadFile(expectedPath)
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}
	if stdout != string(expected) {
		t.Fatalf("JSON report mismatch\nactual:\n%s\nexpected:\n%s", stdout, string(expected))
	}
}

func TestLintSIL003FixtureExpectedJSONReport(t *testing.T) {
	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd returned error: %v", err)
	}
	repoRoot := filepath.Clean(filepath.Join("..", ".."))
	if err := os.Chdir(repoRoot); err != nil {
		t.Fatalf("Chdir returned error: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(originalWD); err != nil {
			t.Fatalf("restore working directory: %v", err)
		}
	})

	mapping := "fixtures/dynamic-templates/sil003-missing-match-mapping-type/mapping-missing-match-mapping-type.json"
	expectedPath := "fixtures/dynamic-templates/sil003-missing-match-mapping-type/expected-missing-match-mapping-type.json"

	code, stdout, stderr := executeForTest("lint", "--mapping", mapping, "--format", "json")
	if code != exitSuccess {
		t.Fatalf("Execute returned %d, want %d; stdout=%s stderr=%s", code, exitSuccess, stdout, stderr)
	}

	expected, err := os.ReadFile(expectedPath)
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}
	if stdout != string(expected) {
		t.Fatalf("JSON report mismatch\nactual:\n%s\nexpected:\n%s", stdout, string(expected))
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

func TestLintSmallMappingReturnsSuccessWithoutFindings(t *testing.T) {
	path := writeTempFile(t, "mapping.json", `{"properties":{"status":{"type":"keyword"}}}`)

	code, stdout, stderr := executeForTest("lint", "--mapping", path)
	if code != exitSuccess {
		t.Fatalf("Execute returned %d, want %d; stdout=%s stderr=%s", code, exitSuccess, stdout, stderr)
	}
	if strings.Contains(stdout, "SIL001") {
		t.Fatalf("stdout contains unexpected SIL001 finding: %s", stdout)
	}
}

func TestLintMappingOverLimitReturnsFindingsExitCode(t *testing.T) {
	path := writeTempFile(t, "mapping.json", mappingJSONWithFields(1000))

	code, stdout, stderr := executeForTest("lint", "--mapping", path)
	if code != exitFindings {
		t.Fatalf("Execute returned %d, want %d; stdout=%s stderr=%s", code, exitFindings, stdout, stderr)
	}
	if !strings.Contains(stdout, "SIL001") {
		t.Fatalf("stdout %q does not contain SIL001", stdout)
	}
	if !strings.Contains(stdout, "total fields") {
		t.Fatalf("stdout %q does not contain total fields message", stdout)
	}
}

func TestLintMappingNearLimitDefaultFailOnErrorReturnsSuccess(t *testing.T) {
	path := writeTempFile(t, "mapping.json", mappingJSONWithFields(800))

	code, stdout, stderr := executeForTest("lint", "--mapping", path)
	if code != exitSuccess {
		t.Fatalf("Execute returned %d, want %d; stdout=%s stderr=%s", code, exitSuccess, stdout, stderr)
	}
	if !strings.Contains(stdout, "SIL001") {
		t.Fatalf("stdout %q does not contain SIL001 warning", stdout)
	}
}

func TestLintMappingNearLimitFailOnWarningReturnsFindingsExitCode(t *testing.T) {
	path := writeTempFile(t, "mapping.json", mappingJSONWithFields(800))

	code, stdout, stderr := executeForTest("lint", "--mapping", path, "--fail-on", "warning")
	if code != exitFindings {
		t.Fatalf("Execute returned %d, want %d; stdout=%s stderr=%s", code, exitFindings, stdout, stderr)
	}
	if !strings.Contains(stdout, "SIL001") {
		t.Fatalf("stdout %q does not contain SIL001 warning", stdout)
	}
}

func TestLintRootDynamicTrueWarningReturnsSuccessByDefault(t *testing.T) {
	path := writeTempFile(t, "mapping.json", mappingJSONWithRootDynamic(true))

	code, stdout, stderr := executeForTest("lint", "--mapping", path)
	if code != exitSuccess {
		t.Fatalf("Execute returned %d, want %d; stdout=%s stderr=%s", code, exitSuccess, stdout, stderr)
	}
	if !strings.Contains(stdout, "SIL002") {
		t.Fatalf("stdout %q does not contain SIL002", stdout)
	}
	if !strings.Contains(stdout, "warning") {
		t.Fatalf("stdout %q does not contain warning", stdout)
	}
}

func TestLintWrappedRootDynamicTrueUsesWrappedJSONPointer(t *testing.T) {
	path := writeTempFile(t, "wrapped.json", `{
  "mappings": {
    "dynamic": true,
    "properties": {
      "a": {
        "type": "keyword"
      }
    }
  }
}`)

	code, stdout, stderr := executeForTest("lint", "--mapping", path)
	if code != exitSuccess {
		t.Fatalf("Execute returned %d, want %d; stdout=%s stderr=%s", code, exitSuccess, stdout, stderr)
	}
	if !strings.Contains(stdout, "#/mappings/dynamic") {
		t.Fatalf("stdout %q does not contain wrapped dynamic JSON pointer", stdout)
	}
	if strings.Contains(stdout, "wrapped.json#/dynamic") {
		t.Fatalf("stdout contains raw mapping dynamic pointer for wrapped mapping: %s", stdout)
	}
}

func TestLintRootDynamicTrueFailOnWarningReturnsFindingsExitCode(t *testing.T) {
	path := writeTempFile(t, "mapping.json", mappingJSONWithRootDynamic(true))

	code, stdout, stderr := executeForTest("lint", "--mapping", path, "--fail-on", "warning")
	if code != exitFindings {
		t.Fatalf("Execute returned %d, want %d; stdout=%s stderr=%s", code, exitFindings, stdout, stderr)
	}
	if !strings.Contains(stdout, "SIL002") {
		t.Fatalf("stdout %q does not contain SIL002", stdout)
	}
}

func TestLintRootDynamicFalseDoesNotEmitSIL002(t *testing.T) {
	path := writeTempFile(t, "mapping.json", mappingJSONWithRootDynamic(false))

	code, stdout, stderr := executeForTest("lint", "--mapping", path)
	if code != exitSuccess {
		t.Fatalf("Execute returned %d, want %d; stdout=%s stderr=%s", code, exitSuccess, stdout, stderr)
	}
	if strings.Contains(stdout, "SIL002") {
		t.Fatalf("stdout contains unexpected SIL002 finding: %s", stdout)
	}
}

func TestLintFormatJSONWithSIL002Finding(t *testing.T) {
	path := writeTempFile(t, "mapping.json", mappingJSONWithRootDynamic(true))

	code, stdout, stderr := executeForTest("lint", "--mapping", path, "--format", "json")
	if code != exitSuccess {
		t.Fatalf("Execute returned %d, want %d; stderr=%s", code, exitSuccess, stderr)
	}

	var result model.RunResult
	if err := json.Unmarshal([]byte(stdout), &result); err != nil {
		t.Fatalf("stdout is not valid JSON: %v\n%s", err, stdout)
	}
	if result.Summary.FindingsTotal != 1 {
		t.Fatalf("findings_total = %d, want 1", result.Summary.FindingsTotal)
	}
	if result.Summary.Warning != 1 {
		t.Fatalf("summary.warning = %d, want 1", result.Summary.Warning)
	}
	if result.Summary.ExitCode != exitSuccess {
		t.Fatalf("summary.exit_code = %d, want %d", result.Summary.ExitCode, exitSuccess)
	}
	if len(result.Findings) != 1 {
		t.Fatalf("findings length = %d, want 1", len(result.Findings))
	}
	finding := result.Findings[0]
	if finding.ID != "SIL002" {
		t.Fatalf("finding ID = %q, want SIL002", finding.ID)
	}
	if finding.Severity != model.SeverityWarning {
		t.Fatalf("finding severity = %q, want %q", finding.Severity, model.SeverityWarning)
	}
	if finding.Confidence != model.ConfidenceMedium {
		t.Fatalf("finding confidence = %q, want %q", finding.Confidence, model.ConfidenceMedium)
	}
	if finding.Determinism != model.DeterminismHeuristic {
		t.Fatalf("finding determinism = %q, want %q", finding.Determinism, model.DeterminismHeuristic)
	}
}

func TestLintFormatJSONWithWrappedSIL002Finding(t *testing.T) {
	path := writeTempFile(t, "wrapped.json", `{
  "mappings": {
    "dynamic": true,
    "properties": {
      "status": {
        "type": "keyword"
      }
    }
  }
}`)

	code, stdout, stderr := executeForTest("lint", "--mapping", path, "--format", "json")
	if code != exitSuccess {
		t.Fatalf("Execute returned %d, want %d; stderr=%s", code, exitSuccess, stderr)
	}

	var result model.RunResult
	if err := json.Unmarshal([]byte(stdout), &result); err != nil {
		t.Fatalf("stdout is not valid JSON: %v\n%s", err, stdout)
	}
	if result.Summary.Warning != 1 {
		t.Fatalf("summary.warning = %d, want 1", result.Summary.Warning)
	}
	if len(result.Findings) != 1 {
		t.Fatalf("findings length = %d, want 1", len(result.Findings))
	}
	finding := result.Findings[0]
	if finding.ID != "SIL002" {
		t.Fatalf("finding ID = %q, want SIL002", finding.ID)
	}
	if finding.JSONPointer != "/mappings/dynamic" {
		t.Fatalf("finding JSON pointer = %q, want /mappings/dynamic", finding.JSONPointer)
	}
	if finding.Severity != model.SeverityWarning {
		t.Fatalf("finding severity = %q, want %q", finding.Severity, model.SeverityWarning)
	}
}

func TestLintDynamicTemplateMissingMatchMappingTypeWarningReturnsSuccessByDefault(t *testing.T) {
	path := writeTempFile(t, "mapping.json", mappingJSONWithDynamicTemplate(false))

	code, stdout, stderr := executeForTest("lint", "--mapping", path)
	if code != exitSuccess {
		t.Fatalf("Execute returned %d, want %d; stdout=%s stderr=%s", code, exitSuccess, stdout, stderr)
	}
	if !strings.Contains(stdout, "SIL003") {
		t.Fatalf("stdout %q does not contain SIL003", stdout)
	}
	if !strings.Contains(stdout, "warning") {
		t.Fatalf("stdout %q does not contain warning", stdout)
	}
}

func TestLintDynamicTemplateMissingMatchMappingTypeFailOnWarningReturnsFindingsExitCode(t *testing.T) {
	path := writeTempFile(t, "mapping.json", mappingJSONWithDynamicTemplate(false))

	code, stdout, stderr := executeForTest("lint", "--mapping", path, "--fail-on", "warning")
	if code != exitFindings {
		t.Fatalf("Execute returned %d, want %d; stdout=%s stderr=%s", code, exitFindings, stdout, stderr)
	}
	if !strings.Contains(stdout, "SIL003") {
		t.Fatalf("stdout %q does not contain SIL003", stdout)
	}
}

func TestLintDynamicTemplateWithMatchMappingTypeDoesNotEmitSIL003(t *testing.T) {
	path := writeTempFile(t, "mapping.json", mappingJSONWithDynamicTemplate(true))

	code, stdout, stderr := executeForTest("lint", "--mapping", path)
	if code != exitSuccess {
		t.Fatalf("Execute returned %d, want %d; stdout=%s stderr=%s", code, exitSuccess, stdout, stderr)
	}
	if strings.Contains(stdout, "SIL003") {
		t.Fatalf("stdout contains unexpected SIL003 finding: %s", stdout)
	}
}

func TestLintFormatJSONWithSIL003Finding(t *testing.T) {
	path := writeTempFile(t, "mapping.json", mappingJSONWithDynamicTemplate(false))

	code, stdout, stderr := executeForTest("lint", "--mapping", path, "--format", "json")
	if code != exitSuccess {
		t.Fatalf("Execute returned %d, want %d; stderr=%s", code, exitSuccess, stderr)
	}

	var result model.RunResult
	if err := json.Unmarshal([]byte(stdout), &result); err != nil {
		t.Fatalf("stdout is not valid JSON: %v\n%s", err, stdout)
	}
	if result.Summary.Warning != 1 {
		t.Fatalf("summary.warning = %d, want 1", result.Summary.Warning)
	}
	if result.Summary.ExitCode != exitSuccess {
		t.Fatalf("summary.exit_code = %d, want %d", result.Summary.ExitCode, exitSuccess)
	}
	if len(result.Findings) != 1 {
		t.Fatalf("findings length = %d, want 1", len(result.Findings))
	}
	finding := result.Findings[0]
	if finding.ID != "SIL003" {
		t.Fatalf("finding ID = %q, want SIL003", finding.ID)
	}
	if finding.Severity != model.SeverityWarning {
		t.Fatalf("finding severity = %q, want %q", finding.Severity, model.SeverityWarning)
	}
	if finding.Confidence != model.ConfidenceMedium {
		t.Fatalf("finding confidence = %q, want %q", finding.Confidence, model.ConfidenceMedium)
	}
	if finding.Determinism != model.DeterminismHeuristic {
		t.Fatalf("finding determinism = %q, want %q", finding.Determinism, model.DeterminismHeuristic)
	}
	if finding.Category != "dynamic-templates" {
		t.Fatalf("finding category = %q, want dynamic-templates", finding.Category)
	}
	if finding.JSONPointer != "/dynamic_templates/0/strings_as_keywords" {
		t.Fatalf("finding JSON pointer = %q, want /dynamic_templates/0/strings_as_keywords", finding.JSONPointer)
	}
}

func TestLintInvalidFailOnReturnsUsageError(t *testing.T) {
	path := writeTempFile(t, "mapping.json", `{"properties":{"status":{"type":"keyword"}}}`)

	code, _, stderr := executeForTest("lint", "--mapping", path, "--fail-on", "banana")
	if code != exitUsage {
		t.Fatalf("Execute returned %d, want %d", code, exitUsage)
	}
	if !strings.Contains(stderr, "invalid --fail-on") {
		t.Fatalf("stderr %q does not explain invalid fail-on", stderr)
	}
}

func TestLintInvalidJSONDoesNotRunRules(t *testing.T) {
	path := writeTempFile(t, "mapping.json", `{"properties":`)

	code, stdout, stderr := executeForTest("lint", "--mapping", path)
	if code != exitInput {
		t.Fatalf("Execute returned %d, want %d; stdout=%s stderr=%s", code, exitInput, stdout, stderr)
	}
	if strings.Contains(stdout, "SIL001") {
		t.Fatalf("stdout contains rule finding despite parse error: %s", stdout)
	}
	if strings.Contains(stdout, "SIL002") {
		t.Fatalf("stdout contains rule finding despite parse error: %s", stdout)
	}
	if strings.Contains(stdout, "SIL003") {
		t.Fatalf("stdout contains rule finding despite parse error: %s", stdout)
	}
}

func TestDiffFieldTypeChangedReturnsFindingsExitCode(t *testing.T) {
	base, current := writeDiffMappingFiles(t, `{"properties":{"status":{"type":"keyword"}}}`, `{"properties":{"status":{"type":"long"}}}`)

	code, stdout, stderr := executeForTest("diff", "--base", base, "--current", current)
	if code != exitFindings {
		t.Fatalf("Execute returned %d, want %d; stdout=%s stderr=%s", code, exitFindings, stdout, stderr)
	}
	for _, text := range []string{"DIF001", "status", "keyword", "long"} {
		if !strings.Contains(stdout, text) {
			t.Fatalf("stdout %q does not contain %q", stdout, text)
		}
	}
	if !strings.Contains(stdout, "current.json") {
		t.Fatalf("stdout %q does not contain current file name", stdout)
	}
}

func TestDiffDirectorySameRelativePathEmitsDIF001(t *testing.T) {
	root := t.TempDir()
	baseDir := filepath.Join(root, "old")
	currentDir := filepath.Join(root, "new")
	writeFileAt(t, baseDir, filepath.Join("schemas", "mapping.json"), `{"properties":{"status":{"type":"keyword"}}}`)
	writeFileAt(t, currentDir, filepath.Join("schemas", "mapping.json"), `{"properties":{"status":{"type":"long"}}}`)

	code, stdout, stderr := executeForTest("diff", "--base", baseDir, "--current", currentDir)
	if code != exitFindings {
		t.Fatalf("Execute returned %d, want %d; stdout=%s stderr=%s", code, exitFindings, stdout, stderr)
	}
	for _, text := range []string{"DIF001", filepath.Join("schemas", "mapping.json"), "status", "keyword", "long"} {
		if !strings.Contains(stdout, text) {
			t.Fatalf("stdout %q does not contain %q", stdout, text)
		}
	}
}

func TestDiffDIF001FixtureEmitsFinding(t *testing.T) {
	base := fixturePath("diff", "dif001-field-type-changed", "base")
	current := fixturePath("diff", "dif001-field-type-changed", "current")

	code, stdout, stderr := executeForTest("diff", "--base", base, "--current", current)
	if code != exitFindings {
		t.Fatalf("Execute returned %d, want %d; stdout=%s stderr=%s", code, exitFindings, stdout, stderr)
	}
	for _, text := range []string{"DIF001", "mapping.json", "status", "keyword", "long"} {
		if !strings.Contains(stdout, text) {
			t.Fatalf("stdout %q does not contain %q", stdout, text)
		}
	}
}

func TestDiffDIF001FixtureJSONMatchesExpectedFindingFields(t *testing.T) {
	base := fixturePath("diff", "dif001-field-type-changed", "base")
	current := fixturePath("diff", "dif001-field-type-changed", "current")
	expected := readExpectedFinding(t, fixturePath("diff", "dif001-field-type-changed", "expected.finding.json"))

	code, stdout, stderr := executeForTest("diff", "--base", base, "--current", current, "--format", "json")
	if code != exitFindings {
		t.Fatalf("Execute returned %d, want %d; stdout=%s stderr=%s", code, exitFindings, stdout, stderr)
	}

	var result model.RunResult
	if err := json.Unmarshal([]byte(stdout), &result); err != nil {
		t.Fatalf("stdout is not valid JSON: %v\n%s", err, stdout)
	}
	if result.Summary.FindingsTotal != 1 {
		t.Fatalf("findings_total = %d, want 1", result.Summary.FindingsTotal)
	}
	if result.Summary.Error != 1 {
		t.Fatalf("summary.error = %d, want 1", result.Summary.Error)
	}
	if result.Summary.ExitCode != exitFindings {
		t.Fatalf("summary.exit_code = %d, want %d", result.Summary.ExitCode, exitFindings)
	}
	if len(result.Diagnostics) != 0 {
		t.Fatalf("diagnostics length = %d, want 0", len(result.Diagnostics))
	}
	if result.Diagnostics == nil {
		t.Fatal("diagnostics is nil, want empty slice")
	}
	if len(result.Findings) != 1 {
		t.Fatalf("findings length = %d, want 1", len(result.Findings))
	}
	assertFindingMatchesExpected(t, result.Findings[0], expected)
}

func TestDiffNoChangesFixtureReturnsSuccess(t *testing.T) {
	base := fixturePath("diff", "no-changes", "base")
	current := fixturePath("diff", "no-changes", "current")

	code, stdout, stderr := executeForTest("diff", "--base", base, "--current", current)
	if code != exitSuccess {
		t.Fatalf("Execute returned %d, want %d; stdout=%s stderr=%s", code, exitSuccess, stdout, stderr)
	}
	if !strings.Contains(stdout, "no diagnostics or findings") {
		t.Fatalf("stdout %q does not report clean diff", stdout)
	}
}

func TestDiffDirectoryDifferentFilenamesDoesNotAlignAsTypeChange(t *testing.T) {
	root := t.TempDir()
	baseDir := filepath.Join(root, "old")
	currentDir := filepath.Join(root, "new")
	writeFileAt(t, baseDir, "legacy.json", `{"properties":{"status":{"type":"keyword"}}}`)
	writeFileAt(t, currentDir, "mapping.json", `{"properties":{"status":{"type":"long"}}}`)

	code, stdout, stderr := executeForTest("diff", "--base", baseDir, "--current", currentDir)
	if code != exitSuccess {
		t.Fatalf("Execute returned %d, want %d; stdout=%s stderr=%s", code, exitSuccess, stdout, stderr)
	}
	if strings.Contains(stdout, "DIF001") {
		t.Fatalf("stdout contains unexpected DIF001 for different relative filenames: %s", stdout)
	}
	if !strings.Contains(stdout, "no diagnostics or findings") {
		t.Fatalf("stdout %q does not report clean diff", stdout)
	}
}

func TestDiffFileVsDirectoryDoesNotPanicOrForceAlignment(t *testing.T) {
	root := t.TempDir()
	base := filepath.Join(root, "base.json")
	currentDir := filepath.Join(root, "current")
	writeFileAt(t, root, "base.json", `{"properties":{"status":{"type":"keyword"}}}`)
	writeFileAt(t, currentDir, "mapping.json", `{"properties":{"status":{"type":"long"}}}`)

	code, stdout, stderr := executeForTest("diff", "--base", base, "--current", currentDir)
	if code != exitSuccess {
		t.Fatalf("Execute returned %d, want %d; stdout=%s stderr=%s", code, exitSuccess, stdout, stderr)
	}
	if strings.Contains(stdout, "DIF001") {
		t.Fatalf("stdout contains unexpected DIF001 for file-vs-directory diff: %s", stdout)
	}
}

func TestDiffNoChangesReturnsSuccess(t *testing.T) {
	base, current := writeDiffMappingFiles(t, `{"properties":{"status":{"type":"keyword"}}}`, `{"properties":{"status":{"type":"keyword"}}}`)

	code, stdout, stderr := executeForTest("diff", "--base", base, "--current", current)
	if code != exitSuccess {
		t.Fatalf("Execute returned %d, want %d; stdout=%s stderr=%s", code, exitSuccess, stdout, stderr)
	}
	if !strings.Contains(stdout, "no diagnostics or findings") {
		t.Fatalf("stdout %q does not report clean diff", stdout)
	}
}

func TestDiffHelpReturnsSuccess(t *testing.T) {
	code, stdout, stderr := executeForTest("diff", "--help")
	if code != exitSuccess {
		t.Fatalf("Execute returned %d, want %d; stdout=%s stderr=%s", code, exitSuccess, stdout, stderr)
	}
	help := stdout + stderr
	if !strings.Contains(help, "--base") {
		t.Fatalf("help output %q does not contain --base", help)
	}
	if !strings.Contains(help, "--current") {
		t.Fatalf("help output %q does not contain --current", help)
	}
}

func TestDiffRejectsPositionalArgs(t *testing.T) {
	base, current := writeDiffMappingFiles(t, `{"properties":{"status":{"type":"keyword"}}}`, `{"properties":{"status":{"type":"long"}}}`)

	code, _, stderr := executeForTest("diff", "--base", base, "--current", current, "extra")
	if code != exitUsage {
		t.Fatalf("Execute returned %d, want %d", code, exitUsage)
	}
	if !strings.Contains(stderr, "does not accept positional arguments") {
		t.Fatalf("stderr %q does not explain positional args", stderr)
	}
}

func TestDiffFormatJSONWithDIF001Finding(t *testing.T) {
	base, current := writeDiffMappingFiles(t, `{"properties":{"status":{"type":"keyword"}}}`, `{"properties":{"status":{"type":"long"}}}`)

	code, stdout, stderr := executeForTest("diff", "--base", base, "--current", current, "--format", "json")
	if code != exitFindings {
		t.Fatalf("Execute returned %d, want %d; stderr=%s", code, exitFindings, stderr)
	}

	var result model.RunResult
	if err := json.Unmarshal([]byte(stdout), &result); err != nil {
		t.Fatalf("stdout is not valid JSON: %v\n%s", err, stdout)
	}
	if result.Summary.FindingsTotal != 1 {
		t.Fatalf("findings_total = %d, want 1", result.Summary.FindingsTotal)
	}
	if result.Summary.Error != 1 {
		t.Fatalf("summary.error = %d, want 1", result.Summary.Error)
	}
	if result.Summary.ExitCode != exitFindings {
		t.Fatalf("summary.exit_code = %d, want %d", result.Summary.ExitCode, exitFindings)
	}
	if len(result.Findings) != 1 {
		t.Fatalf("findings length = %d, want 1", len(result.Findings))
	}
	finding := result.Findings[0]
	if finding.ID != "DIF001" {
		t.Fatalf("finding ID = %q, want DIF001", finding.ID)
	}
	if finding.Severity != model.SeverityError {
		t.Fatalf("finding severity = %q, want %q", finding.Severity, model.SeverityError)
	}
	if finding.Category != "schema-diff" {
		t.Fatalf("finding category = %q, want schema-diff", finding.Category)
	}
	if result.Diagnostics == nil {
		t.Fatal("diagnostics is nil, want empty slice")
	}
	if len(result.Diagnostics) != 0 {
		t.Fatalf("diagnostics length = %d, want 0", len(result.Diagnostics))
	}
}

func TestDiffOutputWritesJSONReport(t *testing.T) {
	base, current := writeDiffMappingFiles(t, `{"properties":{"status":{"type":"keyword"}}}`, `{"properties":{"status":{"type":"long"}}}`)
	output := filepath.Join(t.TempDir(), "report.json")

	code, stdout, stderr := executeForTest("diff", "--base", base, "--current", current, "--format", "json", "--output", output)
	if code != exitFindings {
		t.Fatalf("Execute returned %d, want %d; stdout=%s stderr=%s", code, exitFindings, stdout, stderr)
	}
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty when --output is used", stdout)
	}

	content, err := os.ReadFile(output)
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}
	var result model.RunResult
	if err := json.Unmarshal(content, &result); err != nil {
		t.Fatalf("report is not valid JSON: %v\n%s", err, string(content))
	}
	if len(result.Findings) != 1 {
		t.Fatalf("findings length = %d, want 1", len(result.Findings))
	}
	if result.Findings[0].ID != "DIF001" {
		t.Fatalf("finding ID = %q, want DIF001", result.Findings[0].ID)
	}
	if result.Summary.ExitCode != exitFindings {
		t.Fatalf("summary.exit_code = %d, want %d", result.Summary.ExitCode, exitFindings)
	}
}

func TestDiffFailOnCriticalReturnsSuccessForDIF001Error(t *testing.T) {
	base, current := writeDiffMappingFiles(t, `{"properties":{"status":{"type":"keyword"}}}`, `{"properties":{"status":{"type":"long"}}}`)

	code, stdout, stderr := executeForTest("diff", "--base", base, "--current", current, "--fail-on", "critical", "--format", "json")
	if code != exitSuccess {
		t.Fatalf("Execute returned %d, want %d; stdout=%s stderr=%s", code, exitSuccess, stdout, stderr)
	}
	var result model.RunResult
	if err := json.Unmarshal([]byte(stdout), &result); err != nil {
		t.Fatalf("stdout is not valid JSON: %v\n%s", err, stdout)
	}
	if result.Summary.ExitCode != exitSuccess {
		t.Fatalf("summary.exit_code = %d, want %d", result.Summary.ExitCode, exitSuccess)
	}
	if len(result.Findings) != 1 || result.Findings[0].ID != "DIF001" {
		t.Fatalf("expected DIF001 finding despite non-failing threshold, got %#v", result.Findings)
	}
}

func TestDiffMissingBaseReturnsUsageError(t *testing.T) {
	current := writeTempFile(t, "current.json", `{"properties":{"status":{"type":"long"}}}`)

	code, _, stderr := executeForTest("diff", "--current", current)
	if code != exitUsage {
		t.Fatalf("Execute returned %d, want %d", code, exitUsage)
	}
	if !strings.Contains(stderr, "requires --base") {
		t.Fatalf("stderr %q does not explain missing base", stderr)
	}
}

func TestDiffMissingCurrentReturnsUsageError(t *testing.T) {
	base := writeTempFile(t, "base.json", `{"properties":{"status":{"type":"keyword"}}}`)

	code, _, stderr := executeForTest("diff", "--base", base)
	if code != exitUsage {
		t.Fatalf("Execute returned %d, want %d", code, exitUsage)
	}
	if !strings.Contains(stderr, "requires --current") {
		t.Fatalf("stderr %q does not explain missing current", stderr)
	}
}

func TestDiffInvalidFormatReturnsUsageError(t *testing.T) {
	base, current := writeDiffMappingFiles(t, `{"properties":{"status":{"type":"keyword"}}}`, `{"properties":{"status":{"type":"long"}}}`)

	code, _, stderr := executeForTest("diff", "--base", base, "--current", current, "--format", "xml")
	if code != exitUsage {
		t.Fatalf("Execute returned %d, want %d", code, exitUsage)
	}
	if !strings.Contains(stderr, "invalid --format") {
		t.Fatalf("stderr %q does not explain invalid format", stderr)
	}
}

func TestDiffInvalidFailOnReturnsUsageError(t *testing.T) {
	base, current := writeDiffMappingFiles(t, `{"properties":{"status":{"type":"keyword"}}}`, `{"properties":{"status":{"type":"long"}}}`)

	code, _, stderr := executeForTest("diff", "--base", base, "--current", current, "--fail-on", "banana")
	if code != exitUsage {
		t.Fatalf("Execute returned %d, want %d", code, exitUsage)
	}
	if !strings.Contains(stderr, "invalid --fail-on") {
		t.Fatalf("stderr %q does not explain invalid fail-on", stderr)
	}
}

func TestDiffInvalidBaseJSONShortCircuitsDiff(t *testing.T) {
	base, current := writeDiffMappingFiles(t, `{"properties":`, `{"properties":{"status":{"type":"long"}}}`)

	code, stdout, stderr := executeForTest("diff", "--base", base, "--current", current)
	if code != exitInput {
		t.Fatalf("Execute returned %d, want %d; stdout=%s stderr=%s", code, exitInput, stdout, stderr)
	}
	if !strings.Contains(stdout, "base.json") {
		t.Fatalf("stdout %q does not contain base file name", stdout)
	}
	if strings.Contains(stdout, "DIF001") {
		t.Fatalf("stdout contains diff finding despite parse error: %s", stdout)
	}
}

func TestDiffInvalidCurrentJSONShortCircuitsDiff(t *testing.T) {
	base, current := writeDiffMappingFiles(t, `{"properties":{"status":{"type":"keyword"}}}`, `{"properties":`)

	code, stdout, stderr := executeForTest("diff", "--base", base, "--current", current)
	if code != exitInput {
		t.Fatalf("Execute returned %d, want %d; stdout=%s stderr=%s", code, exitInput, stdout, stderr)
	}
	if !strings.Contains(stdout, "current.json") {
		t.Fatalf("stdout %q does not contain current file name", stdout)
	}
	if strings.Contains(stdout, "DIF001") {
		t.Fatalf("stdout contains diff finding despite parse error: %s", stdout)
	}
}

func TestDiffBothInvalidJSONReportsBothFiles(t *testing.T) {
	base, current := writeDiffMappingFiles(t, `{"properties":`, `{"properties":`)

	code, stdout, stderr := executeForTest("diff", "--base", base, "--current", current)
	if code != exitInput {
		t.Fatalf("Execute returned %d, want %d; stdout=%s stderr=%s", code, exitInput, stdout, stderr)
	}
	if !strings.Contains(stdout, "base.json") {
		t.Fatalf("stdout %q does not contain base file name", stdout)
	}
	if !strings.Contains(stdout, "current.json") {
		t.Fatalf("stdout %q does not contain current file name", stdout)
	}
	if strings.Contains(stdout, "DIF001") {
		t.Fatalf("stdout contains diff finding despite parse error: %s", stdout)
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

func writeDiffMappingFiles(t *testing.T, baseContent string, currentContent string) (string, string) {
	t.Helper()
	root := t.TempDir()
	writeFileAt(t, root, "base.json", baseContent)
	writeFileAt(t, root, "current.json", currentContent)
	return filepath.Join(root, "base.json"), filepath.Join(root, "current.json")
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

type expectedFindingFixture struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	Severity        string   `json:"severity"`
	Category        string   `json:"category"`
	File            string   `json:"file"`
	JSONPointer     string   `json:"json_pointer"`
	MessageContains []string `json:"message_contains"`
}

func readExpectedFinding(t *testing.T, path string) expectedFindingFixture {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}
	var expected expectedFindingFixture
	if err := json.Unmarshal(content, &expected); err != nil {
		t.Fatalf("expected finding fixture is invalid JSON: %v", err)
	}
	return expected
}

func assertFindingMatchesExpected(t *testing.T, finding model.Finding, expected expectedFindingFixture) {
	t.Helper()
	if finding.ID != expected.ID {
		t.Fatalf("finding ID = %q, want %q", finding.ID, expected.ID)
	}
	if finding.Name != expected.Name {
		t.Fatalf("finding name = %q, want %q", finding.Name, expected.Name)
	}
	if string(finding.Severity) != expected.Severity {
		t.Fatalf("finding severity = %q, want %q", finding.Severity, expected.Severity)
	}
	if finding.Category != expected.Category {
		t.Fatalf("finding category = %q, want %q", finding.Category, expected.Category)
	}
	if finding.File != expected.File {
		t.Fatalf("finding file = %q, want %q", finding.File, expected.File)
	}
	if finding.JSONPointer != expected.JSONPointer {
		t.Fatalf("finding JSON pointer = %q, want %q", finding.JSONPointer, expected.JSONPointer)
	}
	for _, text := range expected.MessageContains {
		if !strings.Contains(finding.Message, text) {
			t.Fatalf("finding message %q does not contain %q", finding.Message, text)
		}
	}
}

func fixturePath(parts ...string) string {
	all := append([]string{"..", "..", "fixtures"}, parts...)
	return filepath.Join(all...)
}

func mappingJSONWithFields(count int) string {
	var builder strings.Builder
	builder.WriteString(`{"properties":{`)
	for i := 0; i < count; i++ {
		if i > 0 {
			builder.WriteByte(',')
		}
		fmt.Fprintf(&builder, "%q:{\"type\":\"keyword\"}", fmt.Sprintf("field_%04d", i))
	}
	builder.WriteString(`}}`)
	return builder.String()
}

func wrappedMappingJSONWithFields(count int) string {
	var builder strings.Builder
	builder.WriteString(`{"mappings":{"properties":{`)
	for i := 0; i < count; i++ {
		if i > 0 {
			builder.WriteByte(',')
		}
		fmt.Fprintf(&builder, "%q:{\"type\":\"keyword\"}", fmt.Sprintf("field_%04d", i))
	}
	builder.WriteString(`}}}`)
	return builder.String()
}

func mappingJSONWithRootDynamic(enabled bool) string {
	return fmt.Sprintf(`{"dynamic":%t,"properties":{"status":{"type":"keyword"}}}`, enabled)
}

func mappingJSONWithDynamicTemplate(includeMatchMappingType bool) string {
	matchMappingType := ""
	if includeMatchMappingType {
		matchMappingType = `"match_mapping_type":"string",`
	}
	return fmt.Sprintf(`{"dynamic_templates":[{"strings_as_keywords":{%s"mapping":{"type":"keyword"}}}]}`, matchMappingType)
}
