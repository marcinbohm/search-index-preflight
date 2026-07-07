package parser

import "github.com/marcinbohm/search-index-preflight/internal/model"

func Parse(source model.Source, kind model.DocumentKind, content []byte) model.RawDocument {
	// TODO: Add YAML parsing when config/schema YAML support is implemented.
	if kind == model.DocumentKindSampleDocs {
		return ParseJSONL(source, content)
	}
	return ParseJSON(source, kind, content)
}

func errorDiagnostic(source model.Source, line int, message string) model.Diagnostic {
	return model.Diagnostic{
		Severity: model.SeverityError,
		File:     source.RelativePath,
		Line:     line,
		Message:  message,
	}
}
