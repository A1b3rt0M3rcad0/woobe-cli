package controlplane

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/output"
	"net/url"
	"reflect"
	"regexp"
	"strings"
)

type Pages struct {
	Pages              []any  `json:"pages"`
	Count              int    `json:"page_count"`
	TraversalComplete  bool   `json:"traversal_complete"`
	CollectionComplete string `json:"collection_complete"`
}

var linkPattern = regexp.MustCompile(`<([^<>]+)>\s*((?:;[^<>]*)?)`)

func nextLink(header string) (string, error) {
	next := ""
	if strings.TrimSpace(header) == "" {
		return next, nil
	}
	matches := linkPattern.FindAllStringSubmatch(header, -1)
	if len(matches) == 0 {
		return "", output.New(9, "unsupported pagination Link header")
	}
	for _, m := range matches {
		params := strings.Split(strings.TrimSuffix(strings.TrimSpace(m[2]), ","), ";")
		for _, p := range params {
			k, v, ok := strings.Cut(strings.TrimSpace(p), "=")
			if !ok || strings.ToLower(k) != "rel" {
				continue
			}
			v = strings.Trim(strings.TrimSpace(v), `"`)
			for _, rel := range strings.Fields(v) {
				if strings.EqualFold(rel, "next") {
					if next != "" {
						return "", output.New(9, "multiple pagination next links")
					}
					next = m[1]
				}
			}
		}
	}
	return next, nil
}

// Follow only same-route, same-origin advertised next links. Cursor formats are opaque.
func (c *Client) RequestPages(ctx context.Context, path string, q url.Values, maxPages int) (Pages, error) {
	out := Pages{Pages: []any{}, CollectionComplete: "not_verified"}
	if maxPages < 1 || maxPages > 100 {
		return out, output.New(2, "max-pages must be between 1 and 100")
	}
	if e := ValidatePath(path); e != nil {
		return out, e
	}
	base, e := url.Parse(c.Base + path)
	if e != nil {
		return out, output.New(2, "invalid page URL")
	}
	base.RawQuery = q.Encode()
	current := base
	seen := map[string]bool{}
	bytes := 0
	for {
		if seen[current.String()] {
			return out, output.New(9, "pagination cycle detected")
		}
		seen[current.String()] = true
		v, h, e := c.Request(ctx, "GET", path, current.Query(), nil)
		if e != nil {
			return out, e
		}
		encoded, _ := json.Marshal(v)
		bytes += len(encoded)
		if bytes > 64<<20 {
			return out, output.New(9, "aggregate pages exceed 64 MiB")
		}
		out.Pages = append(out.Pages, v)
		out.Count++
		next, e := nextLink(strings.Join(h.Values("Link"), ","))
		if e != nil {
			return out, e
		}
		if next == "" {
			out.TraversalComplete = true
			return out, nil
		}
		if out.Count >= maxPages {
			return out, output.New(9, "pagination reached max-pages before terminal page")
		}
		ref, e := url.Parse(next)
		if e != nil || ref.User != nil || ref.Fragment != "" {
			return out, output.New(9, "invalid pagination target")
		}
		target := current.ResolveReference(ref)
		if target.Scheme != base.Scheme || target.Host != base.Host || target.EscapedPath() != base.EscapedPath() {
			return out, output.New(9, "pagination must remain on the same origin and route")
		}
		for key, values := range q {
			if strings.HasSuffix(key, "_id") && !reflect.DeepEqual(values, target.Query()[key]) {
				return out, output.New(9, fmt.Sprintf("pagination changed scoped query %s", key))
			}
		}
		current = target
	}
}
