package cli

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBodyValidationPreventsSelectedWrite(t *testing.T) {
	reads, writes := 0, 0
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/openapi.json" {
			reads++
			w.Write([]byte(`{"paths":{"/ai/agents":{"post":{"requestBody":{"content":{"application/json":{"schema":{"required":["name"]}}}}}}}}`))
			return
		}
		writes++
		w.Write([]byte(`{"id":"a"}`))
	}))
	defer s.Close()
	for _, tc := range []struct {
		body string
		code int
	}{{`{}`, 2}, {`{"name":"A"}`, 0}} {
		code, v := invoke(t, []string{"project", "agent", "create", "--file", "-", "--api-url", s.URL, "--validate-body"}, tc.body)
		if code != tc.code {
			t.Fatal(code, v)
		}
		if code != 0 && v["error"].(map[string]any)["write_outcome"] != "not_attempted" {
			t.Fatal(v)
		}
	}
	if reads != 2 || writes != 1 {
		t.Fatal(reads, writes)
	}
	_ = strings.TrimSpace
}

func TestValidationFlagIsNeverSilentlyIgnored(t *testing.T) {
	for _, args := range [][]string{{"request", "POST", "/anything"}, {"runtime", "target", "run", "alias"}, {"project", "agent", "get", "a"}, {"project", "knowledge", "document", "upload"}} {
		code, v := invoke(t, append(args, "--validate-body"), "")
		if code != 9 {
			t.Fatal(args, code, v)
		}
	}
}
