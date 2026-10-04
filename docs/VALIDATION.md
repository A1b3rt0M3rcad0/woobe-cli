# Local validation — 2026-10-04

Compiler: Go 1.22.2, Linux amd64. Dependencies verified by `go mod verify`.

| Check | Result |
| --- | --- |
| `gofmt` and `git diff --check` | Passed |
| `go vet ./...` | Passed |
| `go test -race -coverprofile=coverage.out ./...` | Passed, 22 test functions plus table subtests |
| Total statement coverage | 58.7%; below the final design acceptance gate |
| CLI statement coverage | 54.8% |
| HTTP transport statement coverage | 81.4% |
| Native executable discovery/schema/parent-ID smoke | Passed |
| Linux amd64 and arm64 packaging | Passed |
| macOS amd64 and arm64 packaging | Passed |
| Windows amd64 and arm64 packaging | Passed |
| SHA256 verification of six archives | Passed |
| GitHub CI | Not run: branch publication blocked |
| Real Woobe E2E | Not run: no deployed test instance/credentials |

The tests verify HTTP authorization error mapping, not the backend authorization engine. They also verify no redirected credential request, no administrative write retry, omitted/null PATCH fields, explicit parent flags, private credential files, cross-origin cookie separation, manifest cycles/references, upload form scope, checkpoint skipping and uncertain-write blocking.

Cross-build success verifies compilation/package generation. It does not certify OS keychain support, Windows human-session protection or server interoperability. The complete design remains open as tracked in `STATUS.md`.
