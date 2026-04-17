# Tutorial: Menjalankan Wangsa di Lokal

## Prasyarat

Pastikan sudah terinstall:

| Tool | Cek | Install |
|------|-----|---------|
| Go 1.21+ | `go version` | https://go.dev/dl |
| Node.js 18+ | `node --version` | https://nodejs.org |
| PostgreSQL 14+ | `psql --version` | https://postgresql.org/download |
| Git | `git --version` | https://git-scm.com |

---

## Langkah 1 — Extract dan masuk ke folder

```bash
# Setelah extract ZIP
cd wangsa_app
```

Struktur yang harus ada:
```
wangsa_app/
├── backend/
├── frontend/
├── Makefile
└── TUTORIAL_LOKAL.md
```

---

## Langkah 2 — Setup PostgreSQL

### Opsi A: PostgreSQL sudah terinstall di komputer

```bash
# Buat database dan user (jalankan sekali)
make db-create

# Ini setara dengan:
# createdb wangsa
# psql wangsa -c "CREATE USER wangsa WITH PASSWORD 'wangsa';"
# psql wangsa -c "GRANT ALL PRIVILEGES ON DATABASE wangsa TO wangsa;"
```

### Opsi B: Pakai Docker (tidak perlu install PostgreSQL)

```bash
docker run -d \
  --name wangsa-postgres \
  -e POSTGRES_DB=wangsa \
  -e POSTGRES_USER=wangsa \
  -e POSTGRES_PASSWORD=wangsa \
  -p 5432:5432 \
  postgres:16-alpine

# Verifikasi jalan:
docker ps | grep wangsa-postgres
```

### Opsi C: Homebrew (macOS)

```bash
brew install postgresql@16
brew services start postgresql@16
make db-create
```

---

## Langkah 3 — Setup environment backend

```bash
cd backend
cp .env.example .env
```

Isi `.env` untuk development lokal (default sudah cocok):

```env
PORT=8080
DATABASE_URL=postgres://wangsa:wangsa@localhost:5432/wangsa?sslmode=disable
JWT_SECRET=dev-secret-min-32-chars-change-in-prod-ok
ACCESS_TOKEN_MINUTES=15
COOKIE_DOMAIN=localhost
PRODUCTION=false
FRONTEND_URL=http://localhost:5173
```

> **Catatan:** `JWT_SECRET` minimal 32 karakter. Server akan menolak start jika kurang.

---

## Langkah 4 — Install dependencies

```bash
# Dari root folder (bukan backend/)
make install

# Ini akan:
# 1. go mod tidy && go mod download (generate go.sum)
# 2. npm install di folder frontend/
```

Jika `make` tidak tersedia (Windows), jalankan manual:

```bash
cd backend && go mod tidy && go mod download && cd ..
cd frontend && npm install && cd ..
```

---

## Langkah 5 — Jalankan

Buka **dua terminal**:

### Terminal 1 — Backend

```bash
make dev-backend
# atau manual:
cd backend && go run ./cmd/server/main.go
```

Output yang diharapkan:
```
✅  Database migrations applied
🚀  Wangsa API listening on http://localhost:8080
🌐  CORS allowed for: http://localhost:5173
```

### Terminal 2 — Frontend

```bash
make dev-frontend
# atau manual:
cd frontend && npm run dev
```

Output yang diharapkan:
```
  VITE v5.x.x  ready in 300 ms
  ➜  Local:   http://localhost:5173/
```

---

## Langkah 6 — Verifikasi

1. Buka **http://localhost:5173**
2. Klik **Daftar** → isi form → Register
3. Pengguna pertama otomatis jadi **Super Admin**
4. Cek tab Network di browser DevTools → cookie `wangsa_token` harus ada

### Test API langsung (opsional)

```bash
# Health check
curl http://localhost:8080/health
# → {"status":"ok","service":"wangsa","db":"connected"}

# Register
curl -c cookies.txt -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Test User","email":"test@example.com","password":"password123"}'
```

---

## Troubleshooting

### ❌ "JWT_SECRET must be at least 32 characters"

Edit `backend/.env` dan ganti `JWT_SECRET` dengan string minimal 32 karakter.
Generate: `openssl rand -hex 32`

### ❌ "Database connection failed" / "dial tcp refused"

PostgreSQL belum jalan. Cek:
```bash
# macOS/Linux
pg_isready -h localhost -p 5432

# Docker
docker ps | grep postgres

# Jika Docker container mati:
docker start wangsa-postgres
```

### ❌ "role 'wangsa' does not exist"

Buat ulang user:
```bash
psql postgres -c "CREATE USER wangsa WITH PASSWORD 'wangsa';"
psql postgres -c "CREATE DATABASE wangsa OWNER wangsa;"
```

### ❌ "go: command not found"

Go belum terinstall atau belum di PATH. Download dari https://go.dev/dl lalu:
```bash
# Tambah ke ~/.bashrc atau ~/.zshrc:
export PATH=$PATH:/usr/local/go/bin
source ~/.zshrc
```

### ❌ Frontend 404 saat akses `/silsilah` langsung

Ini normal di development karena Vite dev server sudah handle routing.
Di production pastikan `vercel.json` atau `_redirects` sudah ada (sudah ada di project ini).

### ❌ Cookie tidak tersimpan / selalu redirect ke login

Pastikan:
- `PRODUCTION=false` di `.env` (agar Secure flag tidak aktif di HTTP)
- Browser tidak block third-party cookie
- Akses via `http://localhost:5173`, bukan IP langsung

### ❌ "pq: SSL is not enabled on the server"

Ganti connection string:
```env
DATABASE_URL=postgres://wangsa:wangsa@localhost:5432/wangsa?sslmode=disable
```
(pastikan ada `?sslmode=disable`)

### ❌ go-sqlite3 error / CGO

Proyek ini sudah migrasi ke PostgreSQL (`lib/pq`). Tidak butuh CGO atau GCC.
Pastikan tidak ada sisa `go-sqlite3` di go.mod.

---

## Reset database (jika perlu mulai ulang)

```bash
make db-reset
# Ini drop database lalu buat ulang dari awal
# Schema di-apply otomatis saat backend pertama start
```

---

## Hot reload (opsional, lebih nyaman development)

### Backend dengan Air

```bash
go install github.com/cosmtrek/air@latest
cd backend
air  # baca konfigurasi dari .air.toml jika ada, otherwise auto-detect
```

### Atau tambah `.air.toml` di folder backend:

```toml
root = "."
cmd = "go build -o ./tmp/server ./cmd/server/main.go"
bin = "./tmp/server"
include_ext = ["go"]
exclude_dir = ["tmp", "vendor"]
```

---

## Struktur folder penting

```
backend/
├── cmd/server/main.go          ← Entry point
├── config/config.go            ← Baca env vars
├── internal/
│   ├── delivery/http/
│   │   ├── router.go           ← Semua route
│   │   ├── handler/            ← HTTP handlers
│   │   └── middleware/         ← Auth, rate limit, logger
│   ├── domain/                 ← Struct + interfaces
│   ├── usecase/                ← Business logic
│   ├── repository/             ← SQL queries
│   └── pkg/
│       ├── database/
│       │   ├── postgres.go     ← Koneksi DB + migration
│       │   └── schema.sql      ← DDL (embedded ke binary)
│       ├── jwt/                ← Access + refresh token
│       ├── response/           ← Standar JSON response
│       └── sanitize/           ← Validasi input (XSS prevention)

frontend/src/
├── api/axios.js                ← HTTP client + auto-refresh interceptor
├── stores/                     ← Pinia state (auth, family, kas, events, dll)
├── composables/                ← useToast, useDirtyGuard, useZoomPan
├── utils/                      ← format.js (fmtRp, fmtDate), sanitize.js
├── components/                 ← Reusable UI components
├── views/                      ← Halaman (Dashboard, Silsilah, Kas, dll)
└── router/index.js             ← Vue Router + session guard
```

---

## Make commands tersedia

```bash
make help          # Lihat semua command
make install       # Install semua dependency
make dev-backend   # Jalankan backend (port 8080)
make dev-frontend  # Jalankan frontend (port 5173)
make db-create     # Buat database lokal
make db-reset      # Drop + buat ulang database
make build         # Build production binary + frontend dist
make vet           # Go static analysis
make fmt           # Format Go code
make clean         # Hapus build artifacts
```
