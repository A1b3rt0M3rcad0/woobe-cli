# Woobe CLI

External Go client for the Woobe public control plane and runtime. Server-side policy is authoritative.

```sh
go build -o bin/woobe ./cmd/woobe
bin/woobe help --output json
bin/woobe schema --command "project agent create"
```

Unified discovery exposes 237 executable handlers including 188 HTTP operations, explicit context and credential selection, Agent/Network configuration and release operations, upload, remote manifest comparison, explicit-ID checkpoint reconciliation, resource-kind manifests, selective field capture and bounded advertised pagination, authoritative input validation and partial resource export.

**The full design is not complete.** See [implementation status](docs/STATUS.md), [usage](docs/USAGE.md), [operation catalog](docs/OPERATIONS.md) and the [remaining requirements](docs/REMAINING.md) and [canonical design](docs/PLAN.md). Proposed server capabilities require OpenAPI advertisement and return unsupported errors when absent. No permissions are granted by the client.

Validation: `make check`. Distribution: `bash scripts/package.sh VERSION`.

Full-plan delivery completeness: **37.6% (38/101 deliveries)**; partial items receive no credit. See [audited assessment](docs/COMPLETENESS.md) for evidence and real-server acceptance limits.

Migration from the inspected Python prototype: [docs/MIGRATION.md](docs/MIGRATION.md). Context integrity and bounded schema validation are documented in [docs/USAGE.md](docs/USAGE.md).
