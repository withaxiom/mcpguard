package config

import (
	"os"
	"path/filepath"
)

// CursorParser handles Cursor IDE's MCP configuration.
type CursorParser struct{}

func (p *CursorParser) Name() string {
	return "cursor"
}

func (p *CursorParser) DefaultPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	// Cursor stores MCP config in ~/.cursor/mcp.json
	return filepath.Join(home, ".cursor", "mcp.json")
}

// cursorRaw represents the raw JSON structure of Cursor's MCP config.
type cursorRaw struct {
	MCPServers map[string]cursorServerRaw `json:"mcpServers"`
}

type cursorServerRaw struct {
	Command  string            `json:"command"`
	Args     []string          `json:"args"`
	Env      map[string]string `json:"env"`
	URL      string            `json:"url"`
	Disabled bool              `json:"disabled"`
}

func (p *CursorParser) Parse(path string) (*MCPConfig, error) {
	var raw cursorRaw
	if err := parseJSONFile(path, &raw); err != nil {
		return nil, err
	}

	servers := make(map[string]MCPServer, len(raw.MCPServers))
	for name, srv := range raw.MCPServers {
		servers[name] = MCPServer{
			Name:     name,
			Command:  srv.Command,
			Args:     srv.Args,
			Env:      srv.Env,
			URL:      srv.URL,
			Disabled: srv.Disabled,
		}
	}

	return &MCPConfig{
		Source:  p.Name(),
		Path:    path,
		Servers: servers,
	}, nil
}
