package rules

import (
	"errors"
	"reflect"
	"testing"

	"github.com/marcinbohm/search-index-preflight/internal/model"
)

func TestRunExecutesRulesInSortedIDOrder(t *testing.T) {
	var order []string
	registry, err := NewRegistry(
		testRule{metadata: validMetadata("SIL010"), onCheck: func() { order = append(order, "SIL010") }},
		testRule{metadata: validMetadata("SIL002"), onCheck: func() { order = append(order, "SIL002") }},
		testRule{metadata: validMetadata("SIL001"), onCheck: func() { order = append(order, "SIL001") }},
	)
	if err != nil {
		t.Fatalf("NewRegistry returned error: %v", err)
	}

	if _, err := Run(Context{}, registry, RunRequest{Corpus: model.Corpus{}}); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	want := []string{"SIL001", "SIL002", "SIL010"}
	if !reflect.DeepEqual(order, want) {
		t.Fatalf("execution order = %v, want %v", order, want)
	}
}

func TestRunAggregatesFindings(t *testing.T) {
	registry, err := NewRegistry(
		testRule{
			metadata: validMetadata("SIL001"),
			findings: []model.Finding{
				{ID: "SIL001", Message: "first"},
			},
		},
		testRule{
			metadata: validMetadata("SIL002"),
			findings: []model.Finding{
				{ID: "SIL002", Message: "second"},
			},
		},
	)
	if err != nil {
		t.Fatalf("NewRegistry returned error: %v", err)
	}

	result, err := Run(Context{}, registry, RunRequest{Corpus: model.Corpus{}})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if len(result.Findings) != 2 {
		t.Fatalf("Run returned %d findings, want 2", len(result.Findings))
	}
	if result.Findings[0].ID != "SIL001" || result.Findings[1].ID != "SIL002" {
		t.Fatalf("Findings = %#v, want SIL001 then SIL002", result.Findings)
	}
}

func TestRunReturnsRuleError(t *testing.T) {
	registry, err := NewRegistry(
		testRule{metadata: validMetadata("SIL001"), err: errors.New("rule failed")},
	)
	if err != nil {
		t.Fatalf("NewRegistry returned error: %v", err)
	}

	if _, err := Run(Context{}, registry, RunRequest{Corpus: model.Corpus{}}); err == nil {
		t.Fatal("Run returned nil error for rule error")
	}
}

func TestRunRejectsNilRegistry(t *testing.T) {
	if _, err := Run(Context{}, nil, RunRequest{Corpus: model.Corpus{}}); err == nil {
		t.Fatal("Run returned nil error for nil registry")
	}
}
