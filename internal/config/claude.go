package config

import (
	"path/filepath"
)

// ClaudeDesktopParser handles Claude Desktop's MCP configuration.
type ClaudeDesktopParser struct{}

func (p *ClaudeDesktopParser) Name() string {
	return "claude-desktop"
}

func (p *ClaudeDesktopParser) DefaultPath() string {
	return filepath.Join(configDir(), "Claude", "claude_desktop_config.json")
}

// claudeDesktopRaw represents the raw JSON structure of Claude Desktop config.
type claudeDesktopRaw struct {
	MCPServers map[string]claudeServerRaw `json:"mcpServers"`
}

type claudeServerRaw struct {
	Command     string            `json:"command"`
	Args        []string          `json:"args"`
	Env         map[string]string `json:"env"`
	URL         string            `json:"url"`
	Disabled    bool              `json:"disabled"`
	AutoApprove []string          `json:"autoApprove"`
}

func (p *ClaudeDesktopParser) Parse(path string) (*MCPConfig, error) {
	var raw claudeDesktopRaw
	if err := parseJSONFile(path, &raw); err != nil {
		return nil, err
	}

	servers := make(map[string]MCPServer, len(raw.MCPServers))
	for name, srv := range raw.MCPServers {
		servers[name] = MCPServer{
			Name:        name,
			Command:     srv.Command,
			Args:        srv.Args,
			Env:         srv.Env,
			URL:         srv.URL,
			Disabled:    srv.Disabled,
			AutoApprove: srv.AutoApprove,
		}
	}

	return &MCPConfig{
		Source:  p.Name(),
		Path:    path,
		Servers: servers,
	}, nil
}
