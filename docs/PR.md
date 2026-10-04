# Title

feat: implement Go control-plane CLI foundation and resource commands

# Description

Introduces the external `woobe` executable in Go with workspace/project control, Agent and Network configuration and releases, Tools HTTP/MCP, Knowledge upload, provider/model/Skill/ChatSurface management, key lifecycle, explicit runtime SDK credentials, JSON discovery and deterministic configuration manifests with checkpoints.

The implementation consists of 30 commits on one feature branch. Its registry contains 188 HTTP operations; local and SDK commands are additional. PATCH input preserves omitted/null fields, issued secrets require exclusive private output files, administrative writes are never automatically retried, and unknown manifest writes block resume.

**Draft: this is the initial implementation, not full design acceptance.** Authority categories/introspection require canonical backend contracts; Control Keys depend on an unintegrated backend branch; semantic resource diff/reconciliation and complete field schemas remain. `docs/STATUS.md` tracks each incomplete requirement and separates client fixtures from real authorization evidence. Keep further work in this same branch/PR.

Validation: local `go vet`, race tests, module verification, native build and packaging for Linux/macOS/Windows on amd64/arm64. No real Woobe E2E, GitHub CI result, release or merge is claimed.
