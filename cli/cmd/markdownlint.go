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
	markdownlintConfig string
	markdownlintFix    bool
)

// markdownlintCmd represents the markdownlint command
var markdownlintCmd = &cobra.Command{
	Use:   "markdownlint [path]",
	Short: "Run markdownlint on Markdown files",
	Long: `Run markdownlint to check Markdown files for style and syntax issues.

markdownlint checks your Markdown files against a set of rules to ensure
consistent formatting and style. It outputs results in an interactive TUI
by default, or can output JSON for integration with other tools.

By default, markdownlint scans the docs/ directory. You can specify a
different path as an argument.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runMarkdownlint,
	Example: `  # Scan default docs/ directory with TUI
  marvin markdownlint

  # Scan specific directory
  marvin markdownlint ./content

  # Output JSON only
  marvin markdownlint --json

  # Disable TUI, show plain text
  marvin markdownlint --no-tui

  # Use custom markdownlint config
  marvin markdownlint --config .markdownlint.yaml

  # Automatically fix issues where possible
  marvin markdownlint --fix`,
}

func init() {
	rootCmd.AddCommand(markdownlintCmd)

	// Command-specific flags
	markdownlintCmd.Flags().StringVar(&markdownlintConfig, "config", "",
		"markdownlint config file path (default: auto-detect .markdownlint.yaml)")
	markdownlintCmd.Flags().BoolVar(&markdownlintFix, "fix", false,
		"Automatically fix issues where possible")
}

func runMarkdownlint(cmd *cobra.Command, args []string) error {
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
	installed, markdownlintPath, err := detector.IsInstalled("markdownlint")
	if !installed {
		fmt.Println(detector.GetInstallInstructions("markdownlint"))
		return fmt.Errorf("markdownlint not found")
	}

	if verbose {
		fmt.Printf("Found markdownlint at: %s\n", markdownlintPath)
	}

	// 3. Create and run checker
	markdownlintChecker := checker.NewMarkdownlintChecker(markdownlintConfig, markdownlintFix, markdownlintPath)

	// Validate checker
	if err := markdownlintChecker.Validate(); err != nil {
		return fmt.Errorf("markdownlint validation failed: %w", err)
	}

	if verbose {
		fmt.Println("Running markdownlint check...")
	}

	result, err := markdownlintChecker.Check(cmd.Context(), checker.CheckOptions{
		Path:       path,
		ConfigFile: markdownlintConfig,
	})
	if err != nil {
		return fmt.Errorf("markdownlint check failed: %w", err)
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
