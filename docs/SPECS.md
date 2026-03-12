# Monolith Go Framework — Full Specification

A modern fullstack Go framework inspired by Django, Spring Boot, Laravel, Phoenix, and Ruby on Rails. **Magical but fast**: minimal boilerplate, convention over configuration, and a single coherent spec for backend, frontend, CLI, and deployment.

---

## 1. Overview & Principles

| Principle | Description |
|-----------|-------------|
| **Go-only backend** | No Node in the backend; created via framework CLI with optional React UI. |
| **Entry: main.go** | No `init()` for routes. Entry point is `main.go` at project root; routes registered there or in imported packages. |
| **Handler return: Result** | Handlers return a **Result** type: `Ok(value)` for success, `Err(statusCode, message)` for errors. No raw `error`. |
| **Handler args** | Every handler receives `(app App, db DB, integrations Integrations)`. App = request/response, auth, cache, SSE, WS. DB = ORM wrapper. Integrations = 3rd-party (Storage, Payment, Mail, Maps). |
| **Config in one place** | All configuration in **settings.go**; secrets in **.env**. Provider choice (DB, Storage, Payment, Cache, Auth) is magic through settings. |
| **No explicit JSON** | Return `Ok(user)` or `Ok(users)`; framework serializes to JSON. Return `Ok("text")` for plain text. |

---

## 2. Project Structure & Dev Workflow

- **Root**: `main.go`, `settings.go`, `.env` (gitignored), `go.mod`, established folders (`cmd/`, `internal/`, `app/` for frontend, etc.).
- **Backend**: Run with **Air** (hot reload). `air` or `make dev` starts the Go server.
- **Frontend**: Optional React app under `app/`. In dev, run in **another terminal** (e.g. `cd app && npm run dev`). Backend and frontend are separate processes.
- **CLI** (when creating project): Prompt **"Add React UI? (y/n)"** to scaffold the `app/` with Vite, React, Tailwind, Zustand.

---

## 3. Handlers & Result Type

- **Signature**: `func(app App, db DB, integrations Integrations) Result`.
- **Return**:
  - `Ok(value)` — value can be `string` (→ text/plain), struct or slice (→ JSON). Framework sets 200 and body.
  - `Err(statusCode, message)` — framework sets HTTP status and error body.
- **Example**:

```go
func main() {
	Get("/health", func(app App, db DB, integrations Integrations) Result {
		return Ok("ok")
	})
	Get("/users", func(app App, db DB, integrations Integrations) Result {
		repo := NewUserRepo(db)
		users, _ := repo.List(app.Request().Context())
		return Ok(users)
	})
	Get("/me", func(app App, db DB, integrations Integrations) Result {
		app.Request().AllowedRoles = []string{"user", "admin"}
		user := app.Request().Auth()
		if user == nil {
			return Err(401, "Unauthorized")
		}
		return Ok(user)
	})
	// framework runs server from here
}
```

---

## 4. Database (ORM)

- **DB is an ORM wrapper**, not the raw driver. Default: **SQLite**; Postgres/MySQL via settings.
- **Repos**: Create repositories that take `DB` (e.g. `NewUserRepo(db)`). Inside the handler, pass `db` to the repo; repo uses ORM for queries.
- **CRUD**: ORM exposes Create, GetByID, Update, Delete, List (with filters) for models.
- **Joins**: Recursively **optimal left joins** to avoid N+1 (e.g. load user with orders in one go).
- **Raw SQL**: `db.Exec(ctx, rawSQL, args...)`, `db.Query`, `db.QueryRow` when needed.
- Config (DSN, driver) in **settings.go**; secrets (e.g. DB URL) in **.env**.

---

## 5. Integrations (Provider Interfaces)

All 3rd-party access via the **`integrations`** handler arg. Each area is an **interface** with multiple implementations; provider chosen in **settings**.

### 5.1 Storage

- **Interface**: same API for all providers.
- **Implementations**:
  - **Local** (default): folder on disk; path in settings.
  - **S3**: bucket, region, credentials from settings/.env.
  - **Cloudinary**: cloud name, API key/secret from settings/.env.
- **API**: `integrations.Storage().Presign(bucket, key, ttl)`, `Put`, `Get`, etc.

### 5.2 Payment

- **Interface**: `Pay(details)`, `Receive()` (webhooks).
- **Implementations**: **Stripe**, **Flutterwave**, **Pesapal**. Provider in settings.

### 5.3 Mail, Maps

- **Mail**: send transactional email; provider/config in settings.
- **Maps**: Google Maps (or other); config in settings.

---

## 6. Settings & Secrets

- **settings.go**: Single file (or package) that defines all config: DB, Storage provider/paths, Payment provider, Cache (Redis), Auth (Google, JWT), feature flags.
- **.env**: Secrets only (API keys, DB URL, JWT secret, Google client secret). Loaded by settings; never committed. Framework uses env for production (e.g. Render).
- **Demo** (all options): DB.Driver, DB.DSN | Storage.Provider (local|s3|cloudinary), LocalPath, S3.*, Cloudinary.* | Payment.Provider (stripe|flutterwave|pesapal) | Cache.Backend (redis|memory), Redis.URL | Auth.Google.ClientID/ClientSecret, Auth.JWT.Secret. All secrets from .env.

---

## 7. Auth

- **Google SSO ships by default**. Auth routes: `auth.Signup`, `auth.Login`, `auth.Logout`. Google OAuth is built in; config (client ID/secret) in settings/.env.
- **JWT** for API auth. **Role-based access**: per endpoint, set **`app.Request().AllowedRoles = []string{"admin", "user"}`** or **`app.Request().AllowedUserIDs = []string{"id1", "id2"}`**. Framework validates JWT and enforces role/user before or after handler.
- Current user: **`app.Request().Auth()`** (nil if unauthenticated).

---

## 8. Caching

- **Redis** for caching (config in settings). In-memory fallback if Redis not set.
- **API**: `app.Response().EnableCache(TTLSeconds, cacheType)` where cacheType comes from settings (e.g. `"redis"`).

---

## 9. CLI

- **Create project**: Option to add React UI (y/n).
- **Artisan-style commands**:
  - **create feature** — scaffold feature (handlers, routes, optional repo).
  - **create component** — scaffold React component in `app/`.
  - **sync-client** — sync backend types to `app/src/types.ts`.
  - **benchmark** — run benchmark suite (compare with other frameworks).
- Custom user commands can be registered with the CLI.

---

## 10. Frontend (SPA)

- **Location**: `app/`. **Next.js-style** structure: file-based or folder-based routes, layouts, pages.
- **Stack**: Vite, React, Tailwind, Zustand. Client-side only.
- **Landing**: Framework template includes an **Airbnb-class landing page** (polished, modern).
- **Type sync**: Structs embedding `Base{}` are synced to `app/src/types.ts` via **sync-client** or dev watcher.
- **Deployment**: Build (`npm run build`); serve static assets from backend or CDN. **Render**: backend as web service; static site or same service for SPA.

---

## 11. Deployment (Render, SPA)

- **Backend**: Deploy as a web service (e.g. Render). Start command runs the compiled Go binary; env vars from Render.
- **SPA**: Build frontend; serve from backend’s static route or as a static site. Document Render build/start and env for both.

---

## 12. Benchmarking

- Framework provides a **benchmark** CLI command.
- **Comparison targets**: **Phoenix**, **Django**, **Spring Boot**, **Laravel**, **Ruby on Rails** (requests/sec, latency).
- Docs include a **benchmarking section** with results (or placeholder table) so users see where the stack sits relative to these frameworks.

---

## 13. Quick Reference

| Topic | Spec |
|-------|------|
| Entry | `main.go` at root; no `init()` for routes |
| Handler | `func(app App, db DB, integrations Integrations) Result` |
| Success | `return Ok(value)` (string → text, struct/slice → JSON) |
| Error | `return Err(statusCode, message)` |
| DB | ORM wrapper; repos use `db`; CRUD, left joins, `db.Exec(raw)` |
| Storage | Interface: local (default), S3, Cloudinary — settings |
| Payment | Interface: Stripe, Flutterwave, Pesapal — settings |
| Config | settings.go + .env for secrets |
| Auth | Google SSO default; JWT; AllowedRoles / AllowedUserIDs |
| Cache | Redis (settings); EnableCache(TTL, cacheType) |
| CLI | create feature, create component, sync-client, benchmark; React UI option |
| Dev | Air (backend) + frontend in another terminal |
| Deploy | Render; SPA static + backend service |

---

## 14. Benchmarking Comparison (Target Frameworks)

| Framework | Language | Typical use case | Relative focus in spec |
|-----------|----------|------------------|-------------------------|
| **Phoenix** | Elixir | Concurrency, real-time | SSE, WebSockets, performance |
| **Django** | Python | Batteries-included, ORM, admin | ORM, auth, settings, CLI |
| **Spring Boot** | Java | Enterprise, ecosystem | Structure, profiles, integrations |
| **Laravel** | PHP | DX, Artisan, Blade | CLI, migrations, integrations, auth |
| **Ruby on Rails** | Ruby | Convention, rapid dev | Convention, magical defaults, landing |

The Monolith Go framework aims to match or exceed these in **developer experience** (magic, conventions, CLI) and **performance** (Go + optional benchmarks vs Phoenix/Django/Spring/Laravel/Rails). The **benchmark** command and docs provide reproducible comparisons (requests/sec, p50/p99 latency) against representative apps in each stack.
