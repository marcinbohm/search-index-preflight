package model

type IndexTemplate struct {
	Name          string
	Source        Source
	IndexPatterns []string
	Priority      *int
	ComposedOf    []string
	Template      TemplateBody
	DataStream    bool
	Meta          map[string]any
	Diagnostics   []Diagnostic
}

type ComponentTemplate struct {
	Name        string
	Source      Source
	Template    TemplateBody
	Version     *int
	Meta        map[string]any
	Diagnostics []Diagnostic
}

type TemplateBody struct {
	Settings map[string]any
	Mappings *Mapping
	Aliases  map[string]any
}
