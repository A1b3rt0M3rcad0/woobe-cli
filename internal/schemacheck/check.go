package schemacheck

import (
	"encoding/json"
	"fmt"
	"math/big"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"
)

type Error struct {
	Path, Rule  string
	Unsupported bool
}

func (e *Error) Error() string      { return fmt.Sprintf("schema check at %s: %s", e.Path, e.Rule) }
func fail(p, r string) error        { return &Error{Path: p, Rule: r} }
func unsupported(p, r string) error { return &Error{Path: p, Rule: r, Unsupported: true} }

// Check implements a bounded, explicit schema subset. Unknown assertions fail closed.
func Check(schema, value, document any) error {
	if e := inspect(schema, document, "$", 0, map[string]bool{}); e != nil {
		return e
	}
	return check(schema, value, document, "$", 0)
}
func equal(a, b any) bool {
	if x, y := number(a), number(b); x != nil || y != nil {
		return x != nil && y != nil && x.Cmp(y) == 0
	}
	switch x := a.(type) {
	case map[string]any:
		y, ok := b.(map[string]any)
		if !ok || len(x) != len(y) {
			return false
		}
		for k, v := range x {
			w, ok := y[k]
			if !ok || !equal(v, w) {
				return false
			}
		}
		return true
	case []any:
		y, ok := b.([]any)
		if !ok || len(x) != len(y) {
			return false
		}
		for i, v := range x {
			if !equal(v, y[i]) {
				return false
			}
		}
		return true
	}
	x, _ := json.Marshal(a)
	y, _ := json.Marshal(b)
	return string(x) == string(y)
}
func number(v any) *big.Rat {
	switch x := v.(type) {
	case json.Number:
		r, _ := new(big.Rat).SetString(x.String())
		return r
	case float64:
		r, _ := new(big.Rat).SetString(fmt.Sprint(x))
		return r
	}
	return nil
}
func typeOK(t string, v any) bool {
	switch t {
	case "null":
		return v == nil
	case "object":
		_, ok := v.(map[string]any)
		return ok
	case "array":
		_, ok := v.([]any)
		return ok
	case "string":
		_, ok := v.(string)
		return ok
	case "boolean":
		_, ok := v.(bool)
		return ok
	case "number":
		return number(v) != nil
	case "integer":
		n := number(v)
		return n != nil && n.IsInt()
	}
	return false
}
func check(raw, v, doc any, p string, depth int) error {
	if depth > 64 {
		return unsupported(p, "schema depth exceeds 64")
	}
	if b, ok := raw.(bool); ok {
		if !b {
			return fail(p, "false schema")
		}
		return nil
	}
	s, ok := raw.(map[string]any)
	if !ok {
		return unsupported(p, "schema must be object or boolean")
	}
	if e := declarations(s, p); e != nil {
		return e
	}
	allowed := map[string]bool{}
	for _, k := range strings.Fields("$schema $id title description default examples example deprecated readOnly writeOnly discriminator xml externalDocs type nullable required properties additionalProperties items minItems maxItems minProperties maxProperties dependentRequired dependentSchemas uniqueItems minLength maxLength pattern minimum maximum exclusiveMinimum exclusiveMaximum multipleOf enum const $ref allOf anyOf oneOf not") {
		allowed[k] = true
	}
	for k := range s {
		if !allowed[k] && !strings.HasPrefix(k, "x-") {
			return unsupported(p, "unsupported keyword: "+k)
		}
	}
	if ref, ok := s["$ref"].(string); ok {
		target, e := Resolve(doc, ref)
		if e != nil {
			return e
		}
		if e = check(target, v, doc, p, depth+1); e != nil {
			return e
		}
	}
	for _, key := range []string{"allOf", "anyOf", "oneOf"} {
		if raw, ok := s[key]; ok {
			branches, ok := raw.([]any)
			if !ok || len(branches) == 0 {
				return unsupported(p, "invalid composition")
			}
			matches := 0
			for _, b := range branches {
				e := check(b, v, doc, p, depth+1)
				if x, ok := e.(*Error); ok && x.Unsupported {
					return e
				}
				if e == nil {
					matches++
				}
			}
			if key == "allOf" && matches != len(branches) || key == "anyOf" && matches == 0 || key == "oneOf" && matches != 1 {
				return fail(p, key+" mismatch")
			}
		}
	}
	if b, ok := s["not"]; ok {
		e := check(b, v, doc, p, depth+1)
		if x, ok := e.(*Error); ok && x.Unsupported {
			return e
		}
		if e == nil {
			return fail(p, "not mismatch")
		}
	}
	if v == nil && s["nullable"] == true {
		return nil
	}
	if typ, ok := s["type"]; ok {
		if e := validType(typ, p); e != nil {
			return e
		}
		matched := false
		switch t := typ.(type) {
		case string:
			matched = typeOK(t, v)
		case []any:
			for _, t := range t {
				if name, ok := t.(string); ok && typeOK(name, v) {
					matched = true
				}
			}
		default:
			return unsupported(p, "invalid type declaration")
		}
		if !matched {
			return fail(p, "type mismatch")
		}
	}
	if c, ok := s["const"]; ok && !equal(c, v) {
		return fail(p, "const mismatch")
	}
	if e, ok := s["enum"].([]any); ok {
		match := false
		for _, x := range e {
			match = match || equal(x, v)
		}
		if !match {
			return fail(p, "enum mismatch")
		}
	}
	if obj, ok := v.(map[string]any); ok {
		if e := bounds(s, "minProperties", "maxProperties", len(obj), p); e != nil {
			return e
		}
		if req, ok := s["required"].([]any); ok {
			for _, r := range req {
				name, ok := r.(string)
				if !ok {
					return unsupported(p, "invalid required field")
				}
				if _, ok := obj[name]; !ok {
					return fail(p, "required field absent: "+name)
				}
			}
		}
		if deps, ok := s["dependentRequired"].(map[string]any); ok {
			for k, raw := range deps {
				if _, exists := obj[k]; exists {
					for _, name := range raw.([]any) {
						if _, exists := obj[name.(string)]; !exists {
							return fail(p, "dependent required field absent")
						}
					}
				}
			}
		}
		if deps, ok := s["dependentSchemas"].(map[string]any); ok {
			for k, sub := range deps {
				if _, exists := obj[k]; exists {
					if e := check(sub, v, doc, p, depth+1); e != nil {
						return e
					}
				}
			}
		}
		props, _ := s["properties"].(map[string]any)
		keys := []string{}
		for k := range obj {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			if sub, ok := props[k]; ok {
				if e := check(sub, obj[k], doc, p+"."+k, depth+1); e != nil {
					return e
				}
			} else if sub, ok := s["additionalProperties"]; ok {
				if e := check(sub, obj[k], doc, p+"."+k, depth+1); e != nil {
					return e
				}
			}
		}
	}
	if a, ok := v.([]any); ok {
		if e := bounds(s, "minItems", "maxItems", len(a), p); e != nil {
			return e
		}
		seen := []any{}
		for i, x := range a {
			if s["uniqueItems"] == true {
				for _, prior := range seen {
					if equal(prior, x) {
						return fail(p, "duplicate array item")
					}
				}
				seen = append(seen, x)
			}
			if sub, ok := s["items"]; ok {
				if e := check(sub, x, doc, fmt.Sprintf("%s[%d]", p, i), depth+1); e != nil {
					return e
				}
			}
		}
	}
	if text, ok := v.(string); ok {
		if e := bounds(s, "minLength", "maxLength", utf8.RuneCountInString(text), p); e != nil {
			return e
		}
		if pattern, ok := s["pattern"].(string); ok {
			r, e := regexp.Compile(pattern)
			if e != nil {
				return unsupported(p, "pattern incompatible with RE2")
			}
			if !r.MatchString(text) {
				return fail(p, "pattern mismatch")
			}
		}
	}
	if n := number(v); n != nil {
		for _, k := range []string{"minimum", "maximum", "exclusiveMinimum", "exclusiveMaximum", "multipleOf"} {
			limit, ok := s[k]
			if !ok {
				continue
			}
			l := number(limit)
			if l == nil {
				return unsupported(p, "numeric assertion unsupported: "+k)
			}
			c := n.Cmp(l)
			bad := k == "minimum" && c < 0 || k == "maximum" && c > 0 || k == "exclusiveMinimum" && c <= 0 || k == "exclusiveMaximum" && c >= 0
			if k == "multipleOf" {
				if l.Sign() <= 0 {
					return unsupported(p, "invalid multipleOf")
				}
				bad = !new(big.Rat).Quo(n, l).IsInt()
			}
			if bad {
				return fail(p, k+" constraint")
			}
		}
	}
	return nil
}
func bounds(s map[string]any, min, max string, n int, p string) error {
	for _, k := range []string{min, max} {
		if v, ok := s[k]; ok {
			limit := number(v)
			if limit == nil || !limit.IsInt() {
				return unsupported(p, "invalid "+k)
			}
			c := new(big.Rat).SetInt64(int64(n)).Cmp(limit)
			if k == min && c < 0 || k == max && c > 0 {
				return fail(p, k+" constraint")
			}
		}
	}
	return nil
}
