
### Bind a Provider before the first Agent push

# Managed development

For coding assistants, use `woobe skill install --agent codex` or `claude`;
the [Woobe CLI skill](AGENT_SKILL.md) is embedded and installs offline.
It follows these same development flows; it does not upload a runtime Skill or
replace the optional registry. UUID API/runtime commands still work without it.

Local descriptors use readable YAML blocks for ordinary multiline instructions.
Strings containing tabs or line separators that YAML would normalize use quoted
escapes to preserve their exact value. Do not replace those escapes with literal
indentation unless you intend to change the instructions.

Create/register a public Provider descriptor with `spec.provider` (and optional
base URL); credentials belong to the destination, never the author YAML.

```sh
woobe provider list
woobe resources bind provider '@openai-main' CREDENTIAL_UUID
woobe agent '@support' create --yes
```

Binding verifies the credential is active, belongs to the current Project and
matches the Provider type. It is stored under private `.state`, scoped to API,
Workspace and Project. Repeating the same binding is safe; rebinding or assigning
the same native credential to a second local Provider is rejected. `--dry-run`
performs the read and validation without persisting. Providers captured by `pull`
already have this binding.

### Surface authoring and lifecycle

A Surface has one `spec.target_ref` to a registered Agent or Network. Pull or
create that target first and publish its Production Release. The CLI keeps the
Surface identity once created; configuration pushes preserve its frozen target
Release and never activate or refresh it implicitly.

```sh
woobe resources create surface chat --file surface.yaml
woobe surface '@chat' create --yes
woobe surface '@chat' diff
woobe surface '@chat' push --yes
woobe surface '@chat' pull
woobe resources push --kind surface --dry-run
woobe surface activate '@chat' --file empty.yaml --yes
woobe surface disable '@chat' --yes
```

For an existing native Surface, use
`woobe surface SURFACE_UUID pull --alias chat`; its target must already be pulled
in this connection and Project. Pull merges independent local/remote edits and
stops at conflicting fields. Push compares the captured timestamp and sends
`If-Match`; the server rechecks it while holding the native row lock. A target
change requires cloning/creating another Surface.

Surface settings follow their native lifecycle: configuration changes can affect
an active Surface. Disable it first when changes must stay offline. Agent/Network
pushes always target Draft. Surface bulk dry runs report identity/CAS validation;
the native API validates its payload on apply. They do not claim an executed
server package plan. Bulk traversal runs target roots before their Surfaces.

If a Surface create times out, its durable local intent blocks automatic retries.
Inspect `surface list`, then reconcile explicitly with
`surface UUID pull --alias chat`. No keys, session tokens or access secrets enter
Surface author files.

Author `target_ref` and other dependency references use stable registry **keys**,
shown by `resources list`; terminal references use `@alias`, UUID or path.
UUID letter case does not change identity; alias matching remains case-sensitive.
Changing an alias does not rewrite dependency identity. Provider bind also records its
native identity so a later target pull reuses that same Provider rather than
creating another definition.

### Create local artifacts with support files

`resources create` copies the descriptor and its declared Skill/Knowledge files
from the descriptor's folder. Missing files, links, paths outside that folder and
paths through `.state`, `.woobe-config` or `.git` fail before author files or the
registry change. Register existing artifacts when their files are already under
the configured root.

```sh
woobe resources create skill procedure --file ./author/skill.yaml
woobe resources clone skill '@procedure' --alias procedure-copy
woobe resources move skill '@procedure' --path procedures/skill.yaml
```

Registered descriptors can be folders or individual YAML/JSON files. Pull keeps
their paths and places declared support files beside the descriptor, including
after a move. Paths must work on all supported platforms: Windows device names,
trailing dots/spaces, alternate data streams and private metadata directories
are rejected on every OS. Automatically generated aliases avoid device names.

### Reconcile pending writes

If a managed write times out or loses its response, keep its private checkpoint
and run `woobe agent '@support' reconcile` (or the corresponding Network command).
Reconciliation observes the original operation; it does not upload or apply a
second write and still works if local YAML was edited or became invalid.
For a running operation, the returned `next_command` is
`woobe package status OPERATION_ID`; add `--wait` to observe until terminal.
Run the managed `reconcile` again afterward to accept completed local bases and
clear the pending state before another push.

### Diagnose a development connection

`context use` selects the global default. A project's `.woobe-config` can pin a
different connection; an explicit `--context NAME` takes precedence.

```sh
woobe package doctor --output json
woobe package doctor --context local --output json
woobe agent AGENT_UUID pull --alias support --context local
```

The read-only diagnostic reports the effective API/context, expected and server
catalog digests, authority-fingerprint validity and managed-development support.
An incompatible result exits with code 10 and retains the evidence. It does not
create private state or display credentials or authority fingerprints.

`catalog_status: windows_crlf` identifies a server hashing a Windows CRLF checkout
without normalization. Update Woobe's backend with the LF-normalized digest fix
and restart the API. Reinstalling the CLI or rerunning `init` cannot change that
server hash. Genuine catalog mismatches still require matching backend/CLI schemas.

`PACKAGE_EXPORT_INCOMPLETE` reports whether a ModelSpec, Provider credential
binding, matching Provider, or Project Environment field is missing. These are
classified from exact fixed owner reasons; arbitrary server messages and protected
input values remain excluded. Draft pulls require the backend fix that captures
the live Agent definition even when a saved legacy Draft has no ModelSpec.
Historical environment and explicit snapshot exports remain exact.

In PowerShell, quote comma-separated field lists:

```powershell
woobe agent get AGENT_UUID --fields "id,name,provider,model,provider_credential_id" --output compact
```

Isolated ASaC Drafts use explicit `draft open/show/push/checkpoint/reconcile`
operations. They preserve the native singleton Draft and Production while storing
a verified closed artifact. See [ASaC](ASAC.md#isolated-remote-drafts) for generation
preconditions, checkpoint integrity, capability negotiation and uncertain writes.

Local sealed revisions can restore author files with
`woobe agent '@support' checkout --revision REVISION_ID` (Network has the same
command). This requires a clean source checkpoint, preserves origin/bindings and
protects shared dependencies. See [ASaC source restoration](ASAC.md#restore-local-author-source)
for the source retention and server verification boundaries.
