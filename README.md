# TK2 Basdat Backend

Backend service for a service-marketplace application (home-service booking). It powers user onboarding, service discovery, ordering, payments, worker job flow, and feedback — designed to showcase strong database-driven features for recruiters.

## Why this project stands out
- **End-to-end marketplace flow**: authentication, service browsing, ordering, payments, and reviews in one backend.
- **Rich business logic**: wallet (MyPay) with top-up/transfer/withdrawal, voucher & promo discounts, and worker job lifecycle.
- **Relational data focus**: PostgreSQL with UUID identifiers and transactional operations across users, orders, and payments.
- **Production-style API surface**: clear HTTP endpoints for each domain (users, jobs, payments, testimonials).

## Key Features
- **Auth & user management**: register, login, get/update user profile.
- **Service catalog**: homepage data, subkategori, and ordering flow.
- **Jobs for workers**: available jobs, pick job, update job status.
- **MyPay wallet**: balance, history, top-up, transfer, withdrawal, and payment processing.
- **Promotions**: voucher & promo listing and purchase.
- **Testimonials**: create, list, delete feedback.

## Tech Stack
- **Go 1.21** using `net/http`
- **PostgreSQL** via `lib/pq`
- **UUID** identifiers via `google/uuid`

## API Highlights
Some of the main endpoints (see `main.go` for full list):
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
This project demonstrates building a real-world backend with complex data flows: multi-role users, transactional payments, and marketplace operations. It’s a solid showcase of practical backend engineering with relational data, APIs, and business logic.
