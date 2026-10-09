# Recovery without data loss

| Failure | Next step |
| --- | --- |
| Missing private native bindings after clone | `woobe agent '@support' bindings recover --dry-run`, then recover after review (or Network). Entire closure must have accepted server evidence; explicitly bind Providers afterward. Does not recover an accepted YAML base |
| Missing isolated Draft selection after clone | `woobe agent '@support' draft select NAME_OR_UUID` (or Network); matching origin required; no remote mutation or file checkout |
| Rebase comment conflict | Resolve the reported author comments explicitly; checkpoint the resolved content rather than dropping either edit |
| Missing registered descriptor | `woobe config check`, then inspect the named path; restore accidental deletions or preview `woobe resources prune --dry-run` |
| Several intentionally deleted folders | Review the preview, then `woobe resources prune --yes`; only absent descriptors are unregistered |
| Missing Model still used by an Agent | Restore the Model or edit/remove that consumer's reference; pruning must not create a broken graph |
| Unavailable/unsafe local file | Inspect permissions, links and the reported path; do not disable safety checks |
| Package schema/fingerprint mismatch | `woobe package doctor --output compact`; align server/client, restart updated backend when applicable |
| `PACKAGE_EXPORT_INCOMPLETE` | Inspect the narrowed Model/Provider/environment diagnostic; fix the actual source definition rather than installing repeatedly |
| Pull merge conflict | Resolve only listed local/remote fields; preserve author files and the accepted base |
| Push response lost or timed out | Keep `.state/`; `woobe agent '@support' reconcile` or Network equivalent before another write |
| Import uncertain | Keep its checkpoint; `woobe package status --checkpoint PATH --wait`, then reconcile the original operation |
| Pending reconciliation points to status | Observe the given operation until terminal, then rerun managed reconcile to accept bases |
| Publication or permission denied | Stop and report the actual scope/required permission; `--yes` grants no authority |

`resources unregister KIND '@alias'` removes one local registry entry, including
an absent root descriptor, after validating remaining refs. It retains files and
private native bindings and does not archive/delete server resources. Prune does
not delete unused but present dependencies or accept malformed/linked files as
missing. Do not delete `.woobe-config` or `.state/` as a generic repair.

Use only read-only diagnosis before choosing a repair. Do not copy full private
config/state, CLI Keys, provider secrets or sensitive runtime output into reports.
Confirm operation identity/revision/evidence instead of guessing from names.

Selecting a Draft recovers only private Draft observation. It does not recover
ordinary native push bindings or Provider credentials. Reconcile an unresolved
original write before switching selection. Do not use `create` as a binding repair
for a known same-origin root: clone with a new UID only for intentional duplication.


### Accepted native generations

After a successful native `push` or `reconcile`, the CLI checks every approved
resource against the original operation, definition and native identity before
updating local bases. When the server advertises accepted binding generations,
it uses the generation recorded in the native transaction, even if a later edit
has changed the current generation. Incomplete, missing or superseded evidence
leaves the checkpoint pending; preserve it and investigate with `reconcile`.
Do not delete private state to bypass this check or resend the accepted write.

For an uncertain `test --candidate`, use `draft reconcile` to recover the original Evaluation UUID. Then inspect `evaluation UUID` or `evaluation UUID reconcile --yes`. These commands never execute another case; do not retry the test to infer its outcome.

For exact Candidate publication, retain the original pending operation. An unknown acceptance or exit 10 after local receipt failure requires `woobe agent '@support' draft reconcile`; never repeat publish to repair its local mirror.

Network Candidate publication can remain accepted with partial immutable constituent copies. Do not repeat publish. Inspect `woobe network NETWORK_UUID publication PUBLICATION_UUID`, explicitly resume with `reconcile --yes` or cancel with `cancel --notes REASON --yes`, then run the original local `draft reconcile` to mirror terminal evidence. Cancellation preserves unused immutable child Releases and cannot undo published Networks. See [lifecycle](lifecycle.md).

Lease acceptance uses private pending operations. Read the original with `lease
reconcile` or `draft reconcile`; never repeat acquire after an uncertain response.
Expired/superseded proofs are intentionally retained and fenced by the server.
Explicit release or a deliberate new acquisition is required to change workflow
ownership; `lease show` never adopts a reservation.

For an unresolved exact Release selection, use `woobe agent '@support' deployment
reconcile` (or Network). It reads the recorded original operation and mirrors the
receipt before clearing pending state. Do not repeat apply. Native UUID recovery
requires the original `--operation UUID`; no project configuration is needed.


Managed lifecycle rejection on legacy Agent/Network promotion requires the
ASaC Candidate/evaluation/publication/deployment path. Do not change policy or
repeat legacy writes automatically. A confirmed stale policy generation needs
a fresh `woobe project lifecycle show` and an intentional new policy operation;
an uncertain policy write needs `woobe project lifecycle operation OPERATION_UUID`.

A completed ASaC operation seals its private pending checkpoint as an empty
object (`{}`). This is a cleared marker and does not block `bindings recover`.
A pending operation, partial record or malformed/unreadable checkpoint still
blocks recovery. Reconcile the original operation; do not delete an uncertain
checkpoint or manufacture a new operation ID to bypass it.

For server-signed export and offline verification with a separately pinned root,
read [signatures](signatures.md). Unsigned history integrity checks do not
establish offline server authority.

For `git-attest` uncertainty, inspect `woobe project agent git-attestation get RESOURCE_UUID ORIGINAL_OPERATION_UUID` (Network equivalent). Correct trust/source issues before another operation; never repeat deployment to repair Git evidence.
