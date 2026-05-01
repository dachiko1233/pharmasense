# PharmaSense — Progress Log

## Status: Ready to Run

`docker-compose up --build` should bring the full stack up. Demo data seeds automatically on first run.

---

## What Has Been Completed

### Backend (Go)
- Full PostgreSQL schema (7 tables: pharmacies, users, products, inventory_batches, sales, risk_assessments, alert_actions)
- Embedded migrations (no external tool needed)
- Auto-seeding: 150+ realistic products, 500 batches (100 CRITICAL / 125 HIGH / 100 MEDIUM / 175 LOW), ~10K sales records with seasonal patterns
- JWT auth (bcrypt passwords, 24h tokens)
- Risk engine: CRITICAL/HIGH/MEDIUM/LOW classification with suggested discounts and estimated loss
- All API handlers: auth, pharmacy, products, inventory, risk dashboard, alerts, reports
- CSV export (inventory) and CSV import (preview + actual DB insert)
- Auth middleware accepts Bearer header **and** `?token=` query param (needed for file downloads)

### Frontend (Next.js 15, TypeScript)
- Login page with demo credential hint
- Dashboard: 4 KPI cards + Expiry Timeline chart + Top Risk chart + Recent Critical Alerts
- Inventory page: search, filter (risk/category/supplier/expiry range), sort, paginate, CSV export
- Alerts page: tabs (Critical/High/All), action dialog (Apply Discount, Transfer, Return, Dismiss)
- Reports page: Savings over time chart + Category breakdown chart
- Import page: drag-and-drop CSV upload → preview → confirm import (fully wired to backend)
- Settings page: pharmacy info edit (admin only) + user list
- Full i18n: English + Greek (both translation files complete)
- shadcn-style UI components, Tailwind v4, Recharts, Sonner toasts

### Infrastructure
- `docker-compose.yml`: postgres + adminer (port 8080) + backend (port 3001) + frontend (port 3000)
- Backend Dockerfile: Go 1.25-alpine multi-stage build
- Frontend Dockerfile: Node 20-alpine multi-stage standalone build
- Makefile with dev/build/clean/logs/db-shell targets

---

## Demo Credentials

- Admin: `admin@pharmasense.cy` / `Demo1234!`
- Staff: `staff@pharmasense.cy` / `Demo1234!`

---

## Known Issues / Remaining Work

### Not Yet Implemented
- Actual "recalculate" button feedback in UI (backend triggers async recalc, no polling indicator)
- PDF export in reports page (button exists, shows "coming soon" toast)
- Language switcher in login page (routing is set up, but no switcher UI on the login form itself)

### Low Priority Polish
- The import page's "Download Template" link calls the backend at `http://localhost:3001` — works when Docker is running locally; breaks if accessing from a remote machine
- Mobile sidebar toggle (hamburger menu not implemented — sidebar is always visible)

---

## Architecture Decisions

- **No sqlc**: sqlc was in the spec but the project uses raw pgx queries with manual structs. This was intentional to keep the project self-contained without a code generation step.
- **Embedded migrations**: Uses Go's `embed.FS` to bundle the migration SQL directly into the binary — no external migration tool needed.
- **Single-file seeder**: The entire seed logic is in `cmd/seed/seed.go` and runs automatically if the `pharmacies` table is empty and `SEED_ON_EMPTY=true`.
- **Risk calculation inline**: The `services.Calculate()` function is a pure function (not a method) that can be called from both the RiskEngine and the CSV import handler.
