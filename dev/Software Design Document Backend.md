# **Software Design Document (SDD) \- Backend Keping**

**Proyek:** Keping (Financial Tracker)  
**Platform:** REST API Backend  
**Stack Teknologi:** Go (Golang), Gin Framework, PostgreSQL (Supabase), Gemini 1.5 Flash Vision API.  
**Arsitektur:** Clean Architecture

## **1\. Pendahuluan**

Backend Keping berfungsi sebagai pusat pemrosesan data (otak) dari aplikasi. Berbeda dengan sistem pelacak keuangan pada umumnya, backend ini dioptimalkan untuk memproses string teks mentah (dari notifikasi Android) dan memproses gambar (OCR struk via AI), lalu menerjemahkannya menjadi entitas finansial dan mengalokasikannya ke target tabungan (*wishlist*) secara otomatis.

## **2\. Struktur Arsitektur (Clean Architecture)**

Proyek akan dibagi menjadi empat lapisan utama untuk memisahkan *concern* dan memudahkan pengujian (*testing*):  
keping-backend/  
├── cmd/  
│   └── api/  
│       └── main.go                 \# Entry point aplikasi  
├── internal/  
│   ├── domain/                     \# Layer 1: Entitas Bisnis & Interface (Structs)  
│   ├── delivery/  
│   │   ├── http/  
│   │   │   ├── handler/            \# Layer 4: Gin Handlers (Menerima Request)  
│   │   │   ├── middleware/         \# Autentikasi JWT, Rate Limiting  
│   │   │   └── router.go           \# Definisi endpoint API  
│   ├── usecase/                    \# Layer 3: Logika Bisnis (Regex, Kalkulasi Wishlist, Integrasi AI)  
│   └── repository/                 \# Layer 2: Akses Data (GORM/PostgreSQL Query ke Supabase)  
├── pkg/  
│   ├── ai/                         \# Wrapper client untuk Google Gemini API  
│   ├── logger/                     \# Utilitas logging  
│   └── response/                   \# Standarisasi format JSON response  
├── .env  
└── go.mod

## **3\. Skema Database (PostgreSQL \- Supabase)**

Terdapat 3 tabel utama yang saling berelasi.

### **A. Tabel users**

Tabel ini biasanya di-handle oleh Supabase Auth, namun kita menyimpan ekstensi profil pengguna di schema public.

> * id (UUID, Primary Key)  
> * email (String)  
> * name (String)  
> * created\_at (Timestamp)

### **B. Tabel transactions**

Menyimpan seluruh riwayat pemasukan dan pengeluaran hasil tangkapan notifikasi atau OCR.

> * id (UUID, Primary Key)  
> * user\_id (UUID, Foreign Key)  
> * source (String) \- Enum: NOTIFICATION, OCR\_SCAN, MANUAL  
> * type (String) \- Enum: CREDIT (Masuk), DEBIT (Keluar)  
> * amount (Numeric/Decimal)  
> * description (String)  
> * raw\_text (Text) \- Teks asli notifikasi (untuk keperluan *debugging* jika *regex* gagal)  
> * timestamp (Timestamp)

### **C. Tabel wishlists**

Menyimpan target tabungan pengguna.

> * id (UUID, Primary Key)  
> * user\_id (UUID, Foreign Key)  
> * name (String) \- Misal: "Tiket Konser"  
> * target\_amount (Numeric)  
> * saved\_amount (Numeric) \- Default 0  
> * auto\_save\_percentage (Numeric) \- Misal: 10 (artinya 10%)  
> * icon\_name (String) \- Nama referensi ikon vektor  
> * label\_color (String) \- Kode warna Hex (misal: \#2A9D8F)  
> * is\_completed (Boolean) \- Default False

## **4\. Alur Logika Sistem (Core Workflows)**

### **4.1 Modul Ingesti Notifikasi (Regex Engine)**

Saat aplikasi Flutter mengirim teks notifikasi ke backend:

> 1. **Delivery Layer:** Menerima *payload* dan memvalidasi INGESTION\_SECRET\_KEY.  
> 2. **Usecase Layer:** Menggunakan Regex (Regular Expression) untuk mencari nominal.  
   * Jika paket asal \= id.dana dan teks mengandung "berhasil melakukan isi saldo", status \= CREDIT.  
> 3. **Repository Layer:** Menyimpan hasil ekstraksi ke tabel transactions.  
> 4. **Wishlist Trigger:** Memicu kalkulasi alokasi otomatis (lihat 4.3).

### **4.2 Modul OCR & Struk (Gemini API)**

Saat pengguna mengunggah gambar struk:

> 1. **Delivery Layer:** Menerima file gambar (multipart/form-data).  
> 2. **Usecase (AI Wrapper):** Mengirim gambar ke Gemini Flash Vision dengan *System Prompt* ketat untuk memaksa output JSON baku:  
>    { "amount": 50000, "type": "DEBIT", "merchant": "Kopi Kenangan" }

> 3. **Repository Layer:** Menyimpan ke tabel transactions dengan source OCR\_SCAN.

### **4.3 Kalkulator Alokasi Otomatis (Auto-Save Engine)**

Setiap kali ada transaksi bertipe CREDIT (uang masuk), sistem melakukan kalkulasi berikut di Usecase Layer:  
Untuk setiap *wishlist* aktif milik *user* tersebut, jumlah uang yang ditabung dihitung menggunakan formula:  
$$\\text{Alokasi} \= \\text{Amount} \\times \\left( \\frac{\\text{Auto-Save Percentage}}{100} \\right)$$  
Sistem kemudian memperbarui kolom saved\_amount di *database* secara transaksional agar tidak terjadi *race condition*.

## **5\. Spesifikasi API Endpoint**

Semua *response* mengikuti format standar: { "success": bool, "message": string, "data": object }.

### **5.1 Endpoint Ingesti (Otomatisasi)**

> * POST /api/v1/ingest/notification  
  * **Fungsi:** Menerima teks notifikasi latar belakang dari Android.  
  * **Body:** { "package\_name": "id.dana", "title": "...", "content": "..." }  
> * POST /api/v1/ingest/scan  
  * **Fungsi:** Mengunggah dan membaca gambar struk via AI.  
  * **Body:** Form-data file (image/jpeg).

### **5.2 Endpoint Transaksi & Wishlist (Transaksional)**

> * GET /api/v1/transactions  
  * **Fungsi:** Mengambil riwayat transaksi.  
> * POST /api/v1/wishlists  
  * **Fungsi:** Membuat target tabungan baru.  
  * **Body:** { "name": "HP Baru", "target\_amount": 5000000, "auto\_save\_percentage": 10, "icon\_name": "smartphone", "label\_color": "\#4A90E2" }  
> * GET /api/v1/wishlists  
  * **Fungsi:** Mengambil daftar progres *wishlist*.

### **5.3 Endpoint Infrastruktur**

> * GET /api/v1/health  
  * **Fungsi:** *Keep-alive ping* untuk mencegah server Go tidur dan mencegah Supabase melakukan *auto-pause*. Melakukan instruksi SELECT 1 ke database.

## **6\. Keamanan & Performa**

> * **Autentikasi:** Endpoint (selain /health) dilindungi oleh *middleware* JWT Auth yang menvalidasi token dari Flutter.  
> * **Rate Limiting:** Menggunakan algoritma *Token Bucket* dengan batas RATE\_LIMIT\_RPS=5 untuk mencegah eksploitasi API.  
> * **Transaction Lock:** Pembaruan saldo *wishlist* wajib dibungkus dalam *database transaction* (db.Transaction di GORM) agar data keuangan tidak tumpang tindih jika ada dua notifikasi masuk secara bersamaan.

## **7\. Strategi Pengujian (Testing Strategy & Unit Test)**

Karena kita menggunakan *Clean Architecture*, setiap layer dapat diuji secara independen (terisolasi) tanpa harus terhubung ke database asli.

### **7.1 Framework & Perangkat (Tools)**

> * **Testing Engine:** Bawaan standar Go (testing package).  
> * **Assertion & Mocking:** github.com/stretchr/testify (menyediakan fungsi assert dan pembuatan mock objek).  
> * **Mock Generator:** mockery untuk menghasilkan implementasi palsu dari antarmuka (interface) Layer Repository.

### **7.2 Spesifikasi Pengujian per Layer**

> * **Usecase Layer (Prioritas Utama):**  
  * **Regex Parser Test:** Menguji berbagai variasi teks notifikasi dari DANA/BRImo (berhasil, gagal, typo) dan memastikan *output* nominal dan tipe transaksinya akurat 100%.  
  * **Auto-Save Engine Test:** Menguji keakuratan perhitungan matematis alokasi tabungan (memastikan tidak ada kesalahan pembulatan desimal).  
  * Pada layer ini, koneksi ke *database* dipalsukan (*mocked*).  
> * **Repository Layer:**  
  * Menggunakan DATA-DOG/go-sqlmock untuk memastikan sintaks SQL GORM (seperti *Transaction Lock*) tergenerate dengan benar tanpa perlu menembak Supabase sungguhan.  
> * **Delivery Layer (HTTP Handlers):**  
  * Menggunakan package net/http/httptest bawaan Go untuk mensimulasikan *request* HTTP (POST/GET) dan memvalidasi format respon JSON dan kode status HTTP (misal: 200 OK, 400 Bad Request).

### **7.3 Target *Code Coverage***

> * Modul **Regex Parser** dan **Kalkulator Wishlist** diwajibkan memiliki *test coverage* minimum 90%, karena merupakan inti penggerak otomasi aplikasi ini.