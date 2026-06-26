# Heritage Weaver AI

Personal mythology and ancestral wisdom platform focused on preserving family stories and turning them into culturally grounded guidance.

**Stack**: React/Next.js frontend · Go backend · PostgreSQL

---

## Documentation

| Document | Description |
| -------- | ----------- |
| [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md) | Full system architecture, stack decisions, API design, data models |
| [`docs/IMPLEMENTATION_PLAN.md`](docs/IMPLEMENTATION_PLAN.md) | Phased roadmap, current status, next tasks, dependencies |

---

## Current State (Phase 0 ✅ — Foundation Complete)

### Backend

- Go API with graceful shutdown and env-based config
- PostgreSQL via `pgx` / `database/sql`
- Migration runner (`go run ./cmd/migrate`)
- `GET /health` with database ping
- Users CRUD (5 endpoints)
- Family Members CRUD (5 endpoints)
- Shared JSON helpers, validation, error responses
- Postgres error classification (unique/foreign key violations)

### Database

10 tables: `users`, `family_members`, `relationships`, `stories`, `wisdom_extracts`, `myth_chapters`, `ancestor_profiles`, `guidance_sessions`, `media_assets`, `privacy_settings`

---

## Run Locally

```bash
# Start Postgres
docker compose up -d

# Create env file
cp backend/.env.example backend/.env

# Run migrations
cd backend && go run ./cmd/migrate

# Start API
go run ./cmd/api

# Check health
curl http://localhost:8080/health
```

---

## Project Layout

```
heritage-beaver/
├── backend/             # Go API
│   ├── cmd/api/         # Server entrypoint
│   ├── cmd/migrate/     # Migration runner
│   ├── internal/        # All application code
│   │   ├── config/      # Environment configuration
│   │   ├── domain/      # Core domain models (structs)
│   │   ├── http/        # Router + HTTP handlers
│   │   └── store/       # PostgreSQL repositories
│   └── migrations/      # SQL migration files
├── frontend/            # Next.js app (coming soon)
├── docs/                # Architecture + implementation plan
└── docker-compose.yml   # PostgreSQL for local dev
```
