# **Product Requirements Document (PRD)**

**Nama Produk:** Keping  
**Platform:** Mobile App (Android \- Flutter) & Backend API (Go/Gin)  
**Tahap Pengembangan:** Minimum Viable Product (MVP) \- Lingkup Indie/Pertemanan  
**Tagline:** *"Lacak otomatis. Tabung tanpa sadar."*

## **1\. Ringkasan Eksekutif (Executive Summary)**

**Keping** adalah aplikasi pelacak keuangan personal yang dirancang khusus untuk orang-orang yang malas mencatat pengeluaran secara manual. Berbeda dengan aplikasi keuangan pada umumnya yang menuntut kedisiplinan *input* data, Keping bekerja secara pasif di latar belakang (*background*). Aplikasi ini membaca notifikasi transaksi (BRImo, DANA, GoPay, ShopeePay) dan membedah tangkapan layar (struk/bukti transfer) menggunakan AI, lalu mengalokasikan sisa uang ke daftar impian (*wishlist*) secara otomatis.

## **2\. Masalah & Solusi**

### **2.1 Masalah (Pain Points)**

> * **Kelelahan Mencatat (*Tracking Fatigue*):** Pengguna sering lupa mencatat pengeluaran kecil (seperti jajan kopi atau parkir), sehingga di akhir bulan saldo tidak sinkron dengan catatan.  
> * **Integrasi Bank yang Rumit & Berisiko:** Menghubungkan akun bank langsung via bot/scraping berisiko membuat akun pengguna diblokir oleh pihak bank.  
> * **Menabung Menjadi Beban:** Mengisihkan uang untuk *wishlist* terasa berat jika dilakukan secara manual dan dalam nominal besar sekaligus.

### **2.2 Solusi (Pendekatan "Keping")**

> * **Zero Manual Input:** Menangkap nominal uang keluar/masuk secara seketika melalui *Android Notification Listener*.  
> * **AI Fallback:** Untuk transfer via ATM atau laptop, pengguna cukup *upload screenshot* bukti transfer, dan AI (Gemini Vision) akan mengekstrak datanya dalam 2 detik.  
> * **Auto-Wishlist Allocation:** Setiap ada uang masuk, sistem otomatis memotong X% secara virtual untuk mengisi progres *wishlist* tanpa memindahkan uang fisik pengguna.

## **3\. Ruang Lingkup (Project Scope)**

Untuk menjaga fokus pada fase MVP, berikut adalah batasan ruang lingkup pengembangan:

### **3.1 In-Scope (Termasuk dalam MVP)**

> * **Platform Mobile:** Hanya Android (karena membutuhkan izin *Notification Listener*).  
> * **Target E-Wallet/Bank:** DANA, GoPay, ShopeePay, dan BRImo.  
> * **Metode Input:**  
  1. Otomatis via penangkapan Push Notification (Lokal Android).  
  2. Semi-otomatis via *upload screenshot* / foto struk OCR (Gemini Vision API).  
> * **Fitur Utama:** Manajemen *Wishlist* dengan kustomisasi visual (Ikon Vektor & Warna Kustom), Auto-allocation persentase uang, Histori Transaksi, dan AI Advisor Banner.  
> * **Desain:** Hanya *Light Mode* (Mode Terang).  
> * **Infrastruktur:** Backend API (Go/Gin) dan Database (Supabase PostgreSQL).

### **3.2 Out-of-Scope (Ditunda untuk Versi Masa Depan)**

> * **Grafik & Diagram Analitik Detail:** Visualisasi data seperti *pie chart* pengeluaran bulanan atau grafik batang ditiadakan sementara agar fokus pada *core value* tabungan otomatis.  
> * **Platform iOS (iPhone):** Dikeluarkan dari MVP karena sistem operasi iOS memblokir akses aplikasi pihak ketiga untuk membaca notifikasi sistem.  
> * **Integrasi API Bank Resmi (Open Banking):** Tidak menggunakan API aggregator (seperti Brick/Ayoconnect) untuk menghindari biaya per-transaksi dan kewajiban legalitas PT/CV.  
> * **Multi-User / Tabungan Bersama:** MVP hanya fokus pada *single-player mode* (lacak keuangan pribadi).  
> * **Dark Mode:** Ditiadakan sementara untuk mempercepat proses *slicing* UI di Flutter.  
> * **Ekspor Laporan (PDF/CSV):** Belum menjadi prioritas.

## **4\. Target Pengguna (User Persona)**

> * **Demografi:** Dewasa muda (18 \- 30 tahun), pekerja kantoran awal, atau mahasiswa.  
> * **Perilaku Finansial:** Sering bertransaksi menggunakan *e-wallet* atau QRIS (DANA, ShopeePay, GoPay, BCA/BRI Mobile). Sangat bergantung pada ponsel.  
> * **Tujuan Utama:** Ingin memantau ke mana uang pergi dan ingin menabung untuk barang impian (seperti HP baru, tiket konser, atau liburan) tanpa perlu pusing melihat angka yang rumit.

## **5\. Spesifikasi Fitur Inti (Core Features)**

### **5.1. Modul Ingesti Otomatis (The "Listener")**

> * **Fungsi:** Berjalan sebagai *background service* di Android.  
> * **Cara Kerja:** Menangkap *push notification* dari paket aplikasi tertentu (misal: id.dana, id.co.bri.brimo).  
> * **Output:** Mengirim teks mentah notifikasi ke backend Go secara *real-time* via *endpoint* aman (dilengkapi INGESTION\_SECRET\_KEY).

### **5.2. Modul AI Receipt & Screenshot Scanner**

> * **Fungsi:** Mengolah gambar bukti transaksi yang diunggah pengguna.  
> * **Cara Kerja:** Backend menerima gambar, lalu memanggil Google Gemini 1.5 Flash Vision API dengan instruksi ketat untuk mengembalikan JSON (Nominal, Tipe Transaksi, Nama Merchant/Tujuan).

### **5.3. Modul Kategori & Regex Parsing (Backend)**

> * **Fungsi:** Mengubah teks notifikasi menjadi data terstruktur.  
> * **Cara Kerja:** Backend Go menggunakan pola Regex untuk mendeteksi kata kunci (misal: "berhasil", "pembayaran QRIS", "Rp50.000") dan mengkategorikan transaksi (Makan, Transportasi, Hiburan) secara heuristik.

### **5.4. Modul Manajemen Wishlist & Alokasi Otomatis**

> * **Kustomisasi Kategori & Visual:** Saat membuat target *wishlist* baru, pengguna dapat memilih Kategori, mengatur Ikon Vektor (dari pustaka ikon aplikasi), dan memilih Warna Label (*hex color*) untuk menyesuaikan tampilan kartu *wishlist* mereka.  
> * **Kalkulator Alokasi Otomatis:** Sistem yang menghitung potongan secara virtual.  
> * **Skenario:** Pengguna menargetkan "Tiket Konser Rp1.500.000" (Ikon: Tiket, Warna: Ungu). Pengguna mengatur potongan otomatis 10%. Ketika ada notifikasi uang masuk Rp500.000 dari gaji/transfer, Keping otomatis mencatat Rp50.000 ke dalam progres tiket konser tersebut.

### **5.5. Modul AI Financial Advisor (Smart Insight)**

> * **Fungsi:** Memberikan teguran halus, pujian, atau peringatan di layar beranda berdasarkan pola pengeluaran harian/mingguan.  
> * **Cara Kerja:** Backend (Cron Job) menjalankan analisis seminggu sekali atau saat pengeluaran tiba-tiba melonjak. Data ringkasan dikirim ke Gemini dengan *prompt* khusus. Hasil teks disimpan di *database* agar tidak membebani pemanggilan API secara *real-time*.

## **6\. Arsitektur & Tech Stack**

> * **Pendekatan Desain Sistem:** Mengadopsi **Clean Architecture** secara ketat pada kedua sisi (Frontend dan Backend) untuk memastikan kode modular, mudah diuji (*testable*), dan mudah dikembangkan.  
> * **Frontend (Mobile App):** Flutter (Dart)  
  * *Arsitektur:* Clean Architecture (Domain, Data, Presentation Layers).  
  * *Modul Notifikasi:* Menggunakan *package* flutter\_local\_notifications dan *Notification Listener Service* khusus Android.  
  * *State Management:* Riverpod.  
> * **Backend (REST API):** Go (Golang) \+ Gin Framework  
  * *Arsitektur:* Clean Architecture (Handler/Delivery, Use Case, Repository Layers).  
  * Ringan, memori sangat kecil, dan cocok untuk menangani *regex parsing* berkecepatan tinggi.  
> * **Database & Auth:** Supabase (PostgreSQL terkelola)  
  * Autentikasi menggunakan Supabase Auth (Email/Password atau Google Sign-In).  
> * **Kecerdasan Buatan (AI OCR & Advisor):** Google Gemini 1.5 Flash API.

## **7\. Filosofi Desain UI/UX**

> * **Pendekatan Visual:** *Clean, Minimalist, Data-Driven*. Mengedepankan ruang putih (*whitespace*) agar aplikasi tidak terasa sesak. Penggunaan Ikon Vektor dan warna label pada *wishlist* dirancang untuk memberikan sentuhan personal tanpa merusak estetika tata letak.  
> * **Warna:** Eksklusif **Mode Terang (Light Mode)**.  
  * *Background:* Putih bersih (\#FFFFFF) atau abu-abu sangat muda (\#F8F9FA).  
  * *Teks:* Abu-abu arang / hitam lembut (\#1F2937) untuk keterbacaan optimal.  
  * *Aksen:* Hijau Mint / Emerald untuk pemasukan, dan Merah Karang (*Coral*) untuk pengeluaran.  
> * **Beban Kognitif (Cognitive Load):** Karena aplikasi ini bekerja otomatis, layar utama (Dashboard) tidak boleh penuh sesak. Hanya tampilkan 4 elemen utama:  
  1. **AI Advisor Banner:** Kartu kecil di atas yang berisi teks interaktif peringatan/saran AI.  
  2. **Total Saldo Berjalan:** Hasil kalkulasi estimasi.  
  3. **Progres Wishlist Utama:** Bentuk *circular bar* atau *progress bar* tebal, menampilkan Ikon Vektor dan Warna yang dipilih pengguna.  
  4. **Feed Aktivitas:** Riwayat transaksi terbaru (seperti *timeline* media sosial).

## **8\. Kriteria Kesuksesan (Success Metrics)**

> * **Akurasi Parsing:** Minimal 95% notifikasi dari DANA dan BRImo berhasil dibaca nominalnya dengan benar oleh Regex Engine.  
> * **Kecepatan OCR:** AI dapat mengenali dan memasukkan data dari *screenshot* bukti transfer ke dalam *database* dalam waktu kurang dari 4 detik.  
> * **Frictionless Experience:** Pengguna tidak perlu membuka aplikasi untuk mencatat pengeluaran; mereka hanya membuka aplikasi untuk melihat laporan dan progres *wishlist*.