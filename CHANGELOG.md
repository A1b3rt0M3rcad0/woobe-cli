# Changelog

## 0.1.0 (prepared, unreleased)

- Independent Semantic Version in `VERSION`, mirrored in the npm manifest and
  embedded in versioned native builds with their source commit.
- `woobe-cli` npm tarball with the `woobe` command and six native executables,
  usable through npm/npx without Go, installation scripts or binary downloads.
- Linux/macOS/Windows x64 and arm64 archives, shared checksum metadata and
  byte-for-byte verification against the npm package.
- Tag releases require matching canonical versions and a commit integrated in
  master, reuse complete CI validation and publish the validated artifacts.
- Optional npm trusted publishing with provenance after maintainer setup.
- Installation, Workspace Control Key and context setup documentation.

The existing CLI operations and functional limitations are described in
[STATUS.md](docs/STATUS.md). This version does not claim completion of the full
CLI design or live acceptance of every operation.
