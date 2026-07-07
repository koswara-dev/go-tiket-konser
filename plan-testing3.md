# Plan Testing 3

Dokumen ini berisi rencana pengujian untuk fitur SSE (Server-Sent Events) Notifications dan WebSocket Chat pada sistem e-ticket konser.

## 1. Testing SSE Notifications (Server-Sent Events)
Menguji fungsionalitas pengiriman notifikasi real-time dari server ke klien (berbasis stream).

- **Positive Test 1 (Establish Connection):**
  - Klien melakukan koneksi stream ke `GET /api/v1/notifications/stream` dengan JWT Token Customer 1.
  - Ekspektasi: Koneksi berhasil di-upgrade ke SSE, mengembalikan event `info` dengan pesan `"SSE connection established"`.

- **Positive Test 2 (Broadcast Notification on Concert Creation):**
  - Customer 1 tetap terhubung ke stream SSE.
  - Admin membuat konser baru via `POST /api/v1/concerts`.
  - Ekspektasi: Customer 1 menerima notifikasi baru berisi informasi detail konser baru tersebut dengan event `message` (broadcast).

- **Positive Test 3 (Targeted Notification on Booking Success):**
  - Customer 1 tetap terhubung ke stream SSE.
  - Customer 1 melakukan pemesanan tiket via `POST /api/v1/bookings`.
  - Ekspektasi: Customer 1 menerima notifikasi transaksi sukses bertipe data `Booking Berhasil` (targeted).

- **Negative Test (Unauthorized SSE Connection):**
  - Akses `GET /api/v1/notifications/stream` tanpa menyertakan JWT token.
  - Ekspektasi: Response status `401 Unauthorized` dengan pesan `"Invalid or missing token"`.

## 2. Testing WebSocket Chat (Real-time Messaging)
Menguji keandalan sistem chat real-time antara Customer dan Admin melalui protokol WebSocket.

- **Positive Test 1 (Customer & Admin Handshake):**
  - Customer 1 melakukan koneksi WebSocket ke `/api/v1/chat/ws?token=<customer_token>`.
  - Admin melakukan koneksi WebSocket ke `/api/v1/chat/ws?token=<admin_token>&room_id=<customer1_user_id>`.
  - Ekspektasi: Kedua koneksi berhasil di-upgrade (HTTP 101 Switching Protocols).

- **Positive Test 2 (Bidirectional Chatting):**
  - Customer 1 mengirim pesan: `{"message": "Halo, saya butuh bantuan terkait tiket saya."}`.
  - Admin menerima pesan real-time tersebut dari room `<customer1_user_id>`.
  - Admin membalas pesan: `{"message": "Halo, ada yang bisa kami bantu?"}`.
  - Customer 1 menerima balasan real-time tersebut.
  - Ekspektasi: Kedua pesan berhasil terkirim, diterima, dan disimpan ke MongoDB.

- **Positive Test 3 (Get Room Messages & Rooms List History):**
  - Mengakses endpoint HTTP `GET /api/v1/chat/rooms` (Admin only). Ekspektasi: `200 OK` dengan daftar room chat aktif beserta pesan terakhir.
  - Mengakses endpoint HTTP `GET /api/v1/chat/rooms/{roomId}/messages` (Admin & Pemilik Room). Ekspektasi: `200 OK` dengan riwayat pesan lengkap.

- **Negative Test 1 (Unauthenticated Connection):**
  - Mencoba koneksi WebSocket `/api/v1/chat/ws` tanpa parameter token.
  - Ekspektasi: Upgrade gagal, mengembalikan status HTTP `401 Unauthorized`.

- **Negative Test 2 (IDOR on Room Messages Access):**
  - Customer 2 mencoba mengakses riwayat pesan Customer 1 via `GET /api/v1/chat/rooms/{customer1_user_id}/messages` dengan token Customer 2.
  - Ekspektasi: Response status `403 Forbidden` dengan pesan `"Akses ditolak: Anda tidak dapat mengakses percakapan ini"`.
