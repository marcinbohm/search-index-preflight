package rules

import (
	"fmt"

	"github.com/marcinbohm/search-index-preflight/internal/model"
)

type sil002RootDynamicEnabled struct{}

func NewSIL002() Rule {
	return sil002RootDynamicEnabled{}
}

func (r sil002RootDynamicEnabled) Metadata() Metadata {
	return Metadata{
		ID:          "SIL002",
		Name:        "root-dynamic-enabled",
		Category:    "dynamic-mapping",
		Description: "Detects mappings/templates where root-level dynamic mapping is explicitly enabled.",
		Severity:    model.SeverityWarning,
		Confidence:  model.ConfidenceMedium,
		Determinism: model.DeterminismHeuristic,
	}
}

func (r sil002RootDynamicEnabled) Check(ctx Context, corpus model.Corpus) ([]model.Finding, error) {
	var findings []model.Finding
	for _, mapping := range corpus.Mappings {
		if finding, ok := r.checkMapping(mapping, model.AppendJSONPointer(mapping.JSONPointer, "dynamic"), "Root dynamic mapping is explicitly enabled. Unexpected fields may expand the mapping."); ok {
			findings = append(findings, finding)
		}
	}
	for _, template := range corpus.IndexTemplates {
		if template.Template.Mappings == nil {
			continue
		}
		if finding, ok := r.checkMapping(*template.Template.Mappings, "/template/mappings/dynamic", "Root dynamic mapping is explicitly enabled in this template mapping. New fields may be added to indices created from this template."); ok {
			findings = append(findings, finding)
		}
	}
	for _, template := range corpus.ComponentTemplates {
		if template.Template.Mappings == nil {
			continue
		}
		if finding, ok := r.checkMapping(*template.Template.Mappings, "/template/mappings/dynamic", "Root dynamic mapping is explicitly enabled in this template mapping. New fields may be added to indices created from this template."); ok {
			findings = append(findings, finding)
		}
	}
	return findings, nil
}

func (r sil002RootDynamicEnabled) checkMapping(mapping model.Mapping, pointer string, message string) (model.Finding, bool) {
	if mapping.Dynamic != model.DynamicSettingTrue {
		return model.Finding{}, false
	}

	metadata := r.Metadata()
	return model.Finding{
		ID:          metadata.ID,
		Name:        metadata.Name,
		Severity:    metadata.Severity,
		Confidence:  metadata.Confidence,
		Category:    metadata.Category,
		Determinism: metadata.Determinism,
		File:        mapping.Source.RelativePath,
		JSONPointer: pointer,
		Message:     message,
		Remediation: "Use explicit mappings for known fields. Consider dynamic: strict or dynamic: false for controlled schemas, or scope dynamic behavior to known safe objects. Keep dynamic enabled only when the expansion risk is intentional and reviewed.",
		Fingerprint: fmt.Sprintf("%s:%s:%s", metadata.ID, mapping.Source.RelativePath, pointer),
	}, true
}
