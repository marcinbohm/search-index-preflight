package rules

import (
	"testing"

	"github.com/marcinbohm/search-index-lint/internal/model"
)

func TestSIL002StandaloneMappingDynamicTrue(t *testing.T) {
	findings := runSIL002(t, model.Mapping{
		Source:  testSource("mapping.json"),
		Dynamic: model.DynamicSettingTrue,
	})

	requireSIL002Finding(t, findings, "mapping.json", "/dynamic")
}

func TestSIL002WrappedMappingDynamicTrue(t *testing.T) {
	findings := runSIL002(t, model.Mapping{
		Source:      testSource("wrapped.json"),
		JSONPointer: "/mappings",
		Dynamic:     model.DynamicSettingTrue,
	})

	requireSIL002Finding(t, findings, "wrapped.json", "/mappings/dynamic")
}

func TestSIL002StandaloneMappingDynamicFalse(t *testing.T) {
	findings := runSIL002(t, model.Mapping{
		Source:  testSource("mapping.json"),
		Dynamic: model.DynamicSettingFalse,
	})
	if len(findings) != 0 {
		t.Fatalf("SIL002 returned findings %#v, want none", findings)
	}
}

func TestSIL002StandaloneMappingDynamicStrict(t *testing.T) {
	findings := runSIL002(t, model.Mapping{
		Source:  testSource("mapping.json"),
		Dynamic: model.DynamicSettingStrict,
	})
	if len(findings) != 0 {
		t.Fatalf("SIL002 returned findings %#v, want none", findings)
	}
}

func TestSIL002StandaloneMappingDynamicUnspecified(t *testing.T) {
	findings := runSIL002(t, model.Mapping{
		Source:  testSource("mapping.json"),
		Dynamic: model.DynamicSettingUnspecified,
	})
	if len(findings) != 0 {
		t.Fatalf("SIL002 returned findings %#v, want none", findings)
	}
}

func TestSIL002IndexTemplateMappingDynamicTrue(t *testing.T) {
	mapping := model.Mapping{
		Source:  testSource("index-template.json"),
		Dynamic: model.DynamicSettingTrue,
	}

	findings, err := NewSIL002().Check(Context{}, model.Corpus{
		IndexTemplates: []model.IndexTemplate{
			{
				Source: testSource("index-template.json"),
				Template: model.TemplateBody{
					Mappings: &mapping,
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("SIL002 returned error: %v", err)
	}

	requireSIL002Finding(t, findings, "index-template.json", "/template/mappings/dynamic")
}

func TestSIL002ComponentTemplateMappingDynamicTrue(t *testing.T) {
	mapping := model.Mapping{
		Source:  testSource("component-template.json"),
		Dynamic: model.DynamicSettingTrue,
	}

	findings, err := NewSIL002().Check(Context{}, model.Corpus{
		ComponentTemplates: []model.ComponentTemplate{
			{
				Source: testSource("component-template.json"),
				Template: model.TemplateBody{
					Mappings: &mapping,
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("SIL002 returned error: %v", err)
	}

	requireSIL002Finding(t, findings, "component-template.json", "/template/mappings/dynamic")
}

func TestSIL002ChildObjectDynamicTrueDoesNotTrigger(t *testing.T) {
	findings := runSIL002(t, model.Mapping{
		Source:  testSource("mapping.json"),
		Dynamic: model.DynamicSettingUnspecified,
		Properties: []model.Field{
			{
				Name:    "metadata",
				Path:    "metadata",
				Dynamic: model.DynamicSettingTrue,
				Type:    "object",
			},
		},
	})
	if len(findings) != 0 {
		t.Fatalf("SIL002 returned findings %#v, want none", findings)
	}
}

func runSIL002(t *testing.T, mapping model.Mapping) []model.Finding {
	t.Helper()
	findings, err := NewSIL002().Check(Context{}, model.Corpus{Mappings: []model.Mapping{mapping}})
	if err != nil {
		t.Fatalf("SIL002 returned error: %v", err)
	}
	return findings
}

func requireSIL002Finding(t *testing.T, findings []model.Finding, file string, pointer string) {
	t.Helper()
	if len(findings) != 1 {
		t.Fatalf("SIL002 returned %d findings, want 1: %#v", len(findings), findings)
	}
	finding := findings[0]
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
	if finding.File != file {
		t.Fatalf("finding file = %q, want %q", finding.File, file)
	}
	if finding.JSONPointer != pointer {
		t.Fatalf("finding JSON pointer = %q, want %q", finding.JSONPointer, pointer)
	}
}
