package rules

import (
	"fmt"

	"github.com/marcinbohm/search-index-preflight/internal/model"
)

type RunRequest struct {
	Corpus model.Corpus
}

type RunResult struct {
	Findings    []model.Finding
	Diagnostics []model.Diagnostic
}

func Run(ctx Context, registry *Registry, request RunRequest) (RunResult, error) {
	if registry == nil {
		return RunResult{}, fmt.Errorf("registry is nil")
	}

	var result RunResult
	for _, rule := range registry.List() {
		findings, err := rule.Check(ctx, request.Corpus)
		if err != nil {
			// Rule execution errors are returned for now. Diagnostic policy should be
			// decided when real rule execution is wired into the CLI.
			return result, fmt.Errorf("rule %s: %w", rule.Metadata().ID, err)
		}
		result.Findings = append(result.Findings, findings...)
	}
	return result, nil
}
