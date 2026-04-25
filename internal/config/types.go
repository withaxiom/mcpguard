package config

// MCPConfig represents a parsed MCP configuration from any supported client.
type MCPConfig struct {
	Source  string              `json:"source"`  // e.g. "cursor", "claude-desktop"
	Path   string              `json:"path"`     // filesystem path to config file
	Servers map[string]MCPServer `json:"servers"` // server name → config
}

// MCPServer represents a single MCP server entry.
type MCPServer struct {
	Name        string            `json:"name"`
	Command     string            `json:"command"`
	Args        []string          `json:"args,omitempty"`
	Env         map[string]string `json:"env,omitempty"`
	URL         string            `json:"url,omitempty"`          // for HTTP-based servers
	Disabled    bool              `json:"disabled,omitempty"`
	AutoApprove []string          `json:"autoApprove,omitempty"` // auto-approved tool patterns
}

// ScanResult contains the output of scanning a config.
type ScanResult struct {
	Config      MCPConfig        `json:"config"`
	Fingerprint string           `json:"fingerprint"`  // SHA-256 of the config file
	RedactedEnv map[string][]string `json:"redacted_env"` // server → redacted env var names
	Warnings    []Warning        `json:"warnings,omitempty"`
}

// Warning represents a security concern found during scanning.
type Warning struct {
	Severity string `json:"severity"` // "info", "warning", "critical"
	Server   string `json:"server"`   // which server triggered it
	Message  string `json:"message"`
	Detail   string `json:"detail,omitempty"`
}

// DiffResult represents changes between two scan snapshots.
type DiffResult struct {
	Added    []string `json:"added,omitempty"`
	Removed  []string `json:"removed,omitempty"`
	Modified []DiffEntry `json:"modified,omitempty"`
}

// DiffEntry shows what changed in a specific server.
type DiffEntry struct {
	Server  string `json:"server"`
	Field   string `json:"field"`
	Before  string `json:"before"`
	After   string `json:"after"`
}
