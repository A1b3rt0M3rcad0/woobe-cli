# CLI workflows

## Build and discovery

```sh
go build -o bin/woobe ./cmd/woobe
bin/woobe help --output json
bin/woobe schema --command "project agent create" --kind input
bin/woobe server-schema --command "project agent create" --api-url https://woobe.example
bin/woobe completion bash
```

Discovery covers all executable handlers and flags; local schemas describe invocation and transport/runtime output. Use the server schema to discover the actual DTO. Domains remain validated and authorized by the server.

## Context and credentials

```sh
woobe context create minha-woobe --api-url https://YOUR_WOOBE_API
woobe context use minha-woobe
woobe auth login --cli-key
woobe project agent list
```

Paste the key at the masked prompt; do not put it in arguments. The CLI validates
it against this connection, discovers its Workspace and selects a single eligible
Project automatically. For several Projects, choose by name/number. Switch later
with `woobe context project select` or `woobe context project select production`.
No Project grants means Workspace-only authentication. No UUIDs or environment
variables are required. Linux/macOS store credentials in private files; Windows
uses the current user's Credential Manager. Reopening the terminal retains the
connection. Each named context has its own key and selected Project.

For automation, pipe a secret from your secret manager into
`woobe auth login --cli-key --stdin --select-project production --no-input`.
Without a selection, multiple Projects produce `project_selection_required`.
`woobe auth status` inspects the current key and scope;
`woobe auth logout --cli-key` clears local authentication without server revocation.
Changing the API URL requires login again. Runtime Keys remain separate.

Precedence: explicit flag, corresponding `WOOBE_*` environment variable, selected context, documented default. `WOOBE_CONFIG` selects the config file. Scope variables are `WOOBE_WORKSPACE_ID` and `WOOBE_PROJECT_ID`. Credentials may be imported from stdin on all supported systems; environment fallbacks are `WOOBE_CONTROL_KEY` and `WOOBE_RUNTIME_KEY`. Runtime never falls back to an administrative key. Secret values have no command-line flag.

Select a different Workspace and the previously selected Project is cleared by `context set`. Context references are local preferences; the server resolves and verifies resource ownership. Direct known resource IDs can be read without first listing every Project.

Human login accepts a JSON file or stdin containing the server's login fields:

```sh
woobe auth login --file - < login.json
woobe auth refresh
woobe auth status
woobe auth logout
```

## Configuration and publication

```sh
woobe project agent create --file agent.json
woobe project agent update AGENT --file patch.json --if-match ETAG
woobe project agent prompt create --agent AGENT --file prompt.json
woobe project agent contract create --agent AGENT --file contract.json
woobe project agent model-config create --agent AGENT --file model.json
woobe project network create --file network-metadata.json
woobe project network draft update NETWORK --file network-definition.json
woobe project agent release create --agent AGENT --file release.json
woobe project agent release promote RELEASE --agent AGENT --file promotion.json --yes
```

IDs are positional in route order, or explicit parent flags. PATCH bytes are preserved: omitted fields stay omitted; explicit null remains null. `--if-match` forwards a server ETag; it does not implement revision protection when the server ignores that header.

`--dry-run` renders a redacted request without sending it. `--yes` acknowledges the specified execution/publication/destructive action; it never changes permissions or selects another credential. No interactive prompt is used.

```sh
woobe project api-key create --file key.json --secret-file ./issued-key.json
woobe project knowledge document upload --collection COLLECTION --document-file lesson.pdf
woobe project agent export AGENT --destination agent-projection.json
```

Exports are explicitly incomplete authorized projections, not apply-ready patches. API issuance destinations must not exist; a failed attempt can leave an empty file. Reconcile before deciding whether another issuance is necessary.

## Runtime

```sh
woobe runtime target run tutor --target-kind agent --runtime-credential runtime --file run.json --yes
woobe runtime target stream tutor --target-kind agent --runtime-credential runtime --file run.json --yes --output jsonl --timeout 10m
woobe runtime target observe tutor RUN --target-kind agent --runtime-credential runtime --output jsonl
woobe runtime target active tutor SESSION --runtime-credential runtime
woobe runtime target cancel tutor RUN --runtime-credential runtime --yes
```

Run JSON: `{"input":"Explain this topic","options":{"SessionID":"..."}}`. Omit options to let the runtime establish a Session. Option names currently follow the pinned SDK's Go JSON field names. For Network targets select `--target-kind network` explicitly. Local interruption stops observation; it does not issue remote cancellation. Administrative Runs use the separate `runtime run` family. All five SDK handlers honor `--dry-run` without network execution.

## Manifest validation, plan and checkpoint

```sh
woobe manifest validate --file testdata/agent-network.manifest.json
woobe manifest plan --file testdata/agent-network.manifest.json
woobe manifest apply --file resources.json --yes --checkpoint ./apply-checkpoint.json
```

Document fields: `schema_version`, `workspace_id`, `project_id`, `steps`. Each step declares `id`, canonical `command`, positional `args`, JSON `body`, `depends_on`, optional `if_match`. Declare dependencies explicitly. Exact references such as `${steps.agent.id}` resolve persisted redacted results; interpolation inside larger strings is unsupported.

Match manifest Workspace/Project to the selected context. Preserve the same manifest and checkpoint to resume. Committed and reconciled steps are skipped; `unknown`/`in_flight` blocks execution. Never delete the checkpoint and blindly replay creation after a lost response.

The explicit-step format composes existing server operations. It is not the complete semantic resource reconciler specified in the design. `manifest diff --file resources.json` reads compatible update resources and compares supplied top-level fields. Creation existence and unresolved dependencies remain unevaluated.

For uncertain writes: `woobe manifest reconcile STEP --resource-id ID --file resources.json --checkpoint ./apply-checkpoint.json --yes`. The same manifest, origin, workspace, project and credential reference are required. An authorized GET must match nonempty supplied desired fields. Evidence proves observed state, not original-write attribution. Resume apply using the same checkpoint. Concurrent use is rejected; remove a stale `.lock` only after confirming no process uses it.

## Output and advertised extensions

The default `--output auto` shows concise text in an interactive terminal and
preserves the existing JSON envelope in pipes. `--output text` or `table` selects
relevant rows/fields. AI agents can use `--output compact`, optionally with
`--fields id,name,status`. `--wide` shows all data fields; `--output json` retains
the complete envelope. Error and collection evidence remain visible. See
[format contracts and examples](OUTPUT.md).

Proposed extension endpoints require the exact route/method in `/openapi.json`;
actual requests still require server authorization.

## Exit codes

| Code | Meaning |
| --- | --- |
| 0 | Success |
| 1 | Internal client error |
| 2 | Invalid usage or input |
| 3 | Authentication |
| 4 | Authorization denied by server |
| 5 | Not found/inaccessible |
| 6 | Revision/state conflict |
| 7 | Connection/uncertain write |
| 8 | Deadline |
| 9 | Unsupported capability/protocol |
| 10 | Partial composed result |
| 11 | Remote Run failure/cancellation |
| 130 | Local interruption |

JSONL runtime stdout contains semantic SDK events; client errors go to stderr. Other JSON output is one envelope. Credentials and administrative secret fields are redacted; semantic runtime event payloads are preserved.

## Input validation and checkpoint inspection

```sh
woobe schema --command "manifest validate" --kind document
woobe validate-input --command "project agent create" --file agent.json
woobe manifest status --checkpoint ./apply-checkpoint.json
woobe doctor
```

`validate-input` sends only an OpenAPI GET; it does not submit the supplied body. It supports object/array/primitive types, required/properties/additionalProperties, sizes, enum/const, numeric bounds with exact rationals, RE2-compatible patterns, local references and composition. Formats `uuid` (canonical hyphenated hex, case insensitive), `date` (valid full calendar date) and `date-time` (RFC 3339, required offset, fractions and leap-second minute checks) are validated. Other formats, unsupported assertions, external references and excessive schema depth return exit 9. Advertised OpenAPI versions must be 3.1.0–3.1.2; unknown or 3.0 dialects are refused instead of interpreted as 3.1. An absent version retains the existing explicit subset. Legacy `nullable` is refused in explicitly advertised 2020-12/OpenAPI 3.1 schemas, including unused properties; use a type union or anyOf with null. `$schema` and `jsonSchemaDialect`, when provided, must select JSON Schema 2020-12 or the OpenAPI 3.1 base dialect.

Body validation applies **request direction**: a supplied `readOnly` field returns exit 2, including fields reached through references, `allOf`, nested objects and array items. Required read-only properties in the same object schema can be omitted; annotation lookup follows local references and `allOf` under that property's schema. `writeOnly` fields remain permitted and retain required/type rules. Contradictory annotations and direction annotations inside predicates (`not`, `if`, `contains`, `propertyNames`) return exit 9. Arbitrary cross-branch required/annotation projection is not inferred. Neutral internal schema checks retain their required semantics.

The same checks run in `validate-input`, `manifest preflight`, opt-in HTTP writes and opt-in manifest apply, after dependency values resolve. Invalid inputs stop before mutation and preserve `not_attempted` recovery evidence. Success reports request direction and explicitly leaves path/query validation unevaluated. The shared body validator also rejects duplicate fields, multiple JSON values and excessive input depth. Domain validation and authorization still occur on the server. Validation is explicit; writes do not silently add schema-fetch requirements.

JSON inputs reject duplicate fields, more than one value and nesting beyond 128. Manifests require object configuration bodies, exclude read steps and literal sensitive fields (including credential `value`), and reject duplicate/self dependencies. IDs match `[a-zA-Z0-9_-]{1,128}`. Canonical hashing ignores object-key order/formatting and preserves array order, null and numeric values. Updated hash/fingerprint rules can reject old checkpoints; preserve them for manual reconciliation instead of deleting/replaying.

Control Key fingerprints bind checkpoint recovery to actual key material, including environment keys. Human-session principal binding still needs authoritative server identity. `doctor` runs instance, identity and OpenAPI reads independently; partial diagnostics include evidence and exit 10. Advertised routes never imply effective authority.

## Resource projections and distribution

```sh
woobe export NETWORK --command "project network get" --project PROJECT --destination network.json
woobe export TRACE --command "project trace get" --destination trace.json
woobe manifest export --resource agent --id AGENT --destination agent.json
```

Generic export accepts canonical resource `get` operations with exact route-order IDs. Manifest shortcuts cover agent, network, skill-version, knowledge collection and surface. Every projection remains `complete:false` and `apply_ready:false`, with observed ETag/request ID when available. Unsupported Tool getters are not fabricated. Secret fields are redacted before file/output rendering. Projection export does not guarantee pagination or round-trip apply.

Packages include the executable, README, usage and `manifest.schema.json`. `artifacts.json` records version, source commit, compiler, OS/architecture, byte size and SHA256 for exactly six archives. `python3 scripts/verify_artifacts.py` verifies contents, schemas and checksums before release publication. No release/tag is created by this PR.

## Resource manifests — schema v2

```sh
woobe manifest kinds
woobe schema --command "manifest validate" --kind document --manifest-version 2
woobe manifest validate --file testdata/resources.manifest.json
woobe manifest compile --file testdata/resources.manifest.json
woobe manifest plan --file testdata/resources.manifest.json
woobe manifest apply --file resources.json --project PROJECT --checkpoint ./resources.checkpoint.json --yes
woobe manifest apply --file updates.json --project PROJECT --checkpoint ./updates.checkpoint.json --yes --skip-unchanged
```

The v2 document requires `schema_version:"2"`, `project_id` and `resources`; `workspace_id` is optional. Each resource has `key`, `kind`, explicit `action` and object `spec`, with optional `resource_id`, `parents`, `depends_on` and `if_match`. `create` forbids remote ID/ETag; `update` requires an explicit remote ID. The client never guesses identity by name or implicitly creates a missing update target.

| Kind | Supported actions | Parent / target |
| --- | --- | --- |
| Agent | create, update | update target is Agent ID; creation receives document project_id |
| Network | create, update | update target is Network ID |
| AgentPrompt | create | parents.agent |
| AgentContract | create, update | parents.agent; update target is Contract ID |
| AgentModelConfig | create | parents.agent |
| NetworkDraft | update | resource_id is its Network ID |

Whole-value references use `${resources.KEY.id}` for parent/target IDs. Referenced resource kinds must match their consumer (e.g. AgentPrompt parent must reference Agent, NetworkDraft target must reference Network). References in `spec` may address returned fields. Dependencies remain explicitly declared. Embedded interpolation and the steps namespace are rejected in v2. Actual spec DTOs/bindings still follow server contracts; the example demonstrates client composition, not a certified Network draft DTO.

Offline compilation emits a v1 execution document and hash. Both formats reuse the same diff/apply/checkpoint/reconcile machinery; semantically identical compiled plans have the same execution identity. Keep context and original resource intent stable during recovery. Explicit-ID uncertain-create reconciliation also works for v2 without replaying the creation.

`--skip-unchanged` performs an authorized GET for compatible update targets, comparing only supplied fields. Matching state is checkpointed as `unchanged`, making results available to dependent resources; resume skips that observation just as it skips completed writes. It is a recorded observation, not continuous drift detection. Changed state proceeds to the ordinary PATCH/PUT; server-side revision enforcement remains authoritative. An explicit `if_match` requires an observed matching ETag before this shortcut. Missing compatible getters (including NetworkDraft) return exit 9; default apply remains available without the option. Create intents are always explicit creations and have no name-based existence check.

Packages contain both `manifest.schema.json` (steps v1) and `resources.schema.json` (resources v2). Manifest resource coverage and semantic upsert/existence reconciliation remain limited. Portable Agent/Network YAML export/import and protected bindings use the separate Package V1 workflow below.

## Captura seletiva e paginação anunciada

```sh
woobe manifest capture --kind Agent --id AGENT --project PROJECT --field name --field description --destination update.json
woobe manifest capture --kind AgentContract --id CONTRACT --parent agent=AGENT --project PROJECT --field name --destination contract-update.json
woobe manifest apply --file update.json --project PROJECT --checkpoint update.checkpoint.json --yes
woobe request-pages /ai/agents --query project_id=PROJECT --max-pages 20
```

`capture` exige um tipo com update/get compatíveis, Project explícito e campos de primeiro nível escolhidos. Somente esses campos entram em `spec`; nulos explícitos são preservados. Identidade/contexto, segredos, campos ausentes/inacessíveis e projeções de outro ID/Project são recusados. O ETag observado entra em `if_match` quando disponível, mas a execução da condição continua sendo responsabilidade do servidor. O arquivo de destino contém um manifesto v2 bruto; stdout também informa seleção/proveniência. Isso permite editar e reaplicar a configuração escolhida; não promete exportação completa nem validação integral dos DTOs. Tool e NetworkDraft não possuem getter compatível no catálogo atual; tipos imutáveis não têm capture.

O catálogo v2 agora contém 13 tipos: Agent, Network, AgentPrompt, AgentContract, AgentModelConfig, NetworkDraft, Project, Tool, KnowledgeCollection, KnowledgeDocument, Skill, SkillVersion e ChatSurface. Project aceita apenas update e o ID deve ser o Project do documento. SkillVersion exige parent `skill` e referências de IDs do tipo Skill. As demais ações e campos seguem `manifest kinds` e o schema distribuído; não há emissão declarativa de credenciais ou categorias propostas.

`request-pages` e os comandos paginados suportam contratos revisados de cursor e histórico, além de Link. Use `--all`, `--limit`, `--cursor` ou `--before-revision` conforme o comando. A saída informa coleção inteira, páginas restantes, parcial ou não verificada; falhas preservam páginas recebidas. Consulte [PAGINATION.md](PAGINATION.md) para protocolos, limites, metadados e exemplos.

Falhas de apply agora incluem checkpoint, contagens por estado, etapa interrompida e `checkpoint_saved`. A ausência de persistência é explícita; uma resposta de sucesso não autoriza repetir uma criação. Binários empacotados incluem a revisão fonte em `version`, junto de compilador/OS/arquitetura.

## Medida de completude

`python3 scripts/completeness.py` valida a correspondência com todas as 101 entregas do §18 e verifica o relatório gerado. `docs/COMPLETENESS.md` documenta estado/evidência por item e a fórmula. Esse percentual inclui backend e E2E; não deriva de commits, número de comandos ou cobertura de código.

## Context integrity and extended schema validation

`context show NAME` resolves the selected profile. Explicit empty `--project`, `--workspace`, `--credential` and `--runtime-credential` override inherited values for that invocation. Changing Workspace clears an inherited Project unless a Project is supplied explicitly. Persist clearing with `context unset NAME workspace`, `context unset NAME project`, or `context unset NAME credential runtime-credential`; stored secrets are preserved. Attach a separate runtime reference with `context runtime-credential attach NAME --runtime-credential REF` and detach it with `context runtime-credential detach NAME`.

Config v1 is bounded to 1 MiB, rejects unknown/duplicate fields, dangling active context names and invalid context URLs/names, and validates before atomic replacement. See [MIGRATION.md](MIGRATION.md) for Python prototype command/credential/output migration.

The explicit `validate-input` subset additionally supports min/maxProperties, dependentRequired, dependentSchemas, propertyNames, RE2 patternProperties, if/then/else, prefixItems and contains with min/maxContains. Numeric enum/const and uniqueness use exact mathematical equality; nullable alters type acceptance without bypassing enum/composition constraints. Local JSON pointers support escaped property names, URI fragment decoding and array indices. Schema inspection includes inactive properties/branches; an unsupported rule there returns exit 9 rather than valid=true.

Limits: schema depth 64, input depth 128, shared work budget 100,000 nodes/evaluations, numeric representation at most 4,096 characters and exponent magnitude at most 4,096. Inputs exceeding evaluator bounds return exit 9. Format, external references/anchors, schema resource IDs, unevaluated-property/item semantics and unsupported assertion dialects remain outside the subset. readOnly/writeOnly and documentation keywords are annotations here, not request-direction policy. Server authorization, domain validation and create/PATCH rules remain authoritative. No automatic schema probe is added to writes.

## Body preflight, opt-in write validation and recovery integrity

```bash
woobe manifest preflight --file project.json --project PROJECT
woobe manifest preflight --file project.json --project PROJECT --require-complete
woobe project agent create --file agent.json --validate-body
woobe manifest apply --file project.json --project PROJECT --checkpoint apply.json --yes --validate-body
```

`manifest preflight` reads one OpenAPI snapshot and checks each configuration body against its advertised operation. It never executes the manifest. Body values depending on returned step fields are deferred; operation availability and supported schema declarations are still inspected. `data.complete` and `meta.complete` are false when any body is deferred or fails. `--require-complete` returns exit 9 for deferrals or a dry-run without evaluation. Invalid inputs return exit 2, unsupported contracts exit 9, with per-step evidence. Path/query schemas, authorization, resource existence and domain rules are not evaluated by body preflight.

`--validate-body` is optional and available for canonical JSON HTTP mutations and manifest apply. The default write path preserves its previous behavior. Generic request, multipart upload, SDK runtime, reads and unrelated local commands reject this flag rather than ignoring it. Validation failure marks the selected write `not_attempted`. Secret emission reserves its destination only after validation succeeds. Dry-run never fetches a schema or sends a write.

For manifest apply, one OpenAPI snapshot is fetched after local/context/checkpoint checks. Each pending body's references are resolved, then validation runs before unchanged-state observation or `in_flight` checkpointing. A failure records a resumable `not_attempted` state with a structured cause and partial report (exit 10), preserving earlier committed steps. Unknown/in-flight HTTP attempts still block resume until explicit reconciliation. Validation of a future dependent body cannot be guaranteed before its required values exist.

Validation reports include `schema_sha256`: SHA-256 of the parsed OpenAPI document serialized as sorted-key JSON with preserved numbers. Supply that digest with `--schema-sha256 DIGEST` to `validate-input`, manifest preflight, or a write/apply using `--validate-body`. A different snapshot returns exit 6 before the selected write. This identifies a client snapshot; it is not a server revision lock and does not freeze the server between discovery and mutation. The explicit supported subset and all documented schema limits still apply.

Checkpoint reads are bounded to 32 MiB. Resume/reconcile validates step IDs, result ownership, reconciliation state, saved results for terminal steps and completed dependencies against the same manifest. Returned numbers are preserved exactly through checkpointing and body references. Corrupt/incompatible checkpoints are refused; status remains an offline inspection tool. `not_attempted` can be retried by apply after correcting the prerequisite without replaying committed steps.

HTTP responses reject duplicate fields/nesting ambiguity. Canonical resource IDs must be unambiguous single segments; decoded path traversal, separators, nested escapes and control characters are refused. A canonical scope query must occur once and match the selected Workspace/Project. Generic authorized HTTP still delegates resource authorization to the server.

The default administrative transport pools read connections but sends mutations over fresh HTTP/1 connections with replay disabled, including requests carrying Idempotency-Key. This trades write connection reuse for predictable single-attempt behavior; idempotency remains a server feature, not a retry instruction. Runtime SDK calls inherit this HTTP client. Replacing the HTTP transport in an embedding application requires preserving this guarantee explicitly. Checkpoint writes also enforce 32 MiB before replacing an existing file; a committed remote write followed by a failed checkpoint save remains a partial outcome requiring recovery.

State comparisons in `manifest diff`, `apply --skip-unchanged` and `reconcile` use exact decimal JSON number semantics: `1`, `1.0` and `1e0` are equal, including inside supplied objects/arrays. Large integers remain distinct without conversion to floating point. Array order, missing fields, explicit nulls and string/number types remain significant. Nested supplied objects are compared in full; this is not a recursive partial PATCH interpretation. Comparison shares a 100,000-node budget across supplied fields, supports depth 128 and numbers of at most 4096 characters with exponent magnitude at most 4096. Unsupported comparisons return exit 9; they never justify an unchanged checkpoint or reconciliation. Expected `if_match` in diff also requires an observed matching ETag. These observations do not lock subsequent server state or establish original-write attribution.

## Advertised path/query validation

`--validate-parameters` validates the selected canonical HTTP operation's path
and query values against its advertised OpenAPI 3.1 schema before the operation.
It can be combined with `--validate-body` for writes; neither changes authority.
Dry-run remains network-free and performs neither advertised validation.

```sh
woobe project agent update AGENT_UUID --file agent.json --validate-body --validate-parameters
woobe runtime agent sessions AGENT_UUID --project PROJECT_UUID --query environment=staging --limit 1 --all --validate-parameters
woobe validate-input --command 'project agent get' --validate-parameters --path-param agent_id=AGENT_UUID
woobe manifest preflight --file resources.json --validate-parameters --require-complete
woobe manifest apply --file resources.json --checkpoint apply.json --yes --validate-body --validate-parameters
```

Parameter-only `validate-input` does not require a body and explicitly reports
`body_validation: not_evaluated`. The parameter subset supports strings,
UUID/date/date-time formats, exact JSON numbers, booleans and form arrays
(repeated values for explode=true; one comma-separated value for explode=false).
Path-level parameters are replaced by operation-level declarations with the
same name/location. Local references and unambiguous nullable/composed primitive
types are supported. Canonical/server template names may differ when method and
literal path segments match exactly and identify only one advertised route.

Missing required parameters, duplicate scalar values and invalid values return
exit 2. Unadvertised parameters, unsupported serialization/types/dialects and
ambiguous templates return exit 9. Header/cookie schemas and authorization are
not evaluated by this flag. Schema digests can pin the discovery snapshot.

Manifest apply validates resolved parameters before recording an in-flight
write; failures retain `not_attempted` checkpoints. Preflight defers parameters
that depend on prior step results and records parameter/body evidence separately.
Manifest parameter validation uses each operation's context-derived query scope;
global `--query` overrides are refused because they are not propagated to steps.
Unknown server query keys are never approved by guessing their semantics.


## Workspace category definitions

`AuthorityCategory` is a Workspace-owned resource kind. A category-only resource
manifest requires `workspace_id`; `project_id` is optional. Mixed documents still
require the owning scope of every kind. Creating a definition does not issue keys
or assign authority. Editing a definition requires an explicit ID and strong
`if_match`; it does not migrate existing grants to the new revision.

```json
{
  "schema_version": "2",
  "workspace_id": "WORKSPACE_UUID",
  "resources": [{
    "key": "reader", "kind": "AuthorityCategory", "action": "create",
    "spec": {"name": "custom-reader", "scope": "project", "permissions": ["agent:read"]}
  }]
}
```

```sh
woobe manifest preflight --workspace WORKSPACE_UUID --file category.json --validate-parameters --require-complete
woobe manifest apply --workspace WORKSPACE_UUID --file category.json --checkpoint category.checkpoint.json --validate-body --validate-parameters --yes
woobe workspace authority category diff CATEGORY_UUID --workspace WORKSPACE_UUID --from-revision 1 --to-revision 2
woobe manifest capture --kind AuthorityCategory --id CATEGORY_UUID --workspace WORKSPACE_UUID --field name --field permissions --destination category-edit.json
```

The diff reads exactly the two selected immutable definitions and reports supplied
configuration fields. It does not calculate every affected assignment or apply a
migration. `--dry-run` reports that comparison was not evaluated and performs no
network requests. Permission/condition order does not create a change; an absent
restriction and an empty restrictive set remain distinct. Unknown condition
semantics are unsupported rather than silently compared as understood authority.

Capture copies only selected fields and the observed ETag. Category specs exclude
assignment, grant, identity, status and revision fields; scope is immutable during
update. Server OpenAPI advertisement is checked before category writes. A saved
successful checkpoint resumes without another create; a lost checkpoint does not
make category creation idempotent or infer existence by name. Full uncertain-write
reconciliation and applied-category grant migration remain incomplete.

## Protected manifest credentials

Sensitive values may use an exact object `{"$secret_ref":"NAME"}` in Tool or
provider-credential configuration. Import the value with `woobe auth credential import
--name NAME --stdin < PRIVATE_INPUT` first; lookup uses the protected store beside the
selected config. Providers are private POSIX files and Windows Credential Manager.
Environment-variable interpolation, prefixes, remote stores and step-result
secrets are not supported. The stored value is the complete field value.

```json
{"schema_version":"1","project_id":"PROJECT_UUID","steps":[
  {"id":"provider","command":"project provider-credential create","body":{
    "project_id":"PROJECT_UUID","provider":"custom","name":"Provider",
    "secret":{"$secret_ref":"provider-secret"},"metadata":{}}}
]}
```

Compile and plan preserve reference markers without resolving them. Preflight
reports `deferred_protected_credential`; `--require-complete` refuses that deferred
state. Apply resolves a private in-memory snapshot, then validates the wire body
when `--validate-body` is selected. Checkpoints bind reference names to value
fingerprints and contain no resolved values. Changed, unavailable or added
credentials refuse resume before writes; use the original snapshot to finish an
existing checkpoint. Dry-run does not resolve these values or save a checkpoint.

Successful steps resume without another write. Uncertain protected writes remain
uncertain: redacted getters cannot certify their secret values. Reconciliation
and `--skip-unchanged` refuse them. Output and checkpoint results redact resolved
values even when a response echoes them under an unrelated field. This does not
make unsupported server creates idempotent after losing a checkpoint.

## Canonical MCP permission inputs

`project tool mcp set` and `bulk` always validate `allow`, `deny` or `review`, UUID
provider identity, bounded names/list sizes, duplicate names, surrounding
whitespace and unknown fields. This check also runs for dry-run, validate-input
and resolved manifest writes; invalid inputs stop before writes. Server schemas
remain optional additional validation and authorization remains on the backend.

The backend rejects undiscovered remote tool names before applying any part of a
selective update and requires a valid MCP discovery. Bulk affects only that
provider's discovered tools. `review` keeps a tool unavailable; it does not create
a per-call approval dialog. Rediscovery preserves saved modes by remote name and
new tools start in review. Existing immutable release snapshots are unaffected.


## Portable composition

Use `woobe package validate`, `plan`, `import`, `status`, `resume`, `cancel`, and
`export agent|network` for complete native author definitions. A Package requires
the matching backend capability catalog; it does not fall back to Manifest.
`--destination` optionally overrides the export directory (default: ./normalized-resource-name), while `--output` selects text/JSON/JSONL. Export accepts a name or UUID; `--env` defaults to draft, staging/production resolve current versions, and release requires `--version`. Import checkpoints are automatic unless explicitly specified. Use `package bindings SOURCE` for a destination YAML template.
Knowledge can be rebuilt from portable documents or resolved through an external
binding. See [PACKAGES.md](PACKAGES.md) for flags, lifecycle, protection, recovery
and the supported CLI/server pair. Existing resource projections and Manifest
semantic-upsert limitations remain described above.
