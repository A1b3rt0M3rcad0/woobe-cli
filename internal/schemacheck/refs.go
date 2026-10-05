package schemacheck

import "strings"

func Resolve(doc any, ref string) (any, error) {
	if !strings.HasPrefix(ref, "#/") {
		return nil, unsupported("$", "external schema references are not fetched")
	}
	cur := doc
	for _, part := range strings.Split(strings.TrimPrefix(ref, "#/"), "/") {
		part = strings.ReplaceAll(strings.ReplaceAll(part, "~1", "/"), "~0", "~")
		obj, ok := cur.(map[string]any)
		if !ok {
			return nil, unsupported("$", "unresolved local schema reference")
		}
		cur, ok = obj[part]
		if !ok {
			return nil, unsupported("$", "unresolved local schema reference")
		}
	}
	return cur, nil
}
