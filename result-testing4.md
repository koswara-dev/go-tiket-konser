# Result Concurrency Test (Ticket War Simulation)

Dokumen ini berisi hasil pengujian konkurensi tinggi untuk memvalidasi proteksi *race condition* dan keandalan kuota tiket under high load.

## Skenario Pengujian
- **Target Kategori Tiket:** Gold (Price: Rp 1.000.000)
- **Kuota Tiket Awal:** 500
- **Jumlah Request Konkuren:** 1000 request dari user berbeda secara serentak
- **Mekanisme Proteksi:** SQL Pessimistic Locking (`SELECT ... FOR UPDATE`)

## Hasil Pengujian
- **Total Request Terkirim:** 1000
- **Transaksi Sukses (201 Created):** 500 (Sesuai dengan kuota awal)
- **Transaksi Gagal (400 Bad Request):** 500 (Karena kuota habis)
- **Status Response Lainnya:** 0
- **Sisa Kuota Tiket Akhir di DB:** 0
- **Rata-rata Latency Request:** 7523.49 ms

### Rincian Kegagalan (Error Messages)
- `kuota tiket 'Gold' tidak mencukupi (Tersisa: 0, Permintaan: 1)`: 500 kali

## Analisis & Kesimpulan
1. **Pessimistic Locking Sukses:** Proteksi race condition menggunakan `SELECT ... FOR UPDATE` berhasil menahan transaksi ganda. Tepat 500 tiket terpesan dan sisa kuota akhir adalah 0. Tidak ada kasus *over-selling* (kuota tidak menjadi negatif).
2. **Koneksi Database Stabil:** Pengaturan connection pool baru (`MaxOpenConns=50`, `MaxIdleConns=10`) berhasil menangani lonjakan kueri konkuren secara stabil tanpa ada error timeout koneksi database.
3. **Sistem Terbukti Aman:** backend dapat diandalkan untuk skenario *ticket war* dengan keamanan data kuota tetap 100% konsisten.
