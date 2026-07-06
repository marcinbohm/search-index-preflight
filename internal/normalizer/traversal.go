package normalizer

import "github.com/marcinbohm/search-index-lint/internal/model"

func WalkFields(corpus Corpus, visit func(model.FieldVisit)) {
	for _, mapping := range corpus.Mappings {
		walkMappingFields(mapping, fieldVisitContext{
			origin:        model.FieldOriginMapping,
			mappingSource: mapping.Source,
		}, visit)
	}
	for _, template := range corpus.IndexTemplates {
		if template.Template.Mappings == nil {
			continue
		}
		walkMappingFields(*template.Template.Mappings, fieldVisitContext{
			origin:            model.FieldOriginIndexTemplate,
			mappingSource:     template.Template.Mappings.Source,
			indexTemplateName: template.Name,
		}, visit)
	}
	for _, template := range corpus.ComponentTemplates {
		if template.Template.Mappings == nil {
			continue
		}
		walkMappingFields(*template.Template.Mappings, fieldVisitContext{
			origin:                model.FieldOriginComponentTemplate,
			mappingSource:         template.Template.Mappings.Source,
			componentTemplateName: template.Name,
		}, visit)
	}
}

func CollectFields(corpus Corpus) []model.FieldVisit {
	var visits []model.FieldVisit
	WalkFields(corpus, func(visit model.FieldVisit) {
		visits = append(visits, visit)
	})
	return visits
}

type fieldVisitContext struct {
	origin                model.FieldOrigin
	mappingSource         model.Source
	indexTemplateName     string
	componentTemplateName string
}

func walkMappingFields(mapping model.Mapping, context fieldVisitContext, visit func(model.FieldVisit)) {
	walkFields(mapping.Properties, context, model.FieldRoleProperty, visit)
	walkFields(mapping.RuntimeFields, context, model.FieldRoleRuntimeField, visit)
}

func walkFields(fields []model.Field, context fieldVisitContext, role model.FieldRole, visit func(model.FieldVisit)) {
	for _, field := range fields {
		visit(model.FieldVisit{
			Origin:                context.origin,
			Role:                  role,
			Source:                field.Source,
			MappingSource:         context.mappingSource,
			IndexTemplateName:     context.indexTemplateName,
			ComponentTemplateName: context.componentTemplateName,
			Field:                 field,
			Path:                  field.Path,
			JSONPointer:           field.JSONPointer,
		})
		walkFields(field.Properties, context, model.FieldRoleProperty, visit)
		walkFields(field.Fields, context, model.FieldRoleMultiField, visit)
	}
}
