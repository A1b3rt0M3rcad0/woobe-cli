package cli

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func TestManifestParameterValidationStopsBeforeWrite(t *testing.T) {
	writes := 0
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/openapi.json" {
			w.Write([]byte(`{"openapi":"3.1.0","paths":{"/ai/agents/{agent_id}":{"patch":{"parameters":[{"name":"agent_id","in":"path","required":true,"schema":{"type":"string","format":"uuid"}}]}}}}`))
			return
		}
		writes++
		w.Write([]byte(`{"id":"x"}`))
	}))
	defer s.Close()
	cp := filepath.Join(t.TempDir(), "cp.json")
	code, v := invoke(t, []string{"manifest", "apply", "--yes", "--checkpoint", cp, "--file", "-", "--api-url", s.URL, "--validate-parameters"}, `{"schema_version":"1","steps":[{"id":"a","command":"project agent update","args":["invalid"],"body":{"name":"A"}}]}`)
	if code != 10 || writes != 0 {
		t.Fatal(code, v, writes)
	}
	code, v = invoke(t, []string{"manifest", "status", "--checkpoint", cp}, "")
	if code != 0 {
		t.Fatal(code, v)
	}
	data := v["data"].(map[string]any)
	if data["checkpoint"].(map[string]any)["steps"].(map[string]any)["a"] != "not_attempted" {
		t.Fatal(v)
	}
}
