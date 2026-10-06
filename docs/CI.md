# CI do Woobe CLI

O workflow `CLI` roda em PRs, pushes na master e nas branches `feat/**`/`docs/**`, além de execução manual. O conjunto desta entrega foi integrado na master; versões geradas pelo CI continuam sendo builds de desenvolvimento.

| Job | Aceite |
| --- | --- |
| `test` | Versão canônica/npm e testes de versionamento, gofmt, auditoria de 101 requisitos, módulos, vet, suíte Go com race/coverage, fuzz de schema por 10 segundos com dois workers, build e discovery. |
| `package` | Seis arquivos Linux/macOS/Windows × amd64/arm64, pacote npm com os mesmos binários, schemas v1/v2, SHA256SUMS e artifacts.json. Conteúdo, alvos, nomes, permissões executáveis e SHA de origem devem coincidir. |
| `native-smoke` (três runners) | Baixar os pacotes do mesmo run e executar o alvo correspondente ao host; verificar versão/commit/OS/arch, discovery, schemas, manifest local válido, recusa de JSON inválido, duas páginas por cursor e metadados de coleção parcial em fixture HTTP loopback. Instalar o tarball npm offline com `--ignore-scripts`, executar o comando instalado e `npm exec`, conferir stdin e preservar código de erro do binário. |
| `ci` | Todos os três grupos anteriores devem terminar em success. Failure, cancelled ou skipped impedem o sucesso deste job agregador. |

Os seis jobs concretos são `test`, `package`, três instâncias de `native-smoke` e `ci`. Empacotamento substitui a matriz anterior que apenas compilava seis executáveis. O workflow mantém permissions `contents: read`, timeouts e cancelamento de execuções anteriores do mesmo workflow/event/ref.

O job `ci` é o check estável para configurar nas regras da branch. Sua presença no workflow não configura automaticamente uma proteção obrigatória do GitHub; essas regras dependem da administração do repositório.

## Evidências remotas

- `cli-validation`: coverage.out e discovery.json.
- `woobe-distribution`: seis arquivos compactados, tarball npm, SHA256SUMS e artifacts.json.
- `native-smoke-<runner>`: relatórios JSON do binário e do npm com identidade da origem, alvo efetivamente executado e checks.

O upload falha se os arquivos não existirem. A retenção solicitada é de 90 dias, sujeita às políticas do GitHub. Essas evidências ficam associadas ao run; não são um release permanente. A versão de desenvolvimento contém run number e attempt. Em PRs, o checkout e o pacote correspondem ao commit de integração temporário do GitHub; em push da master, correspondem ao commit efetivo da master. Nenhum job publica tags ou releases.

## Reprodução

```sh
make check
bash scripts/package.sh 0.0.0-dev
python3 scripts/verify_artifacts.py --commit "$(git rev-parse HEAD)"
python3 scripts/smoke_artifacts.py --commit "$(git rev-parse HEAD)" --report native-smoke.json
python3 scripts/smoke_npm.py --commit "$(git rev-parse HEAD)" --report npm-smoke.json
```

O smoke exige um host Linux, macOS ou Windows em amd64/arm64 e executa somente seu próprio alvo. A matriz hospedada executa três combinações concretas de OS/arquitetura, que constam nos relatórios; ela não afirma execução nativa de todos os seis alvos. Não há validação de providers de keychain, sessão protegida Windows, login real, permissões do servidor ou runtime remoto. A evidência de backend está no PR #177; providers nativos e a matriz completa permanecem no planejamento.

O workflow também pode ser reutilizado por **Release CLI**, com a versão exata
de `VERSION`. A publicação depende de todos os gates e usa os artefatos já
validados. A configuração externa da publicação npm está descrita em
[INSTALLATION.md](INSTALLATION.md); preparar o PR não publica releases ou pacotes.

## Pagination continuation evidence — 2026-10-05

The code revision `659baea23458380b84bc7059d2a544c9c06fc1f5` passed all six jobs in [CLI run 37385398090](https://github.com/A1b3rt0M3rcad0/woobe-cli/actions/runs/37385398090). Stored artifacts include six distribution targets, cli-validation and native smoke reports from Linux/macOS/Windows, including body pagination and partial collection checks.

[Woobe PR #177](https://github.com/A1b3rt0M3rcad0/woobe/pull/177) pins that code revision. Its required live CLI workflow passed at backend `8ab9d13953c1cd77462af85a3a3a957395141edc`; other backend workflow results and user approval remain recorded in the PR. The backend pin can remain on the reviewed code revision when a later client commit only updates documentation/evidence.

Request-validation continuation adds loopback native checks for request direction, UUID/date-time and mutation refusal. These are client/package evidence; live server DTO validation remains a separate backend contract.
