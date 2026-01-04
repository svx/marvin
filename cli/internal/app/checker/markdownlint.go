package checker

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	"github.com/svx/marvin/cli/internal/pkg/models"
)

// MarkdownlintChecker implements the Checker interface for markdownlint
type MarkdownlintChecker struct {
	configFile        string
	fix               bool
	markdownlintPath  string
}

// MarkdownlintOutput represents markdownlint's JSON output format
// The output is a map where keys are file paths and values are arrays of issues
type MarkdownlintOutput map[string][]MarkdownlintIssue

// MarkdownlintIssue represents a single markdownlint issue
type MarkdownlintIssue struct {
	LineNumber       int      `json:"lineNumber"`
	RuleNames        []string `json:"ruleNames"`
	RuleDescription  string   `json:"ruleDescription"`
	RuleInformation  string   `json:"ruleInformation"`
	ErrorDetail      *string  `json:"errorDetail"`
	ErrorContext     *string  `json:"errorContext"`
	ErrorRange       []int    `json:"errorRange"`
}

// NewMarkdownlintChecker creates a new markdownlint checker
func NewMarkdownlintChecker(configFile string, fix bool, markdownlintPath string) *MarkdownlintChecker {
	if markdownlintPath == "" {
		markdownlintPath = "markdownlint"
	}
	return &MarkdownlintChecker{
		configFile:       configFile,
		fix:              fix,
		markdownlintPath: markdownlintPath,
	}
}

// Name returns the checker name
func (c *MarkdownlintChecker) Name() string {
	return "markdownlint"
}

// Validate validates the checker configuration
func (c *MarkdownlintChecker) Validate() error {
	// Check if markdownlint is available
	if _, err := exec.LookPath(c.markdownlintPath); err != nil {
		return fmt.Errorf("markdownlint not found in PATH: %w", err)
	}
	return nil
}

// Check runs markdownlint and returns the results
func (c *MarkdownlintChecker) Check(ctx context.Context, opts CheckOptions) (*models.Result, error) {
	// Build command arguments
	args := []string{"--json"}

	// Add config file if specified
	if c.configFile != "" {
		args = append(args, "--config", c.configFile)
	}

	// Add fix flag if enabled
	if c.fix {
		args = append(args, "--fix")
	}

	// Add path to check
	args = append(args, opts.Path)

	// Add any extra arguments
	args = append(args, opts.ExtraArgs...)

	// Execute markdownlint
	cmd := exec.CommandContext(ctx, c.markdownlintPath, args...)
	output, err := cmd.Output()

	// markdownlint returns non-zero exit code when issues are found, which is expected
	// We only care about actual execution errors
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			// If there's stderr output, it might be a real error
			if len(exitErr.Stderr) > 0 {
				return nil, fmt.Errorf("markdownlint execution failed: %s", string(exitErr.Stderr))
			}
			// Otherwise, it's just issues found, continue processing
		} else {
			return nil, fmt.Errorf("failed to execute markdownlint: %w", err)
		}
	}

	// Parse markdownlint output
	var markdownlintOutput MarkdownlintOutput
	if err := json.Unmarshal(output, &markdownlintOutput); err != nil {
		return nil, fmt.Errorf("failed to parse markdownlint output: %w", err)
	}

	// Transform to our Result format
	result := c.transformResult(markdownlintOutput, opts.Path)

	return result, nil
}

// transformResult converts markdownlint output to our unified Result format
func (c *MarkdownlintChecker) transformResult(markdownlintOutput MarkdownlintOutput, path string) *models.Result {
	result := &models.Result{
		Checker:   "markdownlint",
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

	// Process each file's issues
	for file, issues := range markdownlintOutput {
		result.Summary.TotalFiles++

		if len(issues) > 0 {
			filesWithIssues[file] = true
		}

		for _, issue := range issues {
			// Get the primary rule name (first in the array)
			ruleName := "unknown"
			if len(issue.RuleNames) > 0 {
				ruleName = issue.RuleNames[0]
			}

			// Build the message
			message := issue.RuleDescription
			if issue.ErrorDetail != nil && *issue.ErrorDetail != "" {
				message = fmt.Sprintf("%s: %s", message, *issue.ErrorDetail)
			}

			// Determine column from ErrorRange if available
			column := 0
			if len(issue.ErrorRange) >= 1 {
				column = issue.ErrorRange[0]
			}

			// Build context from ErrorContext if available
			context := ""
			if issue.ErrorContext != nil {
				context = *issue.ErrorContext
			}

			modelIssue := models.Issue{
				File:     file,
				Line:     issue.LineNumber,
				Column:   column,
				Severity: "warning", // markdownlint treats all issues as warnings by default
				Message:  message,
				Rule:     ruleName,
				Context:  context,
			}

			result.Issues = append(result.Issues, modelIssue)
			result.Summary.TotalIssues++
			result.Summary.WarningCount++
		}
	}

	result.Summary.FilesWithIssues = len(filesWithIssues)

	// Add metadata
	result.Metadata["config_file"] = c.configFile
	result.Metadata["fix_enabled"] = c.fix

	return result
}
