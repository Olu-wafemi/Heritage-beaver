# Hearthside

Every family gathers at the fire. Hearthside keeps your family's stories, proverbs, and voices — and weaves them into guidance for the generations to come.

A personal mythology platform: tell true stories in your own tongue, draw out the wisdom inside them, and gather that wisdom into chapters of a living family myth.

---

## What it does

| Stage | What you do | What Hearthside does |
|-------|-------------|----------------------|
| **The People** | Name your family, living and remembered | Holds the tree |
| **The Bonds** | Tie them together — parent, sibling, elder | Shapes the branches |
| **The Stories** | Tell what happened, the way you heard it | Keeps the telling |
| **The Wisdom** | — | Listens for proverbs, warnings, and values — in any language |
| **The Myth** | Gather stories into a chapter | Weaves them into mythic narrative |

---

## Stack

| Layer | Choice |
|-------|--------|
| Frontend | Next.js 16 · React 19 · Tailwind 4 · Zustand |
| Backend | Go 1.26 · `net/http` · `pgx` · goose migrations |
| Database | PostgreSQL 16 (local install or Docker) |
| Email | Resend (prod) or Gmail SMTP (free dev) |
| Intelligence | OpenCode Zen relay, keyless (`laguna-s-2.1-free` by default) |

---

## Quick start

Requires Go 1.26+, Node 20+, and PostgreSQL.

```bash
# 1. Clone
git clone https://github.com/oluwafemiomotoso/heritage-beaver.git
cd heritage-beaver

# 2. Database
# Option A — local Postgres you already run (no Docker)
createdb hearthside  # or heritage_weaver — match DATABASE_URL

# Option B — Docker
docker compose up -d

# 3. Environment
cp backend/.env.example backend/.env
# fill DATABASE_URL, JWT_SECRET — and for email either:
#   Resend: RESEND_API_KEY + EMAIL_FROM (verify a domain at resend.com/domains)
#   or Gmail: SMTP_HOST/SMTP_USER/SMTP_PASS + EMAIL_FROM (free, App Password)

# 4. Migrate & run
cd backend && go run ./cmd/migrate
go run ./cmd/api
# → heritage backend listening on :8080 (development)

# 5. Frontend
cd ../frontend && npm install && npm run dev
# → http://localhost:3000
```

Health check: `curl http://localhost:8080/health`

### Environment

| Variable | Required | What |
|----------|----------|------|
| `DATABASE_URL` | yes | `postgres://user:pass@localhost:5432/hearthside?sslmode=disable` |
| `JWT_SECRET` | yes | Random 32+ bytes — never commit the real one |
| `RESEND_API_KEY` | email | Resend key (prod) |
| `SMTP_HOST/USER/PASS` | email alt | Gmail SMTP (`smtp.gmail.com`) — free dev alternative |
| `EMAIL_FROM` | email | e.g. `Hearthside <hello@yourdomain>` |
| `APP_BASE_URL` | email | Frontend URL for verification links (`http://localhost:3000`) |
| `LLM_API_URL/LLM_MODEL` | no | Defaults to keyless OpenCode Zen |

See `backend/.env.example` and `frontend/.env.example`.

---

## API

All routes except `POST /auth/*` and `GET /health` require `Authorization: Bearer <token>`.

```
POST   /auth/register            → 201 { user, message } + confirmation email
POST   /auth/verify-email        → { user, token, refresh_token }
POST   /auth/resend-verification
POST   /auth/login               → { user, token, refresh_token } (verified only)
POST   /auth/refresh             → rotating refresh token
POST   /auth/logout

GET    /family/tree
POST   /family-members    GET /family-members    GET/PATCH/DELETE /family-members/{id}
POST   /relationships     GET  /relationships    GET/PATCH/DELETE /relationships/{id}
POST   /stories           GET  /stories          GET/PATCH/DELETE /stories/{id}
POST   /stories/{id}/process-wisdom   GET /wisdom-extracts[?story_id=]
POST   /mythology/chapters  GET /mythology/chapters  GET/PATCH /mythology/chapters/{id}
GET    /users/{id}  PATCH /users/{id}  DELETE /users/{id}   (self only)
GET    /health
```

Ownership is enforced at the repository layer — every read, write, and delete is scoped to the authenticated user.

---

## Project layout

```
hearthside/
├── backend/
│   ├── cmd/api/               # HTTP server
│   ├── cmd/migrate/           # goose runner
│   ├── internal/
│   │   ├── auth/              # JWT + password + refresh tokens
│   │   ├── config/            # env
│   │   ├── domain/            # User, FamilyMember, Story, Wisdom, MythChapter …
│   │   ├── http/              # router, middleware, handlers
│   │   ├── llm/               # OpenCode client (keyless)
│   │   ├── mail/              # Resend + Gmail SMTP
│   │   └── store/postgres/    # repositories
│   └── migrations/            # 0001_init … 0004_email_verification
├── frontend/
│   ├── app/(auth)/            # login, register, verify-email
│   ├── app/dashboard/         # hearth, family, relationships, stories, wisdom, mythology
│   ├── lib/                   # api, auth, config, types
│   └── store/                 # Zustand auth store
└── docs/                      # architecture, implementation plan
```

---

## Documentation

- [System architecture](docs/ARCHITECTURE.md)
- [Implementation plan](docs/IMPLEMENTATION_PLAN.md)

MIT.
