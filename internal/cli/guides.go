package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"strings"
)

type commandGuide struct{ Summary, Details, Examples string }

func packageVisibleFlag(path, name string) bool {
	if !strings.HasPrefix(path, "package ") {
		return true
	}
	switch name {
	case "file", "query", "secret-file", "if-match", "runtime-credential":
		return false
	case "idempotency-key":
		return path == "package import"
	case "validate-body", "validate-parameters", "schema-sha256":
		return path != "package validate" && path != "package bindings"
	case "dry-run":
		return path == "package validate" || path == "package bindings" || path == "package plan" || path == "package import"
	}
	return true
}

func (a *App) installCommandGuides() {
	guides := map[string]commandGuide{
		"package":              {"Export and import portable YAML packages", "Export -> validate -> bindings (if needed) -> plan (optional review) -> import -> status. Uses the selected connection/project. Imports create resources in draft by default. Requires a compatible Woobe Package backend.", "woobe package export agent \"Support Agent\"\nwoobe package validate ./support-agent --locked\nwoobe package bindings ./support-agent --destination ./destination.yaml\nwoobe package import ./support-agent --bindings ./destination.yaml --wait"},
		"package validate":     {"Validate a YAML package offline", "SOURCE accepts a package directory, its woobe.yaml, or a .tar.gz/.tgz archive. --locked verifies the exported inventory: edited files cannot match their old lock. YAML descriptors use one document, without anchors, aliases or custom tags. Local validity does not establish server authorization.", "woobe package validate ./support-agent\nwoobe package validate ./support-agent/woobe.yaml --locked"},
		"package bindings":     {"Generate destination bindings YAML offline", "Outputs a template based on spec.requires. Replace every placeholder with destination identities/private credential references. No secrets are exported. Save bindings outside the package directory. Credentials can be listed with project provider-credential list. The template is not ready to import until filled.", "woobe package bindings ./support-agent --destination ./destination.yaml\nwoobe project provider-credential list\nwoobe package plan ./support-agent --bindings ./destination.yaml"},
		"package plan":         {"Review import effects and destination bindings", "Does not materialize resources. Every requirement in woobe.yaml must have a matching binding. --bind credential.ALIAS=UUID maps an alias to an existing destination credential. --bindings also supports project_environment, secrets and knowledge. Save bindings/plans outside the package. Default lifecycle: draft. release requires --release-notes; production also requires --reason.", "woobe package plan ./support-agent --bindings ./destination.yaml --save-plan ./approved-plan.json\nwoobe package plan ./support-agent --bind credential.primary-key=UUID\nwoobe package plan ./support-agent --bindings ./destination.yaml --dry-run"},
		"package import":       {"Import a portable package into the selected project", "Creates destination resources, draft by default. SOURCE and --plan-file are exclusive. Checkpoints are automatic in private CLI state, bound to the API/project/principal/input/lifecycle. Repeating the same input reconciles the original operation, including after completion; use a new explicit --checkpoint for an intentional second import. An uncertain acceptance is never replayed. Keep bindings outside the package. Saved plans cannot be changed with lifecycle/bindings flags. A timeout does not cancel the remote operation.", "woobe package import ./support-agent --wait\nwoobe package import ./support-agent --bindings ./destination.yaml --wait\nwoobe package import --plan-file ./approved-plan.json --wait\nwoobe package status --checkpoint PATH_FROM_IMPORT_OUTPUT --wait"},
		"package status":       {"Observe an existing import", "Use an operation ID or its private checkpoint. --wait observes until terminal state. Local timeout does not cancel the operation. A checkpoint with unknown acceptance performs lookup, never another Apply.", "woobe package status OPERATION_ID --wait\nwoobe package status --checkpoint ./import-checkpoint.json --wait"},
		"package resume":       {"Resume an import waiting on dependencies", "Use the observed revision with --revision; a checkpoint can provide its last revision. This continues the existing operation, not a new import.", "woobe package resume OPERATION_ID --revision 3\nwoobe package resume --checkpoint ./import-checkpoint.json"},
		"package cancel":       {"Request cancellation of an existing import", "Cancellation does not delete shared resources. Inspect the final status and retained inventory.", "woobe package cancel OPERATION_ID\nwoobe package status OPERATION_ID --wait"},
		"project agent export": {"Export a partial JSON projection for inspection", "This projection is incomplete and not importable. For a portable YAML package with dependencies, use woobe package export agent NAME_OR_UUID. --destination is a file for this command; --output selects terminal formatting.", "woobe agent export UUID --destination ./agent-projection.json\nwoobe package export agent UUID\nwoobe package export agent \"Support Agent\" --env production"},
		"export":               {"Export a partial resource projection", "For portable YAML composition, use woobe package export agent|network NAME_OR_UUID. --output selects terminal formatting; --destination selects a file/directory according to the command.", "woobe package export agent UUID\nwoobe package export network UUID --env staging"},
		"manifest":             {"Compose explicit API resource operations", "Manifest JSON operations and portable YAML Packages are separate formats. For moving an Agent/Network with its dependencies, use package export/import.", "woobe manifest --help\nwoobe package --help"},
		"context":              {"Manage named Woobe connections", "Each connection keeps its own API URL and CLI key/project selection. Create a connection, select it, then authenticate. A key with one eligible project selects it automatically.", "woobe context create local --api-url http://localhost:8000\nwoobe context use local\nwoobe auth login --cli-key\nwoobe context project select UUID"},
		"auth login":           {"Authenticate the selected Woobe connection", "--cli-key reads a masked terminal prompt. --cli-key --stdin reads a key from stdin for automation. The key is stored privately and resolves workspace/project grants. Never put the key in command arguments.", "woobe context use local\nwoobe auth login --cli-key\nwoobe auth status"},
		"auth status":          {"Show authentication and selected project", "Shows the selected connection and project without revealing the CLI key. --wide provides grant metadata.", "woobe auth status\nwoobe auth status --output compact"},
		"schema":               {"Inspect a command input/output schema", "Use the canonical command path. This describes CLI inputs; server-schema provides advertised API payload schemas. API --file accepts JSON; portable package descriptors accept YAML/JSON.", "woobe schema --command \"project agent create\" --kind input\nwoobe server-schema --output json"},
		"help":                 {"Discover focused command guidance", "Use COMMAND --help for readable examples. Use help CANONICAL_COMMAND --output compact for concise structured guidance. --wide/--output json includes the full command metadata.", "woobe package export agent --help\nwoobe help package export agent --output compact\nwoobe agent create --help"},
	}
	for _, kind := range []string{"agent", "network"} {
		guides["package export "+kind] = commandGuide{"Export a portable " + kind + " YAML package", "NAME_OR_UUID accepts an exact name or UUID in the selected project. Default environment: draft. staging/production capture the currently assigned version; release requires --version. No fallback to another environment. List native release versions with woobe " + kind + " " + map[string]string{"agent": "release", "network": "version"}[kind] + " list UUID. --destination defaults to ./normalized-resource-name and never overwrites existing paths. Package name/version are derived; override with --name/--package-version. --output selects response formatting, not a path. --snapshot-id and --source are advanced legacy selectors.", "woobe package export " + kind + " \"Support\"\nwoobe package export " + kind + " UUID --env staging\nwoobe package export " + kind + " UUID --env release --version v1.0.20261007.01\nwoobe package export " + kind + " UUID --env production\nwoobe package export " + kind + " UUID --destination ./backup"}
		guides["project "+kind] = commandGuide{"Manage project " + kind + "s", "Use get UUID to inspect a resource. Use package export for portable YAML with dependencies. Output: --fields for selected fields, --wide for details, --output compact for concise JSON, --output json for the full contract.", "woobe " + kind + " list\nwoobe " + kind + " get UUID\nwoobe package export " + kind + " UUID"}
		guides["project "+kind+" list"] = commandGuide{"List " + kind + "s in the selected project", "Names may be duplicated; use the UUID for an unambiguous target.", "woobe " + kind + " list\nwoobe " + kind + " list --fields name,id --output compact"}
		guides["project "+kind+" get"] = commandGuide{"Inspect one " + kind, "Use --wide for complete configuration, or --fields to select useful fields.", "woobe " + kind + " get UUID\nwoobe " + kind + " get UUID --wide"}
		for _, action := range []string{"create", "update"} {
			suffix := ""
			if action == "update" {
				suffix = " UUID"
			}
			guides["project "+kind+" "+action] = commandGuide{action + " a project " + kind, "--file is a JSON API payload, not a portable YAML package. Consult the advertised server schema for authoritative fields. --dry-run shows the request without executing it; --validate-body validates the server payload contract.", "woobe " + kind + " " + action + suffix + " --file ./payload.json --dry-run\nwoobe " + kind + " " + action + suffix + " --file ./payload.json --validate-body"}
		}
	}
	guides["project agent release list"] = commandGuide{"List native Agent release versions", "Use the exact version value with package export --env release --version. Details are available through --wide or --output json.", "woobe agent release list UUID\nwoobe package export agent UUID --env release --version VERSION_FROM_LIST"}
	guides["project network version list"] = commandGuide{"List native Network versions", "Uses the selected project automatically. Only kind=release versions can be selected with --env release; staging/production export current assignments directly.", "woobe network version list UUID\nwoobe package export network UUID --env release --version VERSION_FROM_LIST"}
	guides["context create"] = commandGuide{"Create a named Woobe connection", "Supply a name and API origin; then select the connection and log in with its CLI key.", "woobe context create local --api-url http://localhost:8000\nwoobe context use local\nwoobe auth login --cli-key"}
	guides["context use"] = commandGuide{"Select a Woobe connection", "Uses the connection's saved API origin, CLI key and project.", "woobe context list\nwoobe context use local\nwoobe auth status"}
	guides["context project select"] = commandGuide{"Select a project granted to this connection", "Accepts eligible project name, slug or UUID. Does not require entering workspace identity.", "woobe context project select UUID\nwoobe context project select chatbot"}
	guides["context show"] = commandGuide{"Inspect the selected connection", "Connection settings are separate from the private key store.", "woobe context show\nwoobe context show local"}
	var walk func(*cobra.Command)
	walk = func(cmd *cobra.Command) {
		path := strings.TrimPrefix(cmd.CommandPath(), "woobe ")
		if guide, ok := guides[path]; ok {
			cmd.Short = guide.Summary
			cmd.Long = guide.Summary + ".\n\n" + guide.Details
			cmd.Example = guide.Examples
		}
		if strings.HasPrefix(path, "package") {
			cmd.SetUsageFunc(packageUsage)
		}
		for _, child := range cmd.Commands() {
			walk(child)
		}
	}
	walk(a.Root)
}

func packageUsage(cmd *cobra.Command) error {
	fmt.Fprintf(cmd.OutOrStderr(), "Usage:\n  %s\n", cmd.UseLine())
	if cmd.HasAvailableSubCommands() {
		fmt.Fprintln(cmd.OutOrStderr(), "\nCommands:")
		for _, child := range cmd.Commands() {
			if !child.Hidden {
				fmt.Fprintf(cmd.OutOrStderr(), "  %-12s %s\n", child.Name(), child.Short)
			}
		}
	}
	if cmd.Example != "" {
		fmt.Fprintf(cmd.OutOrStderr(), "\nExamples:\n%s\n", cmd.Example)
	}
	fmt.Fprintf(cmd.OutOrStderr(), "\nFlags:\n%s", cmd.LocalFlags().FlagUsages())
	flags := pflag.NewFlagSet("applicable", pflag.ContinueOnError)
	path := strings.TrimPrefix(cmd.CommandPath(), "woobe ")
	cmd.InheritedFlags().VisitAll(func(flag *pflag.Flag) {
		if packageVisibleFlag(path, flag.Name) {
			flags.AddFlag(flag)
		}
	})
	if flags.HasFlags() {
		fmt.Fprintf(cmd.OutOrStderr(), "\nGlobal flags applicable to this command:\n%s", flags.FlagUsages())
	}
	return nil
}
