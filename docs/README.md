# Monolith Go Framework — Documentation

A Go-only web framework inspired by Django, Spring Boot, Laravel, Phoenix, and Ruby on Rails. **Magical but fast**: minimal boilerplate, handler returns `Result`, ORM with repos, provider-based integrations (Storage, Payment), config in settings + .env.

**Full specification**: see **[SPECS.md](./SPECS.md)** for the complete spec, CLI, deployment, and benchmarking vs Phoenix, Django, Spring Boot, Laravel, Rails.

---

## Table of contents

1. [Overview](#1-overview)
2. [Getting started](#2-getting-started)
3. [Routing and handlers](#3-routing-and-handlers)
4. [Handler capabilities](#4-handler-capabilities)
5. [Auth](#5-auth)
6. [Frontend app structure](#6-frontend-app-structure)
7. [Type sync (backend → client)](#7-type-sync-backend--client)
8. [Integrations](#8-integrations)
9. [Settings](#9-settings)
10. [CLI, deployment, benchmarking](#10-cli-deployment-benchmarking)

---

## 1. Overview

- **Stack**: Go only (backend). Created via `npx create-go-app@latest`.
- **Design**: “Magical but fast” — concise API, no `app.run`, no explicit JSON. Return Go types; the framework serializes.
- **Frontend**: Lives under `app/` — Vite, React, Tailwind, Zustand; client-side only.

### Principles

| Principle | Meaning |
|-----------|--------|
| Entry | **main.go** at root (no `init()` for routes). Use **Air** for backend dev; frontend in another terminal. |
| Handler return | **Result** type: `Ok(value)` for success, `Err(status, message)` for errors. |
| Handler args | `(app App, db DB, integrations Integrations)` — DB is ORM wrapper; integrations = Storage, Payment, Mail, Maps. |
| Config | **settings.go** + **.env** for secrets. All provider choice and magic in one place. |
| Built-in auth | **Google SSO by default**; JWT with **AllowedRoles** / **AllowedUserIDs** per endpoint. |
| Type sync | Structs embedding `Base{}` → `app/src/types.ts` via **sync-client** or dev watcher. |

### Minimal example

```go
// main.go at project root
func main() {
	Get("/health", func(app App, db DB, integrations Integrations) Result {
		return Ok("ok")
	})
	// framework runs server
}
```

---

## 2. Getting started

1. Create the app via the framework CLI; choose **Add React UI? (y/n)**.
2. Entry is **main.go** at root; register routes there (or in packages main imports). Use **Air** to run the backend; run the frontend in another terminal (`cd app && npm run dev`).
3. Handlers return **Result**: `Ok(value)` or `Err(statusCode, message)`. Config and secrets in **settings.go** and **.env**.

---

## 3. Routing and handlers

### API

- **Verbs**: `Get`, `Post`, `Put`, `Delete`, etc. — e.g. `Get(pathname, handler)`.
- **Handler**: `func(app App, db DB, integrations Integrations) Result`. Receive **App**, **DB** (ORM wrapper), **Integrations**. Return **Result** via `Ok(value)` or `Err(status, message)`.
- **Groups**: `All("/prefix", func() { ... })` with nested `Get`, `Post`, etc.

### Database (ORM)

- **DB is an ORM wrapper** (not raw DB). Default: **SQLite**; override in settings. Create **repos** that take `db`; use ORM CRUD, optimal left joins, and `db.Exec(ctx, rawSQL, ...)` for raw SQL. Config in settings; secrets in .env.

### Examples

```go
Get("/ping", func(app App, db DB, integrations Integrations) Result {
	return Ok("pong")
})

All("/api/v1", func() {
	Get("/", func(app App, db DB, integrations Integrations) Result {
		app.Response().EnableCache(78, "redis")
		return Ok("Hi")
	})
	Get("/users", func(app App, db DB, integrations Integrations) Result {
		repo := NewUserRepo(db)
		users, _ := repo.List(app.Request().Context())
		return Ok(users)
	})
	Post("/users", func(app App, db DB, integrations Integrations) Result {
		// create user; return Ok(user) or Err(...)
		return Ok(user)
	})
})
```

---

## 4. Handler capabilities

### Response

- **Handlers return `Result`**: `Ok(value)` for success (string → text/plain, struct/slice → JSON), `Err(statusCode, message)` for errors.
- **Cache**: **Redis** by default (settings). `app.Response().EnableCache(TTLSeconds, cacheType)`.
- **Role-based access**: set `app.Request().AllowedRoles = []string{"admin","user"}` or `app.Request().AllowedUserIDs = []string{"id1"}`; framework enforces JWT and role/user.

### WebSockets

- Enable in a handler with `app.EnableWebsockets()`; framework handles upgrade and lifecycle.
- **Helpers** (in WS-enabled handler or connection lifecycle):

| Helper | Usage |
|--------|--------|
| `app.WS().Send(conn, data)` | Send to this connection. |
| `app.WS().Broadcast(channel, data)` | Send to all in a room/channel. |
| `app.WS().Subscribe(conn, channel)` | Subscribe connection to a channel. |
| `app.WS().OnMessage(conn, fn)` | Handle incoming message. |
| `app.WS().OnConnect(conn, fn)` | On connect. |
| `app.WS().OnDisconnect(conn, fn)` | On disconnect. |

### Code examples

```go
Get("/api/users", func(app App, db DB, integrations Integrations) Result {
	app.Response().EnableCache(300, "redis")
	users := []User{}
	return Ok(users)
})

Get("/api/config", func(app App, db DB, integrations Integrations) Result {
	app.Response().EnableCache(300, "redis")
	return Ok(Config{Theme: "dark"})
})

Get("/ping", func(app App, db DB, integrations Integrations) Result {
	return Ok("pong")
})

// Role-based: only these roles allowed
Get("/admin", func(app App, db DB, integrations Integrations) Result {
	app.Request().AllowedRoles = []string{"admin"}
	return Ok("Admin only")
})

Get("/ws", func(app App, db DB, integrations Integrations) Result {
	app.EnableWebsockets()
	return Ok(nil)
})

// WS helpers (in connection lifecycle)
app.WS().Send(conn, []byte(`{"type":"pong"}`))
app.WS().Broadcast("room:lobby", []byte(`{"event":"update"}`))
app.WS().Subscribe(conn, "room:lobby")
app.WS().OnMessage(conn, func(msg []byte) { /* handle */ })
app.WS().OnConnect(conn, func() { /* join default room */ })
app.WS().OnDisconnect(conn, func() { /* cleanup */ })
```

---

## 5. Auth

- **Google SSO ships by default**. Use **auth.Signup**, **auth.Login**, **auth.Logout**; Google OAuth is built in. Config (client ID/secret) in settings and **.env**.
- **JWT** for API. **Role-based access**: set **`app.Request().AllowedRoles = []string{"admin", "user"}`** or **`app.Request().AllowedUserIDs = []string{"id1"}`** per endpoint; framework validates JWT and enforces.
- Current user: **`app.Request().Auth()`** (nil when unauthenticated).

### Using auth in handlers

```go
Get("/me", func(app App, db DB, integrations Integrations) Result {
	app.Request().AllowedRoles = []string{"user", "admin"}
	user := app.Request().Auth()
	if user == nil {
		return Err(401, "Unauthorized")
	}
	return Ok(user)
})

Get("/dashboard", func(app App, db DB, integrations Integrations) Result {
	user := app.Request().Auth()
	if user != nil {
		return Ok(Dashboard{Greeting: "Hello, " + user.Name})
	}
	return Ok("Please sign in")
})
```

---

## 6. Frontend app structure

- **Location**: All UI under **`app/`**. **Next.js-style** structure (file/folder routes, layouts). Target an **Airbnb-class landing** (polished, modern).
- **Stack**: **Vite**, **React**, **Tailwind**, **Zustand**. Client-side SPA; deploy static build (e.g. Render). See [SPECS.md](./SPECS.md) and deployment.

### Generated app folder format

```
app/
  index.html
  package.json
  vite.config.ts
  tailwind.config.js
  tsconfig.json
  public/
  src/
    main.tsx
    App.tsx
    types.ts            # generated from backend (see Type sync)
    layout/
      Layout.tsx
      Header.tsx
    pages/
      Home.tsx
      dashboard/
        Dashboard.tsx
    components/
      ui/
      forms/
    stores/              # Zustand
      authStore.ts
      appStore.ts
  app/                   # optional Next-like route folder
    layout.tsx
    page.tsx
    dashboard/
      page.tsx
```

### Zustand + Tailwind

```tsx
// stores/authStore.ts
import { create } from 'zustand';

export const useAuthStore = create((set) => ({
  user: null,
  setUser: (user) => set({ user }),
  logout: () => set({ user: null }),
}));

// components/Header.tsx
export function Header() {
  const user = useAuthStore((s) => s.user);
  return (
    <header className="flex items-center justify-between p-4 bg-slate-800 text-white">
      <span>{user?.name ?? 'Guest'}</span>
    </header>
  );
}
```

- **`src/types.ts`**: Generated from backend structs; do not hand-edit synced types.

---

## 7. Type sync (backend → client)

- Structs that should be synced **must embed `Base{}`**.
- In **dev**, the framework watches the backend and **generates TypeScript** into **`app/src/types.ts`**. Any change to a monitored struct updates `types.ts`.

### Backend

```go
type User struct {
	Base
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Task struct {
	Base
	Title  string `json:"title"`
	Done   bool   `json:"done"`
	UserID int64  `json:"user_id"`
}
```

### Generated `app/src/types.ts` (example)

```ts
// AUTO-GENERATED; do not edit synced types manually

export interface User {
  id: string;
  name: string;
  email: string;
  created_at: string;
}

export interface Task {
  id: string;
  title: string;
  done: boolean;
  user_id: number;
  created_at: string;
}
```

### Workflow

- Run the app in dev; type-sync watcher runs and scans for structs embedding `Base{}`.
- On change, `app/src/types.ts` is regenerated.

---

## 8. Integrations

- Handlers receive **`integrations`** (with app, db). All 3rd-party access via **`integrations.Storage()`**, **`integrations.Payment()`**, **`integrations.Mail()`**, **`integrations.Maps()`**. Provider choice in **settings**; secrets in **.env**.
- **Storage** — **interface**: **local folder** (default), **S3**, **Cloudinary**. Same API; provider in settings.
- **Payment** — **interface**: **Stripe**, **Flutterwave**, **Pesapal**. `integrations.Payment().Pay(...)`, `.Receive(...)` for webhooks.
- **Mail**, **Maps**: same pattern; config in settings.

### Examples (conceptual)

```go
// Storage — provider (local | s3 | cloudinary) from settings
Get("/upload", func(app App, db DB, integrations Integrations) Result {
	url := integrations.Storage().Presign("my-bucket", "key", 3600)
	return Ok(UploadResponse{UploadURL: url})
})

Post("/invite", func(app App, db DB, integrations Integrations) Result {
	integrations.Mail().Send("user@example.com", "invite", Map{"link": "..."})
	return Ok("OK")
})

Post("/checkout", func(app App, db DB, integrations Integrations) Result {
	result, err := integrations.Payment().Pay(PayDetails{Amount: 1000, Currency: "USD", ...})
	if err != nil {
		return Err(400, err.Error())
	}
	return Ok(result)
})

Get("/events", func(app App, db DB, integrations Integrations) Result {
	app.EnableSSE()
	return Ok(nil)
})
```

---

## 9. Settings

- **All configuration in `settings.go`**; **secrets in .env** (API keys, DB URL, JWT secret, Google client secret). Framework loads at startup. Single source of truth for:
  - **DB**: driver, DSN (from .env)
  - **Storage**: provider (`local` | `s3` | `cloudinary`), local path, S3/Cloudinary keys (from .env)
  - **Payment**: provider (Stripe | Flutterwave | Pesapal), keys from .env
  - **Cache**: backend (`redis` | memory), Redis URL from .env
  - **Auth**: Google client ID/secret, JWT secret from .env
- **Demo**: See [SPECS.md](./SPECS.md) for a full list of options and a benchmarking comparison (Django, Phoenix, Spring Boot, Laravel, Rails).

---

## 10. CLI, deployment, benchmarking

- **CLI**: On create, **Add React UI? (y/n)**. Commands: **create feature**, **create component**, **sync-client** (sync types to `app/src/types.ts`), **benchmark** (compare with Phoenix, Django, Spring Boot, Laravel, Rails).
- **Deployment**: **Render** — backend as web service; SPA built and served as static assets. Use .env / Render env for secrets.
- **Benchmarking**: Framework benchmark command and [SPECS.md](./SPECS.md) compare requests/sec and latency vs top frameworks.

---

## Quick reference

| Need | Use |
|------|-----|
| Entry | **main.go** at root; Air for backend; frontend in another terminal |
| Handler | `func(app App, db DB, integrations Integrations) Result` |
| Success | `return Ok(value)` (string → text, struct/slice → JSON) |
| Error | `return Err(statusCode, message)` |
| DB | ORM wrapper; repos take `db`; CRUD, left joins, `db.Exec(raw)`; config in settings |
| Storage | Interface: **local** (default), **S3**, **Cloudinary** — settings |
| Payment | Interface: Stripe, Flutterwave, Pesapal — settings |
| Config / secrets | **settings.go** + **.env** |
| Cache | **Redis** (settings); `app.Response().EnableCache(ttl, "redis")` |
| Auth | **Google SSO default**; JWT; `app.Request().AllowedRoles`, `AllowedUserIDs` |
| Current user | `app.Request().Auth()` |
| CLI | create feature, create component, sync-client, benchmark; React UI option |
| Deploy | Render; SPA static + backend |
| Full spec | **[SPECS.md](./SPECS.md)** |
