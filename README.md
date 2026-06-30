# golang-rbac

Boilerplate REST API Golang dengan Role-Based Access Control (RBAC) sederhana,
JWT authentication, dan arsitektur clean (domain - usecase - repository - delivery).

## Struktur Folder

```
golang-rbac/
├── cmd/
│   └── main.go                      # Entry point aplikasi
├── internal/
│   ├── config/
│   │   └── config.go                # Load .env & konfigurasi aplikasi (DSN, JWT secret, port)
│   ├── delivery/
│   │   └── http/                    # Handler HTTP (Router, Request/Response, Middleware)
│   │       ├── middleware/
│   │       │   └── auth_middleware.go
│   │       └── user_handler.go
│   ├── domain/                      # Blueprint: Entity, Interface, & Model (Core Business)
│   │   └── user.go
│   ├── repository/                  # Akses ke Database (SQL, NoSQL, Cache)
│   │   └── user_repository.go
│   └── usecase/
│       └── user_usecase.go          # Logika Bisnis (Business Logic / Service)
├── migrations/
│   ├── 000001_init_schema.up.sql    # File migrasi (create table & seed data)
│   └── 000001_init_schema.down.sql  # Rollback (drop table)
├── pkg/
│   └── utils/
│       └── json_response.go         # Helper response JSON standar
├── go.mod
├── go.sum
└── .env                              # DB_DSN, JWT_SECRET, PORT
```

## Menjalankan

1. Buat database MySQL lalu jalankan migration di `migrations/000001_init_schema.up.sql`.
2. Copy `.env.example` jadi `.env` lalu sesuaikan nilainya:
   ```
   cp .env.example .env
   ```
   ```
   DB_DSN=root:password@tcp(127.0.0.1:3306)/golang_rbac_db
   JWT_SECRET=ganti_dengan_secret_yang_kuat_dan_acak
   PORT=8080
   ```
   `.env` sudah masuk `.gitignore` jadi aman dari ke-push ke GitHub.
3. Install dependency & jalankan:
   ```
   go mod tidy
   go run ./cmd
   ```

## Endpoint

| Method | Endpoint        | Akses                         | Keterangan                              |
| ------ | --------------- | ----------------------------- | --------------------------------------- |
| POST   | `/register`     | Public                        | Daftar user baru (role default: Viewer) |
| POST   | `/login`        | Public                        | Login, mengembalikan JWT token          |
| GET    | `/me`           | Admin, Editor, Viewer (login) | Profil user yang sedang login           |
| GET    | `/admin`        | Admin                         | Contoh halaman khusus Admin             |
| GET    | `/edit-article` | Admin, Editor                 | Contoh halaman Admin & Editor           |
| GET    | `/view`         | Admin, Editor, Viewer         | Contoh halaman semua role               |

Untuk endpoint terproteksi, kirim header:

```
Authorization: Bearer <token>
```

## Role Default (seed data)

| ID  | Role   |
| --- | ------ |
| 1   | Admin  |
| 2   | Editor |
| 3   | Viewer |
