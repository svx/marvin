package output

import (
	"fmt"
	"io"

	"github.com/svx/marvin/cli/internal/pkg/models"
)

// Formatter defines the interface for formatting check results
type Formatter interface {
	// Format writes formatted results to the writer
	Format(result *models.Result, w io.Writer) error
}

// PlainTextFormatter formats results as plain text
type PlainTextFormatter struct{}

// NewPlainTextFormatter creates a new plain text formatter
func NewPlainTextFormatter() *PlainTextFormatter {
	return &PlainTextFormatter{}
}

// Format writes the result as plain text
func (f *PlainTextFormatter) Format(result *models.Result, w io.Writer) error {
	// Header
	fmt.Fprintf(w, "Marvin - %s Results\n", result.Checker)
	fmt.Fprintf(w, "═══════════════════════════════════════════════════════════\n\n")

	// Summary
	fmt.Fprintf(w, "Summary:\n")
	fmt.Fprintf(w, "  Path: %s\n", result.Path)
	fmt.Fprintf(w, "  Files Scanned: %d\n", result.Summary.TotalFiles)
	fmt.Fprintf(w, "  Files with Issues: %d\n", result.Summary.FilesWithIssues)
	fmt.Fprintf(w, "  Total Issues: %d", result.Summary.TotalIssues)

	if result.Summary.ErrorCount > 0 || result.Summary.WarningCount > 0 || result.Summary.InfoCount > 0 {
		fmt.Fprintf(w, " (")
		parts := []string{}
		if result.Summary.ErrorCount > 0 {
			parts = append(parts, fmt.Sprintf("%d errors", result.Summary.ErrorCount))
		}
		if result.Summary.WarningCount > 0 {
			parts = append(parts, fmt.Sprintf("%d warnings", result.Summary.WarningCount))
		}
		if result.Summary.InfoCount > 0 {
			parts = append(parts, fmt.Sprintf("%d suggestions", result.Summary.InfoCount))
		}
		for i, part := range parts {
			if i > 0 {
				fmt.Fprintf(w, ", ")
			}
			fmt.Fprintf(w, "%s", part)
		}
		fmt.Fprintf(w, ")")
	}
	fmt.Fprintf(w, "\n\n")

	// Issues
	if len(result.Issues) > 0 {
		fmt.Fprintf(w, "Issues:\n")
		fmt.Fprintf(w, "───────────────────────────────────────────────────────────\n\n")

		for _, issue := range result.Issues {
			// File location
			fmt.Fprintf(w, "%s:%d:%d\n", issue.File, issue.Line, issue.Column)

			// Severity and rule
			fmt.Fprintf(w, "[%s] %s\n", issue.Severity, issue.Rule)

			// Message
			fmt.Fprintf(w, "%s\n", issue.Message)

			// Context (if available)
			if issue.Context != "" {
				fmt.Fprintf(w, "Context: %s\n", issue.Context)
			}

			fmt.Fprintf(w, "\n")
		}
	} else {
		fmt.Fprintf(w, "No issues found! ✓\n\n")
	}

	return nil
}
