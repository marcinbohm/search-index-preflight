package diff

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/marcinbohm/search-index-preflight/internal/model"
	"github.com/marcinbohm/search-index-preflight/internal/normalizer"
)

func TestCompareNoChanges(t *testing.T) {
	base := corpusWithMapping("mapping.json", "", property("status", "keyword"))
	current := corpusWithMapping("mapping.json", "", property("status", "keyword"))

	result := mustCompare(t, base, current)

	if len(result.FieldChanges) != 0 {
		t.Fatalf("expected no changes, got %#v", result.FieldChanges)
	}
}

func TestCompareFieldAdded(t *testing.T) {
	base := corpusWithMapping("mapping.json", "", property("status", "keyword"))
	current := corpusWithMapping(
		"mapping.json",
		"",
		property("status", "keyword"),
		property("service.name", "keyword"),
	)

	change := singleChange(t, mustCompare(t, base, current))

	assertChange(t, change, ChangeFieldAdded, ResourceMapping, "service.name", model.FieldRoleProperty)
	assertAfterSnapshot(t, change, "service.name", model.FieldRoleProperty, "keyword", "/properties/service.name")
}

func TestCompareFieldRemoved(t *testing.T) {
	base := corpusWithMapping(
		"mapping.json",
		"",
		property("status", "keyword"),
		property("service.name", "keyword"),
	)
	current := corpusWithMapping("mapping.json", "", property("status", "keyword"))

	change := singleChange(t, mustCompare(t, base, current))

	assertChange(t, change, ChangeFieldRemoved, ResourceMapping, "service.name", model.FieldRoleProperty)
	assertBeforeSnapshot(t, change, "service.name", model.FieldRoleProperty, "keyword", "/properties/service.name")
}

func TestCompareFieldTypeChanged(t *testing.T) {
	base := corpusWithMapping("mapping.json", "", property("status", "keyword"))
	current := corpusWithMapping("mapping.json", "", property("status", "long"))

	change := singleChange(t, mustCompare(t, base, current))

	assertChange(t, change, ChangeFieldTypeChanged, ResourceMapping, "status", model.FieldRoleProperty)
	assertBeforeSnapshot(t, change, "status", model.FieldRoleProperty, "keyword", "/properties/status")
	assertAfterSnapshot(t, change, "status", model.FieldRoleProperty, "long", "/properties/status")
}

func TestCompareMultiFieldAdded(t *testing.T) {
	base := corpusWithMapping("mapping.json", "", property("message", "text"))
	current := corpusWithMapping("mapping.json", "", textWithMultiField("message", "keyword", "keyword"))

	change := singleChange(t, mustCompare(t, base, current))

	assertChange(t, change, ChangeFieldAdded, ResourceMapping, "message.keyword", model.FieldRoleMultiField)
	assertAfterSnapshot(t, change, "message.keyword", model.FieldRoleMultiField, "keyword", "/properties/message/fields/keyword")
}

func TestCompareRuntimeFieldAdded(t *testing.T) {
	base := corpusWithMapping("mapping.json", "")
	current := corpusWithMapping("mapping.json", "", runtimeField("day_of_week", "keyword"))

	change := singleChange(t, mustCompare(t, base, current))

	assertChange(t, change, ChangeFieldAdded, ResourceMapping, "day_of_week", model.FieldRoleRuntimeField)
	assertAfterSnapshot(t, change, "day_of_week", model.FieldRoleRuntimeField, "keyword", "/runtime/day_of_week")
}

func TestCompareWrappedMappingResourcePointer(t *testing.T) {
	base := corpusWithMapping("mapping.json", "/mappings")
	current := corpusWithMapping("mapping.json", "/mappings", property("status", "keyword"))

	change := singleChange(t, mustCompare(t, base, current))

	if change.Resource.JSONPointer != "/mappings" {
		t.Fatalf("expected resource pointer /mappings, got %q", change.Resource.JSONPointer)
	}
}

func TestCompareIndexTemplateMappingFieldTypeChanged(t *testing.T) {
	base := corpusWithIndexTemplate("template.json", property("status", "keyword"))
	current := corpusWithIndexTemplate("template.json", property("status", "long"))

	change := singleChange(t, mustCompare(t, base, current))

	assertChange(t, change, ChangeFieldTypeChanged, ResourceIndexTemplate, "status", model.FieldRoleProperty)
	if change.Resource.JSONPointer != "/template/mappings" {
		t.Fatalf("expected resource pointer /template/mappings, got %q", change.Resource.JSONPointer)
	}
}

func TestCompareComponentTemplateMappingFieldAdded(t *testing.T) {
	base := corpusWithComponentTemplate("component.json")
	current := corpusWithComponentTemplate("component.json", property("status", "keyword"))

	change := singleChange(t, mustCompare(t, base, current))

	assertChange(t, change, ChangeFieldAdded, ResourceComponentTemplate, "status", model.FieldRoleProperty)
}

func TestCompareDuplicateStandaloneMappingResourceInBase(t *testing.T) {
	base := model.Corpus{
		Mappings: []model.Mapping{
			mapping("mapping.json", "", property("status", "keyword")),
			mapping("mapping.json", "", property("service.name", "keyword")),
		},
	}

	result, err := Compare(base, model.Corpus{})

	assertDuplicateResourceError(t, err, ResourceID{Kind: ResourceMapping, File: "mapping.json"})
	if len(result.FieldChanges) != 0 {
		t.Fatalf("expected no partial changes on duplicate resource error, got %#v", result.FieldChanges)
	}
}

func TestCompareDuplicateStandaloneMappingResourceInCurrent(t *testing.T) {
	current := model.Corpus{
		Mappings: []model.Mapping{
			mapping("mapping.json", "", property("status", "keyword")),
			mapping("mapping.json", "", property("service.name", "keyword")),
		},
	}

	_, err := Compare(model.Corpus{}, current)

	assertDuplicateResourceError(t, err, ResourceID{Kind: ResourceMapping, File: "mapping.json"})
}

func TestCompareAllowsSameFileWithDifferentJSONPointers(t *testing.T) {
	base := model.Corpus{
		Mappings: []model.Mapping{
			mapping("mapping.json", "", property("status", "keyword")),
			mapping("mapping.json", "/mappings", property("status", "keyword")),
		},
	}
	current := base

	result := mustCompare(t, base, current)

	if len(result.FieldChanges) != 0 {
		t.Fatalf("expected no changes, got %#v", result.FieldChanges)
	}
}

func TestCompareAllowsSameFileAndPointerWithDifferentResourceKind(t *testing.T) {
	base := model.Corpus{
		Mappings: []model.Mapping{
			mapping("mapping.json", "/template/mappings", property("status", "keyword")),
		},
		IndexTemplates: []model.IndexTemplate{
			indexTemplate("mapping.json", property("status", "keyword")),
		},
	}
	current := base

	result := mustCompare(t, base, current)

	if len(result.FieldChanges) != 0 {
		t.Fatalf("expected no changes, got %#v", result.FieldChanges)
	}
}

func TestCompareSamePathDifferentRoleIsNotCollapsed(t *testing.T) {
	base := corpusWithMapping("mapping.json", "", property("message.keyword", "keyword"))
	current := corpusWithMapping("mapping.json", "", textWithMultiField("message", "keyword", "keyword"))

	result := mustCompare(t, base, current)

	assertHasChange(t, result.FieldChanges, ChangeFieldRemoved, "message.keyword", model.FieldRoleProperty)
	assertHasChange(t, result.FieldChanges, ChangeFieldAdded, "message", model.FieldRoleProperty)
	assertHasChange(t, result.FieldChanges, ChangeFieldAdded, "message.keyword", model.FieldRoleMultiField)
	assertNoChange(t, result.FieldChanges, ChangeFieldTypeChanged, "message.keyword")
}

func TestCompareMultiFieldTypeChanged(t *testing.T) {
	base := corpusWithMapping("mapping.json", "", textWithMultiField("message", "keyword", "keyword"))
	current := corpusWithMapping("mapping.json", "", textWithMultiField("message", "keyword", "wildcard"))

	change := singleChange(t, mustCompare(t, base, current))

	assertChange(t, change, ChangeFieldTypeChanged, ResourceMapping, "message.keyword", model.FieldRoleMultiField)
	assertBeforeSnapshot(t, change, "message.keyword", model.FieldRoleMultiField, "keyword", "/properties/message/fields/keyword")
	assertAfterSnapshot(t, change, "message.keyword", model.FieldRoleMultiField, "wildcard", "/properties/message/fields/keyword")
}

func TestCompareRuntimeFieldTypeChanged(t *testing.T) {
	base := corpusWithMapping("mapping.json", "", runtimeField("day_of_week", "keyword"))
	current := corpusWithMapping("mapping.json", "", runtimeField("day_of_week", "long"))

	change := singleChange(t, mustCompare(t, base, current))

	assertChange(t, change, ChangeFieldTypeChanged, ResourceMapping, "day_of_week", model.FieldRoleRuntimeField)
	assertBeforeSnapshot(t, change, "day_of_week", model.FieldRoleRuntimeField, "keyword", "/runtime/day_of_week")
	assertAfterSnapshot(t, change, "day_of_week", model.FieldRoleRuntimeField, "long", "/runtime/day_of_week")
}

func TestCompareEntireResourceAddedWithMultipleFields(t *testing.T) {
	current := corpusWithMapping("mapping.json", "", property("status", "keyword"), property("service.name", "keyword"))

	result := mustCompare(t, model.Corpus{}, current)

	got := changeOrder(result.FieldChanges)
	want := []string{
		"mapping|mapping.json||service.name|property|field_added",
		"mapping|mapping.json||status|property|field_added",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected changes\ngot:  %#v\nwant: %#v", got, want)
	}
}

func TestCompareEntireResourceRemovedWithMultipleFields(t *testing.T) {
	base := corpusWithMapping("mapping.json", "", property("status", "keyword"), property("service.name", "keyword"))

	result := mustCompare(t, base, model.Corpus{})

	got := changeOrder(result.FieldChanges)
	want := []string{
		"mapping|mapping.json||service.name|property|field_removed",
		"mapping|mapping.json||status|property|field_removed",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected changes\ngot:  %#v\nwant: %#v", got, want)
	}
}

func TestCompareDeterministicOrdering(t *testing.T) {
	base := corpusWithMapping("b.json", "", property("z", "keyword"))
	base.Mappings = append(base.Mappings, mapping("a.json", "", property("b", "keyword")))
	current := corpusWithMapping("b.json", "", property("a", "keyword"))
	current.Mappings = append(current.Mappings, mapping("a.json", "", property("a", "keyword")))

	result := mustCompare(t, base, current)

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

	change := singleChange(t, mustCompare(t, base, current))

	assertChange(t, change, ChangeFieldTypeChanged, ResourceMapping, "status", model.FieldRoleProperty)
}

func mustCompare(t *testing.T, base model.Corpus, current model.Corpus) Result {
	t.Helper()
	result, err := Compare(base, current)
	if err != nil {
		t.Fatalf("Compare returned error: %v", err)
	}
	return result
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

func assertBeforeSnapshot(t *testing.T, change FieldChange, path string, role model.FieldRole, typ string, pointer string) {
	t.Helper()
	if change.Before == nil {
		t.Fatal("expected before snapshot")
	}
	assertSnapshot(t, *change.Before, path, role, typ, pointer)
}

func assertAfterSnapshot(t *testing.T, change FieldChange, path string, role model.FieldRole, typ string, pointer string) {
	t.Helper()
	if change.After == nil {
		t.Fatal("expected after snapshot")
	}
	assertSnapshot(t, *change.After, path, role, typ, pointer)
}

func assertSnapshot(t *testing.T, snapshot FieldSnapshot, path string, role model.FieldRole, typ string, pointer string) {
	t.Helper()
	if snapshot.Path != path {
		t.Fatalf("expected snapshot path %q, got %q", path, snapshot.Path)
	}
	if snapshot.Role != role {
		t.Fatalf("expected snapshot role %q, got %q", role, snapshot.Role)
	}
	if snapshot.Type != typ {
		t.Fatalf("expected snapshot type %q, got %q", typ, snapshot.Type)
	}
	if snapshot.JSONPointer != pointer {
		t.Fatalf("expected snapshot JSON pointer %q, got %q", pointer, snapshot.JSONPointer)
	}
}

func assertDuplicateResourceError(t *testing.T, err error, resource ResourceID) {
	t.Helper()
	if err == nil {
		t.Fatal("expected duplicate resource error")
	}
	var duplicate DuplicateResourceError
	if !errors.As(err, &duplicate) {
		t.Fatalf("expected DuplicateResourceError, got %T: %v", err, err)
	}
	if duplicate.Resource != resource {
		t.Fatalf("expected duplicate resource %#v, got %#v", resource, duplicate.Resource)
	}
	if !strings.Contains(err.Error(), "duplicate resource identity") {
		t.Fatalf("expected duplicate resource message, got %q", err.Error())
	}
}

func assertHasChange(t *testing.T, changes []FieldChange, kind ChangeKind, path string, role model.FieldRole) {
	t.Helper()
	for _, change := range changes {
		if change.Kind == kind && change.Field.Path == path && change.Field.Role == role {
			return
		}
	}
	t.Fatalf("expected %s change for %s/%s in %#v", kind, path, role, changes)
}

func assertNoChange(t *testing.T, changes []FieldChange, kind ChangeKind, path string) {
	t.Helper()
	for _, change := range changes {
		if change.Kind == kind && change.Field.Path == path {
			t.Fatalf("unexpected %s change for %s in %#v", kind, path, changes)
		}
	}
}

func corpusWithMapping(file, pointer string, fields ...model.Field) model.Corpus {
	return model.Corpus{Mappings: []model.Mapping{mapping(file, pointer, fields...)}}
}

func mapping(file, pointer string, fields ...model.Field) model.Mapping {
	return model.Mapping{
		Source:        source(file),
		JSONPointer:   pointer,
		Properties:    properties(fields),
		RuntimeFields: runtimeFields(fields),
	}
}

func corpusWithIndexTemplate(file string, fields ...model.Field) model.Corpus {
	return model.Corpus{
		IndexTemplates: []model.IndexTemplate{
			indexTemplate(file, fields...),
		},
	}
}

func indexTemplate(file string, fields ...model.Field) model.IndexTemplate {
	mapping := mapping(file, "/template/mappings", fields...)
	return model.IndexTemplate{
		Source: source(file),
		Template: model.TemplateBody{
			Mappings: &mapping,
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

func properties(fields []model.Field) []model.Field {
	var filtered []model.Field
	for _, field := range fields {
		if isRuntimeField(field) {
			continue
		}
		filtered = append(filtered, field)
	}
	return filtered
}

func runtimeFields(fields []model.Field) []model.Field {
	var filtered []model.Field
	for _, field := range fields {
		if isRuntimeField(field) {
			filtered = append(filtered, field)
		}
	}
	return filtered
}

func isRuntimeField(field model.Field) bool {
	return strings.HasPrefix(field.JSONPointer, "/runtime/")
}

func property(path string, typ string) model.Field {
	return model.Field{
		Name:        fieldName(path),
		Path:        path,
		Type:        typ,
		Source:      source("mapping.json"),
		JSONPointer: "/properties/" + path,
	}
}

func runtimeField(path string, typ string) model.Field {
	return model.Field{
		Name:        fieldName(path),
		Path:        path,
		Type:        typ,
		Source:      source("mapping.json"),
		JSONPointer: "/runtime/" + path,
	}
}

func textWithMultiField(path string, multiName string, multiType string) model.Field {
	parent := property(path, "text")
	parent.Fields = []model.Field{
		{
			Name:        multiName,
			Path:        path + "." + multiName,
			Type:        multiType,
			Source:      source("mapping.json"),
			JSONPointer: "/properties/" + path + "/fields/" + multiName,
		},
	}
	return parent
}

func fieldName(path string) string {
	index := strings.LastIndexByte(path, '.')
	if index == -1 {
		return path
	}
	return path[index+1:]
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
