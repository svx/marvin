package dashboard

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/svx/marvin/cli/internal/pkg/models"
)

// LoadDashboardData scans the output directory and aggregates all check results
func LoadDashboardData(outputDir string) (*models.DashboardData, error) {
	// Scan directory for JSON files
	files, err := scanResultsDirectory(outputDir)
	if err != nil {
		return nil, fmt.Errorf("failed to scan results directory: %w", err)
	}

	if len(files) == 0 {
		return &models.DashboardData{
			Checkers:      []models.CheckerStats{},
			TotalChecks:   0,
			LatestResults: make(map[string]*models.Result),
			AllResults:    []*models.Result{},
		}, nil
	}

	// Parse all result files
	var allResults []*models.Result
	for _, file := range files {
		result, err := parseResultFile(file)
		if err != nil {
			// Skip files that can't be parsed
			continue
		}
		allResults = append(allResults, result)
	}

	// Aggregate results
	dashboardData := aggregateResults(allResults)

	return dashboardData, nil
}

// scanResultsDirectory finds all JSON files in the results directory
func scanResultsDirectory(dir string) ([]string, error) {
	// Check if directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return []string{}, nil
	}

	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Only include JSON files
		if filepath.Ext(path) == ".json" {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

// parseResultFile reads and parses a single result JSON file
func parseResultFile(path string) (*models.Result, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", path, err)
	}

	var result models.Result
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON from %s: %w", path, err)
	}

	return &result, nil
}

// aggregateResults processes all results and creates dashboard data
func aggregateResults(results []*models.Result) *models.DashboardData {
	// Group results by checker
	checkerResults := make(map[string][]*models.Result)
	for _, result := range results {
		checkerResults[result.Checker] = append(checkerResults[result.Checker], result)
	}

	// Calculate stats for each checker
	var checkerStats []models.CheckerStats
	latestResults := make(map[string]*models.Result)

	for checkerName, results := range checkerResults {
		// Sort by timestamp (newest first)
		sort.Slice(results, func(i, j int) bool {
			return results[i].Timestamp.After(results[j].Timestamp)
		})

		// Get latest result
		latest := results[0]
		latestResults[checkerName] = latest

		// Calculate aggregated stats
		stats := models.CheckerStats{
			Name:         checkerName,
			TotalRuns:    len(results),
			LatestRun:    latest.Timestamp,
			TotalIssues:  0,
			ErrorCount:   0,
			WarningCount: 0,
			InfoCount:    0,
		}

		// Aggregate counts from all runs
		for _, result := range results {
			stats.TotalIssues += result.Summary.TotalIssues
			stats.ErrorCount += result.Summary.ErrorCount
			stats.WarningCount += result.Summary.WarningCount
			stats.InfoCount += result.Summary.InfoCount
		}

		checkerStats = append(checkerStats, stats)
	}

	// Sort checker stats by name for consistent display
	sort.Slice(checkerStats, func(i, j int) bool {
		return strings.ToLower(checkerStats[i].Name) < strings.ToLower(checkerStats[j].Name)
	})

	return &models.DashboardData{
		Checkers:      checkerStats,
		TotalChecks:   len(results),
		LatestResults: latestResults,
		AllResults:    results,
	}
}

// GetLatestResultForChecker returns the most recent result for a specific checker
func GetLatestResultForChecker(data *models.DashboardData, checkerName string) *models.Result {
	return data.LatestResults[checkerName]
}

// GetCheckerNames returns a sorted list of all checker names
func GetCheckerNames(data *models.DashboardData) []string {
	names := make([]string, len(data.Checkers))
	for i, checker := range data.Checkers {
		names[i] = checker.Name
	}
	return names
}

// GetOverallSummary calculates overall summary across all latest results
func GetOverallSummary(data *models.DashboardData) models.Summary {
	summary := models.Summary{
		TotalFiles:      0,
		FilesWithIssues: 0,
		TotalIssues:     0,
		ErrorCount:      0,
		WarningCount:    0,
		InfoCount:       0,
	}

	for _, result := range data.LatestResults {
		summary.TotalFiles += result.Summary.TotalFiles
		summary.FilesWithIssues += result.Summary.FilesWithIssues
		summary.TotalIssues += result.Summary.TotalIssues
		summary.ErrorCount += result.Summary.ErrorCount
		summary.WarningCount += result.Summary.WarningCount
		summary.InfoCount += result.Summary.InfoCount
	}

	return summary
}

// GetLatestCheckTime returns the timestamp of the most recent check
func GetLatestCheckTime(data *models.DashboardData) time.Time {
	var latest time.Time
	for _, result := range data.LatestResults {
		if result.Timestamp.After(latest) {
			latest = result.Timestamp
		}
	}
	return latest
}
