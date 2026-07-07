package diff

import (
	"reflect"
	"testing"

	"github.com/marcinbohm/search-index-preflight/internal/model"
	"github.com/marcinbohm/search-index-preflight/internal/normalizer"
)

func TestCompareNoChanges(t *testing.T) {
	base := corpusWithMapping("mapping.json", "", field("status", model.FieldRoleProperty, "keyword"))
	current := corpusWithMapping("mapping.json", "", field("status", model.FieldRoleProperty, "keyword"))

	result := Compare(base, current)

	if len(result.FieldChanges) != 0 {
		t.Fatalf("expected no changes, got %#v", result.FieldChanges)
	}
}

func TestCompareFieldAdded(t *testing.T) {
	base := corpusWithMapping("mapping.json", "", field("status", model.FieldRoleProperty, "keyword"))
	current := corpusWithMapping(
		"mapping.json",
		"",
		field("status", model.FieldRoleProperty, "keyword"),
		field("service.name", model.FieldRoleProperty, "keyword"),
	)

	change := singleChange(t, Compare(base, current))

	assertChange(t, change, ChangeFieldAdded, ResourceMapping, "service.name", model.FieldRoleProperty)
	if change.After == nil || change.After.Type != "keyword" {
		t.Fatalf("expected after snapshot with keyword type, got %#v", change.After)
	}
}

func TestCompareFieldRemoved(t *testing.T) {
	base := corpusWithMapping(
		"mapping.json",
		"",
		field("status", model.FieldRoleProperty, "keyword"),
		field("service.name", model.FieldRoleProperty, "keyword"),
	)
	current := corpusWithMapping("mapping.json", "", field("status", model.FieldRoleProperty, "keyword"))

	change := singleChange(t, Compare(base, current))

	assertChange(t, change, ChangeFieldRemoved, ResourceMapping, "service.name", model.FieldRoleProperty)
	if change.Before == nil || change.Before.Type != "keyword" {
		t.Fatalf("expected before snapshot with keyword type, got %#v", change.Before)
	}
}

func TestCompareFieldTypeChanged(t *testing.T) {
	base := corpusWithMapping("mapping.json", "", field("status", model.FieldRoleProperty, "keyword"))
	current := corpusWithMapping("mapping.json", "", field("status", model.FieldRoleProperty, "long"))

	change := singleChange(t, Compare(base, current))

	assertChange(t, change, ChangeFieldTypeChanged, ResourceMapping, "status", model.FieldRoleProperty)
	if change.Before == nil || change.Before.Type != "keyword" {
		t.Fatalf("expected before keyword type, got %#v", change.Before)
	}
	if change.After == nil || change.After.Type != "long" {
		t.Fatalf("expected after long type, got %#v", change.After)
	}
}

func TestCompareMultiFieldAdded(t *testing.T) {
	base := corpusWithMapping("mapping.json", "", field("message", model.FieldRoleProperty, "text"))
	current := model.Corpus{
		Mappings: []model.Mapping{
			{
				Source: source("mapping.json"),
				Properties: []model.Field{
					{
						Name:        "message",
						Path:        "message",
						Type:        "text",
						Source:      source("mapping.json"),
						JSONPointer: "/properties/message",
						Fields: []model.Field{
							field("message.keyword", model.FieldRoleMultiField, "keyword"),
						},
					},
				},
			},
		},
	}

	change := singleChange(t, Compare(base, current))

	assertChange(t, change, ChangeFieldAdded, ResourceMapping, "message.keyword", model.FieldRoleMultiField)
}

func TestCompareRuntimeFieldAdded(t *testing.T) {
	base := corpusWithMapping("mapping.json", "")
	current := model.Corpus{
		Mappings: []model.Mapping{
			{
				Source:        source("mapping.json"),
				RuntimeFields: []model.Field{field("day_of_week", model.FieldRoleRuntimeField, "keyword")},
			},
		},
	}

	change := singleChange(t, Compare(base, current))

	assertChange(t, change, ChangeFieldAdded, ResourceMapping, "day_of_week", model.FieldRoleRuntimeField)
}

func TestCompareWrappedMappingResourcePointer(t *testing.T) {
	base := corpusWithMapping("mapping.json", "/mappings")
	current := corpusWithMapping("mapping.json", "/mappings", field("status", model.FieldRoleProperty, "keyword"))

	change := singleChange(t, Compare(base, current))

	if change.Resource.JSONPointer != "/mappings" {
		t.Fatalf("expected resource pointer /mappings, got %q", change.Resource.JSONPointer)
	}
}

func TestCompareIndexTemplateMappingFieldTypeChanged(t *testing.T) {
	base := corpusWithIndexTemplate("template.json", field("status", model.FieldRoleProperty, "keyword"))
	current := corpusWithIndexTemplate("template.json", field("status", model.FieldRoleProperty, "long"))

	change := singleChange(t, Compare(base, current))

	assertChange(t, change, ChangeFieldTypeChanged, ResourceIndexTemplate, "status", model.FieldRoleProperty)
	if change.Resource.JSONPointer != "/template/mappings" {
		t.Fatalf("expected resource pointer /template/mappings, got %q", change.Resource.JSONPointer)
	}
}

func TestCompareComponentTemplateMappingFieldAdded(t *testing.T) {
	base := corpusWithComponentTemplate("component.json")
	current := corpusWithComponentTemplate("component.json", field("status", model.FieldRoleProperty, "keyword"))

	change := singleChange(t, Compare(base, current))

	assertChange(t, change, ChangeFieldAdded, ResourceComponentTemplate, "status", model.FieldRoleProperty)
}

func TestCompareDeterministicOrdering(t *testing.T) {
	base := corpusWithMapping("b.json", "", field("z", model.FieldRoleProperty, "keyword"))
	base.Mappings = append(base.Mappings, mapping("a.json", "", field("b", model.FieldRoleProperty, "keyword")))
	current := corpusWithMapping("b.json", "", field("a", model.FieldRoleProperty, "keyword"))
	current.Mappings = append(current.Mappings, mapping("a.json", "", field("a", model.FieldRoleProperty, "keyword")))

	result := Compare(base, current)

	got := changeOrder(result.FieldChanges)
	want := []string{
		"mapping|a.json||a|property|field_added",
		"mapping|a.json||b|property|field_removed",
		"mapping|b.json||a|property|field_added",
		"mapping|b.json||z|property|field_removed",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected change order\ngot:  %#v\nwant: %#v", got, want)
	}
}

func TestCompareNormalizerIntegration(t *testing.T) {
	base := normalizer.Normalize([]model.RawDocument{rawMappingDocument("mapping.json", map[string]any{
		"properties": map[string]any{
			"status": map[string]any{"type": "keyword"},
		},
	})})
	current := normalizer.Normalize([]model.RawDocument{rawMappingDocument("mapping.json", map[string]any{
		"properties": map[string]any{
			"status": map[string]any{"type": "long"},
		},
	})})

	if len(base.Diagnostics) != 0 {
		t.Fatalf("expected no base diagnostics, got %#v", base.Diagnostics)
	}
	if len(current.Diagnostics) != 0 {
		t.Fatalf("expected no current diagnostics, got %#v", current.Diagnostics)
	}

	change := singleChange(t, Compare(base, current))

	assertChange(t, change, ChangeFieldTypeChanged, ResourceMapping, "status", model.FieldRoleProperty)
}

func singleChange(t *testing.T, result Result) FieldChange {
	t.Helper()
	if len(result.FieldChanges) != 1 {
		t.Fatalf("expected one field change, got %#v", result.FieldChanges)
	}
	return result.FieldChanges[0]
}

func assertChange(t *testing.T, change FieldChange, kind ChangeKind, resourceKind ResourceKind, path string, role model.FieldRole) {
	t.Helper()
	if change.Kind != kind {
		t.Fatalf("expected change kind %q, got %q", kind, change.Kind)
	}
	if change.Resource.Kind != resourceKind {
		t.Fatalf("expected resource kind %q, got %q", resourceKind, change.Resource.Kind)
	}
	if change.Field.Path != path {
		t.Fatalf("expected field path %q, got %q", path, change.Field.Path)
	}
	if change.Field.Role != role {
		t.Fatalf("expected field role %q, got %q", role, change.Field.Role)
	}
}

func corpusWithMapping(file, pointer string, fields ...model.Field) model.Corpus {
	return model.Corpus{Mappings: []model.Mapping{mapping(file, pointer, fields...)}}
}

func mapping(file, pointer string, fields ...model.Field) model.Mapping {
	return model.Mapping{
		Source:        source(file),
		JSONPointer:   pointer,
		Properties:    fieldsByRole(fields, model.FieldRoleProperty),
		RuntimeFields: fieldsByRole(fields, model.FieldRoleRuntimeField),
	}
}

func corpusWithIndexTemplate(file string, fields ...model.Field) model.Corpus {
	mapping := mapping(file, "/template/mappings", fields...)
	return model.Corpus{
		IndexTemplates: []model.IndexTemplate{
			{
				Source: source(file),
				Template: model.TemplateBody{
					Mappings: &mapping,
				},
			},
		},
	}
}

func corpusWithComponentTemplate(file string, fields ...model.Field) model.Corpus {
	mapping := mapping(file, "/template/mappings", fields...)
	return model.Corpus{
		ComponentTemplates: []model.ComponentTemplate{
			{
				Source: source(file),
				Template: model.TemplateBody{
					Mappings: &mapping,
				},
			},
		},
	}
}

func fieldsByRole(fields []model.Field, role model.FieldRole) []model.Field {
	var filtered []model.Field
	for _, field := range fields {
		if role == model.FieldRoleProperty && fieldRole(field) == model.FieldRoleProperty {
			filtered = append(filtered, field)
		}
		if role == model.FieldRoleRuntimeField && fieldRole(field) == model.FieldRoleRuntimeField {
			filtered = append(filtered, field)
		}
	}
	return filtered
}

func fieldRole(field model.Field) model.FieldRole {
	if field.JSONPointer == "/runtime/"+field.Path {
		return model.FieldRoleRuntimeField
	}
	if field.JSONPointer == "/fields/"+field.Name {
		return model.FieldRoleMultiField
	}
	return model.FieldRoleProperty
}

func field(path string, role model.FieldRole, fieldType string) model.Field {
	name := path
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '.' {
			name = path[i+1:]
			break
		}
	}

	pointer := "/properties/" + path
	if role == model.FieldRoleRuntimeField {
		pointer = "/runtime/" + path
	}
	if role == model.FieldRoleMultiField {
		pointer = "/fields/" + name
	}

	return model.Field{
		Name:        name,
		Path:        path,
		Type:        fieldType,
		Source:      source("mapping.json"),
		JSONPointer: pointer,
	}
}

func source(path string) model.Source {
	return model.Source{Path: path, RelativePath: path}
}

func rawMappingDocument(path string, content map[string]any) model.RawDocument {
	return model.RawDocument{
		Kind:    model.DocumentKindMapping,
		Source:  source(path),
		Content: content,
	}
}

func changeOrder(changes []FieldChange) []string {
	order := make([]string, 0, len(changes))
	for _, change := range changes {
		order = append(order, string(change.Resource.Kind)+"|"+change.Resource.File+"|"+change.Resource.JSONPointer+"|"+change.Field.Path+"|"+string(change.Field.Role)+"|"+string(change.Kind))
	}
	return order
}
