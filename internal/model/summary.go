package model

type Summary struct {
	FilesScanned  int `json:"files_scanned"`
	FindingsTotal int `json:"findings_total"`
	Critical      int `json:"critical"`
	Error         int `json:"error"`
	Warning       int `json:"warning"`
	Info          int `json:"info"`
	ExitCode      int `json:"exit_code"`
}

type Tool struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type RunResult struct {
	SchemaVersion string       `json:"schema_version"`
	Tool          Tool         `json:"tool"`
	Summary       Summary      `json:"summary"`
	Findings      []Finding    `json:"findings"`
	Diagnostics   []Diagnostic `json:"diagnostics"`
}
