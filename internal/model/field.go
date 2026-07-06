package model

type Field struct {
	Name        string
	Path        string
	Type        string
	ParentPath  string
	Source      Source
	JSONPointer string
	Parameters  map[string]any
	Properties  []Field
	Fields      []Field
	Dynamic     DynamicSetting
	Enabled     *bool
}
