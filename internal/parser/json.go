package parser

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/marcinbohm/search-index-preflight/internal/model"
)

func ParseJSON(source model.Source, kind model.DocumentKind, content []byte) model.RawDocument {
	document := model.RawDocument{
		Kind:   kind,
		Source: source,
	}

	var parsed any
	decoder := json.NewDecoder(bytes.NewReader(content))
	decoder.UseNumber()
	if err := decoder.Decode(&parsed); err != nil {
		document.Diagnostics = []model.Diagnostic{errorDiagnostic(source, 0, fmt.Sprintf("invalid JSON: %v", err))}
		return document
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		document.Diagnostics = []model.Diagnostic{errorDiagnostic(source, 0, "invalid JSON: multiple JSON values found")}
		return document
	}

	document.Content = parsed
	return document
}
