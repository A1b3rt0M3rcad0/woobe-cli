package cli

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/output"
	"github.com/spf13/cobra"
	"strings"
)

func (a *App) serverOperation(ctx context.Context, command string) (Operation, map[string]any, map[string]any, map[string]any, error) {
	op, ok := a.operation(command)
	if !ok || op.Kind != "http" {
		return op, nil, nil, nil, output.New(2, "server schema requires a canonical HTTP operation")
	}
	doc, e := a.loadServerSchema(ctx)
	if e != nil {
		return op, nil, nil, nil, e
	}
	paths, _ := doc["paths"].(map[string]any)
	path, _ := paths[op.Path].(map[string]any)
	definition, ok := path[strings.ToLower(op.Method)].(map[string]any)
	if !ok {
		return op, nil, nil, nil, output.New(9, "operation not advertised by server")
	}
	return op, doc, path, definition, nil
}
func (a *App) remoteSchemaCommand() {
	var operation string
	c := &cobra.Command{Use: "server-schema", Short: "Retrieve authoritative OpenAPI operation and components", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		op, doc, path, definition, e := a.serverOperation(cmd.Context(), operation)
		if e != nil {
			return e
		}
		return output.Write(a.Out, a.Mode, map[string]any{"operation": op, "definition": definition, "path_parameters": path["parameters"], "components": doc["components"], "openapi_version": doc["openapi"], "source": "server_openapi"}, nil, nil)
	}}
	c.Flags().StringVar(&operation, "command", "", "Canonical HTTP command")
	_ = c.MarkFlagRequired("command")
	a.Root.AddCommand(c)
}

func (a *App) loadServerSchema(ctx context.Context) (map[string]any, error) {
	client, e := a.client()
	if e != nil {
		return nil, e
	}
	v, _, e := client.Request(ctx, "GET", "/openapi.json", nil, nil)
	if e != nil {
		return nil, e
	}
	doc, ok := v.(map[string]any)
	if !ok {
		return nil, output.New(9, "invalid server OpenAPI")
	}
	if a.SchemaSHA != "" && a.SchemaSHA != schemaDigest(doc) {
		return nil, output.New(6, "advertised schema snapshot differs from expected SHA-256")
	}
	return doc, nil
}
func operationDefinition(doc map[string]any, op Operation) (map[string]any, error) {
	paths, _ := doc["paths"].(map[string]any)
	path, _ := paths[op.Path].(map[string]any)
	def, ok := path[strings.ToLower(op.Method)].(map[string]any)
	if !ok {
		return nil, output.New(9, "operation not advertised by server")
	}
	return def, nil
}

func schemaDigest(doc map[string]any) string {
	b, _ := json.Marshal(doc)
	return fmt.Sprintf("%x", sha256.Sum256(b))
}
