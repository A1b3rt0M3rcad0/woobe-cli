# Portable packages

Package V1 is implemented in woobe-cli and Woobe PR #179. These commands require the matching Package capability/schema catalog;
the stable CLI's Manifest commands remain independent.

## Quick workflow

Use the selected connection/project; authenticate with `woobe auth login --cli-key`.
Names and UUIDs are accepted for Agent and Network export. Names must match exactly;
ambiguous names require a terminal choice or an explicit UUID with `--no-input`.

```powershell
woobe package export agent "Orders & Requests Agent"
woobe package export agent 01a00d4a-fa75-7629-a2c6-773c25c6e2ef
woobe package export agent AGENT_UUID --env staging
woobe package export agent AGENT_UUID --env release --version v1.0.20261007.01
woobe package export agent AGENT_UUID --env production
woobe package export network "Customer Support" --env production
```

The default is the current **draft**, never production. Staging/production resolve
current assignments inside the server's frozen capture; release resolves the exact
`--version`. An unavailable environment/version fails without fallback. Names/versions
of the package are derived from captured author metadata. Unversioned definitions use
`0.0.0+<portable-definition-hash>`. `--package-version` overrides package metadata;
`--version` selects a release with `--env release`. Legacy `--source`/`--snapshot-id`
remain compatible, including their old `--version` package metadata meaning.

`--destination` defaults to `./normalized-resource-name`, even when exporting by
UUID. Paths are safe on Linux, macOS and Windows and existing paths are never replaced.
`--name` optionally overrides package metadata and its default directory name.
Exported `woobe.yaml` and component descriptors are readable YAML, with multiline
prompts. YAML and JSON descriptors are accepted. Operational plans/checkpoints/locks
remain JSON. Editing exported files invalidates the old lock; use ordinary validation
for edited author data and `--locked` to verify an unchanged captured package.

```powershell
woobe package validate ./orders-requests-agent --locked
woobe package bindings ./orders-requests-agent --destination ./destination.yaml
# Fill every placeholder using destination identities/private references.
woobe project provider-credential list
woobe package plan ./orders-requests-agent --bindings ./destination.yaml --save-plan ./plan.json
woobe package import --plan-file ./plan.json --wait
# Or plan and apply directly:
woobe package import ./orders-requests-agent --bindings ./destination.yaml --wait
```

If no destination requirements exist, omit `--bindings`. Keep bindings and saved plans
outside the portable package directory. SOURCE accepts a directory, its `woobe.yaml`,
or a `.tar.gz`/`.tgz` archive. The generated bindings template is incomplete until filled;
credentials and secrets are never included in portable author data. The template uses
`vector_snapshot_id` for external knowledge and `protected_ref` for private values.

Imports **create** resources in draft by default, regardless of export origin.
Checkpoints are automatic under the private CLI config directory's `package-imports/`.
They are scoped to API/project/principal/input/lifecycle. Repeating the same command
reconciles the original operation, even after completion, without another Apply. Use
an explicit new `--checkpoint` for an intentional second import. Output includes the
checkpoint path for `package status --checkpoint PATH --wait`. Automatic Windows
state uses an explicit private inheritable DACL; existing directories are checked.

For focused assistance use `COMMAND --help` or `woobe help package export agent
--output compact`. `--output` selects terminal formatting; `--destination` selects
an artifact path. Use `--wide` or `--output json` for full metadata and inventory.
`woobe agent export` and root `export` still return partial JSON projections; use
`package export` for portable YAML composition.

Validate a declared closure locally, without loading config, credentials or HTTP:

```sh
woobe package validate ./support
woobe package validate ./support --locked
```

Plan with an existing credential binding:

```sh
woobe package plan ./support --project PROJECT_ID \
  --bind credential.primary-key=CREDENTIAL_UUID \
  --save-plan ./approved-plan.json
```

Use `--bindings destination.yaml` for an `ImportBindings` document. Each declared
alias must have exactly one destination binding. `--bind credential.ALIAS=UUID`
is repeatable, but duplicate aliases are errors, including aliases already
present in the bindings file. Protected references are not resolved by planning.
The approved plan receipt is private and is created exclusively: an existing
file is never overwritten, including when `--yes` is provided.

```sh
woobe package plan ./support --bindings destination.yaml --dry-run
```

Dry-run captures and validates local bytes and binding coverage. It never opens a
credential store or contacts the API. It reports server semantics/authorization
as unevaluated, and cannot save a server-approved plan or validate remote schemas.

`--validate-body` validates the actual planning JSON against the advertised
OpenAPI request-body contract before submission. `--validate-parameters` validates
the upload and plan path/query schemas. `--schema-sha256` can pin that OpenAPI
snapshot. Binary multipart upload validation remains server-owned.

Planning defaults to `--lifecycle draft`. `release` and `production` require
`--release-notes`; `production` also requires `--reason`. The total `--deadline`
defaults to five minutes and includes upload, planning and HTTP requests.

Observe already accepted operations with `package status OPERATION_ID`, optionally
`--wait --deadline 5m`. An accepted operation is not complete. `package resume
OPERATION_ID --revision REVISION` resumes its observed dependency wait;
`package cancel OPERATION_ID` requests cancellation without deleting shared
resources. Package polling only retries transient GET observations, never POST
mutations.

Apply a captured source or an approved private plan with an optional explicit checkpoint:

```sh
woobe package import ./support --project PROJECT_ID \
  --bind credential.primary-key=CREDENTIAL_UUID \
  --checkpoint ./import-checkpoint.json --wait --wait-timeout 5m
woobe package import --plan-file ./approved-plan.json --project PROJECT_ID \
  --checkpoint ./import-checkpoint.json --wait
```

Source imports also save `CHECKPOINT.plan.json` exclusively so a prepared request
can use its original approved plan after restart. The checkpoint stores identities,
revisions and state, never authorization or resolved protected values. Protected
binding aliases are read once immediately before the single Apply request. The
in-flight checkpoint is fsynced before POST. A lost response is an unknown
outcome; invoking import with the same checkpoint performs lookup/observation
without posting Apply again or rereading protected bindings. A lookup 404 does
not authorize a second Apply. Existing private operational files are not replaced
by `--yes`.

Native lifecycle phases preserve standalone Agent Production when a Network
creates Agent Releases. Qualification covers all four native lifecycles, complete
cross-project Agent/Network roundtrips, protected HTTP bindings, fresh MCP discovery
and portable Knowledge rebuilt in the destination. Runtime execution uses a separate
key after all local package directories have been deleted. The
[acceptance matrix](https://github.com/A1b3rt0M3rcad0/woobe/blob/feat/portable-package-import-export/docs/testing/PACKAGE_ACCEPTANCE.md)
maps T01–T26 to the executable gates; exact paired results are in VALIDATION.md.

Export a captured native definition with all declared support files:

```sh
woobe package export agent AGENT_UUID --source staging \
  --destination ./support --name support --version 1.0.0
woobe package export network NETWORK_UUID --snapshot-id SNAPSHOT_UUID \
  --destination ./network --name network --version 1.0.0
```

Use exactly one of `--source draft|staging|release|production` or `--snapshot-id`.
The server freezes the source once. Download verifies transport bytes, exact
inventory and artifact digest, then publishes the complete locked destination
without replacing any existing path. `--knowledge portable` is the default;
`--knowledge binding` exports a destination requirement and reports that the
package is not self-contained. Authentication values, OAuth state, discovered
MCP schemas and destination IDs are never portable author configuration.

An existing operation can be observed or controlled using only its checkpoint:

```sh
woobe package status --checkpoint ./import-checkpoint.json --wait
woobe package resume --checkpoint ./import-checkpoint.json --revision REVISION
woobe package cancel --checkpoint ./import-checkpoint.json
```

Scope must still match API origin, Project and principal fingerprint. An unknown
checkpoint performs lookup with its original key; it never starts another Apply.
JSONL output has ordered progress records and one final record with the last
known remote state. A local timeout or interrupt does not cancel the operation.

POSIX operational files use mode 0600 and native exclusive locks. Windows uses
native owner/DACL checks and LockFileEx; its parent directory must already be
private so inheritance cannot expose newly created files. Unsupported or unsafe
protection fails before sending Apply. File data is synced before replacement;
Windows does not offer the POSIX directory-fsync guarantee. The six native CI
runners qualify the same Package parsing, filesystem and checkpoint tests.

Apply preserves central exit codes for server rejection, deadline (`8`) and local
interruption (`130`). An uncertain acceptance has domain code
`PACKAGE_OUTCOME_UNKNOWN`; it never permits an automatic second Apply. Polling
honors a positive server delay; otherwise it uses bounded 2–15 second exponential
backoff with jitter. A local timeout does not cancel the remote operation.
