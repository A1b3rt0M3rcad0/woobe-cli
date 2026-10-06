package cli

import "testing"

func TestCategoryManifestOfflineCompilationAndScope(t *testing.T) {
	code, v := invoke(t, []string{"manifest", "compile", "--file", "-"}, `{"schema_version":"2","workspace_id":"w","resources":[{"key":"reader","kind":"AuthorityCategory","action":"create","spec":{"name":"reader","permissions":["agent:read"]}}]}`)
	if code != 0 || v["data"].(map[string]any)["executed"] != false {
		t.Fatal(code, v)
	}
	for _, command := range []string{"workspace authority category clone", "workspace authority category archive"} {
		code, _ := invoke(t, []string{"manifest", "validate", "--file", "-"}, `{"schema_version":"1","workspace_id":"w","steps":[{"id":"c","command":"`+command+`","args":["c"],"body":{}}]}`)
		if code != 9 {
			t.Fatal(command, code)
		}
	}
	code, _ = invoke(t, []string{"manifest", "validate", "--file", "-"}, `{"schema_version":"1","steps":[{"id":"c","command":"workspace authority category create","body":{}}]}`)
	if code != 2 {
		t.Fatal(code)
	}
}
