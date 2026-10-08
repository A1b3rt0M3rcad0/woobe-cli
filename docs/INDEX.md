# Documentation index

Current guides describe the executable and the separately installed assistant
skill. Dated validation/audit records remain historical evidence, not a claim
that every original roadmap criterion is complete.

| Guide | Purpose |
| --- | --- |
| [Repository agent instructions](../AGENTS.md) | Maintain CLI behavior, help, docs and the embedded/npm skill together; regeneration and validation rules |
| [Installation](INSTALLATION.md) | Native/npm/container setup, CLI Keys, automatic release, Trusted Publishing and recovery |
| [Coding assistant skill](AGENT_SKILL.md) | Built-in offline install, Codex/.codex, Claude Code, custom paths, compatible npm installer and publisher bootstrap |
| [Usage](USAGE.md) | Canonical CLI workflows, managed YAML editing and recovery |
| [Development](DEVELOPMENT.md) | Provider bindings, shared dependencies and Surfaces |
| [ASaC tracking/history](ASAC.md) | Immutable checkpoints, guarded source checkout, objects, origin locks, isolated Drafts and scoped remote observations |
| [Portable packages](PACKAGES.md) | Agent/Network export/import, environments, bindings and checkpoints |
| [Command catalog](OPERATIONS.md) | Current discovery-derived executable commands and HTTP routes |
| [Output](OUTPUT.md) | Text/compact/JSON, fields and collection evidence |
| [Output examples](OUTPUT_EXAMPLES.md) | Measured illustrative before/after output |
| [Output audit](OUTPUT_AUDIT.md) | Historical fixture measurements and their limits ([CSV](OUTPUT_AUDIT.csv)) |
| [Pagination](PAGINATION.md) | Reviewed cursor/revision contracts, bounds and partial collection handling |
| [Migration](MIGRATION.md) | Python prototype → Go client and optional skill adoption |
| [CI](CI.md) | Nine-job gate, six-host native/npm/skill validation and publication separation |
| [Validation](VALIDATION.md) | Reproducible current checks and dated backend/client evidence |
| [Status](STATUS.md) | Current release/boundaries and historical implementation checkpoints |
| [Remaining roadmap](REMAINING.md) | Historical full-plan backlog, scope and owner publication steps |
| [Original plan](PLAN.md) | Original roadmap plus assistant skill delivery addendum |
| [Roadmap audit](COMPLETENESS.md) | Original 101-item historical assessment ([ledger](COMPLETENESS.json)) |
| [Backend handoff](BACKEND_HANDOFF.md) | Current API boundary and historical paired integration evidence |
| [PR description](PR.md) | New assistant skill change and review/validation scope |
| [Changelog](../CHANGELOG.md) | Unreleased change and historical distribution notes |

Package source READMEs: [native CLI](../packages/woobe-cli/README.md),
[assistant skill](../packages/woobe-cli-skill/README.md). The npm skill README is
built from `AGENT_SKILL.md`; installed instructions/references are offline under
[skills/woobe-cli](../packages/woobe-cli-skill/skills/woobe-cli/SKILL.md).
