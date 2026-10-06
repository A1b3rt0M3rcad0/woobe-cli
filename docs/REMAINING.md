# Remaining full-plan work

Checkpoint: 2026-10-05. Backend implementation is delivered in [Woobe PR #177](https://github.com/A1b3rt0M3rcad0/woobe/pull/177), with seven workflows passed at `acb71e17`. That PR remains draft and unmerged for user approval. CLI PRs #1–#7 are merged; the body-pagination continuation is published separately in the client repository and consumed by an immutable backend CI pin.

| Priority | Work | Completion criterion |
| --- | --- | --- |
| 1 | Review and integrate backend PR #177 | User approval and merge; retain exact CLI/server SHA and CI evidence. No new backend PR is necessary. |
| 2 | Native credential providers and protected sessions | Native store/read/rotate/delete tests, Windows protected session support and live login/refresh/logout. Preserve explicit administrative/runtime credential choice. |
| 2 | Complete declarative reconciliation | Authorized observe/plan/apply with identity/existence, dependencies, omission-safe projections, revision checks and authoritative uncertain-write recovery; convergent reapply must not duplicate creation. |
| 2 | Complete input validation | Extend the delivered OpenAPI 3.1 subset, UUID/date/date-time and request direction to exact DTO coverage, arbitrary directional compositions, path/query schemas and create/PATCH differences. Unsupported rules must remain explicit. |
| 2 | Stream recovery | Reconnect by event cursor, replay deduplication, gap detection and terminal evidence without reissuing the execution. Current reconnect remains disabled. |
| 2 | Semantic export/import and secret references | Complete authorized resource projections and round trips; category definition import is delivered without implicit grants; applied-grant migration remains; protected provider references must not embed secrets. |
| 3 | Broader compatibility and phase gates | Audit each operation/principal and full phase criterion; qualify supported client/server versions and native protected providers. A deterministic model fixture is not an external-provider qualification. |
| 3 | Release | User-reviewed integration, clean installation by tagged binaries/go install, capabilities/changelog and final acceptance. Development artifacts are not a tagged release. |

Reviewed body pagination is implemented for category list/history/audit and Agent sessions. It supports bounded `--all`, explicit continuation, retained filters, cycles/invalid markers, partial evidence and collection-completeness metadata. Generic Link traversal remains available but does not certify a collection; endpoint additions need reviewed contracts. Semantic export is still incomplete, so roadmap item 6.10 remains partial.

The current [audit](COMPLETENESS.md) replaces historical backend-unavailable assessments. Commit count does not establish completion. Separate repositories require separate commits; backend PR #177 can validate an immutable client branch SHA without another PR.
