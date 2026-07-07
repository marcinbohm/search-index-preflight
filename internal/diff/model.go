package diff

import "github.com/marcinbohm/search-index-preflight/internal/model"

type ChangeKind string

const (
	ChangeFieldAdded       ChangeKind = "field_added"
	ChangeFieldRemoved     ChangeKind = "field_removed"
	ChangeFieldTypeChanged ChangeKind = "field_type_changed"
)

type ResourceKind string

const (
	ResourceMapping           ResourceKind = "mapping"
	ResourceIndexTemplate     ResourceKind = "index_template"
	ResourceComponentTemplate ResourceKind = "component_template"
)

type ResourceID struct {
	Kind        ResourceKind
	File        string
	JSONPointer string
}

type FieldID struct {
	Path string
	Role model.FieldRole
}

type FieldChange struct {
	Kind ChangeKind

	Resource ResourceID
	Field    FieldID

	Before *FieldSnapshot
	After  *FieldSnapshot
}

type FieldSnapshot struct {
	Path        string
	Role        model.FieldRole
	Type        string
	JSONPointer string
}

type Result struct {
	FieldChanges []FieldChange
}
