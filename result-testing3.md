# Result Testing 3

Dokumen ini berisi hasil pengujian yang dilakukan untuk memvalidasi fitur SSE (Server-Sent Events) Notifications dan WebSocket Chat, termasuk pengujian otorisasi dan proteksi IDOR.

## 1. Testing SSE Notifications (Server-Sent Events)

### Positive Test 1: Establish SSE Connection
- **Endpoint:** `GET /api/v1/notifications/stream`
- **Connected:** False
- **Initial Events Received:**
```json
[]
```
### Positive Test 2: Broadcast Notification on Concert Creation
- **Action:** Create Concert `Test Concert SSE 3460`
- **Concert Creation Status:** 201
- **SSE Events Received after Concert Creation:**
```json
[
    {
        "event": "info",
        "data": "SSE connection established"
    },
    {
        "event": "message",
        "data": "{\"id\":\"6a4d0c69642bfee8680d953d\",\"user_id\":\"\",\"title\":\"Konser Baru!\",\"message\":\"Konser baru 'Test Concert SSE 3460' telah ditambahkan di Gelora Bung Karno pada 2027-01-01!\",\"created_at\":\"2026-07-07T21:25:45.3768954+07:00\"}"
    }
]
```
### Positive Test 3: Targeted Notification on Booking Success
- **Action:** Book ticket category Gold
- **Booking Creation Status:** 201
- **Booking Response:**
```json
{
    "data": {
        "id": "f58b47a1-6fd3-4d2c-a4be-d483c9bda3d9",
        "created_at": "2026-07-07T21:25:47.3877199+07:00",
        "updated_at": "2026-07-07T21:25:47.3924933+07:00",
        "deleted_at": null,
        "created_by": "48f308a4-50bd-450d-a971-720214db95e3",
        "booking_code": "TIX-1783434347-900",
        "customer_id": "1e197317-beb0-4425-9249-af314e41966c",
        "customer": {
            "id": "00000000-0000-0000-0000-000000000000",
            "created_at": "0001-01-01T00:00:00Z",
            "updated_at": "0001-01-01T00:00:00Z",
            "deleted_at": null,
            "user_id": "00000000-0000-0000-0000-000000000000",
            "name": "",
            "email": ""
        },
        "total_amount": 1000000,
        "booking_date": "2026-07-07T21:25:46.961912+07:00",
        "details": [
            {
                "id": "374d48cf-5767-425c-a205-cf26bd140d1f",
                "created_at": "2026-07-07T21:25:47.3915821+07:00",
                "updated_at": "2026-07-07T21:25:47.3915821+07:00",
                "deleted_at": null,
                "created_by": "48f308a4-50bd-450d-a971-720214db95e3",
                "booking_id": "f58b47a1-6fd3-4d2c-a4be-d483c9bda3d9",
                "ticket_category_id": "33333333-3333-3333-3333-333333333331",
                "ticket_category": {
                    "id": "00000000-0000-0000-0000-000000000000",
                    "created_at": "0001-01-01T00:00:00Z",
                    "updated_at": "0001-01-01T00:00:00Z",
                    "deleted_at": null,
                    "concert_id": "00000000-0000-0000-0000-000000000000",
                    "name": "",
                    "price": 0,
                    "total_quota": 0,
                    "available_quota": 0
                },
                "quantity": 1,
                "sub_total": 1000000
            }
        ]
    },
    "message": "Pemesanan tiket berhasil dikonfirmasi!",
    "success": true
}
```
- **SSE Events Received after Booking:**
```json
[
    {
        "event": "message",
        "data": "{\"id\":\"6a4d0c6b642bfee8680d953f\",\"user_id\":\"48f308a4-50bd-450d-a971-720214db95e3\",\"title\":\"Booking Berhasil\",\"message\":\"Booking dengan kode TIX-1783434347-900 berhasil dibuat. Total pembayaran: Rp 1000000.00\",\"created_at\":\"2026-07-07T21:25:47.3980548+07:00\"}"
    }
]
```
### Negative Test: Unauthorized SSE Connection
- **Endpoint:** `GET /api/v1/notifications/stream` (Tanpa Token)
- **Status Code:** 401
- **Response:**
```json
{
    "error": "Authorization header is required"
}
```

## 2. Testing WebSocket Chat (Real-time Messaging)

### Positive Test 1: WebSocket Handshake
- **Customer WS Endpoint:** `GET /api/v1/chat/ws?token=<customer_token>`
- **Admin WS Endpoint:** `GET /api/v1/chat/ws?token=<admin_token>&room_id=48f308a4-50bd-450d-a971-720214db95e3`
- **Upgrade Connection Status:** Success

### Positive Test 2: Bidirectional Chatting
- **Pesan Dikirim Customer:** `"Halo Admin, saya butuh bantuan."`
- **Pesan Diterima Admin:**
```json
{
    "id": "6a4d0c6d642bfee8680d9541",
    "room_id": "48f308a4-50bd-450d-a971-720214db95e3",
    "sender_id": "48f308a4-50bd-450d-a971-720214db95e3",
    "sender_name": "Customer SSE One",
    "role": "customer",
    "message": "Halo Admin, saya butuh bantuan.",
    "timestamp": "2026-07-07T21:25:49.4183692+07:00"
}
```
- **Pesan Balasan Admin:** `"Halo Customer, silakan sebutkan keluhan Anda."`
- **Pesan Balasan Diterima Customer:**
```json
{
    "id": "6a4d0c6d642bfee8680d9542",
    "room_id": "48f308a4-50bd-450d-a971-720214db95e3",
    "sender_id": "99999999-9999-9999-9999-999999999999",
    "sender_name": "Admin Konser",
    "role": "admin",
    "message": "Halo Customer, silakan sebutkan keluhan Anda.",
    "timestamp": "2026-07-07T21:25:49.4445057+07:00"
}
```
### Positive Test 3: Get Room Messages & Rooms List History
- **Endpoint:** `GET /api/v1/chat/rooms` (Admin Token)
- **Status Code:** 200
- **Response:**
```json
{
    "success": true,
    "message": "Daftar room chat berhasil diambil",
    "data": [
        {
            "room_id": "48f308a4-50bd-450d-a971-720214db95e3",
            "customer_name": "Customer SSE One",
            "last_message": {
                "id": "6a4d0c6d642bfee8680d9542",
                "room_id": "48f308a4-50bd-450d-a971-720214db95e3",
                "sender_id": "99999999-9999-9999-9999-999999999999",
                "sender_name": "Admin Konser",
                "role": "admin",
                "message": "Halo Customer, silakan sebutkan keluhan Anda.",
                "timestamp": "2026-07-07T14:25:49.444Z"
            }
        }
    ]
}
```

- **Endpoint:** `GET /api/v1/chat/rooms/48f308a4-50bd-450d-a971-720214db95e3/messages` (Customer 1 Token)
- **Status Code:** 200
- **Response:**
```json
{
    "success": true,
    "message": "Riwayat pesan berhasil diambil",
    "data": [
        {
            "id": "6a4d0c6d642bfee8680d9541",
            "room_id": "48f308a4-50bd-450d-a971-720214db95e3",
            "sender_id": "48f308a4-50bd-450d-a971-720214db95e3",
            "sender_name": "Customer SSE One",
            "role": "customer",
            "message": "Halo Admin, saya butuh bantuan.",
            "timestamp": "2026-07-07T14:25:49.418Z"
        },
        {
            "id": "6a4d0c6d642bfee8680d9542",
            "room_id": "48f308a4-50bd-450d-a971-720214db95e3",
            "sender_id": "99999999-9999-9999-9999-999999999999",
            "sender_name": "Admin Konser",
            "role": "admin",
            "message": "Halo Customer, silakan sebutkan keluhan Anda.",
            "timestamp": "2026-07-07T14:25:49.444Z"
        }
    ]
}
```
### Negative Test 1: Unauthenticated WebSocket Connection
- **Endpoint:** `GET /api/v1/chat/ws` (Tanpa Token)
- **Status Connection:** Failed (Expected): Handshake status 401 Unauthorized -+-+- {'content-type': 'application/json; charset=utf-8', 'date': 'Tue, 07 Jul 2026 14:25:49 GMT', 'content-length': '43'} -+-+- b'{"error":"Authorization token is required"}'
### Negative Test 2: IDOR on Room Messages Access
- **Endpoint:** `GET /api/v1/chat/rooms/48f308a4-50bd-450d-a971-720214db95e3/messages` (Customer 2 Token)
- **Status Code:** 403
- **Response:**
```json
{
    "success": false,
    "message": "Akses ditolak: Anda tidak dapat mengakses percakapan ini"
}
```
