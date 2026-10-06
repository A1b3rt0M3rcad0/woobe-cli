# Portable packages

Package support is being implemented in woobe-cli draft PR #10 and Woobe draft
PR #179. These commands require the matching Package capability/schema catalog;
the stable CLI's Manifest commands remain independent.

Validate a declared closure locally, without loading config, credentials or HTTP:

```sh
woobe package validate ./support
woobe package validate ./support --locked
```

Plan with an existing credential binding:

```sh
woobe package plan ./support --project PROJECT_ID \
  --bind credential.primary-key=CREDENTIAL_UUID \
  --save-plan ./approved-plan.json
```

Use `--bindings destination.yaml` for an `ImportBindings` document. Each declared
alias must have exactly one destination binding. `--bind credential.ALIAS=UUID`
is repeatable, but duplicate aliases are errors, including aliases already
present in the bindings file. Protected references are not resolved by planning.
The approved plan receipt is private and is created exclusively: an existing
file is never overwritten, including when `--yes` is provided.

```sh
woobe package plan ./support --bindings destination.yaml --dry-run
```

Dry-run captures and validates local bytes and binding coverage. It never opens a
credential store or contacts the API. It reports server semantics/authorization
as unevaluated, and cannot save a server-approved plan or validate remote schemas.

`--validate-body` validates the actual planning JSON against the advertised
OpenAPI request-body contract before submission. `--validate-parameters` validates
the upload and plan path/query schemas. `--schema-sha256` can pin that OpenAPI
snapshot. Binary multipart upload validation remains server-owned.

Planning defaults to `--lifecycle draft`. `release` and `production` require
`--release-notes`; `production` also requires `--reason`. The total `--deadline`
defaults to five minutes and includes upload, planning and HTTP requests.

Observe already accepted operations with `package status OPERATION_ID`, optionally
`--wait --deadline 5m`. An accepted operation is not complete. `package resume
OPERATION_ID --revision REVISION` resumes its observed dependency wait;
`package cancel OPERATION_ID` requests cancellation without deleting shared
resources. Package polling only retries transient GET observations, never POST
mutations.

Import/Export, checkpoint integration, native lifecycle materialization and the
complete cross-project/runtime qualification are still under implementation.
Private receipts/checkpoints currently return unsupported exit code 9 on Windows;
there is no unprotected plaintext fallback.
