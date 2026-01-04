package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/svx/marvin/cli/internal/app/dashboard"
	"github.com/svx/marvin/cli/internal/pkg/models"
)

// DashboardModel represents the TUI model for the dashboard
type DashboardModel struct {
	data         *models.DashboardData
	selectedTab  int
	viewMode     string // "summary" or "details"
	content      string
	ready        bool
	quitting     bool
}

// ShowDashboard displays the dashboard in an interactive TUI
func ShowDashboard(data *models.DashboardData) error {
	p := tea.NewProgram(initialDashboardModel(data))
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}

func initialDashboardModel(data *models.DashboardData) DashboardModel {
	m := DashboardModel{
		data:        data,
		selectedTab: 0,
		viewMode:    "summary",
		ready:       true,
	}
	m.content = m.renderContent()
	return m
}

func (m DashboardModel) Init() tea.Cmd {
	return nil
}

func (m DashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit
		case "tab", "right":
			// Switch to next checker tab
			if len(m.data.Checkers) > 0 {
				m.selectedTab = (m.selectedTab + 1) % (len(m.data.Checkers) + 1) // +1 for "All" tab
				m.content = m.renderContent()
			}
		case "shift+tab", "left":
			// Switch to previous checker tab
			if len(m.data.Checkers) > 0 {
				m.selectedTab = (m.selectedTab - 1 + len(m.data.Checkers) + 1) % (len(m.data.Checkers) + 1)
				m.content = m.renderContent()
			}
		case "enter":
			// Toggle between summary and details view
			if m.viewMode == "summary" {
				m.viewMode = "details"
			} else {
				m.viewMode = "summary"
			}
			m.content = m.renderContent()
		}
	}
	return m, nil
}

func (m DashboardModel) View() string {
	if m.quitting {
		return ""
	}

	if !m.ready {
		return "Loading..."
	}

	footer := "\n" + footerStyle.Render("Tab/Shift+Tab: switch tabs | Enter: toggle view | q: quit")
	return m.content + footer
}

// renderContent renders the dashboard content based on current state
func (m DashboardModel) renderContent() string {
	var b strings.Builder

	// Title
	title := " Marvin Dashboard - Documentation QA Results "
	b.WriteString(titleStyle.Render(title))
	b.WriteString("\n\n")

	// Render tabs
	b.WriteString(m.renderTabs())
	b.WriteString("\n\n")

	// Render content based on selected tab
	if m.selectedTab == 0 {
		// "All" tab - show overall summary
		b.WriteString(m.renderOverallSummary())
	} else {
		// Specific checker tab
		checkerIndex := m.selectedTab - 1
		if checkerIndex < len(m.data.Checkers) {
			checker := m.data.Checkers[checkerIndex]
			if m.viewMode == "summary" {
				b.WriteString(m.renderCheckerSummary(checker))
			} else {
				b.WriteString(m.renderCheckerDetails(checker))
			}
		}
	}

	return b.String()
}

// renderTabs renders the tab navigation
func (m DashboardModel) renderTabs() string {
	var tabs []string

	// "All" tab
	allTab := "All"
	if m.selectedTab == 0 {
		tabs = append(tabs, selectedTabStyle.Render(allTab))
	} else {
		tabs = append(tabs, tabStyle.Render(allTab))
	}

	// Checker tabs
	for i, checker := range m.data.Checkers {
		tabName := strings.Title(checker.Name)
		if m.selectedTab == i+1 {
			tabs = append(tabs, selectedTabStyle.Render(tabName))
		} else {
			tabs = append(tabs, tabStyle.Render(tabName))
		}
	}

	return strings.Join(tabs, " ")
}

// renderOverallSummary renders the overall summary across all checkers
func (m DashboardModel) renderOverallSummary() string {
	var b strings.Builder

	b.WriteString(sectionStyle.Render("Overall Summary"))
	b.WriteString("\n")

	// Calculate overall stats
	overallSummary := dashboard.GetOverallSummary(m.data)
	latestCheck := dashboard.GetLatestCheckTime(m.data)

	summaryLines := []string{
		fmt.Sprintf("%s %s",
			summaryLabelStyle.Render("Total Checks Run:"),
			summaryValueStyle.Render(fmt.Sprintf("%d", m.data.TotalChecks))),
		fmt.Sprintf("%s %s",
			summaryLabelStyle.Render("Last Check:"),
			summaryValueStyle.Render(formatRelativeTime(latestCheck))),
		fmt.Sprintf("%s %s",
			summaryLabelStyle.Render("Total Files:"),
			summaryValueStyle.Render(fmt.Sprintf("%d", overallSummary.TotalFiles))),
		fmt.Sprintf("%s %s",
			summaryLabelStyle.Render("Files with Issues:"),
			summaryValueStyle.Render(fmt.Sprintf("%d", overallSummary.FilesWithIssues))),
	}

	// Build issues summary
	issuesParts := []string{}
	if overallSummary.ErrorCount > 0 {
		issuesParts = append(issuesParts,
			errorStyle.Render(fmt.Sprintf("%d errors", overallSummary.ErrorCount)))
	}
	if overallSummary.WarningCount > 0 {
		issuesParts = append(issuesParts,
			warningStyle.Render(fmt.Sprintf("%d warnings", overallSummary.WarningCount)))
	}
	if overallSummary.InfoCount > 0 {
		issuesParts = append(issuesParts,
			infoStyle.Render(fmt.Sprintf("%d info", overallSummary.InfoCount)))
	}

	issuesSummary := fmt.Sprintf("%d", overallSummary.TotalIssues)
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

	// Show individual checker summaries
	b.WriteString("\n")
	b.WriteString(sectionStyle.Render("Checkers"))
	b.WriteString("\n\n")

	for _, checker := range m.data.Checkers {
		b.WriteString(m.renderCheckerCard(checker))
		b.WriteString("\n")
	}

	return b.String()
}

// renderCheckerCard renders a summary card for a checker
func (m DashboardModel) renderCheckerCard(checker models.CheckerStats) string {
	var b strings.Builder

	// Card header
	b.WriteString("  " + cardHeaderStyle.Render(strings.Title(checker.Name)) + "\n")

	// Card content
	b.WriteString(fmt.Sprintf("  %s %s\n",
		summaryLabelStyle.Render("Last run:"),
		summaryValueStyle.Render(formatRelativeTime(checker.LatestRun))))

	issuesParts := []string{}
	if checker.ErrorCount > 0 {
		issuesParts = append(issuesParts,
			errorStyle.Render(fmt.Sprintf("%d errors", checker.ErrorCount)))
	}
	if checker.WarningCount > 0 {
		issuesParts = append(issuesParts,
			warningStyle.Render(fmt.Sprintf("%d warnings", checker.WarningCount)))
	}
	if checker.InfoCount > 0 {
		issuesParts = append(issuesParts,
			infoStyle.Render(fmt.Sprintf("%d info", checker.InfoCount)))
	}

	issuesSummary := fmt.Sprintf("%d", checker.TotalIssues)
	if len(issuesParts) > 0 {
		issuesSummary += " (" + strings.Join(issuesParts, ", ") + ")"
	}

	b.WriteString(fmt.Sprintf("  %s %s\n",
		summaryLabelStyle.Render("Total issues:"),
		summaryValueStyle.Render(issuesSummary)))

	b.WriteString(fmt.Sprintf("  %s %s\n",
		summaryLabelStyle.Render("Total runs:"),
		summaryValueStyle.Render(fmt.Sprintf("%d", checker.TotalRuns))))

	return b.String()
}

// renderCheckerSummary renders the summary view for a specific checker
func (m DashboardModel) renderCheckerSummary(checker models.CheckerStats) string {
	var b strings.Builder

	b.WriteString(sectionStyle.Render(fmt.Sprintf("%s Summary", strings.Title(checker.Name))))
	b.WriteString("\n")

	// Get latest result for this checker
	result := dashboard.GetLatestResultForChecker(m.data, checker.Name)
	if result == nil {
		b.WriteString("  No results available\n")
		return b.String()
	}

	summaryLines := []string{
		fmt.Sprintf("%s %s",
			summaryLabelStyle.Render("Path:"),
			summaryValueStyle.Render(result.Path)),
		fmt.Sprintf("%s %s",
			summaryLabelStyle.Render("Last Run:"),
			summaryValueStyle.Render(formatRelativeTime(result.Timestamp))),
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
			infoStyle.Render(fmt.Sprintf("%d info", result.Summary.InfoCount)))
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

	b.WriteString("\n")
	b.WriteString(infoStyle.Render("  Press Enter to view detailed issues"))
	b.WriteString("\n")

	return b.String()
}

// renderCheckerDetails renders the detailed issues view for a specific checker
func (m DashboardModel) renderCheckerDetails(checker models.CheckerStats) string {
	// Get latest result for this checker
	result := dashboard.GetLatestResultForChecker(m.data, checker.Name)
	if result == nil {
		return "  No results available\n"
	}

	// Reuse the existing formatResults function from viewer.go
	return formatResults(result)
}

// formatRelativeTime formats a time as a relative string
func formatRelativeTime(t time.Time) string {
	if t.IsZero() {
		return "never"
	}

	duration := time.Since(t)

	if duration < time.Minute {
		return "just now"
	} else if duration < time.Hour {
		minutes := int(duration.Minutes())
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	} else if duration < 24*time.Hour {
		hours := int(duration.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	} else {
		days := int(duration.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	}
}
