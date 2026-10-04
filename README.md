# Woobe CLI

External Go client for the Woobe public control plane and runtime. Server-side policy is authoritative.

```sh
go build -o bin/woobe ./cmd/woobe
bin/woobe help --output json
bin/woobe schema --command "project agent create"
```

The initial implementation provides 188 registered API operations, additional local/SDK commands, explicit context and credential selection, Agent/Network configuration and release operations, upload and manifest checkpoints.

**The full design is not complete.** See [implementation status](docs/STATUS.md), [usage](docs/USAGE.md), [operation catalog](docs/OPERATIONS.md) and the [canonical design](docs/PLAN.md). Proposed server capabilities return structured unsupported errors. No permissions are granted by the client.

Validation: `make check`. Distribution: `bash scripts/package.sh VERSION`.
