package normalizer

import "github.com/marcinbohm/search-index-lint/internal/model"

func NormalizeIndexTemplate(document model.RawDocument) model.IndexTemplate {
	root, ok := asObject(document.Content)
	if !ok {
		return model.IndexTemplate{
			Source:      document.Source,
			Diagnostics: []model.Diagnostic{diagnostic(document.Source, "index template document must be a JSON object")},
		}
	}

	template := model.IndexTemplate{
		Source:        document.Source,
		IndexPatterns: stringSlice(root["index_patterns"]),
		Priority:      intPointer(root["priority"]),
		ComposedOf:    stringSlice(root["composed_of"]),
		DataStream:    hasKey(root, "data_stream"),
		Meta:          objectMap(root["_meta"]),
	}
	template.Template = normalizeTemplateBody(document.Source, root["template"], "/template", &template.Diagnostics)
	return template
}

func NormalizeComponentTemplate(document model.RawDocument) model.ComponentTemplate {
	root, ok := asObject(document.Content)
	if !ok {
		return model.ComponentTemplate{
			Source:      document.Source,
			Diagnostics: []model.Diagnostic{diagnostic(document.Source, "component template document must be a JSON object")},
		}
	}

	template := model.ComponentTemplate{
		Source:  document.Source,
		Version: intPointer(root["version"]),
		Meta:    objectMap(root["_meta"]),
	}
	template.Template = normalizeTemplateBody(document.Source, root["template"], "/template", &template.Diagnostics)
	return template
}

func normalizeTemplateBody(source model.Source, value any, pointer string, diagnostics *[]model.Diagnostic) model.TemplateBody {
	var body model.TemplateBody
	if value == nil {
		return body
	}
	templateObject, ok := asObject(value)
	if !ok {
		*diagnostics = append(*diagnostics, diagnosticAt(source, pointer, "template must be an object"))
		return body
	}

	body.Settings = objectMap(templateObject["settings"])
	body.Aliases = objectMap(templateObject["aliases"])
	if mappingsObject, ok := asObject(templateObject["mappings"]); ok {
		mapping := normalizeMappingObject(source, mappingsObject, pointerJoin(pointer, "mappings"))
		body.Mappings = &mapping
		*diagnostics = append(*diagnostics, mapping.Diagnostics...)
	} else if _, exists := templateObject["mappings"]; exists {
		*diagnostics = append(*diagnostics, diagnosticAt(source, pointerJoin(pointer, "mappings"), "mappings must be an object"))
	}
	return body
}

func hasKey(object map[string]any, key string) bool {
	_, ok := object[key]
	return ok
}
