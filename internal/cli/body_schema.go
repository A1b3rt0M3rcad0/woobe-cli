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
	if obj, ok := body.(map[string]any); ok {
		if ref, ok := obj["$ref"].(string); ok {
			var e error
			body, e = schemacheck.Resolve(doc, ref)
			if e != nil {
				return schemaError(e)
			}
		}
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
