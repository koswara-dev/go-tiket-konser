import urllib.request
import urllib.error
import json
import random
import time
import threading
import websocket

api_key = "juara-coding-super-secret"
admin_email = "adminkonser@gmail.com"
admin_password = "Indonesia"
base_url = "http://localhost:8080/api/v1"
ws_base_url = "ws://localhost:8080/api/v1"

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

# --- SSE Client Thread helper ---
class SSEClientThread(threading.Thread):
    def __init__(self, url, token):
        super().__init__()
        self.url = url
        self.token = token
        self.headers = {
            "x-api-key": api_key,
            "Authorization": f"Bearer {token}"
        }
        self.events = []
        self.connected = False
        self.running = True
        self.daemon = True

    def run(self):
        req = urllib.request.Request(self.url, headers=self.headers)
        try:
            with urllib.request.urlopen(req) as resp:
                self.connected = True
                current_event = None
                while self.running:
                    line = resp.readline()
                    if not line:
                        break
                    line = line.decode('utf-8').strip()
                    if line.startswith("event:"):
                        current_event = line[6:].strip()
                    elif line.startswith("data:"):
                        data_content = line[5:].strip()
                        self.events.append({
                            "event": current_event,
                            "data": data_content
                        })
                        current_event = None
        except Exception as e:
            pass

# Collect results in markdown
results = """# Result Testing 3

Dokumen ini berisi hasil pengujian yang dilakukan untuk memvalidasi fitur SSE (Server-Sent Events) Notifications dan WebSocket Chat, termasuk pengujian otorisasi dan proteksi IDOR.

"""

def main():
    global results
    print("Starting Test Suite 3 (SSE & WebSockets)...")
    
    # 1. Login Admin
    print("Logging in Admin...")
    status, content = make_request(f"{base_url}/login", "POST", {"email": admin_email, "password": admin_password})
    admin_token = json.loads(content)["data"]["token"]
    
    # 2. Register & Login Customer 1
    cust1_email = f"test_sse_cust1_{random.randint(1000, 9999)}@gmail.com"
    print(f"Registering Customer 1 ({cust1_email})...")
    status, content = make_request(f"{base_url}/register", "POST", {
        "email": cust1_email,
        "password": "password123",
        "full_name": "Customer SSE One"
    })
    cust1_user_id = json.loads(content)["data"]["id"]
    cust1_customer_id = json.loads(content)["data"]["customer_id"]
    
    status, content = make_request(f"{base_url}/login", "POST", {"email": cust1_email, "password": "password123"})
    cust1_token = json.loads(content)["data"]["token"]

    # 3. Register & Login Customer 2
    cust2_email = f"test_sse_cust2_{random.randint(1000, 9999)}@gmail.com"
    print(f"Registering Customer 2 ({cust2_email})...")
    status, content = make_request(f"{base_url}/register", "POST", {
        "email": cust2_email,
        "password": "password123",
        "full_name": "Customer SSE Two"
    })
    status, content = make_request(f"{base_url}/login", "POST", {"email": cust2_email, "password": "password123"})
    cust2_token = json.loads(content)["data"]["token"]

    results += "## 1. Testing SSE Notifications (Server-Sent Events)\n\n"

    # Positive Test 1: Establish SSE Connection
    print("Testing Positive Test 1: Establish SSE Connection...")
    sse_url = f"{base_url}/notifications/stream"
    sse_thread = SSEClientThread(sse_url, cust1_token)
    sse_thread.start()
    
    # Wait for connection to establish
    time.sleep(2)
    
    results += f"""### Positive Test 1: Establish SSE Connection
- **Endpoint:** `GET /api/v1/notifications/stream`
- **Connected:** {sse_thread.connected}
- **Initial Events Received:**
```json
{json.dumps(sse_thread.events, indent=4)}
```
"""

    # Positive Test 2: Broadcast Notification on Concert Creation
    print("Testing Positive Test 2: Broadcast Notification on Concert Creation...")
    concert_title = f"Test Concert SSE {random.randint(1000, 9999)}"
    concert_body = {
        "title": concert_title,
        "description": "SSE notification test concert description",
        "date": "2027-01-01",
        "venue": "Gelora Bung Karno",
        "status": "upcoming"
    }
    
    # Reset thread events to only capture new ones
    sse_thread.events = []
    
    status, content = make_request(f"{base_url}/concerts", "POST", concert_body, headers={
        "Authorization": f"Bearer {admin_token}"
    })
    
    # Wait for SSE to propagate
    time.sleep(2)
    
    results += f"""### Positive Test 2: Broadcast Notification on Concert Creation
- **Action:** Create Concert `{concert_title}`
- **Concert Creation Status:** {status}
- **SSE Events Received after Concert Creation:**
```json
{json.dumps(sse_thread.events, indent=4)}
```
"""

    # Positive Test 3: Targeted Notification on Booking Success
    print("Testing Positive Test 3: Targeted Notification on Booking Success...")
    booking_body = {
        "customer_id": cust1_customer_id,
        "booking_details": [
            {
                "ticket_category_id": "33333333-3333-3333-3333-333333333331", # Gold ticket seeded
                "quantity": 1
            }
        ]
    }
    
    # Reset thread events
    sse_thread.events = []
    
    status, content = make_request(f"{base_url}/bookings", "POST", booking_body, headers={
        "Authorization": f"Bearer {cust1_token}"
    })
    
    # Wait for SSE to propagate
    time.sleep(2)
    
    results += f"""### Positive Test 3: Targeted Notification on Booking Success
- **Action:** Book ticket category Gold
- **Booking Creation Status:** {status}
- **Booking Response:**
```json
{json.dumps(json.loads(content), indent=4)}
```
- **SSE Events Received after Booking:**
```json
{json.dumps(sse_thread.events, indent=4)}
```
"""

    # Stop SSE thread
    sse_thread.running = False

    # Negative Test: Unauthorized SSE Connection
    print("Testing Negative Test: Unauthorized SSE Connection...")
    status, content = make_request(sse_url, "GET")
    
    results += f"""### Negative Test: Unauthorized SSE Connection
- **Endpoint:** `GET /api/v1/notifications/stream` (Tanpa Token)
- **Status Code:** {status}
- **Response:**
```json
{json.dumps(json.loads(content), indent=4) if status != 500 else content}
```
"""

    results += "\n## 2. Testing WebSocket Chat (Real-time Messaging)\n\n"

    # Positive Test 1 & 2: WS Handshake & Bidirectional Chatting
    print("Testing Positive Test 1 & 2: WS Handshake & Bidirectional Chatting...")
    ws_customer_url = f"{ws_base_url}/chat/ws?token={cust1_token}"
    ws_admin_url = f"{ws_base_url}/chat/ws?token={admin_token}&room_id={cust1_user_id}"
    
    handshake_status = "Failed"
    cust_received_msg = None
    admin_received_msg = None
    
    try:
        # 1. Customer connects
        ws_cust = websocket.create_connection(ws_customer_url, header={"x-api-key": api_key})
        # 2. Admin connects to Customer's room
        ws_admin = websocket.create_connection(ws_admin_url, header={"x-api-key": api_key})
        
        handshake_status = "Success"
        
        # 3. Customer sends message
        ws_cust.send(json.dumps({"message": "Halo Admin, saya butuh bantuan."}))
        
        # Admin receives
        admin_received_raw = ws_admin.recv()
        admin_received_msg = json.loads(admin_received_raw)
        
        # 4. Admin replies
        ws_admin.send(json.dumps({"message": "Halo Customer, silakan sebutkan keluhan Anda."}))
        
        # Customer receives
        cust_received_raw = ws_cust.recv()
        # Since first send will broadcast to customer too, let's read until we get admin's message
        cust_received_msg = json.loads(cust_received_raw)
        if cust_received_msg.get("sender_id") == cust1_user_id:
            # That was customer's own broadcast, read next
            cust_received_raw = ws_cust.recv()
            cust_received_msg = json.loads(cust_received_raw)
            
        ws_cust.close()
        ws_admin.close()
    except Exception as e:
        handshake_status = f"Failed: {str(e)}"
        
    results += f"""### Positive Test 1: WebSocket Handshake
- **Customer WS Endpoint:** `GET /api/v1/chat/ws?token=<customer_token>`
- **Admin WS Endpoint:** `GET /api/v1/chat/ws?token=<admin_token>&room_id={cust1_user_id}`
- **Upgrade Connection Status:** {handshake_status}

### Positive Test 2: Bidirectional Chatting
- **Pesan Dikirim Customer:** `"Halo Admin, saya butuh bantuan."`
- **Pesan Diterima Admin:**
```json
{json.dumps(admin_received_msg, indent=4) if admin_received_msg else "Tidak Diterima"}
```
- **Pesan Balasan Admin:** `"Halo Customer, silakan sebutkan keluhan Anda."`
- **Pesan Balasan Diterima Customer:**
```json
{json.dumps(cust_received_msg, indent=4) if cust_received_msg else "Tidak Diterima"}
```
"""

    # Positive Test 3: Get Room Messages & Rooms List History
    print("Testing Positive Test 3: Get Room Messages & Rooms List History...")
    
    # Get rooms list (Admin only)
    rooms_status, rooms_content = make_request(f"{base_url}/chat/rooms", "GET", headers={
        "Authorization": f"Bearer {admin_token}"
    })
    
    # Get messages history
    messages_status, messages_content = make_request(f"{base_url}/chat/rooms/{cust1_user_id}/messages", "GET", headers={
        "Authorization": f"Bearer {cust1_token}"
    })
    
    results += f"""### Positive Test 3: Get Room Messages & Rooms List History
- **Endpoint:** `GET /api/v1/chat/rooms` (Admin Token)
- **Status Code:** {rooms_status}
- **Response:**
```json
{json.dumps(json.loads(rooms_content), indent=4)}
```

- **Endpoint:** `GET /api/v1/chat/rooms/{cust1_user_id}/messages` (Customer 1 Token)
- **Status Code:** {messages_status}
- **Response:**
```json
{json.dumps(json.loads(messages_content), indent=4)}
```
"""

    # Negative Test 1: Unauthenticated WebSocket Connection
    print("Testing Negative Test 1: Unauthenticated WebSocket Connection...")
    ws_unauth_url = f"{ws_base_url}/chat/ws"
    ws_unauth_status = "Failed"
    try:
        ws_temp = websocket.create_connection(ws_unauth_url, header={"x-api-key": api_key}, timeout=2)
        ws_temp.close()
        ws_unauth_status = "Success (Unexpected)"
    except Exception as e:
        ws_unauth_status = f"Failed (Expected): {str(e)}"
        
    results += f"""### Negative Test 1: Unauthenticated WebSocket Connection
- **Endpoint:** `GET /api/v1/chat/ws` (Tanpa Token)
- **Status Connection:** {ws_unauth_status}
"""

    # Negative Test 2: IDOR on Room Messages Access (Customer 2 trying to read Customer 1's messages)
    print("Testing Negative Test 2: IDOR on Room Messages Access...")
    idor_status, idor_content = make_request(f"{base_url}/chat/rooms/{cust1_user_id}/messages", "GET", headers={
        "Authorization": f"Bearer {cust2_token}"
    })
    
    results += f"""### Negative Test 2: IDOR on Room Messages Access
- **Endpoint:** `GET /api/v1/chat/rooms/{cust1_user_id}/messages` (Customer 2 Token)
- **Status Code:** {idor_status}
- **Response:**
```json
{json.dumps(json.loads(idor_content), indent=4)}
```
"""

    # Save results to result-testing3.md
    with open("result-testing3.md", "w", encoding="utf-8") as f:
        f.write(results)
    print("Saved all test results to result-testing3.md successfully!")

if __name__ == "__main__":
    main()
