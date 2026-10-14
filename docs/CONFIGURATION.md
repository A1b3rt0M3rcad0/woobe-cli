# Project configuration

Check the exact target before running it:

```sh
woobe project configuration check agent AGENT_UUID --env staging
woobe project configuration check network NETWORK_UUID --env production
woobe project configuration check provider CREDENTIAL_UUID --env draft
```

Use `--version-id UUID` for an exact retained Agent/Network snapshot. The backend
selects current snapshots otherwise. The command only reads configuration
metadata: it does not test external credentials or start a Run. `blocked` exits
with 6, `unknown` or an unsupported server exits with 9. Unresolved legacy
dependencies must be reviewed; a missing readiness capability is never treated
as success. Inspect `--output json` for the complete dependency diagnostics.

Project configuration separates a key declaration from its environment values.
Development agents can declare and bind keys without handling provider tokens.
An operator with value-write permission configures them in **Project →
Environments**. Saving configuration does not execute, publish or deploy.

These commands require a backend advertising `project_configuration.v1`.
An older server produces a compatibility error before a declaration is sent.
There is no fallback to local environment variables, private files or a local
secret store. Connection authentication continues to use the saved CLI Key.

```sh
woobe project configuration declare PRIMARY_API_KEY --kind secret --format token
woobe project configuration list --env draft
woobe project configuration list --env staging --status missing --search PRIMARY
woobe project configuration usage KEY_UUID
```

`KEY_UUID` is the declaration ID returned by `declare` or `list`. It is not the
secret value. The default Environment is Draft. List results are bounded to 50
keys; use `--limit 1..200` and `--offset` for subsequent pages. Usage has its own
50-consumer pages and `--offset`. Counts describe the selected Project and filters.
An absent local `.woobe-config` does not prevent these remote operations.

Declarations support `secret` or `value` kinds and `string`, `token`, `json`,
`pem` or `url` formats. A declaration sends only its name, kind, format and
optional description. Never put material in the description or a command flag.
These commands reject `--file`; operator value entry belongs in Environments.

The metadata response distinguishes missing, configured, invalid, revoked and
unavailable material. A configured key alone does not prove a target is ready:
its binding, destination approval, selected Environment and exact dependencies
also matter. Metadata reads do not call the Provider or test a token externally.

For a reference-capable backend, a Provider request uses a typed binding:

```yaml
project_id: PROJECT_UUID
provider: openai
name: Primary provider
authentication:
  type: api_key
  fields:
    api_key:
      source: project_environment
      key: PRIMARY_API_KEY
```

Submit it with `woobe project provider-credential create --file provider.yaml`.
`PROJECT_UUID` is the selected Project's real UUID; the key must already be
declared in that Project. The request contains no Provider API key. Creating a
Provider requires both Provider-write and configuration-bind authority. The
server controls environment selection and refuses inline values mixed with a
binding. Grants for one Environment do not imply access to another.

The configuration commands display metadata in text, compact, table or JSON.
`--wide` does not reveal values. `--dry-run` previews the metadata operation
without contacting the server, so it cannot prove compatibility or authorization.

## Bind and test an existing connection

```sh
woobe project configuration bind provider CREDENTIAL_UUID --key PRIMARY_API_KEY
woobe project configuration bind tool TOOL_UUID --type basic --field username=TOOL_USER --field password=TOOL_PASSWORD
woobe project configuration test provider CREDENTIAL_UUID --env staging
woobe project configuration test tool MCP_TOOL_UUID --env draft
```

`bind` reads the current binding revision and submits a compare-and-swap update.
If another editor changes the resource first, the server rejects the write. Keys
must already exist. `--key` is shorthand for an `api_key` reference; Tool bindings
also accept `bearer` (`token`), `basic` (`username`, `password`), `headers`
(repeated header-name references), or explicit `none`. A reference is never a
secret value. Legacy inline credentials must follow the server's migration flow;
binding must not silently discard their current material.

`check` only inspects metadata. **`test` contacts the external system** using
server-resolved values and approved destinations for the selected Environment.
A Provider test discovers available models. An MCP Tool test performs discovery
and saves its result only if the Tool revision and operator authority still
match. It does not execute HTTP Tool actions. Values never return to the CLI.
An inconclusive result exits with 9; a connection that is not ready exits with 6.
Each command requires its advertised capability, so an unsupported server is
rejected before a binding change or external test is requested.

Portable structure and explicit remote reuse are described in [Packages](PACKAGES.md#remote-configuration-packages-schema-20). Structural import does not imply execution readiness.


Knowledge installed structurally retains a `PENDING_PREPARATION` vector snapshot.
After the operator fills the referenced values and approves their use, read its
`build_config_hash` with `woobe project knowledge snapshot get SNAPSHOT_UUID`.
Create a YAML file containing `project_id`, `environment` and
`expected_build_config_hash`, then queue it explicitly:

```sh
woobe project knowledge snapshot prepare SNAPSHOT_UUID --file prepare.yaml
woobe project knowledge snapshot get SNAPSHOT_UUID
```

The queue request checks the exact build intent and configuration before accepting
work. Repeating the same accepted request does not create another build. `BUILDING`
is pending; wait for `READY` before Candidate/release preparation. A failed or
changed intent requires a new snapshot version. The CLI does not supply values.

## Legacy source migration

Migration is a separate administrative operation. Inspect the server's
`migration_source_kinds` capability before choosing a source family. The current
bridge supports Provider credentials and recognized API-key/custom Secret header
templates; OAuth, mixed authentication and route-specific forms may require review.
A completed batch does not mean that every source family in the Project migrated.

The existing generic request command can read the report or submit an explicitly
reviewed migration plan. Replace `PROJECT_UUID` and `MIGRATION_UUID` in these
examples; `migration-plan.yaml` contains `source_kind: secret` and `limit: 100`:

```sh
woobe request GET /core/projects/PROJECT_UUID/configuration/capabilities
woobe request POST /core/projects/PROJECT_UUID/configuration/migrations/preview --file migration-plan.yaml
woobe request PUT /core/projects/PROJECT_UUID/configuration/migrations/MIGRATION_UUID --file migration-plan.yaml
woobe request GET /core/projects/PROJECT_UUID/configuration/migrations/MIGRATION_UUID
```

Preview is read-only; PUT records the chosen inventory page. Cutover additionally
requires the server's reviewed bridge-process inventory, current job revision and
explicit drain attestation. Follow the server migration runbook for advancing
batches. Do not invent process identities or bypass `CONFIG_MIGRATION_WORK_PENDING`:
active/uncertain Runs, vector builds, pending operations and unconsumed plans must
be resolved first.

Migrated header templates retain their author snapshots. Values resolve on the
server for captured consumers/destinations; a new Tool does not inherit authority
by copying a Secret name. Rotation preserves an accepted Run's value version.
New Secrets created through the compatibility API are centralized on their first
commit, with the existing Secret UUID preserved. Adding a legacy Secret name to
a new Tool requires `secret:read` and `project:environment:use-policy:write` for
the compatibility routes in Draft, Staging and Production. A rejected approval
rolls back the Tool write; editing an existing Tool never restores a revoked
policy. Prefer explicit references and approvals for the intended Environment.
MCP creation with Secret templates defers discovery; use the explicit connection
test/preparation after reviewing its destination approval.
Existing persisted-Tool execution and MCP refresh commands also use the central
resolver on a supporting backend. Their legacy contracts select Draft; choose the
configuration preparation command for another Environment. A revoked use grant
is an authorization failure (exit 4), while a blocked readiness `check` exits 6.
After separating Environments, an ambiguous old writer receives
`LEGACY_ENVIRONMENT_REQUIRED`. Configure the intended Environment centrally.
Schema 2.0 exports contain declarations/references, never local alias authority or
material. The CLI must not recreate a secret locally to bypass these diagnostics.

An explicit package binding to the **same existing Provider UUID** can reuse its
central legacy route. Planning records `binding_source: legacy_configuration`;
the stored authentication, captured consumers and environment mapping are
preserved. A direct reference to that key does not inherit legacy authority.
Changing the connection's binding revision invalidates a captured deployment.
A normal value rotation instead requires refreshing the activation plan.

Authenticated HTTP Tools may use parameter placeholders in the path or query,
with a fixed origin. A path such as `/v1/orders/{{order_id}}` requires approval
for the static `/v1/orders/` prefix. Every rendered request is checked again;
placeholders in the origin and requests outside the approved destination fail.

## Project overview

```sh
woobe project configuration summary --env draft
woobe project configuration summary --env production
```

This metadata-only command reports distinct keys, connections and affected
Agent/Network targets for the selected Environment's current selections. It does
not inspect historical Runs or test credentials. Targets without a deployed
selection are counted separately. The backend bounds the inventory to 50 Agents
and 50 Networks; a larger or unresolved observation exits 9 with `unknown` and
`complete: false`. A complete blocked observation exits 6. Use `check` for an
exact target/Candidate before execution. An older server without the optional
`project_summary` capability is rejected before the summary request is sent.
