package cli

import (
	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/manifest"
	"os"
	"path/filepath"
	"testing"
)

func TestCheckpointReadBound(t *testing.T) {
	p := filepath.Join(t.TempDir(), "cp")
	f, e := os.Create(p)
	if e != nil {
		t.Fatal(e)
	}
	if e = f.Truncate(32<<20 + 1); e != nil {
		t.Fatal(e)
	}
	f.Close()
	if _, e = readCheckpoint(p); e == nil {
		t.Fatal("oversize read")
	}
}

func TestCheckpointPlanIDs(t *testing.T) {
	d := manifest.Document{Steps: []manifest.Step{{ID: "a"}}}
	if validateCheckpointPlan(checkpoint{Steps: map[string]string{"other": "unknown"}}, d) == nil {
		t.Fatal("foreign step")
	}
}

func TestCheckpointEvidenceOwnership(t *testing.T) {
	d := manifest.Document{Steps: []manifest.Step{{ID: "a"}}}
	for _, cp := range []checkpoint{{Steps: map[string]string{}, Results: map[string]any{"other": 1}}, {Steps: map[string]string{"a": "committed"}, Reconciliations: map[string]Reconciliation{"a": {}}}} {
		if validateCheckpointPlan(cp, d) == nil {
			t.Fatal(cp)
		}
	}
}
