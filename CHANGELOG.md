# Changelog

## 0.1.0 (prepared, unreleased)

- Independent Semantic Version in `VERSION`, mirrored in the npm manifest and
  embedded in versioned native builds with their source commit.
- `woobe-cli` npm tarball with the `woobe` command and six native executables,
  usable through npm/npx without Go, installation scripts or binary downloads.
- Linux/macOS/Windows x64 and arm64 archives, shared checksum metadata and
  byte-for-byte verification against the npm package.
- Manual Actions releases require matching canonical versions and an exact
  reviewed master commit, reuse complete CI validation and publish both GitHub
  and npm from the validated immutable artifact ID and manifest digest.
- All-destination preflight, deterministic archives, conflict refusal, recovery
  of partial publication and a verified permanent release manifest.
- npm token bootstrap and subsequent trusted publishing with provenance.
- Installation, Workspace Control Key and context setup documentation.

The existing CLI operations and functional limitations are described in
[STATUS.md](docs/STATUS.md). This version does not claim completion of the full
CLI design or live acceptance of every operation.
