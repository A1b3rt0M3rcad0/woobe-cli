# Title

ci: publish verified distribution artifacts and exercise native CLI binaries

# Description

The old workflow cross-built executables but left complete distribution archive verification local. CI now packages all six Linux/macOS/Windows amd64/arm64 targets, verifies source identity/content/checksums, and downloads the same artifacts for native executable smoke on three hosted OS runners. An aggregate check requires tests, packages and all native smoke instances to succeed, including refusal of cancelled/skipped prerequisites. Evidence files are stored in Actions with finite retention.

The native checks validate embedded build identity, essential discovery commands, both manifest schemas, local resource validation and malformed-input exit behavior. They do not certify keychains, real Woobe authorization/runtime or every native architecture. Development CI builds do not publish tags/releases.

PRs #1–#6 are merged into master. This continuation also reconciles the README completeness count and records the workflow/remote artifact contract in CI.md. The full plan remains 39/101 deliveries (38.6%); real-server phase acceptance remains 0/10. Exact CI run results are reported in the GitHub PR metadata.
