# Tests and native lifecycle

Lifecycle operations use saved server snapshots; they do not upload local YAML.
Push and inspect Draft changes first. For a specifically authorized Agent flow:

```sh
woobe agent '@support' stage --yes
woobe agent '@support' test --file ./support-test.yaml --yes
woobe agent '@support' publish --notes "Reviewed support changes" --yes
woobe agent '@support' release list --output compact
woobe agent '@support' activate --version ACTUAL_RELEASE_VERSION --notes "Approved deployment" --yes
```

Stage copies current Draft to Staging. The test executes saved Staging by default,
records a real Run/result and exits 6 on failure. The bundled
[testcase](../assets/agent-test.yaml) shows `message`, optional `expected_output`
and `external_context`; adapt it to the requested acceptance criteria. For Draft
tests use `--env draft`; a Release test needs `--env release --version VERSION`.
Passing local YAML validation is not a successful execution test.

Agent publish creates an immutable Release from Staging; activate selects its
actual returned version. Do not fabricate a version or assume publish activated
the Agent. Rollback requires an explicit older Release, audit reason and `--yes`.
Archive/delete follows native dependency checks and requires authorization;
local `resources unregister` does not delete the remote Agent.

For Networks, `stage` freezes constituent snapshots. **Network publish also
activates its Release in Production**, following Woobe's native behavior. Obtain
authorization for that Production effect before adding `--yes`; do not present
it as a harmless Release-only action. Network constituent Agents are not
automatically activated as standalone Production Agents. Network activate/
rollback selects an explicit immutable version with `--notes` and `--yes`.

Native UUIDs work without `.woobe-config`; aliases/paths need connection-scoped
bindings. Denied operations do not become allowed by supplying `--yes`.

## Exact isolated revision preparation (CLI 0.24+)

For an isolated historical Draft, use `stage --revision REVISION_ID --yes` after
`draft push`, `revision create --message ...` and `draft checkpoint REVISION_ID`.
Inspect first with `--dry-run`. This uses the verified retained executable object
and exact remote Draft generation; later YAML edits are not included. Provider
bindings must already exist. Network and Agent syntax is identical.

Acceptance returns `preparing` and a Candidate UUID. Read
`woobe agent '@support' candidate CANDIDATE_UUID` (or Network) to inspect progress;
use the preparation operation ID with `woobe package status OPERATION_UUID` for
errors/dependencies. Ready means the detached closure is frozen, not that tests
passed or publication occurred. Do not feed it into legacy publish/activate.

If acceptance is uncertain, use `draft reconcile` before another Stage; preserve
its original ID. Stage does not save current YAML, select native Staging, publish
or activate Production. Without `--revision`, Stage retains the legacy behavior
documented above. Never silently fall back when Candidate support is absent.

```sh
woobe agent '@support' stage --revision rv_00000000-0000-0000-0000-000000000000 --dry-run
woobe network '@customer-support' stage --revision rv_00000000-0000-0000-0000-000000000000 --dry-run
```

Replace the revision placeholders with actual registered checkpoints before use.

## Exact Candidate evaluation

After an isolated `stage --revision`, inspect the Candidate until its state is
`ready`. Evaluate its UUID with a versioned YAML suite; current local YAML and
native environment selections cannot replace this frozen snapshot.

```yaml
schema_version: '1.0'
suite_id: support-smoke
suite_version: '1'
dataset_version: '1'
policy_version: all-cases-pass@1
cases:
  - id: greeting
    message: Hello
    # Optional: expected_output is compared to actual structured output,
    # or {answer: ...} for unstructured text. Omit it for execution/schema checks.
```

```sh
woobe agent '@support' test --candidate CANDIDATE_UUID --file suite.yaml --yes
woobe network '@helpdesk' test --candidate CANDIDATE_UUID --file suite.yaml --yes
woobe agent '@support' evaluation EVALUATION_UUID
woobe network '@helpdesk' evaluation EVALUATION_UUID reconcile --yes
```

The suite accepts 1–16 uniquely named cases. Suite, dataset and policy versions
are declared identities; the backend separately verifies content digests. Every
case gets its original native Run identity before execution. Qualification uses
actual terminal Run evidence, exact snapshot identity, output validation and
current binding fingerprints. Missing evidence stays incomplete. Failed cases
exit with code 6; `accepted` and `running` do not mean passed.

The immutable acceptance receipt identifies the evaluation. Inspect its ledger
for the current result. If acceptance is uncertain, use `draft reconcile` before
another write: it recovers the original receipt and never executes a case again.
The evaluation's own `reconcile` reads original Run evidence without provider
calls. Neither command publishes nor changes Production. Do not combine
`test --candidate` with `--env` or `--version`. Native Agent `test --env` remains
available; Network testing requires an exact Candidate. Inspection by native UUID
works without `.woobe-config`; accepting an evaluation uses private project state
to preserve its write identity. Exact Candidate evaluation requires CLI 0.25+
and a backend advertising `asac.candidate_evaluation`.

## Publish an exact evaluated Agent Candidate

```sh
woobe agent '@support' publish --candidate CANDIDATE_UUID --evaluation EVALUATION_UUID --notes 'Validated support update' --yes
woobe agent AGENT_UUID publication PUBLICATION_UUID
woobe agent '@support' history verify
```

This explicit flow requires CLI 0.26+ and a backend advertising Agent in
`asac.candidate_publication_kinds`. It requires a ready Candidate and a passed
Evaluation for the same retained revision, runtime and bindings. The server
rechecks those bindings and publishes the frozen Candidate, even if current
Staging subsequently changed. It never pushes local YAML or activates Production.
Native Release reuse still produces a separate publication provenance receipt.

Acceptance requires the registered development resource and its private state;
inspection by native UUID does not require `.woobe-config`. The CLI persists the
original operation before POST and mirrors its immutable receipt after acceptance.
An unknown remote outcome blocks another write. A local documentary failure after
remote commit exits with code 10 and retains the pending operation. In both cases,
use `woobe agent '@support' draft reconcile`: it reads the original operation and
retries the local mirror, without another publication request.

`history verify` checks receipt integrity and destination identity. Its
`publication_authenticity: integrity_checked_connection_unverified` distinguishes
an authenticated fetch observation from a signed offline attestation. Git claims
remain declared; this flow does not claim verified pipeline provenance.
Legacy publication without `--candidate` retains its native behavior; in
particular, legacy Network publication activates Production. Exact Network publication uses the explicit Candidate flow below.


## Publish and recover an exact Network Candidate

```sh
woobe network '@support-network' publish --candidate CANDIDATE_UUID --evaluation EVALUATION_UUID --notes 'Validated Network update' --yes
woobe network NETWORK_UUID publication PUBLICATION_UUID
woobe network NETWORK_UUID publication PUBLICATION_UUID reconcile --yes
woobe network NETWORK_UUID publication PUBLICATION_UUID cancel --notes 'Discard unfinished preparation' --yes
woobe network '@support-network' draft reconcile
woobe network '@support-network' history verify
```

Use CLI 0.27+ with a backend advertising Network Candidate publication. Prepare
and evaluate a new Candidate with `network-execution-runtime@2`; historical
Candidates retain their original runtime scope. Publication copies or reuses
immutable constituent Agent Releases through their owning module and creates the
Network Release from that exact tested composition. Production, current native
Draft/Staging and standalone Agent environment selections remain unchanged.

The server accepts the original publication before constituent copying. A failure
can leave `state: preparing` with prepared/required counts. The CLI exits 9 and
retains its pending operation, blocking a second publish. Inspect the returned
Publication UUID, then explicitly resume with `publication ... reconcile --yes`.
The server rechecks current permissions, original Evaluation and binding versions;
resume reuses original child-operation receipts rather than copying twice.

Alternatively, cancel an unfinished preparation with a reason. Cancellation seals
it without publishing the Network; already copied immutable Agent Releases remain
retained. Cancellation cannot undo a completed publication. These native UUID
recovery operations do not require `.woobe-config`.

After either terminal outcome, use `draft reconcile` for the original local
operation. It only reads original acceptance/current phase, mirrors the publication
or cancellation receipt, and clears pending state after the mirror succeeds.
It never repeats publication, constituent copying or activation. Documentary
failure exits 10 and remains recoverable; `history verify` checks both receipt
collections. Offline verification proves content integrity, not server authenticity.
