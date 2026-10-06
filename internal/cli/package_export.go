package cli

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/output"
	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/packageapi"
	"github.com/spf13/cobra"
)

func (a *App) packageExportCommands(group *cobra.Command) {
	export := &cobra.Command{Use: "export", Short: "Export a complete portable native snapshot"}
	for _, kind := range []string{"agent", "network"} {
		var destination, source, snapshot, knowledge, name, version string
		var deadline time.Duration
		command := &cobra.Command{Use: kind + " TARGET_ID", Short: "Capture an exact " + kind + " source and all portable dependencies", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
			if err := packageFlags(cmd, "yes", "validate-body", "validate-parameters", "schema-sha256"); err != nil {
				return err
			}
			if a.DryRun || destination == "" || (source == "") == (snapshot == "") || deadline <= 0 || name == "" || version == "" {
				return output.New(2, "Export requires --destination, --name, --version and exactly one of --source or --snapshot-id")
			}
			switch source {
			case "", "draft", "staging", "release", "production":
			default:
				return output.New(2, "Invalid Package export source")
			}
			if knowledge != "portable" && knowledge != "binding" {
				return output.New(2, "Knowledge mode must be portable or binding")
			}
			absolute, err := filepath.Abs(destination)
			if err != nil {
				return output.New(2, "Invalid export destination")
			}
			if _, err = os.Lstat(absolute); err == nil {
				return output.New(2, "Export destination already exists")
			} else if !os.IsNotExist(err) {
				return output.New(2, "Export destination is unavailable")
			}
			// Fail before capture if the destination parent cannot hold private staging.
			parent, err := os.Stat(filepath.Dir(absolute))
			if err != nil || !parent.IsDir() {
				return output.New(2, "Export destination parent must exist")
			}
			ctx, cancel := context.WithTimeout(cmd.Context(), deadline)
			defer cancel()
			client, err := a.packageClient(ctx, "export")
			if err != nil {
				return err
			}
			targetKind := "Agent"
			if kind == "network" {
				targetKind = "Network"
			}
			request := packageapi.ExportRequest{Kind: targetKind, TargetID: args[0], Knowledge: knowledge, Name: name, Version: version}
			if source != "" {
				request.Source = &source
			} else {
				request.SnapshotID = &snapshot
			}
			var schema map[string]any
			if a.ValidateBody || a.ValidateParameters {
				schema, err = a.loadServerSchema(ctx)
				if err != nil {
					return err
				}
				if err = a.validatePackageRequest(schema, "POST", "/export", map[string]string{"project_id": client.ProjectID}, nil, request); err != nil {
					return err
				}
			}
			receipt, err := client.Export(ctx, request)
			if err != nil {
				return err
			}
			if schema != nil {
				if err = a.validatePackageRequest(schema, "GET", "/exports/{export_id}/artifact", map[string]string{"project_id": client.ProjectID, "export_id": receipt.ExportID}, nil, nil); err != nil {
					return err
				}
			}
			bundle, err := client.Download(ctx, receipt)
			if err != nil {
				return packageError(err)
			}
			defer bundle.Close()
			if err = ctx.Err(); err != nil {
				return output.New(8, "Package export deadline expired before publication")
			}
			if err = bundle.Publish(absolute); err != nil {
				return packageError(err)
			}
			return a.emit(map[string]any{"package_schema_version": "1.0", "destination": absolute, "export_id": receipt.ExportID, "artifact_digest": receipt.ArtifactDigest, "transport_digest": receipt.TransportDigest, "inventory": receipt.Inventory, "closure_complete": receipt.ClosureComplete, "self_contained": receipt.SelfContained, "knowledge": receipt.Knowledge})
		}}
		command.Flags().StringVar(&destination, "destination", "", "New directory to publish after complete integrity verification")
		command.Flags().StringVar(&source, "source", "", "Native source: draft, staging, release or production")
		command.Flags().StringVar(&snapshot, "snapshot-id", "", "Exact native snapshot belonging to the target")
		command.Flags().StringVar(&knowledge, "knowledge", "portable", "Knowledge portability: portable or binding")
		command.Flags().StringVar(&name, "name", "", "Portable package name")
		command.Flags().StringVar(&version, "version", "", "Portable package version")
		command.Flags().DurationVar(&deadline, "deadline", 5*time.Minute, "Total capture/download deadline")
		export.AddCommand(command)
	}
	group.AddCommand(export)
}
