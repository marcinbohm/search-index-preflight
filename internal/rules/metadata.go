package rules

import "github.com/marcinbohm/search-index-preflight/internal/model"

type Metadata struct {
	ID          string
	Name        string
	Category    string
	Description string
	Severity    model.Severity
	Confidence  model.Confidence
	Determinism model.Determinism
}
