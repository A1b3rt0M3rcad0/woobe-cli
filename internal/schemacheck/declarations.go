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
