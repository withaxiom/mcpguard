package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Parser is the interface all config parsers implement.
type Parser interface {
	Name() string
	DefaultPath() string
	Parse(path string) (*MCPConfig, error)
}

// AllParsers returns all registered config parsers.
func AllParsers() []Parser {
	return []Parser{
		&ClaudeDesktopParser{},
		&CursorParser{},
	}
}

// AutoDetect finds and parses all MCP configs on the system.
func AutoDetect() ([]*MCPConfig, error) {
	var configs []*MCPConfig
	for _, p := range AllParsers() {
		path := p.DefaultPath()
		if path == "" {
			continue
		}
		if _, err := os.Stat(path); err != nil {
			continue // config file doesn't exist
		}
		cfg, err := p.Parse(path)
		if err != nil {
			return nil, fmt.Errorf("parsing %s config at %s: %w", p.Name(), path, err)
		}
		configs = append(configs, cfg)
	}
	return configs, nil
}

// parseJSONFile reads and unmarshals a JSON file.
func parseJSONFile(path string, v interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading %s: %w", path, err)
	}
	// Strip BOM if present
	data = stripBOM(data)
	// Strip comments (// and /* */) for JSONC support
	data = stripJSONComments(data)
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("parsing JSON in %s: %w", path, err)
	}
	return nil
}

// stripBOM removes UTF-8 BOM if present.
func stripBOM(data []byte) []byte {
	if len(data) >= 3 && data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {
		return data[3:]
	}
	return data
}

// stripJSONComments removes // and /* */ comments from JSON (JSONC support).
func stripJSONComments(data []byte) []byte {
	var result []byte
	i := 0
	inString := false
	for i < len(data) {
		if inString {
			if data[i] == '\\' && i+1 < len(data) {
				result = append(result, data[i], data[i+1])
				i += 2
				continue
			}
			if data[i] == '"' {
				inString = false
			}
			result = append(result, data[i])
			i++
			continue
		}

		if data[i] == '"' {
			inString = true
			result = append(result, data[i])
			i++
			continue
		}

		// Line comment
		if i+1 < len(data) && data[i] == '/' && data[i+1] == '/' {
			for i < len(data) && data[i] != '\n' {
				i++
			}
			continue
		}

		// Block comment
		if i+1 < len(data) && data[i] == '/' && data[i+1] == '*' {
			i += 2
			for i+1 < len(data) && !(data[i] == '*' && data[i+1] == '/') {
				i++
			}
			if i+1 < len(data) {
				i += 2
			}
			continue
		}

		result = append(result, data[i])
		i++
	}
	return result
}

// expandHome expands ~ in paths.
func expandHome(path string) string {
	if strings.HasPrefix(path, "~/") || path == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(home, path[1:])
	}
	return path
}

// configDir returns the appropriate config directory for the OS.
func configDir() string {
	switch runtime.GOOS {
	case "darwin":
		home, _ := os.UserHomeDir()
		return filepath.Join(home, "Library", "Application Support")
	case "linux":
		if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
			return xdg
		}
		home, _ := os.UserHomeDir()
		return filepath.Join(home, ".config")
	case "windows":
		return os.Getenv("APPDATA")
	default:
		home, _ := os.UserHomeDir()
		return filepath.Join(home, ".config")
	}
}
