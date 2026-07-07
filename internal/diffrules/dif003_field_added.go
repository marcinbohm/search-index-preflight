package diffrules

import (
	"fmt"

	"github.com/marcinbohm/search-index-preflight/internal/diff"
	"github.com/marcinbohm/search-index-preflight/internal/model"
)

type dif003FieldAdded struct{}

func NewDIF003() Rule {
	return dif003FieldAdded{}
}

func (r dif003FieldAdded) Metadata() Metadata {
	return Metadata{
		ID:          "DIF003",
		Name:        "field-added",
		Category:    "schema-diff",
		Description: "Detects fields that are present in the current schema corpus but were missing from the base schema corpus.",
		Severity:    model.SeverityInfo,
		Confidence:  model.ConfidenceHigh,
		Determinism: model.DeterminismDeterministic,
	}
}

func (r dif003FieldAdded) Check(ctx Context, result diff.Result) ([]model.Finding, error) {
	var findings []model.Finding
	for _, change := range result.FieldChanges {
		if change.Kind != diff.ChangeFieldAdded {
			continue
		}
		findings = append(findings, r.finding(change))
	}
	return findings, nil
}

func (r dif003FieldAdded) finding(change diff.FieldChange) model.Finding {
	metadata := r.Metadata()
	pointer := addedFieldPointer(change)

	return model.Finding{
		ID:          metadata.ID,
		Name:        metadata.Name,
		Severity:    metadata.Severity,
		Confidence:  metadata.Confidence,
		Category:    metadata.Category,
		Determinism: metadata.Determinism,
		File:        change.Resource.File,
		JSONPointer: pointer,
		Message:     fmt.Sprintf("Field %q was added to the current schema.", change.Field.Path),
		Remediation: "Verify that the new field is intentional, explicitly mapped as expected, and does not conflict with field-count limits or downstream consumers.",
		Fingerprint: fmt.Sprintf("%s:%s:%s:%s:%s:%s", metadata.ID, change.Resource.Kind, change.Resource.File, pointer, change.Field.Path, change.Field.Role),
	}
}

func addedFieldPointer(change diff.FieldChange) string {
	if change.After != nil && change.After.JSONPointer != "" {
		return change.After.JSONPointer
	}
	return change.Resource.JSONPointer
}
