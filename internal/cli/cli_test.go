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

func TestRulesListConsoleListsAllPublicRules(t *testing.T) {
	code, stdout, stderr := executeForTest("rules", "list")
	if code != exitSuccess {
		t.Fatalf("Execute returned %d, want %d; stderr=%s", code, exitSuccess, stderr)
	}
	for _, text := range []string{"SIL001", "SIL002", "SIL003", "DIF001", "DIF002", "DIF003", "lint", "diff", "error", "warning", "info"} {
		if !strings.Contains(stdout, text) {
			t.Fatalf("stdout %q does not contain %q", stdout, text)
		}
	}
}

func TestRulesListConsoleOrdersLintRulesBeforeDiffRules(t *testing.T) {
	code, stdout, stderr := executeForTest("rules", "list")
	if code != exitSuccess {
		t.Fatalf("Execute returned %d, want %d; stderr=%s", code, exitSuccess, stderr)
	}
	assertContainsInOrder(t, stdout, "SIL001", "SIL002", "SIL003", "DIF001", "DIF002", "DIF003")
}

func TestRulesListFamilyLint(t *testing.T) {
	code, stdout, stderr := executeForTest("rules", "list", "--family", "lint")
	if code != exitSuccess {
		t.Fatalf("Execute returned %d, want %d; stderr=%s", code, exitSuccess, stderr)
	}
	for _, text := range []string{"SIL001", "SIL002", "SIL003"} {
		if !strings.Contains(stdout, text) {
			t.Fatalf("stdout %q does not contain %q", stdout, text)
		}
	}
	for _, text := range []string{"DIF001", "DIF002", "DIF003"} {
		if strings.Contains(stdout, text) {
			t.Fatalf("stdout %q contains unexpected %q", stdout, text)
		}
	}
}

func TestRulesListFamilyLintFormatJSON(t *testing.T) {
	code, stdout, stderr := executeForTest("rules", "list", "--family", "lint", "--format", "json")
	if code != exitSuccess {
		t.Fatalf("Execute returned %d, want %d; stderr=%s", code, exitSuccess, stderr)
	}

	var output ruleListOutput
	if err := json.Unmarshal([]byte(stdout), &output); err != nil {
		t.Fatalf("stdout is not valid JSON: %v\n%s", err, stdout)
	}
	if len(output.Rules) != 3 {
		t.Fatalf("rules length = %d, want 3", len(output.Rules))
	}
	wantIDs := []string{"SIL001", "SIL002", "SIL003"}
	if got := ruleListIDs(output.Rules); fmt.Sprint(got) != fmt.Sprint(wantIDs) {
		t.Fatalf("rule IDs = %#v, want %#v", got, wantIDs)
	}
	for _, rule := range output.Rules {
		if rule.Family != "lint" {
			t.Fatalf("rule family = %q, want lint for %#v", rule.Family, rule)
		}
		if rule.ID == "" || rule.Name == "" || rule.Category == "" || rule.Severity == "" || rule.Confidence == "" || rule.Determinism == "" || rule.Description == "" {
			t.Fatalf("rule has empty metadata field: %#v", rule)
		}
		if strings.HasPrefix(rule.ID, "DIF") {
			t.Fatalf("lint family output contains diff rule: %#v", rule)
		}
	}
}

func TestRulesListFamilyDiff(t *testing.T) {
	code, stdout, stderr := executeForTest("rules", "list", "--family", "diff")
	if code != exitSuccess {
		t.Fatalf("Execute returned %d, want %d; stderr=%s", code, exitSuccess, stderr)
	}
	for _, text := range []string{"DIF001", "DIF002", "DIF003"} {
		if !strings.Contains(stdout, text) {
			t.Fatalf("stdout %q does not contain %q", stdout, text)
		}
	}
	for _, text := range []string{"SIL001", "SIL002", "SIL003"} {
		if strings.Contains(stdout, text) {
			t.Fatalf("stdout %q contains unexpected %q", stdout, text)
		}
	}
}

func TestRulesListFormatJSON(t *testing.T) {
	code, stdout, stderr := executeForTest("rules", "list", "--format", "json")
	if code != exitSuccess {
		t.Fatalf("Execute returned %d, want %d; stderr=%s", code, exitSuccess, stderr)
	}

	var output ruleListOutput
	if err := json.Unmarshal([]byte(stdout), &output); err != nil {
		t.Fatalf("stdout is not valid JSON: %v\n%s", err, stdout)
	}
	if len(output.Rules) != 6 {
		t.Fatalf("rules length = %d, want 6", len(output.Rules))
	}
	wantIDs := []string{"SIL001", "SIL002", "SIL003", "DIF001", "DIF002", "DIF003"}
	if got := ruleListIDs(output.Rules); fmt.Sprint(got) != fmt.Sprint(wantIDs) {
		t.Fatalf("rule IDs = %#v, want %#v", got, wantIDs)
	}
	families := map[string]bool{}
	for _, rule := range output.Rules {
		if rule.ID == "" || rule.Family == "" || rule.Name == "" || rule.Category == "" || rule.Severity == "" || rule.Confidence == "" || rule.Determinism == "" || rule.Description == "" {
			t.Fatalf("rule has empty metadata field: %#v", rule)
		}
		families[rule.Family] = true
	}
	if !families["lint"] || !families["diff"] {
		t.Fatalf("families = %#v, want lint and diff", families)
	}
}

func TestRulesListFamilyDiffFormatJSON(t *testing.T) {
	code, stdout, stderr := executeForTest("rules", "list", "--family", "diff", "--format", "json")
	if code != exitSuccess {
		t.Fatalf("Execute returned %d, want %d; stderr=%s", code, exitSuccess, stderr)
	}

	var output ruleListOutput
	if err := json.Unmarshal([]byte(stdout), &output); err != nil {
		t.Fatalf("stdout is not valid JSON: %v\n%s", err, stdout)
	}
	if len(output.Rules) != 3 {
		t.Fatalf("rules length = %d, want 3", len(output.Rules))
	}
	wantIDs := []string{"DIF001", "DIF002", "DIF003"}
	if got := ruleListIDs(output.Rules); fmt.Sprint(got) != fmt.Sprint(wantIDs) {
		t.Fatalf("rule IDs = %#v, want %#v", got, wantIDs)
	}
	for _, rule := range output.Rules {
		if rule.Family != "diff" {
			t.Fatalf("rule family = %q, want diff for %#v", rule.Family, rule)
		}
	}
}

func TestRulesListInvalidFormatReturnsUsageError(t *testing.T) {
	code, _, stderr := executeForTest("rules", "list", "--format", "yaml")
	if code != exitUsage {
		t.Fatalf("Execute returned %d, want %d", code, exitUsage)
	}
	if !strings.Contains(stderr, "invalid --format") {
		t.Fatalf("stderr %q does not explain invalid format", stderr)
	}
}

func TestRulesListInvalidFamilyReturnsUsageError(t *testing.T) {
	code, _, stderr := executeForTest("rules", "list", "--family", "all-rules")
	if code != exitUsage {
		t.Fatalf("Execute returned %d, want %d", code, exitUsage)
	}
	if !strings.Contains(stderr, "invalid --family") {
		t.Fatalf("stderr %q does not explain invalid family", stderr)
	}
}

func TestRulesListRejectsPositionalArgs(t *testing.T) {
	code, _, stderr := executeForTest("rules", "list", "extra")
	if code != exitUsage {
		t.Fatalf("Execute returned %d, want %d", code, exitUsage)
	}
	if !strings.Contains(stderr, "does not accept positional arguments") {
		t.Fatalf("stderr %q does not explain positional args", stderr)
	}
}

func TestRulesListHelpReturnsSuccess(t *testing.T) {
	code, stdout, stderr := executeForTest("rules", "list", "--help")
	if code != exitSuccess {
		t.Fatalf("Execute returned %d, want %d; stderr=%s", code, exitSuccess, stderr)
	}
	help := stdout + stderr
	for _, text := range []string{"--format", "--family"} {
		if !strings.Contains(help, text) {
			t.Fatalf("help output %q does not contain %q", help, text)
		}
	}
}

func TestExplainSIL001Console(t *testing.T) {
	code, stdout, stderr := executeForTest("explain", "SIL001")
	if code != exitSuccess {
		t.Fatalf("Execute returned %d, want %d; stderr=%s", code, exitSuccess, stderr)
	}
	for _, text := range []string{"SIL001", "lint", "total-fields-limit-risk", "mapping-limits", "warning", "high", "deterministic", "field-count", "total fields", "error findings when the limit is exceeded"} {
		if !strings.Contains(stdout, text) {
			t.Fatalf("stdout %q does not contain %q", stdout, text)
		}
	}
}

func TestExplainDIF001Console(t *testing.T) {
	code, stdout, stderr := executeForTest("explain", "DIF001")
	if code != exitSuccess {
		t.Fatalf("Execute returned %d, want %d; stderr=%s", code, exitSuccess, stderr)
	}
	for _, text := range []string{"DIF001", "diff", "field-type-changed", "schema-diff", "error", "high", "deterministic"} {
		if !strings.Contains(stdout, text) {
			t.Fatalf("stdout %q does not contain %q", stdout, text)
		}
	}
}

func TestExplainLowercaseRuleID(t *testing.T) {
	code, stdout, stderr := executeForTest("explain", "dif003")
	if code != exitSuccess {
		t.Fatalf("Execute returned %d, want %d; stderr=%s", code, exitSuccess, stderr)
	}
	for _, text := range []string{"DIF003", "field-added", "info"} {
		if !strings.Contains(stdout, text) {
			t.Fatalf("stdout %q does not contain %q", stdout, text)
		}
	}
}

func TestExplainSIL001FormatJSON(t *testing.T) {
	code, stdout, stderr := executeForTest("explain", "SIL001", "--format", "json")
	if code != exitSuccess {
		t.Fatalf("Execute returned %d, want %d; stderr=%s", code, exitSuccess, stderr)
	}

	var output ruleExplainOutput
	if err := json.Unmarshal([]byte(stdout), &output); err != nil {
		t.Fatalf("stdout is not valid JSON: %v\n%s", err, stdout)
	}
	if output.Rule.ID != "SIL001" {
		t.Fatalf("rule ID = %q, want SIL001", output.Rule.ID)
	}
	if output.Rule.Family != "lint" {
		t.Fatalf("rule family = %q, want lint", output.Rule.Family)
	}
	if output.Rule.Severity != "warning" {
		t.Fatalf("rule severity = %q, want warning", output.Rule.Severity)
	}
	if output.Rule.Description == "" {
		t.Fatal("rule description is empty")
	}
	if len(output.Rule.Notes) == 0 {
		t.Fatal("rule notes are empty, want SIL001 conditional severity note")
	}
	if !strings.Contains(output.Rule.Notes[0], "warning findings near the field-count limit") || !strings.Contains(output.Rule.Notes[0], "error findings when the limit is exceeded") {
		t.Fatalf("SIL001 notes = %#v, want conditional severity note", output.Rule.Notes)
	}
}

func TestExplainDIF003FormatJSON(t *testing.T) {
	code, stdout, stderr := executeForTest("explain", "DIF003", "--format", "json")
	if code != exitSuccess {
		t.Fatalf("Execute returned %d, want %d; stderr=%s", code, exitSuccess, stderr)
	}

	var output ruleExplainOutput
	if err := json.Unmarshal([]byte(stdout), &output); err != nil {
		t.Fatalf("stdout is not valid JSON: %v\n%s", err, stdout)
	}
	if output.Rule.ID != "DIF003" {
		t.Fatalf("rule ID = %q, want DIF003", output.Rule.ID)
	}
	if output.Rule.Family != "diff" {
		t.Fatalf("rule family = %q, want diff", output.Rule.Family)
	}
	if output.Rule.Name != "field-added" {
		t.Fatalf("rule name = %q, want field-added", output.Rule.Name)
	}
	if output.Rule.Severity != "info" {
		t.Fatalf("rule severity = %q, want info", output.Rule.Severity)
	}
}

func TestExplainMissingIDReturnsUsageError(t *testing.T) {
	code, _, stderr := executeForTest("explain")
	if code != exitUsage {
		t.Fatalf("Execute returned %d, want %d", code, exitUsage)
	}
	if !strings.Contains(stderr, "requires a rule ID") {
		t.Fatalf("stderr %q does not explain missing rule ID", stderr)
	}
}

func TestExplainUnknownIDReturnsUsageError(t *testing.T) {
	code, _, stderr := executeForTest("explain", "ABC999")
	if code != exitUsage {
		t.Fatalf("Execute returned %d, want %d", code, exitUsage)
	}
	if !strings.Contains(stderr, "unknown rule ID") {
		t.Fatalf("stderr %q does not explain unknown rule ID", stderr)
	}
}

func TestExplainTooManyArgsReturnsUsageError(t *testing.T) {
	code, _, stderr := executeForTest("explain", "SIL001", "DIF001")
	if code != exitUsage {
		t.Fatalf("Execute returned %d, want %d", code, exitUsage)
	}
	if !strings.Contains(stderr, "exactly one rule ID") {
		t.Fatalf("stderr %q does not explain too many args", stderr)
	}
}

func TestExplainInvalidFormatReturnsUsageError(t *testing.T) {
	code, _, stderr := executeForTest("explain", "SIL001", "--format", "yaml")
	if code != exitUsage {
		t.Fatalf("Execute returned %d, want %d", code, exitUsage)
	}
	if !strings.Contains(stderr, "invalid --format") {
		t.Fatalf("stderr %q does not explain invalid format", stderr)
	}
}

func TestExplainHelpReturnsSuccess(t *testing.T) {
	code, stdout, stderr := executeForTest("explain", "--help")
	if code != exitSuccess {
		t.Fatalf("Execute returned %d, want %d; stderr=%s", code, exitSuccess, stderr)
	}
	help := stdout + stderr
	for _, text := range []string{"--format", "<RULE_ID>"} {
		if !strings.Contains(help, text) {
			t.Fatalf("help output %q does not contain %q", help, text)
		}
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

func TestDiffDirectorySameRelativePathEmitsDIF003ForAddedField(t *testing.T) {
	root := t.TempDir()
	baseDir := filepath.Join(root, "old")
	currentDir := filepath.Join(root, "new")
	writeFileAt(t, baseDir, filepath.Join("schemas", "mapping.json"), `{"properties":{"status":{"type":"keyword"}}}`)
	writeFileAt(t, currentDir, filepath.Join("schemas", "mapping.json"), `{"properties":{"status":{"type":"keyword"},"customer_id":{"type":"keyword"}}}`)

	code, stdout, stderr := executeForTest("diff", "--base", baseDir, "--current", currentDir)
	if code != exitSuccess {
		t.Fatalf("Execute returned %d, want %d; stdout=%s stderr=%s", code, exitSuccess, stdout, stderr)
	}
	for _, text := range []string{"DIF003", filepath.Join("schemas", "mapping.json"), "customer_id"} {
		if !strings.Contains(stdout, text) {
			t.Fatalf("stdout %q does not contain %q", stdout, text)
		}
	}
	for _, unexpected := range []string{"DIF001", "DIF002"} {
		if strings.Contains(stdout, unexpected) {
			t.Fatalf("stdout contains unexpected %s for clean field addition: %s", unexpected, stdout)
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

func TestDiffDIF002FixtureEmitsWarningFinding(t *testing.T) {
	base := fixturePath("diff", "dif002-field-removed", "base")
	current := fixturePath("diff", "dif002-field-removed", "current")

	code, stdout, stderr := executeForTest("diff", "--base", base, "--current", current)
	if code != exitSuccess {
		t.Fatalf("Execute returned %d, want %d; stdout=%s stderr=%s", code, exitSuccess, stdout, stderr)
	}
	for _, text := range []string{"DIF002", "warning", "legacy_id"} {
		if !strings.Contains(stdout, text) {
			t.Fatalf("stdout %q does not contain %q", stdout, text)
		}
	}
}

func TestDiffDIF002FixtureJSONMatchesExpectedFindingFields(t *testing.T) {
	base := fixturePath("diff", "dif002-field-removed", "base")
	current := fixturePath("diff", "dif002-field-removed", "current")
	expected := readExpectedFinding(t, fixturePath("diff", "dif002-field-removed", "expected.finding.json"))

	code, stdout, stderr := executeForTest("diff", "--base", base, "--current", current, "--format", "json")
	if code != exitSuccess {
		t.Fatalf("Execute returned %d, want %d; stdout=%s stderr=%s", code, exitSuccess, stdout, stderr)
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
	assertFindingMatchesExpected(t, result.Findings[0], expected)
}

func TestDiffDIF003FixtureEmitsInfoFinding(t *testing.T) {
	base := fixturePath("diff", "dif003-field-added", "base")
	current := fixturePath("diff", "dif003-field-added", "current")

	code, stdout, stderr := executeForTest("diff", "--base", base, "--current", current)
	if code != exitSuccess {
		t.Fatalf("Execute returned %d, want %d; stdout=%s stderr=%s", code, exitSuccess, stdout, stderr)
	}
	for _, text := range []string{"DIF003", "info", "customer_id"} {
		if !strings.Contains(stdout, text) {
			t.Fatalf("stdout %q does not contain %q", stdout, text)
		}
	}
}

func TestDiffDIF003FixtureJSONMatchesExpectedFindingFields(t *testing.T) {
	base := fixturePath("diff", "dif003-field-added", "base")
	current := fixturePath("diff", "dif003-field-added", "current")
	expected := readExpectedFinding(t, fixturePath("diff", "dif003-field-added", "expected.finding.json"))

	code, stdout, stderr := executeForTest("diff", "--base", base, "--current", current, "--format", "json")
	if code != exitSuccess {
		t.Fatalf("Execute returned %d, want %d; stdout=%s stderr=%s", code, exitSuccess, stdout, stderr)
	}

	var result model.RunResult
	if err := json.Unmarshal([]byte(stdout), &result); err != nil {
		t.Fatalf("stdout is not valid JSON: %v\n%s", err, stdout)
	}
	if result.Summary.FindingsTotal != 1 {
		t.Fatalf("findings_total = %d, want 1", result.Summary.FindingsTotal)
	}
	if result.Summary.Info != 1 {
		t.Fatalf("summary.info = %d, want 1", result.Summary.Info)
	}
	if result.Summary.Warning != 0 {
		t.Fatalf("summary.warning = %d, want 0", result.Summary.Warning)
	}
	if result.Summary.Error != 0 {
		t.Fatalf("summary.error = %d, want 0", result.Summary.Error)
	}
	if result.Summary.ExitCode != exitSuccess {
		t.Fatalf("summary.exit_code = %d, want %d", result.Summary.ExitCode, exitSuccess)
	}
	if len(result.Findings) != 1 {
		t.Fatalf("findings length = %d, want 1", len(result.Findings))
	}
	assertFindingMatchesExpected(t, result.Findings[0], expected)
}

func TestDiffMixedFixtureEmitsAllDiffRules(t *testing.T) {
	base := fixturePath("diff", "mixed-field-changes", "base")
	current := fixturePath("diff", "mixed-field-changes", "current")

	code, stdout, stderr := executeForTest("diff", "--base", base, "--current", current)
	if code != exitFindings {
		t.Fatalf("Execute returned %d, want %d; stdout=%s stderr=%s", code, exitFindings, stdout, stderr)
	}
	for _, text := range []string{"DIF001", "DIF002", "DIF003", "error", "warning", "info", "status", "legacy_id", "customer_id"} {
		if !strings.Contains(stdout, text) {
			t.Fatalf("stdout %q does not contain %q", stdout, text)
		}
	}
}

func TestDiffMixedFixtureJSONMatchesExpectedFindingFields(t *testing.T) {
	base := fixturePath("diff", "mixed-field-changes", "base")
	current := fixturePath("diff", "mixed-field-changes", "current")
	expected := readExpectedFindings(t, fixturePath("diff", "mixed-field-changes", "expected.findings.json"))

	code, stdout, stderr := executeForTest("diff", "--base", base, "--current", current, "--format", "json")
	if code != exitFindings {
		t.Fatalf("Execute returned %d, want %d; stdout=%s stderr=%s", code, exitFindings, stdout, stderr)
	}

	var result model.RunResult
	if err := json.Unmarshal([]byte(stdout), &result); err != nil {
		t.Fatalf("stdout is not valid JSON: %v\n%s", err, stdout)
	}
	if result.Summary.FindingsTotal != 3 {
		t.Fatalf("findings_total = %d, want 3", result.Summary.FindingsTotal)
	}
	if result.Summary.Error != 1 {
		t.Fatalf("summary.error = %d, want 1", result.Summary.Error)
	}
	if result.Summary.Warning != 1 {
		t.Fatalf("summary.warning = %d, want 1", result.Summary.Warning)
	}
	if result.Summary.Info != 1 {
		t.Fatalf("summary.info = %d, want 1", result.Summary.Info)
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
	if len(result.Findings) != len(expected) {
		t.Fatalf("findings length = %d, want %d", len(result.Findings), len(expected))
	}
	for i := range expected {
		assertFindingMatchesExpected(t, result.Findings[i], expected[i])
	}
}

func TestDiffMixedFixtureFailOnCriticalReportsFindingsWithoutFailing(t *testing.T) {
	base := fixturePath("diff", "mixed-field-changes", "base")
	current := fixturePath("diff", "mixed-field-changes", "current")

	code, stdout, stderr := executeForTest("diff", "--base", base, "--current", current, "--fail-on", "critical")
	if code != exitSuccess {
		t.Fatalf("Execute returned %d, want %d; stdout=%s stderr=%s", code, exitSuccess, stdout, stderr)
	}
	for _, text := range []string{"DIF001", "DIF002", "DIF003"} {
		if !strings.Contains(stdout, text) {
			t.Fatalf("stdout %q does not contain %q", stdout, text)
		}
	}
}

func TestDiffMixedFixtureOutputWritesJSONReport(t *testing.T) {
	base := fixturePath("diff", "mixed-field-changes", "base")
	current := fixturePath("diff", "mixed-field-changes", "current")
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
	if result.Summary.Error != 1 {
		t.Fatalf("summary.error = %d, want 1", result.Summary.Error)
	}
	if result.Summary.Warning != 1 {
		t.Fatalf("summary.warning = %d, want 1", result.Summary.Warning)
	}
	if result.Summary.Info != 1 {
		t.Fatalf("summary.info = %d, want 1", result.Summary.Info)
	}
	if result.Summary.ExitCode != exitFindings {
		t.Fatalf("summary.exit_code = %d, want %d", result.Summary.ExitCode, exitFindings)
	}
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
	if !strings.Contains(stdout, "DIF002") {
		t.Fatalf("stdout %q does not report removed field warning for unmatched base resource", stdout)
	}
	if !strings.Contains(stdout, "DIF003") {
		t.Fatalf("stdout %q does not report added field info for unmatched current resource", stdout)
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

func TestDiffFieldRemovedReturnsWarningByDefault(t *testing.T) {
	base, current := writeDiffMappingFiles(t, `{"properties":{"status":{"type":"keyword"},"legacy_id":{"type":"keyword"}}}`, `{"properties":{"status":{"type":"keyword"}}}`)

	code, stdout, stderr := executeForTest("diff", "--base", base, "--current", current)
	if code != exitSuccess {
		t.Fatalf("Execute returned %d, want %d; stdout=%s stderr=%s", code, exitSuccess, stdout, stderr)
	}
	for _, text := range []string{"DIF002", "legacy_id", "warning"} {
		if !strings.Contains(stdout, text) {
			t.Fatalf("stdout %q does not contain %q", stdout, text)
		}
	}
}

func TestDiffFieldRemovedFailsWithFailOnWarning(t *testing.T) {
	base, current := writeDiffMappingFiles(t, `{"properties":{"status":{"type":"keyword"},"legacy_id":{"type":"keyword"}}}`, `{"properties":{"status":{"type":"keyword"}}}`)

	code, stdout, stderr := executeForTest("diff", "--base", base, "--current", current, "--fail-on", "warning")
	if code != exitFindings {
		t.Fatalf("Execute returned %d, want %d; stdout=%s stderr=%s", code, exitFindings, stdout, stderr)
	}
	if !strings.Contains(stdout, "DIF002") {
		t.Fatalf("stdout %q does not contain DIF002", stdout)
	}
}

func TestDiffFormatJSONWithDIF002Finding(t *testing.T) {
	base, current := writeDiffMappingFiles(t, `{"properties":{"status":{"type":"keyword"},"legacy_id":{"type":"keyword"}}}`, `{"properties":{"status":{"type":"keyword"}}}`)

	code, stdout, stderr := executeForTest("diff", "--base", base, "--current", current, "--format", "json")
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
	if result.Summary.Error != 0 {
		t.Fatalf("summary.error = %d, want 0", result.Summary.Error)
	}
	if result.Summary.ExitCode != exitSuccess {
		t.Fatalf("summary.exit_code = %d, want %d", result.Summary.ExitCode, exitSuccess)
	}
	if len(result.Findings) != 1 {
		t.Fatalf("findings length = %d, want 1", len(result.Findings))
	}
	finding := result.Findings[0]
	if finding.ID != "DIF002" {
		t.Fatalf("finding ID = %q, want DIF002", finding.ID)
	}
	if finding.Severity != model.SeverityWarning {
		t.Fatalf("finding severity = %q, want %q", finding.Severity, model.SeverityWarning)
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

func TestDiffFieldAddedReturnsInfoByDefault(t *testing.T) {
	base, current := writeDiffMappingFiles(t, `{"properties":{"status":{"type":"keyword"}}}`, `{"properties":{"status":{"type":"keyword"},"customer_id":{"type":"keyword"}}}`)

	code, stdout, stderr := executeForTest("diff", "--base", base, "--current", current)
	if code != exitSuccess {
		t.Fatalf("Execute returned %d, want %d; stdout=%s stderr=%s", code, exitSuccess, stdout, stderr)
	}
	for _, text := range []string{"DIF003", "customer_id", "info"} {
		if !strings.Contains(stdout, text) {
			t.Fatalf("stdout %q does not contain %q", stdout, text)
		}
	}
}

func TestDiffFieldAddedDoesNotFailWithFailOnWarning(t *testing.T) {
	base, current := writeDiffMappingFiles(t, `{"properties":{"status":{"type":"keyword"}}}`, `{"properties":{"status":{"type":"keyword"},"customer_id":{"type":"keyword"}}}`)

	code, stdout, stderr := executeForTest("diff", "--base", base, "--current", current, "--fail-on", "warning")
	if code != exitSuccess {
		t.Fatalf("Execute returned %d, want %d; stdout=%s stderr=%s", code, exitSuccess, stdout, stderr)
	}
	if !strings.Contains(stdout, "DIF003") {
		t.Fatalf("stdout %q does not contain DIF003", stdout)
	}
}

func TestDiffFieldAddedFailsWithFailOnInfo(t *testing.T) {
	base, current := writeDiffMappingFiles(t, `{"properties":{"status":{"type":"keyword"}}}`, `{"properties":{"status":{"type":"keyword"},"customer_id":{"type":"keyword"}}}`)

	code, stdout, stderr := executeForTest("diff", "--base", base, "--current", current, "--fail-on", "info")
	if code != exitFindings {
		t.Fatalf("Execute returned %d, want %d; stdout=%s stderr=%s", code, exitFindings, stdout, stderr)
	}
	if !strings.Contains(stdout, "DIF003") {
		t.Fatalf("stdout %q does not contain DIF003", stdout)
	}
}

func TestDiffFormatJSONWithDIF003Finding(t *testing.T) {
	base, current := writeDiffMappingFiles(t, `{"properties":{"status":{"type":"keyword"}}}`, `{"properties":{"status":{"type":"keyword"},"customer_id":{"type":"keyword"}}}`)

	code, stdout, stderr := executeForTest("diff", "--base", base, "--current", current, "--format", "json")
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
	if result.Summary.Info != 1 {
		t.Fatalf("summary.info = %d, want 1", result.Summary.Info)
	}
	if result.Summary.Warning != 0 {
		t.Fatalf("summary.warning = %d, want 0", result.Summary.Warning)
	}
	if result.Summary.Error != 0 {
		t.Fatalf("summary.error = %d, want 0", result.Summary.Error)
	}
	if result.Summary.ExitCode != exitSuccess {
		t.Fatalf("summary.exit_code = %d, want %d", result.Summary.ExitCode, exitSuccess)
	}
	if len(result.Findings) != 1 {
		t.Fatalf("findings length = %d, want 1", len(result.Findings))
	}
	finding := result.Findings[0]
	if finding.ID != "DIF003" {
		t.Fatalf("finding ID = %q, want DIF003", finding.ID)
	}
	if finding.Severity != model.SeverityInfo {
		t.Fatalf("finding severity = %q, want %q", finding.Severity, model.SeverityInfo)
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

func TestDiffFormatJSONWithAllDiffRules(t *testing.T) {
	base, current := writeDiffMappingFiles(t, `{"properties":{"status":{"type":"keyword"},"legacy_id":{"type":"keyword"}}}`, `{"properties":{"status":{"type":"long"},"customer_id":{"type":"keyword"}}}`)

	code, stdout, stderr := executeForTest("diff", "--base", base, "--current", current, "--format", "json")
	if code != exitFindings {
		t.Fatalf("Execute returned %d, want %d; stdout=%s stderr=%s", code, exitFindings, stdout, stderr)
	}

	var result model.RunResult
	if err := json.Unmarshal([]byte(stdout), &result); err != nil {
		t.Fatalf("stdout is not valid JSON: %v\n%s", err, stdout)
	}
	if result.Summary.FindingsTotal != 3 {
		t.Fatalf("findings_total = %d, want 3", result.Summary.FindingsTotal)
	}
	if result.Summary.Error != 1 {
		t.Fatalf("summary.error = %d, want 1", result.Summary.Error)
	}
	if result.Summary.Warning != 1 {
		t.Fatalf("summary.warning = %d, want 1", result.Summary.Warning)
	}
	if result.Summary.Info != 1 {
		t.Fatalf("summary.info = %d, want 1", result.Summary.Info)
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
	if len(result.Findings) != 3 {
		t.Fatalf("findings length = %d, want 3", len(result.Findings))
	}
	wantIDs := []string{"DIF001", "DIF002", "DIF003"}
	for i, want := range wantIDs {
		if result.Findings[i].ID != want {
			t.Fatalf("finding IDs = %#v, want order %#v", []string{result.Findings[0].ID, result.Findings[1].ID, result.Findings[2].ID}, wantIDs)
		}
	}
}

func TestDiffConsoleWithAllDiffRules(t *testing.T) {
	base, current := writeDiffMappingFiles(t, `{"properties":{"status":{"type":"keyword"},"legacy_id":{"type":"keyword"}}}`, `{"properties":{"status":{"type":"long"},"customer_id":{"type":"keyword"}}}`)

	code, stdout, stderr := executeForTest("diff", "--base", base, "--current", current)
	if code != exitFindings {
		t.Fatalf("Execute returned %d, want %d; stdout=%s stderr=%s", code, exitFindings, stdout, stderr)
	}
	for _, text := range []string{"DIF001", "DIF002", "DIF003", "error", "warning", "info", "status", "legacy_id", "customer_id"} {
		if !strings.Contains(stdout, text) {
			t.Fatalf("stdout %q does not contain %q", stdout, text)
		}
	}
}

func TestLintDoesNotEmitDiffRules(t *testing.T) {
	path := writeTempFile(t, "mapping.json", `{"properties":{"status":{"type":"keyword"}}}`)

	code, stdout, stderr := executeForTest("lint", "--mapping", path)
	if code != exitSuccess {
		t.Fatalf("Execute returned %d, want %d; stdout=%s stderr=%s", code, exitSuccess, stdout, stderr)
	}
	if strings.Contains(stdout, "DIF001") {
		t.Fatalf("stdout contains unexpected DIF001 from lint: %s", stdout)
	}
	if strings.Contains(stdout, "DIF002") {
		t.Fatalf("stdout contains unexpected DIF002 from lint: %s", stdout)
	}
	if strings.Contains(stdout, "DIF003") {
		t.Fatalf("stdout contains unexpected DIF003 from lint: %s", stdout)
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

func TestDiffOutputWritesWarningOnlyDIF002JSONReport(t *testing.T) {
	base := fixturePath("diff", "dif002-field-removed", "base")
	current := fixturePath("diff", "dif002-field-removed", "current")
	output := filepath.Join(t.TempDir(), "report.json")

	code, stdout, stderr := executeForTest("diff", "--base", base, "--current", current, "--format", "json", "--output", output)
	if code != exitSuccess {
		t.Fatalf("Execute returned %d, want %d; stdout=%s stderr=%s", code, exitSuccess, stdout, stderr)
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
	if result.Findings[0].ID != "DIF002" {
		t.Fatalf("finding ID = %q, want DIF002", result.Findings[0].ID)
	}
	if result.Findings[0].Severity != model.SeverityWarning {
		t.Fatalf("finding severity = %q, want %q", result.Findings[0].Severity, model.SeverityWarning)
	}
	if result.Summary.FindingsTotal != 1 {
		t.Fatalf("findings_total = %d, want 1", result.Summary.FindingsTotal)
	}
	if result.Summary.Warning != 1 {
		t.Fatalf("summary.warning = %d, want 1", result.Summary.Warning)
	}
	if result.Summary.Error != 0 {
		t.Fatalf("summary.error = %d, want 0", result.Summary.Error)
	}
	if result.Summary.ExitCode != exitSuccess {
		t.Fatalf("summary.exit_code = %d, want %d", result.Summary.ExitCode, exitSuccess)
	}
}

func TestDiffOutputWritesInfoOnlyDIF003JSONReport(t *testing.T) {
	base := fixturePath("diff", "dif003-field-added", "base")
	current := fixturePath("diff", "dif003-field-added", "current")
	output := filepath.Join(t.TempDir(), "report.json")

	code, stdout, stderr := executeForTest("diff", "--base", base, "--current", current, "--format", "json", "--output", output)
	if code != exitSuccess {
		t.Fatalf("Execute returned %d, want %d; stdout=%s stderr=%s", code, exitSuccess, stdout, stderr)
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
	if result.Findings[0].ID != "DIF003" {
		t.Fatalf("finding ID = %q, want DIF003", result.Findings[0].ID)
	}
	if result.Findings[0].Severity != model.SeverityInfo {
		t.Fatalf("finding severity = %q, want %q", result.Findings[0].Severity, model.SeverityInfo)
	}
	if result.Summary.FindingsTotal != 1 {
		t.Fatalf("findings_total = %d, want 1", result.Summary.FindingsTotal)
	}
	if result.Summary.Info != 1 {
		t.Fatalf("summary.info = %d, want 1", result.Summary.Info)
	}
	if result.Summary.Warning != 0 {
		t.Fatalf("summary.warning = %d, want 0", result.Summary.Warning)
	}
	if result.Summary.Error != 0 {
		t.Fatalf("summary.error = %d, want 0", result.Summary.Error)
	}
	if result.Summary.ExitCode != exitSuccess {
		t.Fatalf("summary.exit_code = %d, want %d", result.Summary.ExitCode, exitSuccess)
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
	if strings.Contains(stdout, "DIF002") {
		t.Fatalf("stdout contains diff finding despite parse error: %s", stdout)
	}
	if strings.Contains(stdout, "DIF003") {
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
	if strings.Contains(stdout, "DIF002") {
		t.Fatalf("stdout contains diff finding despite parse error: %s", stdout)
	}
	if strings.Contains(stdout, "DIF003") {
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
	if strings.Contains(stdout, "DIF002") {
		t.Fatalf("stdout contains diff finding despite parse error: %s", stdout)
	}
	if strings.Contains(stdout, "DIF003") {
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

type ruleListOutput struct {
	Rules []ruleListItemForTest `json:"rules"`
}

type ruleListItemForTest struct {
	ID          string `json:"id"`
	Family      string `json:"family"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Severity    string `json:"severity"`
	Confidence  string `json:"confidence"`
	Determinism string `json:"determinism"`
	Description string `json:"description"`
}

type ruleExplainOutput struct {
	Rule ruleExplainItemForTest `json:"rule"`
}

type ruleExplainItemForTest struct {
	ID          string   `json:"id"`
	Family      string   `json:"family"`
	Name        string   `json:"name"`
	Category    string   `json:"category"`
	Severity    string   `json:"severity"`
	Confidence  string   `json:"confidence"`
	Determinism string   `json:"determinism"`
	Description string   `json:"description"`
	Notes       []string `json:"notes"`
}

func ruleListIDs(rules []ruleListItemForTest) []string {
	ids := make([]string, 0, len(rules))
	for _, rule := range rules {
		ids = append(ids, rule.ID)
	}
	return ids
}

func assertContainsInOrder(t *testing.T, text string, tokens ...string) {
	t.Helper()
	offset := 0
	for _, token := range tokens {
		index := strings.Index(text[offset:], token)
		if index < 0 {
			t.Fatalf("text %q does not contain %q after offset %d", text, token, offset)
		}
		offset += index + len(token)
	}
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

func readExpectedFindings(t *testing.T, path string) []expectedFindingFixture {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}
	var expected struct {
		Findings []expectedFindingFixture `json:"findings"`
	}
	if err := json.Unmarshal(content, &expected); err != nil {
		t.Fatalf("expected findings fixture is invalid JSON: %v", err)
	}
	return expected.Findings
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
