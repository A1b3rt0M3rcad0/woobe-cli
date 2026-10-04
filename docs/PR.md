# Title

feat: implement Go CLI with complete command discovery and recoverable manifests

# Description

Introduces the Go `woobe` client for workspace/project administration, Agent/Network configuration and releases, HTTP/MCP Tools, Knowledge upload, provider/model/Skill/ChatSurface management, key lifecycle and runtime SDK execution.

67 commits on one feature branch. Unified discovery exposes 230 executable handlers including 188 HTTP operations, flags and input/output descriptors. Manifests support dependency references, comparison of supplied update fields, exclusive checkpoints and explicit-ID reconciliation after uncertain writes without replaying creation. Writes never automatically retry; PATCH omission/null is preserved; issued secrets require exclusive private files. Runtime dry-run makes no network calls. Table rendering preserves JSON contracts. Proposed endpoints require OpenAPI route/method advertisement before execution.

**Draft: full design acceptance remains open.** Canonical backend Control Key/authority/delegation integration, complete semantic resource-kind reconciliation, full schema semantics, native keychains and real Woobe E2E remain. See `docs/STATUS.md`; all further CLI work stays in this branch/PR.

Validation: Go 1.27.1 module verification, vet, race tests (66 test functions, 69.6% total coverage), native build and six Linux/macOS/Windows amd64/arm64 packages with verified checksums. The previous 37-commit CI passed all seven jobs; final continuation checks are attached to the head. No real Woobe E2E, merge or release is claimed.

The 30-commit continuation adds strict JSON/manifests, key-bound recovery and offline status, exact numeric preservation, bounded authoritative body validation, generic partial exports, independent doctor checks and verified schema-bearing artifacts. `docs/REMAINING.md` records every remaining acceptance category.
