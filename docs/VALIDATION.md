# Validation — 2026-10-04

Go 1.27.1, Linux amd64. Module verification, formatting, diff checks, vet, race tests and native build passed. 89 test functions plus subtests. Statement coverage: 71.2% total, 69.9% CLI, 85.9% transport, 89.2% strict JSON, 82.9% manifest, 64.5% bounded schema validator.

Discovery exposes 234 executable handlers, including 188 HTTP operations. Six Linux/macOS/Windows amd64/arm64 packages built with verified SHA256, exact archive contents and embedded manifest schemas. Version-scoped `artifacts.json` is verified against archive bytes. The previous 37-commit CI passed all seven jobs (run 37243940695); final continuation CI is attached to the same PR.

The 30-commit continuation adds tests for duplicate/ambiguous input, secret-free manifests, canonical hashes, dependency IDs, large-number preservation, checkpoint parsing and key binding, malformed-success uncertainty, HEAD/read headers, network-free dry-runs, versioned schema discovery, authoritative body validation/local references/composition, partial projections and diagnostics. A composed client fixture creates Agent/prompt/contract/model/Network/draft, resumes without replay, rejects changed recovery credentials and distinguishes denied editor publication from explicitly selected publisher success.

Real Woobe E2E was not run: no deployed test instance/credentials. Fixtures prove client behavior, not backend authority/tenancy or real DTO compatibility. Code coverage is measured evidence, not a numerical acceptance gate defined by the plan. Operation coverage, authorization matrices, real E2E and native credential protection/OS execution remain open. See `REMAINING.md`.

Resource continuation: all six kind/action mappings are checked against executable HTTP handlers. Tests cover v2 compile validation, large integers, typed IDs, unsupported actions, exact dependency references, compile-to-v1 execution identity, authorized diff, unchanged observation/ETag requirements, checkpoint resume and uncertain creation reconciliation without replay. Both schema files are verified inside all six packages. Final continuation CI is attached to PR #1.

Latest continuation tests selective capture/reapply with only requested fields, omission/null/revision preservation, unavailable or secret fields, mismatched identity/scope, new typed Skill/Project intents, same-route advertised pagination, cycle/scope/foreign-link rejection, partial redaction, dry-run preflight, build identity and composed-step progress reporting. COMPLETENESS.json covers all 101 original roadmap deliveries exactly; its report is regenerated and checked in CI.


## Portable Package qualification in draft PR #10

The implementation uses the exact Package schema catalog from Woobe draft PR
#179. Current local checks cover transport/inventory proof, bounded archives,
private checkpoints, lookup after lost Apply responses, public alias masking,
ordered JSONL completion, backend-compatible Unicode binding digests and race
detection. Windows protection is implemented with native ACL/lock APIs and is
cross-compiled locally; native execution is a separate required CI result.
The workflow runs native Package tests on Linux/macOS/Windows for amd64/arm64,
plus parser/path fuzzing. Cross-compilation alone is not native qualification.
