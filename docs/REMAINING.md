# Remaining full-plan work after 67 commits

All further CLI implementation stays in `feat/go-control-plane-cli`, PR #1. Commit count is not an acceptance criterion. This is an implementation checkpoint, not completion or production certification.

| Priority | Work | Boundary / evidence needed |
| --- | --- | --- |
| 1 | Reconcile canonical Control Keys, authority categories, grants/revisions/delegation, capabilities and effective-authority contracts | Backend `woobe`; CLI has 5 branch-dependent and 12 proposed operations, not proof they are integrated. OpenAPI advertisement only proves route discovery. |
| 1 | Verify tenant isolation, human viewer, delegation ceilings, target/environment restrictions and production denial | Real backend authorization matrix with explicit editor/publisher/runtime principals. Fixtures cannot certify server rules. |
| 1 | Execute full Agent/Network creation/configuration/release/test/publication/runtime workflow against real Woobe | Deployed test instance and valid scoped credentials required. Client fixture covers configuration and explicit publication selection, not actual server DTOs/releases/runtime. |
| 2 | Complete semantic resource-kind manifest format and reconciliation | Current schema is explicit steps. Need resource identity/existence checks, complete authorized projections, revision-aware updates, kind-specific dependency rules, import/export round trips and omission-safe completeness tracking. |
| 2 | Automatic uncertain-write attribution/idempotency | Explicit-ID authorized reads verify supplied fields, not original-write attribution. Server operation IDs/idempotency or equivalent authoritative evidence required. |
| 2 | Complete schema semantics and write-time validation | Current bounded subset rejects unsupported rules. Need supported OpenAPI/JSON Schema dialect strategy, format/pattern compatibility, full DTO coverage, validation for route/query inputs and create/PATCH differences. |
| 2 | Complete pagination and stream recovery | Need verified endpoint pagination/completeness contracts, bounded `--all`, reconnect/event replay and gap scenarios. No guessed cursor protocol or automatic write retry. |
| 2 | Native credential providers and sessions | OS keychains, Windows protected human-session storage, real login/refresh/logout scenarios and authoritative human-principal checkpoint binding. POSIX fallback and environment credentials exist. |
| 3 | Quantitative acceptance and distribution verification on actual OSes | Total statement coverage 69.6%, CLI 68.6%; the plan does not specify a numerical code-coverage threshold. Its operation/principal/E2E gates remain unmet. Config package lacks direct coverage; credentials 48.6%. Cross-build/package checks are not platform execution/keychain tests. |
| 3 | Release readiness | Complete real E2E and design gates before draft removal/tag/release. No merge/release requested or performed. |

Backend changes cannot be included in the same GitHub PR as this separate CLI repository. Remaining backend work must be delivered in its own repository and consumed here; duplicating backend policy in the CLI does not satisfy acceptance.
