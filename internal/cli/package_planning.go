package cli

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"
	"time"

	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/output"
	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/packageapi"
	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/packagefmt"
	"github.com/spf13/cobra"
)

func packageLifecycle(lifecycle, notes, reason string) error {
	switch lifecycle {
	case "draft", "staging", "release", "production":
	default:
		return output.New(2, "Invalid Package lifecycle")
	}
	if (lifecycle == "release" || lifecycle == "production") && strings.TrimSpace(notes) == "" {
		return output.New(2, "Release and production require --release-notes")
	}
	if lifecycle == "production" && strings.TrimSpace(reason) == "" {
		return output.New(2, "Production requires --reason")
	}
	if len(notes) > 4000 || len(reason) > 1000 {
		return output.New(2, "Package lifecycle audit fields exceed their limits")
	}
	return nil
}

// Package requests use the advertised HTTP contracts while keeping domain
// decisions in the server's owner modules.
func (a *App) validatePackageRequest(doc map[string]any, method, suffix string, params map[string]string, query url.Values, body any) error {
	op := Operation{Method: method, Path: "/projects/{project_id}/packages" + suffix}
	if a.ValidateParameters {
		if err := validateOperationParameters(doc, op, params, query); err != nil {
			return err
		}
	}
	if a.ValidateBody && method != "GET" && method != "HEAD" {
		def, err := operationDefinition(doc, op)
		if err != nil {
			return err
		}
		var encoded []byte
		if body != nil {
			encoded, err = json.Marshal(body)
			if err != nil {
				return output.New(2, "Invalid Package request")
			}
		}
		return validateBodySchema(doc, def, encoded)
	}
	return nil
}

func (a *App) packagePlanningCommands(group *cobra.Command) {
	var bindingsPath, savePath, lifecycle, notes, reason string
	var shortcuts []string
	var locked bool
	var deadline time.Duration
	command := &cobra.Command{Use: "plan SOURCE", Short: "Plan native Package effects and destination bindings without materializing resources", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		if err := packageFlags(cmd, "yes", "validate-body", "validate-parameters", "schema-sha256"); err != nil {
			return err
		}
		if err := packageLifecycle(lifecycle, notes, reason); err != nil {
			return err
		}
		if deadline <= 0 {
			return output.New(2, "Package deadline must be positive")
		}
		if a.DryRun && (a.ValidateBody || a.ValidateParameters || savePath != "") {
			return output.New(2, "Offline dry-run cannot validate remote schemas or save an approved plan")
		}
		bundle, err := loadPackage(args[0], locked)
		if err != nil {
			return err
		}
		defer bundle.Close()
		if err = excludePackageOperationalFiles(args[0], bundle, bindingsPath, savePath); err != nil {
			return err
		}
		bindings, err := packagefmt.LoadBindings(bindingsPath, shortcuts, bundle.Graph)
		if err != nil {
			return packageError(err)
		}
		if a.DryRun {
			return a.emit(packagePublic(map[string]any{"package_schema_version": "1.0", "artifact_digest": bundle.ArtifactDigest, "inventory": bundle.Inventory, "bindings": bindings, "lifecycle": lifecycle, "materialization_order": bundle.Graph.Order, "executed": false, "authorization": "not_evaluated", "semantic_validation": "server_required"}))
		}
		ctx, cancel := context.WithTimeout(cmd.Context(), deadline)
		defer cancel()
		client, err := a.packageClient(ctx, "plan")
		if err != nil {
			return err
		}
		var schema map[string]any
		if a.ValidateBody || a.ValidateParameters {
			schema, err = a.loadServerSchema(ctx)
			if err != nil {
				return err
			}
		}
		if a.ValidateParameters {
			query := url.Values{}
			if locked {
				query.Set("locked", "true")
			}
			// Upload is binary multipart. Parameter validation remains authoritative;
			// JSON body validation applies to the subsequent planning request.
			if err = validateOperationParameters(schema, Operation{Method: "POST", Path: "/projects/{project_id}/packages/uploads"}, map[string]string{"project_id": client.ProjectID}, query); err != nil {
				return err
			}
		}
		// Validate the plan contract before upload using a syntactically valid UUID;
		// the accepted upload identity is validated again before the plan POST.
		request := packageapi.PlanRequest{Mode: "create", UploadID: "00000000-0000-4000-8000-000000000000", ArtifactDigest: bundle.ArtifactDigest, Bindings: bindings, Lifecycle: lifecycle}
		if notes != "" {
			request.ReleaseNotes = &notes
		}
		if reason != "" {
			request.Reason = &reason
		}
		if schema != nil {
			if err = a.validatePackageRequest(schema, "POST", "/plan", map[string]string{"project_id": client.ProjectID}, nil, request); err != nil {
				return err
			}
		}
		upload, err := client.Upload(ctx, bundle, locked)
		if err != nil {
			return err
		}
		request.UploadID = upload.UploadID
		if schema != nil {
			if err = a.validatePackageRequest(schema, "POST", "/plan", map[string]string{"project_id": client.ProjectID}, nil, request); err != nil {
				return err
			}
		}
		approved, err := client.Plan(ctx, request)
		if err != nil {
			return err
		}
		if savePath != "" {
			capabilities, err := client.Capabilities(ctx)
			if err != nil {
				return err
			}
			receipt := packageapi.PlanReceipt{Format: "woobe-package-plan-receipt", SchemaVersion: "1.0", PrincipalFingerprint: capabilities.PrincipalFingerprint, Plan: approved}
			if err = receipt.Save(savePath); err != nil {
				return err
			}
		}
		return a.emit(packagePublic(approved))
	}}
	command.Flags().StringVar(&bindingsPath, "bindings", "", "Destination ImportBindings file; protected references remain unresolved")
	command.Flags().StringArrayVar(&shortcuts, "bind", nil, "Bind credential.ALIAS=UUID without overriding another binding")
	command.Flags().StringVar(&savePath, "save-plan", "", "Exclusive private destination for the approved plan receipt")
	command.Flags().BoolVar(&locked, "locked", false, "Require the captured package's inventory lock")
	command.Flags().StringVar(&lifecycle, "lifecycle", "draft", "Explicit destination lifecycle: draft, staging, release, production")
	command.Flags().StringVar(&notes, "release-notes", "", "Audit notes required for release and production")
	command.Flags().StringVar(&reason, "reason", "", "Activation reason required for production")
	command.Flags().DurationVar(&deadline, "deadline", 5*time.Minute, "Total upload and planning deadline, including HTTP requests")
	group.AddCommand(command)
}
