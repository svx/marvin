package checker

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	"github.com/svx/marvin/cli/internal/pkg/models"
)

// ValeChecker implements the Checker interface for Vale
type ValeChecker struct {
	configFile      string
	minAlertLevel   string
	valePath        string
	glob            string
}

// ValeOutput represents Vale's JSON output format
type ValeOutput map[string][]ValeAlert

// ValeAlert represents a single Vale alert
type ValeAlert struct {
	Check       string `json:"Check"`
	Description string `json:"Description"`
	Line        int    `json:"Line"`
	Link        string `json:"Link"`
	Message     string `json:"Message"`
	Severity    string `json:"Severity"`
	Span        []int  `json:"Span"`
	Match       string `json:"Match"`
}

// NewValeChecker creates a new Vale checker
func NewValeChecker(configFile, minAlertLevel, valePath, glob string) *ValeChecker {
	if minAlertLevel == "" {
		minAlertLevel = "suggestion"
	}
	if valePath == "" {
		valePath = "vale"
	}
	return &ValeChecker{
		configFile:    configFile,
		minAlertLevel: minAlertLevel,
		valePath:      valePath,
		glob:          glob,
	}
}

// Name returns the checker name
func (c *ValeChecker) Name() string {
	return "vale"
}

// Validate validates the checker configuration
func (c *ValeChecker) Validate() error {
	// Check if vale is available
	if _, err := exec.LookPath(c.valePath); err != nil {
		return fmt.Errorf("vale not found in PATH: %w", err)
	}
	return nil
}

// Check runs Vale and returns the results
func (c *ValeChecker) Check(ctx context.Context, opts CheckOptions) (*models.Result, error) {
	// Build command arguments
	args := []string{"--output=JSON"}

	// Add config file if specified
	if c.configFile != "" {
		args = append(args, "--config="+c.configFile)
	}

	// Add min alert level
	if c.minAlertLevel != "" {
		args = append(args, "--minAlertLevel="+c.minAlertLevel)
	}

	// Add glob pattern if specified
	if c.glob != "" {
		args = append(args, "--glob="+c.glob)
	}

	// Add path to check
	args = append(args, opts.Path)

	// Add any extra arguments
	args = append(args, opts.ExtraArgs...)

	// Execute Vale
	cmd := exec.CommandContext(ctx, c.valePath, args...)
	output, err := cmd.Output()
	
	// Vale returns non-zero exit code when issues are found, which is expected
	// We only care about actual execution errors
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			// If there's stderr output, it might be a real error
			if len(exitErr.Stderr) > 0 {
				return nil, fmt.Errorf("vale execution failed: %s", string(exitErr.Stderr))
			}
			// Otherwise, it's just issues found, continue processing
		} else {
			return nil, fmt.Errorf("failed to execute vale: %w", err)
		}
	}

	// Parse Vale output
	var valeOutput ValeOutput
	if err := json.Unmarshal(output, &valeOutput); err != nil {
		return nil, fmt.Errorf("failed to parse vale output: %w", err)
	}

	// Transform to our Result format
	result := c.transformResult(valeOutput, opts.Path)
	
	return result, nil
}

// transformResult converts Vale output to our unified Result format
func (c *ValeChecker) transformResult(valeOutput ValeOutput, path string) *models.Result {
	result := &models.Result{
		Checker:   "vale",
		Timestamp: time.Now(),
		Path:      path,
		Summary: models.Summary{
			TotalFiles:      0,
			FilesWithIssues: 0,
			TotalIssues:     0,
			ErrorCount:      0,
			WarningCount:    0,
			InfoCount:       0,
		},
		Issues:   []models.Issue{},
		Metadata: make(map[string]interface{}),
	}

	filesWithIssues := make(map[string]bool)

	// Process each file's alerts
	for file, alerts := range valeOutput {
		result.Summary.TotalFiles++
		
		if len(alerts) > 0 {
			filesWithIssues[file] = true
		}

		for _, alert := range alerts {
			issue := models.Issue{
				File:     file,
				Line:     alert.Line,
				Column:   0, // Vale doesn't provide column in the same way
				Severity: normalizeSeverity(alert.Severity),
				Message:  alert.Message,
				Rule:     alert.Check,
				Context:  alert.Match,
			}

			// Set column from Span if available
			if len(alert.Span) >= 1 {
				issue.Column = alert.Span[0]
			}

			result.Issues = append(result.Issues, issue)
			result.Summary.TotalIssues++

			// Count by severity
			switch issue.Severity {
			case "error":
				result.Summary.ErrorCount++
			case "warning":
				result.Summary.WarningCount++
			case "info", "suggestion":
				result.Summary.InfoCount++
			}
		}
	}

	result.Summary.FilesWithIssues = len(filesWithIssues)

	// Add metadata
	result.Metadata["config_file"] = c.configFile
	result.Metadata["min_alert_level"] = c.minAlertLevel

	return result
}

// normalizeSeverity converts Vale severity to our standard format
func normalizeSeverity(severity string) string {
	switch severity {
	case "error":
		return "error"
	case "warning":
		return "warning"
	case "suggestion":
		return "info"
	default:
		return "info"
	}
}
