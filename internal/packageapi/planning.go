package packageapi

import (
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"time"

	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/output"
	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/packagebundle"
)

var digestPattern = regexp.MustCompile(`^[a-f0-9]{64}$`)

func validExpiry(value string) bool {
	expires, err := time.Parse(time.RFC3339Nano, value)
	return err == nil && expires.After(time.Now())
}

type Upload struct {
	PackageSchemaVersion string                        `json:"package_schema_version"`
	UploadID             string                        `json:"upload_id"`
	ProjectID            string                        `json:"project_id"`
	ArtifactDigest       string                        `json:"artifact_digest"`
	Inventory            []packagebundle.InventoryFile `json:"inventory"`
	ExpiresAt            string                        `json:"expires_at"`
}

func (c *Client) Upload(ctx context.Context, bundle *packagebundle.Bundle, locked bool) (Upload, error) {
	reader, writer := io.Pipe()
	defer reader.Close()
	form := multipart.NewWriter(writer)
	completed := make(chan error, 1)
	go func() {
		part, err := form.CreateFormFile("file", "package.tar.gz")
		if err == nil {
			err = bundle.Archive(part, true)
		}
		if err == nil {
			err = form.Close()
		}
		_ = writer.CloseWithError(err)
		completed <- err
	}()
	query := url.Values{}
	if locked {
		query.Set("locked", "true")
	}
	value, _, err := c.Control.RequestReader(ctx, http.MethodPost, c.path("/uploads"), query, reader, form.FormDataContentType())
	_ = reader.Close()
	captureErr := <-completed
	if err != nil {
		return Upload{}, err
	}
	if captureErr != nil {
		return Upload{}, output.New(2, "Captured package changed during upload")
	}
	envelope, ok := value.(map[string]any)
	if !ok || envelope["success"] != true {
		return Upload{}, output.New(9, "Invalid Package upload envelope")
	}
	encoded, _ := json.Marshal(envelope["data"])
	var result Upload
	if json.Unmarshal(encoded, &result) != nil || result.PackageSchemaVersion != "1.0" || result.ProjectID != c.ProjectID ||
		!identifier.MatchString(result.UploadID) || result.ArtifactDigest != bundle.ArtifactDigest || !reflect.DeepEqual(result.Inventory, bundle.Inventory) || !validExpiry(result.ExpiresAt) {
		return Upload{}, output.New(9, "Server accepted a different Package artifact")
	}
	return result, nil
}

type PlanRequest struct {
	Mode           string         `json:"mode"`
	UploadID       string         `json:"upload_id"`
	ArtifactDigest string         `json:"artifact_digest"`
	Bindings       map[string]any `json:"bindings"`
	Lifecycle      string         `json:"lifecycle"`
	ReleaseNotes   *string        `json:"release_notes"`
	Reason         *string        `json:"reason"`
}
type Effect struct {
	Key                 string   `json:"key"`
	Owner               string   `json:"owner"`
	Action              string   `json:"action"`
	Component           string   `json:"component,omitempty"`
	Permission          string   `json:"permission"`
	RequiredPermissions []string `json:"required_permissions,omitempty"`
	State               string   `json:"state"`
}
type Plan struct {
	APIOrigin            string         `json:"api_origin"`
	PackageSchemaVersion string         `json:"package_schema_version"`
	PlanID               string         `json:"plan_id"`
	PlanDigest           string         `json:"plan_digest"`
	ProjectID            string         `json:"project_id"`
	UploadID             string         `json:"upload_id"`
	ArtifactDigest       string         `json:"artifact_digest"`
	DefinitionDigest     string         `json:"definition_digest"`
	CapabilitiesDigest   string         `json:"capabilities_digest"`
	BindingsDigest       string         `json:"bindings_digest"`
	Bindings             map[string]any `json:"bindings"`
	Lifecycle            string         `json:"lifecycle"`
	ReleaseNotes         *string        `json:"release_notes"`
	Reason               *string        `json:"reason"`
	Effects              []Effect       `json:"effects"`
	RequiredPermissions  []string       `json:"required_permissions"`
	CreatedAt            string         `json:"created_at"`
	ExpiresAt            string         `json:"expires_at"`
}

func (c *Client) Plan(ctx context.Context, request PlanRequest) (Plan, error) {
	var result Plan
	err := c.request(ctx, http.MethodPost, "/plan", request, &result)
	if err != nil {
		return result, err
	}
	if !identifier.MatchString(result.PlanID) || result.ProjectID != c.ProjectID || result.UploadID != request.UploadID ||
		result.ArtifactDigest != request.ArtifactDigest || result.Lifecycle != request.Lifecycle || !digestPattern.MatchString(result.PlanDigest) || !digestPattern.MatchString(result.DefinitionDigest) ||
		!digestPattern.MatchString(result.CapabilitiesDigest) || !digestPattern.MatchString(result.BindingsDigest) ||
		result.APIOrigin != c.Control.Base || !reflect.DeepEqual(result.Bindings, request.Bindings) ||
		!reflect.DeepEqual(result.ReleaseNotes, request.ReleaseNotes) || !reflect.DeepEqual(result.Reason, request.Reason) ||
		!validExpiry(result.ExpiresAt) || len(result.Effects) == 0 {
		return Plan{}, output.New(9, "Server returned a different Package plan identity")
	}
	return result, nil
}
