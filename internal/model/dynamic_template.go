package model

type DynamicTemplate struct {
	Name             string
	Source           Source
	JSONPointer      string
	Match            string
	Unmatch          string
	PathMatch        string
	PathUnmatch      string
	MatchMappingType string
	Mapping          map[string]any
}
