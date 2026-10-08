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
