# ASaC tracking and immutable local revisions

Managed pull bootstraps `agent.lock.yaml` or `network.lock.yaml` in the registered
resource directory. Flat descriptor paths use a sibling `<stem>.asac/` directory
to avoid collisions. Tracking records resource UID, origin destination and native
snapshot, tracking environment, working revision and portable requirements. It
excludes credentials, export tickets, leases and private CAS observations. Repeated
pulls/checkpoints preserve origin. Unproven origin is `unknown_origin`.

```sh
woobe agent '@support' revision create --message 'Correct citation instructions'
woobe agent '@support' history
woobe agent '@support' heads
woobe agent '@support' revision show REVISION_ID
woobe agent '@support' revision diff REV_A REV_B
woobe agent '@support' history verify
```

Networks use the same commands. Checkpoints have random opaque IDs, explicit
parents and a closed package captured through the existing compiler.
`--parent REVISION_ID` can name multiple existing parents for already resolved
content; it does not automatically merge anything. Same content with different
lineage remains a different revision. Sealed records are append-only. Objects
live at `.woobe/objects/sha256/<artifact-digest>.tar.gz` and use normal package
inventory/locked-archive validation.

`definition_digest` reflects dependencies and support bytes, independently of
aliases and destination credential IDs. `record_digest` includes lineage and
identity but excludes itself. Shared Go/Python fixtures define the type-preserving
`woobe-asac@1.0` canonicalization. Existing package artifact digests are unchanged.
Missing parents, missing objects and runtime executability are separate facts.
Local history does not promise remote freshness. A clone can inspect/checkpoint
using tracking requirements without copying `.state`; remote writes still require
validated bindings and authority.

## Observe the remote

```sh
woobe agent '@support' status --remote
woobe agent AGENT_UUID current --env production
woobe agent '@support' history fetch
```

`current` is always remote and accepts UUID without a local registry. Status stays
local unless `--remote` is supplied, and distinguishes `local_modified`,
`draft_changed` and `tracked_ref_changed`. Observing changed Production does not
rewrite origin. Legacy selection reports `generation_status: legacy_unavailable`;
a release ID is not an environment CAS generation.

Fetch traverses stable watermarks for publication/deployment streams and stores
destination-scoped immutable receipts without editing author files or Production.
Dry-run validates pages without writing receipts. Hashes verify integrity;
authenticated observations are not signed offline attestations or verified Git
provenance.

These commands do not change native lifecycle effects: Agent publish creates a
Release; legacy Network publish also activates Production. Checkpoints do not
grant authority, create isolated Draft lines or provide deployment fencing.
Check actual command/server compatibility before using a newer contract.
