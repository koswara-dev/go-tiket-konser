# Plan Testing

Rencana pengujian untuk sistem e-ticket konser:

## 1. Testing Role (Otorisasi RBAC)
Menguji keamanan API endpoint yang dibatasi oleh role tertentu (misalnya, endpoint Admin).
- **Positive Test:** Melakukan request ke `GET /api/v1/concerts` menggunakan token Admin (`adminkonser@gmail.com`). Ekspektasi: `200 OK`.
- **Negative Test:** Melakukan request ke `GET /api/v1/concerts` menggunakan token Customer. Ekspektasi: `403 Forbidden` dengan pesan `"Akses ditolak: role tidak sesuai."`.

## 2. Testing Post Data Booking Tiket (2 Customer)
Menguji fungsionalitas pemesanan tiket oleh dua customer berbeda.
- **Customer 1 Booking (Positive Test):**
  - Register customer 1, login untuk mendapatkan token JWT dan `customer_id`.
  - Melakukan `POST /api/v1/bookings` untuk memesan 1 tiket kategori Gold (ID 1).
  - Ekspektasi: `201 Created` dengan detail pemesanan.
- **Customer 2 Booking (Positive Test):**
  - Register customer 2, login untuk mendapatkan token JWT dan `customer_id`.
  - Melakukan `POST /api/v1/bookings` untuk memesan 2 tiket kategori Silver (ID 2).
  - Ekspektasi: `201 Created` dengan detail pemesanan.
- **Negative Test (Quota Exceeded):**
  - Melakukan `POST /api/v1/bookings` dengan kuantitas tiket melebihi kuota yang tersedia (misalnya, 1000 tiket).
  - Ekspektasi: `400 Bad Request` dengan error kuota tidak mencukupi.

## 3. Testing Proteksi IDOR (Insecure Direct Object Reference)
Menguji mitigasi IDOR pada endpoint pengambilan detail booking.
- **IDOR Test:**
  - Customer 2 mencoba mengakses detail booking milik Customer 1 (`GET /api/v1/bookings/{booking_id_customer_1}`) menggunakan token Customer 2.
  - Ekspektasi: `404 Not Found` (menyembunyikan faktur dari akses tidak sah demi keamanan informasi).
