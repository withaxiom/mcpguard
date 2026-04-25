package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	outputJSON bool
	verbose    bool
	configPath string
)

var rootCmd = &cobra.Command{
	Use:   "mcpguard",
	Short: "MCP configuration security scanner",
	Long: `MCPGuard scans your MCP (Model Context Protocol) server configurations
for security issues, tracks changes over time, and helps you maintain
a secure AI tooling setup.

Supports: Claude Desktop, Cursor, and more.

Quick start:
  mcpguard scan          # Scan all detected MCP configs
  mcpguard scan --json   # Output as JSON
  mcpguard diff          # Compare current config to last snapshot`,
}

// Execute runs the root command. Cobra prints the error itself when a
// command's RunE returns one, so we just propagate the exit code here.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&outputJSON, "json", false, "Output results as JSON")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "Path to specific config file to scan")
}
