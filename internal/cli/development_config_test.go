package cli

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestMalformedProjectConfigDoesNotAffectDirectCommands(t *testing.T) {
	directory := t.TempDir()
	t.Chdir(directory)
	os.WriteFile(filepath.Join(directory, ".woobe-config"), []byte("broken: ["), 0600)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte(`{"success":true,"data":[]}`)) }))
	defer server.Close()
	for _, args := range [][]string{{"version"}, {"project", "agent", "list", "--project", "project", "--api-url", server.URL}, {"project", "provider-model", "list", "--project", "project", "--api-url", server.URL}, {"help", "--output", "json"}} {
		code, result := invoke(t, args, "")
		if code != 0 {
			t.Fatal(code, result)
		}
	}
	code, result := invoke(t, []string{"config", "check"}, "")
	if code != 2 {
		t.Fatal(code, result)
	}
}

func TestDevelopmentInitAndOfflineValidation(t *testing.T) {
	t.Chdir(t.TempDir())
	code, result := invoke(t, []string{"init", "--root", "definitions"}, "")
	if code != 0 {
		t.Fatal(code, result)
	}
	code, result = invoke(t, []string{"config", "check"}, "")
	if code != 0 || result["data"].(map[string]any)["valid"] != true {
		t.Fatal(code, result)
	}
	code, result = invoke(t, []string{"init"}, "")
	if code != 2 {
		t.Fatal("overwrote config", code, result)
	}
	data, err := os.ReadFile(filepath.Join("definitions", ".gitignore"))
	if err != nil || string(data) != ".state/\n" {
		t.Fatal(string(data), err)
	}
}
