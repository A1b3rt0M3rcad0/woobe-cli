package cli

import (
	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/controlplane"
	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/output"
	"github.com/spf13/cobra"
)

type partialPages struct {
	Data    any
	Message string
	Code    int
}

func (p *partialPages) Error() string { return p.Message }
func (a *App) requestPagesCommand() {
	var limit int
	c := &cobra.Command{Use: "request-pages <path>", Short: "Read pages through advertised same-route next links", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		if e := controlplane.ValidatePath(args[0]); e != nil {
			return e
		}
		if limit < 1 || limit > 100 {
			return output.New(2, "max-pages must be between 1 and 100")
		}
		if a.DryRun {
			return a.emit(map[string]any{"method": "GET", "path": args[0], "executed": false})
		}
		client, e := a.client()
		if e != nil {
			return e
		}
		q, e := a.query()
		if e != nil {
			return e
		}
		pages, e := client.RequestPages(cmd.Context(), args[0], q, limit)
		if e != nil {
			if pages.Count == 0 {
				return e
			}
			code := 10
			if output.Normalize(e).Code == 130 || output.Normalize(e).Code == 8 {
				code = output.Normalize(e).Code
			}
			return &partialPages{Code: code, Data: pages, Message: "page traversal stopped: " + output.Normalize(e).Message}
		}
		return a.emit(pages)
	}}
	c.Flags().IntVar(&limit, "max-pages", 20, "Bounded page count (1 to 100)")
	a.Root.AddCommand(c)
}
