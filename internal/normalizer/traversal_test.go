package normalizer

import (
	"reflect"
	"testing"

	"github.com/marcinbohm/search-index-lint/internal/model"
)

func TestCollectFieldsStandaloneMapping(t *testing.T) {
	corpus := Corpus{
		Mappings: []model.Mapping{normalizedMapping(t, "mapping.json", `{
  "properties": {
    "status": {
      "type": "keyword"
    },
    "message": {
      "type": "text",
      "fields": {
        "keyword": {
          "type": "keyword"
        }
      }
    },
    "user": {
      "type": "object",
      "properties": {
        "id": {
          "type": "keyword"
        }
      }
    }
  },
  "runtime": {
    "day": {
      "type": "keyword"
    }
  }
}`)},
	}

	visits := CollectFields(corpus)
	got := visitPathRoles(visits)
	want := []pathRole{
		{path: "message", role: model.FieldRoleProperty},
		{path: "message.keyword", role: model.FieldRoleMultiField},
		{path: "status", role: model.FieldRoleProperty},
		{path: "user", role: model.FieldRoleProperty},
		{path: "user.id", role: model.FieldRoleProperty},
		{path: "day", role: model.FieldRoleRuntimeField},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("CollectFields returned path/roles %#v, want %#v", got, want)
	}
	for _, visit := range visits {
		if visit.Origin != model.FieldOriginMapping {
			t.Fatalf("Origin = %q, want %q", visit.Origin, model.FieldOriginMapping)
		}
		if visit.Source.RelativePath != "mapping.json" {
			t.Fatalf("Source.RelativePath = %q, want mapping.json", visit.Source.RelativePath)
		}
	}
}

func TestCollectFieldsIndexTemplate(t *testing.T) {
	document := rawJSONDocument(t, model.DocumentKindIndexTemplate, "index-template.json", `{
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
	template := NormalizeIndexTemplate(document)
	if len(template.Diagnostics) != 0 {
		t.Fatalf("NormalizeIndexTemplate returned diagnostics: %#v", template.Diagnostics)
	}

	visits := CollectFields(Corpus{IndexTemplates: []model.IndexTemplate{template}})
	if len(visits) != 1 {
		t.Fatalf("CollectFields returned %d visits, want 1", len(visits))
	}
	visit := visits[0]
	if visit.Origin != model.FieldOriginIndexTemplate {
		t.Fatalf("Origin = %q, want %q", visit.Origin, model.FieldOriginIndexTemplate)
	}
	if visit.Role != model.FieldRoleProperty {
		t.Fatalf("Role = %q, want %q", visit.Role, model.FieldRoleProperty)
	}
	if visit.Source.RelativePath != "index-template.json" {
		t.Fatalf("Source.RelativePath = %q, want index-template.json", visit.Source.RelativePath)
	}
	if visit.IndexTemplateName != "" {
		t.Fatalf("IndexTemplateName = %q, want empty", visit.IndexTemplateName)
	}
}

func TestCollectFieldsComponentTemplate(t *testing.T) {
	document := rawJSONDocument(t, model.DocumentKindComponentTemplate, "component-template.json", `{
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
	template := NormalizeComponentTemplate(document)
	if len(template.Diagnostics) != 0 {
		t.Fatalf("NormalizeComponentTemplate returned diagnostics: %#v", template.Diagnostics)
	}

	visits := CollectFields(Corpus{ComponentTemplates: []model.ComponentTemplate{template}})
	if len(visits) != 1 {
		t.Fatalf("CollectFields returned %d visits, want 1", len(visits))
	}
	visit := visits[0]
	if visit.Origin != model.FieldOriginComponentTemplate {
		t.Fatalf("Origin = %q, want %q", visit.Origin, model.FieldOriginComponentTemplate)
	}
	if visit.Role != model.FieldRoleProperty {
		t.Fatalf("Role = %q, want %q", visit.Role, model.FieldRoleProperty)
	}
	if visit.ComponentTemplateName != "" {
		t.Fatalf("ComponentTemplateName = %q, want empty", visit.ComponentTemplateName)
	}
}

func TestCountFields(t *testing.T) {
	corpus := Corpus{
		Mappings: []model.Mapping{normalizedMapping(t, "mapping.json", `{
  "properties": {
    "status": {
      "type": "keyword"
    },
    "message": {
      "type": "text",
      "fields": {
        "keyword": {
          "type": "keyword"
        }
      }
    },
    "user": {
      "type": "object",
      "properties": {
        "id": {
          "type": "keyword"
        }
      }
    }
  },
  "runtime": {
    "day": {
      "type": "keyword"
    }
  }
}`)},
	}

	stats := CountFields(corpus)
	if stats.Properties != 4 {
		t.Fatalf("Properties = %d, want 4", stats.Properties)
	}
	if stats.MultiFields != 1 {
		t.Fatalf("MultiFields = %d, want 1", stats.MultiFields)
	}
	if stats.RuntimeFields != 1 {
		t.Fatalf("RuntimeFields = %d, want 1", stats.RuntimeFields)
	}
	if stats.TotalFields != 6 {
		t.Fatalf("TotalFields = %d, want 6", stats.TotalFields)
	}
}

type pathRole struct {
	path string
	role model.FieldRole
}

func visitPathRoles(visits []model.FieldVisit) []pathRole {
	values := make([]pathRole, 0, len(visits))
	for _, visit := range visits {
		values = append(values, pathRole{path: visit.Path, role: visit.Role})
	}
	return values
}

func normalizedMapping(t *testing.T, path, content string) model.Mapping {
	t.Helper()
	mapping := NormalizeMapping(rawJSONDocument(t, model.DocumentKindMapping, path, content))
	if len(mapping.Diagnostics) != 0 {
		t.Fatalf("NormalizeMapping returned diagnostics: %#v", mapping.Diagnostics)
	}
	return mapping
}
