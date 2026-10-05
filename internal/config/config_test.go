package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func loadText(t *testing.T, s string) (Config, error) {
	t.Helper()
	p := filepath.Join(t.TempDir(), "config.json")
	if e := os.WriteFile(p, []byte(s), 0600); e != nil {
		t.Fatal(e)
	}
	return Load(p)
}
func TestRejectAmbiguousConfig(t *testing.T) {
	for _, s := range []string{`{"version":1,"version":1}`, `{"version":1,"contexts":{"a":{"api_url":"http://localhost","api_url":"http://other"}}}`, `{"version":1} {}`} {
		if _, e := loadText(t, s); e == nil {
			t.Fatal(s)
		}
	}
}

func TestConfigSizeLimit(t *testing.T) {
	if _, e := loadText(t, strings.Repeat(" ", 1<<20+1)); e == nil {
		t.Fatal("oversize accepted")
	}
}

func TestUnknownConfigFields(t *testing.T) {
	for _, s := range []string{`{"version":1,"credentail":"wrong"}`, `{"version":1,"contexts":{"a":{"api_url":"http://localhost","project":"wrong"}}}`} {
		if _, e := loadText(t, s); e == nil {
			t.Fatal(s)
		}
	}
}
