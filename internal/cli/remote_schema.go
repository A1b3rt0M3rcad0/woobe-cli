package cli

import (
	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/output"
	"github.com/spf13/cobra"
	"strings"
)

func (a *App) remoteSchemaCommand() {
	var operation string
	c := &cobra.Command{Use: "server-schema", Short: "Retrieve the authoritative OpenAPI operation and component schemas", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		var selected *Operation
		for _, op := range a.Registry {
			if op.Command == operation {
				x := op
				selected = &x
				break
			}
		}
		if selected == nil {
			return output.New(2, "unknown operation")
		}
		client, e := a.client()
		if e != nil {
			return e
		}
		v, _, e := client.Request(cmd.Context(), "GET", "/openapi.json", nil, nil)
		if e != nil {
			return e
		}
		doc, ok := v.(map[string]any)
		if !ok {
			return output.New(9, "server OpenAPI document unavailable")
		}
		paths, ok := doc["paths"].(map[string]any)
		if !ok {
			return output.New(9, "server OpenAPI paths unavailable")
		}
		path, ok := paths[selected.Path].(map[string]any)
		if !ok {
			return output.New(9, "operation not advertised by server")
		}
		definition, ok := path[strings.ToLower(selected.Method)]
		if !ok {
			return output.New(9, "method not advertised by server")
		}
		return output.Write(a.Out, a.Mode, map[string]any{"operation": selected, "definition": definition, "components": doc["components"], "source": "server_openapi"}, nil, nil)
	}}
	c.Flags().StringVar(&operation, "command", "", "Canonical command")
	a.Root.AddCommand(c)
}
