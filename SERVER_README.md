# Qwixx – Go + Templ + Datastar

This is the server-rendered rewrite of the Qwixx scorecard app.

## Stack

- **Backend:** Go standard library (`net/http`)
- **Templates:** [templ](https://templ.guide)
- **Client reactivity:** [Datastar](https://data-star.dev) (SSE, no app framework)
- **Styles:** Tailwind CSS (same classes as the Svelte app)

## State and persistence

- **No server-side persistence:** Game state is not stored in a database.
- **Client resilience:** State is kept in the browser:
  - **localStorage:** The full game state is saved on every change (via `data-on-signal-patch`) and restored on load (via `data-init`).
  - **Server role:** Each action (set/unset cross, missed roll, clear) is sent to the server; the server applies the change and returns the new state via Datastar’s `PatchSignals`. The client then updates the store and persists to localStorage.

So the app survives server restarts and client refreshes because the single source of truth is the client store, with the server only validating and returning updated state.

## Prerequisites

- Go 1.26+
- Node/pnpm (for Tailwind build)
- [templ](https://templ.guide/docs/install) CLI: `go install github.com/a-h/templ/cmd/templ@latest`

## Build and run

**Regenerate templ + Tailwind (run before building the Go server):**
- `make generate` — runs both
- `make templ` — templ only
- `make tailwind` — Tailwind only
- `pnpm run generate` — same as `make generate` (uses pnpm)

**Build and run the server:**
- `make build` — regenerate assets, then `go build -o qwixx ./cmd/qwixx`
- `make run` — regenerate assets, then `go run ./cmd/qwixx`
- Or manually: `go build ./cmd/qwixx && ./qwixx` (after `make generate`)

4. Open **http://localhost:3000**

## Project layout

- `cmd/qwixx/main.go` – Entrypoint, routes, layout with Datastar script and Tailwind
- `internal/game/state.go` – Game state struct and score/lock logic
- `internal/handlers/handlers.go` – Datastar action handlers (set/unset cross, missed, clear-all)
- `internal/views/*.templ` – Page, rows, cells, dialogs, icons
- `internal/views/helpers.go` – CSS classes and Datastar expression strings
- `cmd/qwixx/static/` – Embedded static assets (input.css, dist.css from Tailwind build)

## API (Datastar actions)

All actions expect the current store (e.g. `game`, `confirmRow`, `confirmCol`, …) in the request body (Datastar sends it by default). Responses are SSE `datastar-patch-signals` with the updated state.

- `POST /action/set-cross?row=&col=` – Set a cross at (row, col)
- `POST /action/unset-cross?row=&col=` – Clear that cross; also clears confirm dialog state
- `POST /action/set-missed?i=` – Mark missed roll for color `i`
- `POST /action/unset-missed?i=` – Clear that missed roll
- `POST /action/clear-all` – Clear scorecard and confirm dialog state
