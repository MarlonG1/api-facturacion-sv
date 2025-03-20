package models

type TestResult struct {
	Success  bool            `json:"success"`
	Tests    []ComponentTest `json:"tests"`
	Duration int64           `json:"duration_ms"`
}

type ComponentTest struct {
	Name     string `json:"name"`
	Success  bool   `json:"success"`
	Duration int64  `json:"duration_ms"`
}
