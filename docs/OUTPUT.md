# Terminal and agent output

## Assistant guidance

`woobe skill install/status/uninstall` uses normal text/compact/JSON output;
installation destinations, version and drift evidence remain visible.

The embedded [assistant skill](AGENT_SKILL.md) selects compact output and exact
fields for inspection, and reads only the reference needed for the task. Full
JSON remains available for structured processing. It never hides incomplete
collection/write evidence or substitutes presentation for server authorization.

Woobe separates the information needed for the next action from the complete API
response. Presentation never changes the request, authorization, saved credential
or resource. The CLI keeps full IDs and every returned row.

## Choose a format

| Format | Behavior | Use |
| --- | --- | --- |
| `auto` (default) | Text on an interactive terminal; existing JSON envelope when redirected or captured | Terminal use and existing scripts |
| `text` / `table` | Relevant columns for lists, labeled fields for details, concise authentication and diagnostic results | People |
| `compact` | One JSON object with summarized `data`; `error` and nontrivial `meta` when present | AI agents and small machine responses |
| `json` | Existing complete versioned envelope, including context, response and metadata | Scripts needing the full contract |
| `jsonl` | Existing streaming protocol / envelopes | Runtime streams and Package observation |

`auto` detects the output destination, not the reader. An AI tool typically
captures output; set `--output compact` explicitly, or set `WOOBE_OUTPUT=compact`
in its environment. A command's `--output` flag overrides that environment value.
Pipes remain JSON unless you explicitly request another format.

```powershell
woobe agent list
woobe agent list --output compact
woobe agent list --output json

# Give an AI process concise responses for its subsequent CLI commands.
$env:WOOBE_OUTPUT = "compact"
# Restore automatic presentation when finished.
Remove-Item Env:WOOBE_OUTPUT
```

## Select information

The default Agent list shows name, status, model and ID. Networks show name,
status and ID. Authentication identifies the selected connection, Workspace and
Project without dumping every grant. `doctor` reports each check and route counts
without repeating the entire OpenAPI schema. Command discovery lists the command,
method, scope, effect and availability rather than all inherited flags per row.

Use `--fields` to select a comma-separated set of paths from the redacted response
data **before** the default summary is applied. Dot-separated paths refer to nested
objects; output keys retain those dotted names. A missing field becomes JSON `null`
or `-` in text. Field names are case-sensitive. `--fields` takes precedence over
`--wide`, and works with `text`, `table` and `compact`, not the complete JSON/JSONL
contracts. The choice is local presentation; it does not filter the API query.

```powershell
woobe agent list --output text --fields name,status
woobe agent list --output compact --fields id,name
woobe agent get AGENT_ID --output text --fields name,description,system_prompt
woobe auth status --output compact --fields context,identity.workspace.name,project_id
woobe agent list --wide
```

`--wide` retains all response fields in a concise format, with envelope wrappers
removed. `--output json` keeps the complete, existing envelope. Indented terminal
schemas retain their definitions; plans, manifests and Package recovery evidence
also keep their details by default. Unknown response shapes fall back to their
data rather than producing an empty summary.

## Failures and pagination

Compact failures include `error.exit_code`, message and available HTTP status,
domain code, request ID, diagnostics and write outcome. Exit codes are unchanged.
An uncertain mutation remains `write_outcome: unknown`; a short display never
makes it safe to repeat a write. Partial results remain visible alongside the
error, with `meta.complete: false` and continuation/collection evidence.

Pagination markers and collection completeness remain available in concise
output. A successful command does not establish that an entire collection was
retrieved. `--fields` cannot remove error or pagination metadata. Runtime event
streams continue to require `--output jsonl`.

## Give an AI agent a short instruction

```text
Use the saved Woobe connection. Request --output compact for normal commands.
For lists, request only the fields needed for the next action, such as
--fields id,name,status. Use --wide or --output json when full details are needed.
Preserve exit codes, write_outcome and pagination metadata. Do not repeat an
uncertain write. Consult the specific command/schema instead of all discovery.
```

See [generated before/after comparisons](OUTPUT_EXAMPLES.md). Regenerate them
from the repository root with:

```sh
go build -o bin/woobe ./cmd/woobe
python3 scripts/output_examples.py --binary bin/woobe
```

The [complete command audit](OUTPUT_AUDIT.md) lists every discovered operation,
with sizes and scenario limitations; the [CSV](OUTPUT_AUDIT.csv) is sortable.
HTTP measurements use a loopback presentation fixture and do not validate real
backend writes or authoritative DTOs. Entries without a dedicated valid scenario
are marked as not measured; scripts, aliases and JSONL streams are classified
separately. Schemas and recovery data intentionally retain their information.

```sh
python3 scripts/output_audit.py --binary bin/woobe
```

Resource IDs do not shadow global configuration/credential flags. For example,
`project agent model-config get` accepts a positional configuration ID or
`--config-id`; `--config` continues to select the local CLI configuration file.
