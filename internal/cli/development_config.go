package cli

import (
	"os"
	"path/filepath"

	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/config"
	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/devworkspace"
	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/output"
	"github.com/spf13/cobra"
)

func (a *App) developmentConfig(required bool) (*devworkspace.Config, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	path, err := devworkspace.Discover(cwd, a.ProjectConfig, a.NoProjectConfig)
	if err != nil {
		return nil, output.New(2, err.Error())
	}
	if path == "" {
		if required {
			return nil, output.New(2, "No development registry found; run woobe init for managed YAML development, or use standalone package export/import")
		}
		return nil, nil
	}
	c, err := devworkspace.Load(path)
	if err != nil {
		return nil, output.New(2, "Invalid development config: "+err.Error())
	}
	return c, nil
}

func (a *App) developmentConfigCommands() {
	a.developmentManageCommands()
	a.developmentBulkCommands()
	var root string
	init := &cobra.Command{Use: "init", Short: "Initialize an optional local development registry", Args: cobra.NoArgs, Example: "woobe init --context local --root .woobe", Long: "Create .woobe-config explicitly. Ordinary API, runtime and authentication commands remain independent. Provider files contain public connection intent and credential references only.", RunE: func(cmd *cobra.Command, _ []string) error {
		if a.ProjectConfig != "" || a.NoProjectConfig {
			return output.New(2, "init creates .woobe-config in the current directory")
		}
		user, err := config.Load(a.ConfigPath)
		if err != nil {
			return output.New(2, "Cannot load private connection config")
		}
		context := a.ContextName
		if context == "" {
			context = os.Getenv("WOOBE_CONTEXT")
		}
		if context == "" {
			context = user.Current
		}
		if context != "" {
			if _, ok := user.Contexts[context]; !ok {
				return output.New(2, "Unknown connection; create it with context create or initialize without a context")
			}
		}
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		c, err := devworkspace.Create(cwd, root, context)
		if err != nil {
			return output.New(2, err.Error())
		}
		if err := os.MkdirAll(filepath.Join(c.RootPath(), ".state"), 0700); err != nil {
			return err
		}
		ignore := filepath.Join(c.RootPath(), ".gitignore")
		f, err := os.OpenFile(ignore, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
		if err == nil {
			_, err = f.WriteString(".state/\n")
			closeErr := f.Close()
			if err == nil {
				err = closeErr
			}
		} else if os.IsExist(err) {
			err = nil
		}
		if err != nil {
			return err
		}
		return a.emit(map[string]any{"initialized": true, "config": c.File, "root": c.RootPath(), "context": c.Context, "registry_id": c.RegistryID})
	}}
	init.Flags().StringVar(&root, "root", ".woobe", "Artifact root, relative to the project config or absolute")
	a.Root.AddCommand(init)
	group := a.group("config")
	for _, action := range []string{"show", "check"} {
		action := action
		group.AddCommand(&cobra.Command{Use: action, Short: action + " the optional development configuration", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := a.developmentConfig(true)
			if err != nil {
				return err
			}
			if action == "show" {
				return a.emit(map[string]any{"config": c, "file": c.File, "root": c.RootPath(), "scope": "development_only"})
			}
			graph, err := devworkspace.LoadGraph(c)
			if err != nil {
				return output.New(2, err.Error())
			}
			return a.emit(map[string]any{"valid": true, "file": c.File, "root": c.RootPath(), "resources": len(graph.Nodes), "authorization": "not_evaluated", "executed": false})
		}})
	}
	resources := a.group("resources")
	resources.AddCommand(&cobra.Command{Use: "list", Short: "List registered local development resources", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := a.developmentConfig(true)
		if err != nil {
			return err
		}
		return a.emit(map[string]any{"resources": c.Resources, "root": c.RootPath()})
	}})
	resources.AddCommand(&cobra.Command{Use: "used-by KIND REFERENCE", Short: "Inspect direct and indirect local consumers", Args: cobra.ExactArgs(2), Example: "woobe resources used-by provider \"@openai-main\"", RunE: func(cmd *cobra.Command, args []string) error {
		c, err := a.developmentConfig(true)
		if err != nil {
			return err
		}
		resource, err := c.Resolve(args[0], args[1])
		if err != nil {
			return output.New(2, err.Error())
		}
		graph, err := devworkspace.LoadGraph(c)
		if err != nil {
			return output.New(2, err.Error())
		}
		return a.emit(map[string]any{"resource": resource, "consumers": graph.UsedBy(resource.Key), "scope": "local", "remote_impact": "server_required"})
	}})
}
