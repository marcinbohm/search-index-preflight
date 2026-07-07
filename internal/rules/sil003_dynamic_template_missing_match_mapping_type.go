package rules

import (
	"fmt"

	"github.com/marcinbohm/search-index-preflight/internal/model"
)

type sil003DynamicTemplateMissingMatchMappingType struct{}

func NewSIL003() Rule {
	return sil003DynamicTemplateMissingMatchMappingType{}
}

func (r sil003DynamicTemplateMissingMatchMappingType) Metadata() Metadata {
	return Metadata{
		ID:          "SIL003",
		Name:        "dynamic-template-missing-match-mapping-type",
		Category:    "dynamic-templates",
		Description: "Detects dynamic templates that omit match_mapping_type, which may make the template apply more broadly than intended.",
		Severity:    model.SeverityWarning,
		Confidence:  model.ConfidenceMedium,
		Determinism: model.DeterminismHeuristic,
	}
}

func (r sil003DynamicTemplateMissingMatchMappingType) Check(ctx Context, corpus model.Corpus) ([]model.Finding, error) {
	var findings []model.Finding
	for _, mapping := range corpus.Mappings {
		findings = append(findings, r.checkMapping(mapping)...)
	}
	for _, template := range corpus.IndexTemplates {
		if template.Template.Mappings == nil {
			continue
		}
		findings = append(findings, r.checkMapping(*template.Template.Mappings)...)
	}
	for _, template := range corpus.ComponentTemplates {
		if template.Template.Mappings == nil {
			continue
		}
		findings = append(findings, r.checkMapping(*template.Template.Mappings)...)
	}
	return findings, nil
}

func (r sil003DynamicTemplateMissingMatchMappingType) checkMapping(mapping model.Mapping) []model.Finding {
	var findings []model.Finding
	for _, template := range mapping.DynamicTemplates {
		if template.HasMatchMappingType {
			continue
		}
		findings = append(findings, r.finding(template))
	}
	return findings
}

func (r sil003DynamicTemplateMissingMatchMappingType) finding(template model.DynamicTemplate) model.Finding {
	metadata := r.Metadata()
	message := "A dynamic template does not declare match_mapping_type. It may apply more broadly than intended."
	if template.Name != "" {
		message = fmt.Sprintf("Dynamic template %q does not declare match_mapping_type. It may apply more broadly than intended.", template.Name)
	}

	return model.Finding{
		ID:          metadata.ID,
		Name:        metadata.Name,
		Severity:    metadata.Severity,
		Confidence:  metadata.Confidence,
		Category:    metadata.Category,
		Determinism: metadata.Determinism,
		File:        template.Source.RelativePath,
		JSONPointer: template.JSONPointer,
		Message:     message,
		Remediation: "Add match_mapping_type when the template is intended for a specific detected field type, or document why broad matching is intentional.",
		Fingerprint: fmt.Sprintf("%s:%s:%s", metadata.ID, template.Source.RelativePath, template.JSONPointer),
	}
}
