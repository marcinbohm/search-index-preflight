package input

import "github.com/marcinbohm/search-index-preflight/internal/model"

type Source struct {
	Path         string
	RelativePath string
	Kind         model.DocumentKind
	Content      []byte
}
