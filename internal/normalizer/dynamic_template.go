package normalizer

import (
	"strconv"

	"github.com/marcinbohm/search-index-preflight/internal/model"
)

func normalizeDynamicTemplates(source model.Source, value any, pointer string) ([]model.DynamicTemplate, []model.Diagnostic) {
	if value == nil {
		return nil, nil
	}
	items, ok := value.([]any)
	if !ok {
		return nil, []model.Diagnostic{diagnosticAt(source, pointer, "dynamic_templates must be an array")}
	}

	templates := make([]model.DynamicTemplate, 0, len(items))
	var diagnostics []model.Diagnostic
	for index, item := range items {
		itemPointer := pointerJoin(pointer, indexString(index))
		templateObject, ok := asObject(item)
		if !ok {
			diagnostics = append(diagnostics, diagnosticAt(source, itemPointer, "dynamic template entry must be an object"))
			continue
		}
		for _, name := range sortedKeys(templateObject) {
			definitionPointer := pointerJoin(itemPointer, name)
			definition, ok := asObject(templateObject[name])
			if !ok {
				diagnostics = append(diagnostics, diagnosticAt(source, definitionPointer, "dynamic template definition must be an object"))
				continue
			}

			template := model.DynamicTemplate{
				Name:                name,
				Source:              source,
				JSONPointer:         definitionPointer,
				Match:               stringValue(definition["match"]),
				Unmatch:             stringValue(definition["unmatch"]),
				PathMatch:           stringValue(definition["path_match"]),
				PathUnmatch:         stringValue(definition["path_unmatch"]),
				MatchMappingType:    stringValue(definition["match_mapping_type"]),
				HasMatchMappingType: hasKey(definition, "match_mapping_type"),
				Mapping:             objectMap(definition["mapping"]),
			}
			templates = append(templates, template)
		}
	}
	return templates, diagnostics
}

func stringValue(value any) string {
	text, ok := value.(string)
	if !ok {
		return ""
	}
	return text
}

func indexString(index int) string {
	return strconv.Itoa(index)
}
