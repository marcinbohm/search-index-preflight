package normalizer

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/marcinbohm/search-index-preflight/internal/model"
	"github.com/marcinbohm/search-index-preflight/internal/parser"
)

func TestNormalizeRawMapping(t *testing.T) {
	document := rawJSONDocument(t, model.DocumentKindMapping, "mapping.json", `{
  "dynamic": "strict",
  "runtime": {
    "day": {
      "type": "keyword",
      "script": {
        "source": "emit(doc['@timestamp'].value.dayOfWeekEnum.toString())"
      }
    }
  },
  "properties": {
    "status": {
      "type": "keyword",
      "ignore_above": 256
    },
    "message": {
      "type": "text",
      "fields": {
        "keyword": {
          "type": "keyword",
          "ignore_above": 1024
        }
      }
    },
    "user": {
      "type": "object",
      "dynamic": true,
      "properties": {
        "id": {
          "type": "keyword"
        }
      }
    }
  },
  "dynamic_templates": [
    {
      "strings_as_keyword": {
        "match_mapping_type": "string",
        "match": "*",
        "mapping": {
          "type": "keyword"
        }
      }
    }
  ]
}`)

	mapping := NormalizeMapping(document)
	if len(mapping.Diagnostics) != 0 {
		t.Fatalf("NormalizeMapping returned diagnostics: %#v", mapping.Diagnostics)
	}
	if mapping.Dynamic != model.DynamicSettingStrict {
		t.Fatalf("Dynamic = %q, want %q", mapping.Dynamic, model.DynamicSettingStrict)
	}
	if mapping.JSONPointer != "" {
		t.Fatalf("mapping JSON pointer = %q, want empty root pointer", mapping.JSONPointer)
	}

	status := requireField(t, mapping.Properties, "status")
	if status.Path != "status" {
		t.Fatalf("status path = %q, want status", status.Path)
	}
	if got := numberString(status.Parameters["ignore_above"]); got != "256" {
		t.Fatalf("status ignore_above = %q, want 256", got)
	}

	message := requireField(t, mapping.Properties, "message")
	keyword := requireField(t, message.Fields, "keyword")
	if keyword.Path != "message.keyword" {
		t.Fatalf("keyword path = %q, want message.keyword", keyword.Path)
	}
	if got := numberString(keyword.Parameters["ignore_above"]); got != "1024" {
		t.Fatalf("message.keyword ignore_above = %q, want 1024", got)
	}

	user := requireField(t, mapping.Properties, "user")
	if user.Dynamic != model.DynamicSettingTrue {
		t.Fatalf("user dynamic = %q, want %q", user.Dynamic, model.DynamicSettingTrue)
	}
	userID := requireField(t, user.Properties, "id")
	if userID.Path != "user.id" {
		t.Fatalf("user.id path = %q, want user.id", userID.Path)
	}
	if userID.ParentPath != "user" {
		t.Fatalf("user.id parent = %q, want user", userID.ParentPath)
	}

	if len(mapping.DynamicTemplates) != 1 {
		t.Fatalf("dynamic template count = %d, want 1", len(mapping.DynamicTemplates))
	}
	template := mapping.DynamicTemplates[0]
	if template.Name != "strings_as_keyword" {
		t.Fatalf("dynamic template name = %q, want strings_as_keyword", template.Name)
	}
	if template.Match != "*" || template.MatchMappingType != "string" {
		t.Fatalf("dynamic template fields = %#v", template)
	}
	if !template.HasMatchMappingType {
		t.Fatal("dynamic template HasMatchMappingType = false, want true")
	}
	if template.Mapping["type"] != "keyword" {
		t.Fatalf("dynamic template mapping type = %#v, want keyword", template.Mapping["type"])
	}

	day := requireField(t, mapping.RuntimeFields, "day")
	if day.Type != "keyword" {
		t.Fatalf("runtime field type = %q, want keyword", day.Type)
	}
	if _, ok := day.Parameters["script"]; !ok {
		t.Fatal("runtime field script parameter was not preserved")
	}
}

func TestNormalizeWrappedMapping(t *testing.T) {
	document := rawJSONDocument(t, model.DocumentKindMapping, "mapping.json", `{
  "mappings": {
    "properties": {
      "status": {
        "type": "keyword"
      }
    }
  }
}`)

	mapping := NormalizeMapping(document)
	if len(mapping.Diagnostics) != 0 {
		t.Fatalf("NormalizeMapping returned diagnostics: %#v", mapping.Diagnostics)
	}
	if mapping.JSONPointer != "/mappings" {
		t.Fatalf("mapping JSON pointer = %q, want /mappings", mapping.JSONPointer)
	}
	status := requireField(t, mapping.Properties, "status")
	if status.JSONPointer != "/mappings/properties/status" {
		t.Fatalf("status JSON pointer = %q, want /mappings/properties/status", status.JSONPointer)
	}
}

func TestNormalizeDynamicTemplatesPreserveOrder(t *testing.T) {
	document := rawJSONDocument(t, model.DocumentKindMapping, "mapping.json", `{
  "dynamic_templates": [
    {
      "first_template": {
        "match": "first_*",
        "mapping": {
          "type": "keyword"
        }
      }
    },
    {
      "second_template": {
        "match": "second_*",
        "mapping": {
          "type": "text"
        }
      }
    }
  ]
}`)

	mapping := NormalizeMapping(document)
	if len(mapping.DynamicTemplates) != 2 {
		t.Fatalf("dynamic template count = %d, want 2", len(mapping.DynamicTemplates))
	}
	if mapping.DynamicTemplates[0].Name != "first_template" || mapping.DynamicTemplates[1].Name != "second_template" {
		t.Fatalf("dynamic template order = %#v", mapping.DynamicTemplates)
	}
}

func TestNormalizeDynamicTemplateMissingMatchMappingType(t *testing.T) {
	document := rawJSONDocument(t, model.DocumentKindMapping, "mapping.json", `{
  "dynamic_templates": [
    {
      "strings_as_keywords": {
        "mapping": {
          "type": "keyword"
        }
      }
    }
  ]
}`)

	mapping := NormalizeMapping(document)
	if len(mapping.Diagnostics) != 0 {
		t.Fatalf("NormalizeMapping returned diagnostics: %#v", mapping.Diagnostics)
	}
	if len(mapping.DynamicTemplates) != 1 {
		t.Fatalf("dynamic template count = %d, want 1", len(mapping.DynamicTemplates))
	}
	template := mapping.DynamicTemplates[0]
	if template.HasMatchMappingType {
		t.Fatal("dynamic template HasMatchMappingType = true, want false")
	}
	if template.JSONPointer != "/dynamic_templates/0/strings_as_keywords" {
		t.Fatalf("dynamic template JSON pointer = %q, want /dynamic_templates/0/strings_as_keywords", template.JSONPointer)
	}
}

func TestNormalizeDynamicTemplateEscapedNamePointer(t *testing.T) {
	document := rawJSONDocument(t, model.DocumentKindMapping, "mapping.json", `{
  "dynamic_templates": [
    {
      "service/name~template": {
        "mapping": {
          "type": "keyword"
        }
      }
    }
  ]
}`)

	mapping := NormalizeMapping(document)
	if len(mapping.Diagnostics) != 0 {
		t.Fatalf("NormalizeMapping returned diagnostics: %#v", mapping.Diagnostics)
	}
	if len(mapping.DynamicTemplates) != 1 {
		t.Fatalf("dynamic template count = %d, want 1", len(mapping.DynamicTemplates))
	}
	template := mapping.DynamicTemplates[0]
	if template.JSONPointer != "/dynamic_templates/0/service~1name~0template" {
		t.Fatalf("dynamic template JSON pointer = %q, want /dynamic_templates/0/service~1name~0template", template.JSONPointer)
	}
}

func TestNormalizeMappingMalformedPropertiesReturnsDiagnostic(t *testing.T) {
	document := rawJSONDocument(t, model.DocumentKindMapping, "mapping.json", `{
  "properties": []
}`)

	mapping := NormalizeMapping(document)
	requireDiagnosticContaining(t, mapping.Diagnostics, "properties must be an object")
}

func TestNormalizeMappingMalformedDynamicTemplatesReturnsDiagnostic(t *testing.T) {
	document := rawJSONDocument(t, model.DocumentKindMapping, "mapping.json", `{
  "dynamic_templates": {}
}`)

	mapping := NormalizeMapping(document)
	requireDiagnosticContaining(t, mapping.Diagnostics, "dynamic_templates must be an array")
}

func TestNormalizeMappingNonObjectDocumentReturnsDiagnostic(t *testing.T) {
	document := rawJSONDocument(t, model.DocumentKindMapping, "mapping.json", `[]`)

	mapping := NormalizeMapping(document)
	requireDiagnosticContaining(t, mapping.Diagnostics, "mapping document must be a JSON object")
}

func rawJSONDocument(t *testing.T, kind model.DocumentKind, path, content string) model.RawDocument {
	t.Helper()
	document := parser.ParseJSON(model.Source{Path: path, RelativePath: path}, kind, []byte(content))
	if len(document.Diagnostics) != 0 {
		t.Fatalf("ParseJSON returned diagnostics: %#v", document.Diagnostics)
	}
	return document
}

func requireField(t *testing.T, fields []model.Field, name string) model.Field {
	t.Helper()
	for _, field := range fields {
		if field.Name == name {
			return field
		}
	}
	t.Fatalf("field %q not found in %#v", name, fields)
	return model.Field{}
}

func numberString(value any) string {
	number, ok := value.(json.Number)
	if !ok {
		return ""
	}
	return number.String()
}

func requireDiagnosticContaining(t *testing.T, diagnostics []model.Diagnostic, want string) {
	t.Helper()
	if len(diagnostics) == 0 {
		t.Fatalf("got no diagnostics, want one containing %q", want)
	}
	for _, diagnostic := range diagnostics {
		if strings.Contains(diagnostic.Message, want) {
			return
		}
	}
	t.Fatalf("diagnostics %#v do not contain message %q", diagnostics, want)
}
