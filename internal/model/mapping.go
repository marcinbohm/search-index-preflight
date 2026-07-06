package model

type DynamicSetting string

const (
	DynamicSettingUnspecified DynamicSetting = "unspecified"
	DynamicSettingTrue        DynamicSetting = "true"
	DynamicSettingFalse       DynamicSetting = "false"
	DynamicSettingStrict      DynamicSetting = "strict"
	DynamicSettingRuntime     DynamicSetting = "runtime"
)

type Mapping struct {
	Source           Source
	Dynamic          DynamicSetting
	DateDetection    *bool
	NumericDetection *bool
	Properties       []Field
	DynamicTemplates []DynamicTemplate
	RuntimeFields    []Field
	Meta             map[string]any
	Diagnostics      []Diagnostic
}
