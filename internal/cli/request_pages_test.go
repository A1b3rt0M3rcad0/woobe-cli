package cli

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestPagesPartialOutput(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Link", `<?page=2>; rel="next"`)
		_, _ = w.Write([]byte(`{"token":"secret"}`))
	}))
	defer s.Close()
	code, v := invoke(t, []string{"request-pages", "/items", "--api-url", s.URL, "--max-pages", "1"}, "")
	if code != 10 || v["success"] != false {
		t.Fatal(v)
	}
	page := v["data"].(map[string]any)["pages"].([]any)[0].(map[string]any)
	if page["token"] != "[REDACTED]" {
		t.Fatal("secret leaked")
	}
}
