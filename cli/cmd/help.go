package cmd

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	// Styles for help output
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("99")).
			MarginBottom(1)

	sectionStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("86")).
			MarginTop(1)

	commandStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("212"))

	descStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("246"))
)

// helpCmd represents the help command
var helpCmd = &cobra.Command{
	Use:   "help [command]",
	Short: "Help about any command",
	Long: `Help provides help for any command in the application.
Simply type marvin help [path to command] for full details.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			showMainHelp()
		} else {
			// Let Cobra handle specific command help
			cmd.Root().SetArgs(append(args, "--help"))
			cmd.Root().Execute()
		}
	},
}

func init() {
	rootCmd.AddCommand(helpCmd)
}

func showMainHelp() {
	// Title
	fmt.Println(titleStyle.Render("Marvin - Documentation QA Tool"))

	// Description
	fmt.Println("A documentation quality assurance tool with an interactive TUI.")
	fmt.Println()

	// Usage
	fmt.Println(sectionStyle.Render("Usage:"))
	fmt.Println("  marvin [command] [flags]")
	fmt.Println()

	// Available Commands
	fmt.Println(sectionStyle.Render("Available Commands:"))
	
	commands := []struct {
		name string
		desc string
	}{
		{"vale", "Run Vale prose linting on documentation"},
		{"markdownlint", "Run markdownlint on Markdown files"},
		{"dashboard", "View aggregated results from all checks"},
		{"help", "Help about any command"},
	}

	for _, c := range commands {
		fmt.Printf("  %s  %s\n", 
			commandStyle.Render(fmt.Sprintf("%-12s", c.name)),
			descStyle.Render(c.desc))
	}
	fmt.Println()

	// Global Flags
	fmt.Println(sectionStyle.Render("Global Flags:"))
	fmt.Println("  --output-dir string   Output directory for JSON results (default \".marvin/results\")")
	fmt.Println("  --no-tui              Disable TUI, output plain text to stdout")
	fmt.Println("  --json                Output raw JSON to stdout (implies --no-tui)")
	fmt.Println("  --verbose             Enable verbose logging")
	fmt.Println("  --config string       Path to config file (default \".marvin.yaml\")")
	fmt.Println("  -h, --help            Help for marvin")
	fmt.Println("  -v, --version         Version for marvin")
	fmt.Println()

	// Examples
	fmt.Println(sectionStyle.Render("Examples:"))
	fmt.Println("  # Run Vale on default docs/ directory")
	fmt.Println("  marvin vale")
	fmt.Println()
	fmt.Println("  # Run markdownlint on default docs/ directory")
	fmt.Println("  marvin markdownlint")
	fmt.Println()
	fmt.Println("  # View dashboard with all check results")
	fmt.Println("  marvin dashboard")
	fmt.Println()
	fmt.Println("  # Run Vale on specific directory")
	fmt.Println("  marvin vale ./content")
	fmt.Println()
	fmt.Println("  # Output JSON only")
	fmt.Println("  marvin vale --json")
	fmt.Println()
	fmt.Println("  # Get help for a specific command")
	fmt.Println("  marvin help vale")
	fmt.Println()

	// Footer
	fmt.Println("Use \"marvin [command] --help\" for more information about a command.")
}
