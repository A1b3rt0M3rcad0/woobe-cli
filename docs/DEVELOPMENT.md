
### Bind a Provider before the first Agent push

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
