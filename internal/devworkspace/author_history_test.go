package devworkspace

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func authorState() *State {
	return &State{Requirements: map[string]any{"credentials": []any{map[string]any{"ref": "primary", "provider": "custom"}}, "project_environment": []any{map[string]any{"key": "PACKAGE_TEST_TOKEN", "kind": "secret"}}, "secrets": []any{map[string]any{"ref": "package-tool-token"}}}, Bindings: map[string]Binding{}}
}

func TestCheckoutRestoresAuthorBytesWithoutChangingOrigin(t *testing.T) {
	c, g := completeGraph(t)
	r := g.Nodes["support"].Resource
	state := authorState()
	file := filepath.Join(c.RootPath(), g.Nodes[r.Key].Descriptor)
	original, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	original = append([]byte("# preserved author comment\n"), original...)
	if err = os.WriteFile(file, original, 0600); err != nil {
		t.Fatal(err)
	}
	g, err = LoadGraph(c)
	if err != nil {
		t.Fatal(err)
	}
	first, err := g.CreateRevision(r, state, "before", nil)
	if err != nil {
		t.Fatal(err)
	}
	changed := append([]byte("# second comment\n"), original...)
	if err = os.WriteFile(file, changed, 0600); err != nil {
		t.Fatal(err)
	}
	g, err = LoadGraph(c)
	if err != nil {
		t.Fatal(err)
	}
	second, err := g.CreateRevision(r, state, "after", nil)
	if err != nil {
		t.Fatal(err)
	}
	if first.AuthorDigest == second.AuthorDigest || first.DefinitionDigest != second.DefinitionDigest {
		t.Fatal("source comments and executable identities conflated")
	}
	before, err := c.ReadTracking(r)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = g.CheckoutRevision(r, first); err != nil {
		t.Fatal(err)
	}
	got, _ := os.ReadFile(file)
	if string(got) != string(original) {
		t.Fatal("author bytes not restored")
	}
	after, err := c.ReadTracking(r)
	if err != nil || after.Origin != before.Origin || after.Working != first.ID {
		t.Fatal(after, err)
	}
	// Even comment-only edits are protected, without a force bypass.
	if err = os.WriteFile(file, append(got, []byte("\n# uncheckpointed\n")...), 0600); err != nil {
		t.Fatal(err)
	}
	g, err = LoadGraph(c)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = g.CheckoutRevision(r, second); err == nil || !strings.Contains(err.Error(), "local author edits") {
		t.Fatal("uncheckpointed edits accepted", err)
	}
}

func TestCheckoutRejectsTamperedSourceBeforeWriting(t *testing.T) {
	c, g := completeGraph(t)
	r := g.Nodes["support"].Resource
	record, err := g.CreateRevision(r, authorState(), "before", nil)
	if err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(c.RootPath(), "objects", "author", strings.TrimPrefix(record.AuthorDigest, "sha256:")+".json")
	if err = os.WriteFile(file, []byte(`{}`), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err = g.CheckoutRevision(r, record); err == nil {
		t.Fatal("tampered author object accepted")
	}
}

func TestNetworkCheckoutRestoresEntireClosedAuthorGraph(t *testing.T) {
	c, g := completeGraph(t)
	r := g.Nodes["support-network"].Resource
	state := authorState()
	first, err := g.CreateRevision(r, state, "network original", nil)
	if err != nil {
		t.Fatal(err)
	}
	node := g.Nodes["support"]
	file := filepath.Join(c.RootPath(), node.Descriptor)
	original, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	changed := append([]byte("# network dependency edit\n"), original...)
	if err = os.WriteFile(file, changed, 0600); err != nil {
		t.Fatal(err)
	}
	g, err = LoadGraph(c)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = g.CreateRevision(r, state, "network dependency", nil); err != nil {
		t.Fatal(err)
	}
	if _, err = g.CheckoutRevision(r, first); err != nil {
		t.Fatal(err)
	}
	got, _ := os.ReadFile(file)
	if string(got) != string(original) {
		t.Fatal("Network dependency source not restored")
	}
}

func TestHistoryDistinguishesMissingAuthorFromMissingExecutable(t *testing.T) {
	c, g := completeGraph(t)
	r := g.Nodes["support"].Resource
	record, err := g.CreateRevision(r, authorState(), "source", nil)
	if err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(c.RootPath(), "objects", "author", strings.TrimPrefix(record.AuthorDigest, "sha256:")+".json")
	if err = os.Remove(file); err != nil {
		t.Fatal(err)
	}
	report, err := c.VerifyHistory(r)
	if err != nil {
		t.Fatal(err)
	}
	if report["objects_available"] != true || report["author_sources_available"] != false || len(report["missing_author_objects"].([]string)) != 1 {
		t.Fatal(report)
	}
}
