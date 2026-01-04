package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/svx/marvin/cli/internal/app/dashboard"
	"github.com/svx/marvin/cli/internal/app/tui"
)

// dashboardCmd represents the dashboard command
var dashboardCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "View aggregated results from all checks",
	Long: `Display a dashboard showing results from all documentation QA checks.

The dashboard reads all check results from the output directory and displays
them in an interactive TUI with summary statistics and drill-down capabilities.

You can navigate between different checkers using Tab/Shift+Tab and view
detailed results by pressing Enter.`,
	RunE: runDashboard,
	Example: `  # Show dashboard with all results
  marvin dashboard

  # Show dashboard with custom output directory
  marvin dashboard --output-dir ./custom-results`,
}

func init() {
	rootCmd.AddCommand(dashboardCmd)
}

func runDashboard(cmd *cobra.Command, args []string) error {
	if verbose {
		fmt.Printf("Loading results from: %s\n", outputDir)
	}

	// 1. Load dashboard data from output directory
	data, err := dashboard.LoadDashboardData(outputDir)
	if err != nil {
		return fmt.Errorf("failed to load dashboard data: %w", err)
	}

	if verbose {
		fmt.Printf("Loaded %d check results\n", data.TotalChecks)
	}

	// Check if there are any results
	if data.TotalChecks == 0 {
		fmt.Println("No check results found.")
		fmt.Printf("Run 'marvin vale' or 'marvin markdownlint' to generate results.\n")
		fmt.Printf("Results are stored in: %s\n", outputDir)
		return nil
	}

	// 2. Display dashboard in TUI
	if err := tui.ShowDashboard(data); err != nil {
		return fmt.Errorf("failed to show dashboard: %w", err)
	}

	return nil
}
