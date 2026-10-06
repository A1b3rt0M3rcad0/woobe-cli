package cli

import (
	"errors"
	"os"
	"strings"

	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/output"
	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/packagebundle"
	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/packagefmt"
	"github.com/spf13/cobra"
)

type packageFailure struct{ Diagnostic *packagefmt.Diagnostic }

func (e *packageFailure) Error() string { return e.Diagnostic.Error() }

func packageError(err error) error {
	var diagnostic *packagefmt.Diagnostic
	if !errors.As(err, &diagnostic) {
		diagnostic = &packagefmt.Diagnostic{Code: "PACKAGE_PATH_INVALID", Message: "Package input could not be read safely"}
	}
	return &packageFailure{diagnostic}
}

func packageFlags(cmd *cobra.Command, allowed ...string) error {
	permitted := map[string]bool{}
	for _, name := range allowed {
		permitted[name] = true
	}
	for _, name := range []string{"file", "query", "secret-file", "if-match", "idempotency-key", "schema-sha256", "validate-body", "validate-parameters", "yes", "runtime-credential"} {
		if cmd.Flags().Changed(name) && !permitted[name] {
			return output.New(2, "--"+name+" is incompatible with this package command")
		}
	}
	return nil
}

func loadPackage(source string, locked bool) (*packagebundle.Bundle, error) {
	if strings.HasSuffix(source, ".tar.gz") || strings.HasSuffix(source, ".tgz") {
		file, err := os.Open(source)
		if err != nil {
			return nil, packageError(err)
		}
		defer file.Close()
		bundle, err := packagebundle.ReceiveArchive(file, locked)
		if err != nil {
			return nil, packageError(err)
		}
		return bundle, nil
	}
	bundle, err := packagebundle.Load(source, locked)
	if err != nil {
		return nil, packageError(err)
	}
	return bundle, nil
}

func (a *App) packageCommands() {
	group := a.group("package")
	var locked bool
	validate := &cobra.Command{
		Use: "validate SOURCE", Short: "Validate a portable package locally without credentials or HTTP",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := packageFlags(cmd); err != nil {
				return err
			}
			bundle, err := loadPackage(args[0], locked)
			if err != nil {
				return err
			}
			defer bundle.Close()
			return a.emit(map[string]any{
				"package_schema_version": "1.0", "schema_catalog_sha256": packagefmt.CatalogDigest(),
				"valid": true, "artifact_digest": bundle.ArtifactDigest,
				"inventory": bundle.Inventory, "entrypoint": packagefmt.Object(bundle.Graph.Manifest["spec"])["entrypoint"],
				"materialization_order": bundle.Graph.Order, "diagnostics": []any{},
				"authorization": "not_evaluated", "semantic_validation": "server_required", "executed": false,
			})
		},
	}
	validate.Flags().BoolVar(&locked, "locked", false, "Require and verify woobe.lock.json without rewriting it")
	group.AddCommand(validate)
}
