package checker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
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
// The output is an array of issue objects
type MarkdownlintOutput []MarkdownlintIssue

// MarkdownlintIssue represents a single markdownlint issue
type MarkdownlintIssue struct {
	FileName         string                 `json:"fileName"`
	LineNumber       int                    `json:"lineNumber"`
	RuleNames        []string               `json:"ruleNames"`
	RuleDescription  string                 `json:"ruleDescription"`
	RuleInformation  string                 `json:"ruleInformation"`
	ErrorDetail      *string                `json:"errorDetail"`
	ErrorContext     *string                `json:"errorContext"`
	ErrorRange       []int                  `json:"errorRange"`
	FixInfo          map[string]interface{} `json:"fixInfo"` // Can be null or an object with lineNumber and insertText
	Severity         string                 `json:"severity"`
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
	// Note: markdownlint-cli2 uses --json, markdownlint-cli uses --json
	args := []string{}
	
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
	
	// Add JSON output format
	args = append(args, "--json")

	// Add any extra arguments
	args = append(args, opts.ExtraArgs...)

	// Execute markdownlint
	cmd := exec.CommandContext(ctx, c.markdownlintPath, args...)
	
	// Capture stdout and stderr separately
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	err := cmd.Run()

	// markdownlint returns non-zero exit code when issues are found, which is expected
	// markdownlint outputs JSON to stderr (not stdout)
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			// Exit error is expected when issues are found
			// Continue processing - we'll check if stderr has valid JSON
		} else {
			return nil, fmt.Errorf("failed to execute markdownlint: %w", err)
		}
	}

	// markdownlint outputs JSON to stderr, not stdout
	// Try stderr first (where JSON output goes), then stdout as fallback
	output := stderr.Bytes()
	if len(output) == 0 {
		output = stdout.Bytes()
	}
	
	// If we still have no output, that's an error
	if len(output) == 0 {
		return nil, fmt.Errorf("markdownlint produced no output")
	}
	
	// Check if output looks like JSON
	trimmed := strings.TrimSpace(string(output))
	if !strings.HasPrefix(trimmed, "{") && !strings.HasPrefix(trimmed, "[") {
		// This doesn't look like JSON, it's probably an error message
		return nil, fmt.Errorf("markdownlint execution failed: %s", trimmed)
	}
	
	// Parse markdownlint output
	var markdownlintOutput MarkdownlintOutput
	if err := json.Unmarshal(output, &markdownlintOutput); err != nil {
		return nil, fmt.Errorf("failed to parse markdownlint output: %w (output: %s)", err, string(output))
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
	filesProcessed := make(map[string]bool)

	// Process each issue
	for _, issue := range markdownlintOutput {
		// Track files
		if !filesProcessed[issue.FileName] {
			filesProcessed[issue.FileName] = true
			result.Summary.TotalFiles++
		}
		filesWithIssues[issue.FileName] = true

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

		// Map severity - markdownlint can have "error" or "warning"
		severity := "warning"
		if issue.Severity == "error" {
			severity = "error"
		}

		modelIssue := models.Issue{
			File:     issue.FileName,
			Line:     issue.LineNumber,
			Column:   column,
			Severity: severity,
			Message:  message,
			Rule:     ruleName,
			Context:  context,
		}

		result.Issues = append(result.Issues, modelIssue)
		result.Summary.TotalIssues++
		
		// Count by severity
		if severity == "error" {
			result.Summary.ErrorCount++
		} else {
			result.Summary.WarningCount++
		}
	}

	result.Summary.FilesWithIssues = len(filesWithIssues)

	// Add metadata
	result.Metadata["config_file"] = c.configFile
	result.Metadata["fix_enabled"] = c.fix

	return result
}
