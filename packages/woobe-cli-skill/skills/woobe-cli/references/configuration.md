# Remote Project configuration

Use this flow when developing resources without access to Provider/Tool secrets.
It requires a backend advertising `project_configuration.v1`.

1. Declare a reference with `woobe project configuration declare PRIMARY_API_KEY
   --kind secret --format token` (one complete command line).
2. Bind resources using typed remote references. Read
   [dependencies](dependencies.md) for the resource workflow. Provider requests
   use `authentication.type: api_key` and
   `authentication.fields.api_key: {source: project_environment, key: PRIMARY_API_KEY}`.
3. An authorized operator fills values and approves the intended destination in
   Project → Environments. Never request a token in chat, print it, or save it
   in author YAML, instructions, command flags or tool transcripts.
4. Inspect `woobe project configuration list --env draft`. Use `--status missing`
   and `--search PRIMARY` to narrow results. Then run
   `woobe project configuration check agent AGENT_UUID --env draft` or
   `woobe project configuration check network NETWORK_UUID --env staging`.
   Replace the UUID placeholders with real target IDs. `blocked` exits 6;
   `unknown` or incompatible capabilities exit 9. Resolve these diagnostics
   before execution. Metadata readiness does not certify an external test.
5. Run only the test/execution authorized by the user after the server accepts
   its configuration. Filling values never automatically reruns or deploys.

`woobe project configuration usage KEY_UUID` reports retained consumers. Replace
`KEY_UUID` with the declaration ID. List and usage support bounded pages with
`--offset`; list also accepts `--limit` up to 200. The default Environment is
Draft. No configuration command needs `.woobe-config` or a local secret store.

If capability negotiation fails, report the incompatibility. Do not fall back
to reading local protected files or using a value from another Environment.
`--dry-run` is offline metadata preview, not evidence of server readiness.

Bind an existing connection with `woobe project configuration bind provider CREDENTIAL_UUID --key PRIMARY_API_KEY`. Tool bindings support explicit api_key, bearer, basic, headers or none; use declared names only. The server compares the binding revision and rejects concurrent edits. Migrate existing inline credentials before replacing their authority.

Use `woobe project configuration test provider CREDENTIAL_UUID --env staging` only when an external connection test is requested. For MCP discovery use `woobe project configuration test tool MCP_TOOL_UUID --env draft`. These make outbound requests and may save MCP discovery; they do not execute an Agent or arbitrary HTTP Tool. `check` remains read-only with no external I/O. Never retrieve a value or substitute a local token to make readiness succeed.
