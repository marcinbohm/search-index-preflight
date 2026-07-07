package diff

import (
	"sort"

	"github.com/marcinbohm/search-index-preflight/internal/model"
)

func collectResources(corpus model.Corpus) map[ResourceID]map[FieldID]FieldSnapshot {
	resources := make(map[ResourceID]map[FieldID]FieldSnapshot)

	for _, mapping := range corpus.Mappings {
		resourceID := ResourceID{
			Kind:        ResourceMapping,
			File:        sourceFile(mapping.Source),
			JSONPointer: mapping.JSONPointer,
		}
		resources[resourceID] = collectMappingFields(mapping)
	}

	for _, template := range corpus.IndexTemplates {
		if template.Template.Mappings == nil {
			continue
		}
		resourceID := ResourceID{
			Kind:        ResourceIndexTemplate,
			File:        sourceFile(template.Source),
			JSONPointer: template.Template.Mappings.JSONPointer,
		}
		resources[resourceID] = collectMappingFields(*template.Template.Mappings)
	}

	for _, template := range corpus.ComponentTemplates {
		if template.Template.Mappings == nil {
			continue
		}
		resourceID := ResourceID{
			Kind:        ResourceComponentTemplate,
			File:        sourceFile(template.Source),
			JSONPointer: template.Template.Mappings.JSONPointer,
		}
		resources[resourceID] = collectMappingFields(*template.Template.Mappings)
	}

	return resources
}

func collectMappingFields(mapping model.Mapping) map[FieldID]FieldSnapshot {
	corpus := model.Corpus{Mappings: []model.Mapping{mapping}}
	fields := make(map[FieldID]FieldSnapshot)
	for _, visit := range model.CollectFields(corpus) {
		fieldID := FieldID{Path: visit.Path, Role: visit.Role}
		fields[fieldID] = FieldSnapshot{
			Path:        visit.Path,
			Role:        visit.Role,
			Type:        visit.Field.Type,
			JSONPointer: visit.JSONPointer,
		}
	}
	return fields
}

func sourceFile(source model.Source) string {
	if source.RelativePath != "" {
		return source.RelativePath
	}
	return source.Path
}

func sortFieldChanges(changes []FieldChange) {
	sort.Slice(changes, func(i, j int) bool {
		left := changes[i]
		right := changes[j]

		if left.Resource.Kind != right.Resource.Kind {
			return left.Resource.Kind < right.Resource.Kind
		}
		if left.Resource.File != right.Resource.File {
			return left.Resource.File < right.Resource.File
		}
		if left.Resource.JSONPointer != right.Resource.JSONPointer {
			return left.Resource.JSONPointer < right.Resource.JSONPointer
		}
		if left.Field.Path != right.Field.Path {
			return left.Field.Path < right.Field.Path
		}
		if left.Field.Role != right.Field.Role {
			return left.Field.Role < right.Field.Role
		}
		return left.Kind < right.Kind
	})
}
