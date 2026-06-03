# TK2 Basdat Backend

This is the backend I built for a service‑marketplace application (home‑service booking). It handles user onboarding, service discovery, ordering, payments, worker job flow, and feedback — a complete, database‑driven system I can confidently explain end‑to‑end.

## Why this project stands out
- **End-to-end marketplace flow**: I covered authentication, service browsing, ordering, payments, and reviews in one backend.
- **Rich business logic**: I implemented a wallet (MyPay) with top‑up/transfer/withdrawal, voucher & promo discounts, and the worker job lifecycle.
- **Relational data focus**: I modeled everything in PostgreSQL with UUIDs and transactional operations across users, orders, and payments.
- **Production-style API surface**: I exposed clear HTTP endpoints for each domain (users, jobs, payments, testimonials).

## Key Features
- **Auth & user management**: register, login, and profile retrieval/updates.
- **Service catalog**: homepage data, subkategori, and ordering flow.
- **Jobs for workers**: available jobs, pick job, update job status.
- **MyPay wallet**: balance, history, top‑up, transfer, withdrawal, and payment processing.
- **Promotions**: voucher & promo listing and purchase.
- **Testimonials**: create, list, delete feedback.

## Tech Stack
- **Go 1.21** using `net/http`
- **PostgreSQL** via `lib/pq`
- **UUID** identifiers via `google/uuid`

## API Highlights
These are some of the main endpoints I expose (see `main.go` for full list):
- `/login`, `/register`, `/getUser`, `/updateUser`
- `/homepage`, `/subkategori`, `/pesan`
- `/mypay/balance`, `/mypay/history`, `/mypay/topup`, `/mypay/transfer`, `/mypay/withdrawal`
- `/jobs/available`, `/jobs/get-job`, `/jobs/job-pekerja-update`
- `/getDiskon`, `/buyVoucher`
- `/createTestimoni`, `/getTestimoni`, `/deleteTestimoni`

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
If we walk through this together, I can explain how each feature maps to real business needs: multi‑role users, transactional payments, and marketplace operations. This project is my strongest demonstration of practical backend engineering with relational data, APIs, and business logic.
