package schemacheck

func inspect(raw, doc any, p string, depth int, refs map[string]bool) error {
	if depth > 64 {
		return unsupported(p, "schema depth exceeds 64")
	}
	if e := check(raw, nil, doc, p, depth); e != nil {
		if x, ok := e.(*Error); ok && x.Unsupported {
			return e
		}
	}
	s, ok := raw.(map[string]any)
	if !ok {
		return nil
	}
	if ref, ok := s["$ref"].(string); ok && !refs[ref] {
		refs[ref] = true
		target, e := Resolve(doc, ref)
		if e != nil {
			return e
		}
		if e = inspect(target, doc, p, depth+1, refs); e != nil {
			return e
		}
	}
	for _, k := range []string{"properties", "dependentSchemas", "patternProperties"} {
		if m, ok := s[k].(map[string]any); ok {
			for name, sub := range m {
				if e := inspect(sub, doc, p+"."+name, depth+1, refs); e != nil {
					return e
				}
			}
		}
	}
	for _, k := range []string{"items", "additionalProperties", "propertyNames", "not"} {
		if sub, ok := s[k]; ok {
			if e := inspect(sub, doc, p, depth+1, refs); e != nil {
				return e
			}
		}
	}
	for _, k := range []string{"allOf", "anyOf", "oneOf"} {
		if a, ok := s[k].([]any); ok {
			for _, sub := range a {
				if e := inspect(sub, doc, p, depth+1, refs); e != nil {
					return e
				}
			}
		}
	}
	return nil
}
