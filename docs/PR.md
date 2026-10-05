# Title

feat: extend CLI contexts and bounded contract validation with migration evidence

# Description

The CLI preserves explicit context selection and prevents a Workspace override from inheriting the previous Project. Configuration loading is strict and bounded; contexts support clearing fields and independent administrative/runtime credential attachments.

This continuation adds 30 commits on an ordered PR stack: 12 context commits in the existing feature PR, 16 schema/test commits in the validation PR, and 2 migration/final-acceptance commits in the last PR. Across the stack there are 113 commits above the initial master, 237 runnable commands and 188 HTTP operations.

Advertised request-body validation now supports additional object/dependency/conditional/array assertions, exact numeric equality, local pointers and bounded work. Inactive unsupported schemas fail closed. Validation makes only an OpenAPI GET and does not authorize or execute the selected write. Full schema dialect/DTO semantics and automatic write-time validation remain open.

Local validation: 122 test functions plus subtests, go test -race, vet, module verification, native build and 76.0% total statement coverage. Bounded schema fuzzing is included in CI. Final stack CI and distribution checks are reported on the PR head.

Draft acceptance remains open: canonical backend authority integration, real Woobe E2E, semantic resource coverage and native credential providers remain. No merge or release is claimed. Audited full-plan completeness is 39/101 delivered (38.6%), 33 partial and 29 pending; no real-server phase acceptance is verified. Python-to-Go migration is documented against the pinned inspected prototype.
