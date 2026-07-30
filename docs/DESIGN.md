# Exqcor v3 — Cross-platform local server + browser clients

## Context

Exqcor is a live-theater "exquisite corpse" playwriting system (see `AGENTS.md`): audience members are called one at a time to write dialogue for ~5 minutes, seeing only the previous writer's last line; the resulting script is printed and cold-read by improv actors. The original Rails 3 app is complete but ancient; `modern/` (Rails 8) rebuilt the data model/admin/API but not the live show flow; `Apple/` is a thin read-only SwiftUI client.

The user has changed direction: instead of native Apple apps, the new system must run on **any laptop (Windows/macOS/Linux)** as a central server, with **all clients (Android/iOS phones, laptops) connecting via web browser**. Production staff install one app, print a **QR code**, and actors/admins/audience scan it to join. `modern/` and `Apple/` stay as reference only; this is a fresh build.

**Decided with user:** Go single-binary server · audience QR gives queue status + turn call + digital program · auto-detect LAN IP (admin picks if ambiguous) · both dedicated writing stations and writers-on-own-phones supported, chosen per production.

## Stack

| Concern | Choice |
|---|---|
| Server | Go, stdlib `net/http` (1.22+ patterns), single executable per OS, `CGO_ENABLED=0` |
| DB | SQLite via `modernc.org/sqlite` (pure Go); WAL mode; hand-rolled migration runner over embedded `NNNN_*.sql` using `PRAGMA user_version` |
| Real-time | **SSE** (stdlib EventSource — auto-reconnect on flaky Wi-Fi; traffic is nearly all server→client; client actions are plain POSTs) |
| QR | `github.com/skip2/go-qrcode`; browser auto-open via `github.com/pkg/browser` |
| Frontend | **Svelte 5 + Vite** plain SPA (no SvelteKit), hash routing (`svelte-spa-router`) so embed.FS needs no history fallback; built into `server/internal/webui/dist/` and embedded via `//go:embed all:dist` |
| Auth | Role tokens embedded in QR URLs → signed HttpOnly role cookie; admin passphrase (localhost = implicit admin); no user accounts |
| Release | GoReleaser + GitHub Actions on tag; matrix linux/amd64+arm64, windows/amd64, mac universal binary |

## Repo layout (new top-level dirs)

```
web/                      # Svelte SPA (vite outDir → ../server/internal/webui/dist)
  src/lib/                # api.js, sse.js, timer.js, stores.js
  src/routes/{admin,write,audience,script}/
  src/styles/print.css
server/                   # Go module
  cmd/exqcor/main.go      # flags (--port/--db/--no-open), startup, auto-open browser
  internal/httpapi/       # routes.go, handlers, sse.go, auth.go
  internal/store/         # queries, migrate.go, migrations/0001_init.sql
  internal/show/          # runtime.go (turn/timer/queue state machine), hub.go (SSE pub/sub)
  internal/netinfo/       # LAN IP detection (private IPv4s; UDP-dial trick for default route)
  internal/qr/  internal/webui/  internal/version/
  Makefile  .goreleaser.yaml
docs/DESIGN.md  docs/SHOW-NIGHT.md
```

Check in a placeholder `webui/dist/index.html` so `go build` works on fresh clone; gitignore the rest of dist.

## Data model (0001_init.sql)

Carries over `modern/` improvements (rename Play→Script, `characters.role` with partial-unique `vosd`, writer tracking) plus new show-runtime tables:

- `scripts(title, description, theme, writing_seconds DEFAULT 300, station_mode 'station'|'byod')`
- `actors(name, bio)`
- `characters(script_id, actor_id, name, description, role 'character'|'vosd')` — unique(script_id,name); partial unique index: one VOSD per script
- `sections(script_id, name, ordering, status 'pending'|'writing'|'complete')`
- `character_sections(character_id, section_id, on_stage)` · `props(section_id, name, description)`
- `sub_sections(section_id, ordering, writer_id, started_at, ends_at, completed_at)` — a writer turn
- `lines(sub_section_id, character_id, text, ordering, created_at)`
- `writers(name, audience_member_id)`
- `audience_members(script_id, number, name, device_token UNIQUE, status 'waiting'|'called'|'writing'|'done'|'skipped', called_at)`
- `settings(key, value)` — admin_pass_hash, role tokens, active_script_id, chosen_ip

Business rules ported from `modern/app/models/` (script.rb, section.rb): script auto-creates VOSD (undeletable); section auto-attaches VOSD + creates sub_section #0 for the admin priming line; script duplication copies characters/sections/on_stage/priming lines (template shows).

## Key mechanics

- **Turn flow**: admin (or station) starts turn → sub_section gets `started_at`/`ends_at`; writer endpoint returns **only the predecessor's last line**, on-stage characters + VOSD chip, and `ends_at`. Each "Add line" POSTs immediately (crash loses at most one keystroke). Countdown computed client-side from `ends_at` + server-time offset; server enforces deadline + 15s grace, then rejects late lines and emits `turn_ended`.
- **SSE hub**: events `queue_changed`, `turn_started/ended`, `your_turn` (targeted), `show_state`; **payloads filtered per role server-side** (writers/audience never see more than the last line); `Last-Event-ID` resume.
- **Persistence rule**: everything show-critical hits SQLite immediately; on boot, recover queue + in-flight turn and re-arm timers — power-cord-yank safe.
- **QR/join**: QR encodes `http://<lan-ip>:8080/j/<role-token>` → sets signed role cookie → redirects to role home. Separate tokens for audience / station / actor (admin needs passphrase on top; tokens regenerable). Audience landing auto-joins the queue (device cookie idempotency) and shows their number huge. Printable QR page with interface picker when multiple LAN IPs.
- **Frontend routes**: `/#/admin` (dashboard), `/#/admin/scripts/:id` (editor), `/#/admin/show` (run-of-show: start turn, call-next, skip), `/#/admin/qr` (print page), `/#/write` (station + BYOD writing UI state machine: idle → writing → grace → handoff), `/#/audience`, `/#/program`, `/#/script/:id[/actor/:id|/character/:id]` with `@media print` playscript CSS (VOSD as italic stage directions, per-actor highlighting).

## Milestones (each independently testable)

1. **M1 Walking skeleton** — server + embedded placeholder page, SQLite + migrations, LAN IP detection, QR endpoint, browser auto-open, `make dev` (vite proxy). *Verify: phone on same Wi-Fi reaches the page via QR URL; test once on Windows early (firewall prompt).*
2. **M2 Production setup** — full admin CRUD + editor UI, VOSD rules, priming sub_section, duplication, admin passphrase/localhost bypass. *Verify: recreate a real past show through the UI; Go table-driven store tests.*
3. **M3 Writing flow (station mode, polling shim)** — turn lifecycle, last-line-only endpoint, character picker, line posting, grace lockout, hand-off. *Verify: 20-second test timer run; check DB ordering.*
4. **M4 Real-time + audience** — SSE hub w/ role filtering, audience join/queue/call-next/skip, "your turn" push, BYOD mode, live admin dashboard, crash-recovery. *Verify: laptop + 2–3 phones + tablet; kill server mid-turn and confirm resume; Wi-Fi toggle reconnect.*
5. **M5 Script output + program** — full/per-actor/per-character views, print CSS, digital program, writing stats. *Verify: Cmd-P → PDF looks like a playscript.*
6. **M6 Show-night hardening** — QR print page + tokens + regeneration, reconnect UX, rate limits/length caps, SHOW-NIGHT.md runbook. *Verify: scan printed paper QRs from a fresh phone; role isolation (audience gets 401 on admin API).*
7. **M7 Packaging** — GoReleaser + GH Actions, per-artifact README (Gatekeeper `xattr -d com.apple.quarantine` / SmartScreen "Run anyway" / firewall-allow notes), v1.0.0. *Verify: download real release artifacts on all three OSes and run a 3-person, 2-act dress rehearsal ending in a printed script.*

## Critical files

- `server/internal/show/runtime.go` — live show state machine (heart of the app)
- `server/internal/store/migrations/0001_init.sql` — schema
- `server/internal/httpapi/routes.go` — API surface + role middleware
- `web/src/routes/write/WritingStation.svelte` — writer turn UX
- `server/cmd/exqcor/main.go` — startup/LAN/embed/auto-open

Reference (read, don't modify): `AGENTS.md`, `modern/db/schema.rb`, `modern/app/models/{script,section,sub_section}.rb`.
