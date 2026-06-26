# Result Testing

Dokumen ini berisi hasil pengujian yang dilakukan untuk memvalidasi peran (RBAC), pembuatan booking tiket untuk 2 customer (termasuk kasus negatif), dan proteksi IDOR.

## 1. Testing Role (Otorisasi RBAC)
### Positive Test: Admin Mengakses Daftar Konser
- **Request:** `GET /api/v1/concerts` dengan Token Admin
- **Status Code:** 200
- **Response:**
```json
[
    {
        "id": 1,
        "title": "Coldplay Music of the Spheres World Tour Jakarta",
        "description": "Konser perdana band asal Inggris, Coldplay, di Indonesia yang memukau ratusan ribu penonton dengan gelang Xyloband yang menyala warna-warni.",
        "date": "2023-11-16T03:00:00+07:00",
        "venue": "Stadion Utama Gelora Bung Karno (SUGBK), Jakarta",
        "status": "completed",
        "created_at": "2026-06-22T19:56:26.295337+07:00",
        "updated_at": "2026-06-22T19:56:26.295337+07:00",
        "deleted_at": "0001-01-01T07:00:00+07:00"
    },
    {
        "id": 2,
        "title": "Blackpink [Born Pink] World Tour Jakarta",
        "description": "Konser megah dari girlgroup K-Pop fenomenal, Blackpink, yang berhasil meremajakan Jakarta menjadi lautan cahaya merah muda selama dua hari berturut-turut.",
        "date": "2023-03-12T02:00:00+07:00",
        "venue": "Stadion Utama Gelora Bung Karno (SUGBK), Jakarta",
        "status": "completed",
        "created_at": "2026-06-22T19:56:26.295337+07:00",
        "updated_at": "2026-06-22T19:56:26.295337+07:00",
        "deleted_at": "0001-01-01T07:00:00+07:00"
    },
    {
        "id": 3,
        "title": "Metallica Live in Jakarta 2013",
        "description": "Konser sejarah kembalinya raja thrash metal dunia ke Indonesia setelah penantian 20 tahun, dihadiri oleh puluhan ribu pecinta musik cadas dari berbagai generasi.",
        "date": "2013-08-26T03:00:00+07:00",
        "venue": "Stadion Utama Gelora Bung Karno (SUGBK), Jakarta",
        "status": "completed",
        "created_at": "2026-06-22T19:56:26.295337+07:00",
        "updated_at": "2026-06-22T19:56:26.295337+07:00",
        "deleted_at": "0001-01-01T07:00:00+07:00"
    },
    {
        "id": 6,
        "title": "Konser Dewa 19",
        "description": "Konser Reuni Dewa 19",
        "date": "2026-07-05T07:00:00+07:00",
        "venue": "Stadion Utama GBK",
        "status": "upcoming",
        "created_at": "2026-06-24T20:54:51.982816+07:00",
        "updated_at": "2026-06-24T20:54:51.982816+07:00",
        "deleted_at": "0001-01-01T00:00:00Z"
    },
    {
        "id": 4,
        "title": "Bruno Mars Live in Jakarta 2026",
        "description": "Konser tur dunia dari solois legendaris Bruno Mars yang membawakan deretan lagu hitsnya dengan koreografi dan vokal yang sangat enerjik.",
        "date": "2026-06-28T03:00:00+07:00",
        "venue": "Jakarta International Stadium (JIS), Jakarta",
        "status": "active",
        "created_at": "2026-06-22T19:56:26.295337+07:00",
        "updated_at": "2026-06-22T19:56:26.295337+07:00",
        "deleted_at": "0001-01-01T07:00:00+07:00"
    },
    {
        "id": 5,
        "title": "Pesta Rakyat Dewa 19 - 30 Tahun Berkarya",
        "description": "Konser selebrasi 3 dekade salah satu band rock terbesar di Indonesia, Dewa 19, yang memboyong 4 vokalis dan 5 drummer dalam satu panggung.",
        "date": "2026-06-28T02:30:00+07:00",
        "venue": "Stadion Utama Gelora Bung Karno (SUGBK), Jakarta",
        "status": "active",
        "created_at": "2026-06-22T19:56:26.295337+07:00",
        "updated_at": "2026-06-22T19:56:26.295337+07:00",
        "deleted_at": "0001-01-01T07:00:00+07:00"
    }
]
```
### Negative Test: Customer Mengakses Daftar Konser
- **Request:** `GET /api/v1/concerts` dengan Token Customer
- **Status Code:** 403
- **Response:**
```json
{
    "message": "Akses ditolak: role tidak sesuai.",
    "success": false
}
```

## 2. Testing Post Data Booking Tiket
### Positive Test: Customer 1 Booking Tiket
- **Customer ID:** 6
- **Request:** `POST /api/v1/bookings` (1 tiket kategori Gold)
- **Status Code:** 201
- **Response:**
```json
{
    "data": {
        "id": 6,
        "booking_code": "TIX-1782482141-300",
        "customer_id": 6,
        "customer": {
            "id": 0,
            "user_id": 0,
            "name": "",
            "email": "",
            "created_at": "0001-01-01T00:00:00Z",
            "updated_at": "0001-01-01T00:00:00Z",
            "deleted_at": null
        },
        "total_amount": 1000000,
        "booking_date": "2026-06-26T20:55:41.609576+07:00",
        "details": [
            {
                "id": 5,
                "booking_id": 6,
                "ticket_category_id": 1,
                "ticket_category": {
                    "id": 0,
                    "concert_id": 0,
                    "name": "",
                    "price": 0,
                    "total_quota": 0,
                    "available_quota": 0,
                    "created_at": "0001-01-01T00:00:00Z",
                    "updated_at": "0001-01-01T00:00:00Z",
                    "deleted_at": "0001-01-01T00:00:00Z"
                },
                "quantity": 1,
                "sub_total": 1000000,
                "created_at": "2026-06-26T20:55:41.6153339+07:00",
                "updated_at": "2026-06-26T20:55:41.6153339+07:00"
            }
        ],
        "created_at": "2026-06-26T20:55:41.6111563+07:00",
        "updated_at": "2026-06-26T20:55:41.6160925+07:00"
    },
    "message": "Pemesanan tiket berhasil dikonfirmasi!",
    "success": true
}
```
### Positive Test: Customer 2 Booking Tiket
- **Customer ID:** 7
- **Request:** `POST /api/v1/bookings` (2 tiket kategori Silver)
- **Status Code:** 201
- **Response:**
```json
{
    "data": {
        "id": 7,
        "booking_code": "TIX-1782482141-900",
        "customer_id": 7,
        "customer": {
            "id": 0,
            "user_id": 0,
            "name": "",
            "email": "",
            "created_at": "0001-01-01T00:00:00Z",
            "updated_at": "0001-01-01T00:00:00Z",
            "deleted_at": null
        },
        "total_amount": 1000000,
        "booking_date": "2026-06-26T20:55:41.627423+07:00",
        "details": [
            {
                "id": 6,
                "booking_id": 7,
                "ticket_category_id": 2,
                "ticket_category": {
                    "id": 0,
                    "concert_id": 0,
                    "name": "",
                    "price": 0,
                    "total_quota": 0,
                    "available_quota": 0,
                    "created_at": "0001-01-01T00:00:00Z",
                    "updated_at": "0001-01-01T00:00:00Z",
                    "deleted_at": "0001-01-01T00:00:00Z"
                },
                "quantity": 2,
                "sub_total": 1000000,
                "created_at": "2026-06-26T20:55:41.6336503+07:00",
                "updated_at": "2026-06-26T20:55:41.6336503+07:00"
            }
        ],
        "created_at": "2026-06-26T20:55:41.6290729+07:00",
        "updated_at": "2026-06-26T20:55:41.6358474+07:00"
    },
    "message": "Pemesanan tiket berhasil dikonfirmasi!",
    "success": true
}
```
### Negative Test: Booking Melebihi Kuota Tersedia
- **Request:** `POST /api/v1/bookings` (1000 tiket kategori Gold)
- **Status Code:** 400
- **Response:**
```json
{
    "error": "kuota tiket 'Gold' tidak mencukupi (Tersisa: 45, Permintaan: 1000)",
    "message": "Reservasi tiket gagal diproses",
    "success": false
}
```

## 3. Testing Proteksi IDOR
### IDOR Test: Customer 2 Mengakses Booking Customer 1
- **Target Booking ID:** 6 (Milik Customer 1)
- **Request:** `GET /api/v1/bookings/6` dengan Token Customer 2
- **Status Code:** 404
- **Response:**
```json
{
    "message": "Faktur pemesanan tiket tidak ditemukan",
    "success": false
}
```
