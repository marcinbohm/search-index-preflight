package model

type FieldOrigin string

const (
	FieldOriginMapping           FieldOrigin = "mapping"
	FieldOriginIndexTemplate     FieldOrigin = "index_template"
	FieldOriginComponentTemplate FieldOrigin = "component_template"
)

type FieldRole string

const (
	FieldRoleProperty     FieldRole = "property"
	FieldRoleMultiField   FieldRole = "multi_field"
	FieldRoleRuntimeField FieldRole = "runtime_field"
)

type FieldVisit struct {
	Origin                FieldOrigin
	Role                  FieldRole
	Source                Source
	MappingSource         Source
	IndexTemplateName     string
	ComponentTemplateName string
	Field                 Field
	Path                  string
	JSONPointer           string
}
