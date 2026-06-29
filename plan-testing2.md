# Plan Testing 2

Dokumen ini berisi rencana pengujian untuk fitur Pagination & Search pada daftar konser serta pengujian mekanisme Force Logout pada 10 akun Customer.

## 1. Testing Pagination & Search (Daftar Konser)
Menguji fungsionalitas pencarian, pembatasan jumlah data per halaman, dan pengurutan pada endpoint `GET /api/v1/concerts`.

- **Positive Test (Pagination & Limit):**
  - Mengirimkan request `GET /api/v1/concerts?page=1&limit=2` dengan menyertakan API Key.
  - Ekspektasi: Response status `200 OK`, jumlah data konser maksimal 2, dan metadata pagination (`page=1`, `limit=2`, `total_data`, `total_page`) sesuai dengan jumlah konser di database.

- **Positive Test (Search):**
  - Mengirimkan request `GET /api/v1/concerts?search=Coldplay` dengan menyertakan API Key.
  - Ekspektasi: Response status `200 OK`, data konser yang dikembalikan hanya yang memiliki judul atau lokasi mengandung kata "Coldplay".

- **Positive Test (Sorting):**
  - Mengirimkan request `GET /api/v1/concerts?sort=date_asc` dengan menyertakan API Key.
  - Ekspektasi: Response status `200 OK`, daftar konser terurut berdasarkan tanggal pelaksanaan dari yang terlama ke terbaru.

- **Negative Test (Invalid Query Parameter):**
  - Mengirimkan request `GET /api/v1/concerts?page=-1` atau `GET /api/v1/concerts?sort=invalid_sort` dengan menyertakan API Key.
  - Ekspektasi: Response status `400 Bad Request` dengan pesan error validasi query parameter.

## 2. Testing Force Logout (10 Akun Customer)
Menguji keandalan token blacklist dengan membuat 10 akun customer baru, melakukan login untuk mendapatkan token JWT, memverifikasi akses profile, melakukan logout (inaktivasi token), dan memastikan token tersebut sudah tidak dapat digunakan kembali.

- **Alur Pengujian per Akun (1 s.d. 10):**
  1. **Register:** Melakukan registrasi akun customer baru dengan email unik (`cust_logout_X@gmail.com`). Ekspektasi: `201 Created`.
  2. **Login:** Melakukan login untuk mendapatkan token JWT. Ekspektasi: `200 OK`.
  3. **Check Profile (Sebelum Logout):** Mengakses `GET /api/v1/profile` menggunakan JWT token. Ekspektasi: `200 OK` dengan data profil user.
  4. **Force Logout:** Mengakses `POST /api/v1/logout` menggunakan JWT token. Ekspektasi: `200 OK` dengan pesan logout berhasil.
  5. **Check Profile (Setelah Logout/Blacklisted):** Mengakses kembali `GET /api/v1/profile` menggunakan token yang sama. Ekspektasi: `401 Unauthorized` dengan pesan `"Token is blacklisted"`.
