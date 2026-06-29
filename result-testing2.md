# Result Testing 2

Dokumen ini berisi hasil pengujian yang dilakukan untuk memvalidasi fitur Pagination & Search pada daftar konser (termasuk 1 kasus negatif), serta pengujian mekanisme Force Logout pada 10 akun Customer.

## 1. Testing Pagination & Search (Daftar Konser)
### Positive Test: Pagination & Limit
- **Request:** `GET /api/v1/concerts?page=1&limit=2`
- **Status Code:** 200
- **Response:**
```json
{
    "success": true,
    "message": "Data berhasil diambil",
    "data": [
        {
            "id": 21,
            "title": "World Tour Taylor Swift - Live in Surabaya Part 15",
            "description": "Konser megah dan eksklusif bersama Taylor Swift di Surabaya.",
            "date": "2026-07-21",
            "venue": "Stadion Gelora Bung Tomo",
            "status": "upcoming",
            "created_at": "2026-06-29T20:55:05.786845+07:00",
            "updated_at": "2026-06-29T20:55:05.786845+07:00"
        },
        {
            "id": 20,
            "title": "World Tour Bruno Mars - Live in Surabaya Part 14",
            "description": "Konser megah dan eksklusif bersama Bruno Mars di Surabaya.",
            "date": "2026-07-20",
            "venue": "Stadion Gelora Bung Tomo",
            "status": "upcoming",
            "created_at": "2026-06-29T20:55:05.677619+07:00",
            "updated_at": "2026-06-29T20:55:05.677619+07:00"
        }
    ],
    "meta": {
        "page": 1,
        "limit": 2,
        "total_data": 21,
        "total_page": 11
    }
}
```
### Positive Test: Search (Coldplay)
- **Request:** `GET /api/v1/concerts?search=Coldplay`
- **Status Code:** 200
- **Response:**
```json
{
    "success": true,
    "message": "Data berhasil diambil",
    "data": [
        {
            "id": 19,
            "title": "World Tour Coldplay - Live in Jakarta Part 13",
            "description": "Konser megah dan eksklusif bersama Coldplay di Gelora Bung Karno.",
            "date": "2026-07-19",
            "venue": "Gelora Bung Karno",
            "status": "upcoming",
            "created_at": "2026-06-29T20:55:05.567717+07:00",
            "updated_at": "2026-06-29T20:55:05.567717+07:00"
        },
        {
            "id": 1,
            "title": "Coldplay Music of the Spheres World Tour Jakarta",
            "description": "Konser perdana band asal Inggris, Coldplay, di Indonesia yang memukau ratusan ribu penonton dengan gelang Xyloband yang menyala warna-warni.",
            "date": "2023-11-16",
            "venue": "Stadion Utama Gelora Bung Karno (SUGBK), Jakarta",
            "status": "completed",
            "created_at": "2026-06-22T19:56:26.295337+07:00",
            "updated_at": "2026-06-22T19:56:26.295337+07:00"
        }
    ],
    "meta": {
        "page": 1,
        "limit": 10,
        "total_data": 2,
        "total_page": 1
    }
}
```
### Positive Test: Sorting (date_asc)
- **Request:** `GET /api/v1/concerts?sort=date_asc`
- **Status Code:** 200
- **Response:**
```json
{
    "success": true,
    "message": "Data berhasil diambil",
    "data": [
        {
            "id": 3,
            "title": "Metallica Live in Jakarta 2013",
            "description": "Konser sejarah kembalinya raja thrash metal dunia ke Indonesia setelah penantian 20 tahun, dihadiri oleh puluhan ribu pecinta musik cadas dari berbagai generasi.",
            "date": "2013-08-26",
            "venue": "Stadion Utama Gelora Bung Karno (SUGBK), Jakarta",
            "status": "completed",
            "created_at": "2026-06-22T19:56:26.295337+07:00",
            "updated_at": "2026-06-22T19:56:26.295337+07:00"
        },
        {
            "id": 2,
            "title": "Blackpink [Born Pink] World Tour Jakarta",
            "description": "Konser megah dari girlgroup K-Pop fenomenal, Blackpink, yang berhasil meremajakan Jakarta menjadi lautan cahaya merah muda selama dua hari berturut-turut.",
            "date": "2023-03-12",
            "venue": "Stadion Utama Gelora Bung Karno (SUGBK), Jakarta",
            "status": "completed",
            "created_at": "2026-06-22T19:56:26.295337+07:00",
            "updated_at": "2026-06-22T19:56:26.295337+07:00"
        },
        {
            "id": 1,
            "title": "Coldplay Music of the Spheres World Tour Jakarta",
            "description": "Konser perdana band asal Inggris, Coldplay, di Indonesia yang memukau ratusan ribu penonton dengan gelang Xyloband yang menyala warna-warni.",
            "date": "2023-11-16",
            "venue": "Stadion Utama Gelora Bung Karno (SUGBK), Jakarta",
            "status": "completed",
            "created_at": "2026-06-22T19:56:26.295337+07:00",
            "updated_at": "2026-06-22T19:56:26.295337+07:00"
        },
        {
            "id": 5,
            "title": "Pesta Rakyat Dewa 19 - 30 Tahun Berkarya",
            "description": "Konser selebrasi 3 dekade salah satu band rock terbesar di Indonesia, Dewa 19, yang memboyong 4 vokalis dan 5 drummer dalam satu panggung.",
            "date": "2026-06-28",
            "venue": "Stadion Utama Gelora Bung Karno (SUGBK), Jakarta",
            "status": "active",
            "created_at": "2026-06-22T19:56:26.295337+07:00",
            "updated_at": "2026-06-22T19:56:26.295337+07:00"
        },
        {
            "id": 4,
            "title": "Bruno Mars Live in Jakarta 2026",
            "description": "Konser tur dunia dari solois legendaris Bruno Mars yang membawakan deretan lagu hitsnya dengan koreografi dan vokal yang sangat enerjik.",
            "date": "2026-06-28",
            "venue": "Jakarta International Stadium (JIS), Jakarta",
            "status": "active",
            "created_at": "2026-06-22T19:56:26.295337+07:00",
            "updated_at": "2026-06-22T19:56:26.295337+07:00"
        },
        {
            "id": 6,
            "title": "Konser Dewa 19",
            "description": "Konser Reuni Dewa 19",
            "date": "2026-07-05",
            "venue": "Stadion Utama GBK",
            "status": "upcoming",
            "created_at": "2026-06-24T20:54:51.982816+07:00",
            "updated_at": "2026-06-24T20:54:51.982816+07:00"
        },
        {
            "id": 7,
            "title": "World Tour Bruno Mars - Live in Jakarta Part 1",
            "description": "Konser megah dan eksklusif bersama Bruno Mars di ICE BSD.",
            "date": "2026-07-07",
            "venue": "ICE BSD",
            "status": "active",
            "created_at": "2026-06-29T20:55:04.223186+07:00",
            "updated_at": "2026-06-29T20:55:04.223186+07:00"
        },
        {
            "id": 8,
            "title": "World Tour Taylor Swift - Live in Jakarta Part 2",
            "description": "Konser megah dan eksklusif bersama Taylor Swift di Jakarta International Stadium.",
            "date": "2026-07-08",
            "venue": "Jakarta International Stadium",
            "status": "active",
            "created_at": "2026-06-29T20:55:04.352137+07:00",
            "updated_at": "2026-06-29T20:55:04.352137+07:00"
        },
        {
            "id": 9,
            "title": "World Tour Tulus - Live in Jakarta Part 3",
            "description": "Konser megah dan eksklusif bersama Tulus di Jakarta International Stadium.",
            "date": "2026-07-09",
            "venue": "Jakarta International Stadium",
            "status": "upcoming",
            "created_at": "2026-06-29T20:55:04.465606+07:00",
            "updated_at": "2026-06-29T20:55:04.465606+07:00"
        },
        {
            "id": 10,
            "title": "World Tour Tulus - Live in Jakarta Part 4",
            "description": "Konser megah dan eksklusif bersama Tulus di Jakarta International Stadium.",
            "date": "2026-07-10",
            "venue": "Jakarta International Stadium",
            "status": "active",
            "created_at": "2026-06-29T20:55:04.576337+07:00",
            "updated_at": "2026-06-29T20:55:04.576337+07:00"
        }
    ],
    "meta": {
        "page": 1,
        "limit": 10,
        "total_data": 21,
        "total_page": 3
    }
}
```
### Negative Test: Pagination dengan Page Negatif
- **Request:** `GET /api/v1/concerts?page=-1`
- **Status Code:** 400
- **Response:**
```json
{
    "success": false,
    "message": "Validasi parameter pencarian gagal",
    "data": "Key: 'ConcertQueryRequest.Page' Error:Field validation for 'Page' failed on the 'gte' tag"
}
```

## 2. Testing Force Logout (10 Akun Customer)
Melakukan pengujian force logout untuk memastikan token JWT yang sudah di-blacklist/logout tidak dapat digunakan kembali.

### Akun Customer 1 (cust_logout_1_888403@gmail.com)
- **Langkah 1: Registrasi Akun**
  - **Endpoint:** `POST /api/v1/register`
  - **Status Code:** 201
  - **Response:** `{"data":{"id":19,"full_name":"Customer Logout Test 1","email":"cust_logout_1_888403@gmail.com","role":"customer","customer_id":10},"message":"User registered successfully","success":true}`
- **Langkah 2: Login Akun**
  - **Endpoint:** `POST /api/v1/login`
  - **Status Code:** 200
  - **Response:** `{"data":{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxOSwiY3VzdG9tZXJfaWQiOjEwLCJlbWFpbCI6ImN1c3RfbG9nb3V0XzFfODg4NDAzQGdtYWlsLmNvbSIsInJvbGUiOiJjdXN0b21lciIsImV4cCI6MTc4MjgyOTM0NCwiaWF0IjoxNzgyNzQyOTQ0fQ.1PRKCI_PJEKItdthiDubUO11faJwisQe_0iWU8dPkNs"},"message":"Login berhasil","success":true}`
- **Langkah 3: Akses Profil Sebelum Logout**
  - **Endpoint:** `GET /api/v1/profile`
  - **Status Code:** 200
  - **Response:** `{"data":{"id":19,"full_name":"Customer Logout Test 1","email":"cust_logout_1_888403@gmail.com","role":"customer","customer_id":10},"message":"User profile retrieved successfully","success":true}`
- **Langkah 4: Force Logout**
  - **Endpoint:** `POST /api/v1/logout`
  - **Status Code:** 200
  - **Response:** `{"message":"User logged out successfully","success":true}`
- **Langkah 5: Akses Profil Setelah Logout (Blacklisted)**
  - **Endpoint:** `GET /api/v1/profile`
  - **Status Code:** 401
  - **Response:** `{"error":"Token is blacklisted"}`

### Akun Customer 2 (cust_logout_2_270018@gmail.com)
- **Langkah 1: Registrasi Akun**
  - **Endpoint:** `POST /api/v1/register`
  - **Status Code:** 201
  - **Response:** `{"data":{"id":20,"full_name":"Customer Logout Test 2","email":"cust_logout_2_270018@gmail.com","role":"customer","customer_id":11},"message":"User registered successfully","success":true}`
- **Langkah 2: Login Akun**
  - **Endpoint:** `POST /api/v1/login`
  - **Status Code:** 200
  - **Response:** `{"data":{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyMCwiY3VzdG9tZXJfaWQiOjExLCJlbWFpbCI6ImN1c3RfbG9nb3V0XzJfMjcwMDE4QGdtYWlsLmNvbSIsInJvbGUiOiJjdXN0b21lciIsImV4cCI6MTc4MjgyOTM0NCwiaWF0IjoxNzgyNzQyOTQ0fQ.6PX7aWbsoUhoOB0vXrkMylgsWVaMkZrwiotzKSuvP-0"},"message":"Login berhasil","success":true}`
- **Langkah 3: Akses Profil Sebelum Logout**
  - **Endpoint:** `GET /api/v1/profile`
  - **Status Code:** 200
  - **Response:** `{"data":{"id":20,"full_name":"Customer Logout Test 2","email":"cust_logout_2_270018@gmail.com","role":"customer","customer_id":11},"message":"User profile retrieved successfully","success":true}`
- **Langkah 4: Force Logout**
  - **Endpoint:** `POST /api/v1/logout`
  - **Status Code:** 200
  - **Response:** `{"message":"User logged out successfully","success":true}`
- **Langkah 5: Akses Profil Setelah Logout (Blacklisted)**
  - **Endpoint:** `GET /api/v1/profile`
  - **Status Code:** 401
  - **Response:** `{"error":"Token is blacklisted"}`

### Akun Customer 3 (cust_logout_3_613142@gmail.com)
- **Langkah 1: Registrasi Akun**
  - **Endpoint:** `POST /api/v1/register`
  - **Status Code:** 201
  - **Response:** `{"data":{"id":21,"full_name":"Customer Logout Test 3","email":"cust_logout_3_613142@gmail.com","role":"customer","customer_id":12},"message":"User registered successfully","success":true}`
- **Langkah 2: Login Akun**
  - **Endpoint:** `POST /api/v1/login`
  - **Status Code:** 200
  - **Response:** `{"data":{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyMSwiY3VzdG9tZXJfaWQiOjEyLCJlbWFpbCI6ImN1c3RfbG9nb3V0XzNfNjEzMTQyQGdtYWlsLmNvbSIsInJvbGUiOiJjdXN0b21lciIsImV4cCI6MTc4MjgyOTM0NCwiaWF0IjoxNzgyNzQyOTQ0fQ.zHBMwlvQK_PE11VVGz2Jaabl2D3FfVAp_1Fv9MY_QFk"},"message":"Login berhasil","success":true}`
- **Langkah 3: Akses Profil Sebelum Logout**
  - **Endpoint:** `GET /api/v1/profile`
  - **Status Code:** 200
  - **Response:** `{"data":{"id":21,"full_name":"Customer Logout Test 3","email":"cust_logout_3_613142@gmail.com","role":"customer","customer_id":12},"message":"User profile retrieved successfully","success":true}`
- **Langkah 4: Force Logout**
  - **Endpoint:** `POST /api/v1/logout`
  - **Status Code:** 200
  - **Response:** `{"message":"User logged out successfully","success":true}`
- **Langkah 5: Akses Profil Setelah Logout (Blacklisted)**
  - **Endpoint:** `GET /api/v1/profile`
  - **Status Code:** 401
  - **Response:** `{"error":"Token is blacklisted"}`

### Akun Customer 4 (cust_logout_4_852223@gmail.com)
- **Langkah 1: Registrasi Akun**
  - **Endpoint:** `POST /api/v1/register`
  - **Status Code:** 201
  - **Response:** `{"data":{"id":22,"full_name":"Customer Logout Test 4","email":"cust_logout_4_852223@gmail.com","role":"customer","customer_id":13},"message":"User registered successfully","success":true}`
- **Langkah 2: Login Akun**
  - **Endpoint:** `POST /api/v1/login`
  - **Status Code:** 200
  - **Response:** `{"data":{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyMiwiY3VzdG9tZXJfaWQiOjEzLCJlbWFpbCI6ImN1c3RfbG9nb3V0XzRfODUyMjIzQGdtYWlsLmNvbSIsInJvbGUiOiJjdXN0b21lciIsImV4cCI6MTc4MjgyOTM0NSwiaWF0IjoxNzgyNzQyOTQ1fQ.89JpFvT5j9sy_P09OAqKEPR1zP4AmglSXEm8N5_GE2w"},"message":"Login berhasil","success":true}`
- **Langkah 3: Akses Profil Sebelum Logout**
  - **Endpoint:** `GET /api/v1/profile`
  - **Status Code:** 200
  - **Response:** `{"data":{"id":22,"full_name":"Customer Logout Test 4","email":"cust_logout_4_852223@gmail.com","role":"customer","customer_id":13},"message":"User profile retrieved successfully","success":true}`
- **Langkah 4: Force Logout**
  - **Endpoint:** `POST /api/v1/logout`
  - **Status Code:** 200
  - **Response:** `{"message":"User logged out successfully","success":true}`
- **Langkah 5: Akses Profil Setelah Logout (Blacklisted)**
  - **Endpoint:** `GET /api/v1/profile`
  - **Status Code:** 401
  - **Response:** `{"error":"Token is blacklisted"}`

### Akun Customer 5 (cust_logout_5_431650@gmail.com)
- **Langkah 1: Registrasi Akun**
  - **Endpoint:** `POST /api/v1/register`
  - **Status Code:** 201
  - **Response:** `{"data":{"id":23,"full_name":"Customer Logout Test 5","email":"cust_logout_5_431650@gmail.com","role":"customer","customer_id":14},"message":"User registered successfully","success":true}`
- **Langkah 2: Login Akun**
  - **Endpoint:** `POST /api/v1/login`
  - **Status Code:** 200
  - **Response:** `{"data":{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyMywiY3VzdG9tZXJfaWQiOjE0LCJlbWFpbCI6ImN1c3RfbG9nb3V0XzVfNDMxNjUwQGdtYWlsLmNvbSIsInJvbGUiOiJjdXN0b21lciIsImV4cCI6MTc4MjgyOTM0NSwiaWF0IjoxNzgyNzQyOTQ1fQ.nXgFkJDTi-KHlpDXNRL3iWuCOFdOImGjbRUzFmOc1o4"},"message":"Login berhasil","success":true}`
- **Langkah 3: Akses Profil Sebelum Logout**
  - **Endpoint:** `GET /api/v1/profile`
  - **Status Code:** 200
  - **Response:** `{"data":{"id":23,"full_name":"Customer Logout Test 5","email":"cust_logout_5_431650@gmail.com","role":"customer","customer_id":14},"message":"User profile retrieved successfully","success":true}`
- **Langkah 4: Force Logout**
  - **Endpoint:** `POST /api/v1/logout`
  - **Status Code:** 200
  - **Response:** `{"message":"User logged out successfully","success":true}`
- **Langkah 5: Akses Profil Setelah Logout (Blacklisted)**
  - **Endpoint:** `GET /api/v1/profile`
  - **Status Code:** 401
  - **Response:** `{"error":"Token is blacklisted"}`

### Akun Customer 6 (cust_logout_6_672848@gmail.com)
- **Langkah 1: Registrasi Akun**
  - **Endpoint:** `POST /api/v1/register`
  - **Status Code:** 201
  - **Response:** `{"data":{"id":24,"full_name":"Customer Logout Test 6","email":"cust_logout_6_672848@gmail.com","role":"customer","customer_id":15},"message":"User registered successfully","success":true}`
- **Langkah 2: Login Akun**
  - **Endpoint:** `POST /api/v1/login`
  - **Status Code:** 200
  - **Response:** `{"data":{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyNCwiY3VzdG9tZXJfaWQiOjE1LCJlbWFpbCI6ImN1c3RfbG9nb3V0XzZfNjcyODQ4QGdtYWlsLmNvbSIsInJvbGUiOiJjdXN0b21lciIsImV4cCI6MTc4MjgyOTM0NSwiaWF0IjoxNzgyNzQyOTQ1fQ.iZ3DhmE2ZWb5ihd_RuR18jzkhLtjDB7DDS2zKpOW1o4"},"message":"Login berhasil","success":true}`
- **Langkah 3: Akses Profil Sebelum Logout**
  - **Endpoint:** `GET /api/v1/profile`
  - **Status Code:** 200
  - **Response:** `{"data":{"id":24,"full_name":"Customer Logout Test 6","email":"cust_logout_6_672848@gmail.com","role":"customer","customer_id":15},"message":"User profile retrieved successfully","success":true}`
- **Langkah 4: Force Logout**
  - **Endpoint:** `POST /api/v1/logout`
  - **Status Code:** 200
  - **Response:** `{"message":"User logged out successfully","success":true}`
- **Langkah 5: Akses Profil Setelah Logout (Blacklisted)**
  - **Endpoint:** `GET /api/v1/profile`
  - **Status Code:** 401
  - **Response:** `{"error":"Token is blacklisted"}`

### Akun Customer 7 (cust_logout_7_776344@gmail.com)
- **Langkah 1: Registrasi Akun**
  - **Endpoint:** `POST /api/v1/register`
  - **Status Code:** 201
  - **Response:** `{"data":{"id":25,"full_name":"Customer Logout Test 7","email":"cust_logout_7_776344@gmail.com","role":"customer","customer_id":16},"message":"User registered successfully","success":true}`
- **Langkah 2: Login Akun**
  - **Endpoint:** `POST /api/v1/login`
  - **Status Code:** 200
  - **Response:** `{"data":{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyNSwiY3VzdG9tZXJfaWQiOjE2LCJlbWFpbCI6ImN1c3RfbG9nb3V0XzdfNzc2MzQ0QGdtYWlsLmNvbSIsInJvbGUiOiJjdXN0b21lciIsImV4cCI6MTc4MjgyOTM0NiwiaWF0IjoxNzgyNzQyOTQ2fQ.tRvOWLg0SDphkCmVegXKfACKiC7ic0trTmszhjSkOk8"},"message":"Login berhasil","success":true}`
- **Langkah 3: Akses Profil Sebelum Logout**
  - **Endpoint:** `GET /api/v1/profile`
  - **Status Code:** 200
  - **Response:** `{"data":{"id":25,"full_name":"Customer Logout Test 7","email":"cust_logout_7_776344@gmail.com","role":"customer","customer_id":16},"message":"User profile retrieved successfully","success":true}`
- **Langkah 4: Force Logout**
  - **Endpoint:** `POST /api/v1/logout`
  - **Status Code:** 200
  - **Response:** `{"message":"User logged out successfully","success":true}`
- **Langkah 5: Akses Profil Setelah Logout (Blacklisted)**
  - **Endpoint:** `GET /api/v1/profile`
  - **Status Code:** 401
  - **Response:** `{"error":"Token is blacklisted"}`

### Akun Customer 8 (cust_logout_8_806250@gmail.com)
- **Langkah 1: Registrasi Akun**
  - **Endpoint:** `POST /api/v1/register`
  - **Status Code:** 201
  - **Response:** `{"data":{"id":26,"full_name":"Customer Logout Test 8","email":"cust_logout_8_806250@gmail.com","role":"customer","customer_id":17},"message":"User registered successfully","success":true}`
- **Langkah 2: Login Akun**
  - **Endpoint:** `POST /api/v1/login`
  - **Status Code:** 200
  - **Response:** `{"data":{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyNiwiY3VzdG9tZXJfaWQiOjE3LCJlbWFpbCI6ImN1c3RfbG9nb3V0XzhfODA2MjUwQGdtYWlsLmNvbSIsInJvbGUiOiJjdXN0b21lciIsImV4cCI6MTc4MjgyOTM0NiwiaWF0IjoxNzgyNzQyOTQ2fQ.Gr_rVFCRFSLX1i1VfZ2DdWkNiRDca_fEQJHh2mHFZRM"},"message":"Login berhasil","success":true}`
- **Langkah 3: Akses Profil Sebelum Logout**
  - **Endpoint:** `GET /api/v1/profile`
  - **Status Code:** 200
  - **Response:** `{"data":{"id":26,"full_name":"Customer Logout Test 8","email":"cust_logout_8_806250@gmail.com","role":"customer","customer_id":17},"message":"User profile retrieved successfully","success":true}`
- **Langkah 4: Force Logout**
  - **Endpoint:** `POST /api/v1/logout`
  - **Status Code:** 200
  - **Response:** `{"message":"User logged out successfully","success":true}`
- **Langkah 5: Akses Profil Setelah Logout (Blacklisted)**
  - **Endpoint:** `GET /api/v1/profile`
  - **Status Code:** 401
  - **Response:** `{"error":"Token is blacklisted"}`

### Akun Customer 9 (cust_logout_9_618600@gmail.com)
- **Langkah 1: Registrasi Akun**
  - **Endpoint:** `POST /api/v1/register`
  - **Status Code:** 201
  - **Response:** `{"data":{"id":27,"full_name":"Customer Logout Test 9","email":"cust_logout_9_618600@gmail.com","role":"customer","customer_id":18},"message":"User registered successfully","success":true}`
- **Langkah 2: Login Akun**
  - **Endpoint:** `POST /api/v1/login`
  - **Status Code:** 200
  - **Response:** `{"data":{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyNywiY3VzdG9tZXJfaWQiOjE4LCJlbWFpbCI6ImN1c3RfbG9nb3V0XzlfNjE4NjAwQGdtYWlsLmNvbSIsInJvbGUiOiJjdXN0b21lciIsImV4cCI6MTc4MjgyOTM0NiwiaWF0IjoxNzgyNzQyOTQ2fQ.hHp0MBg9xbTXaeZPZqXDaWgx6p2fGcRoG0yi4ttDmOI"},"message":"Login berhasil","success":true}`
- **Langkah 3: Akses Profil Sebelum Logout**
  - **Endpoint:** `GET /api/v1/profile`
  - **Status Code:** 200
  - **Response:** `{"data":{"id":27,"full_name":"Customer Logout Test 9","email":"cust_logout_9_618600@gmail.com","role":"customer","customer_id":18},"message":"User profile retrieved successfully","success":true}`
- **Langkah 4: Force Logout**
  - **Endpoint:** `POST /api/v1/logout`
  - **Status Code:** 200
  - **Response:** `{"message":"User logged out successfully","success":true}`
- **Langkah 5: Akses Profil Setelah Logout (Blacklisted)**
  - **Endpoint:** `GET /api/v1/profile`
  - **Status Code:** 401
  - **Response:** `{"error":"Token is blacklisted"}`

### Akun Customer 10 (cust_logout_10_695254@gmail.com)
- **Langkah 1: Registrasi Akun**
  - **Endpoint:** `POST /api/v1/register`
  - **Status Code:** 201
  - **Response:** `{"data":{"id":28,"full_name":"Customer Logout Test 10","email":"cust_logout_10_695254@gmail.com","role":"customer","customer_id":19},"message":"User registered successfully","success":true}`
- **Langkah 2: Login Akun**
  - **Endpoint:** `POST /api/v1/login`
  - **Status Code:** 200
  - **Response:** `{"data":{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyOCwiY3VzdG9tZXJfaWQiOjE5LCJlbWFpbCI6ImN1c3RfbG9nb3V0XzEwXzY5NTI1NEBnbWFpbC5jb20iLCJyb2xlIjoiY3VzdG9tZXIiLCJleHAiOjE3ODI4MjkzNDYsImlhdCI6MTc4Mjc0Mjk0Nn0.lp6ZEKKcc6UtLA3g1e6kph4-j_Ji_q8OYzA37jZXe2U"},"message":"Login berhasil","success":true}`
- **Langkah 3: Akses Profil Sebelum Logout**
  - **Endpoint:** `GET /api/v1/profile`
  - **Status Code:** 200
  - **Response:** `{"data":{"id":28,"full_name":"Customer Logout Test 10","email":"cust_logout_10_695254@gmail.com","role":"customer","customer_id":19},"message":"User profile retrieved successfully","success":true}`
- **Langkah 4: Force Logout**
  - **Endpoint:** `POST /api/v1/logout`
  - **Status Code:** 200
  - **Response:** `{"message":"User logged out successfully","success":true}`
- **Langkah 5: Akses Profil Setelah Logout (Blacklisted)**
  - **Endpoint:** `GET /api/v1/profile`
  - **Status Code:** 401
  - **Response:** `{"error":"Token is blacklisted"}`

