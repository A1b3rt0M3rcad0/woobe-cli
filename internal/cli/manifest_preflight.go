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
		doc, e := a.loadServerSchema(cmd.Context())
		if e != nil {
			return e
		}
		steps, _ := d.Order()
		rows := []map[string]any{}
		var failure *output.Error
		complete := true
		for _, s := range steps {
			var body any
			if len(s.Body) > 0 {
				dec := json.NewDecoder(bytes.NewReader(s.Body))
				dec.UseNumber()
				if e = dec.Decode(&body); e != nil {
					return output.New(2, "invalid manifest body")
				}
				body, e = manifest.Resolve(body, map[string]any{})
				if e != nil {
					complete = false
					rows = append(rows, map[string]any{"step_id": s.ID, "body_validation": "deferred_dependency_result", "write_executed": false})
					continue
				}
				s.Body, _ = json.Marshal(body)
			}
			op, _ := a.operation(s.Command)
			def, e := operationDefinition(doc, op)
			if e == nil {
				e = validateBodySchema(doc, def, s.Body)
			}
			if e != nil {
				complete = false
				cause := output.Normalize(e)
				if failure == nil || cause.Code == 9 {
					failure = cause
				}
				rows = append(rows, map[string]any{"step_id": s.ID, "body_validation": "failed", "error": cause, "write_executed": false})
				continue
			}
			rows = append(rows, map[string]any{"step_id": s.ID, "body_validation": "valid", "write_executed": false})
		}
		data := map[string]any{"manifest_hash": d.Hash(), "operations": rows, "complete": complete, "authorization": "not_evaluated", "path_query_validation": "not_evaluated", "executed": false}
		if requireComplete && !complete && failure == nil {
			failure = output.New(9, "dependency body schemas require execution results before complete validation")
		}
		if failure != nil {
			return &preflightFailure{Data: data, Cause: failure}
		}
		return a.emit(data)
	}}
	c.Flags().BoolVar(&requireComplete, "require-complete", false, "Fail when dependency body values cannot yet be validated")
	g.AddCommand(c)
}

type preflightFailure struct {
	Data  map[string]any
	Cause *output.Error
}

func (e *preflightFailure) Error() string { return "manifest body preflight failed" }
