# TK2 Basdat Backend

This is the backend I built for a service‑marketplace application (home‑service booking). I wrote it in Go and focused on a clean, straightforward backend that I can explain clearly.

## Tech Stack & Approach
- **Go 1.21** using `net/http`
- **PostgreSQL** with **raw SQL** (no ORM)
- **database/sql** + `lib/pq` for queries
- **UUID** identifiers via `google/uuid`

What I want you to notice:
- I write the SQL myself to control the schema and queries.
- I keep the API simple and clear, without heavy frameworks.

## Setup & Run
1. **Install Go 1.21+**
2. **Configure PostgreSQL**
   - Update the connection constants in `main.go` (`host`, `port`, `user`, `password`, `dbname`).
3. **Install dependencies**
   ```bash
   go mod download
   ```
4. **Run the server**
   ```bash
   go run main.go
   ```
5. Server runs at **http://localhost:8080**

## Notes for Recruiters
I can walk through the architecture, the database design, and how I wire up endpoints to SQL. The main point is that I understand the backend deeply because I built it without hiding the database behind an ORM.
