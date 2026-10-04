# Woobe CLI

External Go client for the Woobe public control plane and runtime. Server-side policy is authoritative.

```sh
go build -o bin/woobe ./cmd/woobe
bin/woobe help --output json
bin/woobe schema --command "project agent create"
```

Unified discovery exposes 226 executable handlers including 188 HTTP operations, explicit context and credential selection, Agent/Network configuration and release operations, upload, remote manifest comparison and explicit-ID checkpoint reconciliation.

**The full design is not complete.** See [implementation status](docs/STATUS.md), [usage](docs/USAGE.md), [operation catalog](docs/OPERATIONS.md) and the [canonical design](docs/PLAN.md). Proposed server capabilities require OpenAPI advertisement and return unsupported errors when absent. No permissions are granted by the client.

Validation: `make check`. Distribution: `bash scripts/package.sh VERSION`.
