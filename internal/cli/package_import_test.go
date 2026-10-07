package cli

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/packageapi"
	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/packagecheckpoint"
	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/packagefmt"
)

func TestPackageImportLostResponseReconcilesWithoutResending(t *testing.T) {
	t.Setenv("WOOBE_CONTROL_KEY", "test-control")
	var posts atomic.Int32
	var acceptedKey string
	digest := strings.Repeat("a", 64)
	operation := map[string]any{"package_schema_version": "1.0", "operation_id": "operation", "project_id": "project", "artifact_digest": digest, "state": "accepted", "terminal": false, "revision": 1, "phases_completed": 0, "phases_total": 3}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := map[string]any{}
		switch {
		case strings.HasSuffix(r.URL.Path, "/capabilities"):
			data = map[string]any{"package_schema_version": "1.0", "schema_catalog_sha256": packagefmt.CatalogDigest(), "supported_operations": []string{"apply"}, "principal_fingerprint": strings.Repeat("b", 64)}
		case strings.HasSuffix(r.URL.Path, "/apply"):
			posts.Add(1)
			acceptedKey = r.Header.Get("Idempotency-Key")
			connection, _, err := w.(http.Hijacker).Hijack()
			if err != nil {
				t.Error(err)
				return
			}
			_ = connection.Close()
			return
		case strings.HasSuffix(r.URL.Path, "/operations/lookup"):
			var request struct {
				Key string `json:"idempotency_key"`
			}
			_ = json.NewDecoder(r.Body).Decode(&request)
			if request.Key != acceptedKey {
				t.Error("lookup changed the accepted key")
			}
			data = operation
		default:
			t.Error(r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"success": true, "data": data})
	}))
	defer server.Close()
	directory := t.TempDir()
	planPath := filepath.Join(directory, "plan.json")
	checkpoint := filepath.Join(directory, "checkpoint.json")
	receipt := packageapi.PlanReceipt{Format: "woobe-package-plan-receipt", SchemaVersion: "1.0", PrincipalFingerprint: strings.Repeat("b", 64), Plan: packageapi.Plan{PackageSchemaVersion: "1.0", APIOrigin: server.URL, ProjectID: "project", UploadID: "upload", PlanID: "plan", PlanDigest: digest, ArtifactDigest: digest, DefinitionDigest: digest, CapabilitiesDigest: digest, BindingsDigest: digest, Lifecycle: "draft", CreatedAt: time.Now().UTC().Format(time.RFC3339Nano), ExpiresAt: time.Now().Add(time.Hour).UTC().Format(time.RFC3339Nano), Effects: []packageapi.Effect{{Key: "agent:x", Owner: "agent", Action: "prepare_inactive_agent", Permission: "agent:write"}}, Bindings: map[string]any{"format": "woobe-package", "schema_version": "1.0", "kind": "ImportBindings", "metadata": map[string]any{"name": "Destination"}, "spec": map[string]any{}}}}
	receipt.Plan.BindingsDigest, _ = packageapi.BindingsDigest(receipt.Plan.Bindings)
	if err := receipt.Save(planPath); err != nil {
		t.Fatal(err)
	}
	args := []string{"package", "import", "--plan-file", planPath, "--checkpoint", checkpoint, "--api-url", server.URL, "--project", "project"}
	code, result := invoke(t, args, "")
	if code != 9 || posts.Load() != 1 {
		t.Fatal(code, result, posts.Load())
	}
	// Reconciliation succeeds even if the protected store would now be unavailable.
	code, result = invoke(t, args, "")
	if code != 0 || posts.Load() != 1 || result["meta"].(map[string]any)["complete"] != false {
		t.Fatal(code, result, posts.Load())
	}
	store, err := packagecheckpoint.Open(checkpoint)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	cp, err := store.Read()
	if err != nil || cp.OperationID != "operation" || cp.IdempotencyKey != acceptedKey || cp.State != packagecheckpoint.Observing {
		t.Fatal(cp, err)
	}
}

func TestPackageImportOfflineAndFlagErrorsNeverContactServer(t *testing.T) {
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { requests.Add(1) }))
	defer server.Close()
	code, result := invoke(t, []string{"package", "import", localPackageFixture(t), "--bind", packageBindingShortcut, "--dry-run", "--api-url", server.URL, "--credential", "missing"}, "")
	if code != 0 || requests.Load() != 0 || result["data"].(map[string]any)["executed"] != false {
		t.Fatal(code, result)
	}
	code, result = invoke(t, []string{"package", "import", localPackageFixture(t), "--plan-file", "missing", "--checkpoint", "missing", "--api-url", server.URL}, "")
	if code != 2 || requests.Load() != 0 {
		t.Fatal(code, result)
	}
}
