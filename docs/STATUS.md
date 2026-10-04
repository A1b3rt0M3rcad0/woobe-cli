# Implementation status

Date: 2026-10-04. Branch: `feat/go-control-plane-cli`. This is the initial 30-commit implementation, **not completion of the full design**. Keep all follow-up CLI work in this branch and one PR.

## Delivered client behavior

- Go 1.22 minimum; Cobra v1.8.1; `woobe-sdk-go` v0.1.0 resolves the inspected commit `5a78817a64dc5dcb15aa1d38ec54d4c289f7f56c`.
- 188 registry-backed API operations: 171 observed in the inspected master, 5 branch-dependent Control Key routes, 12 proposed extension contracts. Additional local context, credential, upload, export, runtime SDK and manifest commands exist outside that HTTP registry.
- Workspace membership/invites; Project/access/environment keys; Agent prompts/contracts/model configs; Agent release lifecycle/tests; multiple Networks, drafts/promotions/activations/environments; HTTP/MCP Tools; Knowledge/Documents/snapshots; providers/models; Skills; ChatSurfaces/keys; Runs/traces/usage.
- `request` sends only origin-relative paths, forbids redirects and never automatically retries a mutation.
- Explicit parent identifiers: either positional IDs in route order, or flags such as `--agent`, `--release`, `--network`. Help lists route order. `agent`, `network` and `control-key` aliases forward to canonical handlers.
- JSON envelopes, documented exit codes, redacted administrative output, JSONL semantic Runtime v2 events. `table` currently means indented JSON, not tabular columns.
- Atomic local config; explicit administrative and runtime credential references; private files on POSIX; origin-scoped cookie persistence with CSRF. POSIX credential storage is a fallback, **not OS keychain integration**. Windows private-file credential import/read is intentionally unsupported; environment credentials remain available.
- Credential issuance reserves an exclusive 0600 destination before requesting the write. The complete issuance response goes into that file, ordinary output is redacted. A failed/uncertain issuance leaves the destination for reconciliation; it never silently emits a second key.
- Document upload uses a multipart pipe rather than buffering the document.
- Configuration manifests support dependency order, exact `${steps.<id>.<field>}` references, persisted redacted results and checkpoints. Publication, execution, key issuance and deletion require separate commands.
- A completed checkpoint skips committed steps. Unknown/in-flight writes block resume until externally reconciled. This is client checkpointing, not server idempotency or transactional apply.
- `help --output json` provides operation metadata. `schema` describes the local transport contract and **does not claim exact domain field validation**. `server-schema --command ...` reads the server's OpenAPI operation and components when available.
- CI definition: formatting, module verification, vet, race tests, native build and six cross-build targets. Tag workflow builds archives and checksums; no release/tag is created by this change.

## Full-plan requirements still open

| Requirement | Current evidence / remaining work |
| --- | --- |
| Control Key acceptance in current master | Design pins a divergent backend feature branch. Five client operations exist, but current-master integration is unverified. |
| Authority categories, materialized grants, revision history and delegation | Twelve proposed operations return exit 9 rather than pretending the API exists. Generic authorized `request` can consume an independently implemented contract. Locate and integrate the canonical backend implementation. |
| Effective authority, access checks and capabilities | Require server contracts. Permission annotations are historical catalog hints, not live authority or a complete per-route policy map. |
| Human viewer enforcement and delegation matrix | Server responsibility; local HTTP fixtures cannot prove tenancy, principal policy or permission ceilings. |
| Runtime key constraints and production negative tests | SDK adapter exists; real Agent/Network authorization, environment and target tests remain. |
| Full resource-kind declarative format from design | The current format is an explicit `steps` manifest. Semantic resource reconciliation, remote diff, revision-aware plan, export/import for every kind and full completeness tracking remain. `manifest diff` returns exit 9. |
| Automatic reconciliation of uncertain writes | Checkpoint refuses replay. Server lookups/idempotency and an evidence-based reconciliation command remain. |
| All discovery from one registry | HTTP commands use the registry; local, SDK, upload and export commands still need unified metadata/schema registration. |
| Canonical field-level schemas and constraints | Server-schema offers authoritative discovery when enabled; embedded full DTO schemas and local field validators remain. |
| Pagination and long streams | Manual repeatable query parameters supported; complete endpoint-specific pagination and SDK reconnect/gap scenarios remain. Reconnect is disabled rather than accepting unverified replay behavior. |
| Login/refresh/logout E2E | Cookie/CSRF persistence tested with fixtures, not a deployed Woobe installation. No automatic refresh/retry path. |
| OS credential protection | Native keychains and Windows protected session storage remain. |
| Real Woobe E2E and final coverage gates | No deployment or real instance credential supplied. Race tests/cross-builds are client evidence only. |
| Supported production toolchain | CI pins the reproducible locally verified Go 1.22.2 compiler. Move to a currently supported Go toolchain and rerun the matrix before release. |

The design includes changes in `woobe` as well as `woobe-cli`. A GitHub PR belongs to one repository; this CLI PR cannot integrate backend changes into a different repository. Do not import the Python server or duplicate authority enforcement inside the CLI to make that dependency disappear.

## Publication state

The remote repository was discovered empty. Local `master` contains only commit `0017916` (`chore: initialize CLI repository`, empty tree); the implementation branch contains 30 subsequent commits. The attempted initial default-branch push was rejected by automatic approval review because the requested workflow was a feature branch and PR. The feature-branch push separately lacked authenticated Git transport. No PR, merge, tag or release has been created.

Once the minimal remote base is explicitly authorized, initialize it, publish this same branch through the authenticated GitHub connector and open **one draft PR** using `docs/PR.md`. Keep it draft until the remaining full-plan acceptance conditions are met.
