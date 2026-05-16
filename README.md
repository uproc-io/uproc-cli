# Bizzmod CLI (Go)

Minimal CLI to authenticate and call UProc External API endpoints (`/api/v1/external/*`).

## Requirements

- Go 1.22+

## Setup

Credentials are managed in `config.yml` (project root) using profiles.
Use `uproc processes login --profile <name> --use` to create/update a profile.

## Build and run

```bash
go mod tidy
go build -o uproc
./uproc --help
./uproc --version
```

Or run directly:

```bash
go run . --help
```

## Distribution

This CLI is configured for multi-platform binary distribution using GoReleaser.

Targets:
- Linux: `amd64`, `arm64`
- macOS: `amd64`, `arm64`
- Windows: `amd64`, `arm64`

Packaging:
- GitHub Releases artifacts + checksums
- Homebrew tap formula update
- Scoop manifest update

Important:
- Homebrew formula repository is `uproc-io/homebrew-uproc`.

### Release process

1. Run local checks:

```bash
gofmt -w .
go vet ./...
go test ./...
```

2. Optional local release dry-run:

```bash
goreleaser release --snapshot --clean
```

3. Create and push a version tag:

```bash
git tag v0.1.0
git push origin v0.1.0
```

Or use the automatic tag script (patch by default):

```bash
./release_tag.sh
# optional:
./release_tag.sh --minor
./release_tag.sh --major
```

GitHub Actions (`.github/workflows/release.yml`) will publish the release.

If Homebrew formula is not updated automatically, run:

```bash
./update_homebrew_uproc.sh --tag vX.Y.Z --push-pr
```

This updates `Formula/uproc.rb` in `uproc-io/homebrew-uproc` with the release version and checksums from `checksums.txt`.

### Check release in GitHub with `gh`

```bash
# Replace with the tag you want to verify
TAG="vX.Y.Z"

# 1) Check the tag exists in GitHub
gh api "repos/uproc-io/uproc.cli/git/ref/tags/${TAG}" --jq '.ref'

# 2) Check GitHub Release exists for that tag
gh release view "${TAG}" --repo uproc-io/uproc.cli \
  --json tagName,isDraft,isPrerelease,publishedAt,url
```

Expected:
- Step 1 returns `refs/tags/vX.Y.Z`
- Step 2 returns release metadata and URL

### Install or update CLI with Homebrew

```bash
# First-time setup (tap + install)
brew tap uproc-io/uproc
brew install uproc

# Update to latest available version
brew update
brew upgrade uproc
```

## Commands

### Auth

```bash
uproc processes login --profile mcolomer@local --use
```

Stores credentials in `./config.yml` under the selected profile.

`login` reads credentials in this order:
- command arguments (optional, still supported)
- existing values from the selected profile
- interactive prompt step-by-step for all values (shows current value as default)

`CUSTOMER_DOMAIN` must be the customer domain identifier (not a URL).

Example:

```bash
uproc processes login --profile mcolomer@local --use
```

`login` always lets you review values and keep/update each one.
When any value changes, CLI validates credentials by calling `/api/v1/external/modules` before saving.

### Raw external request

```bash
uproc processes request <METHOD> <PATH> [JSON_BODY]
```

Example:

```bash
uproc processes request GET /api/v1/external/modules
```

Output is always rendered as readable tables/lists (never raw JSON).
When backend response includes `{ success, message, data }`, CLI prints only `data`.

### Module commands

```bash
uproc processes module list
uproc processes module get <module_slug>
uproc processes module overview <module_slug> [kpis|charts|tables]
uproc processes module collections <module_slug>
uproc processes module collection <module_slug> <collection_name> [--page 1 --sort-field key --sort-order asc --filter-field key --filter-value val]
uproc processes module data <module_slug> <collection_name> [--page 1 --sort-field key --sort-order asc --filter-field key --filter-value val]
uproc processes module settings-tabs <module_slug>
uproc processes module settings-tab <module_slug> <tab_key>
uproc processes module upload <module_slug> <collection_name> <file_path>
uproc processes module upload <module_slug> <collection_name> "*.pdf"
uproc processes module webhook <module_slug> <collection_name> <payload_json>
uproc processes module submit-public-form <customer_domain> <form_slug> <payload_json>
```

`module upload` accepts one or more file paths and glob masks. When a mask matches multiple files, CLI uploads each file and prints per-file progress and result.

### Admin commands

```bash
uproc processes admin users list [--customer-id 1]
uproc processes admin users get <user_id>

uproc processes admin customers list
uproc processes admin customers get <customer_id>

uproc processes admin credentials list [--customer-id 1 --category ai --type api_key]
uproc processes admin credentials get <credential_id>

uproc processes admin modules list
uproc processes admin modules get <module_slug>

uproc processes admin tickets list
uproc processes admin tickets get <ticket_id>

uproc processes admin logs --module-slug <module_slug> [--level all --page 1]
uproc processes admin ai-requests [--customer-id 1 --module-slug financial-reconciliation --page 1 --limit 25]
uproc processes admin changelog
```

Admin create/update subcommands are currently hidden from help output.
Admin create/update commands run interactive contract mode (contracts fetched from API):

```bash
uproc processes admin users create
uproc processes admin users update
uproc processes admin customers create
uproc processes admin customers update
uproc processes admin credentials create
uproc processes admin credentials update
uproc processes admin tickets create
uproc processes admin tickets update
```

All admin commands use external API endpoints under `/api/v1/external/admin/*`, except ticket commands that use `/api/v1/external/tickets/*`.
Admin list output uses backend list contracts (`/api/v1/external/admin/contracts/<resource>/list`) to keep visible columns aligned with Admin UI tables.

### Interactive mode

```bash
uproc processes interactive
```

Inside interactive mode, run commands without the binary name:

```text
uproc> module list
uproc> module get order-track
uproc> request GET /api/v1/external/modules
uproc> help
uproc> exit
```

### Install plan (dry-run)

```bash
uproc processes install <CUSTOMER_API_KEY> --dry-run
```

This command fetches `/api/v1/external/install` and shows the full installation plan (release versions, required services, and ordered steps) without executing changes on the server.

### Update check (dry-run only)

```bash
uproc processes update check <CUSTOMER_API_KEY>
```

This command validates update readiness using `/api/v1/external/install?dry_run=true` plus local read-only checks (docker, dokploy, required services, required env vars, and health endpoints). It never executes deployment/apply actions.

## Notes

- All calls send headers required by backend external auth:
  - `x-api-key`
  - `x-customer-domain`
  - `x-user-email`
- `request` allows calling any current/future external endpoint without waiting for a dedicated subcommand.
- CLI output is always displayed in list/table format (never JSON output).
