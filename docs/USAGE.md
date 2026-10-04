# CLI workflows

## Build and discovery

```sh
go build -o bin/woobe ./cmd/woobe
bin/woobe help --output json
bin/woobe schema --command "project agent create" --kind input
bin/woobe server-schema --command "project agent create" --api-url https://woobe.example
bin/woobe completion bash
```

Discovery covers all executable handlers and flags; local schemas describe invocation and transport/runtime output. Use the server schema to discover the actual DTO. Domains remain validated and authorized by the server.

## Context and credentials

```sh
woobe auth credential import --name editor --stdin < editor.key
woobe context create dev --api-url http://localhost:8000
woobe context credential attach dev --credential editor
woobe context use dev
woobe context set --workspace WORKSPACE --project PROJECT
woobe project agent list
```

Precedence: explicit flag, corresponding `WOOBE_*` environment variable, selected context, documented default. `WOOBE_CONFIG` selects the config file. Scope variables are `WOOBE_WORKSPACE_ID` and `WOOBE_PROJECT_ID`. Credentials may be imported from stdin on POSIX; environment fallbacks are `WOOBE_CONTROL_KEY` and `WOOBE_RUNTIME_KEY`. Runtime never falls back to an administrative key. Secret values have no command-line flag.

Select a different Workspace and the previously selected Project is cleared by `context set`. Context references are local preferences; the server resolves and verifies resource ownership. Direct known resource IDs can be read without first listing every Project.

Human login accepts a JSON file or stdin containing the server's login fields:

```sh
woobe auth login --file - < login.json
woobe auth refresh
woobe auth status
woobe auth logout
```

## Configuration and publication

```sh
woobe project agent create --file agent.json
woobe project agent update AGENT --file patch.json --if-match ETAG
woobe project agent prompt create --agent AGENT --file prompt.json
woobe project agent contract create --agent AGENT --file contract.json
woobe project agent model-config create --agent AGENT --file model.json
woobe project network create --file network-metadata.json
woobe project network draft update NETWORK --file network-definition.json
woobe project agent release create --agent AGENT --file release.json
woobe project agent release promote RELEASE --agent AGENT --file promotion.json --yes
```

IDs are positional in route order, or explicit parent flags. PATCH bytes are preserved: omitted fields stay omitted; explicit null remains null. `--if-match` forwards a server ETag; it does not implement revision protection when the server ignores that header.

`--dry-run` renders a redacted request without sending it. `--yes` acknowledges the specified execution/publication/destructive action; it never changes permissions or selects another credential. No interactive prompt is used.

```sh
woobe project api-key create --file key.json --secret-file ./issued-key.json
woobe project knowledge document upload --collection COLLECTION --document-file lesson.pdf
woobe project agent export AGENT --destination agent-projection.json
```

Exports are explicitly incomplete authorized projections, not apply-ready patches. API issuance destinations must not exist; a failed attempt can leave an empty file. Reconcile before deciding whether another issuance is necessary.

## Runtime

```sh
woobe runtime target run tutor --target-kind agent --runtime-credential runtime --file run.json --yes
woobe runtime target stream tutor --target-kind agent --runtime-credential runtime --file run.json --yes --output jsonl --timeout 10m
woobe runtime target observe tutor RUN --target-kind agent --runtime-credential runtime --output jsonl
woobe runtime target active tutor SESSION --runtime-credential runtime
woobe runtime target cancel tutor RUN --runtime-credential runtime --yes
```

Run JSON: `{"input":"Explain this topic","options":{"SessionID":"..."}}`. Omit options to let the runtime establish a Session. Option names currently follow the pinned SDK's Go JSON field names. For Network targets select `--target-kind network` explicitly. Local interruption stops observation; it does not issue remote cancellation. Administrative Runs use the separate `runtime run` family. All five SDK handlers honor `--dry-run` without network execution.

## Manifest validation, plan and checkpoint

```sh
woobe manifest validate --file testdata/agent-network.manifest.json
woobe manifest plan --file testdata/agent-network.manifest.json
woobe manifest apply --file resources.json --yes --checkpoint ./apply-checkpoint.json
```

Document fields: `schema_version`, `workspace_id`, `project_id`, `steps`. Each step declares `id`, canonical `command`, positional `args`, JSON `body`, `depends_on`, optional `if_match`. Declare dependencies explicitly. Exact references such as `${steps.agent.id}` resolve persisted redacted results; interpolation inside larger strings is unsupported.

Match manifest Workspace/Project to the selected context. Preserve the same manifest and checkpoint to resume. Committed and reconciled steps are skipped; `unknown`/`in_flight` blocks execution. Never delete the checkpoint and blindly replay creation after a lost response.

The explicit-step format composes existing server operations. It is not the complete semantic resource reconciler specified in the design. `manifest diff --file resources.json` reads compatible update resources and compares supplied top-level fields. Creation existence and unresolved dependencies remain unevaluated.

For uncertain writes: `woobe manifest reconcile STEP --resource-id ID --file resources.json --checkpoint ./apply-checkpoint.json --yes`. The same manifest, origin, workspace, project and credential reference are required. An authorized GET must match nonempty supplied desired fields. Evidence proves observed state, not original-write attribution. Resume apply using the same checkpoint. Concurrent use is rejected; remove a stale `.lock` only after confirming no process uses it.

## Output and advertised extensions

`--output table` renders rows/fields; JSON envelopes remain stable. Proposed extension endpoints require the exact route/method in `/openapi.json`; actual requests still require server authorization.

## Exit codes

| Code | Meaning |
| --- | --- |
| 0 | Success |
| 1 | Internal client error |
| 2 | Invalid usage or input |
| 3 | Authentication |
| 4 | Authorization denied by server |
| 5 | Not found/inaccessible |
| 6 | Revision/state conflict |
| 7 | Connection/uncertain write |
| 8 | Deadline |
| 9 | Unsupported capability/protocol |
| 10 | Partial composed result |
| 11 | Remote Run failure/cancellation |
| 130 | Local interruption |

JSONL runtime stdout contains semantic SDK events; client errors go to stderr. Other JSON output is one envelope. Credentials and administrative secret fields are redacted; semantic runtime event payloads are preserved.

## Input validation and checkpoint inspection

```sh
woobe schema --command "manifest validate" --kind document
woobe validate-input --command "project agent create" --file agent.json
woobe manifest status --checkpoint ./apply-checkpoint.json
woobe doctor
```

`validate-input` sends only an OpenAPI GET; it does not submit the supplied body. It supports object/array/primitive types, required/properties/additionalProperties, sizes, enum/const, numeric bounds with exact rationals, RE2-compatible patterns, local references and composition. Unsupported assertions (including format), external references and excessive schema depth return exit 9. Domain validation and authorization still occur on the server. Validation is explicit; writes do not silently add schema-fetch requirements.

JSON inputs reject duplicate fields, more than one value and nesting beyond 128. Manifests require object configuration bodies, exclude read steps and literal sensitive fields (including credential `value`), and reject duplicate/self dependencies. IDs match `[a-zA-Z0-9_-]{1,128}`. Canonical hashing ignores object-key order/formatting and preserves array order, null and numeric values. Updated hash/fingerprint rules can reject old checkpoints; preserve them for manual reconciliation instead of deleting/replaying.

Control Key fingerprints bind checkpoint recovery to actual key material, including environment keys. Human-session principal binding still needs authoritative server identity. `doctor` runs instance, identity and OpenAPI reads independently; partial diagnostics include evidence and exit 10. Advertised routes never imply effective authority.

## Resource projections and distribution

```sh
woobe export NETWORK --command "project network get" --project PROJECT --destination network.json
woobe export TRACE --command "project trace get" --destination trace.json
woobe manifest export --resource agent --id AGENT --destination agent.json
```

Generic export accepts canonical resource `get` operations with exact route-order IDs. Manifest shortcuts cover agent, network, skill-version, knowledge collection and surface. Every projection remains `complete:false` and `apply_ready:false`, with observed ETag/request ID when available. Unsupported Tool getters are not fabricated. Secret fields are redacted before file/output rendering. Projection export does not guarantee pagination or round-trip apply.

Packages include the executable, README, usage and `manifest.schema.json`. `artifacts.json` records version, source commit, compiler, OS/architecture, byte size and SHA256 for exactly six archives. `python3 scripts/verify_artifacts.py` verifies contents, schemas and checksums before release publication. No release/tag is created by this PR.
