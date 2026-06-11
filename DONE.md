# DONE

Completed CLI work, grouped by date.

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
