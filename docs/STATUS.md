# Implementation status

Updated 2026-10-05 (America/Sao_Paulo). The canonical Go client is in `A1b3rt0M3rcad0/woobe-cli`; implementation PRs #1–#7 are merged at `a9c6ace053ea1f3e57abbb5d372f77010211ad56`. This continuation is on `feat/cli-pagination-contracts`. No new PR, merge, tag or release is claimed.

## Backend acceptance

[Woobe PR #177](https://github.com/A1b3rt0M3rcad0/woobe/pull/177) at `77c53832f0e5b35488d1574b3cf62777486f5189` is **draft and unmerged**, awaiting user approval. Seven backend workflows passed on that SHA. Required tests compile the pinned Go client and run API, workers, PostgreSQL, Redis, RabbitMQ and MongoDB with a deterministic model provider.

Verified scenarios include Project/provider/Agent/Network configuration, staging and production publication, Agent and Network rollback, runtime execution with separate keys, category lifecycle, revision-fixed grants, stale ETags, constrained authority, administrative replay, key rotation and old-key refusal. This evidence replaces the older claim that backend acceptance was entirely unavailable. It does not imply approval, merged availability or complete acceptance of every roadmap phase.

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
