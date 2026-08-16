# Review Log

Items needing human review: technical debt, workarounds, pending decisions, scope cuts.

## Entry Format

- Date: `YYYY-MM-DD`
- Area: short scope name
- Item: what needs review
- Reason: why it was deferred or needs attention
- Impact: functional/tech/ops impact
- Status: `pending` | `planned` | `won't_do`

## Entries

- Date: `2026-08-16` · Area: `self-update` (Windows)
  - Item: En Windows no se puede reemplazar un ejecutable en uso; `replaceExecutable` deja el binario nuevo como `<uproc>.exe.new` e imprime instrucciones para el swap manual.
  - Reason: Reemplazo atómico del exe en ejecución requiere un shim/rename póstumo; se decidió no añadir esa complejidad en v1 (uso principal macOS/Linux).
  - Impact: En Windows el auto-update no queda aplicado hasta cerrar/renombrar manualmente (o re-ejecutar tras el cierre).
  - Status: `pending`

- Date: `2026-08-16` · Area: `self-update` (integridad)
  - Item: Los releases de `uproc.cli` no llevan firma criptográfica (cosign/GPG); la verificación es SHA256 sobre HTTPS desde `checksums.txt`.
  - Reason: GoReleaser no tiene `signs` configurado; añadirlo es un cambio de release (cosign/GPG) que se deja como follow-up.
  - Impact: Protección contra corrupción/errores, no contra un repo GitHub comprometido.
  - Status: `planned`

- Date: `2026-08-16` · Area: `self-update` (rate limit)
  - Item: La GitHub API no autenticada tiene límite de ~60 requests/h por IP; el aviso proactivo podría agotarlo en máquinas compartidas.
  - Reason: El cooldown de 24h (`update-check.json`) mitiga el uso normal; no se autentica la petición.
  - Impact: En condiciones normales inapreciable; el fallo se silencia y no bloquea ningún comando.
  - Status: `won't_do` (aceptado)

