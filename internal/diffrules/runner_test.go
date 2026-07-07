package diffrules

import (
	"fmt"
	"strings"
	"testing"

	"github.com/marcinbohm/search-index-preflight/internal/diff"
	"github.com/marcinbohm/search-index-preflight/internal/model"
)

func TestBuiltinRegistryContainsDIF001AndDIF002(t *testing.T) {
	registry, err := BuiltinRegistry()
	if err != nil {
		t.Fatalf("BuiltinRegistry returned error: %v", err)
	}

	rules := registry.List()
	if len(rules) != 2 {
		t.Fatalf("expected two built-in diff rules, got %d", len(rules))
	}
	ids := []string{rules[0].Metadata().ID, rules[1].Metadata().ID}
	want := []string{"DIF001", "DIF002"}
	for i := range want {
		if ids[i] != want[i] {
			t.Fatalf("rule IDs = %#v, want %#v", ids, want)
		}
	}
}

func TestRegistryRejectsDuplicateIDs(t *testing.T) {
	_, err := NewRegistry(fakeRule{id: "DIF001"}, fakeRule{id: "DIF001"})
	if err == nil {
		t.Fatal("expected duplicate ID error")
	}
}

func TestRegistryRejectsNilRule(t *testing.T) {
	_, err := NewRegistry(nil)
	if err == nil {
		t.Fatal("expected nil rule error")
	}
	if !strings.Contains(err.Error(), "diff rule is nil") {
		t.Fatalf("expected nil diff rule message, got %q", err.Error())
	}
}

func TestRegistryRejectsIncompleteMetadata(t *testing.T) {
	valid := validFakeMetadata()
	tests := []struct {
		name     string
		metadata Metadata
	}{
		{
			name: "empty ID",
			metadata: func() Metadata {
				metadata := valid
				metadata.ID = ""
				return metadata
			}(),
		},
		{
			name: "empty name",
			metadata: func() Metadata {
				metadata := valid
				metadata.Name = ""
				return metadata
			}(),
		},
		{
			name: "empty category",
			metadata: func() Metadata {
				metadata := valid
				metadata.Category = ""
				return metadata
			}(),
		},
		{
			name: "empty severity",
			metadata: func() Metadata {
				metadata := valid
				metadata.Severity = ""
				return metadata
			}(),
		},
		{
			name: "empty confidence",
			metadata: func() Metadata {
				metadata := valid
				metadata.Confidence = ""
				return metadata
			}(),
		},
		{
			name: "empty determinism",
			metadata: func() Metadata {
				metadata := valid
				metadata.Determinism = ""
				return metadata
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewRegistry(fakeRule{metadata: tt.metadata})
			if err == nil {
				t.Fatal("expected metadata validation error")
			}
		})
	}
}

func TestRunRejectsNilRegistry(t *testing.T) {
	_, err := Run(Context{}, nil, RunRequest{})
	if err == nil {
		t.Fatal("expected nil registry error")
	}
	if !strings.Contains(err.Error(), "diff rule registry is nil") {
		t.Fatalf("expected nil registry message, got %q", err.Error())
	}
}

func TestRunExecutesBuiltinDiffRules(t *testing.T) {
	registry, err := BuiltinRegistry()
	if err != nil {
		t.Fatalf("BuiltinRegistry returned error: %v", err)
	}

	result, err := Run(Context{}, registry, RunRequest{Result: diff.Result{
		FieldChanges: []diff.FieldChange{
			fieldTypeChanged("status", model.FieldRoleProperty, "keyword", "long", "/properties/status", "/properties/status"),
			fieldRemoved("legacy_id", model.FieldRoleProperty, "keyword", "/properties/legacy_id"),
		},
	}})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if len(result.Findings) != 2 {
		t.Fatalf("expected two findings, got %#v", result.Findings)
	}
	if result.Findings[0].ID != "DIF001" {
		t.Fatalf("expected DIF001 finding, got %q", result.Findings[0].ID)
	}
	if result.Findings[1].ID != "DIF002" {
		t.Fatalf("expected DIF002 finding, got %q", result.Findings[1].ID)
	}
}

func TestRunGroupsFindingsByRuleOrder(t *testing.T) {
	registry, err := BuiltinRegistry()
	if err != nil {
		t.Fatalf("BuiltinRegistry returned error: %v", err)
	}

	result, err := Run(Context{}, registry, RunRequest{Result: diff.Result{
		FieldChanges: []diff.FieldChange{
			fieldRemoved("legacy_id", model.FieldRoleProperty, "keyword", "/properties/legacy_id"),
			fieldTypeChanged("status", model.FieldRoleProperty, "keyword", "long", "/properties/status", "/properties/status"),
		},
	}})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if len(result.Findings) != 2 {
		t.Fatalf("expected two findings, got %#v", result.Findings)
	}
	if result.Findings[0].ID != "DIF001" || result.Findings[1].ID != "DIF002" {
		t.Fatalf("expected grouped-by-rule order DIF001 then DIF002, got %#v", result.Findings)
	}
}

func TestDiffCompareToDIF001Integration(t *testing.T) {
	base := corpusWithMapping(property("status", "keyword"))
	current := corpusWithMapping(property("status", "long"))

	diffResult, err := diff.Compare(base, current)
	if err != nil {
		t.Fatalf("diff.Compare returned error: %v", err)
	}
	registry, err := BuiltinRegistry()
	if err != nil {
		t.Fatalf("BuiltinRegistry returned error: %v", err)
	}

	runResult, err := Run(Context{}, registry, RunRequest{Result: diffResult})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if len(runResult.Findings) != 1 {
		t.Fatalf("expected one finding, got %#v", runResult.Findings)
	}
	if runResult.Findings[0].ID != "DIF001" {
		t.Fatalf("expected DIF001, got %q", runResult.Findings[0].ID)
	}
}

func TestDiffCompareToDIF002Integration(t *testing.T) {
	base := corpusWithMapping(property("status", "keyword"), property("legacy_id", "keyword"))
	current := corpusWithMapping(property("status", "keyword"))

	diffResult, err := diff.Compare(base, current)
	if err != nil {
		t.Fatalf("diff.Compare returned error: %v", err)
	}
	registry, err := BuiltinRegistry()
	if err != nil {
		t.Fatalf("BuiltinRegistry returned error: %v", err)
	}

	runResult, err := Run(Context{}, registry, RunRequest{Result: diffResult})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if len(runResult.Findings) != 1 {
		t.Fatalf("expected one finding, got %#v", runResult.Findings)
	}
	if runResult.Findings[0].ID != "DIF002" {
		t.Fatalf("expected DIF002, got %q", runResult.Findings[0].ID)
	}
}

type fakeRule struct {
	id       string
	metadata Metadata
}

func (r fakeRule) Metadata() Metadata {
	if r.metadata.ID != "" ||
		r.metadata.Name != "" ||
		r.metadata.Category != "" ||
		r.metadata.Severity != "" ||
		r.metadata.Confidence != "" ||
		r.metadata.Determinism != "" {
		return r.metadata
	}
	metadata := validFakeMetadata()
	if r.id != "" {
		metadata.ID = r.id
	}
	return metadata
}

func validFakeMetadata() Metadata {
	return Metadata{
		ID:          "FAKE001",
		Name:        "fake",
		Category:    "test",
		Description: "fake test rule",
		Severity:    model.SeverityWarning,
		Confidence:  model.ConfidenceHigh,
		Determinism: model.DeterminismDeterministic,
	}
}

func (r fakeRule) Check(ctx Context, result diff.Result) ([]model.Finding, error) {
	return nil, nil
}

func corpusWithMapping(fields ...model.Field) model.Corpus {
	return model.Corpus{
		Mappings: []model.Mapping{
			{
				Source:      model.Source{Path: "mapping.json", RelativePath: "mapping.json"},
				Properties:  fields,
				JSONPointer: "",
			},
		},
	}
}

func property(path string, typ string) model.Field {
	return model.Field{
		Name:        path,
		Path:        path,
		Type:        typ,
		Source:      model.Source{Path: "mapping.json", RelativePath: "mapping.json"},
		JSONPointer: fmt.Sprintf("/properties/%s", path),
	}
}
