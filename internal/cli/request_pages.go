package cli

import (
	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/output"
	"github.com/spf13/cobra"
)

type partialPages struct {
	Data    any
	Message string
}

func (p *partialPages) Error() string { return p.Message }
func (a *App) requestPagesCommand() {
	var limit int
	c := &cobra.Command{Use: "request-pages <path>", Short: "Read pages through advertised same-route next links", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
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
			return &partialPages{Data: pages, Message: "page traversal stopped: " + output.Normalize(e).Message}
		}
		return a.emit(pages)
	}}
	c.Flags().IntVar(&limit, "max-pages", 20, "Bounded page count (1 to 100)")
	a.Root.AddCommand(c)
}
