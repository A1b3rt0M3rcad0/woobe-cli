package cli

import (
	"strings"

	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/devworkspace"
	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/output"
	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/packageapi"
	"github.com/spf13/cobra"
)

func (a *App) asacHydrate(cmd *cobra.Command, c *devworkspace.Config, state *devworkspace.State, resource *devworkspace.Resource, record *devworkspace.Revision) error {
	client, err := a.packageClient(cmd.Context(), "registry")
	if err != nil {
		return err
	}
	capabilities, err := client.Capabilities(cmd.Context())
	if err != nil {
		return err
	}
	if !capabilities.ASaC.RetainedObjectHydration {
		return output.New(9, "Server does not advertise retained executable object hydration")
	}
	id := state.Bindings[resource.UID].ResourceID
	if id == "" {
		tracking, err := c.ReadTracking(*resource)
		if err != nil {
			return output.New(2, "Object hydration needs a native binding or tracking origin")
		}
		if tracking.Origin.Target != strings.TrimRight(client.Control.Base, "/") || tracking.Origin.Workspace != a.Workspace || tracking.Origin.Project != a.Project {
			return output.New(9, "Tracking origin belongs to another selected destination")
		}
		id = tracking.Origin.ResourceID
	}
	if err := resourceID(id); err != nil {
		return err
	}
	bundle, err := client.HydrateRevision(cmd.Context(), packageapi.RetainedArtifactReceipt{ProjectID: a.Project, Kind: resource.Kind, ResourceID: id, ResourceUID: resource.UID, RevisionID: record.ID, RecordDigest: record.RecordDigest, ArtifactDigest: record.ArtifactDigest, DefinitionDigest: record.DefinitionDigest})
	if err != nil {
		return err
	}
	defer bundle.Close()
	if a.DryRun {
		err = devworkspace.ValidateExecutableObject(*resource, *record, bundle)
	} else {
		err = c.StoreExecutableObject(*resource, *record, bundle)
	}
	if err != nil {
		return output.New(9, err.Error())
	}
	source := "unavailable"
	if _, err := c.ReadAuthorObject(*resource, record); err == nil {
		source = "available_locally"
	}
	return a.emit(map[string]any{"resource_uid": resource.UID, "revision_id": record.ID, "artifact_digest": record.ArtifactDigest, "executable_object": "verified", "hydrated": !a.DryRun, "author_object": source, "bindings_restored": false, "author_files_changed": false, "working_head_changed": false, "production_changed": false, "executed": !a.DryRun})
}
