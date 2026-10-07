package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"strings"
)

func (a *App) aliasCommands() {
	for name, target := range map[string][]string{"agent": {"project", "agent"}, "network": {"project", "network"}, "control-key": {"workspace", "control-key"}} {
		name, target := name, target
		a.Root.AddCommand(&cobra.Command{Use: name, Short: "Alias for canonical command", DisableFlagParsing: true, RunE: func(cmd *cobra.Command, args []string) error {
			child := New(a.In, a.Out, a.Err)
			// Persistent flags supplied before an alias belong to its parent.
			// Forward them so presentation and connection flags have one meaning.
			forwarded := []string{}
			a.Root.PersistentFlags().Visit(func(flag *pflag.Flag) {
				if flag.Name == "query" {
					for _, value := range a.Query {
						forwarded = append(forwarded, "--query", value)
					}
				} else if flag.Name == "fields" {
					forwarded = append(forwarded, "--fields="+strings.Join(a.OutputFields, ","))
				} else {
					forwarded = append(forwarded, "--"+flag.Name+"="+flag.Value.String())
				}
			})
			child.Root.SetArgs(append(append(forwarded, target...), args...))
			err := child.Root.ExecuteContext(cmd.Context())
			a.Mode, a.OutputFields, a.OutputWide = child.Mode, child.OutputFields, child.OutputWide
			a.Workspace, a.Project, a.outputCommand = child.Workspace, child.Project, child.outputCommand
			return err
		}})
	}
}
