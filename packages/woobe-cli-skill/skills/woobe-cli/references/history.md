# ASaC origin and history

These commands need a CLI build with ASaC history support. Check targeted help
first; do not infer availability from the presence of a runtime Release command.
Remote reads also need the server's ASaC read contract. A legacy server can reject
them; do not silently downgrade that observation to a verified current selection.

Keep lockfiles, revisions and referenced objects with author files. Keep `.state`
private. Lockfile origin is historical; observing a changed Production must not
rewrite it. A revision is a local checkpoint, not publication or approval.

```sh
woobe agent '@support' revision create --message 'Correct citation instructions'
woobe agent '@support' history
woobe agent '@support' history verify
woobe agent '@support' status --remote
woobe agent AGENT_UUID current --env production
woobe agent '@support' history fetch
```

Networks use the same vocabulary. Local history works without remote access.
`history fetch` adds immutable receipts without replacing YAML or Production;
it verifies traversal through a stable scoped watermark. History metadata,
object availability and runtime executability are different facts.

Compare `local_modified`, `draft_changed` and `tracked_ref_changed` separately.
Offline status is `unverified`. `legacy_unavailable` environment generation is
not a CAS token. Hashes prove integrity, not server signatures or verified Git
provenance. Missing objects or unknown origin require explicit recovery; never
invent parents, commits, bindings or activation authority.

History support does not change legacy publication effects. Agent publish creates
a Release; Network publish also activates Production. Checkpoint creation does
not create isolated Draft lines or provide fenced deployment guarantees.
