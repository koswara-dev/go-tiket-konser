# Go Ticket Concert API

Go Ticket Concert API adalah backend service berbasis RESTful API yang dirancang untuk mengelola pemesanan tiket konser musik secara real-time dan aman. Proyek ini mengimplementasikan arsitektur bersih (*Clean Architecture*) serta perlindungan terhadap beban tinggi (*high concurrency / ticket war*).

---

## 🚀 Tech Stack

Proyek ini dibangun menggunakan teknologi modern untuk performa tinggi, keandalan, dan skalabilitas:

- **Bahasa Pemrograman**: [Go (Golang)](https://golang.org/) v1.25.0
- **Web Framework**: [Gin Gonic](https://github.com/gin-gonic/gin) - Framework HTTP minimalis dan berkinerja tinggi.
- **Relational Database & ORM**: [PostgreSQL](https://www.postgresql.org/) dengan [GORM](https://gorm.io/) - Untuk data transaksi terstruktur, relasional, dan ACID compliance.
- **NoSQL Database**: [MongoDB](https://www.mongodb.org/) - Digunakan untuk menyimpan Audit Logs (logging request/response), riwayat notifikasi, dan riwayat chat WebSocket.
- **In-Memory Cache & Session**: [Redis](https://redis.io/) - Digunakan untuk token blacklist (saat logout) dan pencatatan API Rate Limiting.
- **Autentikasi**: JWT (JSON Web Tokens) menggunakan [golang-jwt](https://github.com/golang-jwt/jwt) untuk sesi user yang aman.
- **Real-Time Communication**:
  - **WebSockets** (via [gorilla/websocket](https://github.com/gorilla/websocket)) - Untuk fitur live support chat antara admin dan customer.
  - **Server-Sent Events (SSE)** - Untuk *streaming* real-time notifikasi ke customer.
- **Log & Monitoring**: [Logrus](https://github.com/sirupsen/logrus) untuk structured logging dan [Lumberjack](https://github.com/natefinch/lumberjack) untuk log rotation.
- **API Documentation**: [Swagger / Swaggo](https://github.com/swaggo/swag) - Dokumentasi API interaktif.
- **Validasi**: [go-playground/validator](https://github.com/go-playground/validator) - Untuk validasi data request.

---

## ✨ Fitur Utama

1. **Authentication & Authorization (RBAC)**:
   - Registrasi, Login, dan Verifikasi OTP via Email (SMTP).
   - Role-Based Access Control (Admin & Customer).
   - Sesi aman menggunakan JWT Token dengan mekanisme logout instan via blacklist di Redis.
2. **Manajemen Konser (CRUD)**:
   - Pengelolaan konser oleh Admin.
   - Fitur unggah/upload thumbnail (image) dan berkas panduan/rules (PDF) ke penyimpanan lokal.
3. **Manajemen Kategori Tiket**:
   - Pembuatan kategori tiket (seperti Gold, Silver, dll.) dengan harga dan kuota yang terintegrasi ke konser tertentu.
4. **Pemesanan Tiket (Booking / Transaction)**:
   - Alur booking terproteksi dengan mitigasi *race condition* / *overselling* saat *ticket war* menggunakan skema SQL **Pessimistic Locking (`SELECT ... FOR UPDATE`)** dalam transaksi database database.
5. **Real-time Live Chat Support**:
   - Menggunakan WebSocket untuk menghubungkan customer dengan admin secara instan.
   - Pesan chat dan room disimpan secara persisten di MongoDB.
6. **Real-time SSE Notifications**:
   - Pengiriman notifikasi pemesanan dan update konser langsung ke user secara real-time.
7. **Security & Reliability**:
   - **CORS Protection**: Konfigurasi CORS dinamis yang mengizinkan semua domain (`*`) pada lingkungan *development* dan membatasi hanya untuk `https://domain.id` pada lingkungan *production*.
   - **API Rate Limiter**: Membatasi jumlah request per IP menggunakan Redis Pipeline guna mencegah serangan brute-force atau spamming.
   - **API Key Protection**: Pengamanan endpoint menggunakan token header `x-api-key`.
   - **Audit Logs**: Setiap request dan response krusial dicatat dan disimpan di database MongoDB untuk keperluan forensik/keamanan.

---

## 📁 Struktur Proyek (Clean Architecture)

```text
go-tiket-konser/
├── config/             # Konfigurasi koneksi DB (Postgres, Mongo, Redis)
├── docs/               # Dokumentasi Swagger yang di-generate otomatis
├── dto/                # Data Transfer Objects (Payload input & output API)
├── handler/            # Controller / HTTP Handlers (Gin)
├── logs/               # Lokasi file log aplikasi (lumberjack)
├── middleware/         # Middleware Gin (JWT, ApiKey, Rate Limiter, Audit Log)
├── models/             # GORM Database Models & Structs
├── repository/         # Data Access Object / Database Queries
├── routes/             # Defini rute API v1 & WebSocket/SSE
├── service/            # Business Logic Layer (Layanan inti aplikasi)
├── uploads/            # Direktori file media statis (thumbnail, PDF)
├── utils/              # Helper utilitas seperti Logger
├── main.go             # Entrypoint aplikasi utama
└── .env                # File konfigurasi environment variabel
```

---

## 🛠️ Persyaratan Sistem (Prerequisites)

Sebelum menjalankan aplikasi, pastikan Anda telah memasang:
- **Go** (versi 1.25.0 atau lebih tinggi)
- **PostgreSQL** (port default `5434` / sesuaikan di `.env`)
- **Redis** (port default `6379`)
- **MongoDB** (port default `27017`)
- **Python** atau **PowerShell** (opsional, untuk menjalankan skrip pengujian)

---

## ⚙️ Konfigurasi Environment (`.env`)

Salin file `.env` di direktori utama proyek dan sesuaikan nilainya:

```ini
APP_PORT=8080
DB_HOST=localhost
DB_PORT=5434
DB_USER=postgres
DB_PASS=secret45
DB_NAME=eticketdb

SMTP_HOST=sandbox.smtp.mailtrap.io
SMTP_PORT=587
SMTP_USER=your_smtp_user
SMTP_PASSWORD=your_smtp_password
EMAIL_FROM=noreply@tiket-konser.id

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

MONGO_URI=mongodb://localhost:27017
MONGO_DB=eticketdb

APP_ENV=development
```

---

## 🏃 Cara Menjalankan Aplikasi

### 1. Jalankan Database & Services
Pastikan PostgreSQL, Redis, dan MongoDB sudah menyala di sistem Anda.

### 2. Jalankan Server Aplikasi

Anda dapat langsung menjalankan file `main.go`:
```bash
go run main.go
```

Atau menggunakan [Air](https://github.com/cosmtrek/air) untuk *live reload* / hot reloading saat terjadi perubahan kode:
```bash
air
```

Saat pertama kali dijalankan, sistem secara otomatis akan melakukan:
- **Auto Migration**: Membuat tabel database PostgreSQL berdasarkan model GORM.
- **Seeding Data**: Menambahkan 5 data konser awal, 2 kategori tiket, 2 data customer uji coba, dan 1 data admin default.

---

## 📖 Dokumentasi API (Swagger)

Saat server berjalan di mode `development`, dokumentasi API Swagger interaktif dapat diakses pada tautan berikut:

👉 **[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)**

---

## 🧪 Pengujian (Testing)

### 1. Integrasi & Fitur Utama (RBAC, IDOR, Booking)
Kami menyediakan skrip otomatis untuk memvalidasi fitur-fitur penting dalam backend ini. Jalankan melalui PowerShell:
```powershell
./run_tests.ps1
```
Hasil pengujian akan tercatat pada berkas `result-testing.md`.

### 2. Pengujian Konkurensi & Beban Tinggi (Ticket War Simulation)
Untuk menguji keandalan penanganan pesanan secara massal dan konkuren:
```bash
python run_concurrency_test.py
```
Skrip ini menyimulasikan 200 user melakukan pesanan tiket secara bersamaan dalam kurun waktu 5 detik pada sisa kuota tiket 50. Output uji coba dapat dipantau di konsol dan laporan pengujian tercatat di `result-concurrency.md`.
