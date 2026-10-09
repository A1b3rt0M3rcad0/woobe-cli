# Portable transfer

Use package export/import for portable composition. Root `export` and
`agent export` are partial JSON projections, not apply-ready packages.

```sh
woobe package export agent AGENT_UUID
woobe package export network NETWORK_UUID --env production
woobe package export agent AGENT_UUID --env release --version VERSION
woobe package validate ./exported-agent --locked
woobe package bindings ./exported-agent --destination ./destination.yaml
woobe package plan ./exported-agent --bindings ./destination.yaml --save-plan ./plan.json
woobe package import --plan-file ./plan.json --wait
```

Export also accepts an exact name. It defaults to current Draft and creates
`./normalized-resource-name`; existing destinations are not replaced. Do not
require snapshot UUIDs. `--version` selects a Release with `--env release`;
`--package-version` overrides portable metadata only.

Inspect declared requirements and fill the generated binding template with
destination native identities/private references. Omit `--bindings` when there
are no requirements. Keep bindings/plans outside the portable directory and
never put raw credentials in them. SOURCE also accepts `woobe.yaml` or tar.gz.

A present lock is always verified, even without `--locked`. That flag additionally
requires a lock. Do not edit the captured copy or remove its lock to hide drift.
Create a separate editable author, edit it, then produce a new sealed copy:

```sh
woobe package edit ./captured --destination ./author
woobe package validate ./author
woobe package seal ./author --destination ./sealed
woobe package validate ./sealed --locked
```

`edit` verifies the sealed source before copying and omits the old lock only in
the new author directory. `seal` validates the current author and writes a new
inventory lock in another directory. Neither overwrites a destination or mutates
the source; server semantics and destination permissions still require a plan.

Import creates Draft resources even when exported from Production; it is not the
managed update flow. Editing an existing resource uses pull/push. Repeating an
import reconciles its checkpoint, not an intentional second create. Keep the
checkpoint and use `woobe package status --checkpoint PATH --wait` after uncertainty.
