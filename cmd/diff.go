package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/withaxiom/mcpguard/internal/config"
	diffpkg "github.com/withaxiom/mcpguard/internal/diff"
)

var (
	diffBefore string
	diffAfter  string
)

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Compare MCP config changes",
	Long: `Diff compares MCP configurations to detect changes.

By default, it compares the current live config against the last saved snapshot.
You can also compare two specific snapshot files.

Examples:
  mcpguard diff                                    # Current vs last snapshot
  mcpguard diff --before snap1.json --after snap2.json  # Two snapshots
  mcpguard scan --save && mcpguard diff             # Save then diff next time`,
	RunE: runDiff,
}

func init() {
	diffCmd.Flags().StringVar(&diffBefore, "before", "", "Path to 'before' snapshot file")
	diffCmd.Flags().StringVar(&diffAfter, "after", "", "Path to 'after' snapshot file")
	rootCmd.AddCommand(diffCmd)
}

func runDiff(cmd *cobra.Command, args []string) error {
	bold := color.New(color.Bold)
	green := color.New(color.FgGreen, color.Bold)
	red := color.New(color.FgRed, color.Bold)
	yellow := color.New(color.FgYellow)
	dim := color.New(color.Faint)

	// If explicit before/after snapshots provided
	if diffBefore != "" && diffAfter != "" {
		beforeSnap, err := diffpkg.LoadSnapshot(diffBefore)
		if err != nil {
			return fmt.Errorf("loading 'before' snapshot: %w", err)
		}
		afterSnap, err := diffpkg.LoadSnapshot(diffAfter)
		if err != nil {
			return fmt.Errorf("loading 'after' snapshot: %w", err)
		}
		return showDiff(beforeSnap.Config, afterSnap.Config, beforeSnap.Source, bold, green, red, yellow, dim)
	}

	// Default: compare current config vs saved snapshot
	configs, err := config.AutoDetect()
	if err != nil {
		return err
	}

	if len(configs) == 0 {
		fmt.Println("No MCP configurations found.")
		return nil
	}

	snapshotDir := mcpguardDir()
	anyDiff := false

	for _, cfg := range configs {
		snapPath := filepath.Join(snapshotDir, fmt.Sprintf("%s-latest.json", cfg.Source))
		snap, err := diffpkg.LoadSnapshot(snapPath)
		if err != nil {
			if os.IsNotExist(err) || err != nil {
				if !outputJSON {
					dim.Printf("No saved snapshot for %s. Run 'mcpguard scan --save' first.\n", cfg.Source)
				}
				continue
			}
		}

		if !outputJSON {
			bold.Printf("\n📊 Diff: %s\n", cfg.Source)
			dim.Printf("   Snapshot: %s\n", snap.Timestamp)
			fmt.Println()
		}

		if err := showDiff(snap.Config, cfg, cfg.Source, bold, green, red, yellow, dim); err != nil {
			return err
		}
		anyDiff = true
	}

	if !anyDiff && !outputJSON {
		fmt.Println("No snapshots found to compare against.")
		fmt.Println("Run 'mcpguard scan --save' to create a baseline snapshot.")
	}

	return nil
}

func showDiff(before, after *config.MCPConfig, source string, bold, green, red, yellow, dim *color.Color) error {
	result := diffpkg.Configs(before, after)

	if outputJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(result)
	}

	if !result.HasChanges() {
		green.Println("   ✅ No changes detected")
		fmt.Println()
		return nil
	}

	if len(result.Added) > 0 {
		green.Printf("   ➕ Added servers (%d):\n", len(result.Added))
		for _, name := range result.Added {
			fmt.Printf("      + %s\n", name)
		}
		fmt.Println()
	}

	if len(result.Removed) > 0 {
		red.Printf("   ➖ Removed servers (%d):\n", len(result.Removed))
		for _, name := range result.Removed {
			fmt.Printf("      - %s\n", name)
		}
		fmt.Println()
	}

	if len(result.Modified) > 0 {
		yellow.Printf("   ✏️  Modified (%d changes):\n", len(result.Modified))
		for _, entry := range result.Modified {
			fmt.Printf("      ~ %s.%s\n", entry.Server, entry.Field)
			red.Printf("        - %s\n", entry.Before)
			green.Printf("        + %s\n", entry.After)
		}
		fmt.Println()
	}

	return nil
}
