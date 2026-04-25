package scanner

import (
	"fmt"
	"strings"

	"github.com/withaxiom/mcpguard/internal/config"
)

// SecurityRules defines what the scanner checks for.
var SecurityRules = []Rule{
	{Name: "exposed-secrets", Check: checkExposedSecrets},
	{Name: "wildcard-autoapprove", Check: checkWildcardAutoApprove},
	{Name: "shell-injection", Check: checkShellInjection},
	{Name: "suspicious-url", Check: checkSuspiciousURL},
	{Name: "disabled-servers", Check: checkDisabledServers},
}

// Rule is a security check that runs against a config.
type Rule struct {
	Name  string
	Check func(cfg *config.MCPConfig) []config.Warning
}

// Scan runs all security rules against a config.
func Scan(cfg *config.MCPConfig) []config.Warning {
	var warnings []config.Warning
	for _, rule := range SecurityRules {
		warnings = append(warnings, rule.Check(cfg)...)
	}
	return warnings
}

// checkExposedSecrets looks for API keys/tokens in env vars.
func checkExposedSecrets(cfg *config.MCPConfig) []config.Warning {
	var warnings []config.Warning
	sensitivePatterns := []string{"api_key", "apikey", "secret", "token", "password", "passwd", "credential", "private_key"}

	for name, srv := range cfg.Servers {
		for envKey := range srv.Env {
			lower := strings.ToLower(envKey)
			for _, pattern := range sensitivePatterns {
				if strings.Contains(lower, pattern) {
					warnings = append(warnings, config.Warning{
						Severity: "warning",
						Server:   name,
						Message:  fmt.Sprintf("Environment variable '%s' appears to contain a secret", envKey),
						Detail:   "Secrets in config files can be leaked through version control or backups. Consider using a secrets manager.",
					})
					break
				}
			}
		}
	}
	return warnings
}

// checkWildcardAutoApprove flags servers with broad auto-approve patterns.
func checkWildcardAutoApprove(cfg *config.MCPConfig) []config.Warning {
	var warnings []config.Warning
	for name, srv := range cfg.Servers {
		for _, pattern := range srv.AutoApprove {
			if pattern == "*" || pattern == "**" {
				warnings = append(warnings, config.Warning{
					Severity: "critical",
					Server:   name,
					Message:  fmt.Sprintf("Wildcard auto-approve pattern '%s' grants unrestricted tool access", pattern),
					Detail:   "This allows the MCP server to execute any tool without user confirmation. Remove or restrict the pattern.",
				})
			}
		}
	}
	return warnings
}

// checkShellInjection looks for potentially dangerous command patterns.
func checkShellInjection(cfg *config.MCPConfig) []config.Warning {
	var warnings []config.Warning

	// Patterns that almost always indicate shell composition, not data.
	// We deliberately exclude single '|' and ';' from arg checks because they
	// appear legitimately in URLs (query separators, MIME types) and produce
	// noisy false positives.
	commandPatterns := []string{"&&", "||", ";", "|", "`", "$(", "eval ", "exec "}
	argPatterns := []string{"&&", "||", "`", "$(", "eval ", "exec "}

	for name, srv := range cfg.Servers {
		// The command field should be a single executable (npx, node, python, ...).
		// Any shell metacharacter here is high-confidence injection.
		for _, pattern := range commandPatterns {
			if strings.Contains(srv.Command, pattern) {
				warnings = append(warnings, config.Warning{
					Severity: "critical",
					Server:   name,
					Message:  fmt.Sprintf("Command contains shell meta-character '%s' — possible injection risk", pattern),
					Detail:   fmt.Sprintf("Command: %s", srv.Command),
				})
				break
			}
		}
		// Args are passed as a slice to exec (no shell), so injection only
		// matters if the consumer re-evaluates them. Flag only the patterns
		// that strongly imply intent to compose shell commands.
		for _, arg := range srv.Args {
			if looksLikeURL(arg) {
				continue
			}
			for _, pattern := range argPatterns {
				if strings.Contains(arg, pattern) {
					warnings = append(warnings, config.Warning{
						Severity: "warning",
						Server:   name,
						Message:  fmt.Sprintf("Argument contains shell meta-character '%s'", pattern),
						Detail:   fmt.Sprintf("Arg: %s", arg),
					})
					break
				}
			}
		}
	}
	return warnings
}

// looksLikeURL returns true if the argument resembles a URL or connection
// string. URL query strings legitimately use '&' and ';' as separators,
// so we skip them to avoid false positives on database/API connection args.
func looksLikeURL(arg string) bool {
	lower := strings.ToLower(arg)
	prefixes := []string{
		"http://", "https://", "ws://", "wss://",
		"postgres://", "postgresql://", "mysql://", "mongodb://", "mongodb+srv://",
		"redis://", "rediss://", "amqp://", "amqps://", "file://",
	}
	for _, p := range prefixes {
		if strings.HasPrefix(lower, p) {
			return true
		}
	}
	return false
}

// checkSuspiciousURL flags non-localhost URLs.
func checkSuspiciousURL(cfg *config.MCPConfig) []config.Warning {
	var warnings []config.Warning
	for name, srv := range cfg.Servers {
		if srv.URL == "" {
			continue
		}
		lower := strings.ToLower(srv.URL)
		if !strings.Contains(lower, "localhost") && !strings.Contains(lower, "127.0.0.1") && !strings.Contains(lower, "::1") {
			warnings = append(warnings, config.Warning{
				Severity: "info",
				Server:   name,
				Message:  "Server connects to a remote URL",
				Detail:   fmt.Sprintf("URL: %s — Verify this is a trusted endpoint.", srv.URL),
			})
		}
		if strings.HasPrefix(lower, "http://") {
			warnings = append(warnings, config.Warning{
				Severity: "warning",
				Server:   name,
				Message:  "Server uses unencrypted HTTP connection",
				Detail:   fmt.Sprintf("URL: %s — Consider using HTTPS.", srv.URL),
			})
		}
	}
	return warnings
}

// checkDisabledServers notes disabled servers (informational).
func checkDisabledServers(cfg *config.MCPConfig) []config.Warning {
	var warnings []config.Warning
	for name, srv := range cfg.Servers {
		if srv.Disabled {
			warnings = append(warnings, config.Warning{
				Severity: "info",
				Server:   name,
				Message:  "Server is disabled",
				Detail:   "Consider removing disabled servers to reduce config surface area.",
			})
		}
	}
	return warnings
}
