package cli

import (
	"context"
	"time"

	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/output"
	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/packageapi"
	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/packagefmt"
	"github.com/spf13/cobra"
)

type packageOperationFailure struct {
	Operation packageapi.Operation
	Cause     *output.Error
}

func (e *packageOperationFailure) Error() string { return e.Cause.Error() }

func (a *App) packageClient(ctx context.Context, operation string) (*packageapi.Client, error) {
	control, err := a.client()
	if err != nil {
		return nil, err
	}
	client, err := packageapi.New(control, a.Project)
	if err != nil {
		return nil, err
	}
	capabilities, err := client.Capabilities(ctx)
	if err != nil {
		return nil, err
	}
	if capabilities.SchemaCatalogSHA256 != packagefmt.CatalogDigest() {
		return nil, output.New(9, "Server Package schema catalog is incompatible with this CLI")
	}
	for _, supported := range capabilities.SupportedOperations {
		if supported == operation {
			return client, nil
		}
	}
	return nil, output.New(9, "Server does not advertise the requested Package operation")
}

func (a *App) emitPackageOperation(result packageapi.Operation, err error) error {
	if err != nil {
		return &packageOperationFailure{result, output.Normalize(err)}
	}
	if result.State == "failed" {
		return &packageOperationFailure{result, output.New(7, "Package operation failed; inspect its diagnostics and retained inventory")}
	}
	return output.WriteWithMeta(a.Out, a.Mode, result, map[string]string{"project_id": result.ProjectID}, nil, map[string]any{"complete": result.Terminal})
}

func (a *App) packageOperationCommands(group *cobra.Command) {
	var wait bool
	var deadline time.Duration
	status := &cobra.Command{Use: "status OPERATION_ID", Short: "Observe a durable Package operation", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		if err := packageFlags(cmd); err != nil {
			return err
		}
		if a.DryRun {
			return output.New(2, "Package status requires a remote observation")
		}
		if deadline <= 0 {
			return output.New(2, "Package deadline must be positive")
		}
		ctx, cancel := context.WithTimeout(cmd.Context(), deadline)
		defer cancel()
		client, err := a.packageClient(ctx, "status")
		if err != nil {
			return err
		}
		result, err := client.Status(ctx, args[0])
		if err != nil {
			return err
		}
		if wait {
			result, err = client.Wait(ctx, result, nil)
		}
		return a.emitPackageOperation(result, err)
	}}
	status.Flags().BoolVar(&wait, "wait", false, "Observe until terminal state within the total deadline")
	status.Flags().DurationVar(&deadline, "deadline", 5*time.Minute, "Total observation deadline, including HTTP requests")
	group.AddCommand(status)

	cancel := &cobra.Command{Use: "cancel OPERATION_ID", Short: "Request cancellation without deleting shared resources", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		if err := packageFlags(cmd); err != nil {
			return err
		}
		if a.DryRun {
			return output.New(2, "Package cancellation requires an explicit operation")
		}
		client, err := a.packageClient(cmd.Context(), "cancel")
		if err != nil {
			return err
		}
		result, err := client.Cancel(cmd.Context(), args[0])
		if err != nil {
			return err
		}
		return a.emitPackageOperation(result, nil)
	}}
	group.AddCommand(cancel)

	var revision int64
	resume := &cobra.Command{Use: "resume OPERATION_ID", Short: "Resume an existing Package dependency wait", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		if err := packageFlags(cmd); err != nil {
			return err
		}
		if a.DryRun || revision <= 0 {
			return output.New(2, "Resume requires --revision from the observed operation")
		}
		client, err := a.packageClient(cmd.Context(), "resume")
		if err != nil {
			return err
		}
		result, err := client.Resume(cmd.Context(), args[0], revision)
		if err != nil {
			return err
		}
		return a.emitPackageOperation(result, nil)
	}}
	resume.Flags().Int64Var(&revision, "revision", 0, "Revision of the dependency wait approved for resumption")
	group.AddCommand(resume)
}
