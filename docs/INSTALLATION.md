# Installation and CLI releases

The CLI stays in this repository and has its own Semantic Version, independent
of the Woobe backend. `VERSION` is authoritative; the npm manifest mirrors it.
The first prepared version is **0.1.0**. A prepared version is not a published
release: the commands below that use npm or GitHub URLs become available after
publication. Build and install locally until then.

## Install with npm

Requires Node.js 22 or newer; Go is not required.

```sh
npm install --global woobe-cli
woobe version
woobe help
```

Without a global installation:

```sh
npx --package=woobe-cli woobe version
```

For reproducible automation, pin the package:

```sh
npx --package=woobe-cli@0.1.0 woobe version
```

The tarball contains Linux, macOS and Windows binaries for x64 and arm64. It has
no dependencies or installation scripts and does not download executables at
installation or runtime. `npm install --ignore-scripts` is supported.

GitHub Releases also contains the npm tarball. It can be installed before a
registry publication is configured:

```sh
npm install --global https://github.com/A1b3rt0M3rcad0/woobe-cli/releases/download/v0.1.0/woobe-cli-0.1.0.tgz
woobe version
```

## Download without Node or Go

Open [GitHub Releases](https://github.com/A1b3rt0M3rcad0/woobe-cli/releases) and
select your platform. The executable is `woobe` on Linux/macOS and `woobe.exe`
on Windows. `amd64` means x64; Apple Silicon uses `darwin_arm64`.

| Platform | Archive for 0.1.0 |
| --- | --- |
| Linux x64 | `woobe_0.1.0_linux_amd64.tar.gz` |
| Linux arm64 | `woobe_0.1.0_linux_arm64.tar.gz` |
| macOS Intel | `woobe_0.1.0_darwin_amd64.tar.gz` |
| macOS Apple Silicon | `woobe_0.1.0_darwin_arm64.tar.gz` |
| Windows x64 | `woobe_0.1.0_windows_amd64.zip` |
| Windows arm64 | `woobe_0.1.0_windows_arm64.zip` |

Download `SHA256SUMS` with the archive and compare its SHA-256 before extraction
(`sha256sum` on Linux, `shasum -a 256` on macOS, `Get-FileHash -Algorithm SHA256`
on Windows). Extract the archive, put the executable in a directory on `PATH`,
then run `woobe version`. The archives also contain usage instructions and the
two manifest schemas. No administrator privileges are required.
All archives and the npm package include the pinned dependencies' and Go
toolchain's license notices in `THIRD_PARTY_NOTICES.txt`.

## Connect Woobe

Create a Control Key at **Workspace settings → CLI Keys** with the Workspace
permissions and Project grants your integration needs. Import it from stdin:

```sh
woobe auth credential import --name workspace-cli --stdin < workspace.key
woobe context create woobe --api-url https://YOUR_WOOBE_API
woobe context credential attach woobe --credential workspace-cli
woobe context use woobe
woobe context set --workspace WORKSPACE_ID --project PROJECT_ID
woobe project agent list
```

Use the backend API URL, not the frontend URL. Runtime credentials are separate.
See [usage](USAGE.md) for secret handling, runtime and manifests. Packaging does
not change the [functional coverage](STATUS.md) or server authorization rules.

## Build locally

```sh
make build
bin/woobe version
make package
npm install --global ./dist/woobe-cli-0.1.0.tgz
woobe version
```

Building requires the Go toolchain declared in `go.mod`, Python 3, Node.js/npm,
and the existing archive tools (`tar`, `zip`, `sha256sum`). `make build` embeds
the canonical version and source commit. A bare `go build` retains the `dev`
version for an unversioned development build.

## Prepare the next version

Use `MAJOR.MINOR.PATCH`, or a prerelease such as `0.2.0-rc.1`. Build metadata and
the `v` prefix are not part of canonical `VERSION`. Increment PATCH for fixes,
MINOR for additions and MAJOR for incompatible public changes; before 1.0,
incompatible changes increment MINOR and must be documented in the changelog.

```sh
python3 scripts/version.py --set 0.1.1
python3 scripts/version.py
make check
make package
```

Update `CHANGELOG.md` and commit the version change in a PR. Existing releases
are immutable: use a new version for different bytes. Merge and tagging remain
explicit maintainer actions; the version script does neither.

After the reviewed commit is integrated in `master`, tagging that exact commit
as `v<VERSION>` starts **Release CLI**. The workflow rejects a tag/version
mismatch or a commit outside `master`, runs the full reusable CI (including
native and npm installation smoke tests on three OSes), then publishes those
exact validated artifacts. It does not rebuild them in the publication job.
`artifacts.json` records the version, commit and SHA-256 of every artifact;
the npm binaries must match the corresponding native archives byte for byte.

CI branch/PR packages use `0.0.0-ci.<run>.<attempt>` and remain Actions artifacts.
They do not publish tags, GitHub Releases or npm packages.

## npm registry publication

The npm name is `woobe-cli`. The source manifest is private so an accidental
publish from `packages/woobe-cli` cannot publish a package without binaries.
Publish only the generated, verified `dist/woobe-cli-<VERSION>.tgz`.

For a new npm package, a maintainer must perform its first registry publication
with npm authentication after the GitHub Release validation, for example:

```sh
npm publish ./dist/woobe-cli-0.1.0.tgz --ignore-scripts --access public
```

Then configure the npm package's **Trusted Publisher** for this GitHub repository
and workflow file `release.yml`, and set the repository variable
`WOOBE_PUBLISH_NPM=true`. Later tag releases use npm 11's OIDC trusted publishing
and provenance, without an npm token stored in the repository. Stable versions
use npm's `latest` tag; prereleases use `next` and GitHub's prerelease flag.
If the variable is unset, GitHub binary/tarball publication still works and the
npm registry job is skipped. Registry ownership/configuration is external to
this code change; no package or release is published by preparing this PR.

This repository currently supplies no software license file. Packaging does
not introduce a license grant; the npm metadata stays `UNLICENSED`.
