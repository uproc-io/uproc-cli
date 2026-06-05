# TODO

Use this file to track pending CLI work, partially applied requirements, blockers, and follow-ups.

## Active

- Validate whether admin usage commands also need CSV/export support once backend exposes an external export route.

## Planned

No planned CLI items recorded yet.

## Blocked

No blocked CLI items recorded yet.

## Done Recently

- Added `uproc processes admin usage list` and `uproc processes admin usage summary`, aligned with the external admin usage endpoints and filters.
- Added repository-level TODO tracking policy for CLI agent workflows.
- Updated `module submit-public-form` to use the canonical `form-generator` public route and synced the backend API/CLI docs accordingly.
- Added `forms submit-public` as the canonical CLI business verb for public forms, while keeping `module submit-public-form` as a deprecated compatibility alias.
- Added the next `forms` lifecycle business verbs in CLI: `publish`, `archive`, `restore`, and `mark-submission-processed`.
- Completed the forms CLI mini-batch with `archive-submission`.
- Added `candidate`, `support`, and `approval` CLI business-verb groups aligned with existing backend workflows.
- Added `campaign`, `contract`, and `order` CLI business-verb groups aligned with existing backend workflows.
- Added `email`, `process`, and `signals` CLI business-verb groups aligned with existing backend workflows.
- Added `editorial`, `signing`, and `tax` CLI business-verb groups aligned with existing backend workflows.
- Added `documents`, `inventory`, and `orders-ingest` CLI business-verb groups aligned with existing backend workflows.
- Added `cases`, `invoice`, and `sync` CLI business-verb groups aligned with existing backend workflows.
- Added `leads`, `prospecting`, and `reconciliation` CLI business-verb groups aligned with existing backend workflows.
- Added `chat` and `invoice-lines` CLI business-verb groups aligned with existing backend workflows.
- Extended `leads` with `send-proposal` aligned with the existing backend workflow.
- Extended `invoice` with `get-pdf` aligned with the existing backend workflow.
- Extended `leads` with `list` aligned with the existing backend collection read flow.
- Added business-verb list/read commands across the curated CLI groups using backend collection metadata.
