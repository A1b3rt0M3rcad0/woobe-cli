package cli

import (
	"github.com/spf13/cobra"
)

func (a *App) aliasCommands() {
	for name, target := range map[string][]string{"agent": {"project", "agent"}, "network": {"project", "network"}, "control-key": {"workspace", "control-key"}} {
		name, target := name, target
		a.Root.AddCommand(&cobra.Command{Use: name, Short: "Alias for canonical command", DisableFlagParsing: true, RunE: func(cmd *cobra.Command, args []string) error {
			child := New(a.In, a.Out, a.Err)
			child.Root.SetArgs(append(target, args...))
			return child.Root.ExecuteContext(cmd.Context())
		}})
	}
}
