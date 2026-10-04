# Validation — 2026-10-04

Go 1.27.1, Linux amd64. Module verification, formatting, diff checks, vet, race tests and native build passed. 36 test functions plus subtests; statement coverage 64.6% total, 63.2% CLI, 81.4% transport. Coverage remains below full design acceptance.

Unified native discovery exposes 226 entries. Six Linux/macOS/Windows amd64/arm64 packages built; all archive SHA256 checks passed. Initial 30-commit PR CI passed all seven jobs (run 37242768678); follow-up head checks are available on PR #1.

New tests cover uncertain creation recovery by explicit-ID authorized GET, resumed dependent creation without replay, field-only remote diff, checkpoint locking, extension advertisement and authorization, complete executable discovery/schema, table redaction, SDK error mapping and zero network calls for runtime dry-run. Earlier transport, credential, upload, manifest and PATCH tests remain.

Real Woobe E2E was not run: no deployed test instance/credentials. Fixtures prove client behavior, not backend authority/tenancy enforcement. Cross-builds prove compilation, not native keychain or server compatibility. See `STATUS.md` for remaining full-plan requirements.
