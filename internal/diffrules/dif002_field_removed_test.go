package diffrules

import (
	"strings"
	"testing"

	"github.com/marcinbohm/search-index-preflight/internal/diff"
	"github.com/marcinbohm/search-index-preflight/internal/model"
)

func TestDIF002EmitsFindingForFieldRemoved(t *testing.T) {
	findings := runDIF002(t, diff.Result{
		FieldChanges: []diff.FieldChange{
			fieldRemoved("legacy_id", model.FieldRoleProperty, "keyword", "/properties/legacy_id"),
		},
	})

	if len(findings) != 1 {
		t.Fatalf("expected one finding, got %#v", findings)
	}
	finding := findings[0]
	if finding.ID != "DIF002" {
		t.Fatalf("expected ID DIF002, got %q", finding.ID)
	}
	if finding.Name != "field-removed" {
		t.Fatalf("expected name field-removed, got %q", finding.Name)
	}
	if finding.Severity != model.SeverityWarning {
		t.Fatalf("expected severity warning, got %q", finding.Severity)
	}
	if finding.Confidence != model.ConfidenceHigh {
		t.Fatalf("expected confidence high, got %q", finding.Confidence)
	}
	if finding.Determinism != model.DeterminismDeterministic {
		t.Fatalf("expected deterministic finding, got %q", finding.Determinism)
	}
	if finding.Category != "schema-diff" {
		t.Fatalf("expected category schema-diff, got %q", finding.Category)
	}
	if finding.File != "mapping.json" {
		t.Fatalf("expected file mapping.json, got %q", finding.File)
	}
	if finding.JSONPointer != "/properties/legacy_id" {
		t.Fatalf("expected JSON pointer /properties/legacy_id, got %q", finding.JSONPointer)
	}
	if !strings.Contains(finding.Message, "legacy_id") {
		t.Fatalf("expected message %q to contain legacy_id", finding.Message)
	}
	if finding.Remediation == "" {
		t.Fatal("expected remediation")
	}
	expectedFingerprint := "DIF002:mapping:mapping.json:/properties/legacy_id:legacy_id:property"
	if finding.Fingerprint != expectedFingerprint {
		t.Fatalf("expected fingerprint %q, got %q", expectedFingerprint, finding.Fingerprint)
	}
}

func TestDIF002IgnoresFieldAdded(t *testing.T) {
	findings := runDIF002(t, diff.Result{
		FieldChanges: []diff.FieldChange{
			{
				Kind:     diff.ChangeFieldAdded,
				Resource: diff.ResourceID{Kind: diff.ResourceMapping, File: "mapping.json"},
				Field:    diff.FieldID{Path: "status", Role: model.FieldRoleProperty},
				After:    snapshot("status", model.FieldRoleProperty, "keyword", "/properties/status"),
			},
		},
	})

	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %#v", findings)
	}
}

func TestDIF002IgnoresFieldTypeChanged(t *testing.T) {
	findings := runDIF002(t, diff.Result{
		FieldChanges: []diff.FieldChange{
			fieldTypeChanged("status", model.FieldRoleProperty, "keyword", "long", "/properties/status", "/properties/status"),
		},
	})

	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %#v", findings)
	}
}

func TestDIF002UsesResourcePointerWhenBeforeHasNoPointer(t *testing.T) {
	findings := runDIF002(t, diff.Result{
		FieldChanges: []diff.FieldChange{
			{
				Kind:     diff.ChangeFieldRemoved,
				Resource: diff.ResourceID{Kind: diff.ResourceMapping, File: "mapping.json", JSONPointer: "/mappings"},
				Field:    diff.FieldID{Path: "legacy_id", Role: model.FieldRoleProperty},
				Before:   snapshot("legacy_id", model.FieldRoleProperty, "keyword", ""),
			},
		},
	})

	if len(findings) != 1 {
		t.Fatalf("expected one finding, got %#v", findings)
	}
	if findings[0].JSONPointer != "/mappings" {
		t.Fatalf("expected resource pointer fallback, got %q", findings[0].JSONPointer)
	}
}

func TestDIF002WorksForMultiFieldRemoved(t *testing.T) {
	findings := runDIF002(t, diff.Result{
		FieldChanges: []diff.FieldChange{
			fieldRemoved("message.keyword", model.FieldRoleMultiField, "keyword", "/properties/message/fields/keyword"),
		},
	})

	if len(findings) != 1 {
		t.Fatalf("expected one finding, got %#v", findings)
	}
	if !strings.Contains(findings[0].Message, "message.keyword") {
		t.Fatalf("expected message to contain multi-field path, got %q", findings[0].Message)
	}
	if findings[0].JSONPointer != "/properties/message/fields/keyword" {
		t.Fatalf("expected multi-field pointer, got %q", findings[0].JSONPointer)
	}
}

func TestDIF002WorksForRuntimeFieldRemoved(t *testing.T) {
	findings := runDIF002(t, diff.Result{
		FieldChanges: []diff.FieldChange{
			fieldRemoved("day_of_week", model.FieldRoleRuntimeField, "keyword", "/runtime/day_of_week"),
		},
	})

	if len(findings) != 1 {
		t.Fatalf("expected one finding, got %#v", findings)
	}
	if !strings.Contains(findings[0].Message, "day_of_week") {
		t.Fatalf("expected message to contain runtime field path, got %q", findings[0].Message)
	}
	if findings[0].JSONPointer != "/runtime/day_of_week" {
		t.Fatalf("expected runtime pointer, got %q", findings[0].JSONPointer)
	}
}

func TestDIF002WorksForIndexTemplateFieldRemoved(t *testing.T) {
	findings := runDIF002(t, diff.Result{
		FieldChanges: []diff.FieldChange{
			fieldRemovedInResource(diff.ResourceIndexTemplate, "template.json", "legacy_id", model.FieldRoleProperty, "keyword", "/template/mappings/properties/legacy_id"),
		},
	})

	if len(findings) != 1 {
		t.Fatalf("expected one finding, got %#v", findings)
	}
	if findings[0].File != "template.json" {
		t.Fatalf("expected template file, got %q", findings[0].File)
	}
	if findings[0].JSONPointer != "/template/mappings/properties/legacy_id" {
		t.Fatalf("expected field snapshot pointer, got %q", findings[0].JSONPointer)
	}
	if !strings.Contains(findings[0].Fingerprint, string(diff.ResourceIndexTemplate)) {
		t.Fatalf("expected fingerprint to include resource kind, got %q", findings[0].Fingerprint)
	}
	if !strings.Contains(findings[0].Message, "legacy_id") {
		t.Fatalf("expected message to contain field path, got %q", findings[0].Message)
	}
}

func TestDIF002WorksForComponentTemplateFieldRemoved(t *testing.T) {
	findings := runDIF002(t, diff.Result{
		FieldChanges: []diff.FieldChange{
			fieldRemovedInResource(diff.ResourceComponentTemplate, "component.json", "legacy_id", model.FieldRoleProperty, "keyword", "/template/mappings/properties/legacy_id"),
		},
	})

	if len(findings) != 1 {
		t.Fatalf("expected one finding, got %#v", findings)
	}
	if findings[0].File != "component.json" {
		t.Fatalf("expected component template file, got %q", findings[0].File)
	}
	if findings[0].JSONPointer != "/template/mappings/properties/legacy_id" {
		t.Fatalf("expected field snapshot pointer, got %q", findings[0].JSONPointer)
	}
	if !strings.Contains(findings[0].Fingerprint, string(diff.ResourceComponentTemplate)) {
		t.Fatalf("expected fingerprint to include resource kind, got %q", findings[0].Fingerprint)
	}
	if !strings.Contains(findings[0].Message, "legacy_id") {
		t.Fatalf("expected message to contain field path, got %q", findings[0].Message)
	}
}

func TestDIF002EmitsMultipleFindingsForMultipleRemovedFields(t *testing.T) {
	findings := runDIF002(t, diff.Result{
		FieldChanges: []diff.FieldChange{
			fieldRemoved("legacy_id", model.FieldRoleProperty, "keyword", "/properties/legacy_id"),
			fieldRemoved("old_status", model.FieldRoleProperty, "keyword", "/properties/old_status"),
		},
	})

	if len(findings) != 2 {
		t.Fatalf("expected two findings, got %#v", findings)
	}
	if !strings.Contains(findings[0].Message, "legacy_id") {
		t.Fatalf("expected first finding to preserve input order for legacy_id, got %q", findings[0].Message)
	}
	if !strings.Contains(findings[1].Message, "old_status") {
		t.Fatalf("expected second finding to preserve input order for old_status, got %q", findings[1].Message)
	}
}

func runDIF002(t *testing.T, result diff.Result) []model.Finding {
	t.Helper()
	findings, err := NewDIF002().Check(Context{}, result)
	if err != nil {
		t.Fatalf("DIF002 returned error: %v", err)
	}
	return findings
}

func fieldRemoved(path string, role model.FieldRole, typ string, beforePointer string) diff.FieldChange {
	return fieldRemovedInResource(diff.ResourceMapping, "mapping.json", path, role, typ, beforePointer)
}

func fieldRemovedInResource(resourceKind diff.ResourceKind, file string, path string, role model.FieldRole, typ string, beforePointer string) diff.FieldChange {
	return diff.FieldChange{
		Kind:     diff.ChangeFieldRemoved,
		Resource: diff.ResourceID{Kind: resourceKind, File: file},
		Field:    diff.FieldID{Path: path, Role: role},
		Before:   snapshot(path, role, typ, beforePointer),
	}
}
