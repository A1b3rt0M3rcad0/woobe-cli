# Completude do planejamento integral

**38.6% — 39/101 entregas concluídas; 33 parciais e 29 pendentes.**

Base: todas as 101 entregas das fases 0–9 do §18 de PLAN.md, com peso igual. Concluída=1; parcial=0; pendente=0. A classificação é uma avaliação de engenharia com evidência por item, não estimativa de esforço, cobertura de código ou certificação de produção.

O denominador inclui backend, CLI, documentação e distribuição. Concluída significa implementação entregue no boundary indicado; os aceites de fase que exigem API real/E2E continuam sem comprovação. Nenhuma das 10 fases tem seu aceite integral verificado.

| Fase | Concluídas | Parciais | Pendentes | Entregas concluídas |
| --- | --- | --- | --- | --- |
| 0 | 2 | 3 | 4 | 2/9 |
| 1 | 0 | 0 | 12 | 0/12 |
| 2 | 4 | 7 | 0 | 4/11 |
| 3 | 5 | 3 | 1 | 5/9 |
| 4 | 0 | 2 | 8 | 0/10 |
| 5 | 5 | 5 | 0 | 5/10 |
| 6 | 6 | 4 | 0 | 6/10 |
| 7 | 7 | 3 | 0 | 7/10 |
| 8 | 6 | 2 | 2 | 6/10 |
| 9 | 4 | 4 | 2 | 4/10 |

## Avaliação item a item

| ID | Entrega | Estado | Evidência / limite |
| --- | --- | --- | --- |
| 0.1 | Atualizar a referência da master no início da implementação. | parcial | docs/PLAN.md records inspected revisions; current backend master integration is unverified. |
| 0.2 | Localizar o cadastro de categorias citado pelo usuário, inclusive em outra branch. | pendente | Canonical authority category implementation has not been located/integrated. |
| 0.3 | Documentar quais contratos serão reaproveitados. | concluída | docs/OPERATIONS.md, docs/STATUS.md and SDK/source pins distinguish reusable routes from proposals. |
| 0.4 | Abrir branch de integração a partir da master. | parcial | One CLI branch exists; the separate backend integration branch/delivery is unverified. |
| 0.5 | Portar ControlPlanePrincipal, emissão de control key, grants, guard e testes pertinentes. | pendente | No backend principal/grant/guard port delivered in this repository. |
| 0.6 | Resolver a colisão Alembic com revisão livre e down_revision atual. | pendente | No Alembic migration integration delivered. |
| 0.7 | Preservar ModelSpec e demais mudanças da master. | pendente | Backend master preservation has not been validated by integration. |
| 0.8 | Inicializar o module Go no repositório canônico `woobe-cli`, fixar toolchain/dependências e documentar a migração do protótipo Python. | concluída | go.mod, cmd/woobe, scripts/package.sh and docs/USAGE.md establish canonical Go client and Python-free runtime. |
| 0.9 | Atualizar ADR e documentação para o estado integrado. | parcial | Client status documents exist; integrated backend ADR/state is missing. |
| 1.1 | Registrar cada operação administrativa por operation_id, principal, escopo e efeito. | pendente | Backend implementation and real authorization matrix are not delivered/verified here. |
| 1.2 | Cobrir aliases/rotas legadas ou negar seu uso com Control Key. | pendente | Backend implementation and real authorization matrix are not delivered/verified here. |
| 1.3 | Corrigir pesquisa POST, auth status MCP, catálogos e runtime environments. | pendente | Backend implementation and real authorization matrix are not delivered/verified here. |
| 1.4 | Separar staging, produção e execução de ferramentas. | pendente | Backend implementation and real authorization matrix are not delivered/verified here. |
| 1.5 | Resolver todos os tipos de ID ao Project proprietário. | pendente | Backend implementation and real authorization matrix are not delivered/verified here. |
| 1.6 | Aplicar política de ações a principals humanos; fechar viewer e teto de delegação de admin. | pendente | Backend implementation and real authorization matrix are not delivered/verified here. |
| 1.7 | Preservar comparação de subconjunto por Project e Workspace. | pendente | Backend implementation and real authorization matrix are not delivered/verified here. |
| 1.8 | Especificar delegação de Runtime Keys e política de grants de Project recém-criado. | pendente | Backend implementation and real authorization matrix are not delivered/verified here. |
| 1.9 | Definir revisão/controle de concorrência para grants e alterações sem mecanismo existente. | pendente | Backend implementation and real authorization matrix are not delivered/verified here. |
| 1.10 | Implementar idempotência administrativa nas operações exigidas e consistência entre criação de Project e grant do criador. | pendente | Backend implementation and real authorization matrix are not delivered/verified here. |
| 1.11 | Documentar autoridade/expiração dos descendentes e semântica de revogação. | pendente | Backend implementation and real authorization matrix are not delivered/verified here. |
| 1.12 | Garantir projeção de segredos no servidor. | pendente | Backend implementation and real authorization matrix are not delivered/verified here. |
| 2.1 | Disponibilizar catálogo de permissions/operations/capabilities versionado. | parcial | Local operation catalog exists; effective permissions/capabilities need server contracts. |
| 2.2 | Disponibilizar introspecção própria e preflight sem mutação. | parcial | server-schema, doctor and validate-input perform reads; self authority/preflight contracts remain absent. |
| 2.3 | Implementar a entrada `cmd/woobe` e a árvore Cobra, com divisão em `internal/cli`, application, controlplane, identity, runtime, contracts, config e output. | parcial | Cobra/controlplane/identity/config/output exist; application/contracts/runtime boundaries are not fully split as proposed. |
| 2.4 | Implementar transporte Go, cancelamento, erros tipados, adapters de credenciais e injeção de stdin/stdout/stderr. | concluída | internal/controlplane, credentials, output and root support cancellation/typed errors and injected IO; client tests pass. |
| 2.5 | Implementar contextos locais e resolução determinística. | concluída | Context precedence, named show, explicit empty overrides, Workspace/Project inheritance, strict bounded config and administrative/runtime attachment covered by config and CLI tests; server remains scope authority. |
| 2.6 | Implementar login/refresh/logout de sessão com cookies e CSRF. | parcial | Cookie/CSRF fixtures exist; Windows session storage and real login E2E remain. |
| 2.7 | Implementar importação e referências protegidas de Control/Runtime Keys. | parcial | POSIX credential references/environment fallback exist; native keychains/Windows protection remain. |
| 2.8 | Implementar JSON, tabela, stderr e códigos de saída. | concluída | internal/output and root implement JSON/table and typed exit/stream error behavior, verified in tests. |
| 2.9 | Portar para Go os comportamentos compatíveis de request/control-key, sem dependência de execução do protótipo Python. | parcial | Go request/control-key handlers exist; branch-dependent Control Key server compatibility is unverified. |
| 2.10 | Implementar doctor/version/help e validação de capacidades. | concluída | doctor/version/help and advertised-route diagnostics are implemented with fixtures. |
| 2.11 | Implementar help JSON e exportação de schemas a partir do registro único, com compatibilidade e efeito das operações. | parcial | Cobra-derived discovery includes 238 handlers and 188 HTTP operations. Shared advertised body validation, manifest preflight, opt-in before-write checks and snapshot digest pin are delivered; full schema/DTO/path/query semantics remain partial. |
| 3.1 | Workspace list/get/create/update. | concluída | workspace.go provides observed CRUD handlers and operation discovery. |
| 3.2 | Membros e convites suportados, com limites de papel. | parcial | Member/invite handlers exist; server role ceilings are unverified. |
| 3.3 | Project list/get/create/update/overview. | concluída | project.go provides observed CRUD/overview handlers. |
| 3.4 | Acesso humano e acesso pendente por convite. | concluída | auth/workspace/project handlers cover observed membership/access and invite operations. |
| 3.5 | Control Key CRUD/lifecycle e grants por Project. | parcial | Five Control Key routes are branch-dependent; canonical server acceptance remains unverified. |
| 3.6 | Project API Keys e Runtime Keys. | concluída | keys.go supplies observed key lifecycle handlers and exclusive secret issuance outputs. |
| 3.7 | Environment keys com valores protegidos. | concluída | Environment key handlers use recursive value/secret output redaction and explicit issuance controls. |
| 3.8 | Alterações incrementais de grants com proteção de revisão. | parcial | Expected revision headers exist; authoritative grant concurrency behavior is unverified. |
| 3.9 | Testes de reader mínimo e credencial com permissões diferentes por Project. | pendente | No real reader/multi-Project authority E2E matrix run. |
| 4.1 | Reutilizar ou implementar entidade/versionamento de categoria. | pendente | Canonical category entity/versioning integration unavailable. |
| 4.2 | Validar permissões primitivas e compatibilidade de scope. | pendente | Canonical primitive permission/scope validation unavailable. |
| 4.3 | List/get/create/update/clone/history/archive/delete. | parcial | Proposed category routes are advertisement-gated; full history/clone lifecycle contract is absent. |
| 4.4 | Aplicar categoria/revisão a grants, preservando lista concreta. | pendente | Canonical category-to-grant expansion unavailable. |
| 4.5 | Atribuições humanas, se fizerem parte do modelo existente/escopo integrado. | pendente | Canonical human assignment contract unavailable. |
| 4.6 | Perfis iniciais versionados de reader/editor/builder/publisher/admin. | parcial | Design lists profiles; canonical integrated versioned profile seeding is absent. |
| 4.7 | Plan/diff de atualização de categoria aplicada. | pendente | No authoritative applied-category plan/diff contract. |
| 4.8 | Subconjunto e condições de delegação após expansão. | pendente | No real expanded delegation subset matrix. |
| 4.9 | Garantir que nova revisão não amplia automaticamente credenciais. | pendente | No verified category revision isolation enforcement. |
| 4.10 | Eventos de auditoria de criação, revisão e aplicação. | pendente | No verified category administrative audit events. |
| 5.1 | Agents e configurações principais. | concluída | agent.go provides observed Agent configuration CRUD. |
| 5.2 | Prompts, contratos e model configs. | concluída | agent.go supplies prompt/contract/model configuration handlers; composed fixture passes. |
| 5.3 | Bindings de Tools, Skills e Knowledge. | parcial | Knowledge binding route exists; complete typed Tool/Skill/Knowledge bindings are not validated. |
| 5.4 | Releases/tests/promotion/activation/rollback. | concluída | release.go supplies create/test/promote/activation/rollback handlers with explicit execution/publication confirmation. |
| 5.5 | Networks múltiplas por Project. | concluída | network.go uses explicit Network IDs under Project; composed fixture covers separate Network creation. |
| 5.6 | Draft, nodes/bindings e revisão esperada. | parcial | Draft update/revision headers exist; full nodes/bindings DTO/reconciliation remains. |
| 5.7 | Versions, previews, staging e production activations. | concluída | network.go provides observed version/promotion/activation commands. |
| 5.8 | Ambientes de Agent/Network com política adequada. | parcial | Environment handlers exist; server policy by environment is unverified. |
| 5.9 | Lifecycle de exclusão/arquivamento conforme suporte real. | parcial | Generic requests exist; complete resource lifecycle coverage/semantics remain unverified. |
| 5.10 | Diff/export de configurações autorizado. | parcial | Authorized projections/capture/field diff exist; complete semantic export for all configuration is absent. |
| 6.1 | Tools HTTP/MCP, lifecycle, uso, testes e execução explícita. | concluída | tool.go provides HTTP/MCP lifecycle/test/execute handlers with explicit execution effects. |
| 6.2 | Discovery/refresh e autenticação MCP suportada. | concluída | Observed MCP discovery/refresh/auth handlers exist; real deployment acceptance remains a separate gate. |
| 6.3 | Permissões MCP allow/deny/review. | parcial | MCP permission endpoints exist; canonical allow/deny/review conditions are not validated. |
| 6.4 | Segredos e credentials com projeção correta. | parcial | Client redaction/secret issuance is tested; server projection and native storage remain unverified. |
| 6.5 | Collections, Documents, upload/search e Vector Snapshots. | concluída | knowledge.go/upload.go supply collection/document/search/snapshot handlers and multipart fixture. |
| 6.6 | Providers/models e discovery separado da persistência. | concluída | catalog.go separates observed provider/model discovery/read/write handlers. |
| 6.7 | Catálogos e compatibilidade de ModelSpec/estratégias. | parcial | Catalogs exist; full strategy/ModelSpec compatibility matrix is absent. |
| 6.8 | Skills e versões. | concluída | surface.go provides Skill/version handlers; resource manifests include typed SkillVersion parents. |
| 6.9 | ChatSurfaces, lifecycle e Access Keys. | concluída | Surface lifecycle/access-key handlers exist; configuration capture supports ChatSurface updates. |
| 6.10 | Paginação/export conforme capacidade real. | parcial | request-pages follows bounded Link next routes and export/capture exists; endpoint pagination/completeness remains unverified. |
| 7.1 | Runtime público por SDK/adapter, com credencial própria. | concluída | runtime.go uses pinned Go SDK with separately selected runtime credentials. |
| 7.2 | Agent e Network Run, HTTP e streaming. | concluída | Agent/Network run/stream adapters use pinned SDK; real target E2E is separately unverified. |
| 7.3 | Contratos de output/External Context. | parcial | Semantic event output exists; full External Context/output contract validation remains incomplete. |
| 7.4 | Active Run, observe/reattach e cancelamento. | concluída | SDK ActiveRun/ObserveRun/CancelRun adapters are implemented and discovered. |
| 7.5 | Sessions/mensagens e histórico administrativo autorizado. | concluída | Observed session/message/admin run handlers are provided. |
| 7.6 | Rerun explícito e reset de contexto suportado. | concluída | Explicit rerun/reset endpoints are registered; they do not silently retry writes. |
| 7.7 | Traces, invocations, snapshots e diagnostics. | concluída | Observed trace/invocation/snapshot/diagnostic handlers are provided. |
| 7.8 | Usage/metrics existentes. | concluída | Observed usage summary/daily handlers are provided. |
| 7.9 | Auditoria administrativa, sem confundir com runtime tracing. | parcial | Audit routes require proposed backend capabilities; real administrative audit contract remains absent. |
| 7.10 | Interrupção, deadline, gap e degradação sem duplicação de Run. | parcial | Cancellation/deadline are supported; reconnect/gap replay scenarios are not implemented/certified. |
| 8.1 | Schema de manifest versionado. | concluída | Versioned step v1/resource v2 schemas are packaged and tested. |
| 8.2 | Validate/diff/plan/export. | parcial | Validate/plan/field diff/projection/capture are implemented; complete semantic export/import is absent. |
| 8.3 | Apply por operação e dependências. | concluída | Explicit actions compile to dependency-ordered steps; apply/ref/recovery fixtures pass. |
| 8.4 | IDs/revisões esperadas e revalidação antes de mutar. | parcial | Expected If-Match forwarded, optional --validate-body checks resolved bodies before write attempts, schema digest pin and read-only preflight exist. Full route/query/DTO/domain/revision semantics and real-backend revalidation remain partial. |
| 8.5 | Checkpoint e retomada por reconciliação. | concluída | Locked key/context-bound checkpoint resume, strict bounded reads/saves, exact returned numeric references, plan/result/dependency integrity and explicit-ID reconciliation are delivered. Original-write attribution/server idempotency and human principal binding remain partial. |
| 8.6 | Relatório parcial/skipped e código de saída próprio. | concluída | Partial reports include checkpoint state/counts, stopped step, save status and cause. not_attempted steps are resumable; committed steps skipped and uncertain attempts require reconciliation. Fixture tests prove dependent validation cannot send a refused write. |
| 8.7 | Importação de categorias sem concessão implícita. | pendente | AuthorityCategory declarative import is unavailable pending canonical server contract. |
| 8.8 | Referências protegidas para credenciais. | pendente | Manifests reject literal secrets; protected secret-provider references are not implemented. |
| 8.9 | Nenhuma publicação/prune implícitos. | concluída | Manifest validation rejects publication/execution/key issuance/delete; there is no implicit prune. |
| 8.10 | Não usar export incompleto como PATCH destrutivo. | concluída | Projection documents are not accepted as manifests; capture selects explicit readable fields and tests preserve omissions/null. |
| 9.1 | Build do executável Go `woobe` e versão semântica com revisão identificável. | concluída | Go build/version exposes source commit, compiler and OS/architecture; packaging injects release version and commit. |
| 9.2 | Instalação limpa a partir de binários publicados e `go install` por tag. | pendente | No published tag/release or fresh go install by tag validated. |
| 9.3 | CI do CLI e suites integradas do backend. | parcial | CLI formatting/module/vet/race/build and bounded schema fuzz CI are implemented; six distribution targets are packaged, source/checksum/schema verified, and native executable smoke runs on Linux/macOS/Windows. Required aggregate job fails on failed/cancelled/skipped groups. Backend integrated suites and real Woobe E2E remain unverified. |
| 9.4 | Documentação de cada comando, permissão, exemplo e capacidade mínima. | parcial | Generated catalog/usage exist; all domain examples/effective permission/capability contracts are incomplete. |
| 9.5 | Completion e smoke de help. | concluída | Cobra completion/help discovery and native build smoke pass. |
| 9.6 | Matriz de compatibilidade cliente/servidor/protocolo. | parcial | SDK/toolchain/protocol pins exist; real server compatibility matrix remains unverified. |
| 9.7 | Binários por sistema/arquitetura declarados, checksums SHA-256 e manifesto dos artefatos. | concluída | Six OS/arch packages include schemas, SHA256 and artifact metadata; CI validates and stores archives and exercises packaged executables on three native runner OSes. Artifacts have finite retention and are development builds, not a tagged release. |
| 9.8 | Guia de migração do CLI mínimo Python para o CLI Go. | concluída | docs/MIGRATION.md maps pinned Python prototype flags, credentials, installation, output, key lifecycle and explicit compatibility gaps to the Go client; no automatic protected-session/checkpoint migration or release is claimed. |
| 9.9 | Changelog e notas de release com capabilities realmente entregues. | parcial | PR/status notes exist; final release changelog/capability notes have not been published. |
| 9.10 | Smoke em instalação Woobe com infraestrutura real e targets de teste. | pendente | No deployed Woobe instance/credentials for real infrastructure smoke. |
