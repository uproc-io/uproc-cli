# Changelog

All notable changes in `bizzmod-cli` should be documented in this file.

## 2026-05-16

### Added
- Added `uproc processes module submit-public-form <customer_domain> <form_slug> <payload_json>` command to submit public Form Management payloads through `/api/v1/external/public/modules/form-management/forms/{customer_domain}/{form_slug}/submit`.
- Added payload parsing/validation helpers for the new command, including explicit invalid JSON and empty payload error paths.
- Added tests for command registration, payload parsing, and public submission path building.

### Changed
- Updated command docs in `README.md` to include the new module command.

### Verification
- `go test ./...`
- `go vet ./...`
- `go build -o uproc`

## 2026-04-16

### Added
- Added module settings read-only commands:
  - `uproc processes module settings-tabs <module_slug>`
  - `uproc processes module settings-tab <module_slug> <tab_key>`
  These commands consume new external read-only settings endpoints and do not expose any settings mutation workflow.
- Added command tree test coverage for module settings commands in `cmd/processes/module_test.go`.

### Added
- Added new root command groups `operations` and `data` with explicit under-construction messaging and docs link (`https://uproc.io`).
- Added `config path` command to print the resolved config file location.
- Added `profile` command group (`list`, `show`, `use`) for profile management.

### Changed
- Reorganized command implementation by moving process-related commands into `cmd/processes/*` and exposing them under `uproc processes ...`.
- Updated root command tree to include `processes`, `profile`, `operations`, and `data` groups.
- Moved interactive shell under processes (`uproc processes interactive`) and updated prompt/help text.
- Renamed distributed binary and CI/build references to `uproc`.
- Switched auth/config flow to `config.yml` profiles only (project-local by default), removing `.env` credential fallback and `UPROC_CONFIG` override behavior.
- Updated root help to show which `config.yml` file is being used while keeping normal command output clean.
- Updated README usage examples to grouped command syntax (`uproc processes ...`) and profile-based login.
- CLI requests now send explicit client-identification headers (`x-client-app: uproc-cli`, `User-Agent: uproc-cli/1.0`) so backend can enforce CLI-specific permissions without affecting non-CLI integrations.
- Updated login interactive/base-URL validation messaging to show `UPROC_PROCESSES_API_URL` while keeping internal config storage as `api_url`.
- Fixed GoReleaser Homebrew repository target to `uproc-io/homebrew-uproc`.
- Added CLI version support (`uproc --version` and `uproc version`) with release-time injection via GoReleaser ldflags.

### Added
- Added `update_homebrew_uproc.sh` to manually update `homebrew-uproc/Formula/uproc.rb` with release version and checksums when automatic formula updates fail.

### Verification
- `go test ./...`
- `go vet ./...`
- `go build -o uproc`

## 2026-04-15

### Added
- Added `release_tag.sh` to automatically compute and push the next semantic version tag (`vX.Y.Z`) with `patch` as default and optional `--minor` / `--major` bump modes.
- Added glob pattern support to `module upload`, allowing masks like `"*.pdf"` and multiple file inputs in a single command.

### Added
- Added a new `admin` command group with external-only management operations for users, customers, credentials, and modules.
- Added `list`, `get`, `create`, and `update` subcommands for `users`, `customers`, and `credentials`, plus `list` and `get` for `modules`, all routed to `/api/v1/external/admin/*`.
- Added `admin tickets list` and `admin tickets get <ticket_id>` mapped to `/api/v1/external/tickets/all` and `/api/v1/external/tickets/{ticket_id}/detail`.
- Added `admin logs` mapped to `/api/v1/external/admin/logs`.
- Added `admin ai-requests` mapped to `/api/v1/external/admin/ai-requests`.
- Added `admin changelog` mapped to `/api/v1/external/admin/changelog`.
- Added interactive `admin create/update` flows for users, customers, credentials, and tickets, driven by contracts fetched from backend API (`/api/v1/external/admin/contracts/{resource}/{action}`).
- Added `admin tickets create` and `admin tickets update` commands with interactive contract support.

### Changed
- Root command registration now includes `admin` to expose dedicated admin workflows without relying on generic raw requests.
- CLI documentation now includes explicit admin command usage examples and payload patterns.
- Admin `create/update` subcommands for `users`, `customers`, and `credentials` are now hidden from CLI help while remaining available.
- Admin read commands now rely on the same external auth header model as other CLI commands (`x-api-key`, `x-customer-domain`, `x-user-email`) after backend auth alignment; no dedicated superadmin-email header is required.
- CLI output rendering now prints only response data payloads (`data`, `data.rows`, or equivalent), hiding response envelope keys like `success`, `message`, `columns`, and `header`.
- Admin create/update commands no longer accept inline payload args; they always run interactive mode and fetch the contract from API before prompting.
- Admin list commands (`users`, `customers`, `credentials`, `tickets`) now fetch list contracts and enforce backend-provided `visible_fields` so CLI list output matches admin table-visible columns.
- `login` interactive flow now runs step-by-step for all credentials when no args are provided, showing current values as defaults and allowing keep-or-update per field.
- `login` now validates updated credentials by calling `/api/v1/external/modules` whenever any value changes, and only then persists local config.
- Replaced `module kpis` with `module overview <module_slug> [kpis|charts|tables]` and added terminal-adapted overview rendering.
- Overview output now renders KPIs, charts (including donut/pie visual summaries), and tables in a readable terminal format.
- Overview rendering now strictly consumes backend overview sections (`kpis`, `charts`, `tables`) and no longer depends on legacy `labels`/`series` fallbacks.
- `module upload` now prints per-file progress and result (`uploading`, `OK`, `FAIL`) and shows a final uploaded summary.

### Verification
- `./release_tag.sh --help`
- `go test ./...`
- `go vet ./...`
- `go build -o uproc`

## 2026-03-28

### Added
- Added `bizzmod update check <CUSTOMER_API_KEY>` as a dry-run-only verification command that fetches `/api/v1/external/install?dry_run=true` and runs local read-only checks for required binaries, services, env vars, and health endpoints.

### Changed
- Root command registration now includes the new `update` command group.
- CLI documentation now includes update-check usage and explicit no-apply behavior.

### Verification
- `go test ./...`
- `go vet ./...`
- `go build -o uproc`
