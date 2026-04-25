package fingerprint

import (
	"regexp"
	"strings"

	"github.com/withaxiom/mcpguard/internal/config"
)

// sensitivePatterns matches common secret/token env var names.
var sensitivePatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(api[_-]?key)`),
	regexp.MustCompile(`(?i)(secret)`),
	regexp.MustCompile(`(?i)(token)`),
	regexp.MustCompile(`(?i)(password|passwd|pwd)`),
	regexp.MustCompile(`(?i)(credential)`),
	regexp.MustCompile(`(?i)(auth)`),
	regexp.MustCompile(`(?i)(private[_-]?key)`),
	regexp.MustCompile(`(?i)(access[_-]?key)`),
	regexp.MustCompile(`(?i)(connection[_-]?string)`),
	regexp.MustCompile(`(?i)(database[_-]?url)`),
}

// IsSensitiveKey returns true if an env var name looks like it holds a secret.
func IsSensitiveKey(key string) bool {
	for _, p := range sensitivePatterns {
		if p.MatchString(key) {
			return true
		}
	}
	return false
}

// RedactValue masks a sensitive value, showing only the first 4 chars.
func RedactValue(value string) string {
	if len(value) <= 4 {
		return "****"
	}
	return value[:4] + strings.Repeat("*", len(value)-4)
}

// RedactEnv returns a copy of env vars with sensitive values redacted.
func RedactEnv(env map[string]string) map[string]string {
	redacted := make(map[string]string, len(env))
	for k, v := range env {
		if IsSensitiveKey(k) {
			redacted[k] = RedactValue(v)
		} else {
			redacted[k] = v
		}
	}
	return redacted
}

// FindSensitiveEnvVars returns the names of env vars that look like secrets
// for each server in the config.
func FindSensitiveEnvVars(cfg *config.MCPConfig) map[string][]string {
	result := make(map[string][]string)
	for name, srv := range cfg.Servers {
		var sensitive []string
		for k := range srv.Env {
			if IsSensitiveKey(k) {
				sensitive = append(sensitive, k)
			}
		}
		if len(sensitive) > 0 {
			result[name] = sensitive
		}
	}
	return result
}
