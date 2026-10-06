package cli

import (
	"encoding/json"
	"io"
	"net/url"
	"strings"

	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/output"
	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/schemacheck"
)

type advertisedParameter struct {
	name, location string
	required       bool
	schema         any
	style          string
	explode        bool
}

// Operation declarations replace path-level parameters with the same name/location.
func parameterDeclarations(doc, path, def map[string]any) ([]advertisedParameter, error) {
	parameters := []advertisedParameter{}
	indexes := map[string]int{}
	for _, source := range []map[string]any{path, def} {
		raw, exists := source["parameters"]
		if !exists {
			continue
		}
		list, ok := raw.([]any)
		if !ok || len(list) > 1024 {
			return nil, output.New(9, "invalid or excessive OpenAPI parameters")
		}
		seen := map[string]bool{}
		for _, raw := range list {
			resolved, e := resolveOpenAPIObject(doc, raw)
			if e != nil {
				return nil, e
			}
			obj := resolved.(map[string]any)
			name, ok := obj["name"].(string)
			if !ok || name == "" {
				return nil, output.New(9, "invalid parameter name")
			}
			location, ok := obj["in"].(string)
			if !ok || location != "path" && location != "query" && location != "header" && location != "cookie" {
				return nil, output.New(9, "invalid parameter location")
			}
			key := location + ":" + name
			if seen[key] {
				return nil, output.New(9, "duplicate OpenAPI parameter declaration")
			}
			seen[key] = true
			required := false
			if v, exists := obj["required"]; exists {
				var ok bool
				required, ok = v.(bool)
				if !ok {
					return nil, output.New(9, "invalid parameter required declaration")
				}
			}
			if location == "path" && !required {
				return nil, output.New(9, "path parameter must be required")
			}
			p := advertisedParameter{name: name, location: location, required: required, schema: obj["schema"]}
			p.style = "form"
			p.explode = true
			if location == "path" {
				p.style = "simple"
				p.explode = false
			}
			if v, exists := obj["style"]; exists {
				var ok bool
				p.style, ok = v.(string)
				if !ok {
					return nil, output.New(9, "invalid parameter style")
				}
			}
			if v, exists := obj["explode"]; exists {
				var ok bool
				p.explode, ok = v.(bool)
				if !ok {
					return nil, output.New(9, "invalid parameter explode")
				}
			}
			if _, exists := obj["content"]; exists {
				return nil, output.New(9, "content parameter serialization is not supported")
			}
			if p.schema == nil {
				return nil, output.New(9, "parameter schema is not advertised")
			}
			if i, exists := indexes[key]; exists {
				parameters[i] = p
			} else {
				indexes[key] = len(parameters)
				parameters = append(parameters, p)
			}
		}
	}
	return parameters, nil
}

func validateParameterSchema(doc, path, def map[string]any, values map[string]string, query url.Values) error {
	if e := validateSchemaDialect(doc); e != nil {
		return e
	}
	parameters, e := parameterDeclarations(doc, path, def)
	if e != nil {
		return e
	}
	known := map[string]bool{}
	for _, p := range parameters {
		if p.location != "path" && p.location != "query" {
			continue
		}
		var raw []string
		if p.location == "path" {
			if v, ok := values[p.name]; ok {
				raw = []string{v}
			}
		} else {
			known[p.name] = true
			raw = query[p.name]
		}
		if len(raw) == 0 {
			if p.required {
				return output.New(2, "required "+p.location+" parameter is missing: "+p.name)
			}
			continue
		}
		v, e := decodeParameter(p, raw, doc)
		if e != nil {
			return e
		}
		if e := schemacheck.Check(p.schema, v, doc); e != nil {
			return schemaError(e)
		}
	}
	for name := range query {
		if !known[name] {
			return output.New(9, "query parameter is not advertised: "+name)
		}
	}
	for name := range values {
		found := false
		for _, p := range parameters {
			found = found || p.location == "path" && p.name == name
		}
		if !found {
			return output.New(9, "path parameter is not advertised: "+name)
		}
	}
	return nil
}

func decodeParameter(p advertisedParameter, raw []string, doc map[string]any) (any, error) {
	if p.location == "path" && p.style != "simple" || p.location == "query" && p.style != "form" {
		return nil, output.New(9, "unsupported parameter serialization style")
	}
	if len(raw) != 1 {
		return nil, output.New(2, "scalar parameter must occur exactly once: "+p.name)
	}
	schema := p.schema
	for depth := 0; depth < 16; depth++ {
		obj, ok := schema.(map[string]any)
		if !ok {
			if b, ok := schema.(bool); ok && b {
				return raw[0], nil
			}
			return nil, output.New(9, "unsupported parameter schema")
		}
		if ref, ok := obj["$ref"].(string); ok {
			var e error
			schema, e = schemacheck.Resolve(doc, ref)
			if e != nil {
				return nil, schemaError(e)
			}
			continue
		}
		typ, _ := obj["type"].(string)
		switch typ {
		case "", "string":
			return raw[0], nil
		case "boolean":
			if raw[0] == "true" {
				return true, nil
			}
			if raw[0] == "false" {
				return false, nil
			}
			return nil, output.New(2, "boolean parameter must be true or false: "+p.name)
		case "integer", "number":
			d := json.NewDecoder(strings.NewReader(raw[0]))
			d.UseNumber()
			var value any
			if e := d.Decode(&value); e != nil {
				return nil, output.New(2, "invalid numeric parameter: "+p.name)
			}
			if _, ok := value.(json.Number); !ok {
				return nil, output.New(2, "invalid numeric parameter: "+p.name)
			}
			var extra any
			if d.Decode(&extra) != io.EOF {
				return nil, output.New(2, "ambiguous numeric parameter: "+p.name)
			}
			return value, nil
		default:
			return nil, output.New(9, "unsupported parameter schema type")
		}
	}
	return nil, output.New(9, "parameter reference chain exceeds 16")
}
