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
```

`module upload` accepts one or more file paths and glob masks. When a mask matches multiple files, CLI uploads each file and prints per-file progress and result.

### Forms commands

```bash
uproc processes forms list [--page 1 --sort-field created_at --sort-order desc --filter-field status --filter-value published]
uproc processes forms list-fields [--page 1]
uproc processes forms list-submissions [--page 1]
uproc processes forms submit-public <customer_domain> <form_slug> <payload_json>
uproc processes forms publish <form_id>
uproc processes forms archive <form_id>
uproc processes forms archive-submission <submission_id>
uproc processes forms restore <form_id>
uproc processes forms mark-submission-processed <submission_id>
```

`forms submit-public` is the canonical business-verb command for public form submissions and posts to the public route under `form-generator`.

The lifecycle commands call the generic external action routes for `form-generator`:
- `publish <form_id>`
- `archive <form_id>`
- `archive-submission <submission_id>`
- `restore <form_id>`
- `mark-submission-processed <submission_id>`

Compatibility alias (deprecated):

```bash
uproc processes module submit-public-form <customer_domain> <form_slug> <payload_json>
```

### Candidate commands

```bash
uproc processes candidate list-profiles [--page 1]
uproc processes candidate list-job-openings [--page 1]
uproc processes candidate list-applications [--page 1]
uproc processes candidate list-evaluations [--page 1]
uproc processes candidate list-stage-events [--page 1]
uproc processes candidate create-profile <item_json>
uproc processes candidate create-job-opening <item_json>
uproc processes candidate create-application <item_json>
uproc processes candidate move-stage <application_id> <stage>
uproc processes candidate update-status <application_id> <status>
uproc processes candidate create-evaluation <item_json>
```

These commands wrap the existing `candidate-evaluation` business verbs.

### Support commands

```bash
uproc processes support list [--page 1]
uproc processes support create-ticket <item_json>
uproc processes support assign-ticket <ticket_id> <assignee>
uproc processes support reply-ticket <ticket_id> <message>
uproc processes support mark-resolved <ticket_id>
uproc processes support close-ticket <ticket_id>
uproc processes support reopen-ticket <ticket_id>
```

These commands wrap the existing `customer-care` business verbs.

### Approval commands

```bash
uproc processes approval list [--page 1]
uproc processes approval approve <request_id>
uproc processes approval reject <request_id>
uproc processes approval reassign <request_id> <approver> [note]
uproc processes approval cancel <request_id>
```

These commands wrap the existing `approval-management` business verbs.

### Campaign commands

```bash
uproc processes campaign list [--page 1]
uproc processes campaign list-audiences [--page 1]
uproc processes campaign preview-audience <campaign_id> [limit]
uproc processes campaign add-audience <campaign_id> [mode]
uproc processes campaign pause <campaign_id>
uproc processes campaign activate <campaign_id>
```

These commands wrap the existing `campaign-automation` business verbs.

### Contract commands

```bash
uproc processes contract list [--page 1]
uproc processes contract list-expiring [--page 1]
uproc processes contract list-by-counterparty [--page 1]
uproc processes contract renew <contract_id>
uproc processes contract terminate <contract_id>
uproc processes contract update <contract_id> <data_json>
```

These commands wrap the existing `contract-lifecycle` business verbs.

### Order commands

```bash
uproc processes order list [--page 1]
uproc processes order mark-received <order_id>
uproc processes order cancel <order_id>
uproc processes order send-reminder <order_id>
```

These commands wrap the existing `order-track` business verbs.

### Email commands

```bash
uproc processes email list [--page 1]
uproc processes email mark-processed <email_id>
uproc processes email archive <email_id>
```

These commands wrap the existing `email-assistant` business verbs.

### Process commands

```bash
uproc processes process list [--page 1]
uproc processes process retry-step <process_id>
uproc processes process reassign-owner <process_id>
uproc processes process cancel <process_id>
```

These commands wrap the existing `process-visibility` business verbs.

### Signals commands

```bash
uproc processes signals list [--page 1]
uproc processes signals list-executions [--page 1]
uproc processes signals list-activations [--page 1]
uproc processes signals approve <signal_id>
uproc processes signals discard <signal_id>
uproc processes signals mark-pending-review <signal_id>
uproc processes signals activate <signal_id>
uproc processes signals close <signal_id>
```

These commands wrap the existing `market-signals` business verbs.

### Editorial commands

```bash
uproc processes editorial list-opportunities [--page 1]
uproc processes editorial list-projects [--page 1]
uproc processes editorial list-articles [--page 1]
uproc processes editorial list-combined [--page 1]
uproc processes editorial generate-proposal <opportunity_id>
uproc processes editorial generate-article <opportunity_id>
uproc processes editorial publish <opportunity_id>
uproc processes editorial schedule <opportunity_id>
uproc processes editorial discard <opportunity_id>
```

These commands wrap the existing `editorial-engine` business verbs.

### Signing commands

```bash
uproc processes signing list [--page 1]
uproc processes signing cancel <request_id>
uproc processes signing reopen <request_id>
uproc processes signing send-reminder <request_id>
uproc processes signing sync-status <request_id>
```

These commands wrap the existing `document-signing` business verbs.

### Tax commands

```bash
uproc processes tax list [--page 1]
uproc processes tax generate <report_id>
uproc processes tax recalculate <report_id>
uproc processes tax validate <report_id>
uproc processes tax export <report_id>
```

These commands wrap the existing `tax-reporting` business verbs.

### Documents commands

```bash
uproc processes documents list [--page 1]
uproc processes documents mark-ready <document_id>
uproc processes documents mark-review <document_id>
uproc processes documents archive <document_id>
uproc processes documents restore <document_id>
uproc processes documents regenerate <document_id>
```

These commands wrap the existing `document-generator` business verbs.

### Inventory commands

```bash
uproc processes inventory list [--page 1]
uproc processes inventory mark-received <order_id>
uproc processes inventory cancel <order_id>
uproc processes inventory send-reminder <order_id>
```

These commands wrap the existing `inventory-planning` business verbs.

### Orders Ingest commands

```bash
uproc processes orders-ingest list [--page 1]
uproc processes orders-ingest list-emails [--page 1]
uproc processes orders-ingest reprocess <order_id>
uproc processes orders-ingest validate <order_id>
uproc processes orders-ingest send-to-erp <order_id>
```

These commands wrap the existing `order-ingest` business verbs.

### Cases commands

```bash
uproc processes cases list [--page 1]
uproc processes cases list-by-status [--page 1]
uproc processes cases list-by-type [--page 1]
uproc processes cases add-note <case_id> <content> [created_by]
uproc processes cases close <case_id>
uproc processes cases reopen <case_id>
```

These commands wrap the existing `case-lifecycle` business verbs.

### Invoice commands

```bash
uproc processes invoice list [--page 1]
uproc processes invoice issue <invoice_id>
uproc processes invoice rectify <invoice_id> [reason]
uproc processes invoice send <invoice_id> [email] [subject] [message]
uproc processes invoice get-pdf <invoice_id>
```

These commands wrap the existing safe `invoice-generator` business verbs for already-created invoices.

### Invoice lines commands

```bash
uproc processes invoice-lines list [--page 1]
uproc processes invoice-lines add <invoice_id> <concept> [quantity] [unit_price] [tax_rate] [sort_order]
uproc processes invoice-lines update <invoice_id> <line_id> [concept] [quantity] [unit_price] [tax_rate] [sort_order]
uproc processes invoice-lines delete <invoice_id> <line_id>
```

These commands wrap the existing safe `invoice-generator` invoice line verbs.

### Sync commands

```bash
uproc processes sync list-workflows [--page 1]
uproc processes sync list-runs [--page 1]
uproc processes sync list-records [--page 1]
uproc processes sync run <workflow_id>
uproc processes sync preview <workflow_id> [limit]
uproc processes sync dry-run <workflow_id> [limit]
```

These commands wrap the existing `data-sync` business verbs.

### Leads commands

```bash
uproc processes leads list [--page 1 --sort-field created_at --sort-order desc --filter-field status --filter-value qualified]
uproc processes leads generate-proposal <lead_id> [template_id] [title] [description] [output_format]
uproc processes leads send-proposal <lead_id> <mailbox_id> <to_email> <subject> <body> [proposal_url]
uproc processes leads rerun-intelligence <lead_id>
```

These commands wrap the existing safe `lead-management` workflow verbs.

### Prospecting commands

```bash
uproc processes prospecting list-strategies [--page 1]
uproc processes prospecting list-opportunities [--page 1]
uproc processes prospecting list-prospects [--page 1]
uproc processes prospecting list-executions [--page 1]
uproc processes prospecting run-discovery <strategy_id> [company] [domain]
uproc processes prospecting send-to-leads <opportunity_id>
```

These commands wrap the existing `lead-prospecting` workflow verbs.

### Reconciliation commands

```bash
uproc processes reconciliation list-entries [--page 1]
uproc processes reconciliation list-extracts [--page 1]
uproc processes reconciliation list-exports [--page 1]
uproc processes reconciliation list-matches [--page 1]
uproc processes reconciliation reconcile [process_id]
```

This command wraps the existing `financial-reconciliation` workflow verb.

### Chat commands

```bash
uproc processes chat list [--page 1]
uproc processes chat ask <domain> <question> [context] [channel] [sender_id] [origin_session_id]
```

This command wraps the existing `data-chatbot` workflow verb.

All business-verb list commands use the same read flags as the generic module collection reader:
- `--page`
- `--sort-field`
- `--sort-order`
- `--filter-field`
- `--filter-value`

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
uproc processes admin usage list [--customer-id 1 --module-slug financial-reconciliation --source all --from-date 2026-06-01 --to-date 2026-06-05 --page 1 --limit 25]
uproc processes admin usage summary [--customer-id 1 --module-slug financial-reconciliation --source mcp --from-date 2026-06-01 --to-date 2026-06-05]
uproc processes admin changelog
```

These commands wrap the existing external admin read endpoints.

```bash
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
