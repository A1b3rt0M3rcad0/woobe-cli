package controlplane

import (
	"context"
	"encoding/json"
	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/output"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestTransportContract(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer fixture" {
			t.Error("missing bearer")
		}
		if r.URL.Query().Get("project_id") != "p" {
			t.Error("missing project")
		}
		if r.Header.Get("If-Match") != "rev1" {
			t.Error("missing condition")
		}
		var m map[string]any
		_ = json.NewDecoder(r.Body).Decode(&m)
		if v, ok := m["prompt"]; !ok || v != nil {
			t.Error("null lost")
		}
		if _, ok := m["name"]; ok {
			t.Error("omitted field sent")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":"a"}}`))
	}))
	defer s.Close()
	c, _ := New(s.URL, "fixture", time.Second)
	c.Headers.Set("If-Match", "rev1")
	_, _, e := c.Request(context.Background(), "PATCH", "/ai/agents/a", url.Values{"project_id": {"p"}}, []byte(`{"prompt":null}`))
	if e != nil {
		t.Fatal(e)
	}
}
func TestStatusCodes(t *testing.T) {
	for status, code := range map[int]int{401: 3, 403: 4, 404: 5, 409: 6, 412: 6, 422: 2, 501: 9} {
		status, code := status, code
		t.Run(http.StatusText(status), func(t *testing.T) {
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-Request-ID", "req")
				w.WriteHeader(status)
				_, _ = w.Write([]byte(`{"secret":"must-not-leak"}`))
			}))
			defer s.Close()
			c, _ := New(s.URL, "", time.Second)
			_, _, e := c.Request(context.Background(), "POST", "/x", nil, nil)
			v := output.Normalize(e)
			if v.Code != code || v.RequestID != "req" || v.Message == "must-not-leak" {
				t.Fatalf("%+v", v)
			}
		})
	}
}
func TestNoRedirectOrWriteRetry(t *testing.T) {
	count := 0
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count++
		w.Header().Set("Location", "https://example.invalid")
		w.WriteHeader(307)
	}))
	defer s.Close()
	c, _ := New(s.URL, "fixture", time.Second)
	_, _, e := c.Request(context.Background(), "POST", "/x", nil, nil)
	if e == nil || count != 1 {
		t.Fatal("redirect or write replay")
	}
}
func TestRejectForeignPaths(t *testing.T) {
	c, _ := New("https://example.invalid", "", time.Second)
	for _, p := range []string{"https://evil.invalid", "//evil.invalid", "/x?secret=y", "/x#y", "/x\\y"} {
		if _, _, e := c.Request(context.Background(), "GET", p, nil, nil); output.Normalize(e).Code != 2 {
			t.Fatal(p)
		}
	}
}
func TestCancellation(t *testing.T) {
	c, _ := New("http://127.0.0.1:1", "", time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, _, e := c.Request(ctx, "PATCH", "/x", nil, []byte(`{}`))
	v := output.Normalize(e)
	if v.Code != 130 || v.Outcome != "unknown" {
		t.Fatalf("%+v", v)
	}
}
