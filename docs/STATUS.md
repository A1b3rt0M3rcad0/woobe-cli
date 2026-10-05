# Implementation status

Date: 2026-10-05. The implementation contains 113 commits above the initial master across an ordered PR stack: `feat/go-control-plane-cli` (foundation/context integrity), `feat/cli-schema-validation` (bounded contract validation), then `docs/cli-migration-acceptance` (migration and final validation evidence). This is **not completion of the full design**. The latest user instruction authorizes multiple PRs; no merge or release is claimed.

## Delivered client behavior

- Go 1.22 language minimum; Go 1.27.1 toolchain/CI; Cobra v1.8.1; `woobe-sdk-go` v0.1.0 resolves the inspected commit `5a78817a64dc5dcb15aa1d38ec54d4c289f7f56c`.
- 188 registry-backed API operations: 171 observed in the inspected master, 5 branch-dependent Control Key routes, 12 proposed extension contracts. Unified discovery covers all 237 executable handlers, including the 188 HTTP operations.
- Workspace membership/invites; Project/access/environment keys; Agent prompts/contracts/model configs; Agent release lifecycle/tests; multiple Networks, drafts/promotions/activations/environments; HTTP/MCP Tools; Knowledge/Documents/snapshots; providers/models; Skills; ChatSurfaces/keys; Runs/traces/usage.
- `request` sends only origin-relative paths, forbids redirects and never automatically retries a mutation.
- Explicit parent identifiers: either positional IDs in route order, or flags such as `--agent`, `--release`, `--network`. Help lists route order. `agent`, `network` and `control-key` aliases forward to canonical handlers.
- JSON envelopes, documented exit codes, redacted administrative output, JSONL semantic Runtime v2 events. `table` renders sorted tabular columns/fields with administrative secret redaction.
- Atomic local config; explicit administrative and runtime credential references; private files on POSIX; origin-scoped cookie persistence with CSRF. POSIX credential storage is a fallback, **not OS keychain integration**. Windows private-file credential import/read is intentionally unsupported; environment credentials remain available.
- Credential issuance reserves an exclusive 0600 destination before requesting the write. The complete issuance response goes into that file, ordinary output is redacted. A failed/uncertain issuance leaves the destination for reconciliation; it never silently emits a second key.
- Document upload uses a multipart pipe rather than buffering the document.
- Configuration manifests support dependency order, exact `${steps.<id>.<field>}` references, persisted redacted results and checkpoints. Publication, execution, key issuance and deletion require separate commands.
- A completed checkpoint skips committed steps. Unknown/in-flight writes block resume until explicit-ID authorized-read reconciliation; checkpoint locking prevents concurrent apply/reconcile. `manifest status` inspects progress offline; recovery binds to the selected key fingerprint. Legacy checkpoints with different canonical hashes or without a required key fingerprint are refused, not silently migrated. This is client checkpointing, not server idempotency or transactional apply.
- `help --output json` provides operation metadata. `schema` describes the local transport contract and **does not claim exact domain field validation**. `server-schema --command ...` reads the server's OpenAPI operation and components when available.
- CI definition: formatting, module verification, vet, race tests, native build and six cross-build targets. Tag workflow builds six archives with usage and manifest schema, verifies archive contents/checksums, and produces `artifacts.json`; no release/tag is created by this change.

Additional delivery: strict bounded duplicate-free JSON; secret-free manifest bodies, canonical plan hashes, reference-addressable IDs and preserved large integers; malformed successful write responses retain uncertainty and request IDs; HEAD support and mutation-only header separation; zero-network proposed dry-runs; authoritative input validation; partial identity/instance/route diagnostics.

## Full-plan requirements still open

| Requirement | Current evidence / remaining work |
| --- | --- |
| Control Key acceptance in current master | Design pins a divergent backend feature branch. Five client operations exist, but current-master integration is unverified. |
| Authority categories, materialized grants, revision history and delegation | Twelve proposed operations require exact route/method advertisement in server OpenAPI; absent contracts return exit 9 before writing. Advertisement does not grant authority. Generic authorized `request` can consume an independently implemented contract. Locate and integrate the canonical backend implementation. |
| Effective authority, access checks and capabilities | Require server contracts. Permission annotations are historical catalog hints, not live authority or a complete per-route policy map. |
| Human viewer enforcement and delegation matrix | Server responsibility; local HTTP fixtures cannot prove tenancy, principal policy or permission ceilings. |
| Runtime key constraints and production negative tests | SDK adapter exists; real Agent/Network authorization, environment and target tests remain. |
| Full resource-kind declarative format from design | Both explicit steps v1 and resource intents v2 are supported. V2 compiles thirteen kinds with explicit create/update, typed parent/target references and dependencies. Opt-in authorized reads checkpoint unchanged compatible updates. Remote diff compares supplied top-level update fields through compatible authorized GETs and checks observed ETags. Creation existence and unresolved dependencies remain unevaluated. Full semantic reconciliation, generic canonical-get projections and `manifest export` are available but explicitly partial/non-apply-ready. Semantic export/import and completeness tracking remain. |
| Automatic reconciliation of uncertain writes | Explicit-ID authorized-read reconciliation checks nonempty supplied expected fields and persists evidence. It proves observed state, not attribution of the lost write; automatic attribution and server idempotency remain. |
| All discovery from one registry | Delivered for all 237 executable handlers and their flags/input/output descriptors. Domain DTO validation remains server-derived. |
| Canonical field-level schemas and constraints | Server-schema offers authoritative discovery when enabled; `validate-input` reads the authoritative request-body schema and checks an explicit bounded subset, including local references and composition. Unsupported assertions/external references return exit 9; full DTO/schema semantics and automatic validation before every write remain. |
| Pagination and long streams | Manual repeatable query parameters supported; complete endpoint-specific pagination and SDK reconnect/gap scenarios remain. Reconnect is disabled rather than accepting unverified replay behavior. |
| Login/refresh/logout E2E | Cookie/CSRF persistence tested with fixtures, not a deployed Woobe installation. No automatic refresh/retry path. |
| OS credential protection | Native keychains and Windows protected session storage remain. |
| Real Woobe E2E and final coverage gates | No deployment or real instance credential supplied. Race tests/cross-builds are client evidence only. |
| Supported production toolchain | Go 1.27.1 pinned; local vet/race tests and six cross-builds passed. CI uses the same compiler. |

The design includes changes in `woobe` as well as `woobe-cli`. A GitHub PR belongs to one repository; this CLI PR cannot integrate backend changes into a different repository. Do not import the Python server or duplicate authority enforcement inside the CLI to make that dependency disappear.

## Publication state

One branch, `feat/go-control-plane-cli`, and one draft PR: https://github.com/A1b3rt0M3rcad0/woobe-cli/pull/1. The minimal master base was initialized after user authorization. The previous 37-commit head passed all seven CI jobs (run 37243940695). The 30-commit continuation is validated on the same PR. Follow-up work stays in this PR. No merge, tag or release. Keep draft until full-plan acceptance.

Resource continuation: schema v2, kinds/compile discovery, explicit identity and typed dependencies, unchanged-update observations and both packaged schemas. No implicit upsert or complete export/import is claimed.

Latest continuation: 13 typed resource kinds, explicit field capture/update round trips, bounded Link pagination, source revision in binaries and structured partial apply evidence. Full-plan delivery assessment is 38/101 (37.6%), with 33 partial and 30 pending items; see COMPLETENESS.md for all evidence and unverified phase acceptance.


Latest 30-commit continuation: strict bounded config, explicit context clearing/runtime attachments, Workspace override inheritance fixes and named show; expanded body-schema subset with inactive-branch inspection, exact numeric equality, recursive local pointers and bounded evaluator; read-only HTTP acceptance fixtures, fuzzing and Python-to-Go migration guide. Local validation: 122 test functions plus subtests, race/vet/module checks, 76.0% total statement coverage (CLI 71.8%, config 75.8%, schema 84.3%). Cross-build is not execution on each OS. Real backend phase acceptance remains 0/10; migration guide delivery raises the full-plan count to 39/101 (38.6%).
