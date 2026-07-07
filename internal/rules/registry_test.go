package rules

import (
	"testing"

	"github.com/marcinbohm/search-index-preflight/internal/model"
)

type testRule struct {
	metadata Metadata
	findings []model.Finding
	err      error
	onCheck  func()
}

func (r testRule) Metadata() Metadata {
	return r.metadata
}

func (r testRule) Check(ctx Context, corpus model.Corpus) ([]model.Finding, error) {
	if r.onCheck != nil {
		r.onCheck()
	}
	return r.findings, r.err
}

func TestRegistryRejectsDuplicateRuleID(t *testing.T) {
	_, err := NewRegistry(
		testRule{metadata: validMetadata("SIL001")},
		testRule{metadata: validMetadata("SIL001")},
	)
	if err == nil {
		t.Fatal("NewRegistry returned nil error for duplicate rule ID")
	}
}

func TestRegistryRejectsNilRule(t *testing.T) {
	if _, err := NewRegistry(nil); err == nil {
		t.Fatal("NewRegistry returned nil error for nil rule")
	}
}

func TestRegistryRejectsIncompleteMetadata(t *testing.T) {
	tests := []struct {
		name     string
		metadata Metadata
	}{
		{name: "empty ID", metadata: Metadata{}},
		{name: "empty name", metadata: Metadata{ID: "SIL001", Category: "metadata", Severity: model.SeverityInfo, Confidence: model.ConfidenceHigh, Determinism: model.DeterminismDeterministic}},
		{name: "empty category", metadata: Metadata{ID: "SIL001", Name: "test-rule", Severity: model.SeverityInfo, Confidence: model.ConfidenceHigh, Determinism: model.DeterminismDeterministic}},
		{name: "empty severity", metadata: Metadata{ID: "SIL001", Name: "test-rule", Category: "metadata", Confidence: model.ConfidenceHigh, Determinism: model.DeterminismDeterministic}},
		{name: "empty confidence", metadata: Metadata{ID: "SIL001", Name: "test-rule", Category: "metadata", Severity: model.SeverityInfo, Determinism: model.DeterminismDeterministic}},
		{name: "empty determinism", metadata: Metadata{ID: "SIL001", Name: "test-rule", Category: "metadata", Severity: model.SeverityInfo, Confidence: model.ConfidenceHigh}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := NewRegistry(testRule{metadata: tt.metadata}); err == nil {
				t.Fatal("NewRegistry returned nil error for incomplete metadata")
			}
		})
	}
}

func TestRegistryListSortedByID(t *testing.T) {
	registry, err := NewRegistry(
		testRule{metadata: validMetadata("SIL010")},
		testRule{metadata: validMetadata("SIL002")},
		testRule{metadata: validMetadata("SIL001")},
	)
	if err != nil {
		t.Fatalf("NewRegistry returned error: %v", err)
	}

	gotRules := registry.List()
	got := make([]string, 0, len(gotRules))
	for _, rule := range gotRules {
		got = append(got, rule.Metadata().ID)
	}

	want := []string{"SIL001", "SIL002", "SIL010"}
	if len(got) != len(want) {
		t.Fatalf("List returned %d rules, want %d: %v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("List returned IDs %v, want %v", got, want)
		}
	}
}

func validMetadata(id string) Metadata {
	return Metadata{
		ID:          id,
		Name:        "test-rule-" + id,
		Category:    "metadata",
		Severity:    model.SeverityInfo,
		Confidence:  model.ConfidenceHigh,
		Determinism: model.DeterminismDeterministic,
	}
}
