# Title

feat: implement Go CLI with complete command discovery and recoverable manifests

# Description

Introduces the Go `woobe` client for workspace/project administration, Agent/Network configuration and releases, HTTP/MCP Tools, Knowledge upload, provider/model/Skill/ChatSurface management, key lifecycle and runtime SDK execution.

37 commits on one feature branch. Unified discovery exposes 226 executable handlers including 188 HTTP operations, flags and input/output descriptors. Manifests support dependency references, comparison of supplied update fields, exclusive checkpoints and explicit-ID reconciliation after uncertain writes without replaying creation. Writes never automatically retry; PATCH omission/null is preserved; issued secrets require exclusive private files. Runtime dry-run makes no network calls. Table rendering preserves JSON contracts. Proposed endpoints require OpenAPI route/method advertisement before execution.

**Draft: full design acceptance remains open.** Canonical backend Control Key/authority/delegation integration, complete semantic resource-kind reconciliation, domain schemas, native keychains and real Woobe E2E remain. See `docs/STATUS.md`; all further CLI work stays in this branch/PR.

Validation: Go 1.27.1 module verification, vet, race tests (36 test functions, 64.6% total coverage), native build and six Linux/macOS/Windows amd64/arm64 packages with verified checksums. Initial 30-commit PR CI passed all seven jobs; follow-up head checks are attached to the PR. No real Woobe E2E, merge or release is claimed.
