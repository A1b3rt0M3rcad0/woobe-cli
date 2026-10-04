# Registered HTTP operations

Generated from the CLI registry. Availability reflects the inspected source, not live capability or successful authorization. Permission entries are historical hints. Local, SDK, upload and export commands are documented separately.

| Command | HTTP | Availability | Permission hint |
| --- | --- | --- | --- |
| `workspace list` | `GET /identity/workspaces` | observed | `` |
| `workspace create` | `POST /identity/workspaces` | observed | `` |
| `workspace get` | `GET /identity/workspaces/{workspace_id}` | observed | `` |
| `workspace update` | `PATCH /identity/workspaces/{workspace_id}` | observed | `` |
| `workspace member list` | `GET /identity/workspaces/{workspace_id}/members` | observed | `` |
| `workspace member remove` | `DELETE /identity/workspaces/{workspace_id}/members/{user_id}` | observed | `` |
| `workspace invite list` | `GET /identity/workspaces/{workspace_id}/invites` | observed | `` |
| `workspace invite create` | `POST /identity/workspaces/{workspace_id}/invites` | observed | `` |
| `workspace invite cancel` | `DELETE /identity/workspaces/{workspace_id}/invites/{invite_id}` | observed | `` |
| `workspace project-access list` | `GET /core/workspaces/{workspace_id}/project-access` | observed | `` |
| `project list` | `GET /core/projects` | observed | `` |
| `project create` | `POST /core/projects` | observed | `` |
| `project get` | `GET /core/projects/{resource_id}` | observed | `` |
| `project update` | `PATCH /core/projects/{resource_id}` | observed | `` |
| `project overview` | `GET /core/projects/{project_id}/overview` | observed | `` |
| `project run list` | `GET /core/projects/{project_id}/runs` | observed | `` |
| `project member list` | `GET /core/projects/{project_id}/members` | observed | `` |
| `project member grant` | `POST /core/projects/{project_id}/members` | observed | `` |
| `project member revoke` | `DELETE /core/projects/{project_id}/members/{user_id}` | observed | `` |
| `project access grant` | `POST /core/projects/{project_id}/access` | observed | `` |
| `project access revoke` | `DELETE /core/projects/{project_id}/access/{subject_type}/{subject_id}` | observed | `` |
| `project env list` | `GET /core/projects/{project_id}/environment-keys` | observed | `project:environment:read` |
| `project env set` | `PUT /core/projects/{project_id}/environment-keys/{key}` | observed | `project:environment:write` |
| `project env unset` | `DELETE /core/projects/{project_id}/environment-keys/{key}` | observed | `project:environment:delete` |
| `project agent list` | `GET /ai/agents` | observed | `agent:read` |
| `project agent create` | `POST /ai/agents` | observed | `agent:write` |
| `project agent get` | `GET /ai/agents/{agent_id}` | observed | `agent:read` |
| `project agent update` | `PATCH /ai/agents/{agent_id}` | observed | `agent:write` |
| `project agent prompt list` | `GET /ai/agents/{agent_id}/prompts` | observed | `agent:read` |
| `project agent prompt create` | `POST /ai/agents/{agent_id}/prompts` | observed | `agent:write` |
| `project agent prompt get` | `GET /ai/agents/{agent_id}/prompts/{version_id}` | observed | `agent:read` |
| `project agent contract list` | `GET /ai/agents/{agent_id}/contracts` | observed | `agent:read` |
| `project agent contract create` | `POST /ai/agents/{agent_id}/contracts` | observed | `agent:write` |
| `project agent contract get` | `GET /ai/agents/{agent_id}/contracts/{contract_id}` | observed | `agent:read` |
| `project agent contract update` | `PATCH /ai/agents/{agent_id}/contracts/{contract_id}` | observed | `agent:write` |
| `project agent model-config list` | `GET /ai/agents/{agent_id}/model-configs` | observed | `agent:read` |
| `project agent model-config create` | `POST /ai/agents/{agent_id}/model-configs` | observed | `agent:write` |
| `project agent model-config get` | `GET /ai/agents/{agent_id}/model-configs/{config_id}` | observed | `agent:read` |
| `project agent release list` | `GET /ai/agents/{agent_id}/releases` | observed | `agent:read` |
| `project agent release create` | `POST /ai/agents/{agent_id}/releases` | observed | `agent:version` |
| `project agent release get` | `GET /ai/agents/{agent_id}/releases/{release_id}` | observed | `agent:read` |
| `project agent release delete` | `DELETE /ai/agents/{agent_id}/releases/{release_id}` | observed | `agent:delete` |
| `project agent release test` | `POST /ai/agents/{agent_id}/release-tests` | observed | `agent:version` |
| `project agent release tests` | `GET /ai/agents/{agent_id}/releases/{release_id}/tests` | observed | `agent:read` |
| `project agent release activations` | `GET /ai/agents/{agent_id}/release-activations` | observed | `agent:read` |
| `project agent environment list` | `GET /ai/agents/{agent_id}/runtime-environments` | observed | `agent:read` |
| `project agent environment update` | `PATCH /ai/agents/{agent_id}/runtime-environments/{environment}` | observed | `agent:write` |
| `project agent release promote` | `POST /ai/agents/{agent_id}/releases/{release_id}/promote` | observed | `agent:production` |
| `project agent release rollback` | `POST /ai/agents/{agent_id}/releases/{release_id}/rollback` | observed | `agent:production` |
| `project network list` | `GET /network/projects/{project_id}/networks` | observed | `network:read` |
| `project network create` | `POST /network/projects/{project_id}/networks` | observed | `network:write` |
| `project network get` | `GET /network/projects/{project_id}/networks/{network_id}` | observed | `network:read` |
| `project network update` | `PATCH /network/projects/{project_id}/networks/{network_id}` | observed | `network:write` |
| `project network draft update` | `PATCH /network/{network_id}/draft` | observed | `network:write` |
| `project network promotion preview` | `POST /network/{network_id}/promotions/preview` | observed | `network:write` |
| `project network promotion create` | `POST /network/{network_id}/promotions` | observed | `network:production` |
| `project network version list` | `GET /network/{network_id}/versions` | observed | `network:read` |
| `project network management versions` | `GET /network/{network_id}/management/versions` | observed | `network:read` |
| `project network management external-context` | `GET /network/{network_id}/management/external-context-contract` | observed | `network:read` |
| `project network activation list` | `GET /network/{network_id}/management/activations` | observed | `network:read` |
| `project network activation create` | `POST /network/{network_id}/production/activations` | observed | `network:production` |
| `project network rollback` | `POST /network/{network_id}/rollback` | observed | `network:production` |
| `project network environment list` | `GET /network/{network_id}/runtime-environments` | observed | `network:read` |
| `project network environment update` | `PATCH /network/{network_id}/runtime-environments/{environment}` | observed | `network:write` |
| `project network session list` | `GET /network/{network_id}/sessions` | observed | `network:read` |
| `project network session create` | `POST /network/{network_id}/sessions` | observed | `network:write` |
| `project network session messages` | `GET /network/{network_id}/sessions/{session_id}/messages` | observed | `network:read` |
| `project network session reset-context` | `POST /network/{network_id}/sessions/{session_id}/reset-context` | observed | `network:write` |
| `project tool list` | `GET /tools/tools` | observed | `tool:read` |
| `project tool create` | `POST /tools/tools` | observed | `tool:write` |
| `project tool update` | `PATCH /tools/tools/{tool_id}` | observed | `tool:write` |
| `project tool delete` | `DELETE /tools/tools/{tool_id}` | observed | `tool:delete` |
| `project tool status` | `PATCH /tools/tools/{tool_id}/status` | observed | `tool:write` |
| `project tool usage-impact` | `GET /tools/tools/{tool_id}/usage-impact` | observed | `tool:read` |
| `project agent tool list` | `GET /tools/agents/{agent_id}/tools` | observed | `agent:read` |
| `project tool test` | `POST /tools/tools/test` | observed | `tool:write` |
| `project tool execute` | `POST /tools/tools/execute` | observed | `tool:write` |
| `project tool mcp discover` | `POST /tools/mcp/discover` | observed | `tool:write` |
| `project tool mcp refresh` | `POST /tools/mcp/refresh` | observed | `tool:write` |
| `project tool mcp set` | `PATCH /tools/mcp/permissions` | observed | `tool:write` |
| `project tool mcp bulk` | `POST /tools/mcp/permissions/bulk` | observed | `tool:write` |
| `project tool mcp oauth start` | `POST /tools/mcp/oauth/start` | observed | `tool:write` |
| `project tool mcp oauth complete` | `POST /tools/mcp/oauth/complete` | observed | `tool:write` |
| `project tool mcp oauth device-start` | `POST /tools/mcp/oauth/device/start` | observed | `tool:write` |
| `project tool mcp oauth device-poll` | `POST /tools/mcp/oauth/device/poll` | observed | `tool:write` |
| `project tool mcp oauth status` | `POST /tools/mcp/oauth/status` | observed | `tool:write` |
| `project tool mcp oauth disconnect` | `POST /tools/mcp/oauth/disconnect` | observed | `tool:write` |
| `project knowledge collection list` | `GET /knowledge/collections` | observed | `knowledge:read` |
| `project knowledge collection create` | `POST /knowledge/collections` | observed | `knowledge:write` |
| `project knowledge collection get` | `GET /knowledge/collections/{collection_id}` | observed | `knowledge:read` |
| `project knowledge collection update` | `PATCH /knowledge/collections/{collection_id}` | observed | `knowledge:write` |
| `project knowledge document list` | `GET /knowledge/documents` | observed | `knowledge:read` |
| `project knowledge document create` | `POST /knowledge/documents` | observed | `knowledge:write` |
| `project knowledge document get` | `GET /knowledge/documents/{document_id}` | observed | `knowledge:read` |
| `project knowledge document delete` | `DELETE /knowledge/documents/{document_id}` | observed | `knowledge:delete` |
| `project knowledge search` | `POST /knowledge/search` | observed | `knowledge:write` |
| `project agent knowledge list` | `GET /knowledge/agent-collections` | observed | `agent:read` |
| `project agent knowledge bind` | `POST /knowledge/agent-collections` | observed | `agent:write` |
| `project agent knowledge unbind` | `DELETE /knowledge/agent-collections/{link_id}` | observed | `agent:delete` |
| `project agent knowledge version-create` | `POST /knowledge/agents/{agent_id}/knowledge-versions` | observed | `agent:write` |
| `project agent knowledge version-list` | `GET /knowledge/agents/{agent_id}/knowledge-versions` | observed | `agent:read` |
| `project knowledge collection overview` | `GET /agent/knowledge/collections/{collection_id}/overview` | observed | `knowledge:read` |
| `project knowledge collection vector-preview` | `POST /agent/knowledge/collections/{collection_id}/vector-preview` | observed | `knowledge:write` |
| `project knowledge snapshot create` | `POST /agent/knowledge/collections/{collection_id}/vector-snapshots` | observed | `knowledge:write` |
| `project knowledge snapshot list` | `GET /agent/knowledge/collections/{collection_id}/vector-snapshots` | observed | `knowledge:read` |
| `project knowledge snapshot get` | `GET /agent/knowledge/vector-snapshots/{snapshot_id}` | observed | `knowledge:read` |
| `project provider-model list` | `GET /ai/provider-models` | observed | `model:read` |
| `project provider-model create` | `POST /ai/provider-models` | observed | `model:write` |
| `project provider-model get` | `GET /ai/provider-models/{provider_model_id}` | observed | `model:read` |
| `project provider-model update` | `PATCH /ai/provider-models/{provider_model_id}` | observed | `model:write` |
| `project provider-model usage` | `GET /ai/provider-models/{provider_model_id}/usage` | observed | `model:read` |
| `project provider-credential list` | `GET /ai/credentials` | observed | `provider:read` |
| `project provider-credential create` | `POST /ai/credentials` | observed | `provider:write` |
| `project provider-credential update` | `PATCH /ai/credentials/{credential_id}` | observed | `provider:write` |
| `project provider-credential rotate` | `POST /ai/credentials/{credential_id}/rotate` | observed | `provider:credential:rotate` |
| `project provider-credential revoke` | `POST /ai/credentials/{credential_id}/revoke` | observed | `provider:credential:revoke` |
| `project provider-credential usage` | `GET /ai/credentials/{credential_id}/usage` | observed | `provider:read` |
| `project catalog context-assemblers` | `GET /ai/context-assemblers` | observed | `` |
| `project catalog execution-strategies` | `GET /ai/execution-strategies` | observed | `` |
| `project catalog knowledge-strategies` | `GET /ai/knowledge-strategies` | observed | `` |
| `project catalog rag-strategies` | `GET /ai/rag-strategies` | observed | `` |
| `project skill list` | `GET /ai/projects/{project_id}/skills` | observed | `skill:read` |
| `project skill create` | `POST /ai/projects/{project_id}/skills` | observed | `skill:write` |
| `project skill version list` | `GET /ai/projects/{project_id}/skills/{skill_id}/versions` | observed | `skill:read` |
| `project skill version create` | `POST /ai/projects/{project_id}/skills/{skill_id}/versions` | observed | `skill:write` |
| `project skill version get` | `GET /ai/projects/{project_id}/skills/versions/{skill_version_id}` | observed | `skill:read` |
| `project surface list` | `GET /chat-surfaces/projects/{project_id}` | observed | `chat_surface:read` |
| `project surface create` | `POST /chat-surfaces/projects/{project_id}` | observed | `chat_surface:write` |
| `project surface get` | `GET /chat-surfaces/{surface_id}` | observed | `chat_surface:read` |
| `project surface update` | `PATCH /chat-surfaces/{surface_id}` | observed | `chat_surface:write` |
| `project surface release` | `POST /chat-surfaces/{surface_id}/release` | observed | `chat_surface:release` |
| `project surface activate` | `POST /chat-surfaces/{surface_id}/activate` | observed | `chat_surface:release` |
| `project surface disable` | `POST /chat-surfaces/{surface_id}/disable` | observed | `chat_surface:write` |
| `project surface archive` | `POST /chat-surfaces/{surface_id}/archive` | observed | `chat_surface:write` |
| `project surface test-sessions` | `POST /chat-surfaces/{surface_id}/test-sessions` | observed | `chat_surface:write` |
| `project surface sessions` | `GET /chat-surfaces/{surface_id}/sessions` | observed | `chat_surface:read` |
| `project surface activity` | `GET /chat-surfaces/{surface_id}/activity` | observed | `chat_surface:read` |
| `workspace control-key list` | `GET /core/workspaces/{workspace_id}/control-keys` | branch-dependent | `keys:read` |
| `workspace control-key create` | `POST /core/workspaces/{workspace_id}/control-keys` | branch-dependent | `keys:create` |
| `workspace control-key update` | `PATCH /core/workspaces/{workspace_id}/control-keys/{control_key_id}` | branch-dependent | `keys:update` |
| `workspace control-key rotate` | `POST /core/workspaces/{workspace_id}/control-keys/{control_key_id}/rotate` | branch-dependent | `keys:rotate` |
| `workspace control-key revoke` | `POST /core/workspaces/{workspace_id}/control-keys/{control_key_id}/revoke` | branch-dependent | `keys:revoke` |
| `project api-key list` | `GET /core/projects/{project_id}/api-keys` | observed | `api_key:read` |
| `project api-key create` | `POST /core/projects/{project_id}/api-keys` | observed | `api_key:write` |
| `project api-key update` | `PATCH /core/projects/{project_id}/api-keys/{api_key_id}` | observed | `api_key:write` |
| `project api-key rotate` | `POST /core/projects/{project_id}/api-keys/{api_key_id}/rotate` | observed | `api_key:write` |
| `project api-key revoke` | `POST /core/projects/{project_id}/api-keys/{api_key_id}/revoke` | observed | `api_key:write` |
| `project surface access-key list` | `GET /chat-surfaces/{surface_id}/access-keys` | observed | `chat_surface:read` |
| `project surface access-key create` | `POST /chat-surfaces/{surface_id}/access-keys` | observed | `chat_surface:write` |
| `project surface access-key rotate` | `POST /chat-surfaces/{surface_id}/access-keys/{api_key_id}/rotate` | observed | `chat_surface:write` |
| `project surface access-key revoke` | `POST /chat-surfaces/{surface_id}/access-keys/{api_key_id}/revoke` | observed | `chat_surface:write` |
| `project runtime-key create` | `POST /core/projects/{project_id}/runtime-keys` | observed | `api_key:write` |
| `permission list` | `GET /identity/permissions` | proposed | `` |
| `permission effective` | `GET /identity/access/me` | proposed | `` |
| `permission check` | `POST /identity/access/check` | proposed | `` |
| `permission explain` | `POST /identity/access/check` | proposed | `` |
| `workspace authority category list` | `GET /identity/workspaces/{workspace_id}/authority-categories` | proposed | `` |
| `workspace authority category get` | `GET /identity/workspaces/{workspace_id}/authority-categories/{category_id}` | proposed | `` |
| `workspace authority category create` | `POST /identity/workspaces/{workspace_id}/authority-categories` | proposed | `` |
| `workspace authority category update` | `PATCH /identity/workspaces/{workspace_id}/authority-categories/{category_id}` | proposed | `` |
| `workspace authority category clone` | `POST /identity/workspaces/{workspace_id}/authority-categories/{category_id}/clone` | proposed | `` |
| `workspace authority category history` | `GET /identity/workspaces/{workspace_id}/authority-categories/{category_id}/versions` | proposed | `` |
| `workspace authority category archive` | `POST /identity/workspaces/{workspace_id}/authority-categories/{category_id}/archive` | proposed | `` |
| `workspace authority category delete` | `DELETE /identity/workspaces/{workspace_id}/authority-categories/{category_id}` | proposed | `` |
| `auth login` | `POST /identity/auth/login` | observed | `` |
| `auth register` | `POST /identity/auth/register` | observed | `` |
| `auth refresh` | `POST /identity/auth/refresh` | observed | `` |
| `auth status` | `GET /identity/users/me` | observed | `` |
| `auth invite list` | `GET /identity/invites/pending` | observed | `` |
| `auth invite accept` | `POST /identity/invites/{invite_id}/accept` | observed | `` |
| `instance status` | `GET /identity/instance/status` | observed | `` |
| `instance bootstrap` | `POST /identity/instance/bootstrap` | observed | `` |
| `runtime run list` | `GET /runtime/runs` | observed | `` |
| `runtime run get` | `GET /runtime/runs/{run_id}` | observed | `` |
| `runtime run cancel` | `POST /runtime/runs/{run_id}/cancel` | observed | `` |
| `runtime run rerun` | `POST /runtime/runs/{run_id}/rerun` | observed | `` |
| `runtime run invocations` | `GET /runtime/runs/{run_id}/model-invocations` | observed | `` |
| `project trace list` | `GET /observability/traces` | observed | `observability:read` |
| `project trace get` | `GET /observability/traces/{trace_id}` | observed | `observability:read` |
| `project usage summary` | `GET /metering/projects/{project_id}/usage` | observed | `usage:read` |
| `project usage daily` | `GET /metering/projects/{project_id}/daily` | observed | `usage:read` |
| `project agent usage daily` | `GET /metering/agents/{agent_id}/daily` | observed | `agent:read` |
| `project network execution diagnostics` | `GET /network/executions/{execution_id}/diagnostics` | observed | `network:read` |
| `project network execution snapshot` | `GET /network/executions/{execution_id}/snapshot` | observed | `network:read` |
| `project network execution events` | `GET /network/executions/{execution_id}/events` | observed | `network:read` |
| `project network execution cancel` | `POST /network/executions/{execution_id}/cancel` | observed | `network:write` |
| `runtime session messages` | `GET /runtime/sessions/{session_id}/messages` | observed | `` |
| `runtime session create` | `POST /runtime/sessions` | observed | `` |
