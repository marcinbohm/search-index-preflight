package normalizer

import (
	"testing"

	"github.com/marcinbohm/search-index-lint/internal/model"
)

func TestNormalizeIndexTemplate(t *testing.T) {
	document := rawJSONDocument(t, model.DocumentKindIndexTemplate, "index-template.json", `{
  "index_patterns": ["logs-*"],
  "priority": 100,
  "composed_of": ["logs-common"],
  "data_stream": {},
  "_meta": {
    "owner": "example"
  },
  "template": {
    "settings": {
      "index.mapping.total_fields.limit": 1000
    },
    "mappings": {
      "dynamic": "strict",
      "properties": {
        "@timestamp": {
          "type": "date"
        },
        "message": {
          "type": "text"
        }
      }
    },
    "aliases": {
      "logs-read": {}
    }
  }
}`)

	template := NormalizeIndexTemplate(document)
	if len(template.Diagnostics) != 0 {
		t.Fatalf("NormalizeIndexTemplate returned diagnostics: %#v", template.Diagnostics)
	}
	if len(template.IndexPatterns) != 1 || template.IndexPatterns[0] != "logs-*" {
		t.Fatalf("IndexPatterns = %#v, want [logs-*]", template.IndexPatterns)
	}
	if template.Priority == nil || *template.Priority != 100 {
		t.Fatalf("Priority = %#v, want 100", template.Priority)
	}
	if len(template.ComposedOf) != 1 || template.ComposedOf[0] != "logs-common" {
		t.Fatalf("ComposedOf = %#v, want [logs-common]", template.ComposedOf)
	}
	if !template.DataStream {
		t.Fatal("DataStream = false, want true")
	}
	if template.Meta["owner"] != "example" {
		t.Fatalf("_meta owner = %#v, want example", template.Meta["owner"])
	}
	if template.Template.Settings["index.mapping.total_fields.limit"] == nil {
		t.Fatal("template settings did not preserve total_fields limit")
	}
	if template.Template.Aliases["logs-read"] == nil {
		t.Fatal("template aliases did not preserve logs-read")
	}
	if template.Template.Mappings == nil {
		t.Fatal("template mappings is nil")
	}
	if template.Template.Mappings.JSONPointer != "/template/mappings" {
		t.Fatalf("mapping JSON pointer = %q, want /template/mappings", template.Template.Mappings.JSONPointer)
	}
	if template.Template.Mappings.Dynamic != model.DynamicSettingStrict {
		t.Fatalf("mapping dynamic = %q, want strict", template.Template.Mappings.Dynamic)
	}
	requireField(t, template.Template.Mappings.Properties, "@timestamp")
	requireField(t, template.Template.Mappings.Properties, "message")
}

func TestNormalizeComponentTemplate(t *testing.T) {
	document := rawJSONDocument(t, model.DocumentKindComponentTemplate, "component-template.json", `{
  "version": 7,
  "_meta": {
    "owner": "example"
  },
  "template": {
    "settings": {
      "analysis": {
        "normalizer": {
          "lowercase": {
            "type": "custom"
          }
        }
      }
    },
    "mappings": {
      "properties": {
        "service.name": {
          "type": "keyword",
          "normalizer": "lowercase"
        }
      }
    },
    "aliases": {
      "service-read": {}
    }
  }
}`)

	template := NormalizeComponentTemplate(document)
	if len(template.Diagnostics) != 0 {
		t.Fatalf("NormalizeComponentTemplate returned diagnostics: %#v", template.Diagnostics)
	}
	if template.Version == nil || *template.Version != 7 {
		t.Fatalf("Version = %#v, want 7", template.Version)
	}
	if template.Meta["owner"] != "example" {
		t.Fatalf("_meta owner = %#v, want example", template.Meta["owner"])
	}
	if template.Template.Settings["analysis"] == nil {
		t.Fatal("template settings did not preserve analysis")
	}
	if template.Template.Aliases["service-read"] == nil {
		t.Fatal("template aliases did not preserve service-read")
	}
	if template.Template.Mappings == nil {
		t.Fatal("template mappings is nil")
	}
	if template.Template.Mappings.JSONPointer != "/template/mappings" {
		t.Fatalf("mapping JSON pointer = %q, want /template/mappings", template.Template.Mappings.JSONPointer)
	}
	field := requireField(t, template.Template.Mappings.Properties, "service.name")
	if field.Parameters["normalizer"] != "lowercase" {
		t.Fatalf("normalizer parameter = %#v, want lowercase", field.Parameters["normalizer"])
	}
}

func TestNormalizeIndexTemplateNonObjectTemplateBodyReturnsDiagnostic(t *testing.T) {
	document := rawJSONDocument(t, model.DocumentKindIndexTemplate, "index-template.json", `{
  "index_patterns": ["logs-*"],
  "template": []
}`)

	template := NormalizeIndexTemplate(document)
	requireDiagnosticContaining(t, template.Diagnostics, "template must be an object")
}
