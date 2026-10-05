package cli

import (
	"encoding/json"
	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/output"
	"testing"
)

func schemaDoc(s string) map[string]any {
	var v map[string]any
	json.Unmarshal([]byte(s), &v)
	return v
}
func TestSharedBodySchema(t *testing.T) {
	def := schemaDoc(`{"requestBody":{"content":{"application/json":{"schema":{"required":["name"]}}}}}`)
	if e := validateBodySchema(nil, def, []byte(`{"name":"A"}`)); e != nil {
		t.Fatal(e)
	}
	if e := validateBodySchema(nil, def, []byte(`{}`)); output.Normalize(e).Code != 2 {
		t.Fatal(e)
	}
}

func TestRequestBodyReferenceChains(t *testing.T) {
	doc := schemaDoc(`{"components":{"requestBodies":{"a":{"$ref":"#/components/requestBodies/b"},"b":{"content":{"application/json":{"schema":true}}}}}}`)
	def := schemaDoc(`{"requestBody":{"$ref":"#/components/requestBodies/a"}}`)
	if e := validateBodySchema(doc, def, []byte(`{}`)); e != nil {
		t.Fatal(e)
	}
	doc["components"] = schemaDoc(`{"requestBodies":{"a":{"$ref":"#/components/requestBodies/a"}}}`)
	if e := validateBodySchema(doc, def, []byte(`{}`)); output.Normalize(e).Code != 9 {
		t.Fatal(e)
	}
}
