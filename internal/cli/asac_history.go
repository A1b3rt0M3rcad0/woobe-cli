package cli

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/devworkspace"
	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/output"
	"github.com/spf13/cobra"
)

func (a *App) asacCommands() {
	for _, kind := range []string{"agent", "network"} {
		parent := a.group("develop " + kind)
		var environment string
		current := &cobra.Command{Use: "current REFERENCE", Short: "Observe the exact remote environment selection", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
			remote, _ := cmd.Flags().GetBool("remote")
			if !remote {
				return output.New(2, "current requires --remote; use status for offline observations")
			}
			data, err := a.asacCurrent(cmd, kind, args[0], environment)
			if err != nil {
				return err
			}
			return a.emit(data)
		}}
		current.Flags().StringVar(&environment, "env", "production", "Environment selection to observe: draft, staging or production")
		current.Flags().Bool("remote", true, "Consult the authoritative server; this command is always remote")
		parent.AddCommand(current)
		for _, action := range []string{"history", "heads", "revision"} {
			var message string
			var parents []string
			command := &cobra.Command{Use: action + " REFERENCE [ACTION] [REVISION...]", Short: "Inspect or seal local immutable ASaC " + action, Args: cobra.MinimumNArgs(1), Long: "Local history does not claim remote freshness. Revision create freezes a closed portable package without publishing, staging or activating it. history verify checks immutable record digests, DAG and object availability. Explicit multiple parents describe resolved content; no content is merged automatically.", RunE: func(cmd *cobra.Command, args []string) error {
				if action == "heads" && len(args) != 1 {
					return output.New(2, "heads expects only one resource reference")
				}
				if action != "revision" && (cmd.Flags().Changed("message") || cmd.Flags().Changed("parent")) {
					return output.New(2, "--message and --parent apply only to revision create")
				}
				if action == "revision" && len(args) > 1 && args[1] != "create" && (cmd.Flags().Changed("message") || cmd.Flags().Changed("parent")) {
					return output.New(2, "--message and --parent apply only to revision create")
				}
				c, err := a.developmentConfig(true)
				if err != nil {
					return err
				}
				unlock, err := c.Lock()
				if err != nil {
					return output.New(2, err.Error())
				}
				defer unlock()
				state := &devworkspace.State{API: "", Bindings: map[string]devworkspace.Binding{}, Requirements: map[string]any{}, Credentials: map[string]string{}}
				if a.ContextName == "" && os.Getenv("WOOBE_CONTEXT") == "" {
					a.ContextName = c.Context
				}
				if connection, resolveErr := a.resolve(); resolveErr == nil {
					state, err = c.ReadState(strings.TrimRight(connection.APIURL, "/"), connection.Workspace, connection.Project)
					if err != nil {
						return output.New(2, err.Error())
					}
				}
				resource, err := developmentReference(c, state, kind, args[0], "")
				if err != nil {
					return err
				}
				if action == "heads" || action == "history" && len(args) == 2 && args[1] == "verify" {
					if len(args) > 2 {
						return output.New(2, "unexpected history arguments")
					}
					data, err := c.VerifyHistory(*resource)
					if err != nil {
						return output.New(9, err.Error())
					}
					return a.emit(data)
				}
				if action == "history" {
					if len(args) == 2 && args[1] == "fetch" {
						return a.asacFetch(cmd, kind, c, state, resource)
					}
					if len(args) != 1 {
						return output.New(2, "Use history, history fetch or history verify")
					}
					records, err := c.Revisions(*resource)
					if err != nil {
						return output.New(9, err.Error())
					}
					return a.emit(map[string]any{"resource_uid": resource.UID, "revisions": records, "scope": "local", "remote_status": "unverified"})
				}
				if len(args) < 2 {
					return output.New(2, "Use revision create, show or diff")
				}
				switch args[1] {
				case "create":
					if len(args) != 2 {
						return output.New(2, "unexpected revision create arguments")
					}
					if a.DryRun {
						return output.New(2, "revision create is a local checkpoint; use validate/diff to preview without writes")
					}
					graph, err := devworkspace.LoadGraph(c)
					if err != nil {
						return output.New(2, err.Error())
					}
					record, err := graph.CreateRevision(*resource, state, message, parents)
					if err != nil {
						return output.New(2, err.Error())
					}
					return a.emit(record)
				case "show":
					if len(args) != 3 {
						return output.New(2, "revision show requires one revision ID")
					}
					record, err := c.ReadRevision(*resource, args[2])
					if err != nil {
						return output.New(9, err.Error())
					}
					return a.emit(record)
				case "diff":
					if len(args) != 4 {
						return output.New(2, "revision diff requires two revision IDs")
					}
					left, err := c.ReadRevision(*resource, args[2])
					if err != nil {
						return output.New(9, err.Error())
					}
					right, err := c.ReadRevision(*resource, args[3])
					if err != nil {
						return output.New(9, err.Error())
					}
					return a.emit(map[string]any{"from": left.ID, "to": right.ID, "definition_changed": left.DefinitionDigest != right.DefinitionDigest, "components_before": left.Components, "components_after": right.Components, "scope": "sealed_composition"})
				default:
					return output.New(2, "Use revision create, show or diff")
				}
			}}
			command.Flags().StringVar(&message, "message", "", "Checkpoint message (required for revision create)")
			command.Flags().StringSliceVar(&parents, "parent", nil, "Explicit sealed parent revision IDs (repeatable); default current working revision")
			parent.AddCommand(command)
		}
	}
}

func (a *App) asacCurrent(cmd *cobra.Command, kind, reference, environment string) (map[string]any, error) {
	if environment != "draft" && environment != "staging" && environment != "production" {
		return nil, output.New(2, "Use draft, staging or production")
	}
	id, err := a.nativeReference(kind+"_id", reference)
	if err != nil {
		return nil, err
	}
	if err = resourceID(id); err != nil {
		return nil, err
	}
	control, err := a.client()
	if err != nil {
		return nil, err
	}
	if a.Project == "" {
		return nil, output.New(2, "Select a Project before querying remote history")
	}
	path := fmt.Sprintf("/projects/%s/%ss/%s/asac/current", url.PathEscape(a.Project), kind, url.PathEscape(id))
	response, _, err := control.Request(cmd.Context(), "GET", path, url.Values{"environment": {environment}}, nil)
	if err != nil {
		return nil, err
	}
	data, err := developmentResponseData(response)
	if err != nil {
		return nil, err
	}
	if data["schema_version"] != "1.0" || data["resource_id"] != id || data["environment"] != environment {
		return nil, output.New(9, "ASaC selection identity mismatch")
	}
	return data, nil
}

func (a *App) asacFetch(cmd *cobra.Command, kind string, c *devworkspace.Config, state *devworkspace.State, resource *devworkspace.Resource) error {
	control, err := a.client()
	if err != nil {
		return err
	}
	id := state.Bindings[resource.UID].ResourceID
	if id == "" {
		tracking, err := c.ReadTracking(*resource)
		if err != nil {
			return output.New(2, "History fetch needs a native binding or tracking origin")
		}
		if tracking.Origin.Target != strings.TrimRight(control.Base, "/") || tracking.Origin.Workspace != a.Workspace || tracking.Origin.Project != a.Project {
			return output.New(9, "Tracking origin belongs to another selected destination")
		}
		id = tracking.Origin.ResourceID
	}
	if err = resourceID(id); err != nil {
		return err
	}
	total := 0
	watermarks := map[string]any{}
	for _, stream := range []string{"releases", "deployments"} {
		cursor := ""
		seen := map[string]bool{}
		var watermark any
		for page := 0; page < 100; page++ {
			path := fmt.Sprintf("/projects/%s/%ss/%s/asac/history", url.PathEscape(a.Project), kind, url.PathEscape(id))
			query := url.Values{"stream": {stream}, "limit": {"100"}}
			if cursor != "" {
				query.Set("cursor", cursor)
			}
			response, _, err := control.Request(cmd.Context(), "GET", path, query, nil)
			if err != nil {
				return err
			}
			data, err := developmentResponseData(response)
			if err != nil {
				return err
			}
			if data["schema_version"] != "1.0" || data["resource_id"] != id || data["stream"] != stream {
				return output.New(9, "ASaC history scope mismatch")
			}
			if page == 0 {
				watermark = data["watermark"]
			} else if fmt.Sprint(watermark) != fmt.Sprint(data["watermark"]) {
				return output.New(9, "ASaC history watermark changed during traversal")
			}
			records, ok := data["records"].([]any)
			if !ok || len(records) > 100 {
				return output.New(9, "Invalid history page")
			}
			for _, raw := range records {
				record, ok := raw.(map[string]any)
				if !ok || record["resource_id"] != id {
					return output.New(9, "Invalid history record identity")
				}
				if !a.DryRun {
					if err = c.StoreHistoryReceipt(*resource, strings.TrimRight(control.Base, "/"), a.Workspace, a.Project, record); err != nil {
						return output.New(9, err.Error())
					}
				}
				total++
			}
			hasMore, ok := data["has_more"].(bool)
			if !ok || data["complete"] != !hasMore {
				return output.New(9, "ASaC history completeness mismatch")
			}
			if !hasMore {
				watermarks[stream] = watermark
				break
			}
			cursor, ok = data["next_cursor"].(string)
			if !ok || cursor == "" || seen[cursor] {
				return output.New(9, "ASaC history cursor missing or repeated")
			}
			seen[cursor] = true
			if page == 99 {
				return output.New(9, "ASaC history exceeds traversal bound; retained receipts do not constitute a complete fetch")
			}
		}
	}
	return a.emit(map[string]any{"resource_uid": resource.UID, "resource_id": id, "fetched": total, "watermarks": watermarks, "metadata_complete": true, "coverage": "authorized_resource_metadata", "objects_available": "not_evaluated", "author_files_changed": false, "production_changed": false, "executed": !a.DryRun, "remote_status": "observed_at_fetch"})
}
