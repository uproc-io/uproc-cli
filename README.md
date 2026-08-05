# Bizzmod CLI (Go)

Minimal CLI to authenticate and call UProc External API endpoints (`/api/v1/external/*`).

## Requirements

- Go 1.22+

## Setup

Credentials are managed in `config.yml` (project root) using profiles.
Use `uproc applications login --profile <name> --use` to create/update a profile.

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
uproc applications login --profile mcolomer@local --use
```

Stores credentials in `./config.yml` under the selected profile.

`login` reads credentials in this order:
- command arguments (optional, still supported)
- existing values from the selected profile
- interactive prompt step-by-step for all values (shows current value as default)

`CUSTOMER_DOMAIN` must be the customer domain identifier (not a URL).

Example:

```bash
uproc applications login --profile mcolomer@local --use
```

`login` always lets you review values and keep/update each one.
When any value changes, CLI validates credentials by calling `/api/v1/external/modules` before saving.

### Raw external request

```bash
uproc applications request <METHOD> <PATH> [JSON_BODY]
```

Example:

```bash
uproc applications request GET /api/v1/external/modules
```

Output is always rendered as readable tables/lists (never raw JSON).
When backend response includes `{ success, message, data }`, CLI prints only `data`.


> **Note**: `uproc applications` was formerly `uproc processes`; the old name still works as a hidden deprecated alias.

### Module commands

```bash
uproc applications module list
uproc applications module get <module_slug>
uproc applications module overview <module_slug> [kpis|charts|tables]
uproc applications module collections <module_slug>
uproc applications module collection <module_slug> <collection_name> [--page 1 --sort-field key --sort-order asc --filter-field key --filter-value val]
uproc applications module data <module_slug> <collection_name> [--page 1 --sort-field key --sort-order asc --filter-field key --filter-value val]
uproc applications module settings-tabs <module_slug>
uproc applications module settings-tab <module_slug> <tab_key>
uproc applications module upload <module_slug> <collection_name> <file_path>
uproc applications module upload <module_slug> <collection_name> "*.pdf"
uproc applications module webhook <module_slug> <collection_name> <payload_json>
```

`module upload` accepts one or more file paths and glob masks. When a mask matches multiple files, CLI uploads each file and prints per-file progress and result.

### Data-process commands

```bash
uproc applications data-process execute <tool> [params_json]
uproc applications data-process batch <tool> <source_entity_id> [column_mapping_json] [fixed_values_json] [output_name]
uproc applications data-process runs
```

`execute` runs a data-process tool inline (no daily inline-call limit; balance and daily budget apply). `batch` queues a background run over a data-management entity.

### Profile (connected user)

```bash
uproc me get
uproc me update first_name=Jane language=es default_app=data-process
```

### Forms commands

```bash
uproc applications forms list [--page 1 --sort-field created_at --sort-order desc --filter-field status --filter-value published]
uproc applications forms list-fields [--page 1]
uproc applications forms list-submissions [--page 1]
uproc applications forms submit-public <customer_domain> <form_slug> <payload_json>
uproc applications forms publish <form_id>
uproc applications forms archive <form_id>
uproc applications forms archive-submission <submission_id>
uproc applications forms restore <form_id>
uproc applications forms mark-submission-processed <submission_id>
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
uproc applications module submit-public-form <customer_domain> <form_slug> <payload_json>
```

### Candidate commands

```bash
uproc applications candidate list-profiles [--page 1]
uproc applications candidate list-job-openings [--page 1]
uproc applications candidate list-applications [--page 1]
uproc applications candidate list-evaluations [--page 1]
uproc applications candidate list-stage-events [--page 1]
uproc applications candidate create-profile <item_json>
uproc applications candidate create-job-opening <item_json>
uproc applications candidate create-application <item_json>
uproc applications candidate move-stage <application_id> <stage>
uproc applications candidate update-status <application_id> <status>
uproc applications candidate create-evaluation <item_json>
```

These commands wrap the existing `candidate-evaluation` business verbs.

### Support commands

```bash
uproc applications support list [--page 1]
uproc applications support create-ticket <item_json>
uproc applications support assign-ticket <ticket_id> <assignee>
uproc applications support reply-ticket <ticket_id> <message>
uproc applications support mark-resolved <ticket_id>
uproc applications support close-ticket <ticket_id>
uproc applications support reopen-ticket <ticket_id>
```

These commands wrap the existing `customer-care` business verbs.

### Approval commands

```bash
uproc applications approval list [--page 1]
uproc applications approval approve <request_id>
uproc applications approval reject <request_id>
uproc applications approval reassign <request_id> <approver> [note]
uproc applications approval cancel <request_id>
```

These commands wrap the existing `approval-management` business verbs.

### Campaign commands

```bash
uproc applications campaign list [--page 1]
uproc applications campaign list-audiences [--page 1]
uproc applications campaign preview-audience <campaign_id> [limit]
uproc applications campaign add-audience <campaign_id> [mode]
uproc applications campaign pause <campaign_id>
uproc applications campaign activate <campaign_id>
```

These commands wrap the existing `campaign-automation` business verbs.

### Contract commands

```bash
uproc applications contract list [--page 1]
uproc applications contract list-expiring [--page 1]
uproc applications contract list-by-counterparty [--page 1]
uproc applications contract renew <contract_id>
uproc applications contract terminate <contract_id>
uproc applications contract update <contract_id> <data_json>
```

These commands wrap the existing `contract-lifecycle` business verbs.

### Order commands

```bash
uproc applications order list [--page 1]
uproc applications order mark-received <order_id>
uproc applications order cancel <order_id>
uproc applications order send-reminder <order_id>
```

These commands wrap the existing `order-track` business verbs.

### Email commands

```bash
uproc applications email list [--page 1]
uproc applications email mark-processed <email_id>
uproc applications email archive <email_id>
```

These commands wrap the existing `email-assistant` business verbs.

### Process commands

```bash
uproc applications process list [--page 1]
uproc applications process retry-step <process_id>
uproc applications process reassign-owner <process_id>
uproc applications process cancel <process_id>
```

These commands wrap the existing `process-visibility` business verbs.

### Signals commands

```bash
uproc applications signals list [--page 1]
uproc applications signals list-executions [--page 1]
uproc applications signals list-activations [--page 1]
uproc applications signals approve <signal_id>
uproc applications signals discard <signal_id>
uproc applications signals mark-pending-review <signal_id>
uproc applications signals activate <signal_id>
uproc applications signals close <signal_id>
```

These commands wrap the existing `market-signals` business verbs.

### Editorial commands

```bash
uproc applications editorial list-opportunities [--page 1]
uproc applications editorial list-projects [--page 1]
uproc applications editorial list-articles [--page 1]
uproc applications editorial list-combined [--page 1]
uproc applications editorial generate-proposal <opportunity_id>
uproc applications editorial generate-article <opportunity_id>
uproc applications editorial publish <opportunity_id>
uproc applications editorial schedule <opportunity_id>
uproc applications editorial discard <opportunity_id>
```

These commands wrap the existing `editorial-engine` business verbs.

### Signing commands

```bash
uproc applications signing list [--page 1]
uproc applications signing cancel <request_id>
uproc applications signing reopen <request_id>
uproc applications signing send-reminder <request_id>
uproc applications signing sync-status <request_id>
```

These commands wrap the existing `document-signing` business verbs.

### Tax commands

```bash
uproc applications tax list [--page 1]
uproc applications tax generate <report_id>
uproc applications tax recalculate <report_id>
uproc applications tax validate <report_id>
uproc applications tax export <report_id>
```

These commands wrap the existing `tax-reporting` business verbs.

### Documents commands

```bash
uproc applications documents list [--page 1]
uproc applications documents mark-ready <document_id>
uproc applications documents mark-review <document_id>
uproc applications documents archive <document_id>
uproc applications documents restore <document_id>
uproc applications documents regenerate <document_id>
```

These commands wrap the existing `document-generator` business verbs.

### Inventory commands

```bash
uproc applications inventory list [--page 1]
uproc applications inventory mark-received <order_id>
uproc applications inventory cancel <order_id>
uproc applications inventory send-reminder <order_id>
```

These commands wrap the existing `inventory-planning` business verbs.

### Orders Ingest commands

```bash
uproc applications orders-ingest list [--page 1]
uproc applications orders-ingest list-emails [--page 1]
uproc applications orders-ingest reprocess <order_id>
uproc applications orders-ingest validate <order_id>
uproc applications orders-ingest send-to-erp <order_id>
```

These commands wrap the existing `order-ingest` business verbs.

### Cases commands

```bash
uproc applications cases list [--page 1]
uproc applications cases list-by-status [--page 1]
uproc applications cases list-by-type [--page 1]
uproc applications cases add-note <case_id> <content> [created_by]
uproc applications cases close <case_id>
uproc applications cases reopen <case_id>
```

These commands wrap the existing `case-lifecycle` business verbs.

### Invoice commands

```bash
uproc applications invoice list [--page 1]
uproc applications invoice issue <invoice_id>
uproc applications invoice rectify <invoice_id> [reason]
uproc applications invoice send <invoice_id> [email] [subject] [message]
uproc applications invoice get-pdf <invoice_id>
```

These commands wrap the existing safe `sales-invoices` business verbs for already-created invoices.

### VeriFactu commands

```bash
uproc applications sales-invoices verifactu resubmit <invoice_id>
uproc applications sales-invoices verifactu backfill
uproc applications sales-invoices verifactu xml <invoice_id>
```

`resubmit` retries the VeriFactu registration of an invoice, `backfill` registers already-issued invoices into the chain (F2 backlog), and `xml` fetches the generated AEAT SF v1.1 file.

### Invoice lines commands

```bash
uproc applications invoice-lines list [--page 1]
uproc applications invoice-lines add <invoice_id> <concept> [quantity] [unit_price] [tax_rate] [sort_order]
uproc applications invoice-lines update <invoice_id> <line_id> [concept] [quantity] [unit_price] [tax_rate] [sort_order]
uproc applications invoice-lines delete <invoice_id> <line_id>
```

These commands wrap the existing safe `sales-invoices` invoice line verbs.

### Purchase invoice commands

```bash
uproc applications purchase-invoices list [--page 1]
uproc applications purchase-invoices list-lines [--page 1]
uproc applications purchase-invoices list-payments [--page 1]
uproc applications purchase-invoices validate <id> [reason]
uproc applications purchase-invoices pay <id> [reason]
uproc applications purchase-invoices assign-from-ingest <invoice_id> [payload_json]
```

These commands wrap the existing `purchase-invoices` business verbs (validate/pay via the module status endpoint, assign via `assign_from_ingest`).

### Sync commands

```bash
uproc applications sync list-workflows [--page 1]
uproc applications sync list-runs [--page 1]
uproc applications sync list-records [--page 1]
uproc applications sync run <workflow_id>
uproc applications sync preview <workflow_id> [limit]
uproc applications sync dry-run <workflow_id> [limit]
uproc applications sync push-to-erp --credential <id> --resource invoices --json '{"number":"H-1","total":121}'
```

These commands wrap the existing `data-sync` business verbs. `push-to-erp` writes a record to an ERP provider using a stored ERP credential (create/update/delete).

### Leads commands

```bash
uproc applications leads list [--page 1 --sort-field created_at --sort-order desc --filter-field status --filter-value qualified]
uproc applications leads generate-proposal <lead_id> [template_id] [title] [description] [output_format]
uproc applications leads send-proposal <lead_id> <mailbox_id> <to_email> <subject> <body> [proposal_url]
uproc applications leads rerun-intelligence <lead_id>
```

These commands wrap the existing safe `lead-management` workflow verbs.

### Prospecting commands

```bash
uproc applications prospecting list-strategies [--page 1]
uproc applications prospecting list-opportunities [--page 1]
uproc applications prospecting list-prospects [--page 1]
uproc applications prospecting list-executions [--page 1]
uproc applications prospecting run-discovery <strategy_id> [company] [domain]
uproc applications prospecting send-to-leads <opportunity_id>
```

These commands wrap the existing `lead-prospecting` workflow verbs.

### Reconciliation commands

```bash
uproc applications reconciliation list-entries [--page 1]
uproc applications reconciliation list-extracts [--page 1]
uproc applications reconciliation list-exports [--page 1]
uproc applications reconciliation list-matches [--page 1]
uproc applications reconciliation reconcile [process_id]
```

This command wraps the existing `financial-reconciliation` workflow verb.

### Chat commands

```bash
uproc applications data-chatbot domains
uproc applications data-chatbot ask <domain> <question> [context] [channel] [sender_id] [origin_session_id]
uproc applications data-chatbot follow-up <origin_session_id> <question> [channel] [domain] [origin_user_id]
uproc applications data-chatbot interactive
uproc applications data-chatbot list [--page 1]
```

These wrap the existing `data-chatbot` workflow verbs:
- `ask` sends a question and returns the response including the `session_id` to resume the conversation.
- `follow-up` continues an existing conversation (session) with a new question; channel/domain/context default to the session's last request.
- `domains` lists the available data domains (channels) with their sample questions.
- `interactive` is a guided flow: pick a domain, select or edit one of the available questions, send it, and keep following up (or switch channel/domain).

All business-verb list commands use the same read flags as the generic module collection reader:
- `--page`
- `--sort-field`
- `--sort-order`
- `--filter-field`
- `--filter-value`

### Admin commands

```bash
uproc applications admin users list [--customer-id 1]
uproc applications admin users get <user_id>

uproc applications admin customers list
uproc applications admin customers get <customer_id>

uproc applications admin credentials list [--customer-id 1 --category ai --type api_key]
uproc applications admin credentials get <credential_id>

uproc applications admin modules list
uproc applications admin modules get <module_slug>

uproc applications admin tickets list
uproc applications admin tickets get <ticket_id>

uproc applications admin logs --module-slug <module_slug> [--level all --page 1]
uproc applications admin ai-requests [--customer-id 1 --module-slug financial-reconciliation --page 1 --limit 25]
uproc applications admin usage list [--customer-id 1 --module-slug financial-reconciliation --source all --from-date 2026-06-01 --to-date 2026-06-05 --page 1 --limit 25]
uproc applications admin usage summary [--customer-id 1 --module-slug financial-reconciliation --source mcp --from-date 2026-06-01 --to-date 2026-06-05]
uproc applications admin changelog
```

These commands wrap the existing external admin read endpoints.

```bash
uproc applications admin changelog
```

Admin create/update subcommands are currently hidden from help output.
Admin create/update commands run interactive contract mode (contracts fetched from API):

```bash
uproc applications admin users create
uproc applications admin users update
uproc applications admin customers create
uproc applications admin customers update
uproc applications admin credentials create
uproc applications admin credentials update
uproc applications admin tickets create
uproc applications admin tickets update
```

All admin commands use external API endpoints under `/api/v1/external/admin/*`, except ticket commands that use `/api/v1/external/tickets/*`.
Admin list output uses backend list contracts (`/api/v1/external/admin/contracts/<resource>/list`) to keep visible columns aligned with Admin UI tables.

### Interactive mode

```bash
uproc applications interactive
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
uproc applications install <CUSTOMER_API_KEY> --dry-run
```

This command fetches `/api/v1/external/install` and shows the full installation plan (release versions, required services, and ordered steps) without executing changes on the server.

### Update check (dry-run only)

```bash
uproc applications update check <CUSTOMER_API_KEY>
```

This command validates update readiness using `/api/v1/external/install?dry_run=true` plus local read-only checks (docker, dokploy, required services, required env vars, and health endpoints). It never executes deployment/apply actions.

## Notes

- All calls send headers required by backend external auth:
  - `x-api-key`
  - `x-customer-domain`
  - `x-user-email`
- `request` allows calling any current/future external endpoint without waiting for a dedicated subcommand.
- CLI output is always displayed in list/table format (never JSON output).
