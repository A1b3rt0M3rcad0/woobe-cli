# Remaining full-plan work after 150 commits

The latest continuation is reviewed in an ordered CLI PR stack: foundation/context integrity, schema validation, then migration/acceptance documentation and recovery/body-preflight hardening. Commit count is not an acceptance criterion. This is an implementation checkpoint, not completion or production certification.

| Priority | Work | Boundary / evidence needed |
| --- | --- | --- |
| 1 | Reconcile canonical Control Keys, authority categories, grants/revisions/delegation, capabilities and effective-authority contracts | Backend `woobe`; CLI has 5 branch-dependent and 12 proposed operations, not proof they are integrated. OpenAPI advertisement only proves route discovery. |
| 1 | Verify tenant isolation, human viewer, delegation ceilings, target/environment restrictions and production denial | Real backend authorization matrix with explicit editor/publisher/runtime principals. Fixtures cannot certify server rules. |
| 1 | Execute full Agent/Network creation/configuration/release/test/publication/runtime workflow against real Woobe | Deployed test instance and valid scoped credentials required. Client fixture covers configuration and explicit publication selection, not actual server DTOs/releases/runtime. |
| 2 | Complete semantic resource-kind manifest format and reconciliation | Steps v1 and 13-kind resource intents v2 exist; explicit create/update compiles to recoverable steps, with opt-in unchanged-update reads. Need resource identity/existence checks, complete authorized projections, revision-aware updates, kind-specific dependency rules, import/export round trips and omission-safe completeness tracking. |
| 2 | Automatic uncertain-write attribution/idempotency | Explicit-ID authorized reads verify supplied fields, not original-write attribution. Server operation IDs/idempotency or equivalent authoritative evidence required. |
| 2 | Complete schema semantics and write-time validation | Current bounded subset rejects unsupported rules. Opt-in write-body validation, pinned snapshots and body-only manifest preflight are delivered. Need supported OpenAPI/JSON Schema dialect strategy, format/pattern compatibility, full DTO coverage, validation for route/query inputs and create/PATCH differences. |
| 2 | Complete pagination and stream recovery | Need verified endpoint pagination/completeness contracts, bounded `--all`, reconnect/event replay and gap scenarios. No guessed cursor protocol or automatic write retry. |
| 2 | Native credential providers and sessions | OS keychains, Windows protected human-session storage, real login/refresh/logout scenarios and authoritative human-principal checkpoint binding. POSIX fallback and environment credentials exist. |
| 3 | Quantitative acceptance and distribution verification on actual OSes | Total statement coverage 77.9%, CLI 74.7%; the plan does not specify a numerical code-coverage threshold. Its operation/principal/E2E gates remain unmet. Config direct coverage is 75.8%, schema 84.3%; credentials remain 48.6%. CI now executes packaged binary smoke on three native OSes; this does not verify keychain/session providers or native execution of all six targets. |
| 3 | Release readiness | Complete real E2E and design gates before draft removal/tag/release. Client PRs are merged into master; no tagged production release is delivered. |

Backend changes cannot be included in the same GitHub PR as this separate CLI repository. Remaining backend work must be delivered in its own repository and consumed here; duplicating backend policy in the CLI does not satisfy acceptance.

Selective field capture now produces update intents that preserve omission/null and observed ETags. Complete projections/semantic imports remain open. Generic Link traversal is implemented, while endpoint pagination/completeness and SDK stream reconnect still remain. Audited delivery completeness: 38.6% (39/101 completed, 33 partial, 29 pending), including separate backend work; see COMPLETENESS.md.

Recovery integrity and resumable pre-write refusals are strengthened; schema snapshot identity does not lock server schema/domain state. Exact bounded JSON number comparison is now shared by diff, unchanged-update observations and reconciliation. Endpoint completeness and authoritative original-write recovery still need acceptance.

The next implementation block is ready to start in the backend repository. Follow [BACKEND_HANDOFF.md](BACKEND_HANDOFF.md) to reconcile canonical contracts before enabling further client behavior. Client readiness to start integration is distinct from release readiness.

All six implementation PRs (#1–#6) were integrated into master on 2026-10-05. The CI distribution continuation records native smoke and remote artifacts in Actions. Merging the client implementation does not satisfy the separate server acceptance gates.
