package cli

import (
	"bytes"
	"encoding/json"
	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/manifest"
	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/output"
	"github.com/spf13/cobra"
)

func (a *App) manifestPreflightCommand(g *cobra.Command) {
	var requireComplete bool
	c := &cobra.Command{Use: "preflight", Short: "Read advertised body schemas and validate manifest inputs without mutations", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		d, e := a.readManifest()
		if e != nil {
			return e
		}
		if e = a.validateManifest(d); e != nil {
			return e
		}
		v, e := a.resolve()
		if e != nil {
			return e
		}
		if d.Workspace != "" && d.Workspace != v.Workspace || d.Project != "" && d.Project != v.Project {
			return output.New(6, "manifest scope must match selected context")
		}
		if a.DryRun {
			data := map[string]any{"manifest_hash": d.Hash(), "complete": false, "body_validation": "not_evaluated_dry_run", "executed": false, "schema_fetched": false}
			if requireComplete {
				return &preflightFailure{Data: data, Cause: output.New(9, "dry-run does not perform body validation")}
			}
			return a.emitPreflight(data, nil)
		}
		doc, e := a.loadServerSchema(cmd.Context())
		if e != nil {
			return e
		}
		steps, _ := d.Order()
		rows := []map[string]any{}
		var failure *output.Error
		complete := true
		failed := func(id string, e error) {
			complete = false
			cause := output.Normalize(e)
			if failure == nil || cause.Code == 9 {
				failure = cause
			}
			rows = append(rows, map[string]any{"step_id": id, "body_validation": "failed", "error": cause, "write_executed": false})
		}
		for _, s := range steps {
			op, _ := a.operation(s.Command)
			def, e := operationDefinition(doc, op)
			if e != nil {
				failed(s.ID, e)
				continue
			}
			var body any
			if len(s.Body) > 0 {
				dec := json.NewDecoder(bytes.NewReader(s.Body))
				dec.UseNumber()
				if e = dec.Decode(&body); e != nil {
					return output.New(2, "invalid manifest body")
				}
				body, e = manifest.Resolve(body, map[string]any{})
				if e != nil {
					if check := validateBodySchema(doc, def, []byte("null")); check != nil && output.Normalize(check).Code == 9 {
						failed(s.ID, check)
						continue
					}
					complete = false
					rows = append(rows, map[string]any{"step_id": s.ID, "body_validation": "deferred_dependency_result", "write_executed": false})
					continue
				}
				s.Body, _ = json.Marshal(body)
			}
			if e = validateBodySchema(doc, def, s.Body); e != nil {
				failed(s.ID, e)
				continue
			}
			rows = append(rows, map[string]any{"step_id": s.ID, "body_validation": "valid", "write_executed": false})
		}
		data := map[string]any{"manifest_hash": d.Hash(), "schema_sha256": schemaDigest(doc), "operations": rows, "complete": complete, "authorization": "not_evaluated", "validation_direction": "request", "path_query_validation": "not_evaluated", "executed": false}
		if requireComplete && !complete && failure == nil {
			failure = output.New(9, "dependency body schemas require execution results before complete validation")
		}
		if failure != nil {
			return &preflightFailure{Data: data, Cause: failure}
		}
		return a.emitPreflight(data, nil)
	}}
	c.Flags().BoolVar(&requireComplete, "require-complete", false, "Fail when dependency body values cannot yet be validated")
	g.AddCommand(c)
}

type preflightFailure struct {
	Data  map[string]any
	Cause *output.Error
}

func (e *preflightFailure) Error() string { return "manifest body preflight failed" }

func (a *App) emitPreflight(data map[string]any, cause *output.Error) error {
	if a.Mode == "table" {
		return output.Write(a.Out, a.Mode, output.Redact(data), nil, cause)
	}
	env := output.Envelope{SchemaVersion: "1", Success: cause == nil, Context: map[string]string{"workspace_id": a.Workspace, "project_id": a.Project}, Data: output.Redact(data), Error: cause, Meta: map[string]any{"complete": data["complete"] == true}}
	return json.NewEncoder(a.Out).Encode(env)
}
