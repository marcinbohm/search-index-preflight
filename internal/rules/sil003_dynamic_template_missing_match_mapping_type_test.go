package rules

import (
	"testing"

	"github.com/marcinbohm/search-index-preflight/internal/model"
)

func TestSIL003StandaloneMappingDynamicTemplateMissingMatchMappingType(t *testing.T) {
	findings := runSIL003(t, model.Corpus{
		Mappings: []model.Mapping{
			{
				Source: testSource("mapping.json"),
				DynamicTemplates: []model.DynamicTemplate{
					testDynamicTemplate("mapping.json", "strings_as_keywords", "/dynamic_templates/0/strings_as_keywords", false),
				},
			},
		},
	})

	requireSIL003Finding(t, findings, "mapping.json", "/dynamic_templates/0/strings_as_keywords")
}

func TestSIL003StandaloneMappingDynamicTemplateWithMatchMappingType(t *testing.T) {
	findings := runSIL003(t, model.Corpus{
		Mappings: []model.Mapping{
			{
				Source: testSource("mapping.json"),
				DynamicTemplates: []model.DynamicTemplate{
					testDynamicTemplate("mapping.json", "strings_as_keywords", "/dynamic_templates/0/strings_as_keywords", true),
				},
			},
		},
	})
	if len(findings) != 0 {
		t.Fatalf("SIL003 returned findings %#v, want none", findings)
	}
}

func TestSIL003WrappedMappingDynamicTemplateMissingMatchMappingType(t *testing.T) {
	findings := runSIL003(t, model.Corpus{
		Mappings: []model.Mapping{
			{
				Source:      testSource("wrapped.json"),
				JSONPointer: "/mappings",
				DynamicTemplates: []model.DynamicTemplate{
					testDynamicTemplate("wrapped.json", "strings_as_keywords", "/mappings/dynamic_templates/0/strings_as_keywords", false),
				},
			},
		},
	})

	requireSIL003Finding(t, findings, "wrapped.json", "/mappings/dynamic_templates/0/strings_as_keywords")
}

func TestSIL003IndexTemplateDynamicTemplateMissingMatchMappingType(t *testing.T) {
	mapping := model.Mapping{
		Source:      testSource("index-template.json"),
		JSONPointer: "/template/mappings",
		DynamicTemplates: []model.DynamicTemplate{
			testDynamicTemplate("index-template.json", "strings_as_keywords", "/template/mappings/dynamic_templates/0/strings_as_keywords", false),
		},
	}

	findings := runSIL003(t, model.Corpus{
		IndexTemplates: []model.IndexTemplate{
			{
				Source: testSource("index-template.json"),
				Template: model.TemplateBody{
					Mappings: &mapping,
				},
			},
		},
	})

	requireSIL003Finding(t, findings, "index-template.json", "/template/mappings/dynamic_templates/0/strings_as_keywords")
}

func TestSIL003ComponentTemplateDynamicTemplateMissingMatchMappingType(t *testing.T) {
	mapping := model.Mapping{
		Source:      testSource("component-template.json"),
		JSONPointer: "/template/mappings",
		DynamicTemplates: []model.DynamicTemplate{
			testDynamicTemplate("component-template.json", "strings_as_keywords", "/template/mappings/dynamic_templates/0/strings_as_keywords", false),
		},
	}

	findings := runSIL003(t, model.Corpus{
		ComponentTemplates: []model.ComponentTemplate{
			{
				Source: testSource("component-template.json"),
				Template: model.TemplateBody{
					Mappings: &mapping,
				},
			},
		},
	})

	requireSIL003Finding(t, findings, "component-template.json", "/template/mappings/dynamic_templates/0/strings_as_keywords")
}

func TestSIL003EscapedDynamicTemplateName(t *testing.T) {
	findings := runSIL003(t, model.Corpus{
		Mappings: []model.Mapping{
			{
				Source: testSource("mapping.json"),
				DynamicTemplates: []model.DynamicTemplate{
					testDynamicTemplate("mapping.json", "service/name~template", "/dynamic_templates/0/service~1name~0template", false),
				},
			},
		},
	})

	requireSIL003Finding(t, findings, "mapping.json", "/dynamic_templates/0/service~1name~0template")
}

func TestSIL003MultipleDynamicTemplatesOnlyOneMissing(t *testing.T) {
	findings := runSIL003(t, model.Corpus{
		Mappings: []model.Mapping{
			{
				Source: testSource("mapping.json"),
				DynamicTemplates: []model.DynamicTemplate{
					testDynamicTemplate("mapping.json", "with_type", "/dynamic_templates/0/with_type", true),
					testDynamicTemplate("mapping.json", "missing_type", "/dynamic_templates/1/missing_type", false),
				},
			},
		},
	})

	requireSIL003Finding(t, findings, "mapping.json", "/dynamic_templates/1/missing_type")
}

func runSIL003(t *testing.T, corpus model.Corpus) []model.Finding {
	t.Helper()
	findings, err := NewSIL003().Check(Context{}, corpus)
	if err != nil {
		t.Fatalf("SIL003 returned error: %v", err)
	}
	return findings
}

func requireSIL003Finding(t *testing.T, findings []model.Finding, file string, pointer string) {
	t.Helper()
	if len(findings) != 1 {
		t.Fatalf("SIL003 returned %d findings, want 1: %#v", len(findings), findings)
	}
	finding := findings[0]
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
	if finding.File != file {
		t.Fatalf("finding file = %q, want %q", finding.File, file)
	}
	if finding.JSONPointer != pointer {
		t.Fatalf("finding JSON pointer = %q, want %q", finding.JSONPointer, pointer)
	}
}

func testDynamicTemplate(path string, name string, pointer string, hasMatchMappingType bool) model.DynamicTemplate {
	template := model.DynamicTemplate{
		Name:                name,
		Source:              testSource(path),
		JSONPointer:         pointer,
		HasMatchMappingType: hasMatchMappingType,
		Mapping:             map[string]any{"type": "keyword"},
	}
	if hasMatchMappingType {
		template.MatchMappingType = "string"
	}
	return template
}
