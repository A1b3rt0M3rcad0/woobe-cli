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
grant authority or provide deployment fencing. Isolated Draft storage is explicit
and described below.
Check actual command/server compatibility before using a newer contract.

## Isolated remote Drafts

```sh
woobe agent '@support' draft list
woobe agent '@support' draft open hotfix --from production
woobe agent '@support' draft select hotfix
woobe agent '@support' draft show
woobe agent '@support' draft push
woobe agent '@support' revision create --message 'Correct citation instructions'
woobe agent '@support' draft checkpoint REVISION_ID
woobe agent '@support' draft reconcile
```

Networks use the same commands. Opening captures an authorized exact source
(`draft`, `staging`, `production`, or `release --version VERSION`) and selects its
remote Draft for subsequent `draft push`. It preserves local author files and
the existing native default Draft. It is not a local checkout: review/diff local
YAML against the chosen source before saving it. Ordinary `push` remains the
existing native Draft synchronization command; the explicit `draft push` saves
an isolated closed object. No publication or activation occurs.

The server must advertise compatible isolated Draft support and the portable
definition digest scope. An unsupported server is refused before opening/saving;
there is no fallback that writes the default Draft instead. Save uses the last
observed Draft generation. `draft show` refreshes it; review remote changes before
retrying a conflict. Checkpoint registration requires an existing local revision
whose exact object matches the saved isolated Draft and current YAML. Changing
YAML after sealing is rejected. Providers outside its closure are excluded.

Draft selection and pending operation payloads live in destination-scoped private
`.state`. An operation ID is persisted before its request. A definite rejection
clears the pending write; an uncertain outcome blocks additional writes until
`draft reconcile` observes the original receipt. Reconcile never repeats POST/PUT.
`not_observed` is incomplete and does not authorize another mutation. Read-only
listing/inspection remains available. A fresh clone can open a Draft through its
matching durable origin; the server verifies authority and logical/native binding.
Use `draft select NAME_OR_UUID` to recover an existing isolated Draft selection
without opening a duplicate. It verifies the logical UID and native root and
refreshes that Draft's observed generation. It writes only private selection
state: author files and server resources remain unchanged. `--dry-run` performs
these reads without saving the selection. An unresolved original write blocks
selection until reconciled. Selection does not recover native push bindings or
Provider credentials; keep/rebind those separately. A known origin cannot use
`create` to duplicate its root merely because private bindings are absent.

Isolated saves currently validate closed package schema/inventory. Output reports
`semantic_validation: stage_required`; runtime candidate qualification and fenced
deployment are separate negotiated contracts. The storage commands alone do not
make an isolated Draft executable or eligible for legacy promotion.

## Restore local author source

```sh
woobe agent '@support' checkout --revision REVISION_ID
woobe network '@customer-support' checkout --revision REVISION_ID
```

New checkpoints retain exact author descriptors (including YAML comments) and
support bytes in bounded immutable `.woobe/objects/author/<digest>.json` objects.
The revision links this source using `author_artifact_digest`; this hash is
separate from the verified portable executable definition. Checkout compiles the
retained source in an isolated scratch registry and checks its definition and
component identities before writing. It preserves current aliases, paths, origin
and private destination bindings. The working revision changes locally; remote
Drafts, Releases and Production do not change.

Checkout requires working files to match their current source checkpoint. It
rejects edits, changed shared dependency files, tampered source objects and
unregistered target files. `--yes` cannot bypass those protections. Checkpoint
edits explicitly before restoring another revision. Older revisions without
retained author source remain readable but cannot use this local restore.

Remote checkpointing negotiates `author_source_retention` and uploads the exact author object before registering its immutable revision. Upload tickets are actor scoped and expire; retained source remains readable through the authorized Agent/Network revision with another permitted CLI Key. Older backends retain declarations only (`declared_not_retained`). The server labels verified canonical custody separately from compiler attestation.

## Reconcile sealed local branches

```sh
woobe agent '@support' rebase --onto REVISION_ID
woobe network '@customer-support' rebase --onto REVISION_ID --message 'Reconcile network changes'
woobe agent '@support' revision merge REV_A REV_B --message 'Resolved combined content'
```

Rebase computes a three-way semantic merge of the current checkpoint and the
selected revision against their known common ancestor. It merges independent
object fields, treats arrays and support files as atomic, validates the complete
result in scratch space, and seals a new revision whose parent is `--onto`.
Original revisions and origin remain unchanged. Surviving YAML field comments
are carried forward when fields are recomposed; a comment-only incoming change
is also retained. Competing comment edits report `/author_comments`; deleting a
component concurrently with an author-source edit reports `/author_source`.
Deleted fields have no place to retain their comments. JSON descriptors remain JSON. Private destination bindings
are unchanged. The result is local, not a remote push or deployment.

Rebase requires clean checkpointed files. Conflicts report component/field paths
with `ASAC_REBASE_CONFLICT` (exit 6) and do not change files or revision history.
Incomplete ancestry is not guessed; obtain missing history first. Multiple common
ancestors require an explicit `--base REVISION_ID`, which must be an ancestor of
both branches. Review `revision diff` before writing; `--yes` does not discard edits.

`revision merge REV_A REV_B` seals the **currently resolved author content** with
both immutable parents. It does not silently choose a winner or merge YAML for
you. Resolve conflicts and validate the resulting files before this checkpoint.
Neither command publishes, stages, activates nor rewrites a shared revision.

Source verification uses the compiler recipe recorded in the author object.
Older objects without a recipe use the original Provider-to-Model lowering;
new objects use `woobe-development-compiler@2.0`, keeping Model endpoint overrides
separate from Provider defaults. Checkout validates the historical object using
its own recipe; a new checkpoint uses the current one. Unknown recipes are
rejected rather than silently reinterpreted.


### Fetch isolated revision metadata

Use `woobe agent @support history fetch --revisions` (or the same Network
command) to append the authorized server revision catalog. It checks immutable
record digests, logical identity, page completeness and a stable watermark.
The default `history fetch` continues to retrieve native release/deployment
receipts. Neither mode changes author YAML, the working head or Production.

Revision metadata and source materialization are separate: this fetch reports
`objects_available: not_downloaded`. A remote source digest is a declaration,
not a retained author object or a server attestation. Local checkout still
requires verified executable and exact author objects. `--dry-run` validates
the catalog without storing metadata. A failed traversal can retain earlier
immutable pages, but never reports that traversal as complete.


### Hydrate executable and exact author objects

```sh
woobe agent @support history fetch --revisions
woobe agent @support revision hydrate rv_00000000-0000-0000-0000-000000000000
```

Replace the revision placeholder with a verified retained ID. Agent and Network
use the same operation. It negotiates backend hydration support, reauthorizes
the retained resource and verifies the receipt, transport SHA-256, inventory,
artifact digest and executable definition before storing an immutable object.
`--dry-run` downloads and verifies without storing. Source YAML, native bindings,
working head and Production are preserved. If source custody is available, hydration verifies its canonical author hash, exact YAML/support inventory and recorded compiler recipe against the executable definition and every component identity before storing source. It reports `author_object: verified_source_and_compilation`; a separate guarded checkout restores working files. An old record without retained source still permits executable hydration and reports source unavailable. Credentials and native bindings remain separate; hydration does not qualify runtime execution.
