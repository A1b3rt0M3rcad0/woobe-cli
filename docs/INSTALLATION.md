# Installation and automatic CLI releases

The CLI stays in this repository and is versioned independently of the Woobe
backend. Distribution uses **GitHub Releases** and **GitHub Packages (GHCR)**: six native archives for Linux,
macOS and Windows, in amd64/x64 and arm64. Go, Node.js and npm are not required
to use the binaries. npm distribution is deferred.

The initial release floor is **0.1.0**. Publication requires the complete CI
gate. Master pushes calculate the next version automatically; an explicit
`v<VERSION>` tag also releases its exact source without requiring a merge.

## Download and run

Open [the latest release](https://github.com/A1b3rt0M3rcad0/woobe-cli/releases/latest)
and download your platform archive plus `SHA256SUMS`.

| Platform | Example archive |
| --- | --- |
| Linux x64 | `woobe_0.1.4_linux_amd64.tar.gz` |
| Linux arm64 | `woobe_0.1.4_linux_arm64.tar.gz` |
| macOS Intel | `woobe_0.1.4_darwin_amd64.tar.gz` |
| macOS Apple Silicon | `woobe_0.1.4_darwin_arm64.tar.gz` |
| Windows x64 | `woobe_0.1.4_windows_amd64.zip` |
| Windows arm64 | `woobe_0.1.4_windows_arm64.zip` |

Verify the downloaded file's SHA-256 against `SHA256SUMS` before extraction.
Use `sha256sum` on Linux, `shasum -a 256` on macOS or
`Get-FileHash -Algorithm SHA256` in PowerShell. Extract, place `woobe` or
`woobe.exe` in a directory on `PATH`, then run:

```sh
woobe version
woobe help
```

The archives also contain usage instructions, both manifest schemas and pinned
dependency/toolchain license notices. No administrator privileges are required.

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

Use the backend API URL. Runtime credentials are separate. See [USAGE.md](USAGE.md)
for secret handling, runtime and manifests. Distribution does not change server
authorization or the [functional coverage](STATUS.md).

## Automatic release

**Every push to `master` starts Release CLI**, including a maintainer's merge of
a reviewed PR. No manual version update, tag push, release button, npm account
or additional release token is required. The workflow never merges PRs and
never pushes changes into `master`.

1. Pin the source to the push's exact SHA, require it to belong to `master`,
   and calculate the version from existing immutable `v*` tags and commit history.
2. Run the full reusable CI: Go/race/fuzz checks, six cross-builds, archive/schema/
   license/checksum validation, and native smoke tests on Linux/macOS/Windows.
3. Download the immutable Actions artifact ID and verify its manifest SHA-256.
   Extract the tested Linux binaries into a non-root, multi-platform GHCR image
   for `linux/amd64` and `linux/arm64`, without rebuilding the executables.
   Verify its digest, platforms and CLI behavior before publication.
4. Preflight both destinations and refuse conflicting tags, assets or image
   digests. Reserve `v<VERSION>` on the exact validated source, promote the image
   to `<VERSION>`, and create or resume a draft GitHub Release. Upload only
   missing files; existing bytes are never overwritten.
5. Download and compare every asset, execute the native CLI again, publish the
   Release and verify all six archives, checksums, artifacts.json and
   release-manifest.json. The permanent manifest records the package digest.
   Promote the verified stable image to `latest`; older recovery cannot move
   either destination's latest version backwards.

The workflow uses the built-in **GITHUB_TOKEN** with job permissions
`contents: write` and `packages: write`. No additional token or npm account is
required. Repository policies must allow these permissions and tag creation.
Explicit owner-created version tags trigger the same gates. Tags created by the
workflow's GITHUB_TOKEN do not trigger another workflow run.

### GitHub Packages

The repository package is `ghcr.io/a1b3rt0m3rcad0/woobe-cli`. After publication:

```sh
docker run --rm ghcr.io/a1b3rt0m3rcad0/woobe-cli:0.1.4 version
docker run --rm ghcr.io/a1b3rt0m3rcad0/woobe-cli:latest help
```

Use a volume at `/data` to persist credentials and contexts. The image runs as
UID 10001, so the directory must be writable by that UID. Package visibility is
a GitHub setting: new GHCR packages may start private even for public repositories.
The owner can change visibility in the package settings; private pulls require
GitHub authentication. Native release downloads are public independently.

### Automatic version calculation

`VERSION` is the reviewed **minimum release version**, initially `0.1.0`.
The reserved private npm source manifest mirrors this floor for future work.
Published versions are authoritative in Git tags and permanent artifact
manifests; automatic increments do not rewrite source files.

- No previous version: use `VERSION`.
- A `feat:` or `feat(scope):` commit since the latest tag: increment MINOR.
- A conventional `type!:` / `type(scope)!:` subject or `BREAKING CHANGE:` /
  `BREAKING-CHANGE:` footer: increment MAJOR; before 1.0, increment MINOR.
- Other changes: increment PATCH. A previous prerelease becomes its stable core
  version when no feature/breaking bump is needed.
- If a reviewed increase in `VERSION` is higher than the calculated version,
  use that floor. No source version edit is needed for ordinary releases.

Examples: `0.1.0` + `fix(cli): ...` → `0.1.1`; `feat(cli): ...` → `0.2.0`;
`1.2.3` + `fix(api)!: ...` → `2.0.0`.

### Retries and manual recovery

Release transactions are serialized and active publication is not cancelled by
later pushes. A rerun for an already tagged source recovers the same version,
even when master has advanced. A pending push superseded by a newer tagged
source is skipped. Tags reserved by a failed upload are never reused for a
later commit; the next push advances the version.

To publish an explicitly selected source without merging, create and push an
annotated `v<VERSION>` tag pointing to it. The tag must match the selected source
and version floor. Never move or reuse an existing tag.

For recovery or an explicit prerelease, **Actions → Release CLI → Run workflow**
remains available on `master`. Leave both fields empty to release its current
source automatically. To repair a specific release, enter its existing version
without `v`; its tag recovers the source. An optional full `revision` can pin the
original SHA. New explicit versions must advance all reserved versions and
respect the source floor. Every recovery still passes the complete CI gate.

Archive timestamps, owners and compression metadata are normalized. Rebuilding
the same source/version with pinned toolchains yields the same bytes. A
conflicting tag or asset is never overwritten; use a new version instead.

## Local build

```sh
make build
bin/woobe version
make package
python3 scripts/verify_artifacts.py --commit "$(git rev-parse HEAD)"
python3 scripts/smoke_artifacts.py --commit "$(git rev-parse HEAD)" --report native-smoke.json
```

Requires the Go toolchain in `go.mod`, Python 3, GNU tar, gzip, zip and sha256sum.
`make build` and `make package` use the source floor locally. Release CI passes
its automatically calculated version into packaging and embeds that version and
the exact source SHA in every binary. A bare `go build` keeps the `dev` version.
Development PR/branch packages use `0.0.0-ci.<run>.<attempt>` and are Actions
artifacts, not public releases.

The private npm wrapper and opt-in packaging helpers are retained for future
work. Current CI and Release CLI do not build, install or publish npm packages,
access an npm registry, use OIDC or require NPM_TOKEN. This repository currently
supplies no software license file; packaging adds dependency notices without
introducing a software license grant.
