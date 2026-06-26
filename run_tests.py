import urllib.request
import urllib.error
import json
import random

api_key = "juara-coding-super-secret"
admin_email = "adminkonser@gmail.com"
admin_password = "Indonesia"

cust1_email = f"customer1_{random.randint(100000, 999999)}@gmail.com"
cust2_email = f"customer2_{random.randint(100000, 999999)}@gmail.com"
cust_password = "password123"

def make_request(url, method="GET", body=None, headers=None):
    req_headers = {
        "x-api-key": api_key,
        "Content-Type": "application/json"
    }
    if headers:
        req_headers.update(headers)
        
    data = None
    if body:
        data = json.dumps(body).encode('utf-8')
        
    req = urllib.request.Request(url, data=data, headers=req_headers, method=method)
    
    try:
        with urllib.request.urlopen(req) as response:
            status_code = response.getcode()
            content = response.read().decode('utf-8')
            return status_code, content
    except urllib.error.HTTPError as e:
        status_code = e.code
        content = e.read().decode('utf-8')
        return status_code, content
    except Exception as e:
        return 500, str(e)

# --- 0. Get Tokens ---
print("Getting Admin Token...")
status, res = make_request("http://localhost:8080/api/v1/login", "POST", {"email": admin_email, "password": admin_password})
admin_token = json.loads(res)["data"]["token"]

print("Registering Customer 1...")
status, res = make_request("http://localhost:8080/api/v1/register", "POST", {"email": cust1_email, "password": cust_password, "full_name": "Customer One"})
cust1_id = json.loads(res)["data"]["customer_id"]

print("Logging in Customer 1...")
status, res = make_request("http://localhost:8080/api/v1/login", "POST", {"email": cust1_email, "password": cust_password})
cust1_token = json.loads(res)["data"]["token"]

print("Registering Customer 2...")
status, res = make_request("http://localhost:8080/api/v1/register", "POST", {"email": cust2_email, "password": cust_password, "full_name": "Customer Two"})
cust2_id = json.loads(res)["data"]["customer_id"]

print("Logging in Customer 2...")
status, res = make_request("http://localhost:8080/api/v1/login", "POST", {"email": cust2_email, "password": cust_password})
cust2_token = json.loads(res)["data"]["token"]


# --- Collect test results ---
results = """# Result Testing

Dokumen ini berisi hasil pengujian yang dilakukan untuk memvalidasi peran (RBAC), pembuatan booking tiket untuk 2 customer (termasuk kasus negatif), dan proteksi IDOR.
"""

# --- 1. Testing Role (RBAC) ---
results += "\n## 1. Testing Role (Otorisasi RBAC)\n"

# Positive Test
print("Testing Role: Admin accessing Concert List...")
status, content = make_request("http://localhost:8080/api/v1/concerts", "GET", headers={"Authorization": f"Bearer {admin_token}"})
parsed_content = json.dumps(json.loads(content), indent=4)
results += f"""### Positive Test: Admin Mengakses Daftar Konser
- **Request:** `GET /api/v1/concerts` dengan Token Admin
- **Status Code:** {status}
- **Response:**
```json
{parsed_content}
```
"""

# Negative Test
print("Testing Role: Customer accessing Concert List...")
status, content = make_request("http://localhost:8080/api/v1/concerts", "GET", headers={"Authorization": f"Bearer {cust1_token}"})
try:
    parsed_content = json.dumps(json.loads(content), indent=4)
except Exception:
    parsed_content = content
results += f"""### Negative Test: Customer Mengakses Daftar Konser
- **Request:** `GET /api/v1/concerts` dengan Token Customer
- **Status Code:** {status}
- **Response:**
```json
{parsed_content}
```
"""


# --- 2. Testing Post Data Booking Tiket (2 Customer) ---
results += "\n## 2. Testing Post Data Booking Tiket\n"

# Customer 1 Positive Test
print("Booking ticket for Customer 1...")
booking1_body = {
    "customer_id": cust1_id,
    "booking_details": [
        {
            "ticket_category_id": 1,
            "quantity": 1
        }
    ]
}
status, content = make_request("http://localhost:8080/api/v1/bookings", "POST", body=booking1_body, headers={"Authorization": f"Bearer {cust1_token}"})
booking1_res = json.loads(content)
booking1_id = booking1_res["data"]["id"]
parsed_content = json.dumps(booking1_res, indent=4)

results += f"""### Positive Test: Customer 1 Booking Tiket
- **Customer ID:** {cust1_id}
- **Request:** `POST /api/v1/bookings` (1 tiket kategori Gold)
- **Status Code:** {status}
- **Response:**
```json
{parsed_content}
```
"""

# Customer 2 Positive Test
print("Booking ticket for Customer 2...")
booking2_body = {
    "customer_id": cust2_id,
    "booking_details": [
        {
            "ticket_category_id": 2,
            "quantity": 2
        }
    ]
}
status, content = make_request("http://localhost:8080/api/v1/bookings", "POST", body=booking2_body, headers={"Authorization": f"Bearer {cust2_token}"})
parsed_content = json.dumps(json.loads(content), indent=4)

results += f"""### Positive Test: Customer 2 Booking Tiket
- **Customer ID:** {cust2_id}
- **Request:** `POST /api/v1/bookings` (2 tiket kategori Silver)
- **Status Code:** {status}
- **Response:**
```json
{parsed_content}
```
"""

# Negative Test: Quota Exceeded
print("Booking ticket negative test: Quota Exceeded...")
booking_neg_body = {
    "customer_id": cust1_id,
    "booking_details": [
        {
            "ticket_category_id": 1,
            "quantity": 1000
        }
    ]
}
status, content = make_request("http://localhost:8080/api/v1/bookings", "POST", body=booking_neg_body, headers={"Authorization": f"Bearer {cust1_token}"})
try:
    parsed_content = json.dumps(json.loads(content), indent=4)
except Exception:
    parsed_content = content

results += f"""### Negative Test: Booking Melebihi Kuota Tersedia
- **Request:** `POST /api/v1/bookings` (1000 tiket kategori Gold)
- **Status Code:** {status}
- **Response:**
```json
{parsed_content}
```
"""


# --- 3. Testing Proteksi IDOR ---
results += "\n## 3. Testing Proteksi IDOR\n"

# Negative IDOR Test
print("Testing IDOR: Customer 2 reading Customer 1's booking...")
status, content = make_request(f"http://localhost:8080/api/v1/bookings/{booking1_id}", "GET", headers={"Authorization": f"Bearer {cust2_token}"})
try:
    parsed_content = json.dumps(json.loads(content), indent=4)
except Exception:
    parsed_content = content

results += f"""### IDOR Test: Customer 2 Mengakses Booking Customer 1
- **Target Booking ID:** {booking1_id} (Milik Customer 1)
- **Request:** `GET /api/v1/bookings/{booking1_id}` dengan Token Customer 2
- **Status Code:** {status}
- **Response:**
```json
{parsed_content}
```
"""

# Save results
with open("result-testing.md", "w", encoding="utf-8") as f:
    f.write(results)

print("Saved results to result-testing.md")
