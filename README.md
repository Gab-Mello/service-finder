# Service Finder

A Go backend API for a service marketplace that connects customers with independent providers.

> Implemented as a college project. The scope reflects that: data and sessions are stored in memory, reset on restart, and the code prioritizes clarity and backend structure over production hardening.

## Features

- User registration and session-based authentication for two roles: **providers** and **customers**
- Service postings with search by city, district, and category
- Order/booking lifecycle: `PENDENTE → ACEITO → EM_ANDAMENTO → CONCLUIDO` (with `CANCELADO` as a terminal state)
- Reviews and ratings left by customers after a completed order
- Provider profiles with expertise, location, contact, and bio

## Tech Stack

- **Language:** Go 1.25.3
- **HTTP:** standard library `net/http`
- **API docs:** Swagger / OpenAPI 2.0 via [`swaggo/swag`](https://github.com/swaggo/swag) + [`swaggo/http-swagger`](https://github.com/swaggo/http-swagger)
- **IDs:** [`google/uuid`](https://github.com/google/uuid)
- **Storage:** in-memory repositories (no database)
- **Architecture:** clean architecture-inspired structure with domain packages under `internal/{user,posting,order,review}`, HTTP handlers under `internal/http`, and domain interfaces under `internal/ports`

## Project Structure

```
service-finder/
├── cmd/api/            # Entry point (main.go)
├── docs/               # Generated Swagger files (swagger.yaml/json)
└── internal/
    ├── auth/           # Session management & password handling
    ├── user/           # User domain
    ├── posting/        # Service posting domain
    ├── order/          # Order/booking domain
    ├── review/         # Review/rating domain
    ├── http/           # HTTP server, handlers, routes, middleware
    └── ports/          # Domain interfaces
```

## Getting Started

Requirements: Go 1.25 or later.

```bash
go mod download
go run ./cmd/api
```

The server listens on `http://localhost:8080`. No environment variables are required.

Swagger UI is available at:

```
http://localhost:8080/swagger/
```

## API Overview

All endpoints are under the base path `/api/v1`. See Swagger for full request/response schemas.

**Auth & Users**
- `POST /users` — register a new user
- `POST /login` / `POST /logout`
- `GET /me` — current user profile
- `PATCH /providers/profile` — update provider profile

**Postings**
- `GET /postings` — search public listings
- `POST /postings` — create (provider only)
- `GET /postings/{id}` / `PATCH /postings/{id}`
- `GET /postings/mine` — provider's own postings
- `POST /postings/{id}/archive`

**Orders**
- `POST /orders` — request a service
- `GET /orders/mine` / `GET /orders/{id}`
- `POST /orders/{id}/accept` · `/start` · `/complete` · `/cancel`

**Reviews**
- `POST /reviews` — create after order is completed
- `PATCH /reviews/{orderId}` — edit within the edit window

**Utility**
- `GET /healthz` — health check

## Notes

- Sessions expire after 5 minutes.
- All data is held in memory and is lost when the process restarts.
