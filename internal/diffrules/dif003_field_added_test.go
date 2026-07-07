package diffrules

import (
	"strings"
	"testing"

	"github.com/marcinbohm/search-index-preflight/internal/diff"
	"github.com/marcinbohm/search-index-preflight/internal/model"
)

func TestDIF003EmitsFindingForFieldAdded(t *testing.T) {
	findings := runDIF003(t, diff.Result{
		FieldChanges: []diff.FieldChange{
			fieldAdded("customer_id", model.FieldRoleProperty, "keyword", "/properties/customer_id"),
		},
	})

	if len(findings) != 1 {
		t.Fatalf("expected one finding, got %#v", findings)
	}
	finding := findings[0]
	if finding.ID != "DIF003" {
		t.Fatalf("expected ID DIF003, got %q", finding.ID)
	}
	if finding.Name != "field-added" {
		t.Fatalf("expected name field-added, got %q", finding.Name)
	}
	if finding.Severity != model.SeverityInfo {
		t.Fatalf("expected severity info, got %q", finding.Severity)
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
	if finding.JSONPointer != "/properties/customer_id" {
		t.Fatalf("expected JSON pointer /properties/customer_id, got %q", finding.JSONPointer)
	}
	if !strings.Contains(finding.Message, "customer_id") {
		t.Fatalf("expected message %q to contain customer_id", finding.Message)
	}
	if finding.Remediation == "" {
		t.Fatal("expected remediation")
	}
	expectedFingerprint := "DIF003:mapping:mapping.json:/properties/customer_id:customer_id:property"
	if finding.Fingerprint != expectedFingerprint {
		t.Fatalf("expected fingerprint %q, got %q", expectedFingerprint, finding.Fingerprint)
	}
}

func TestDIF003IgnoresFieldRemoved(t *testing.T) {
	findings := runDIF003(t, diff.Result{
		FieldChanges: []diff.FieldChange{
			fieldRemoved("legacy_id", model.FieldRoleProperty, "keyword", "/properties/legacy_id"),
		},
	})

	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %#v", findings)
	}
}

func TestDIF003IgnoresFieldTypeChanged(t *testing.T) {
	findings := runDIF003(t, diff.Result{
		FieldChanges: []diff.FieldChange{
			fieldTypeChanged("status", model.FieldRoleProperty, "keyword", "long", "/properties/status", "/properties/status"),
		},
	})

	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %#v", findings)
	}
}

func TestDIF003UsesResourcePointerWhenAfterHasNoPointer(t *testing.T) {
	findings := runDIF003(t, diff.Result{
		FieldChanges: []diff.FieldChange{
			{
				Kind:     diff.ChangeFieldAdded,
				Resource: diff.ResourceID{Kind: diff.ResourceMapping, File: "mapping.json", JSONPointer: "/mappings"},
				Field:    diff.FieldID{Path: "customer_id", Role: model.FieldRoleProperty},
				After:    snapshot("customer_id", model.FieldRoleProperty, "keyword", ""),
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

func TestDIF003WorksForMultiFieldAdded(t *testing.T) {
	findings := runDIF003(t, diff.Result{
		FieldChanges: []diff.FieldChange{
			fieldAdded("message.keyword", model.FieldRoleMultiField, "keyword", "/properties/message/fields/keyword"),
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

func TestDIF003WorksForRuntimeFieldAdded(t *testing.T) {
	findings := runDIF003(t, diff.Result{
		FieldChanges: []diff.FieldChange{
			fieldAdded("day_of_week", model.FieldRoleRuntimeField, "keyword", "/runtime/day_of_week"),
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

func TestDIF003WorksForIndexTemplateFieldAdded(t *testing.T) {
	findings := runDIF003(t, diff.Result{
		FieldChanges: []diff.FieldChange{
			fieldAddedInResource(diff.ResourceIndexTemplate, "template.json", "customer_id", model.FieldRoleProperty, "keyword", "/template/mappings/properties/customer_id"),
		},
	})

	if len(findings) != 1 {
		t.Fatalf("expected one finding, got %#v", findings)
	}
	if findings[0].File != "template.json" {
		t.Fatalf("expected template file, got %q", findings[0].File)
	}
	if findings[0].JSONPointer != "/template/mappings/properties/customer_id" {
		t.Fatalf("expected field snapshot pointer, got %q", findings[0].JSONPointer)
	}
	if !strings.Contains(findings[0].Fingerprint, string(diff.ResourceIndexTemplate)) {
		t.Fatalf("expected fingerprint to include resource kind, got %q", findings[0].Fingerprint)
	}
	if !strings.Contains(findings[0].Message, "customer_id") {
		t.Fatalf("expected message to contain field path, got %q", findings[0].Message)
	}
}

func TestDIF003WorksForComponentTemplateFieldAdded(t *testing.T) {
	findings := runDIF003(t, diff.Result{
		FieldChanges: []diff.FieldChange{
			fieldAddedInResource(diff.ResourceComponentTemplate, "component.json", "customer_id", model.FieldRoleProperty, "keyword", "/template/mappings/properties/customer_id"),
		},
	})

	if len(findings) != 1 {
		t.Fatalf("expected one finding, got %#v", findings)
	}
	if findings[0].File != "component.json" {
		t.Fatalf("expected component template file, got %q", findings[0].File)
	}
	if findings[0].JSONPointer != "/template/mappings/properties/customer_id" {
		t.Fatalf("expected field snapshot pointer, got %q", findings[0].JSONPointer)
	}
	if !strings.Contains(findings[0].Fingerprint, string(diff.ResourceComponentTemplate)) {
		t.Fatalf("expected fingerprint to include resource kind, got %q", findings[0].Fingerprint)
	}
	if !strings.Contains(findings[0].Message, "customer_id") {
		t.Fatalf("expected message to contain field path, got %q", findings[0].Message)
	}
}

func TestDIF003EmitsMultipleFindingsForMultipleAddedFields(t *testing.T) {
	findings := runDIF003(t, diff.Result{
		FieldChanges: []diff.FieldChange{
			fieldAdded("customer_id", model.FieldRoleProperty, "keyword", "/properties/customer_id"),
			fieldAdded("account_id", model.FieldRoleProperty, "keyword", "/properties/account_id"),
		},
	})

	if len(findings) != 2 {
		t.Fatalf("expected two findings, got %#v", findings)
	}
	if !strings.Contains(findings[0].Message, "customer_id") {
		t.Fatalf("expected first finding to preserve input order for customer_id, got %q", findings[0].Message)
	}
	if !strings.Contains(findings[1].Message, "account_id") {
		t.Fatalf("expected second finding to preserve input order for account_id, got %q", findings[1].Message)
	}
}

func runDIF003(t *testing.T, result diff.Result) []model.Finding {
	t.Helper()
	findings, err := NewDIF003().Check(Context{}, result)
	if err != nil {
		t.Fatalf("DIF003 returned error: %v", err)
	}
	return findings
}

func fieldAdded(path string, role model.FieldRole, typ string, afterPointer string) diff.FieldChange {
	return fieldAddedInResource(diff.ResourceMapping, "mapping.json", path, role, typ, afterPointer)
}

func fieldAddedInResource(resourceKind diff.ResourceKind, file string, path string, role model.FieldRole, typ string, afterPointer string) diff.FieldChange {
	return diff.FieldChange{
		Kind:     diff.ChangeFieldAdded,
		Resource: diff.ResourceID{Kind: resourceKind, File: file},
		Field:    diff.FieldID{Path: path, Role: role},
		After:    snapshot(path, role, typ, afterPointer),
	}
}
