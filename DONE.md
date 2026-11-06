# DONE

## 2026-11-06

- **CLI — `admin` como grupo top-level (`uproc admin …`) (HECHO 2026-11-06 · T1 · cli)** — Los recursos admin (`users`, `customers`, `credentials`, `modules`, `tickets`, `logs`, `ai-requests`, `usage`, `audit`, `changelog`) se exponen ahora también en la raíz, en paridad con la sección `/admin` del frontend, junto a `license`. `uproc applications admin …` sigue aceptado (misma implementación vía `NewAdminChildren()`, árboles frescos para cobra). README y docs CLI EN/ES/CA actualizados. Gates: `gofmt`, `go vet`, `go test ./...` verdes (66 tests). `[cli]`

- **CLI admin audit capture — intake automático de tools sin resultado (HECHO 2026-11-06 · T2 · cli)** — Nuevo `uproc applications admin audit capture` que recorre `/api/v1/external/admin/audit` (module/entity/date), detecta `no_result` (`after.matched=false`, `free_reason=no_result` o `reason=no_result`) y escribe una ocurrencia markdown por detección en `<out>/<tool>/occurrences/<fecha>-<trace>.md`. Dedupe por fichero, `--dry-run`, `--max-pages`, y **redacción de email/IP por defecto** (`--include-raw` solo local). Tests: `admin_audit_test.go` (helpers + redacción + dedupe + dry-run) → 72 passed; `gofmt`, `go vet`, `go test ./...` y build verdes. `[cli]`

## 2026-11-06

- **CLI admin audit — recuperar traces por módulo/entidad/resultado (HECHO 2026-11-06 · T1 · cli+back)** — Nuevo `uproc applications admin audit list|show` que replica Admin > Sistema > Auditoría. `list` filtra por `--module-slug`, `--action`, `--actor-type`, `--entity-type`, `--entity-id`, `--trace-id`, `--result`, `--policy-result`, `--from-date/--to-date`, `--customer-id` (superadmin) y paginación; `show <record_id>` devuelve la traza completa con `before`/`after`/`diff` y el payload del run de data-process. Backend: wrapper `/api/v1/external/admin/audit` y `/api/v1/external/admin/audit/{record_id}` reutilizando la lógica y el scoping del router interno (admin → su customer; superadmin → cualquiera). Tests: `test_admin.py` 30 passed; `gofmt`, `go vet`, `go test ./...` y build verdes. `[cli]+[back]`

## 2026-10-31

- **Data Management — comandos de calidad de datos (HECHO · T2 · cli)** — Nuevo `uproc applications data-management quality` con `get`, `analyze`, `rules` (accept/reject de reglas inferidas) y `correct` (preview por defecto, `--apply` para aplicar). Consume los endpoints externos `/api/v1/external/data-management/quality*`; `gofmt`, `go vet`, `go test ./...` y build verdes. `[cli]`


Completed CLI work, grouped by date.

## 2026-09-11

- **DM export report sidecar** — `data-management export-transfer` now parses the export report returned by the API, writes it next to the archive as `<output>.report.json` (mode 0600), and prints a relation summary (`resolved`, `unresolved`, ignored non-positive). The report exposes relation integrity per entity/column and never row values. README and backend CLI docs synchronized. Validation: `go vet ./...`, `go test ./...`, build.

- **CLI API externa directa** — Migrados los seis comandos restantes que usaban `/api/v1/external/mcp/call`: export/import completo de customer, template/bulk-create/reset-password de usuarios y `data column update`. El CLI ya no contiene llamadas MCP; mantiene la misma autorización y confirmaciones al usar rutas externas directas. Validación: `go vet ./...`, `go test ./...` y build.

- **Portable Data Management transfer** — Added `uproc applications data-management export-transfer --output <file>` and `import-transfer --file <file> [--confirm]`. The commands call direct external API endpoints, not MCP: `GET /api/v1/external/data-management/transfer` and `POST /api/v1/external/data-management/transfer/import`. Export writes a mode-0600 portable ZIP; import previews by default and confirms a full replacement limited to the authenticated customer's Data Management schema and rows. README and backend CLI/API/MCP docs synchronized. Validation: `go test ./...`.

- **CI format fix and release-tracking cleanup** — Ran `gofmt -w cmd/processes/leads.go` (map literal alignment), which was making the `ci` "Format check" fail on every push; `gofmt -l .`, `go vet ./...`, and `go test ./...` are now clean. Closed the stale TODO/REVIEW items that the canonical-repo fix and the automated distribution resolved: the `307 location <nil>` asset-upload failure (renamed repo) and the manual Scoop/Homebrew publication.

- **Automated Homebrew/Scoop distribution** — GoReleaser now publishes the CLI to the distribution taps with `UPROC_DISTRIBUTION_TOKEN`. Root cause of stale taps was the Homebrew formula being published to the tap root (`uproc.rb`) because `brews` lacked `directory`; added `directory: Formula`, removed the stray root formula, and synced `Formula/uproc.rb` to the current release. Verified end-to-end with `v0.2.4`: GitHub Release assets, Homebrew `Formula/uproc.rb` (`v0.2.4`), and Scoop `uproc.json` (`v0.2.4`) all updated automatically by the workflow, which finished green.

- **CLI bearer user authentication** — CLI profiles now store `user_api_key` rather than customer key/domain/email fields. `login`, `install`, and `update check` accept a user API key, and every authenticated external request sends `Authorization: Bearer <user-api-key>`. Customer authentication headers are no longer sent. Tests: `go test ./...` passed.

## 2026-08-16

- **Auto-update del CLI (`uproc self-update` + aviso proactivo)** — Nuevo comando `uproc self-update` (flags `--check`, `--version vX.Y.Z`, `--pre`) que consulta el release de `uproc-io/uproc.cli` para la plataforma actual, verifica el SHA256 contra `checksums.txt` y reemplaza el binario en ejecución de forma atómica. Nuevo paquete `internal/update/` (`update.go`, `checksum.go`, `install.go`, `cache.go`): selección de asset `uproc.cli_{version}_{os}_{arch}.tar.gz`/`.zip` según `GOOS`/`GOARCH`, extracción del binario, detección del método de instalación (Homebrew/Scoop → imprime `brew upgrade uproc` / `scoop update uproc` sin tocar el binario; standalone → swap). Aviso proactivo en cada invocación (cache 24h en `os.UserCacheDir()/uproc/update-check.json`, silencioso sin red, opt-out `UPROC_NO_UPDATE_CHECK=1`, skip en builds `dev`). Versión normalizada para semver con `golang.org/x/mod` (v0.19.0, compatible go 1.22). `cmd/root.go` registra el comando y el hook. Tests `internal/update/update_test.go` (normalización, checksums, selección de asset, detección de instalación). Verificado end-to-end contra GitHub (v0.1.0 → v0.1.6). `gofmt`/`go vet`/`go test` ✅. README y AGENTS actualizados.

- **Release v0.1.7 publicada (primera con self-update)** — Tag `v0.1.7` (desde `v0.1.6`, `release_tag.sh`). El CI `release.yml` falló 2× al subir assets (`307 location <nil>` de uploads.github.com; goreleaser dejó drafts sin publicar), así que se publicó manualmente con `gh release create` (cuenta `mcolomer`, admin) usando los assets locales `uproc.cli_0.1.7_{darwin,linux,windows}_{amd64,arm64}` + `checksums.txt` (drafts previos eliminados). Homebrew tap `uproc-io/homebrew-uproc` actualizado a 0.1.7 (`update_homebrew_uproc.sh`); el bucket de Scoop estaba vacío → se creó `uproc.json` 0.1.7 con los hashes de Windows. Verificación end-to-end del auto-update contra la release publicada: binario reportando v0.1.6 → `uproc self-update` descarga v0.1.7, verifica checksum y reemplaza → "Updated uproc from v0.1.6 to v0.1.7"; `--check` → "up to date".

- **Documentación de instalación (Homebrew trust + Scoop)** — `README.md`: en la sección de Homebrew se añade `brew trust uproc-io/uproc` (requerido en Homebrew 5.5+ para cargar taps no oficiales) y se añade la sección "Install or update CLI with Scoop (Windows)" (`scoop bucket add uproc https://github.com/uproc-io/scoop-bucket.git` + `scoop install uproc`, update `scoop update uproc`). Se sincronizó con `home/public/llms-full-cli.txt` y con `back/docs/cli.{en,es,ca}.md` + `back/docs/templates/cli.template.md` (solo la nota de `brew trust`, ya cubrían Scoop).

## 2026-08-14

- **campaign-automation: `sync-audience` y `launch-platforms`** — Nuevos subcomandos `uproc applications campaign sync-audience <campaign_id> <audience_id>` y `launch-platforms <campaign_id>` (envuelven `POST /api/v1/external/modules/campaign-automation/actions/sync_audience` y `launch_platforms`, para sincronizar audiencias como suscriptores y lanzar campañas en plataformas externas como Sendy). README actualizado.

## 2026-08-05

- `README.md` título `# Bizzmod CLI (Go)` → `# Uproc CLI (Go)` (rename de marca; el artefacto distribuido ya era `uproc`).

- Renamed the primary command `uproc processes` → `uproc applications` (kept `uproc processes` as a hidden deprecated alias via `NewProcessesAliasCmd`). Added `uproc applications data-process search [query]` (discover tools) and `uproc applications data-process tool <key>` (full tool spec: params/accepted/format/cost). Synced `README.md` and backend CLI/API/MCP docs (template + en/es/ca).

- Added `uproc processes data-process execute|batch|runs` (inline without daily-call limit; batch over a data-management entity) and `uproc me get|update` (connected user profile via `/api/v1/external/profile`). Synced `README.md` and backend CLI/API docs (template + en/es/ca).

## 2026-06-04

- Added `uproc processes admin usage list` and `uproc processes admin usage summary`, aligned with the external admin usage endpoints and filters.
- Added repository-level TODO tracking policy for CLI agent workflows.
- Updated `module submit-public-form` to use the canonical `form-generator` public route and synced backend API/CLI docs.
- Added `forms submit-public` as the canonical CLI business verb for public forms, while keeping `module submit-public-form` as a deprecated compatibility alias.
- Added the next `forms` lifecycle business verbs: `publish`, `archive`, `restore`, and `mark-submission-processed`.
- Completed the forms CLI mini-batch with `archive-submission`.
- Added `candidate`, `support`, and `approval` CLI business-verb groups.
- Added `campaign`, `contract`, and `order` CLI business-verb groups.
- Added `email`, `process`, and `signals` CLI business-verb groups.
- Added `editorial`, `signing`, and `tax` CLI business-verb groups.
- Added `documents`, `inventory`, and `orders-ingest` CLI business-verb groups.
- Added `cases`, `invoice`, and `sync` CLI business-verb groups.
- Added `leads`, `prospecting`, and `reconciliation` CLI business-verb groups.
- Added `chat` and `invoice-lines` CLI business-verb groups.
- Extended `leads` with `send-proposal` aligned with existing backend workflow.
- Extended `invoice` with `get-pdf` aligned with existing backend workflow.
- Extended `leads` with `list` aligned with existing backend collection read flow.
- Added business-verb list/read commands across the curated CLI groups using backend collection metadata.
