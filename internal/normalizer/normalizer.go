package normalizer

import (
	"fmt"

	"github.com/marcinbohm/search-index-preflight/internal/model"
)

func Normalize(documents []model.RawDocument) model.Corpus {
	var corpus model.Corpus
	for _, document := range documents {
		if len(document.Diagnostics) > 0 {
			continue
		}
		switch detectKind(document) {
		case model.DocumentKindMapping:
			mapping := NormalizeMapping(document)
			corpus.Mappings = append(corpus.Mappings, mapping)
			corpus.Diagnostics = append(corpus.Diagnostics, mapping.Diagnostics...)
		case model.DocumentKindIndexTemplate:
			template := NormalizeIndexTemplate(document)
			corpus.IndexTemplates = append(corpus.IndexTemplates, template)
			corpus.Diagnostics = append(corpus.Diagnostics, template.Diagnostics...)
		case model.DocumentKindComponentTemplate:
			template := NormalizeComponentTemplate(document)
			corpus.ComponentTemplates = append(corpus.ComponentTemplates, template)
			corpus.Diagnostics = append(corpus.Diagnostics, template.Diagnostics...)
		case model.DocumentKindSampleDocs:
			corpus.SampleDocuments = append(corpus.SampleDocuments, document)
		default:
			corpus.Diagnostics = append(corpus.Diagnostics, diagnostic(document.Source, "unknown JSON document kind"))
		}
	}
	return corpus
}

func detectKind(document model.RawDocument) model.DocumentKind {
	if document.Kind != model.DocumentKindUnknown {
		return document.Kind
	}

	root, ok := asObject(document.Content)
	if !ok {
		return model.DocumentKindUnknown
	}
	if _, ok := root["index_patterns"]; ok {
		return model.DocumentKindIndexTemplate
	}
	if _, ok := root["template"]; ok {
		return model.DocumentKindComponentTemplate
	}
	if _, ok := root["mappings"]; ok {
		return model.DocumentKindMapping
	}
	if _, ok := root["properties"]; ok {
		return model.DocumentKindMapping
	}
	return model.DocumentKindUnknown
}

func diagnostic(source model.Source, message string) model.Diagnostic {
	return model.Diagnostic{
		Severity: model.SeverityError,
		File:     source.RelativePath,
		Message:  message,
	}
}

func diagnosticAt(source model.Source, pointer, message string) model.Diagnostic {
	if pointer == "" {
		return diagnostic(source, message)
	}
	return diagnostic(source, fmt.Sprintf("%s: %s", pointer, message))
}

func asObject(value any) (map[string]any, bool) {
	object, ok := value.(map[string]any)
	return object, ok
}
