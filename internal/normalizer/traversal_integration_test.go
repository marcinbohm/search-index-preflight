package normalizer

import (
	"reflect"
	"testing"

	"github.com/marcinbohm/search-index-preflight/internal/model"
	"github.com/marcinbohm/search-index-preflight/internal/parser"
)

func TestNormalizeCorpusOutputWorksWithModelCollectFields(t *testing.T) {
	source := model.Source{Path: "mapping.json", RelativePath: "mapping.json"}
	document := parser.ParseJSON(source, model.DocumentKindMapping, []byte(`{
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
}`))
	if len(document.Diagnostics) != 0 {
		t.Fatalf("ParseJSON returned diagnostics: %#v", document.Diagnostics)
	}

	corpus := Normalize([]model.RawDocument{document})
	if len(corpus.Diagnostics) != 0 {
		t.Fatalf("Normalize returned diagnostics: %#v", corpus.Diagnostics)
	}
	if len(corpus.Mappings) != 1 {
		t.Fatalf("Normalize returned %d mappings, want 1", len(corpus.Mappings))
	}

	visits := model.CollectFields(corpus)
	got := integrationVisitPathRoles(visits)
	want := []integrationPathRole{
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
}

type integrationPathRole struct {
	path string
	role model.FieldRole
}

func integrationVisitPathRoles(visits []model.FieldVisit) []integrationPathRole {
	values := make([]integrationPathRole, 0, len(visits))
	for _, visit := range visits {
		values = append(values, integrationPathRole{path: visit.Path, role: visit.Role})
	}
	return values
}
