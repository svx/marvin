package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Global flags
	outputDir  string
	noTUI      bool
	jsonOutput bool
	verbose    bool
	configFile string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "marvin",
	Short: "Documentation quality assurance tool",
	Long: `Marvin is a documentation QA tool that helps you maintain high-quality
documentation by running various checks like prose linting, markdown linting,
and more.

It provides an interactive TUI for viewing results and can output JSON for
integration with other tools or CI/CD pipelines.`,
	Version: "0.1.0",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVar(&outputDir, "output-dir", ".marvin/results", "Output directory for JSON results")
	rootCmd.PersistentFlags().BoolVar(&noTUI, "no-tui", false, "Disable TUI, output plain text to stdout")
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output raw JSON to stdout (implies --no-tui)")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Enable verbose logging")
	rootCmd.PersistentFlags().StringVar(&configFile, "config", ".marvin.yaml", "Path to config file")
}
