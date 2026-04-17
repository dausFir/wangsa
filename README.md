# 🏠 Wangsa

**Platform digital terpadu untuk keluarga besar** — Silsilah interaktif, kas bersama, kalender acara, dan peta domisili.

---

## 🧰 Tech Stack

| Layer | Teknologi |
|---|---|
| Backend | Go 1.21+, Gin, Clean Architecture |
| Database | **PostgreSQL 14+** |
| Auth | JWT (HttpOnly Cookie, anti-XSS) |
| Frontend | Vue 3, Vite, Pinia, Tailwind CSS |

---

## 📋 Prasyarat (Development)

- Go 1.21+
- Node.js 18+ / npm 9+
- **Database (pilih salah satu):**
  - PostgreSQL 14+ (install langsung)
  - Docker Desktop (untuk menjalankan PostgreSQL di container)

### Install Docker Desktop (untuk Mac)

Jika belum punya Docker, install dulu:

1. **Download:** https://www.docker.com/products/docker-desktop/
2. **Install:** Buka file `.dmg` dan drag ke Applications
3. **Jalankan:** Buka Docker Desktop dari Applications
4. **Cek:** Buka terminal dan ketik `docker --version`

Alternatif via Homebrew:
```bash
brew install --cask docker
# Setelah install, buka Docker Desktop dari Applications
```

### ⚠️ Troubleshooting Docker

**Error: "failed to connect to the docker API"**
```
docker ps
failed to connect to the docker API at unix:///var/run/docker.sock
```

**Solusi:** Docker Desktop belum jalan. Pastikan:

1. **Buka Docker Desktop** dari Applications (ikon whale di dock)
2. **Tunggu sampai status "Running"** (ada tulisan hijau di Docker Desktop)
3. **Cek lagi:** `docker --version` dan `docker ps`

**Jika masih error:**
```bash
# Restart Docker Desktop
# Atau via terminal:
killall Docker && open /Applications/Docker.app
```

---

## 🚀 Development Setup

### 1. Clone dan install dependencies

```bash
git clone <repo>
cd wangsa
make install
```

### 2. Setup PostgreSQL lokal

**Dengan PostgreSQL terinstall:**
```bash
make db-create
# Membuat database 'wangsa' dengan user 'wangsa'
```

**Dengan Docker (Direkomendasikan untuk Mac):**
```bash
# Jalankan PostgreSQL container
docker run -d \
  --name wangsa-db \
  -e POSTGRES_DB=wangsa \
  -e POSTGRES_USER=wangsa \
  -e POSTGRES_PASSWORD=wangsa \
  -p 5432:5432 \
  postgres:16-alpine

# Cek apakah container sudah jalan
docker ps

# Tunggu beberapa detik sampai database siap, lalu lanjut ke langkah 3
```

**Management Docker container:**
```bash
# Stop database
docker stop wangsa-db

# Start lagi
docker start wangsa-db

# Hapus container (jika tidak diperlukan)
docker rm -f wangsa-db

# Lihat logs jika ada masalah
docker logs wangsa-db
```

### 3. Setup environment

```bash
cp backend/.env.example backend/.env
# Edit .env jika perlu — default sudah cocok untuk setup di atas
```

### 4. Jalankan

```bash
# Terminal 1
make dev-backend

# Terminal 2
make dev-frontend
```

Buka: **http://localhost:5173**

---

## 🌐 Deployment ke Production

### Pilihan platform (semua gratis tier tersedia)

| Platform | Database | Backend | Frontend |
|---|---|---|---|
| **Railway** | Railway Postgres | Railway Service | Vercel / Netlify |
| **Render** | Render Postgres | Render Web Service | Vercel / Netlify |
| **Supabase** | Supabase Postgres | Railway / Render | Vercel / Netlify |
| **Neon** | Neon Serverless PG | Railway / Render | Vercel / Netlify |

### Environment variables yang wajib di-set di production

```env
# Database — ambil dari dashboard platform
DATABASE_URL=postgres://user:pass@host:5432/dbname?sslmode=require

# JWT — generate dengan: openssl rand -hex 32
JWT_SECRET=<min 32 karakter acak>

# Server
PRODUCTION=true
COOKIE_DOMAIN=yourdomain.com
FRONTEND_URL=https://yourfrontend.vercel.app
PORT=8080
```

### Deploy ke Railway (termudah)

```bash
# Install Railway CLI
npm install -g @railway/cli

# Login dan init
railway login
railway init

# Tambah PostgreSQL plugin dari dashboard Railway
# Railway otomatis inject DATABASE_URL ke environment

# Deploy
railway up
```

### Deploy backend ke Render

1. Buat **Web Service** di render.com
2. Connect repo GitHub
3. Set **Build Command**: `cd backend && go build -o server ./cmd/server/main.go`
4. Set **Start Command**: `./backend/server`
5. Tambah **PostgreSQL** database dari Render dashboard
6. Set environment variables di Render dashboard

### Deploy frontend ke Vercel

```bash
cd frontend
npm run build

# Install Vercel CLI
npm install -g vercel
vercel --prod
```

Atau connect repo GitHub langsung di vercel.com — auto-deploy setiap push.

**Jangan lupa:** set `VITE_API_URL` di Vercel environment jika backend bukan di domain yang sama, dan update `vite.config.js` proxy target.

---

## 🏗️ Struktur Project

```
wangsa/
├── backend/
│   ├── cmd/server/main.go              ← Entry point + graceful shutdown
│   ├── config/config.go                ← Load env vars (DATABASE_URL, dll)
│   ├── internal/
│   │   ├── delivery/http/
│   │   │   ├── router.go
│   │   │   ├── handler/                ← HTTP handlers
│   │   │   └── middleware/             ← Auth + Role JWT middleware
│   │   ├── domain/                     ← Entities + Repository interfaces
│   │   ├── usecase/                    ← Business logic
│   │   ├── repository/                 ← PostgreSQL implementations
│   │   └── pkg/
│   │       ├── database/
│   │       │   ├── postgres.go         ← Connection pool + migrations
│   │       │   └── schema.sql          ← Embedded PostgreSQL schema
│   │       ├── jwt/
│   │       └── response/
│   ├── go.mod / go.sum
│   └── .env.example
├── frontend/
│   └── src/
│       ├── stores/                     ← Pinia stores
│       ├── utils/format.js             ← Shared formatters
│       ├── composables/                ← useToast, useDirtyGuard, useZoomPan
│       └── views/ + components/
└── Makefile
```

---

## 🔧 Make Commands

```bash
make help          # Semua commands
make install       # Install Go + npm dependencies
make db-create     # Buat database PostgreSQL lokal
make db-reset      # Reset database lokal
make dev-backend   # Jalankan backend (port 8080)
make dev-frontend  # Jalankan frontend (port 5173)
make build         # Build semua untuk production
make vet           # Go vet
make fmt           # gofmt
make clean         # Hapus build artifacts
```

---

## 🛡️ Security Notes

- `JWT_SECRET` minimal 32 karakter — server **tidak akan start** jika kurang
- `PRODUCTION=true` mengaktifkan `Secure` flag pada cookie (wajib HTTPS)
- Semua delete adalah **soft delete** — data tidak hilang permanen
- Setiap perubahan data terekam di tabel `audit_log`
- `version` field di semua tabel mendukung optimistic locking
