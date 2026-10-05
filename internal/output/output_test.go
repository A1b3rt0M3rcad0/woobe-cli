package output

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestNestedTypedSecrets(t *testing.T) {
	v := struct {
		Body json.RawMessage `json:"body"`
	}{json.RawMessage(`{"password":"never-output","nested":{"access_token":"never-output"}}`)}
	b, _ := json.Marshal(Redact(v))
	if strings.Contains(string(b), "never-output") {
		t.Fatal(string(b))
	}
}
