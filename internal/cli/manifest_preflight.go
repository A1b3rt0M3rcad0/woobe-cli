package cli

import (
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
		for _, s := range steps {
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
		return a.emit(map[string]any{"manifest_hash": d.Hash(), "operations": rows, "complete": true, "authorization": "not_evaluated", "path_query_validation": "not_evaluated", "executed": false})
	}})
}
