package model

import "strings"

func AppendJSONPointer(base string, tokens ...string) string {
	parts := make([]string, 0, len(tokens)+1)
	if base != "" {
		parts = append(parts, strings.TrimPrefix(base, "/"))
	}
	for _, token := range tokens {
		parts = append(parts, escapeJSONPointerToken(token))
	}
	return "/" + strings.Join(parts, "/")
}

func escapeJSONPointerToken(token string) string {
	token = strings.ReplaceAll(token, "~", "~0")
	return strings.ReplaceAll(token, "/", "~1")
}
