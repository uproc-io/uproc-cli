# TODO

Use this file to track pending CLI work, partially applied requirements, blockers, and follow-ups.

All items grouped by date under each section.

## Active

### 2026-08-16

- **Fix `cmd/processes/leads.go` gofmt** — El workflow `ci` falla en el "Format check" (`gofmt -l .` → `cmd/processes/leads.go`). Pre-existente (no tocado por self-update); el `gofmt` local (go 1.25) reformatea el fichero pero CI usa go 1.22 — verificar la versión de gofmt correcta antes de formatear (posible discrepancia entre versiones).
- **Investigar fallo de upload en `release.yml` (`307 location <nil>` de uploads.github.com)** — GoReleaser crea la release como draft y falla al subir assets (2 reintentos). Se publicó v0.1.7 manualmente con `gh release create`; decidir si se corrige el workflow (p.ej. pin de versión de goreleaser / workaround) o se documenta el fallback manual en `README.md`.

### 2026-06-11

- Validate whether admin usage commands also need CSV/export support once backend exposes an external export route.

## Planned

No planned CLI items recorded yet.

## Blocked

No blocked CLI items recorded yet.
