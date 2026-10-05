package cli

import (
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
