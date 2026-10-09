# Other operations and economical inspection

Use targeted help/schema discovery, not recursive help scans:

```sh
woobe help "project tool create" --output compact
woobe schema --command "project tool create" --kind input
woobe server-schema --command "project tool create"
woobe validate-input --command "project tool create" --file ./tool.yaml
```

Local invocation schemas describe CLI input/transport; server schemas describe
the actual DTO. Do not assume the portable Agent YAML is a raw Agent PATCH body.
Administrative CRUD accepts `--file` YAML or JSON. Omitted fields and explicit
null have different meanings. Forward exact supported `--if-match` evidence;
never infer a safe overwrite without an advertised precondition.

Workspace/Project authority, provider credentials, tools, Knowledge, runtime
Skills and Surfaces remain canonical CLI operations. Choose the owning domain's
exact command using targeted help and validate the supplied body. Do not invent
a route or overwrite configuration by passing an entire read response as PATCH.

For small inspections use `--output compact --fields "id,name,status"`. CLI
output defaults to text in a terminal and JSON in pipes; automation should choose
explicit formatting. Use JSONL for streaming. Inspect `meta.complete`; a single
page is not a complete collection. Use bounded `--all` only on supported list
commands when the task needs every result, preserving filters and deadlines.

Runtime is a separate requested execution, not configuration validation:

```sh
woobe runtime target run TARGET --target-kind agent --runtime-credential runtime --file ./run.yaml --yes
woobe runtime target observe TARGET RUN --target-kind agent --runtime-credential runtime --output jsonl
```

A Runtime Key does not substitute for an administrative key. Do not replay a Run
to recover a lost response. Runtime input/options use the pinned SDK contract;
consult exact help/schema before composing them.

Manifest supports explicit multi-resource orchestration with typed dependencies:
validate → compile/diff → plan/apply as applicable. Read only the needed schema,
use exact references and keep checkpoints. A partial projection or GET observation
does not prove acceptance of an earlier uncertain create.

A positive integer `fencing_token` is a concurrency counter and remains visible. Treat workflow lease state as private operational state; access tokens/keys and textual token values remain redacted. Do not copy `.state` into Git.
