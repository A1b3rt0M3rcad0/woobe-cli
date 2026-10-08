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
