package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/withaxiom/mcpguard/internal/config"
	diffpkg "github.com/withaxiom/mcpguard/internal/diff"
	"github.com/withaxiom/mcpguard/internal/fingerprint"
	"github.com/withaxiom/mcpguard/internal/scanner"
)

var (
	scanSource   string
	saveSnapshot bool
	listSources  bool
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan MCP configurations for security issues",
	Long: `Scan detects MCP server configurations on your system,
analyzes them for security issues, computes fingerprints,
and optionally saves a snapshot for future diffing.

Examples:
  mcpguard scan                       # auto-detect and scan all configs
  mcpguard scan --source cursor       # scan only Cursor's MCP config
  mcpguard scan --save                # scan and save a snapshot for later diffing
  mcpguard scan --list-sources        # show supported sources and where they live
  mcpguard scan -c path/to/file.json  # scan an explicit config file`,
	RunE: runScan,
}

func init() {
	scanCmd.Flags().StringVarP(&scanSource, "source", "s", "", "Scan only a specific source (cursor, claude-desktop)")
	scanCmd.Flags().BoolVar(&saveSnapshot, "save", false, "Save a snapshot of the current config for future diffing")
	scanCmd.Flags().BoolVar(&listSources, "list-sources", false, "List supported MCP config sources and their default paths, then exit")
	// Don't dump --help on runtime errors (parse failures, missing configs);
	// usage info should only appear for actual misuse of flags.
	scanCmd.SilenceUsage = true
	rootCmd.AddCommand(scanCmd)
}

func runScan(cmd *cobra.Command, args []string) error {
	bold := color.New(color.Bold)
	green := color.New(color.FgGreen, color.Bold)
	yellow := color.New(color.FgYellow, color.Bold)
	red := color.New(color.FgRed, color.Bold)
	cyan := color.New(color.FgCyan)
	dim := color.New(color.Faint)

	if listSources {
		return printSources(bold, dim)
	}

	var configs []*config.MCPConfig

	if configPath != "" {
		// Parse a specific file — try each parser, collecting per-parser errors
		// so the user gets useful context when none succeed.
		var parseErrs []string
		for _, p := range config.AllParsers() {
			cfg, err := p.Parse(configPath)
			if err == nil {
				configs = append(configs, cfg)
				parseErrs = nil
				break
			}
			parseErrs = append(parseErrs, fmt.Sprintf("%s: %v", p.Name(), err))
		}
		if len(configs) == 0 {
			return fmt.Errorf("could not parse %s with any known parser:\n  - %s",
				configPath, strings.Join(parseErrs, "\n  - "))
		}
	} else {
		// Auto-detect all configs
		var err error
		configs, err = config.AutoDetect()
		if err != nil {
			return err
		}
	}

	if len(configs) == 0 {
		fmt.Println("No MCP configurations found.")
		fmt.Println()
		dim.Println("MCPGuard looks for configs in these locations:")
		for _, p := range config.AllParsers() {
			fmt.Printf("  • %s: %s\n", p.Name(), p.DefaultPath())
		}
		return nil
	}

	// Filter by source if specified
	if scanSource != "" {
		var filtered []*config.MCPConfig
		for _, cfg := range configs {
			if cfg.Source == scanSource {
				filtered = append(filtered, cfg)
			}
		}
		configs = filtered
		if len(configs) == 0 {
			return fmt.Errorf("no configs found for source: %s", scanSource)
		}
	}

	type scanConfigOutput struct {
		Source      string            `json:"source"`
		Path       string            `json:"path"`
		FileHash   string            `json:"file_hash"`
		ConfigHash string            `json:"config_hash"`
		Servers    int               `json:"server_count"`
		Warnings   []config.Warning  `json:"warnings"`
	}
	type scanOutput struct {
		Configs []scanConfigOutput `json:"configs"`
	}

	var jsonOutput scanOutput

	for _, cfg := range configs {
		if !outputJSON {
			bold.Printf("\n⚡ Scanning: %s\n", cfg.Source)
			dim.Printf("   Path: %s\n", cfg.Path)
			fmt.Println()
		}

		// Compute fingerprints
		fileHash, err := fingerprint.FileHash(cfg.Path)
		if err != nil {
			return fmt.Errorf("fingerprinting %s: %w", cfg.Path, err)
		}
		configHash, err := fingerprint.ConfigHash(cfg)
		if err != nil {
			return fmt.Errorf("computing config hash: %w", err)
		}

		if !outputJSON {
			dim.Printf("   File SHA-256:   %s\n", fileHash[:16]+"...")
			dim.Printf("   Config SHA-256: %s\n", configHash[:16]+"...")
			fmt.Println()
		}

		// List servers
		serverNames := sortedServerNames(cfg)
		if !outputJSON {
			bold.Printf("   📦 Servers (%d):\n", len(cfg.Servers))
			for _, name := range serverNames {
				srv := cfg.Servers[name]
				status := green.Sprint("●")
				if srv.Disabled {
					status = dim.Sprint("○")
				}
				fmt.Printf("      %s %s", status, name)
				if srv.Command != "" {
					dim.Printf(" → %s", srv.Command)
				} else if srv.URL != "" {
					dim.Printf(" → %s", srv.URL)
				}
				fmt.Println()

				// Show redacted env vars
				if len(srv.Env) > 0 && verbose {
					redacted := fingerprint.RedactEnv(srv.Env)
					for k, v := range redacted {
						dim.Printf("         %s=%s\n", k, v)
					}
				}
			}
			fmt.Println()
		}

		// Run security scan
		warnings := scanner.Scan(cfg)

		if !outputJSON {
			if len(warnings) == 0 {
				green.Println("   ✅ No security issues found")
			} else {
				bold.Printf("   ⚠️  Security Findings (%d):\n", len(warnings))
				for _, w := range warnings {
					var icon string
					var colorFn *color.Color
					switch w.Severity {
					case "critical":
						icon = "🔴"
						colorFn = red
					case "warning":
						icon = "🟡"
						colorFn = yellow
					default:
						icon = "🔵"
						colorFn = cyan
					}
					fmt.Printf("      %s ", icon)
					colorFn.Printf("[%s] ", w.Severity)
					fmt.Printf("%s", w.Message)
					if w.Server != "" {
						dim.Printf(" (server: %s)", w.Server)
					}
					fmt.Println()
					if verbose && w.Detail != "" {
						dim.Printf("         %s\n", w.Detail)
					}
				}
			}
			fmt.Println()
		}

		// Save snapshot if requested
		if saveSnapshot {
			snapshotDir := mcpguardDir()
			if err := os.MkdirAll(snapshotDir, 0755); err != nil {
				return fmt.Errorf("creating snapshot dir: %w", err)
			}
			snap := &diffpkg.Snapshot{
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				Source:    cfg.Source,
				Path:      cfg.Path,
				Hash:      configHash,
				Config:    cfg,
			}
			snapPath := filepath.Join(snapshotDir, fmt.Sprintf("%s-latest.json", cfg.Source))
			if err := diffpkg.SaveSnapshot(snap, snapPath); err != nil {
				return fmt.Errorf("saving snapshot: %w", err)
			}
			if !outputJSON {
				green.Printf("   💾 Snapshot saved: %s\n\n", snapPath)
			}
		}

		jsonOutput.Configs = append(jsonOutput.Configs, scanConfigOutput{
			Source:     cfg.Source,
			Path:       cfg.Path,
			FileHash:   fileHash,
			ConfigHash: configHash,
			Servers:    len(cfg.Servers),
			Warnings:   warnings,
		})
	}

	if outputJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(jsonOutput)
	}

	return nil
}

func sortedServerNames(cfg *config.MCPConfig) []string {
	names := make([]string, 0, len(cfg.Servers))
	for name := range cfg.Servers {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func mcpguardDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".mcpguard")
}

// printSources lists every supported parser, its default config path, and
// whether that path currently exists on disk. It honors --json so scripts
// can discover sources programmatically.
func printSources(bold, dim *color.Color) error {
	type sourceInfo struct {
		Name        string `json:"name"`
		DefaultPath string `json:"default_path"`
		Exists      bool   `json:"exists"`
	}

	parsers := config.AllParsers()
	out := make([]sourceInfo, 0, len(parsers))
	for _, p := range parsers {
		path := p.DefaultPath()
		exists := false
		if path != "" {
			if _, err := os.Stat(path); err == nil {
				exists = true
			}
		}
		out = append(out, sourceInfo{Name: p.Name(), DefaultPath: path, Exists: exists})
	}

	if outputJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(map[string]interface{}{"sources": out})
	}

	bold.Println("Supported MCP config sources:")
	for _, s := range out {
		marker := dim.Sprint("(not found)")
		if s.Exists {
			marker = color.GreenString("(found)")
		}
		fmt.Printf("  • %s — %s %s\n", s.Name, s.DefaultPath, marker)
	}
	fmt.Println()
	dim.Println("Use --source <name> to scan only one of these,")
	dim.Println("or --config <path> to scan an explicit file.")
	return nil
}
