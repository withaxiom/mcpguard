package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version info — set by GoReleaser ldflags at build time.
var (
	Version   = "dev"
	Commit    = "none"
	Date      = "unknown"
	BuiltBy   = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("mcpguard %s\n", Version)
		if verbose {
			fmt.Printf("  commit:   %s\n", Commit)
			fmt.Printf("  built:    %s\n", Date)
			fmt.Printf("  built by: %s\n", BuiltBy)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
