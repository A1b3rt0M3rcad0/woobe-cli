# Installable Woobe CLI skill for coding assistants

Coding assistants currently need to explore CLI help and full responses to learn
Woobe workflows. Add a separate `woobe-cli-skill` npm package with a small portable
Agent Skill, seven offline task references and two YAML templates. It documents
real connection, managed Draft editing, shared dependencies, package transfer,
execution tests, lifecycle, compact output and recovery commands.

`woobe-skill` explicitly installs for Codex (`.agents` and legacy `.codex`), Claude
Code, Copilot, Cursor or a custom skills root, at project/user scope. Status and
dry-run are read-only; updates/removal preserve edited/unmanaged files and other
settings. Receipts, hash checks, preflight, locks, directory swaps and rollback
protect managed installation. The installer has no dependencies, lifecycle scripts,
API calls, key handling or implicit CLI installation.

Distribution builds bind both npm tarballs to the same version/source commit.
The release manifest/checksums include the skill tarball; all six native runners
exercise the actual npm package offline and verify CLI examples. Existing CLI
OIDC publishing stays independent. A gated `npm-skill` job requires the new name's
one-time owner bootstrap, its own Trusted Publisher and repository variable;
GitHub releases include the skill tarball even before enabling npm publication.

Documentation is refreshed across installation/authentication, CI/release,
usage/development, output, migration, package boundaries, current status and
historical evidence. A complete documentation index and regenerated command
catalog distinguish current behavior from original roadmap measurements.

Validation: `make check`, immutable artifact verification, offline native/npm/skill
smokes, installer regressions and all fenced skill command examples. Hosted CI
qualifies Linux/macOS/Windows × amd64/arm64. No backend code changes; no merge,
tag or npm publication is performed by this PR. Merge requires owner approval.
