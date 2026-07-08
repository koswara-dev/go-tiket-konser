# Plan Concurrency & High Load Handling

Dokumen ini berisi rencana optimasi dan pengujian sistem untuk menangani beban tinggi (*high concurrency*) dari banyak user secara bersamaan, khususnya saat war tiket (*ticket war*).

---

## 1. Strategi Penanganan Konkurensi Tinggi

### A. Proteksi Race Condition & Over-Selling (Pessimistic Locking)
Saat ribuan user mencoba melakukan booking kategori tiket yang sama dalam waktu mikrodetik yang berdekatan, terdapat risiko *race condition* di mana tiket terjual melebihi sisa kuota yang tersedia (*over-selling*).
* **Solusi Saat Ini:** Menggunakan kueri SQL `SELECT ... FOR UPDATE` via GORM transaction block pada `service/booking.go`.
* **Detail Kerja:** Baris kategori tiket yang sedang dipesan akan dikunci (*row lock*) dari proses modifikasi transaksi lain hingga transaksi saat ini di-commit atau di-rollback.

### B. Optimalisasi Database Connection Pool
Membatasi dan mengelola koneksi ke database Postgres agar server Go tidak kehabisan *file descriptor* atau mengalami *connection timeout*.
* **Rencana Pengaturan (di `config/db.go`):**
  * `db.SetMaxOpenConns(100)`: Jumlah maksimal koneksi aktif ke database.
  * `db.SetMaxIdleConns(20)`: Jumlah koneksi menganggur yang tetap dipertahankan.
  * `db.SetConnMaxLifetime(1 * time.Hour)`: Durasi maksimal suatu koneksi dapat digunakan kembali.

### C. Pencegahan Goroutine Leak (SSE & WebSocket)
Klien yang terputus secara tidak terduga pada koneksi panjang (*persistent connection*) dapat menyebabkan goroutine tetap hidup (*leak*).
* **Solusi:**
  * **SSE:** Memanfaatkan handler `c.Stream` di Gin yang mendeteksi `c.Writer.CloseNotify()` atau konteks request `c.Request.Context().Done()` untuk memicu pembatalan pendaftaran channel (`Unregister`) secara bersih.
  * **WebSocket:** Menyetel `ReadDeadline` dan `WriteDeadline` serta mengirimkan `PingMessage` secara periodik agar koneksi yang mati di sisi klien dapat dideteksi dan dilepaskan (*cleanup*).

### D. Optimasi Rate Limiter (Pencegahan Memory Leak)
Rate Limiter berbasis memori (`IPRateLimiter`) yang menyimpan IP address di dalam map tanpa batas waktu akan memicu *memory leak* seiring bertambahnya pengunjung unik.
* **Rencana Solusi:** Mengintegrasikan rate limiter berbasis Redis menggunakan skema Token Bucket atau Sliding Window Log, atau menyertakan mekanisme pembersihan periodik (misal: menggunakan `sync.Map` dengan TTL).

---

## 2. Rencana Pengujian Konkurensi (Concurrency Load Test)

Pengujian akan dilakukan menggunakan script Python (`asyncio` / `multiprocessing`) atau tool load testing seperti **k6** untuk menyimulasikan *load* tinggi.

### Skenario Pengujian (Ticket War Simulation)
1. **Kondisi Awal:** Sisa kuota tiket kategori "Gold" diatur sebanyak **50 tiket**.
2. **Simulasi User:** Sebanyak **200 user** terdaftar melakukan request `POST /api/v1/bookings` secara bersamaan (konkuren) dalam rentang waktu **5 detik**.
3. **Ekspektasi Hasil:**
   * Tepat **50 transaksi** pemesanan berhasil dicatat dengan status `201 Created`.
   * **150 transaksi** sisanya ditolak dengan status `400 Bad Request` karena kuota habis.
   * Total tiket yang terpotong di database harus tepat bernilai **50** (kuota sisa menjadi **0**).
   * **Tidak boleh terjadi over-selling** (kuota menjadi negatif).
   * **Tidak boleh terjadi deadlock** pada database Postgres.

### Parameter Target Performa (SLO)
* **Average Response Time:** < 300ms untuk request booking.
* **Error Rate (selain kuota habis):** 0%.
* **Database Connection Utilization:** Tetap stabil di bawah batas maksimum pool yang ditentukan.

---

## 3. Langkah Implementasi Load Test Script
Kita akan membuat script pengujian konkurensi di `run_concurrency_test.py` dengan alur:
1. Registrasi 200 akun customer secara bertahap.
2. Login 200 customer secara bertahap untuk mengumpulkan JWT Token mereka.
3. Menyiapkan request body booking untuk masing-masing user.
4. Menggunakan `asyncio` dan `aiohttp` (atau library multithreading) untuk menembak endpoint `/api/v1/bookings` secara serentak.
5. Memverifikasi status response masing-masing request dan menghitung total transaksi sukses vs gagal.
6. Memverifikasi sisa kuota tiket di database setelah pengujian berakhir.
