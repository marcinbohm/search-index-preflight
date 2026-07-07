package rules

import (
	"fmt"

	"github.com/marcinbohm/search-index-preflight/internal/model"
)

const (
	sil001DefaultTotalFieldsLimit = 1000
	sil001WarningThreshold        = 800
)

type sil001TotalFieldsLimit struct{}

func NewSIL001() Rule {
	return sil001TotalFieldsLimit{}
}

func (r sil001TotalFieldsLimit) Metadata() Metadata {
	return Metadata{
		ID:          "SIL001",
		Name:        "total-fields-limit-risk",
		Category:    "mapping-limits",
		Description: "Detects explicit mappings/templates whose normalized field count approaches or exceeds the default total fields limit.",
		Severity:    model.SeverityWarning,
		Confidence:  model.ConfidenceHigh,
		Determinism: model.DeterminismDeterministic,
	}
}

func (r sil001TotalFieldsLimit) Check(ctx Context, corpus model.Corpus) ([]model.Finding, error) {
	var findings []model.Finding
	for _, mapping := range corpus.Mappings {
		if finding, ok := r.checkMapping(mapping, mapping.JSONPointer); ok {
			findings = append(findings, finding)
		}
	}
	for _, template := range corpus.IndexTemplates {
		if template.Template.Mappings == nil {
			continue
		}
		if finding, ok := r.checkMapping(*template.Template.Mappings, "/template/mappings"); ok {
			findings = append(findings, finding)
		}
	}
	for _, template := range corpus.ComponentTemplates {
		if template.Template.Mappings == nil {
			continue
		}
		if finding, ok := r.checkMapping(*template.Template.Mappings, "/template/mappings"); ok {
			findings = append(findings, finding)
		}
	}
	return findings, nil
}

func (r sil001TotalFieldsLimit) checkMapping(mapping model.Mapping, pointer string) (model.Finding, bool) {
	stats := model.CountMappingFields(mapping)
	if stats.TotalFields < sil001WarningThreshold {
		return model.Finding{}, false
	}

	metadata := r.Metadata()
	severity := model.SeverityWarning
	message := fmt.Sprintf("Mapping has %d normalized fields, approaching the default total fields limit of %d.", stats.TotalFields, sil001DefaultTotalFieldsLimit)
	if stats.TotalFields >= sil001DefaultTotalFieldsLimit {
		severity = model.SeverityError
		message = fmt.Sprintf("Mapping has %d normalized fields, exceeding the default total fields limit of %d.", stats.TotalFields, sil001DefaultTotalFieldsLimit)
	}

	return model.Finding{
		ID:          metadata.ID,
		Name:        metadata.Name,
		Severity:    severity,
		Confidence:  metadata.Confidence,
		Category:    metadata.Category,
		Determinism: metadata.Determinism,
		File:        mapping.Source.RelativePath,
		JSONPointer: pointer,
		Message:     message,
		Remediation: "Reduce explicit field count, restrict dynamic mappings, consider flattened/flat_object only when query semantics fit, split unrelated data into separate indices, or raise index.mapping.total_fields.limit only with operational review.",
		Fingerprint: fmt.Sprintf("%s:%s:%s", metadata.ID, mapping.Source.RelativePath, pointer),
	}, true
}
