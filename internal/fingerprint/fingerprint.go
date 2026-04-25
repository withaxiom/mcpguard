package fingerprint

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/withaxiom/mcpguard/internal/config"
)

// FileHash computes SHA-256 of a file on disk.
func FileHash(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("reading file for fingerprint: %w", err)
	}
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:]), nil
}

// ConfigHash computes a deterministic SHA-256 of the parsed config content.
// This normalizes the config so formatting changes don't affect the hash.
func ConfigHash(cfg *config.MCPConfig) (string, error) {
	// Build a deterministic representation
	type serverEntry struct {
		Name    string   `json:"name"`
		Command string   `json:"command"`
		Args    []string `json:"args"`
		URL     string   `json:"url,omitempty"`
	}

	var entries []serverEntry
	for name, srv := range cfg.Servers {
		entries = append(entries, serverEntry{
			Name:    name,
			Command: srv.Command,
			Args:    srv.Args,
			URL:     srv.URL,
		})
	}
	// Sort for determinism
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name < entries[j].Name
	})

	data, err := json.Marshal(entries)
	if err != nil {
		return "", fmt.Errorf("marshaling config for fingerprint: %w", err)
	}

	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:]), nil
}
