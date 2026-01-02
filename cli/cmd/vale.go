package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/svx/marvin/cli/internal/app/checker"
	"github.com/svx/marvin/cli/internal/app/dependency"
	"github.com/svx/marvin/cli/internal/app/output"
	"github.com/svx/marvin/cli/internal/app/tui"
)

var (
	valeConfig        string
	valeMinAlertLevel string
	valeGlob          string
)

// valeCmd represents the vale command
var valeCmd = &cobra.Command{
	Use:   "vale [path]",
	Short: "Run Vale prose linting",
	Long: `Run Vale prose linting on documentation files.

Vale checks your documentation for style guide violations, grammar issues,
and other prose problems. It outputs results in an interactive TUI by default,
or can output JSON for integration with other tools.

By default, Vale scans the docs/ directory. You can specify a different path
as an argument.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runVale,
	Example: `  # Scan default docs/ directory with TUI
  marvin vale

  # Scan specific directory
  marvin vale ./content

  # Output JSON only
  marvin vale --json

  # Disable TUI, show plain text
  marvin vale --no-tui

  # Use custom Vale config
  marvin vale --config .vale.ini

  # Ignore specific directories
  marvin vale --glob='!node_modules'

  # Ignore multiple patterns
  marvin vale --glob='!{node_modules/*,.vitepress/*}'`,
}

func init() {
	rootCmd.AddCommand(valeCmd)

	// Command-specific flags
	valeCmd.Flags().StringVar(&valeConfig, "config", "", "Vale config file path (default: auto-detect .vale.ini)")
	valeCmd.Flags().StringVar(&valeMinAlertLevel, "min-alert-level", "suggestion", "Minimum alert level (suggestion, warning, error)")
	valeCmd.Flags().StringVar(&valeGlob, "glob", "", "Glob pattern to filter files (e.g., '!node_modules' or '!{dir1/*,dir2/*}')")
}

func runVale(cmd *cobra.Command, args []string) error {
	// 1. Parse arguments and flags
	path := "docs/"
	if len(args) > 0 {
		path = args[0]
	}

	// Check if path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("path does not exist: %s", path)
	}

	if verbose {
		fmt.Printf("Scanning path: %s\n", path)
	}

	// 2. Check dependencies
	detector := dependency.NewMultiDetector()
	installed, valePath, err := detector.IsInstalled("vale")
	if !installed {
		fmt.Println(detector.GetInstallInstructions("vale"))
		return fmt.Errorf("vale not found")
	}

	if verbose {
		fmt.Printf("Found vale at: %s\n", valePath)
	}

	// 3. Create and run checker
	valeChecker := checker.NewValeChecker(valeConfig, valeMinAlertLevel, valePath, valeGlob)
	
	// Validate checker
	if err := valeChecker.Validate(); err != nil {
		return fmt.Errorf("vale validation failed: %w", err)
	}

	if verbose {
		fmt.Println("Running Vale check...")
	}

	result, err := valeChecker.Check(cmd.Context(), checker.CheckOptions{
		Path:       path,
		ConfigFile: valeConfig,
	})
	if err != nil {
		return fmt.Errorf("vale check failed: %w", err)
	}

	// 4. Save results
	writer := output.NewJSONWriter(outputDir)
	outputPath, err := writer.Write(result)
	if err != nil {
		return fmt.Errorf("failed to save results: %w", err)
	}

	if verbose {
		fmt.Printf("Results saved to: %s\n", outputPath)
	}

	// 5. Display output
	if jsonOutput {
		// Output raw JSON to stdout
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(result); err != nil {
			return fmt.Errorf("failed to encode JSON: %w", err)
		}
	} else if noTUI {
		// Output plain text
		formatter := output.NewPlainTextFormatter()
		if err := formatter.Format(result, os.Stdout); err != nil {
			return fmt.Errorf("failed to format output: %w", err)
		}
		fmt.Printf("\nResults saved to: %s\n", outputPath)
	} else {
		// Show TUI
		if err := tui.ShowResults(result); err != nil {
			return fmt.Errorf("failed to show TUI: %w", err)
		}
		fmt.Printf("\nResults saved to: %s\n", outputPath)
	}

	// Exit with non-zero code if there are errors
	if result.Summary.ErrorCount > 0 {
		os.Exit(1)
	}

	return nil
}
