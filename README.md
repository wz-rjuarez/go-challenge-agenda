# Clinic Scheduling Challenge

A monorepo with two services that power a medical scheduling system.

## Architecture

```
services/agenda   — internal gRPC service, owns all scheduling logic and data
services/api      — public HTTP REST API (Gin), calls agenda via gRPC
```

Both services share a single Go module. Proto definitions live in `proto/` and generated code in `gen/`.

## Quick start

```bash
# Generate proto code (requires buf)
make proto

# Run with Docker Compose (SQLite)
make docker-up

# Run with PostgreSQL
make docker-up-postgres

# API available at http://localhost:8080
# Agenda gRPC at localhost:50051
```

## Database Configuration

The system supports both SQLite and PostgreSQL, switchable via environment variables:

**SQLite (default):**
```bash
DB_DRIVER=sqlite3
DB_SOURCE=agenda.db
```

**PostgreSQL:**
```bash
DB_DRIVER=postgres
DB_SOURCE="host=localhost user=agenda password=agenda_pass dbname=agenda port=5432 sslmode=disable"
```

Both databases share the same business logic and repository interfaces. No code changes required to switch between them.

## API routes

| Method | Path | Description |
|--------|------|-------------|
| GET | /v1/doctors | List all doctors |
| GET | /v1/doctors/:id | Get doctor with working hours |
| GET | /v1/doctors/:id/availability?date=YYYY-MM-DD&type=first_visit | Get available slots |
| POST | /v1/reservations | Create a reservation |
| GET | /v1/reservations/:id | Get reservation |
| GET | /v1/reservations?doctor_id=&from=&to= | List reservations |
| DELETE | /v1/reservations/:id | Cancel reservation |
| GET | /v1/users | List all users |
| POST | /v1/users | Create a user |
| GET | /v1/users/:id | Get user |
| PATCH | /v1/users/:id | Update user |
| DELETE | /v1/users/:id | Delete user |
| GET | /v1/users/:id/reservations | List reservations for a user |

## Reservation types

- `first_visit` — 60 minute block
- `follow_up` — 30 minute block

## Tasks

Work through as many issues as you can. They are loosely ordered from more concrete to more open-ended.

- [ ] **Issue #1** — `GET /v1/doctors/:id/availability` returns hardcoded stub data. Wire it to the real usecase.
- [ ] **Issue #2** — There is a bug in reservation conflict detection that allows double-booking in certain cases. Find and fix it.
- [ ] **Issue #3** — `ListReservations` in the SQLite repository is not implemented. Implement it.
- [ ] **Issue #4** — gRPC errors are not properly mapped to HTTP status codes. Fix the mapping in `pkg/errcodes`.
- [ ] **Issue #5** — There are failing tests in the codebase. Make them pass.
- [ ] **Issue #6** — `GET /v1/users/:id/reservations` returns `501 Not Implemented`. Implement it end-to-end across all layers.
- [ ] **Issue #7** — The `api` service usecases are coupled to a concrete type where an interface should be used. Fix it.
- [ ] **Issue #8** — Blocked slots are stored and retrieved, but recurrences are never expanded. Implement recurrence expansion (daily / weekly / monthly until a given date).
- [ ] **Issue #9** — The system only supports SQLite. Add Postgres support, switchable via `DB_DRIVER=postgres` environment variable, without modifying business logic.
- [ ] **Issue #10** — Some HTTP handlers bypass the usecase layer. Fix the layering.
- [ ] **Issue #11** — Blocked slots are not taken into account when computing availability. Fix this.
- [ ] **Issue #12** — Extend the system to support additional service types (labs, therapy) with different slot durations. Design the domain model, update the proto, and expose them through the existing services.
- [ ] **Issue #13** — Add structured logging and at least one meaningful metric without coupling observability concerns to the domain layer.
- [ ] **Issue #14 (optional)** — Consider whether the system would benefit from a dedicated third service. If so, propose the service boundary, define its proto contract, and implement it.
- [ ] **Issue #15** — Write a short ADR (Architecture Decision Record) for one technical decision you made or changed during this challenge.

## Pre-seeded doctors

| ID | Name | Specialty | Working days |
|----|------|-----------|-------------|
| doc-001 | Dr. Ana García | General Practice | Mon–Fri (09–17, Fri until 13) |
| doc-002 | Dr. Luis Mendoza | Cardiology | Mon, Wed, Fri (08–16, Fri until 12) |
| doc-003 | Dr. Sara Patel | Pediatrics | Tue, Thu (10–18) |
