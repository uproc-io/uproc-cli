# DONE

Completed CLI work, grouped by date.

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
