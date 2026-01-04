package models

import "time"

// CheckerStats represents aggregated statistics for a checker
type CheckerStats struct {
	Name         string    `json:"name"`
	TotalRuns    int       `json:"total_runs"`
	LatestRun    time.Time `json:"latest_run"`
	TotalIssues  int       `json:"total_issues"`
	ErrorCount   int       `json:"error_count"`
	WarningCount int       `json:"warning_count"`
	InfoCount    int       `json:"info_count"`
}

// DashboardData represents the aggregated data for the dashboard
type DashboardData struct {
	Checkers      []CheckerStats       `json:"checkers"`
	TotalChecks   int                  `json:"total_checks"`
	LatestResults map[string]*Result   `json:"latest_results"`
	AllResults    []*Result            `json:"all_results"`
}
