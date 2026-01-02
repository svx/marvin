package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/svx/marvin/cli/internal/pkg/models"
)

// Model represents the TUI model
type Model struct {
	result   *models.Result
	content  string
	ready    bool
	quitting bool
}

// ShowResults displays the results in an interactive TUI
func ShowResults(result *models.Result) error {
	p := tea.NewProgram(initialModel(result))
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}

func initialModel(result *models.Result) Model {
	return Model{
		result:  result,
		content: formatResults(result),
		ready:   true,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) View() string {
	if m.quitting {
		return ""
	}

	if !m.ready {
		return "Loading..."
	}

	return m.content + "\n" + footerStyle.Render("Press q to quit")
}

// formatResults formats the result for display
func formatResults(result *models.Result) string {
	var b strings.Builder

	// Title
	title := fmt.Sprintf(" Marvin - %s Results ", strings.Title(result.Checker))
	b.WriteString(titleStyle.Render(title))
	b.WriteString("\n\n")

	// Summary section
	b.WriteString(sectionStyle.Render("Summary"))
	b.WriteString("\n")

	summaryLines := []string{
		fmt.Sprintf("%s %s",
			summaryLabelStyle.Render("Path:"),
			summaryValueStyle.Render(result.Path)),
		fmt.Sprintf("%s %s",
			summaryLabelStyle.Render("Files Scanned:"),
			summaryValueStyle.Render(fmt.Sprintf("%d", result.Summary.TotalFiles))),
		fmt.Sprintf("%s %s",
			summaryLabelStyle.Render("Files with Issues:"),
			summaryValueStyle.Render(fmt.Sprintf("%d", result.Summary.FilesWithIssues))),
	}

	// Build issues summary
	issuesParts := []string{}
	if result.Summary.ErrorCount > 0 {
		issuesParts = append(issuesParts,
			errorStyle.Render(fmt.Sprintf("%d errors", result.Summary.ErrorCount)))
	}
	if result.Summary.WarningCount > 0 {
		issuesParts = append(issuesParts,
			warningStyle.Render(fmt.Sprintf("%d warnings", result.Summary.WarningCount)))
	}
	if result.Summary.InfoCount > 0 {
		issuesParts = append(issuesParts,
			infoStyle.Render(fmt.Sprintf("%d suggestions", result.Summary.InfoCount)))
	}

	issuesSummary := fmt.Sprintf("%d", result.Summary.TotalIssues)
	if len(issuesParts) > 0 {
		issuesSummary += " (" + strings.Join(issuesParts, ", ") + ")"
	}

	summaryLines = append(summaryLines,
		fmt.Sprintf("%s %s",
			summaryLabelStyle.Render("Total Issues:"),
			summaryValueStyle.Render(issuesSummary)))

	for _, line := range summaryLines {
		b.WriteString("  " + line + "\n")
	}

	// Issues section
	b.WriteString("\n")
	b.WriteString(sectionStyle.Render("Issues"))
	b.WriteString("\n")

	if len(result.Issues) == 0 {
		b.WriteString(infoStyle.Render("  ✓ No issues found!"))
		b.WriteString("\n")
	} else {
		// Display issues
		for i, issue := range result.Issues {
			if i > 0 {
				b.WriteString("\n")
			}

			// File location
			location := fmt.Sprintf("%s:%d:%d", issue.File, issue.Line, issue.Column)
			b.WriteString("  " + fileLocationStyle.Render(location) + "\n")

			// Severity and rule
			severityStyle := getSeverityStyle(issue.Severity)
			severityText := fmt.Sprintf("[%s]", issue.Severity)
			b.WriteString("  " + severityStyle.Render(severityText) + " " + ruleStyle.Render(issue.Rule) + "\n")

			// Message
			b.WriteString("  " + messageStyle.Render(issue.Message) + "\n")

			// Context (if available)
			if issue.Context != "" {
				b.WriteString("  " + contextStyle.Render("Context: "+issue.Context) + "\n")
			}
		}
	}

	return b.String()
}
