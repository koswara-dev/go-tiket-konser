import urllib.request
import urllib.error
import json
import random
import time

api_key = "juara-coding-super-secret"
base_url = "http://localhost:8080/api/v1"

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

# --- Collect test results ---
results = """# Result Testing 2

Dokumen ini berisi hasil pengujian yang dilakukan untuk memvalidasi fitur Pagination & Search pada daftar konser (termasuk 1 kasus negatif), serta pengujian mekanisme Force Logout pada 10 akun Customer.

"""

# ==========================================
# 1. TESTING PAGINATION & SEARCH
# ==========================================
results += "## 1. Testing Pagination & Search (Daftar Konser)\n"

# Positive Test 1: Pagination & Limit
print("Testing Pagination & Limit...")
status, content = make_request(f"{base_url}/concerts?page=1&limit=2", "GET")
try:
    parsed = json.loads(content)
    formatted = json.dumps(parsed, indent=4, ensure_ascii=False)
except Exception:
    formatted = content
results += f"""### Positive Test: Pagination & Limit
- **Request:** `GET /api/v1/concerts?page=1&limit=2`
- **Status Code:** {status}
- **Response:**
```json
{formatted}
```
"""

# Positive Test 2: Search
print("Testing Search...")
status, content = make_request(f"{base_url}/concerts?search=Coldplay", "GET")
try:
    parsed = json.loads(content)
    formatted = json.dumps(parsed, indent=4, ensure_ascii=False)
except Exception:
    formatted = content
results += f"""### Positive Test: Search (Coldplay)
- **Request:** `GET /api/v1/concerts?search=Coldplay`
- **Status Code:** {status}
- **Response:**
```json
{formatted}
```
"""

# Positive Test 3: Sorting
print("Testing Sorting...")
status, content = make_request(f"{base_url}/concerts?sort=date_asc", "GET")
try:
    parsed = json.loads(content)
    formatted = json.dumps(parsed, indent=4, ensure_ascii=False)
except Exception:
    formatted = content
results += f"""### Positive Test: Sorting (date_asc)
- **Request:** `GET /api/v1/concerts?sort=date_asc`
- **Status Code:** {status}
- **Response:**
```json
{formatted}
```
"""

# Negative Test: Invalid Query Parameter
print("Testing Invalid Query Parameter...")
status, content = make_request(f"{base_url}/concerts?page=-1", "GET")
try:
    parsed = json.loads(content)
    formatted = json.dumps(parsed, indent=4, ensure_ascii=False)
except Exception:
    formatted = content
results += f"""### Negative Test: Pagination dengan Page Negatif
- **Request:** `GET /api/v1/concerts?page=-1`
- **Status Code:** {status}
- **Response:**
```json
{formatted}
```
"""

# ==========================================
# 2. TESTING FORCE LOGOUT 10 ACCOUNTS
# ==========================================
results += "\n## 2. Testing Force Logout (10 Akun Customer)\n"
results += "Melakukan pengujian force logout untuk memastikan token JWT yang sudah di-blacklist/logout tidak dapat digunakan kembali.\n\n"

# Loop 10 accounts
for i in range(1, 11):
    rand_suffix = random.randint(100000, 999999)
    email = f"cust_logout_{i}_{rand_suffix}@gmail.com"
    password = "password123"
    fullname = f"Customer Logout Test {i}"
    
    print(f"\n--- Testing Account {i}: {email} ---")
    
    # 1. Register
    reg_status, reg_content = make_request(f"{base_url}/register", "POST", {
        "email": email,
        "password": password,
        "full_name": fullname
    })
    
    # 2. Login
    login_status, login_content = make_request(f"{base_url}/login", "POST", {
        "email": email,
        "password": password
    })
    try:
        token = json.loads(login_content)["data"]["token"]
    except Exception:
        token = ""
        
    # 3. Check profile (Before logout)
    prof_before_status, prof_before_content = make_request(f"{base_url}/profile", "GET", headers={
        "Authorization": f"Bearer {token}"
    })
    
    # 4. Logout
    logout_status, logout_content = make_request(f"{base_url}/logout", "POST", headers={
        "Authorization": f"Bearer {token}"
    })
    
    # 5. Check profile (After logout)
    prof_after_status, prof_after_content = make_request(f"{base_url}/profile", "GET", headers={
        "Authorization": f"Bearer {token}"
    })
    
    # Format and save details
    results += f"""### Akun Customer {i} ({email})
- **Langkah 1: Registrasi Akun**
  - **Endpoint:** `POST /api/v1/register`
  - **Status Code:** {reg_status}
  - **Response:** `{reg_content.strip()}`
- **Langkah 2: Login Akun**
  - **Endpoint:** `POST /api/v1/login`
  - **Status Code:** {login_status}
  - **Response:** `{login_content.strip()}`
- **Langkah 3: Akses Profil Sebelum Logout**
  - **Endpoint:** `GET /api/v1/profile`
  - **Status Code:** {prof_before_status}
  - **Response:** `{prof_before_content.strip()}`
- **Langkah 4: Force Logout**
  - **Endpoint:** `POST /api/v1/logout`
  - **Status Code:** {logout_status}
  - **Response:** `{logout_content.strip()}`
- **Langkah 5: Akses Profil Setelah Logout (Blacklisted)**
  - **Endpoint:** `GET /api/v1/profile`
  - **Status Code:** {prof_after_status}
  - **Response:** `{prof_after_content.strip()}`

"""

# Save results to result-testing2.md
with open("result-testing2.md", "w", encoding="utf-8") as f:
    f.write(results)

print("\nSaved all test results to result-testing2.md successfully!")
