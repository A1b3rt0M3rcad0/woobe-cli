# Implementation status

Updated 2026-10-05 (America/Sao_Paulo). The canonical Go client is in `A1b3rt0M3rcad0/woobe-cli`; implementation PRs #1–#7 are merged at `a9c6ace053ea1f3e57abbb5d372f77010211ad56`. This continuation is on `feat/cli-pagination-contracts`; the code revision `659baea23458380b84bc7059d2a544c9c06fc1f5` passed all six [CLI CI jobs](https://github.com/A1b3rt0M3rcad0/woobe-cli/actions/runs/37385398090). No new PR, merge, tag or release is claimed.

## Backend acceptance

[Woobe PR #177](https://github.com/A1b3rt0M3rcad0/woobe/pull/177) at `77c53832f0e5b35488d1574b3cf62777486f5189` is **draft and unmerged**, awaiting user approval. Seven backend workflows passed on that SHA. Required tests compile the pinned Go client and run API, workers, PostgreSQL, Redis, RabbitMQ and MongoDB with a deterministic model provider.

Verified scenarios include Project/provider/Agent/Network configuration, staging and production publication, Agent and Network rollback, runtime execution with separate keys, category lifecycle, revision-fixed grants, stale ETags, constrained authority, administrative replay, key rotation and old-key refusal. This evidence replaces the older claim that backend acceptance was entirely unavailable. It does not imply approval, merged availability or complete acceptance of every roadmap phase.

The pagination/backend continuation is `8ab9d13953c1cd77462af85a3a3a957395141edc`. Its required live CLI tests passed, including all four body-paginated routes and a session reader restricted to the Agent and staging. Catalog revision `2026-10-05.2` includes the session GET under `run:read`; the authority guard evaluates the actual environment query and the default draft filter. Exact results of the remaining workflow groups are maintained in PR #177.

## Delivered client behavior

- Go 1.22 language minimum, Go 1.27.1 toolchain and pinned runtime SDK; 240 executable handlers including 190 HTTP operations.
- Explicit Workspace/Project contexts and separate administrative/runtime credential selection; cookie/CSRF sessions and POSIX credential fallback. Native protected providers and Windows protected sessions remain open.
- Registry-derived help, flags and invocation schemas; server OpenAPI discovery and bounded request-body validation, read-only manifest preflight and opt-in validation before writes.
- Workspace/Project administration; Agent/Network configuration and release lifecycle; Tools HTTP/MCP, Knowledge, providers/models, Skills, ChatSurfaces, key lifecycle, runtime and diagnostics.
- No automatic mutation retry; strong preconditions and idempotency headers are forwarded where supported. Secret issuance requires an exclusive private destination; ordinary outputs are redacted.
- Versioned manifests, typed dependencies, selective capture, exact numeric state comparisons, locked checkpoints, partial reports and explicit-ID reconciliation. A GET observation does not attribute the original uncertain write.
- Reviewed `--all` pagination for category list/history, authority audit and Agent sessions. Three body protocols preserve all filters; bounds cover pages, aggregate bytes and total traversal deadline. Single-page and partial output report collection completeness explicitly. See [PAGINATION.md](PAGINATION.md).
- Six packaged OS/architecture targets; CI executes native packages on Linux/macOS/Windows. Native smoke now exercises body pagination and partial-collection metadata against loopback fixtures; live backend acceptance is separate.

## Remaining work

Native credential providers and sessions, full declarative observe/plan/apply reconciliation, complete input schema/dialect/route/query coverage, stream reconnect/replay/gaps, complete semantic export/import and protected manifest secret references remain. Pagination is implemented for all four reviewed body-paginated administrative routes, but full semantic export remains partial and unknown endpoint protocols are not inferred.

The [101-item audit](COMPLETENESS.md) now includes the separately delivered backend evidence. Its percentage measures delivered roadmap items, not effort or production readiness. No full phase acceptance is asserted merely from individual tests.

## Request validation continuation — 2026-10-05

Shared request-body validation now enforces readOnly/writeOnly request direction, validates UUID/date/date-time formats and refuses unsupported advertised OpenAPI/JSON Schema dialects. Regressions exercise validate-input, opt-in HTTP writes, manifest preflight/apply and checkpoint evidence. Native package smoke also checks valid/invalid UUID and timestamps, read-only refusal and write counts. Path/query and full DTO/dialect/composition coverage remain partial; no roadmap item is promoted solely for this narrower delivery.

## Path/query and DTO continuation — 2026-10-05

Opt-in parameter validation now covers canonical HTTP reads/writes, standalone
read-only validate-input and manifest preflight/apply. It enforces advertised
required fields, primitive formats/types, nullable/composed types, scalar
multiplicity and form arrays; preserves context/route aliases; refuses unknown
parameters and unsupported styles. Inline DTO `$defs` are inspected with valid
document-relative references. Native package smoke checks preflight and validated
pagination, and adversarial parameter fuzzing ran 33,051 cases locally.

The backend continuation adds advertised DTO body contracts to Agent,
provider/model, prompt/contract, model configuration/release-test, Knowledge,
Tools, session/message/job and identity operations, plus consumed Agent-session
query filters. The live acceptance workflow now requires this client validation.
Final remote backend validation is tracked in PR #177; local unit/schema checks
are not a substitute for the required live gate. Full input validation remains
partial: arbitrary dialects/compositions, object/content parameter serialization,
header/cookie coverage and DTO/controller custom validation remain unqualified.
The full-plan delivery count remains 62/101 (61.4%).
