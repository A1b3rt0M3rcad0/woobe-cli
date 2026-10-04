# Validation — 2026-10-04

Go 1.27.1, Linux amd64. Module verification, formatting, diff checks, vet, race tests and native build passed. 66 test functions plus subtests. Statement coverage: 69.6% total, 68.6% CLI, 86.2% transport, 89.2% strict JSON, 86.0% manifest, 64.5% bounded schema validator.

Discovery exposes 230 executable handlers, including 188 HTTP operations. Six Linux/macOS/Windows amd64/arm64 packages built with verified SHA256, exact archive contents and embedded manifest schemas. Version-scoped `artifacts.json` is verified against archive bytes. The previous 37-commit CI passed all seven jobs (run 37243940695); final continuation CI is attached to the same PR.

The 30-commit continuation adds tests for duplicate/ambiguous input, secret-free manifests, canonical hashes, dependency IDs, large-number preservation, checkpoint parsing and key binding, malformed-success uncertainty, HEAD/read headers, network-free dry-runs, versioned schema discovery, authoritative body validation/local references/composition, partial projections and diagnostics. A composed client fixture creates Agent/prompt/contract/model/Network/draft, resumes without replay, rejects changed recovery credentials and distinguishes denied editor publication from explicitly selected publisher success.

Real Woobe E2E was not run: no deployed test instance/credentials. Fixtures prove client behavior, not backend authority/tenancy or real DTO compatibility. Code coverage is measured evidence, not a numerical acceptance gate defined by the plan. Operation coverage, authorization matrices, real E2E and native credential protection/OS execution remain open. See `REMAINING.md`.
