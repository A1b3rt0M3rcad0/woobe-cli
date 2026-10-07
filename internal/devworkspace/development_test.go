package devworkspace

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/packagebundle"
	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/packagefmt"
)

func completeGraph(t *testing.T) (*Config, *Graph) {
	t.Helper()
	data, err := os.ReadFile("../../testdata/package/shared/complete.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixture struct {
		Files map[string]string `json:"files"`
	}
	if err = json.Unmarshal(data, &fixture); err != nil {
		t.Fatal(err)
	}
	c, err := Create(t.TempDir(), ".woobe", "")
	if err != nil {
		t.Fatal(err)
	}
	os.MkdirAll(c.RootPath(), 0700)
	for path, content := range fixture.Files {
		target := filepath.Join(c.RootPath(), filepath.FromSlash(path))
		os.MkdirAll(filepath.Dir(target), 0700)
		os.WriteFile(target, []byte(content), 0600)
	}
	bundle, err := packagebundle.Load(c.RootPath(), false)
	if err != nil {
		t.Fatal(err)
	}
	defer bundle.Close()
	for key, document := range bundle.Graph.Components {
		uid, _ := NewID()
		c.Resources = append(c.Resources, Resource{UID: uid, Kind: packagefmt.Text(document["kind"]), Key: key, Alias: key, Path: bundle.Graph.Paths[key]})
	}
	g, err := LoadGraph(c)
	if err != nil {
		t.Fatal(err)
	}
	return c, g
}
func TestCompilerKeepsClosureSupportTypesAndFrozenNodes(t *testing.T) {
	_, g := completeGraph(t)
	req := map[string]any{"credentials": []any{map[string]any{"ref": "primary", "provider": "custom"}}, "project_environment": []any{map[string]any{"key": "PACKAGE_TEST_TOKEN", "kind": "secret"}}, "secrets": []any{map[string]any{"ref": "package-tool-token"}}}
	bundle, _, err := g.Compile("support-network", req, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer bundle.Close()
	if len(bundle.Graph.Components) != len(g.Nodes) {
		t.Fatal("closure dropped components", len(bundle.Graph.Components), len(g.Nodes))
	}
	if !reflect.DeepEqual(packagefmt.Object(bundle.Graph.Components["support-network"]["spec"])["nodes"], packagefmt.Object(g.Nodes["support-network"].Document["spec"])["nodes"]) {
		t.Fatal("Network node identities changed")
	}
	if _, ok := packagefmt.Object(packagefmt.Object(bundle.Graph.Components["specialist"]["spec"])["model"])["parameters"].(map[string]any); !ok {
		t.Fatal("model parameters lost")
	}
	for _, file := range bundle.Inventory {
		if strings.HasPrefix(file.Path, ".state") {
			t.Fatal("private state exported")
		}
	}
}
func TestThreeWayMergeHandlesIndependentEditsNullDeletionAndAtomicLists(t *testing.T) {
	base := map[string]any{"spec": map[string]any{"name": "before", "description": "before", "nullable": nil, "nodes": []any{"a", "b"}}}
	local := clone(base)
	remote := clone(base)
	packagefmt.Object(local["spec"])["name"] = "local"
	packagefmt.Object(remote["spec"])["description"] = "remote"
	merged, conflicts := Merge(base, local, remote)
	if len(conflicts) != 0 || packagefmt.Object(merged["spec"])["name"] != "local" || packagefmt.Object(merged["spec"])["description"] != "remote" {
		t.Fatal(merged, conflicts)
	}
	delete(packagefmt.Object(local["spec"]), "nullable")
	packagefmt.Object(remote["spec"])["nullable"] = "value"
	_, conflicts = Merge(base, local, remote)
	if len(conflicts) != 1 || conflicts[0].Path != "/spec/nullable" {
		t.Fatal(conflicts)
	}
	local = clone(base)
	remote = clone(base)
	packagefmt.Object(local["spec"])["nodes"] = []any{"a", "b", "c"}
	packagefmt.Object(remote["spec"])["nodes"] = []any{"b", "a"}
	_, conflicts = Merge(base, local, remote)
	if len(conflicts) != 1 || conflicts[0].Path != "/spec/nodes" {
		t.Fatal(conflicts)
	}
}
func TestStateScopesConnectionsAndPreservesNumericRevisions(t *testing.T) {
	c, _ := Create(t.TempDir(), ".woobe", "")
	a, _ := c.ReadState("https://one", "workspace", "project")
	a.Bindings["agent"] = Binding{ResourceID: "native", Revision: json.Number("7"), Base: map[string]any{"temperature": json.Number("0.23")}}
	if err := c.WriteState(a); err != nil {
		t.Fatal(err)
	}
	loaded, err := c.ReadState("https://one", "workspace", "project")
	if err != nil || loaded.Bindings["agent"].Revision != json.Number("7") {
		t.Fatal(loaded, err)
	}
	other, err := c.ReadState("https://two", "workspace", "project")
	if err != nil || len(other.Bindings) != 0 {
		t.Fatal("cross-connection binding leak", other, err)
	}
}
func TestCompilerRejectsEscapedAndSymlinkSupport(t *testing.T) {
	c, g := completeGraph(t)
	spec := packagefmt.Object(g.Nodes["support-skill"].Document["spec"])
	packagefmt.Object(spec["package"])["manifest"] = "../../outside.md"
	if b, _, err := g.Compile("support", nil, nil); err == nil {
		b.Close()
		t.Fatal("escaped support accepted")
	}
	c, g = completeGraph(t)
	spec = packagefmt.Object(g.Nodes["support-skill"].Document["spec"])
	manifest := filepath.Join(c.RootPath(), filepath.FromSlash(packagefmt.Text(packagefmt.Object(spec["package"])["manifest"])))
	outside := filepath.Join(t.TempDir(), "SKILL.md")
	if err := os.WriteFile(outside, []byte("outside"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(manifest); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, manifest); err != nil {
		t.Skip("symlink creation unavailable")
	}
	if b, _, err := g.Compile("support", nil, nil); err == nil {
		b.Close()
		t.Fatal("linked support accepted")
	}
}

func TestLocalTransactionRecoversInterruptedWritesAndPreservesExternalEdits(t *testing.T) {
	c, err := Create(t.TempDir(), ".woobe", "")
	if err != nil {
		t.Fatal(err)
	}
	a := filepath.Join(c.RootPath(), "agents/support/agent.yaml")
	b := filepath.Join(c.RootPath(), ".state/destination.json")
	if err := c.Commit(map[string][]byte{a: []byte("old"), b: []byte("state-old")}); err != nil {
		t.Fatal(err)
	}
	journal := localJournal{RegistryID: c.RegistryID, Changes: []localChange{{Path: a, Before: []byte("old"), Exists: true, After: []byte("new")}, {Path: b, Before: []byte("state-old"), Exists: true, After: []byte("state-new")}}}
	data, _ := json.Marshal(journal)
	if err := privateWrite(c.journalPath(), data); err != nil {
		t.Fatal(err)
	}
	if err := privateWrite(a, []byte("new")); err != nil {
		t.Fatal(err)
	}
	if err := c.Recover(); err != nil {
		t.Fatal(err)
	}
	actual, _ := os.ReadFile(b)
	if string(actual) != "state-new" {
		t.Fatal("interrupted state was not recovered")
	}
	if _, err := os.Stat(c.journalPath()); !os.IsNotExist(err) {
		t.Fatal("completed journal retained")
	}
	if err := privateWrite(c.journalPath(), data); err != nil {
		t.Fatal(err)
	}
	if err := privateWrite(a, []byte("external edit")); err != nil {
		t.Fatal(err)
	}
	if err := c.Recover(); err == nil {
		t.Fatal("external edit overwritten during recovery")
	}
	actual, _ = os.ReadFile(a)
	if string(actual) != "external edit" {
		t.Fatal("external edit lost")
	}
}

func TestDevelopmentLockSerializesProcesses(t *testing.T) {
	c, _ := Create(t.TempDir(), ".woobe", "")
	unlock, err := c.Lock()
	if err != nil {
		t.Fatal(err)
	}
	defer unlock()
	if second, err := c.Lock(); err == nil {
		second()
		t.Fatal("another command acquired the same registry lock")
	}
}

func TestManagedCaptureReusesRegistryAndConflictsWithoutWriting(t *testing.T) {
	source, _ := completeGraph(t)
	bundle, err := packagebundle.Load(source.RootPath(), false)
	if err != nil {
		t.Fatal(err)
	}
	defer bundle.Close()
	target, _ := Create(t.TempDir(), ".woobe", "")
	os.MkdirAll(target.RootPath(), 0700)
	state, _ := target.ReadState("http://local", "workspace", "project")
	captured := map[string]CapturedBinding{}
	for key, d := range bundle.Graph.Components {
		id, _ := NewID()
		captured[key] = CapturedBinding{ResourceID: id, Revision: json.Number("1"), SourceKind: "draft"}
		_ = d
	}
	credential, _ := NewID()
	resource, conflicts, err := target.ImportCapture(bundle, state, captured, map[string]string{"primary": credential}, "support", "")
	if err != nil || len(conflicts) > 0 {
		t.Fatal(resource, conflicts, err)
	}
	count := len(target.Resources)
	if _, conflicts, err = target.ImportCapture(bundle, state, captured, map[string]string{"primary": credential}, "", ""); err != nil || len(conflicts) > 0 || len(target.Resources) != count {
		t.Fatal("duplicated capture", conflicts, err, len(target.Resources), count)
	}
	graph, err := LoadGraph(target)
	if err != nil {
		t.Fatal(err)
	}
	node := graph.Nodes[resource.Key]
	document := clone(node.Document)
	packagefmt.Object(document["metadata"])["description"] = "local"
	data, _ := Encode(document)
	file := graph.Path(node)
	os.WriteFile(file, data, 0600)
	packagefmt.Object(bundle.Graph.Components["support-network"]["metadata"])["description"] = "remote"
	original, _ := os.ReadFile(file)
	if _, conflicts, err = target.ImportCapture(bundle, state, captured, map[string]string{"primary": credential}, "", ""); err != nil || len(conflicts) != 1 {
		t.Fatal(conflicts, err)
	}
	after, _ := os.ReadFile(file)
	if !reflect.DeepEqual(original, after) {
		t.Fatal("conflicting pull changed author files")
	}
}
func TestPrivateStateRefusesLinkedParents(t *testing.T) {
	c, _ := Create(t.TempDir(), ".woobe", "")
	os.MkdirAll(c.RootPath(), 0700)
	if err := os.Symlink(t.TempDir(), filepath.Join(c.RootPath(), ".state")); err != nil {
		t.Skip("platform does not allow symlink creation")
	}
	if _, err := c.ReadState("http://local", "workspace", "project"); err == nil {
		t.Fatal("linked private state accepted")
	}
}
