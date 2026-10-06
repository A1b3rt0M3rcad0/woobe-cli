# Woobe CLI

External Go client for the Woobe public control plane and runtime. Server-side policy is authoritative.

## Install and use

The first version being prepared is **0.1.0**, independent of the backend.
After publication, install with Node.js 22+; no Go toolchain is needed:

```sh
npm install --global woobe-cli
woobe version
woobe help
```

Or run `npx --package=woobe-cli woobe version` without a global install.
Native Linux/macOS/Windows binaries and the npm tarball are distributed through
[GitHub Releases](https://github.com/A1b3rt0M3rcad0/woobe-cli/releases).
See [installation, Workspace setup and versioning](docs/INSTALLATION.md).
These installation endpoints become available when a maintainer publishes a
release; preparing the version/package does not publish them.

## Build from source

```sh
go build -o bin/woobe ./cmd/woobe
bin/woobe help --output json
bin/woobe schema --command "project agent create"
```

Unified discovery exposes 241 executable handlers including 190 HTTP operations, explicit context and credential selection, Agent/Network configuration and release operations, upload, remote manifest comparison, explicit-ID checkpoint reconciliation, resource-kind manifests, selective field capture and bounded reviewed body/Link pagination, authoritative input validation and partial resource export.

**The full design is not complete.** See [implementation status](docs/STATUS.md), [usage](docs/USAGE.md), [operation catalog](docs/OPERATIONS.md) and the [remaining requirements](docs/REMAINING.md) and [canonical design](docs/PLAN.md). Proposed server capabilities require OpenAPI advertisement and return unsupported errors when absent. No permissions are granted by the client.

Validation: `make check`. Distribution: `make package` (uses canonical `VERSION`).

Full-plan delivery completeness: **65.3% (66/101 deliveries)**; partial items receive no credit. See [audited assessment](docs/COMPLETENESS.md) for evidence and real-server acceptance limits.

Migration from the inspected Python prototype: [docs/MIGRATION.md](docs/MIGRATION.md). Context integrity and bounded schema validation are documented in [docs/USAGE.md](docs/USAGE.md).

Canonical HTTP writes and manifest apply can opt into `--validate-body`; `manifest preflight` produces read-only body validation evidence. See [docs/USAGE.md](docs/USAGE.md) for deferred validation, schema digest pins and recovery states.

CI packages six Linux/macOS/Windows amd64/arm64 targets and an npm tarball, verifies their source SHA/checksums/schemas and identical binaries, and exercises native executables plus npm/npx installation on three runner OSes. Distribution archives and verification reports are stored in Actions artifacts; these development builds are not tagged releases. See [CI contract](docs/CI.md).

Backend authority and live Agent/Network acceptance are delivered in [Woobe PR #177](https://github.com/A1b3rt0M3rcad0/woobe/pull/177), which remains draft for user approval. Reviewed pagination contracts and explicit collection evidence are documented in [PAGINATION.md](docs/PAGINATION.md).

Workspace category definition manifests, selective capture and explicit immutable revision diffs are documented in [USAGE.md](docs/USAGE.md). Existing grants remain revision-pinned.
