# Changelog

## Unreleased

- Add guarded semantic local rebase and explicit resolved-content merge checkpoints; preserve old revisions and protect newly referenced local support files.

- Add guarded local Agent/Network revision checkout with exact retained author YAML and support bytes; executable identity is checked before restore.

- Add negotiated Agent/Network isolated Draft storage with generation CAS, exact
  checkpoint registration and read-only reconciliation of uncertain writes.
- Qualify closed-package portable definition digests against shared Go/Python
  fixtures and exclude unrelated Providers from checkpoint dependencies.

- Add Agent/Network origin tracking, immutable local checkpoints and portable
  objects, explicit remote status/current reads and append-only history fetch.

- Preserve tab-indented instructions when writing local YAML descriptors by
  escaping tab-containing strings instead of emitting an unreadable literal block.

- Embed the portable skill in native CLI binaries: `woobe skill install/status/uninstall/agents`, offline without npm/Node; share receipts with the optional standalone installer.

- Add `woobe-cli-skill`, a separate npm package and explicit coding-agent skill installer.
- Support Codex, legacy `.codex`, Claude Code, Copilot, Cursor, user/project scopes and custom roots.
- Preserve local edits/settings with receipts, preflight, locks and rollback.
- Include offline task references/YAML templates and six-host npm skill qualification.
- Add immutable skill release artifacts and independent, opt-in OIDC publishing.
- Refresh installation/authentication, release/CI and historical documentation boundaries.

## 0.13.7

- Recover intentionally deleted local descriptors with safe registry prune.
- Preserve referenced resources, bindings and server identities during local cleanup.

## 0.1.4

- Find draft releases by immutable ID when the published by-tag endpoint returns 404.
- Resume partial uploads and verify all remote assets before finalizing publication.

## 0.1.3

- Support Actions installation tokens whose repository metadata omits user-role permissions.
- Verify repository identity and retain authoritative GitHub write authorization.
- Show publication errors in check annotations and version tags in run titles.

## 0.1.1

- Publish the exact validated CLI to GitHub Packages (GHCR), Linux amd64/arm64.
- Gate image promotion and GitHub Release on source, digest and native smoke checks.
- Explicit version tags enable a first release without merging into master.

## 0.1.0 (distribution foundation)

- Automatic GitHub release on every master push, with immutable source identity
  and SemVer increments calculated from tags and conventional commit history.
- Reviewed source VERSION floor; no automatic source pushes or PR merges.
- Six Linux/macOS/Windows x64 and arm64 archives, schemas, notices and checksums.
- Full reusable CI gate, immutable candidate artifact ID and manifest digest,
  conflict refusal, deterministic archives and partial publication recovery.
- Verified public release and permanent manifest using built-in GITHUB_TOKEN.
- npm distribution deferred; current builds and releases have no npm dependency.
- Installation, Workspace Control Key and context setup documentation.

The existing CLI operations and functional limitations are described in
[STATUS.md](docs/STATUS.md). This version does not claim completion of the full
CLI design or live acceptance of every operation.
