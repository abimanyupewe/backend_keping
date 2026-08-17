# Keping - Backend API

[![Go Version](https://img.shields.io/badge/Go-1.23%2B-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://golang.org/)
[![Framework](https://img.shields.io/badge/Framework-Gin-008080?style=for-the-badge&logo=gin&logoColor=white)](https://gin-gonic.com/)
[![Database](https://img.shields.io/badge/Database-Supabase%20PostgreSQL-3ECF8E?style=for-the-badge&logo=supabase&logoColor=white)](https://supabase.com/)
[![Storage](https://img.shields.io/badge/Storage-Cloudinary-3448C5?style=for-the-badge&logo=cloudinary&logoColor=white)](https://cloudinary.com/)
[![AI Engine](https://img.shields.io/badge/AI%20Vision-Google%20Gemini%201.5%20Flash-8E44AD?style=for-the-badge&logo=google&logoColor=white)](https://ai.google.dev/)
[![Authentication](https://img.shields.io/badge/Auth-Supabase%20JWT-FF4500?style=for-the-badge&logo=jsonwebtokens&logoColor=white)](https://jwt.io/)
[![Architecture](https://img.shields.io/badge/Architecture-Clean%20Architecture-4A90E2?style=for-the-badge)](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

Backend REST API untuk **Keping** — aplikasi pelacak keuangan personal & tabungan otomatis berbasis Android (*Zero Manual Input*). Backend ini dioptimalkan untuk memproses string notifikasi *real-time* (DANA, GoPay, ShopeePay, BRImo) dan tangkapan layar struk/bukti transfer via AI Vision OCR & Cloudinary Storage, serta mengalokasikan tabungan impian (*wishlist*) secara virtual.

---

## Core Value Proposition & Fitur Utama

- **Zero Manual Input Processing:** Menerima *payload push notification* latar belakang dari aplikasi Android dan mengekstrak nominal serta tipe transaksi secara otomatis via **Regex Engine**.
- **AI Receipt Scanner (OCR):** Menggunakan **Google Gemini 1.5 Flash Vision API** untuk membedah bukti transfer/struk pembayaran dan mengembalikan data terstruktur dalam waktu $< 2$ detik.
- **Cloudinary Image Storage:** Gambar struk/bukti transfer disimpan secara otomatis ke **Cloudinary Storage** dan URL gambarnya disimpan pada kolom `image_url` tabel `transactions`.
- **Auto-Wishlist Virtual Allocation:** Setiap kali ada pemasukan (*CREDIT*), sistem secara otomatis menghitung potongan virtual $X\%$ dan menambahkannya ke *wishlist* aktif milik pengguna.
- **Database Transaction Lock:** Menggunakan `FOR UPDATE` Row-level Locking pada PostgreSQL untuk mencegah *race condition* jika beberapa notifikasi masuk secara bersamaan.
- **Stateless Supabase JWT Auth & Rate Limiter:** Dilindungi oleh Middleware verifikasi JWT Supabase dan *Token Bucket Rate Limiter*.

---

## Struktur Arsitektur (Clean Architecture)

Proyek ini menerapkan **Clean Architecture** secara ketat untuk memastikan pemisahan tanggung jawab (*separation of concerns*), keterujian (*testability*), dan kemudahan perawatan (*maintainability*):

```text
backend_keping/
├── cmd/
│   └── api/
│       └── main.go                 # Entrypoint aplikasi (Dependency Injection & Server Setup)
├── internal/
│   ├── domain/                     # Layer 1: Entitas Bisnis, Structs, & Constants
│   │   ├── user.go
│   │   ├── transaction.go
│   │   ├── wishlist.go
│   │   └── ingestion.go
│   ├── repository/                 # Layer 2: Akses Data (GORM / Supabase Queries)
│   │   ├── transaction_repository.go
│   │   └── wishlist_repository.go
│   ├── usecase/                    # Layer 3: Logika Bisnis & Core Engines
│   │   ├── regex_parser.go         # Parsing notifikasi DANA/BRImo/GoPay/ShopeePay
│   │   ├── auto_save_calculator.go # Kalkulasi matematis alokasi wishlist
│   │   ├── ingest_usecase.go       # Orchestrator alur Notifikasi & OCR Cloudinary
│   │   ├── transaction_usecase.go  # Pengelolaan data transaksi
│   │   └── wishlist_usecase.go     # Pengelolaan target wishlist
│   └── delivery/
│       └── http/                   # Layer 4: HTTP Transport Layer (Gin Framework)
│           ├── handler/            # HTTP Handlers (Requests & Responses)
│           ├── middleware/         # Supabase JWT Auth & Rate Limiter
│           └── router.go           # Routing & Grouping Endpoints
├── pkg/                            # Shared Utilities & Third-Party Wrappers
│   ├── ai/                         # Client wrapper Google Gemini 1.5 Flash API
│   ├── storage/                    # Client wrapper Cloudinary Storage API
│   ├── logger/                     # Structured Logger Utility
│   └── response/                   # JSON Standardized Response Formatter
├── postman/                        # Postman Collection JSON & API Testing
│   └── Keping_API.postman_collection.json
├── dev/                            # Dokumentasi PRD & SDD
├── .env.example                    # Template variabel lingkungan
├── go.mod                          # Dependency definitions
└── README.md                       # Dokumentasi resmi proyek
```

---

## Skema Database (Supabase PostgreSQL)

```mermaid
erdiagram
    USERS ||--o{ TRANSACTIONS : "owns"
    USERS ||--o{ WISHLISTS : "manages"

    USERS {
        uuid id PK
        string email
        string name
        timestamp created_at
    }

    TRANSACTIONS {
        uuid id PK
        uuid user_id FK
        string source "NOTIFICATION | OCR_SCAN | MANUAL"
        string type "CREDIT | DEBIT"
        numeric amount
        string description
        text raw_text
        text image_url
        timestamp timestamp
    }

    WISHLISTS {
        uuid id PK
        uuid user_id FK
        string name
        numeric target_amount
        numeric saved_amount
        numeric auto_save_percentage
        string icon_name
        string label_color
        boolean is_completed
    }
```

---

## Panduan Setup & Instalasi

### 1. Prasyarat System

- **Go:** Version `1.23+`
- **Database:** Supabase PostgreSQL Project
- **Storage:** Cloudinary Account ([Cloudinary Dashboard](https://cloudinary.com/))
- **API Key:** Google Gemini AI API Key ([Google AI Studio](https://aistudio.google.com/))

### 2. Clone & Setup Environment

```bash
# Copy file environment template
cp .env.example .env
```

Konfigurasi file `.env`:

```ini
# App & Server Configuration
APP_ENV=development
PORT=8080
GIN_MODE=debug

# Supabase Configuration
SUPABASE_URL="https://your_supabase_project_ref.supabase.co"
SUPABASE_ANON_KEY="your_supabase_anon_key_here"
SUPABASE_SERVICE_ROLE_KEY="your_supabase_service_role_key_here"
SUPABASE_JWT_SECRET="your_supabase_jwt_secret_here"
DB_URL="postgresql://postgres.[PROJECT_REF]:[PASSWORD]@aws-0-[REGION].pooler.supabase.com:6543/postgres?sslmode=require"

# File Storage (Cloudinary)
CLOUDINARY_CLOUD_NAME="your_cloudinary_cloud_name_here"
CLOUDINARY_API_KEY="your_cloudinary_api_key_here"
CLOUDINARY_API_SECRET="your_cloudinary_api_secret_here"
CLOUDINARY_FOLDER="receipts"

# AI Service (Google Gemini Vision)
GEMINI_API_KEY="your_gemini_api_key_here"
GEMINI_MODEL="your_gemini_model_here"

# Security & Ingestion
JWT_SECRET="your_jwt_secret_64_character_hex_here"
JWT_EXPIRATION_HOURS="your_jwt_expiration_hours_here"
INGESTION_SECRET_KEY="your_ingestion_secret_key_here"

# Rate Limit & Upload
RATE_LIMIT_RPS=5
RATE_LIMIT_BURST=10
MAX_UPLOAD_SIZE_MB=5
```

### 3. Install Dependensi & Jalankan Aplikasi

```bash
# Unduh dan rapihkan dependensi Go
go mod tidy

# Menjalankan server lokal
go run cmd/api/main.go
```

---

## Spesifikasi API Endpoints

Format respon seluruh API mengikuti standar:
```json
{
  "success": true,
  "message": "Pesan deskriptif",
  "data": { ... }
}
```

### 1. Infrastructure & Health Check
* `GET /api/v1/health`
  * **Headers:** None
  * **Response (200 OK):**
    ```json
    {
      "success": true,
      "message": "Keping Backend API & Database is Healthy",
      "data": { "database": "CONNECTED", "status": "UP" }
    }
    ```

### 2. Ingestion Module (Android & OCR)
* `POST /api/v1/ingest/notification`
  * **Headers:** 
    * `Authorization: Bearer <JWT_SECRET_TOKEN>`
    * `X-Ingestion-Secret: <INGESTION_SECRET_KEY>`
  * **Body:**
    ```json
    {
      "package_name": "id.dana",
      "title": "Isi Saldo Berhasil",
      "content": "Kamu telah berhasil melakukan isi saldo sebesar Rp500.000 via BRI"
    }
    ```
* `POST /api/v1/ingest/scan`
  * **Headers:** `Authorization: Bearer <JWT_SECRET_TOKEN>`, `X-Ingestion-Secret: <INGESTION_SECRET_KEY>`
  * **Body:** `multipart/form-data` (field: `receipt` berupa image file)

### 3. Transaksi & Wishlists
* `GET /api/v1/transactions?limit=20`
  * **Headers:** `Authorization: Bearer <JWT_SECRET_TOKEN>`
* `POST /api/v1/wishlists`
  * **Headers:** `Authorization: Bearer <JWT_SECRET_TOKEN>`
  * **Body:**
    ```json
    {
      "name": "Tiket Konser",
      "target_amount": 1500000,
      "auto_save_percentage": 10,
      "icon_name": "ticket",
      "label_color": "#8E44AD"
    }
    ```
* `GET /api/v1/wishlists`
  * **Headers:** `Authorization: Bearer <JWT_SECRET_TOKEN>`

---

## Testing & Quality Assurance

### 1. Unit Testing
Proyek ini dilengkapi dengan Unit Test untuk layer kritis (*Regex Engine* dan *Auto-Save Calculator*):

```bash
# Menjalankan seluruh unit test
go test -v ./...

# Mengecek Test Coverage
go test -cover ./internal/usecase/...
```

### 2. Postman Collection
File [postman/Keping_API.postman_collection.json](file:///home/clara/Project/backend_keping/postman/Keping_API.postman_collection.json) dapat langsung di-import ke **Postman**. Collection ini sudah dikelompokkan berdasarkan folder (`Health Check`, `Ingestion Engine`, `Transactions`, `Wishlists`) dan dilengkapi dengan variabel lingkungan (`{{base_url}}`, `{{bearer_token}}`, `{{ingestion_secret}}`).

---

## Lisensi & Kredit

Dipersembahkan untuk proyek **Keping Financial Tracker** — MVP Release 2026.
