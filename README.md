# Woobe CLI

Command-line client for the [Woobe](https://github.com/A1b3rt0M3rcad0/woobe)
control plane and runtime. Manage Workspace and Project resources, validate
manifests, publish Agent and Network releases and execute released targets.
The backend authorizes every operation.

## Install and run

Download your platform archive and `SHA256SUMS` from
[the latest GitHub Release](https://github.com/A1b3rt0M3rcad0/woobe-cli/releases/latest).
The published `v0.1.4` release includes all six targets:

| System | x64 | arm64 |
| --- | --- | --- |
| Linux | `linux_amd64.tar.gz` | `linux_arm64.tar.gz` |
| macOS | `darwin_amd64.tar.gz` | `darwin_arm64.tar.gz` |
| Windows | `windows_amd64.zip` | `windows_arm64.zip` |

Verify the checksum, extract the archive and place `woobe` or `woobe.exe` on
`PATH`. For example, on Linux x64:

```bash
sha256sum --check --ignore-missing SHA256SUMS
tar -xzf woobe_0.1.4_linux_amd64.tar.gz
./woobe version
./woobe help
```

On macOS, compare `shasum -a 256` with the checksum file. On Windows, use
`Get-FileHash -Algorithm SHA256`, extract the ZIP and run `.\woobe.exe version`.
No Go, Node.js or npm installation is required to use the native binaries.

The public [GitHub Packages image](https://github.com/users/A1b3rt0M3rcad0/packages/container/package/woobe-cli)
supports Linux amd64 and arm64:

```bash
docker run --rm ghcr.io/a1b3rt0m3rcad0/woobe-cli:latest version
docker run --rm ghcr.io/a1b3rt0m3rcad0/woobe-cli:0.1.4 help
```

Pin a version for automation. Mount persistent configuration at `/data`, writable
by UID 10001. See [installation and container usage](docs/INSTALLATION.md).

## Connect to Woobe

Create a Control Key at **Workspace settings → CLI Keys → Create CLI key**.
Choose the Workspace role you need and add Project grants separately, selecting
a role for each Project. Set an expiration date and copy the secret before
closing the issuance view. Rotation replaces the secret immediately.

Import a private key file from stdin on Linux or macOS, then select your backend
origin and scope:

```bash
woobe auth credential import --name workspace-cli --stdin < workspace.key
woobe context create woobe --api-url http://localhost:8000
woobe context credential attach woobe --credential workspace-cli
woobe context use woobe
woobe context set --workspace WORKSPACE_ID --project PROJECT_ID
woobe project agent list
```

Keep `workspace.key` private; never commit it or include its value in command-line
arguments. Use the backend API URL, rather than the Web UI or ChatSurface URL.
Windows and CI can supply `WOOBE_CONTROL_KEY` through their environment; protected
local import currently uses the POSIX credential store.

See [the Woobe CLI guide](https://github.com/A1b3rt0M3rcad0/woobe/blob/master/docs/CLI.md)
for calendar expiration, predefined roles, masked input and platform setup.

## Discover and automate

```bash
woobe help --output json
woobe doctor
woobe schema --command "project agent create" --kind input
woobe validate-input --command "project agent create" --file agent.json
woobe manifest validate --file resources.json
woobe manifest plan --file resources.json
woobe manifest apply --file resources.json --checkpoint apply-checkpoint.json --yes
```

Preserve the manifest and checkpoint when resuming. Exports can be incomplete
projections; they are not automatically safe to replay. The [usage guide](docs/USAGE.md)
documents Agent/Network configuration, publication, streaming, protected inputs
and recovery.

Runtime execution uses a separately selected Runtime Key. The CLI never falls
back to a Control Key for Agent or Network execution:

```bash
woobe auth credential import --name runtime --stdin < runtime.key
woobe context runtime-credential attach woobe --runtime-credential runtime
woobe runtime target run TARGET --target-kind agent --file run.json --yes
```

Keep `runtime.key` private. Use `--target-kind network` for a Network.

## Develop and release

Use the Go toolchain pinned in `go.mod`. Development commands:

```bash
make build
bin/woobe help --output json
make check
make package
```

Source builds use the reviewed `VERSION` floor. Published binaries embed their
release version and exact source commit. CLI and platform versions are independent.

The Release CLI workflow calculates SemVer, validates six native archives,
creates an immutable tag, publishes GitHub Releases and the GHCR version, and
updates `latest` for successful stable releases. Once that workflow is integrated
into `master`, every push to `master` starts this complete flow. It never merges
development PRs or rewrites source versions automatically.

Develop in a branch, open a PR to `master`, pass CI and merge the reviewed change.
Use Conventional Commit subjects, including the final squash subject:

| Change | Version increment |
| --- | --- |
| `fix(cli): correct an input` | PATCH |
| `feat(cli): add a command` | MINOR |
| `feat(cli)!: change a contract` | MINOR before 1.0; MAJOR after 1.0 |

Explicit `v<VERSION>` tags can release the exact tagged source through the same
CI gates before master integration. Rerunning a release recovers that version;
it does not create another version. Published tags, assets and image versions
must never be overwritten. See [release calculation and recovery](docs/INSTALLATION.md).
npm distribution is deferred.

## Coverage and documentation

Unified discovery includes 241 handlers and 190 HTTP operations. The full design
remains incomplete: the audited plan contains 66/101 completed deliveries
(65.3%), with partial deliveries receiving no credit. This planning measure does
not replace acceptance of a particular operation or imply full production coverage.

- [Implementation status and limits](docs/STATUS.md)
- [Usage and command catalog](docs/USAGE.md), [operations](docs/OPERATIONS.md)
- [CI and native validation](docs/CI.md)
- [Plan assessment](docs/COMPLETENESS.md) and [remaining requirements](docs/REMAINING.md)
- [Control Plane authority](https://github.com/A1b3rt0M3rcad0/woobe/blob/master/docs/security/CONTROL_PLANE.md)

Backend integration and Workspace CLI Keys are delivered by
[Woobe PR #177](https://github.com/A1b3rt0M3rcad0/woobe/pull/177).
Server capabilities require advertisement and authorization; the CLI cannot grant
permissions the server does not provide.
