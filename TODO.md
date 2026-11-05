# TODO

Use this file to track pending CLI work, partially applied requirements, blockers, and follow-ups.

All items grouped by date under each section.

## Active

### 2026-11-02

- **CLI — superadmin acceso a datos de producción: logs, traces AI, payments (ANÁLISIS · complejidad: T1 · cli)** — **Estado actual**: CLI expone `uproc applications admin logs` (`--module-slug` obligatorio), `uproc applications admin ai-requests` (filtro por customer-id, module-slug), `uproc applications admin usage list|summary` (filtro por customer-id, module-slug), `uproc applications admin users|customers|credentials|modules|tickets|changelog` — todos con role guard admin/superadmin. **GAP**: no existe comando `uproc applications admin payments` ni `uproc applications admin transactions` — el usuario NO puede consultar pagos/transacciones desde CLI. **Veredicto**: el superadmin ya puede leer logs, AI requests/traces, usage stats y datos administrativos vía CLI. Solo le falta acceso a datos de pagos/transactions. [cli]

- **Seguridad CLI — supply chain y CI (PENDIENTE · complejidad: T2 · ~1-2d · cli · plan: `../backend/docs/audits/security-audit-2026-10-31.md`)** — Añadir al pipeline de release: `go vet`, `go test`, `govulncheck`, SBOM y secret scanning; revisar que binarios/artefactos no embeben credenciales. Relacionado con SEC-31-08 (CI de seguridad) y CRA/SBOM. [cli]

- **CLI — Campaign Automation: campañas recurrentes, runs y stats (PENDIENTE · complejidad: T2 · ~2-3d · cli · plan: `../backend/docs/plans/campaign-automation-internal-newsletter.md`)** — Añadir comandos `uproc applications campaign schedule`, `campaign run`, `campaign stats`, `campaign suppressions`; wrappers de los endpoints externos nuevos; output formateado (tabla/list) y flags `--campaign-id`, `--since`, `--channel`, `--provider`. Depende de back F9. [cli]

- **CLI — acceso a /summary y /reports (PENDIENTE · complejidad: T2 · ~1d · cli · depende de back)** — **Estado actual**: CLI (`uproc applications …`) no tiene subcomando para summary ni reports. El back tiene endpoints `/api/v1/summary`, `/api/v1/summary/layout`, `/api/v1/reports/weekly`, `/api/v1/reports/security-report`, y más en `business_channel.py`. **Alcance**: (1) nuevos comandos `uproc applications summary get`, `uproc applications summary layout`, `uproc applications reports get`, `uproc applications reports list-periods`; (2) wrappers de `/api/v1/external/*` existentes; (3) output formateado (tabla/list); (4) flags: `--days`, `--customer-id`, `--kind`, `--period-end`, `--all`. Depende del item back "MCP/CLI — acceso a /summary y /reports". [cli]

## 2026-06-11

- Validate whether admin usage commands also need CSV/export support once backend exposes an external export route.

## Planned

No planned CLI items recorded yet.

## Blocked

No blocked CLI items recorded yet.
