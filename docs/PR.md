# Title

feat: validate manifest bodies before writes and harden recovery transport

# Description

Adds exactly 30 commits after PR #3, for 143 commits above the initial master. Resolved manifest bodies can opt into advertised schema validation before write attempts; read-only preflight reports deferred dependency values, strict completeness and schema snapshot identity. Unsupported validation flags are refused. The backend remains responsible for authorization, route/query schemas and domain rules.

Checkpoint loading/saving is bounded, plan/result/dependency ownership is verified, returned numeric precision is preserved and pre-write failures are resumable with structured partial evidence. Default writes use fresh HTTP/1 connections to block net/http replay even with Idempotency-Key; reads retain pooled connections. Ambiguous response JSON, path segments and conflicting scope queries are refused.

Local validation: 151 test functions plus subtests, race tests, vet, module verification, native build/discovery. Coverage 77.6%; discovery 238 handlers / 188 HTTP operations. Final published-head CI and six-target package verification are reported in PR metadata.

Full-plan completion remains 39/101 (38.6%), 33 partial, 29 pending. This advances partial requirements without claiming their complete backend/domain/E2E acceptance. Canonical backend authority, real Woobe E2E, complete semantic reconciliation and native credential providers remain open. No backend mutation, merge or release is performed. Review order: #1 → #2 → #3 → recovery/preflight PR.
