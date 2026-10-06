package cli

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/packagefmt"
)

func TestPackageStatusKeepsAcceptedOperationIncomplete(t *testing.T) {
	t.Setenv("WOOBE_CONTROL_KEY", "test-control-key")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := map[string]any{"package_schema_version": "1.0"}
		if strings.HasSuffix(r.URL.Path, "/capabilities") {
			data["schema_catalog_sha256"] = packagefmt.CatalogDigest()
			data["supported_operations"] = []string{"status"}
		} else {
			data["operation_id"] = "operation"
			data["project_id"] = "project"
			data["state"] = "accepted"
			data["terminal"] = false
			data["revision"] = 1
			data["phases_completed"] = 0
			data["phases_total"] = 3
			data["artifact_digest"] = strings.Repeat("a", 64)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"success": true, "data": data})
	}))
	defer server.Close()
	code, result := invoke(t, []string{"package", "status", "operation", "--api-url", server.URL, "--project", "project"}, "")
	if code != 0 || result["success"] != true || result["meta"].(map[string]any)["complete"] != false || result["data"].(map[string]any)["state"] != "accepted" {
		t.Fatal(code, result)
	}
}

func TestPackageUnsupportedServerCannotReceiveCancellation(t *testing.T) {
	t.Setenv("WOOBE_CONTROL_KEY", "test-control-key")
	var mutations atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			mutations.Add(1)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"success": true, "data": map[string]any{"package_schema_version": "1.0", "schema_catalog_sha256": "incompatible", "supported_operations": []string{"cancel"}}})
	}))
	defer server.Close()
	code, result := invoke(t, []string{"package", "cancel", "operation", "--api-url", server.URL, "--project", "project"}, "")
	if code != 9 || mutations.Load() != 0 {
		t.Fatal(code, result, mutations.Load())
	}
}

func TestPackageResumeRequiresObservedRevisionBeforeCredentials(t *testing.T) {
	code, result := invoke(t, []string{"package", "resume", "operation", "--credential", "missing"}, "")
	if code != 2 {
		t.Fatal(code, result)
	}
}
