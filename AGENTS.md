# AGENTS.md

## Branch discipline (mandatory)

- **Before making any change in this repository, confirm it is on its development branch and switch if it is not.** The development branch is **`dev`**; if the repository has **no `dev` branch**, use **`main`**.
  ```sh
  git branch --show-current
  git checkout dev   # or `main` when the repo has no `dev` branch; then `git pull --rebase`
  ```
- Commit feature/fix work on the development branch. When a `dev` branch exists, never commit feature work on `main` (`main` only receives merges from `dev`).
- If you find yourself on the wrong branch, move the commit instead of redoing it: `git checkout dev && git cherry-pick <sha> && git push origin dev`, then restore the other branch. Report the branch in the final summary.


This file is the CLI subproject guide for agentic coding in `cli/`.
The `cli/` folder is a standalone repository mounted via symlink in the mono-repo.

--------------------------------------------------------------------------------
Scope and precedence
--------------------------------------------------------------------------------

- Applies only to files under `cli/`.
- When this file conflicts with the mono-repo `AGENTS.md`, this file is authoritative for `cli/` changes.
- Keep changes minimal and aligned with existing CLI repository conventions.

--------------------------------------------------------------------------------
Tracking policy (mandatory)
--------------------------------------------------------------------------------

### TODO.md
- Before starting work, read `TODO.md` in the `cli/` repository root.
- Use `TODO.md` to detect pending work, partially applied requirements, blockers, and follow-ups before making changes.
- When work leaves unfinished scope or reveals new follow-ups, update `TODO.md` in the same change set.
- When completing pending items, move them to `DONE.md` and remove from `TODO.md`.
- **Date grouping**: items within each section must be grouped under `## YYYY-MM-DD` headers.
- Treat `TODO.md` as the operational handoff and pending-work tracker for CLI agent workflows.

### DONE.md
- Maintain `DONE.md` as the permanent record of completed CLI work.
- When closing a TODO, move the entry to `DONE.md` under a `## YYYY-MM-DD` header.
- `TODO.md` should never have a lingering "Done Recently" section.

### REVIEW.md
- Maintain `REVIEW.md` as the canonical log for CLI items needing human review.
- A review item is any intentionally omitted or postponed scope decision, technical debt, workaround, command/flag shortcut, or unresolved design trade-off.
- Every review entry must include: date, area, item, reason, impact, and follow-up status.
- If no review item exists in a change set, do not add a placeholder entry.

--------------------------------------------------------------------------------
Build, lint, and test commands
--------------------------------------------------------------------------------

- Install deps: `go mod tidy`
- Build binary: `go build -o uproc`
- Run without build: `go run . --help`
- Format: `gofmt -w .`
- Lint (if available): `go vet ./...`
- Test: `go test ./...`
- Release local dry-run: `goreleaser release --snapshot --clean`
- Release by tag (CI): push tag `vX.Y.Z`

When adding new CLI commands:
- Prefer wrapping existing `/api/v1/external/*` endpoints directly.
- Keep a generic raw command (currently `request`) to avoid endpoint coverage gaps.
- Update `README.md` command list in the same change set.
- If anything changes under `cmd/processes/*` or `cmd/profile.go` (commands, flags, args, hierarchy, output/help text, examples), you MUST update backend CLI docs policy/docs in the same change set:
  - `../back/AGENTS.md`
  - `../back/docs/templates/cli.template.md`
  - `../back/docs/cli.en.md`
  - `../back/docs/cli.es.md`
  - `../back/docs/cli.ca.md`

Install command policy:
- `install` must consume `/api/v1/external/install` and render a full installation plan.
- `install` supports `--dry-run` to print every step without executing server changes.
- Default/expected operational usage is dry-run preview first; any future execution mode must be explicit and opt-in.

Distribution notes:
- Release automation is defined in `.github/workflows/release.yml` and `.goreleaser.yml`.
- Artifacts are produced for Linux/macOS/Windows on `amd64` + `arm64`.
- Packaging targets include GitHub Releases, Homebrew tap, and Scoop bucket.
- `uproc self-update` (see `cmd/selfupdate.go`, `internal/update/`) downloads the
  release asset for the current platform from `uproc-io/uproc-cli`, verifies the
  checksum and replaces the binary; Homebrew/Scoop installs get an instruction
  instead. A proactive 24h-cached notice runs on every command (opt-out
  `UPROC_NO_UPDATE_CHECK=1`). Keep release asset naming in `.goreleaser.yml`
  (`uproc.cli_{version}_{os}_{arch}.tar.gz`/`.zip`) in sync with `internal/update`.

Authentication UX notes:
- `login` supports args and interactive prompt fallback.
- Credentials are stored in `config.yml` profiles (project-local by default).

--------------------------------------------------------------------------------
Code style and safety
--------------------------------------------------------------------------------

- Prefer explicit, minimal changes over broad refactors.
- Follow the style already present in the touched files.
- Do not commit secrets or environment files.
- Add comments only when needed to explain non-obvious behavior.

--------------------------------------------------------------------------------
Tracking policy
--------------------------------------------------------------------------------

- Work is tracked via `TODO.md` (pending), `DONE.md` (completed), and `REVIEW.md` (deferred decisions).
- No `CHANGELOG.md` is maintained.

--------------------------------------------------------------------------------
Commit, test and push (mandatory)
--------------------------------------------------------------------------------

- **Commit and push happen ONCE, at the END of a plan — after ALL of its steps have been
  applied — never between steps.** Steps that modify code do not commit or push mid-plan;
  when the last step is done, close the whole change set automatically with:
  **gates green → tracking docs updated → commit → push to `main`**. Do not wait for an
  explicit "commit push" request (root `AGENTS.md`, "Commit, Test and Push (mandatory)").
- Gates: `gofmt -w .`, `go vet ./...`, `go test ./...` (and `go build` for the binary).
  Command/flag changes must also refresh `README.md` and, per the rule above, the backend
  CLI docs in the same change set.
- Commit: one coherent change set, only the touched files staged (`git status --short` first),
  descriptive imperative message, no secrets or environment files.
- Push to `main` after the green commit; report the commit hash in the final summary.
