package models

import "time"

// Result represents the output of a documentation QA check
type Result struct {
	Checker   string                 `json:"checker"`
	Timestamp time.Time              `json:"timestamp"`
	Path      string                 `json:"path"`
	Summary   Summary                `json:"summary"`
	Issues    []Issue                `json:"issues"`
	Metadata  map[string]interface{} `json:"metadata"`
}

// Summary contains aggregate statistics about the check results
type Summary struct {
	TotalFiles      int `json:"total_files"`
	FilesWithIssues int `json:"files_with_issues"`
	TotalIssues     int `json:"total_issues"`
	ErrorCount      int `json:"error_count"`
	WarningCount    int `json:"warning_count"`
	InfoCount       int `json:"info_count"`
}

// Issue represents a single documentation issue found by a checker
type Issue struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Column   int    `json:"column"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
	Rule     string `json:"rule"`
	Context  string `json:"context,omitempty"`
}
