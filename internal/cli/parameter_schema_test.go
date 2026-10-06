package cli

import (
	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/output"
	"net/url"
	"testing"
)

func parameter(name, location, typ string, required bool) map[string]any {
	return map[string]any{"name": name, "in": location, "required": required, "schema": map[string]any{"type": typ}}
}
func TestParameterDeclarationsAndPrimitiveValidation(t *testing.T) {
	doc := map[string]any{"openapi": "3.1.0"}
	path := map[string]any{"parameters": []any{parameter("id", "path", "string", true)}}
	def := map[string]any{"parameters": []any{parameter("limit", "query", "integer", false), parameter("enabled", "query", "boolean", false)}}
	for _, tc := range []struct {
		query url.Values
		code  int
	}{
		{url.Values{"limit": {"2"}, "enabled": {"true"}}, 0},
		{url.Values{"limit": {"2.5"}}, 2},
		{url.Values{"limit": {"2", "3"}}, 2},
		{url.Values{"enabled": {"yes"}}, 2},
		{url.Values{"unknown": {"x"}}, 9},
	} {
		e := validateParameterSchema(doc, path, def, map[string]string{"id": "a"}, tc.query)
		if e == nil && tc.code != 0 || e != nil && output.Normalize(e).Code != tc.code {
			t.Fatal(tc, e)
		}
	}
}
func TestOperationParameterOverride(t *testing.T) {
	path := map[string]any{"parameters": []any{parameter("limit", "query", "string", false)}}
	def := map[string]any{"parameters": []any{parameter("limit", "query", "integer", true)}}
	if e := validateParameterSchema(map[string]any{}, path, def, nil, url.Values{"limit": {"3"}}); e != nil {
		t.Fatal(e)
	}
	if e := validateParameterSchema(map[string]any{}, path, def, nil, nil); e == nil || output.Normalize(e).Code != 2 {
		t.Fatal(e)
	}
	def["parameters"] = []any{parameter("limit", "query", "integer", true), parameter("limit", "query", "string", false)}
	if _, e := parameterDeclarations(map[string]any{}, path, def); e == nil {
		t.Fatal("duplicate accepted")
	}
}
