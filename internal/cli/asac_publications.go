package cli

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/devworkspace"
	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/output"
	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/packagefmt"
	"github.com/spf13/cobra"
)

func validatePublicationAcceptance(data map[string]any, operation, id, uid, candidate, evaluation, record, runtime string) error {
	if err := devworkspace.ValidatePublicationReceipt(data); err != nil {
		return output.New(9, "Publication acceptance is uncertain; run draft reconcile: "+err.Error())
	}
	if data["operation_id"] != operation || data["resource_id"] != id || data["resource_uid"] != uid || data["candidate_id"] != candidate || data["evaluation_id"] != evaluation || data["record_digest"] != record || data["runtime_digest"] != runtime {
		return output.New(9, "Publication receipt differs from the original request; run draft reconcile")
	}
	return nil
}

func (a *App) asacPublicationCommands() {
	command := &cobra.Command{Use: "publication REFERENCE PUBLICATION_UUID", Short: "Inspect an immutable exact Candidate publication receipt", Args: cobra.ExactArgs(2), Long: "Read server-owned publication provenance. This is separate from Production selection, native Release reuse and declared Git claims. Native UUID inspection does not require project configuration.", RunE: func(cmd *cobra.Command, args []string) error {
		if !uuidReference(args[1]) {
			return output.New(2, "Use the Publication UUID returned by publish --candidate")
		}
		id, err := a.nativeReference("agent_id", args[0])
		if err != nil {
			return err
		}
		if err = resourceID(id); err != nil {
			return err
		}
		if _, err = a.resolve(); err != nil {
			return err
		}
		path := fmt.Sprintf("/projects/%s/agents/%s/asac/publications/%s", url.PathEscape(a.Project), url.PathEscape(id), url.PathEscape(args[1]))
		if a.DryRun {
			return a.emit(map[string]any{"method": "GET", "path": path, "executed": false})
		}
		control, err := a.client()
		if err != nil {
			return err
		}
		response, _, err := control.Request(cmd.Context(), "GET", path, nil, nil)
		if err != nil {
			return err
		}
		data, err := developmentResponseData(response)
		if err != nil {
			return err
		}
		if err = devworkspace.ValidatePublicationReceipt(data); err != nil {
			return output.New(9, err.Error())
		}
		if data["resource_id"] != id || data["publication_id"] != args[1] || data["project_id"] != a.Project || data["kind"] != "Agent" {
			return output.New(9, "Publication receipt belongs to another resource or destination")
		}
		return a.emit(data)
	}}
	a.group("develop agent").AddCommand(command)
}

func (a *App) publishASaCCandidate(cmd *cobra.Command, kind, reference, candidate, evaluation, reason string) error {
	if !uuidReference(candidate) || !uuidReference(evaluation) {
		return output.New(2, "Exact publication requires --candidate UUID and --evaluation UUID")
	}
	if len(strings.TrimSpace(reason)) < 3 {
		return output.New(2, "Supply a publication reason with --notes (at least three characters)")
	}
	if a.File != "" {
		return output.New(2, "Publication uses exact server Candidate and Evaluation evidence; omit --file")
	}
	if !a.Yes && !a.DryRun {
		return output.New(2, "Exact publication requires --yes")
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
	if a.ContextName == "" && os.Getenv("WOOBE_CONTEXT") == "" {
		a.ContextName = c.Context
	}
	control, err := a.client()
	if err != nil {
		return err
	}
	state, err := c.ReadState(strings.TrimRight(control.Base, "/"), a.Workspace, a.Project)
	if err != nil {
		return output.New(2, err.Error())
	}
	resource, err := developmentReference(c, state, kind, reference, "")
	if err != nil {
		return err
	}
	id := state.Bindings[resource.UID].ResourceID
	if !uuidReference(id) {
		return output.New(2, "Recover the native binding before publication")
	}
	pendingPath := asacPrivatePath(c, state, resource.UID, "pending")
	if raw, readErr := c.ReadOperationalFile(pendingPath, 16<<20); readErr == nil {
		var prior asacPending
		if json.Unmarshal(raw, &prior) != nil || prior.OperationID != "" {
			return output.New(9, "An ASaC write is unresolved; run draft reconcile before publication")
		}
	} else if !os.IsNotExist(readErr) {
		return output.New(2, readErr.Error())
	}
	request := func(method, path string, body any) (map[string]any, error) {
		var encoded []byte
		if body != nil {
			encoded, err = json.Marshal(body)
			if err != nil {
				return nil, err
			}
		}
		response, _, requestErr := control.Request(cmd.Context(), method, path, nil, encoded)
		if requestErr != nil {
			return nil, requestErr
		}
		return developmentResponseData(response)
	}
	base := fmt.Sprintf("/projects/%s/%ss/%s/asac", url.PathEscape(a.Project), kind, url.PathEscape(id))
	capabilities, err := request("GET", fmt.Sprintf("/projects/%s/packages/capabilities", url.PathEscape(a.Project)), nil)
	if err != nil {
		return err
	}
	protocol := packagefmt.Object(capabilities["asac"])
	supported := false
	for _, value := range packagefmt.List(protocol["candidate_publication_kinds"]) {
		if value == resource.Kind {
			supported = true
		}
	}
	if protocol["schema_version"] != "1.0" || !supported {
		return output.New(9, "Backend does not support exact Candidate publication for this kind; update it first")
	}
	frozen, err := request("GET", base+"/candidates/"+url.PathEscape(candidate), nil)
	if err != nil {
		return err
	}
	if frozen["state"] != "ready" || frozen["candidate_id"] != candidate || frozen["resource_id"] != id || frozen["resource_uid"] != resource.UID {
		return output.New(9, "Candidate is not ready or belongs to another resource")
	}
	qualified, err := request("GET", base+"/evaluations/"+url.PathEscape(evaluation), nil)
	if err != nil {
		return err
	}
	if qualified["state"] != "passed" || qualified["evaluation_id"] != evaluation || qualified["candidate_id"] != candidate || qualified["resource_id"] != id || qualified["record_digest"] != frozen["record_digest"] || qualified["runtime_digest"] != frozen["runtime_digest"] {
		return output.New(6, "Evaluation does not qualify this exact Candidate; inspect its original native Runs")
	}
	operation, err := devworkspace.NewID()
	if err != nil {
		return err
	}
	body := map[string]any{"operation_id": operation, "candidate_id": candidate, "evaluation_id": evaluation, "reason": reason}
	path := base + "/publications"
	if a.DryRun {
		return a.emit(map[string]any{"method": "POST", "path": path, "body": body, "executed": false, "production_changed": false})
	}
	pendingBody := map[string]any{}
	for key, value := range body {
		pendingBody[key] = value
	}
	pendingBody["runtime_digest"] = frozen["runtime_digest"]
	raw, _ := json.Marshal(asacPending{OperationID: operation, RecordDigest: packagefmt.Text(frozen["record_digest"]), Method: "POST", Path: path, Body: pendingBody})
	if err = c.WriteOperationalFile(pendingPath, raw); err != nil {
		return output.New(10, "Cannot persist publication identity before acceptance")
	}
	receipt, err := request("POST", path, body)
	if err != nil {
		normalized := output.Normalize(err)
		if normalized.Outcome == "rejected" || normalized.Outcome == "not_attempted" {
			if clearErr := c.WriteOperationalFile(pendingPath, []byte("{}")); clearErr != nil {
				return output.New(10, "Publication rejected; its local receipt could not be cleared")
			}
		}
		return err
	}
	if err = validatePublicationAcceptance(receipt, operation, id, resource.UID, candidate, evaluation, packagefmt.Text(frozen["record_digest"]), packagefmt.Text(frozen["runtime_digest"])); err != nil {
		return err
	}
	if err = c.StorePublicationReceipt(*resource, strings.TrimRight(control.Base, "/"), a.Workspace, a.Project, receipt); err != nil {
		return output.New(10, "Published; documentary receipt is pending. Run draft reconcile without publishing again: "+err.Error())
	}
	if err = c.WriteOperationalFile(pendingPath, []byte("{}")); err != nil {
		return output.New(10, "Published; run draft reconcile before another write")
	}
	return a.emit(receipt)
}
