package diff

import (
	"testing"

	"github.com/withaxiom/mcpguard/internal/config"
)

func TestConfigs_NoChanges(t *testing.T) {
	cfg := &config.MCPConfig{
		Servers: map[string]config.MCPServer{
			"test": {Name: "test", Command: "node"},
		},
	}

	result := Configs(cfg, cfg)
	if result.HasChanges() {
		t.Error("expected no changes when comparing same config")
	}
}

func TestConfigs_AddedServer(t *testing.T) {
	before := &config.MCPConfig{
		Servers: map[string]config.MCPServer{
			"existing": {Name: "existing", Command: "node"},
		},
	}
	after := &config.MCPConfig{
		Servers: map[string]config.MCPServer{
			"existing": {Name: "existing", Command: "node"},
			"new":      {Name: "new", Command: "python"},
		},
	}

	result := Configs(before, after)
	if len(result.Added) != 1 || result.Added[0] != "new" {
		t.Errorf("expected 'new' in added, got %v", result.Added)
	}
}

func TestConfigs_RemovedServer(t *testing.T) {
	before := &config.MCPConfig{
		Servers: map[string]config.MCPServer{
			"existing": {Name: "existing", Command: "node"},
			"removed":  {Name: "removed", Command: "python"},
		},
	}
	after := &config.MCPConfig{
		Servers: map[string]config.MCPServer{
			"existing": {Name: "existing", Command: "node"},
		},
	}

	result := Configs(before, after)
	if len(result.Removed) != 1 || result.Removed[0] != "removed" {
		t.Errorf("expected 'removed' in removed, got %v", result.Removed)
	}
}

func TestConfigs_ModifiedCommand(t *testing.T) {
	before := &config.MCPConfig{
		Servers: map[string]config.MCPServer{
			"test": {Name: "test", Command: "node"},
		},
	}
	after := &config.MCPConfig{
		Servers: map[string]config.MCPServer{
			"test": {Name: "test", Command: "python"},
		},
	}

	result := Configs(before, after)
	if len(result.Modified) != 1 {
		t.Fatalf("expected 1 modification, got %d", len(result.Modified))
	}
	if result.Modified[0].Field != "command" {
		t.Errorf("expected field 'command', got '%s'", result.Modified[0].Field)
	}
	if result.Modified[0].Before != "node" || result.Modified[0].After != "python" {
		t.Errorf("unexpected values: before=%s after=%s", result.Modified[0].Before, result.Modified[0].After)
	}
}

func TestConfigs_EnvChange(t *testing.T) {
	before := &config.MCPConfig{
		Servers: map[string]config.MCPServer{
			"test": {
				Name:    "test",
				Command: "node",
				Env:     map[string]string{"KEY": "old"},
			},
		},
	}
	after := &config.MCPConfig{
		Servers: map[string]config.MCPServer{
			"test": {
				Name:    "test",
				Command: "node",
				Env:     map[string]string{"KEY": "new"},
			},
		},
	}

	result := Configs(before, after)
	if len(result.Modified) != 1 {
		t.Fatalf("expected 1 modification, got %d", len(result.Modified))
	}
	if result.Modified[0].Field != "env.KEY" {
		t.Errorf("expected field 'env.KEY', got '%s'", result.Modified[0].Field)
	}
}
