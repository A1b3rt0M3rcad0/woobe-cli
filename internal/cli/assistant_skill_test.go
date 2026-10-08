package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEmbeddedAssistantSkillOfflineAliasesAndRuntimeBoundary(t *testing.T) {
	root := t.TempDir()
	for _, action := range []string{"install", "status", "uninstall"} {
		code, v := invoke(t, []string{"skill", action, "--agent", "codex,claude", "--project-dir", root, "--output", "json"}, "")
		if code != 0 || v["success"] != true {
			t.Fatal(action, v)
		}
	}
	if _, err := os.Stat(filepath.Join(root, ".agents/skills/woobe-cli")); !os.IsNotExist(err) {
		t.Fatal("uninstall did not remove skill", err)
	}
	code, v := invoke(t, []string{"--dry-run", "skill", "install", "--project-dir", root, "--dry-run", "--output", "json"}, "")
	if code != 0 || v["data"].(map[string]any)["executed"] != false {
		t.Fatal(v)
	}
	code, v = invoke(t, []string{"skills", "agents", "--output", "json"}, "")
	if code != 0 || len(v["data"].([]any)) != 6 {
		t.Fatal(v)
	}
	code, v = invoke(t, []string{"skill", "install", "--project", "not-a-directory", "--output", "json"}, "")
	if code != 2 {
		t.Fatal("Project ID used for skill install", v)
	}
	// Runtime Skills remain their canonical HTTP operations, never redirected.
	a := New(nil, nil, nil)
	if op, ok := a.operation("project skill list"); !ok || op.Kind != "http" {
		t.Fatal("runtime Skill route changed", op, ok)
	}
}
