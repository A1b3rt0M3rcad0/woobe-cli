package schemacheck

import "encoding/json"

func inspect(raw, doc any, p string, depth int, refs map[string]bool, budget *int) error {
	if depth > 64 {
		return unsupported(p, "schema depth exceeds 64")
	}
	*budget--
	if *budget < 0 {
		return unsupported(p, "schema evaluation budget exceeded")
	}
	if _, ok := raw.(bool); ok {
		return nil
	}
	s, ok := raw.(map[string]any)
	if !ok {
		return unsupported(p, "schema must be object or boolean")
	}
	if e := schemaHeader(s, p); e != nil {
		return e
	}
	for _, k := range []string{"enum", "const"} {
		if v, ok := s[k]; ok {
			if e := valueWork(v, 0, budget); e != nil {
				return e
			}
		}
	}
	if ref, ok := s["$ref"].(string); ok && !refs[ref] {
		refs[ref] = true
		target, e := Resolve(doc, ref)
		if e != nil {
			return e
		}
		if e = inspect(target, doc, p, depth+1, refs, budget); e != nil {
			return e
		}
	}
	for _, k := range []string{"properties", "dependentSchemas", "patternProperties"} {
		if m, ok := s[k].(map[string]any); ok {
			for name, sub := range m {
				if e := inspect(sub, doc, p+"."+name, depth+1, refs, budget); e != nil {
					return e
				}
			}
		}
	}
	for _, k := range []string{"items", "additionalProperties", "propertyNames", "contains", "not", "if", "then", "else"} {
		if sub, ok := s[k]; ok {
			if e := inspect(sub, doc, p, depth+1, refs, budget); e != nil {
				return e
			}
		}
	}
	for _, k := range []string{"allOf", "anyOf", "oneOf", "prefixItems"} {
		if a, ok := s[k].([]any); ok {
			for _, sub := range a {
				if e := inspect(sub, doc, p, depth+1, refs, budget); e != nil {
					return e
				}
			}
		}
	}
	return nil
}

func valueWork(v any, depth int, budget *int) error {
	*budget--
	if *budget < 0 || depth > 128 {
		return unsupported("$", "input evaluation budget exceeded")
	}
	switch x := v.(type) {
	case json.Number:
		if number(x) == nil {
			return unsupported("$", "numeric representation exceeds supported precision or exponent")
		}
	case []any:
		for _, v := range x {
			if e := valueWork(v, depth+1, budget); e != nil {
				return e
			}
		}
	case map[string]any:
		for _, v := range x {
			if e := valueWork(v, depth+1, budget); e != nil {
				return e
			}
		}
	}
	return nil
}
