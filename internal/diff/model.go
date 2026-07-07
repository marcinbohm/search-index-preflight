package diff

import (
	"fmt"

	"github.com/marcinbohm/search-index-preflight/internal/model"
)

// ChangeKind identifies the kind of semantic field-level change.
type ChangeKind string

const (
	ChangeFieldAdded       ChangeKind = "field_added"
	ChangeFieldRemoved     ChangeKind = "field_removed"
	ChangeFieldTypeChanged ChangeKind = "field_type_changed"
)

// ResourceKind identifies the mapping-bearing resource kind being compared.
type ResourceKind string

const (
	ResourceMapping           ResourceKind = "mapping"
	ResourceIndexTemplate     ResourceKind = "index_template"
	ResourceComponentTemplate ResourceKind = "component_template"
)

// ResourceID identifies a mapping-bearing resource within a normalized corpus.
type ResourceID struct {
	Kind        ResourceKind
	File        string
	JSONPointer string
}

// FieldID identifies a field by path and role.
type FieldID struct {
	Path string
	Role model.FieldRole
}

// FieldChange describes one field-level difference between two corpora.
type FieldChange struct {
	Kind ChangeKind

	Resource ResourceID
	Field    FieldID

	Before *FieldSnapshot
	After  *FieldSnapshot
}

// FieldSnapshot is the minimal comparable field state captured for a change.
type FieldSnapshot struct {
	Path        string
	Role        model.FieldRole
	Type        string
	JSONPointer string
}

// Result contains the internal semantic diff output.
type Result struct {
	FieldChanges []FieldChange
}

// DuplicateResourceError reports an invalid corpus containing duplicate
// mapping-bearing resource identities.
type DuplicateResourceError struct {
	Resource ResourceID
}

func (e DuplicateResourceError) Error() string {
	return fmt.Sprintf("duplicate resource identity: kind=%s file=%s pointer=%s", e.Resource.Kind, e.Resource.File, e.Resource.JSONPointer)
}
