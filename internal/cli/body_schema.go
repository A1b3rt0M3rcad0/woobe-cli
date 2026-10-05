package cli

import (
	"bytes"
	"encoding/json"
	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/output"
	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/schemacheck"
)

func validateBodySchema(doc, def map[string]any, b []byte) error {
	body, ok := def["requestBody"]
	if !ok {
		return output.New(9, "operation has no advertised request body")
	}
	var e error
	body, e = resolveOpenAPIObject(doc, body)
	if e != nil {
		return e
	}
	obj, _ := body.(map[string]any)
	content, _ := obj["content"].(map[string]any)
	media, _ := content["application/json"].(map[string]any)
	schema, ok := media["schema"]
	if !ok {
		return output.New(9, "operation has no application/json schema")
	}
	var value any
	d := json.NewDecoder(bytes.NewReader(b))
	d.UseNumber()
	if e := d.Decode(&value); e != nil {
		return output.New(2, "invalid JSON body")
	}
	if e := schemacheck.Check(schema, value, doc); e != nil {
		return schemaError(e)
	}
	return nil
}

func resolveOpenAPIObject(doc map[string]any, raw any) (any, error) {
	seen := map[string]bool{}
	for depth := 0; depth < 16; depth++ {
		obj, ok := raw.(map[string]any)
		if !ok {
			return nil, output.New(9, "invalid OpenAPI object")
		}
		v, exists := obj["$ref"]
		if !exists {
			return obj, nil
		}
		ref, ok := v.(string)
		if !ok || seen[ref] {
			return nil, output.New(9, "invalid or cyclic OpenAPI reference")
		}
		seen[ref] = true
		for k := range obj {
			if k != "$ref" && k != "summary" && k != "description" {
				return nil, output.New(9, "unsupported OpenAPI reference siblings")
			}
		}
		var e error
		raw, e = schemacheck.Resolve(doc, ref)
		if e != nil {
			return nil, schemaError(e)
		}
	}
	return nil, output.New(9, "OpenAPI reference chain exceeds 16")
}
