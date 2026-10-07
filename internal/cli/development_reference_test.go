package cli

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/devworkspace"
)

func TestCanonicalAndReferenceFirstAliasesUseScopedNativeBindings(t *testing.T) {
	t.Chdir(t.TempDir())
	t.Setenv("WOOBE_CONTROL_KEY", "test-control")
	requests := []string{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r.URL.Path)
		_ = json.NewEncoder(w).Encode(map[string]any{"success": true, "data": []any{}})
	}))
	defer server.Close()
	c, err := devworkspace.Create(".", ".woobe", "")
	if err != nil {
		t.Fatal(err)
	}
	uid, _ := devworkspace.NewID()
	c.Resources = []devworkspace.Resource{{UID: uid, Kind: "Agent", Key: "support", Alias: "support", Path: "agents/support"}}
	if err := c.Save(); err != nil {
		t.Fatal(err)
	}
	state, err := c.ReadState(server.URL, "workspace", "project")
	if err != nil {
		t.Fatal(err)
	}
	state.Bindings[uid] = devworkspace.Binding{ResourceID: "native-agent"}
	if err := c.WriteState(state); err != nil {
		t.Fatal(err)
	}
	for _, command := range [][]string{{"agent", "@support", "release", "list"}, {"agent", "release", "list", "@support"}, {"agent", "@support", "get"}} {
		args := append([]string{"--api-url", server.URL, "--workspace", "workspace", "--project", "project", "--output", "json"}, command...)
		code, response := invoke(t, args, "")
		if code != 0 {
			t.Fatal(command, code, response)
		}
	}
	if len(requests) != 3 || requests[0] != "/ai/agents/native-agent/releases" || requests[1] != requests[0] || requests[2] != "/ai/agents/native-agent" {
		t.Fatal(requests)
	}
	code, _ := invoke(t, []string{"--api-url", server.URL, "--workspace", "workspace", "--project", "other-project", "agent", "@support", "get", "--output", "json"}, "")
	if code != 2 || len(requests) != 3 {
		t.Fatal("cross-project alias sent a request", code, requests)
	}
}
