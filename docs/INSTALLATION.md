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
are immutable: use a new version for different bytes. Integration remains an
explicit maintainer action; the version script never merges or creates a tag.

## Publish from GitHub Actions

The entrypoint follows Woobe's manual release pattern:

1. Review the product and version changes in PRs. They must already be in
   `master`; the release workflow never merges PRs or changes repository files.
2. Open **Actions → Release CLI → Run workflow**, choose **master**, enter the
   canonical version (for example `0.1.0`) and start the run. Optionally pin the
   exact reviewed 40-character source SHA in `revision`.
3. The workflow resolves that immutable source, requires matching `VERSION`
   and npm metadata, and runs the full reusable CI: Go/race/fuzz validation,
   six cross-builds, and native plus npm installation smoke tests on three OSes.
4. Publication receives the immutable Actions artifact ID and SHA-256 of its
   manifest from CI. It verifies the source/version and every file, then checks
   **both GitHub and npm before writing either destination**. An existing
   version/tag/asset must match the exact validated bytes or source. Unknown
   errors, authentication failures and conflicts stop the transaction.
5. The job publishes the validated npm tarball, downloads and compares the
   registry's exact bytes and exercises the installed command. It then creates
   or resumes a draft GitHub Release at the exact commit, uploads missing files
   without overwriting anything, downloads every asset for comparison, and
   makes the release public. Final success requires both destinations, the
   immutable tag and `release-manifest.json` to match.

There are no tag-push or master-push publication triggers. The **Run workflow**
button appears when this workflow is integrated in the default branch. The old
**Release binaries** workflow only accepted tag pushes; an empty runs list was
expected before its first tag.

GitHub Releases contains six native archives, the npm tarball, `SHA256SUMS`,
`artifacts.json` and the permanent `release-manifest.json`. They are the same
files exercised by CI; the publication job does not rebuild. Stable versions
use npm's `latest`; prereleases use `next` and GitHub's prerelease flag.
New versions must advance the highest published SemVer. A retry of an existing
version never overwrites an npm version, asset, tag or npm distribution tag.

### First npm publication and authentication

The npm name is `woobe-cli`. The source manifest is private so an accidental
publish from `packages/woobe-cli` cannot publish a package without binaries.
Publish only the generated, verified `dist/woobe-cli-<VERSION>.tgz`.

For the **first publication**, configure repository Actions secret **NPM_TOKEN**
with an npm token authorized to create this public package. If the account's
2FA policy requires bypass for automation, grant that capability to this token.
The workflow fails before publication when a new package lacks this credential.
It does not silently skip npm or report a binaries-only release as complete.

After the first publication, configure the npm package's **Trusted Publisher**
for GitHub owner `A1b3rt0M3rcad0`, repository `woobe-cli`, workflow `release.yml`
(no environment name), then remove the bootstrap token. npm 11 uses the job's
OIDC identity and provenance for subsequent releases. Registry ownership and
Trusted Publisher settings remain account configuration. The previous
`WOOBE_PUBLISH_NPM` switch is no longer used: this release publishes both targets.

### Recovery

Re-run failed jobs to reuse the original immutable candidate. For a completely
new run, enter the same version; an existing Git tag recovers its original
commit even if master has advanced. Before tag creation (for example, npm
succeeded but GitHub failed), supply the original `revision` from the failed
run. A partially uploaded draft is repaired only at that exact revision.
Packaging normalizes timestamps, owners and compression metadata so a retry
from the same commit and pinned toolchains produces the same archive bytes.
Conflicts require a new version, never clobbering a release.

CI branch/PR packages use `0.0.0-ci.<run>.<attempt>` and remain Actions artifacts.
They do not publish tags, GitHub Releases or npm packages.

This repository currently supplies no software license file. Packaging does
not introduce a license grant; the npm metadata stays `UNLICENSED`.
