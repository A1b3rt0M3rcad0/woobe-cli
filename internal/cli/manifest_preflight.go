package cli

import (
	"bytes"
	"encoding/json"
	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/manifest"
	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/output"
	"github.com/spf13/cobra"
)

func (a *App) manifestPreflightCommand(g *cobra.Command) {
	g.AddCommand(&cobra.Command{Use: "preflight", Short: "Read advertised body schemas and validate manifest inputs without mutations", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
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
				return e
			}
			rows = append(rows, map[string]any{"step_id": s.ID, "body_validation": "valid", "write_executed": false})
		}
		return a.emit(map[string]any{"manifest_hash": d.Hash(), "operations": rows, "complete": complete, "authorization": "not_evaluated", "path_query_validation": "not_evaluated", "executed": false})
	}})
}
