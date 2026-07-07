package rules

import "github.com/marcinbohm/search-index-preflight/internal/model"

type Rule interface {
	Metadata() Metadata
	Check(ctx Context, corpus model.Corpus) ([]model.Finding, error)
}
