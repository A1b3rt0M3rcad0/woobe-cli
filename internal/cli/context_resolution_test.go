package cli

import (
	"bytes"
	"context"
	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/config"
	"path/filepath"
	"testing"
)

func contextApp(t *testing.T) *App {
	t.Helper()
	a := New(&bytes.Buffer{}, &bytes.Buffer{}, &bytes.Buffer{})
	a.ConfigPath = filepath.Join(t.TempDir(), "config.json")
	if e := config.Save(a.ConfigPath, config.Config{Version: 1, Current: "dev", Contexts: map[string]config.Context{"dev": {APIURL: "http://localhost", Workspace: "w1", Project: "p1"}}}); e != nil {
		t.Fatal(e)
	}
	return a
}
func TestWorkspaceOverrideClearsInheritedProject(t *testing.T) {
	for _, env := range []bool{false, true} {
		a := contextApp(t)
		if env {
			t.Setenv("WOOBE_WORKSPACE_ID", "w2")
		} else {
			a.Workspace = "w2"
		}
		v, e := a.resolve()
		if e != nil || v.Project != "" || v.Workspace != "w2" {
			t.Fatal(v, e)
		}
	}
	a := contextApp(t)
	a.Workspace = "w2"
	a.Project = "p2"
	v, e := a.resolve()
	if e != nil || v.Project != "p2" {
		t.Fatal(v, e)
	}
}
