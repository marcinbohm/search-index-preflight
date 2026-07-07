package rules

import (
	"testing"

	"github.com/marcinbohm/search-index-preflight/internal/model"
)

func TestBuiltinRegistryContainsOnlySIL001SIL002AndSIL003(t *testing.T) {
	registry, err := BuiltinRegistry()
	if err != nil {
		t.Fatalf("BuiltinRegistry returned error: %v", err)
	}

	rules := registry.List()
	if len(rules) != 3 {
		t.Fatalf("BuiltinRegistry returned %d rules, want 3", len(rules))
	}
	if rules[0].Metadata().ID != "SIL001" {
		t.Fatalf("built-in rule ID = %q, want SIL001", rules[0].Metadata().ID)
	}
	if rules[1].Metadata().ID != "SIL002" {
		t.Fatalf("built-in rule ID = %q, want SIL002", rules[1].Metadata().ID)
	}
	if rules[2].Metadata().ID != "SIL003" {
		t.Fatalf("built-in rule ID = %q, want SIL003", rules[2].Metadata().ID)
	}
}

func TestBuiltinRegistryRunExecutesSIL001(t *testing.T) {
	registry, err := BuiltinRegistry()
	if err != nil {
		t.Fatalf("BuiltinRegistry returned error: %v", err)
	}

	result, err := Run(Context{}, registry, RunRequest{
		Corpus: model.Corpus{
			Mappings: []model.Mapping{
				{
					Source:     testSource("mapping.json"),
					Properties: testFields(1000),
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if len(result.Findings) != 1 {
		t.Fatalf("Run returned %d findings, want 1", len(result.Findings))
	}
	if result.Findings[0].ID != "SIL001" {
		t.Fatalf("finding ID = %q, want SIL001", result.Findings[0].ID)
	}
}
