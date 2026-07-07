package parser

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/marcinbohm/search-index-preflight/internal/model"
)

func ParseJSONL(source model.Source, content []byte) model.RawDocument {
	document := model.RawDocument{
		Kind:   model.DocumentKindSampleDocs,
		Source: source,
	}

	var documents []any
	scanner := bufio.NewScanner(bytes.NewReader(content))
	line := 0
	for scanner.Scan() {
		line++
		text := scanner.Text()
		if strings.TrimSpace(text) == "" {
			continue
		}

		var parsed any
		decoder := json.NewDecoder(strings.NewReader(text))
		decoder.UseNumber()
		if err := decoder.Decode(&parsed); err != nil {
			document.Diagnostics = append(document.Diagnostics, errorDiagnostic(source, line, fmt.Sprintf("invalid JSONL: %v", err)))
			continue
		}
		if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
			document.Diagnostics = append(document.Diagnostics, errorDiagnostic(source, line, "invalid JSONL: multiple JSON values found"))
			continue
		}
		documents = append(documents, parsed)
	}
	if err := scanner.Err(); err != nil {
		document.Diagnostics = append(document.Diagnostics, errorDiagnostic(source, line, fmt.Sprintf("invalid JSONL: %v", err)))
	}

	if len(document.Diagnostics) == 0 {
		document.Content = documents
	}
	return document
}
