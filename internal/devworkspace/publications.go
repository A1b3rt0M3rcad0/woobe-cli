package devworkspace

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/asac"
	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/packagefmt"
)

// An authenticated observation and a content digest do not establish offline
// server authenticity. Signing/trust verification is a separate protocol.
func ValidatePublicationReceipt(record map[string]any) error {
	value := clone(record)
	expected := packagefmt.Text(value["receipt_digest"])
	delete(value, "receipt_digest")
	digest, err := asac.Digest("publication-receipt", value)
	if err != nil || digest != expected || record["schema_version"] != "1.0" || record["state"] != "published" || record["published"] != true || record["production_changed"] != false || record["write_outcome"] != "committed" {
		return fmt.Errorf("ASAC_PUBLICATION_RECEIPT_INVALID: integrity or outcome mismatch")
	}
	for _, field := range []string{"publication_id", "operation_id", "resource_uid", "resource_id", "project_id", "candidate_id", "evaluation_id", "release_id"} {
		if !uuidPattern.MatchString(packagefmt.Text(record[field])) {
			return fmt.Errorf("ASAC_PUBLICATION_RECEIPT_INVALID: %s", field)
		}
	}
	if !revisionID.MatchString(packagefmt.Text(record["revision_id"])) {
		return fmt.Errorf("ASAC_PUBLICATION_RECEIPT_INVALID: revision identity")
	}
	for _, field := range []string{"record_digest", "definition_digest", "artifact_digest", "runtime_digest", "binding_digest", "suite_digest", "dataset_digest", "policy_digest"} {
		if !historyDigest.MatchString(packagefmt.Text(record[field])) {
			return fmt.Errorf("ASAC_PUBLICATION_RECEIPT_INVALID: %s", field)
		}
	}
	if record["kind"] != "Agent" && record["kind"] != "Network" {
		return fmt.Errorf("ASAC_PUBLICATION_RECEIPT_INVALID: resource kind")
	}
	actor := packagefmt.Text(record["actor"])
	parts := strings.SplitN(actor, ":", 2)
	if len(parts) != 2 || (parts[0] != "user" && parts[0] != "control") || !uuidPattern.MatchString(parts[1]) {
		return fmt.Errorf("ASAC_PUBLICATION_RECEIPT_INVALID: actor identity")
	}
	if _, err := time.Parse(time.RFC3339Nano, packagefmt.Text(record["created_at"])); err != nil {
		return fmt.Errorf("ASAC_PUBLICATION_RECEIPT_INVALID: timestamp")
	}
	if _, ok := record["reused"].(bool); !ok || record["complete"] != true {
		return fmt.Errorf("ASAC_PUBLICATION_RECEIPT_INVALID: terminal outcome")
	}
	git := packagefmt.Object(record["git"])
	if git["status"] != "unavailable" && git["status"] != "declared" {
		return fmt.Errorf("ASAC_PUBLICATION_RECEIPT_INVALID: unsupported Git attestation")
	}
	allowed := map[string]bool{}
	for _, key := range []string{"schema_version", "kind", "publication_id", "project_id", "resource_id", "resource_uid", "operation_id", "write_outcome", "state", "candidate_id", "evaluation_id", "release_id", "release_version", "reused", "revision_id", "actor", "reason", "created_at", "record_digest", "runtime_digest", "definition_digest", "artifact_digest", "binding_digest", "suite_digest", "dataset_digest", "policy_digest", "git", "published", "production_changed", "complete", "receipt_digest"} {
		allowed[key] = true
	}
	for key := range record {
		if !allowed[key] {
			return fmt.Errorf("ASAC_PUBLICATION_RECEIPT_INVALID: unexpected field")
		}
	}
	return nil
}

func (c *Config) StorePublicationReceipt(resource Resource, target, workspace, project string, record map[string]any) error {
	if err := ValidatePublicationReceipt(record); err != nil {
		return err
	}
	if record["resource_uid"] != resource.UID || record["project_id"] != project || record["kind"] != resource.Kind {
		return fmt.Errorf("ASAC_PUBLICATION_RECEIPT_INVALID: destination identity mismatch")
	}
	namespace, err := asac.Digest("destination", []string{target, workspace, project})
	if err != nil {
		return err
	}
	dir, err := c.historyDirectory(resource)
	if err != nil {
		return err
	}
	wrapper := map[string]any{"format": "woobe-observed-publication", "schema_version": "1.0", "resource_uid": resource.UID, "target": target, "workspace_id": workspace, "project_id": project, "receipt": record, "authenticity": "authenticated_connection_observation"}
	data, err := Encode(wrapper)
	if err != nil {
		return err
	}
	return c.immutableWrite(filepath.Join(dir, "publications", strings.TrimPrefix(namespace, "sha256:")+"_"+packagefmt.Text(record["publication_id"])+".yaml"), data)
}

func (c *Config) VerifyPublicationReceipts(resource Resource) (int, error) {
	dir, err := c.historyDirectory(resource)
	if err != nil {
		return 0, err
	}
	base := filepath.Join(dir, "publications")
	if err = confinedParents(c.RootPath(), filepath.Join(base, "receipt.yaml")); err != nil {
		return 0, err
	}
	files, err := os.ReadDir(base)
	if os.IsNotExist(err) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	if len(files) > 10000 {
		return 0, fmt.Errorf("ASaC publication receipts exceed local record limit")
	}
	count := 0
	for _, file := range files {
		if strings.HasPrefix(file.Name(), ".asac-") {
			continue
		}
		var wrapper struct {
			Format       string         `json:"format"`
			Version      string         `json:"schema_version"`
			UID          string         `json:"resource_uid"`
			Target       string         `json:"target"`
			Workspace    string         `json:"workspace_id"`
			Project      string         `json:"project_id"`
			Receipt      map[string]any `json:"receipt"`
			Authenticity string         `json:"authenticity"`
		}
		if err = readHistoryJSON(c, filepath.Join(base, file.Name()), &wrapper); err != nil {
			return 0, err
		}
		if wrapper.Format != "woobe-observed-publication" || wrapper.Version != "1.0" || wrapper.UID != resource.UID || wrapper.Authenticity != "authenticated_connection_observation" || wrapper.Receipt["project_id"] != wrapper.Project || wrapper.Receipt["resource_uid"] != resource.UID || wrapper.Receipt["kind"] != resource.Kind {
			return 0, fmt.Errorf("ASAC_PUBLICATION_RECEIPT_INVALID: scope mismatch")
		}
		if err = ValidatePublicationReceipt(wrapper.Receipt); err != nil {
			return 0, err
		}
		namespace, digestErr := asac.Digest("destination", []string{wrapper.Target, wrapper.Workspace, wrapper.Project})
		if digestErr != nil {
			return 0, digestErr
		}
		if file.Name() != strings.TrimPrefix(namespace, "sha256:")+"_"+packagefmt.Text(wrapper.Receipt["publication_id"])+".yaml" {
			return 0, fmt.Errorf("ASAC_PUBLICATION_RECEIPT_INVALID: filename identity mismatch")
		}
		count++
	}
	return count, nil
}
