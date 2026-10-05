package schemacheck

import (
	"regexp"
	"strings"
)

func validType(v any, p string) error {
	types := []any{v}
	if a, ok := v.([]any); ok {
		types = a
	}
	if len(types) == 0 {
		return unsupported(p, "empty type declaration")
	}
	seen := map[string]bool{}
	for _, v := range types {
		s, ok := v.(string)
		if !ok || !strings.Contains("|null|object|array|string|boolean|number|integer|", "|"+s+"|") || s == "" || seen[s] {
			return unsupported(p, "invalid type declaration")
		}
		seen[s] = true
	}
	return nil
}

func stringList(v any, p, k string) error {
	a, ok := v.([]any)
	if !ok {
		return unsupported(p, "invalid "+k)
	}
	seen := map[string]bool{}
	for _, v := range a {
		s, ok := v.(string)
		if !ok || seen[s] {
			return unsupported(p, "invalid "+k)
		}
		seen[s] = true
	}
	return nil
}
func declarations(s map[string]any, p string) error {
	if v, ok := s["type"]; ok {
		if e := validType(v, p); e != nil {
			return e
		}
	}
	for _, k := range []string{"nullable", "uniqueItems", "readOnly", "writeOnly", "deprecated"} {
		if v, ok := s[k]; ok {
			if _, ok := v.(bool); !ok {
				return unsupported(p, "invalid "+k)
			}
		}
	}
	if v, ok := s["required"]; ok {
		if e := stringList(v, p, "required"); e != nil {
			return e
		}
	}
	if v, ok := s["properties"]; ok {
		if _, ok := v.(map[string]any); !ok {
			return unsupported(p, "invalid properties")
		}
	}
	if v, ok := s["enum"]; ok {
		a, ok := v.([]any)
		if !ok || len(a) == 0 {
			return unsupported(p, "invalid enum")
		}
	}
	if v, ok := s["$ref"]; ok {
		if _, ok := v.(string); !ok {
			return unsupported(p, "invalid $ref")
		}
	}
	if v, ok := s["pattern"]; ok {
		r, ok := v.(string)
		if !ok {
			return unsupported(p, "invalid pattern")
		}
		if _, e := regexp.Compile(r); e != nil {
			return unsupported(p, "pattern incompatible with RE2")
		}
	}
	for _, k := range []string{"minimum", "maximum", "exclusiveMinimum", "exclusiveMaximum", "multipleOf", "minItems", "maxItems", "minLength", "maxLength"} {
		if v, ok := s[k]; ok {
			n := number(v)
			if n == nil {
				return unsupported(p, "invalid "+k)
			}
			if strings.HasPrefix(k, "min") && k != "minimum" || strings.HasPrefix(k, "max") && k != "maximum" {
				if !n.IsInt() || n.Sign() < 0 {
					return unsupported(p, "invalid "+k)
				}
			}
			if k == "multipleOf" && n.Sign() <= 0 {
				return unsupported(p, "invalid multipleOf")
			}
		}
	}
	return nil
}
