# TK2 Basdat Backend

This is the backend I built for a service‑marketplace application (home‑service booking). It is a single Go service that uses raw SQL with PostgreSQL (no ORM).

## Tech Stack
- **Go 1.21** using `net/http`
- **PostgreSQL**
- **database/sql** + `lib/pq`
- **UUID** identifiers via `google/uuid`

## Architecture
- Single Go server in `main.go`.
- Opens one PostgreSQL connection on startup and shares it across handlers.
- Routes are registered with `http.HandleFunc` and wrapped by a simple CORS middleware.
- Handlers decode JSON, run SQL queries, and return JSON responses.

## Database Design (High Level)
- **Users & roles**: `user`, `pelanggan`, `pekerja`.
- **Service catalog**: `kategori_jasa`, `subkategori_jasa`, `sesi_layanan`, `pekerja_kategori_jasa`.
- **Order flow**: `tr_pemesanan_jasa`, `tr_pemesanan_status`, `status_pesanan`.
- **Wallet & payments**: `tr_mypay`, `kategori_tr_mypay`.
- **Promotions & feedback**: `voucher`, `promo`, `diskon`, `testimoni`.

## Endpoints and SQL
All routes are defined in `main.go`. Each endpoint maps directly to SQL queries in its handler using `QueryRow`, `Query`, or `Exec`. There is no ORM layer, so the SQL is written explicitly in the handlers.

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
