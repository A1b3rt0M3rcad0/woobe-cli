package manifest

import (
	"encoding/json"
	"fmt"
	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/output"
)

func rejectSecrets(v any) error {
	switch x := v.(type) {
	case map[string]any:
		for k, v := range x {
			if output.Sensitive(k) && v != nil {
				return fmt.Errorf("manifest secret values are forbidden; use explicit credential operations")
			}
			if e := rejectSecrets(v); e != nil {
				return e
			}
		}
	case []any:
		for _, v := range x {
			if e := rejectSecrets(v); e != nil {
				return e
			}
		}
	}
	return nil
}
func validateBody(b json.RawMessage) error {
	if len(b) == 0 {
		return nil
	}
	var v any
	if e := json.Unmarshal(b, &v); e != nil {
		return e
	}
	return rejectSecrets(v)
}
