package cli

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/devworkspace"
	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/packagefmt"
)

func TestRetainedHydrationVerifiesBytesAndPreservesWorkingFiles(t *testing.T) {
	for _, mode := range []string{"valid", "dry-run", "tamper", "wrong-record", "unsupported"} {
		t.Run(mode, func(t *testing.T) {
			root, _ := filepath.Abs("../..")
			t.Setenv("WOOBE_TEST_REPOSITORY", root)
			var receipt map[string]any
			var archive []byte
			var digest string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Error("hydration attempted write", r.Method)
				}
				if strings.HasSuffix(r.URL.Path, "/capabilities") {
					_ = json.NewEncoder(w).Encode(map[string]any{"success": true, "data": map[string]any{"package_schema_version": "1.0", "schema_catalog_sha256": packagefmt.CatalogDigest(), "supported_operations": []string{"registry"}, "asac": map[string]any{"retained_object_hydration": mode != "unsupported"}}})
					return
				}
				if strings.HasSuffix(r.URL.Path, "/receipt") {
					value := map[string]any{}
					for k, v := range receipt {
						value[k] = v
					}
					if mode == "wrong-record" {
						value["record_digest"] = "sha256:" + strings.Repeat("f", 64)
					}
					_ = json.NewEncoder(w).Encode(map[string]any{"success": true, "data": value})
					return
				}
				w.Header().Set("Content-Length", fmt.Sprint(len(archive)))
				w.Header().Set("X-Woobe-Transport-Sha256", digest)
				w.Header().Set("X-Woobe-Artifact-Digest", strings.TrimPrefix(receipt["artifact_digest"].(string), "sha256:"))
				w.Header().Set("X-Woobe-Record-Digest", receipt["record_digest"].(string))
				content := append([]byte(nil), archive...)
				if mode == "tamper" {
					content[0] ^= 1
				}
				_, _ = w.Write(content)
			}))
			defer server.Close()
			c, state, resource := asacWorkspace(t, server.URL)
			graph, err := devworkspace.LoadGraph(c)
			if err != nil {
				t.Fatal(err)
			}
			record, err := graph.CreateRevision(resource, state, "Exact hydration", nil)
			if err != nil {
				t.Fatal(err)
			}
			bundle, _, err := graph.Compile(resource.Key, state.Requirements, state.Credentials)
			if err != nil {
				t.Fatal(err)
			}
			defer bundle.Close()
			var encoded bytes.Buffer
			if err := bundle.Archive(&encoded, true); err != nil {
				t.Fatal(err)
			}
			archive = encoded.Bytes()
			hash := sha256.Sum256(archive)
			digest = hex.EncodeToString(hash[:])
			receipt = map[string]any{"package_schema_version": "1.0", "schema_version": "1.0", "kind": resource.Kind, "project_id": "project", "resource_id": "01a0a033-5820-770c-854b-902864857273", "resource_uid": resource.UID, "revision_id": record.ID, "record_digest": record.RecordDigest, "artifact_digest": record.ArtifactDigest, "definition_digest": record.DefinitionDigest, "transport_digest": digest, "size_bytes": len(archive), "inventory": bundle.Inventory, "author_artifact_status": "declared_not_retained", "complete": true}
			object := filepath.Join(c.RootPath(), "objects", "sha256", bundle.ArtifactDigest+".tar.gz")
			if err := os.Remove(object); err != nil {
				t.Fatal(err)
			}
			if err := os.RemoveAll(filepath.Join(c.RootPath(), ".state")); err != nil {
				t.Fatal(err)
			}
			author, _ := c.ResourcePath(resource)
			before, _ := os.ReadFile(author)
			args := []string{"agent", "@support", "revision", "hydrate", record.ID, "--api-url", server.URL, "--workspace", "workspace", "--project", "project", "--output", "json"}
			if mode == "dry-run" {
				args = append(args, "--dry-run")
			}
			code, response := invoke(t, args, "")
			if mode == "valid" || mode == "dry-run" {
				if code != 0 || response["data"].(map[string]any)["executable_object"] != "verified" || response["data"].(map[string]any)["bindings_restored"] != false {
					t.Fatal(code, response)
				}
			} else if code != 9 {
				t.Fatal("accepted invalid hydration", code, response)
			}
			after, _ := os.ReadFile(author)
			if string(after) != string(before) {
				t.Fatal("hydration changed source")
			}
			tracking, err := c.ReadTracking(resource)
			if err != nil || tracking.Working != record.ID {
				t.Fatal("hydration moved head", tracking, err)
			}
			_, err = os.Stat(object)
			if mode == "valid" {
				if err != nil {
					t.Fatal("verified object missing", err)
				}
				code, response = invoke(t, args, "")
				if code != 0 {
					t.Fatal("non-idempotent hydration", code, response)
				}
			} else if !os.IsNotExist(err) {
				t.Fatal("unverified/dry object published", err)
			}
		})
	}
}
