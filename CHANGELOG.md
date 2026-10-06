# Changelog

## 0.1.0 (prepared, unreleased)

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
