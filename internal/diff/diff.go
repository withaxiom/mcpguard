package diff

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/withaxiom/mcpguard/internal/config"
)

// Result represents the differences between two config snapshots.
type Result struct {
	Added    []string `json:"added,omitempty"`
	Removed  []string `json:"removed,omitempty"`
	Modified []Entry  `json:"modified,omitempty"`
}

// Entry shows what changed in a specific server.
type Entry struct {
	Server string `json:"server"`
	Field  string `json:"field"`
	Before string `json:"before"`
	After  string `json:"after"`
}

// Configs compares two MCPConfig structs and returns the differences.
func Configs(before, after *config.MCPConfig) *Result {
	result := &Result{}

	// Find added and modified servers
	for name, afterSrv := range after.Servers {
		beforeSrv, exists := before.Servers[name]
		if !exists {
			result.Added = append(result.Added, name)
			continue
		}
		// Check for modifications
		entries := compareServers(name, beforeSrv, afterSrv)
		result.Modified = append(result.Modified, entries...)
	}

	// Find removed servers
	for name := range before.Servers {
		if _, exists := after.Servers[name]; !exists {
			result.Removed = append(result.Removed, name)
		}
	}

	// Sort for deterministic output
	sort.Strings(result.Added)
	sort.Strings(result.Removed)
	sort.Slice(result.Modified, func(i, j int) bool {
		if result.Modified[i].Server == result.Modified[j].Server {
			return result.Modified[i].Field < result.Modified[j].Field
		}
		return result.Modified[i].Server < result.Modified[j].Server
	})

	return result
}

// compareServers finds field-level differences between two server configs.
func compareServers(name string, before, after config.MCPServer) []Entry {
	var entries []Entry

	if before.Command != after.Command {
		entries = append(entries, Entry{
			Server: name,
			Field:  "command",
			Before: before.Command,
			After:  after.Command,
		})
	}

	beforeArgs := strings.Join(before.Args, " ")
	afterArgs := strings.Join(after.Args, " ")
	if beforeArgs != afterArgs {
		entries = append(entries, Entry{
			Server: name,
			Field:  "args",
			Before: beforeArgs,
			After:  afterArgs,
		})
	}

	if before.URL != after.URL {
		entries = append(entries, Entry{
			Server: name,
			Field:  "url",
			Before: before.URL,
			After:  after.URL,
		})
	}

	if before.Disabled != after.Disabled {
		entries = append(entries, Entry{
			Server: name,
			Field:  "disabled",
			Before: fmt.Sprintf("%v", before.Disabled),
			After:  fmt.Sprintf("%v", after.Disabled),
		})
	}

	// Compare env vars
	entries = append(entries, compareEnv(name, before.Env, after.Env)...)

	return entries
}

// compareEnv detects changes in environment variables.
func compareEnv(serverName string, before, after map[string]string) []Entry {
	var entries []Entry
	allKeys := make(map[string]bool)
	for k := range before {
		allKeys[k] = true
	}
	for k := range after {
		allKeys[k] = true
	}

	for k := range allKeys {
		bv, bOk := before[k]
		av, aOk := after[k]
		if bOk && aOk && bv != av {
			entries = append(entries, Entry{
				Server: serverName,
				Field:  fmt.Sprintf("env.%s", k),
				Before: bv,
				After:  av,
			})
		} else if !bOk && aOk {
			entries = append(entries, Entry{
				Server: serverName,
				Field:  fmt.Sprintf("env.%s", k),
				Before: "(not set)",
				After:  av,
			})
		} else if bOk && !aOk {
			entries = append(entries, Entry{
				Server: serverName,
				Field:  fmt.Sprintf("env.%s", k),
				Before: bv,
				After:  "(removed)",
			})
		}
	}
	return entries
}

// HasChanges returns true if there are any differences.
func (r *Result) HasChanges() bool {
	return len(r.Added) > 0 || len(r.Removed) > 0 || len(r.Modified) > 0
}

// Snapshot represents a saved config state for later diffing.
type Snapshot struct {
	Timestamp string         `json:"timestamp"`
	Source    string         `json:"source"`
	Path     string         `json:"path"`
	Hash     string         `json:"hash"`
	Config   *config.MCPConfig `json:"config"`
}

// SaveSnapshot writes a config snapshot to a JSON file.
func SaveSnapshot(snapshot *Snapshot, outputPath string) error {
	data, err := json.MarshalIndent(snapshot, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling snapshot: %w", err)
	}
	return os.WriteFile(outputPath, data, 0644)
}

// LoadSnapshot reads a config snapshot from a JSON file.
func LoadSnapshot(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading snapshot: %w", err)
	}
	var snapshot Snapshot
	if err := json.Unmarshal(data, &snapshot); err != nil {
		return nil, fmt.Errorf("parsing snapshot: %w", err)
	}
	return &snapshot, nil
}
