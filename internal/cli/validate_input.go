package cli

import (
	"errors"
	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/output"
	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/schemacheck"
	"github.com/spf13/cobra"
)

func (a *App) validateInputCommand() {
	var command string
	c := &cobra.Command{Use: "validate-input", Short: "Validate a JSON body against advertised server schema without executing the operation", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		b, e := a.body(true)
		if e != nil {
			return e
		}
		op, doc, _, def, e := a.serverOperation(cmd.Context(), command)
		if e != nil {
			return e
		}
		if e = validateBodySchema(doc, def, b); e != nil {
			return e
		}
		return a.emit(map[string]any{"valid": true, "operation": op.Command, "source": "server_openapi", "validation": "supported_schema_subset", "authorization": "not_evaluated", "executed": false})
	}}
	c.Flags().StringVar(&command, "command", "", "Canonical HTTP command")
	_ = c.MarkFlagRequired("command")
	a.Root.AddCommand(c)
}
func schemaError(e error) error {
	var s *schemacheck.Error
	if errors.As(e, &s) {
		code := 2
		if s.Unsupported {
			code = 9
		}
		return output.New(code, s.Error())
	}
	return output.New(9, "schema cannot be evaluated")
}
