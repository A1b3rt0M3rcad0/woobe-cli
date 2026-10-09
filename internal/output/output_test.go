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

func TestFencingCounterExceptionNeverExposesTextualTokensOrOtherCredentials(t *testing.T) {
	for _, value := range []any{json.Number("1"), json.Number("9223372036854775807"), int(2), int64(3), uint64(4), float64(5)} {
		data := Redact(map[string]any{"fencing_token": value, "access_token": "private-access", "nested": map[string]any{"provider_token": "private-provider"}}).(map[string]any)
		if data["fencing_token"] != value {
			t.Fatal("Concurrency counter hidden", data)
		}
		raw, _ := json.Marshal(data)
		if strings.Contains(string(raw), "private-") {
			t.Fatal("Credential disclosed", string(raw))
		}
	}
	for _, value := range []any{"private-provider-secret", "123", json.Number("1.2"), json.Number("9223372036854775808"), json.Number("1e3"), float64(1.5), float64(0), int(-1), true, nil, map[string]any{"value": "private-token"}} {
		data := Redact(map[string]any{"fencing_token": value, "Fencing_Token": "private-token"}).(map[string]any)
		if data["fencing_token"] != "[REDACTED]" || data["Fencing_Token"] != "[REDACTED]" {
			t.Fatal("Textual/invalid token exposed", data)
		}
	}
}
