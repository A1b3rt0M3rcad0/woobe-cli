package packageapi

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"

	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/output"
	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/packagebundle"
)

type ExportRequest struct {
	ReleaseVersion string  `json:"release_version,omitempty"`
	Kind           string  `json:"kind"`
	TargetID       string  `json:"target_id"`
	Source         *string `json:"source,omitempty"`
	SnapshotID     *string `json:"snapshot_id,omitempty"`
	Knowledge      string  `json:"knowledge"`
	Name           string  `json:"name,omitempty"`
	Version        string  `json:"version,omitempty"`
}

type ExportReceipt struct {
	Source               map[string]any                `json:"source,omitempty"`
	Name                 string                        `json:"name,omitempty"`
	Version              string                        `json:"version,omitempty"`
	PackageSchemaVersion string                        `json:"package_schema_version"`
	ExportID             string                        `json:"export_id"`
	ProjectID            string                        `json:"project_id"`
	ArtifactDigest       string                        `json:"artifact_digest"`
	TransportDigest      string                        `json:"transport_digest"`
	SizeBytes            int64                         `json:"size_bytes"`
	Inventory            []packagebundle.InventoryFile `json:"inventory"`
	ClosureComplete      bool                          `json:"closure_complete"`
	SelfContained        bool                          `json:"self_contained"`
	Knowledge            string                        `json:"knowledge"`
	ExpiresAt            string                        `json:"expires_at"`
}

func (c *Client) Export(ctx context.Context, request ExportRequest) (ExportReceipt, error) {
	var receipt ExportReceipt
	if (request.Kind != "Agent" && request.Kind != "Network") || !identifier.MatchString(request.TargetID) || (request.Source == nil) == (request.SnapshotID == nil) {
		return receipt, output.New(2, "Export requires a target and exactly one source selector")
	}
	if err := c.request(ctx, http.MethodPost, "/export", request, &receipt); err != nil {
		return receipt, err
	}
	if !identifier.MatchString(receipt.ExportID) || receipt.ProjectID != c.ProjectID || !receipt.ClosureComplete || !digestPattern.MatchString(receipt.ArtifactDigest) || !digestPattern.MatchString(receipt.TransportDigest) || receipt.SizeBytes <= 0 || receipt.SizeBytes > 128<<20 || len(receipt.Inventory) == 0 || !validExpiry(receipt.ExpiresAt) || receipt.Knowledge != request.Knowledge || (request.Knowledge == "portable" && !receipt.SelfContained) {
		return ExportReceipt{}, output.New(9, "Server returned incomplete Package export evidence")
	}
	return receipt, nil
}

// Download never trusts an arbitrary receipt URL and never publishes partial bytes.
func (c *Client) Download(ctx context.Context, receipt ExportReceipt) (*packagebundle.Bundle, error) {
	if receipt.ProjectID != c.ProjectID || !identifier.MatchString(receipt.ExportID) || !receipt.ClosureComplete || receipt.SizeBytes <= 0 || receipt.SizeBytes > 128<<20 || !digestPattern.MatchString(receipt.TransportDigest) || !digestPattern.MatchString(receipt.ArtifactDigest) || !validExpiry(receipt.ExpiresAt) {
		return nil, output.New(9, "Invalid Package export identity")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.Control.Base+c.path("/exports/"+receipt.ExportID+"/artifact"), nil)
	if err != nil {
		return nil, output.New(2, "Invalid Package export request")
	}
	req.Header = c.Control.Headers.Clone()
	req.Header.Del("If-Match")
	req.Header.Del("Idempotency-Key")
	req.Header.Set("Accept", "application/gzip")
	if c.Control.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Control.Token)
	}
	response, err := c.Control.HTTP.Do(req)
	if err != nil {
		return nil, output.New(7, "Package artifact download failed")
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, output.New(7, "Server rejected Package artifact download")
	}
	if c.Control.ResponseHook != nil {
		if err = c.Control.ResponseHook(response); err != nil {
			return nil, output.New(10, "Package response session could not be persisted")
		}
	}
	if response.Header.Get("X-Woobe-Transport-Sha256") != receipt.TransportDigest || response.Header.Get("X-Woobe-Artifact-Digest") != receipt.ArtifactDigest || response.Header.Get("Content-Length") != fmt.Sprint(receipt.SizeBytes) {
		return nil, output.New(9, "Package download headers differ from its receipt")
	}
	archive, err := os.CreateTemp("", "woobe-package-download-")
	if err != nil {
		return nil, output.New(10, "Private Package download staging is unavailable")
	}
	defer os.Remove(archive.Name())
	defer archive.Close()
	digest := sha256.New()
	size, err := io.Copy(io.MultiWriter(archive, digest), io.LimitReader(response.Body, receipt.SizeBytes+1))
	if err != nil || size != receipt.SizeBytes || hex.EncodeToString(digest.Sum(nil)) != receipt.TransportDigest {
		return nil, output.New(9, "Package download failed transport integrity verification")
	}
	if _, err = archive.Seek(0, io.SeekStart); err != nil {
		return nil, output.New(10, "Package staging could not be read")
	}
	bundle, err := packagebundle.ReceiveArchive(archive, true)
	if err != nil {
		return nil, err
	}
	if bundle.ArtifactDigest != receipt.ArtifactDigest || !reflect.DeepEqual(bundle.Inventory, receipt.Inventory) {
		bundle.Close()
		return nil, output.New(9, "Downloaded Package closure differs from its receipt")
	}
	return bundle, nil
}
