package normalizer

import (
	"encoding/json"
	"sort"
	"strings"

	"github.com/marcinbohm/search-index-preflight/internal/model"
)

func NormalizeMapping(document model.RawDocument) model.Mapping {
	root, ok := asObject(document.Content)
	if !ok {
		return model.Mapping{
			Source:      document.Source,
			Dynamic:     model.DynamicSettingUnspecified,
			Diagnostics: []model.Diagnostic{diagnostic(document.Source, "mapping document must be a JSON object")},
		}
	}

	if wrapped, ok := asObject(root["mappings"]); ok {
		return normalizeMappingObject(document.Source, wrapped, "/mappings")
	}
	return normalizeMappingObject(document.Source, root, "")
}

func normalizeMappingObject(source model.Source, object map[string]any, pointer string) model.Mapping {
	mapping := model.Mapping{
		Source:      source,
		JSONPointer: pointer,
		Dynamic:     parseDynamicSetting(object["dynamic"]),
	}
	mapping.DateDetection = boolPointer(object["date_detection"])
	mapping.NumericDetection = boolPointer(object["numeric_detection"])
	mapping.Meta = objectMap(object["_meta"])

	if properties, ok := asObject(object["properties"]); ok {
		mapping.Properties = normalizeFields(source, properties, "", pointerJoin(pointer, "properties"))
	} else if _, exists := object["properties"]; exists {
		mapping.Diagnostics = append(mapping.Diagnostics, diagnosticAt(source, pointerJoin(pointer, "properties"), "properties must be an object"))
	}

	if runtimeFields, ok := asObject(object["runtime"]); ok {
		mapping.RuntimeFields = normalizeFields(source, runtimeFields, "", pointerJoin(pointer, "runtime"))
	} else if _, exists := object["runtime"]; exists {
		mapping.Diagnostics = append(mapping.Diagnostics, diagnosticAt(source, pointerJoin(pointer, "runtime"), "runtime must be an object"))
	}

	templates, diagnostics := normalizeDynamicTemplates(source, object["dynamic_templates"], pointerJoin(pointer, "dynamic_templates"))
	mapping.DynamicTemplates = templates
	mapping.Diagnostics = append(mapping.Diagnostics, diagnostics...)
	return mapping
}

func normalizeFields(source model.Source, fields map[string]any, parentPath, pointer string) []model.Field {
	names := sortedKeys(fields)
	normalized := make([]model.Field, 0, len(names))
	for _, name := range names {
		fieldPointer := pointerJoin(pointer, name)
		field := model.Field{
			Name:        name,
			Path:        joinPath(parentPath, name),
			ParentPath:  parentPath,
			Source:      source,
			JSONPointer: fieldPointer,
			Dynamic:     model.DynamicSettingUnspecified,
		}

		object, ok := asObject(fields[name])
		if !ok {
			field.Parameters = map[string]any{"value": fields[name]}
			normalized = append(normalized, field)
			continue
		}

		if value, ok := object["type"].(string); ok {
			field.Type = value
		}
		field.Dynamic = parseDynamicSetting(object["dynamic"])
		field.Enabled = boolPointer(object["enabled"])
		field.Parameters = fieldParameters(object)

		if properties, ok := asObject(object["properties"]); ok {
			field.Properties = normalizeFields(source, properties, field.Path, pointerJoin(fieldPointer, "properties"))
		}
		if multiFields, ok := asObject(object["fields"]); ok {
			field.Fields = normalizeFields(source, multiFields, field.Path, pointerJoin(fieldPointer, "fields"))
		}

		normalized = append(normalized, field)
	}
	return normalized
}

func fieldParameters(object map[string]any) map[string]any {
	parameters := make(map[string]any)
	for key, value := range object {
		switch key {
		case "type", "properties", "fields", "dynamic", "enabled":
			continue
		default:
			parameters[key] = value
		}
	}
	if len(parameters) == 0 {
		return nil
	}
	return parameters
}

func parseDynamicSetting(value any) model.DynamicSetting {
	switch typed := value.(type) {
	case bool:
		if typed {
			return model.DynamicSettingTrue
		}
		return model.DynamicSettingFalse
	case string:
		switch strings.ToLower(typed) {
		case "true":
			return model.DynamicSettingTrue
		case "false":
			return model.DynamicSettingFalse
		case "strict":
			return model.DynamicSettingStrict
		case "runtime":
			return model.DynamicSettingRuntime
		}
	}
	return model.DynamicSettingUnspecified
}

func boolPointer(value any) *bool {
	typed, ok := value.(bool)
	if !ok {
		return nil
	}
	return &typed
}

func objectMap(value any) map[string]any {
	object, ok := asObject(value)
	if !ok || len(object) == 0 {
		return nil
	}
	return object
}

func stringSlice(value any) []string {
	switch typed := value.(type) {
	case string:
		return []string{typed}
	case []any:
		values := make([]string, 0, len(typed))
		for _, item := range typed {
			if text, ok := item.(string); ok {
				values = append(values, text)
			}
		}
		return values
	default:
		return nil
	}
}

func intPointer(value any) *int {
	switch typed := value.(type) {
	case json.Number:
		number, err := typed.Int64()
		if err != nil {
			return nil
		}
		converted := int(number)
		return &converted
	case float64:
		converted := int(typed)
		if float64(converted) != typed {
			return nil
		}
		return &converted
	default:
		return nil
	}
}

func sortedKeys(values map[string]any) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func joinPath(parentPath, name string) string {
	if parentPath == "" {
		return name
	}
	return parentPath + "." + name
}

func pointerJoin(base string, segments ...string) string {
	return model.AppendJSONPointer(base, segments...)
}
