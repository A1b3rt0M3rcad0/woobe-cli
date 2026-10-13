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
