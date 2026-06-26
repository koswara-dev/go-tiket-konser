import urllib.request
import urllib.error
import json
import random

api_key = "juara-coding-super-secret"
admin_email = "adminkonser@gmail.com"
admin_password = "Indonesia"

cust_a_email = f"customer_a_{random.randint(100000, 999999)}@gmail.com"
cust_b_email = f"customer_b_{random.randint(100000, 999999)}@gmail.com"
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
            return response.getcode(), response.read().decode('utf-8')
    except urllib.error.HTTPError as e:
        return e.code, e.read().decode('utf-8')
    except Exception as e:
        return 500, str(e)

# --- 1. SETUP: REGISTER AND LOGIN ---
print("Logging in Admin...")
_, res = make_request("http://localhost:8080/api/v1/login", "POST", {"email": admin_email, "password": admin_password})
admin_token = json.loads(res)["data"]["token"]

print("Registering Customer A...")
_, res = make_request("http://localhost:8080/api/v1/register", "POST", {"email": cust_a_email, "password": cust_password, "full_name": "Customer A"})
cust_a_res = json.loads(res)["data"]
cust_a_user_id = cust_a_res["id"]
cust_a_cust_id = cust_a_res["customer_id"]

print("Logging in Customer A...")
_, res = make_request("http://localhost:8080/api/v1/login", "POST", {"email": cust_a_email, "password": cust_password})
cust_a_token = json.loads(res)["data"]["token"]

print("Registering Customer B...")
_, res = make_request("http://localhost:8080/api/v1/register", "POST", {"email": cust_b_email, "password": cust_password, "full_name": "Customer B"})
cust_b_res = json.loads(res)["data"]
cust_b_user_id = cust_b_res["id"]
cust_b_cust_id = cust_b_res["customer_id"]

print("Logging in Customer B...")
_, res = make_request("http://localhost:8080/api/v1/login", "POST", {"email": cust_b_email, "password": cust_password})
cust_b_token = json.loads(res)["data"]["token"]


# --- 2. USERS API TESTING ---
print("\n=== TESTING USERS API ===")

# List all users as Customer -> should fail (403)
code, res = make_request("http://localhost:8080/api/v1/users", "GET", headers={"Authorization": f"Bearer {cust_a_token}"})
print(f"List users as Customer: Status={code}, Response={res.strip()}")
assert code == 403, "Expected 403 when customer lists users"

# List all users as Admin -> should succeed (200)
code, res = make_request("http://localhost:8080/api/v1/users", "GET", headers={"Authorization": f"Bearer {admin_token}"})
print(f"List users as Admin: Status={code}")
assert code == 200, "Expected 200 when admin lists users"

# Get own user as Customer A -> should succeed (200)
code, res = make_request(f"http://localhost:8080/api/v1/users/{cust_a_user_id}", "GET", headers={"Authorization": f"Bearer {cust_a_token}"})
print(f"Get own user profile: Status={code}, Response={res.strip()}")
assert code == 200, "Expected 200 when getting own user profile"

# Get user B as Customer A -> should fail (403 - IDOR protection)
code, res = make_request(f"http://localhost:8080/api/v1/users/{cust_b_user_id}", "GET", headers={"Authorization": f"Bearer {cust_a_token}"})
print(f"Get other user profile (IDOR): Status={code}, Response={res.strip()}")
assert code == 403, "Expected 403 when trying to access other user's profile"

# Get user B as Admin -> should succeed (200)
code, res = make_request(f"http://localhost:8080/api/v1/users/{cust_b_user_id}", "GET", headers={"Authorization": f"Bearer {admin_token}"})
print(f"Get other user profile as Admin: Status={code}")
assert code == 200, "Expected 200 when admin gets another user profile"

# Update own profile as Customer A -> should succeed (200)
update_body = {"full_name": "Customer A Edited", "email": cust_a_email, "role": "admin"} # try to elevate role
code, res = make_request(f"http://localhost:8080/api/v1/users/{cust_a_user_id}", "PUT", body=update_body, headers={"Authorization": f"Bearer {cust_a_token}"})
print(f"Update own profile: Status={code}, Response={res.strip()}")
assert code == 200, "Expected 200 when updating own profile"
# Check if role elevation was ignored/rejected
assert json.loads(res)["role"] == "customer", "Customer should not be able to elevate their role"

# Update B profile as Customer A -> should fail (403 - IDOR protection)
code, res = make_request(f"http://localhost:8080/api/v1/users/{cust_b_user_id}", "PUT", body=update_body, headers={"Authorization": f"Bearer {cust_a_token}"})
print(f"Update other profile (IDOR): Status={code}, Response={res.strip()}")
assert code == 403, "Expected 403 when trying to update other user profile"

# Delete user as Customer -> should fail (403)
code, res = make_request(f"http://localhost:8080/api/v1/users/{cust_b_user_id}", "DELETE", headers={"Authorization": f"Bearer {cust_a_token}"})
print(f"Delete user as Customer: Status={code}, Response={res.strip()}")
assert code == 403, "Expected 403 when customer tries to delete a user"


# --- 3. CUSTOMERS API TESTING ---
print("\n=== TESTING CUSTOMERS API ===")

# List all customers as Customer -> should fail (403)
code, res = make_request("http://localhost:8080/api/v1/customers", "GET", headers={"Authorization": f"Bearer {cust_a_token}"})
print(f"List customers as Customer: Status={code}, Response={res.strip()}")
assert code == 403, "Expected 403 when customer lists customers"

# List all customers as Admin -> should succeed (200)
code, res = make_request("http://localhost:8080/api/v1/customers", "GET", headers={"Authorization": f"Bearer {admin_token}"})
print(f"List customers as Admin: Status={code}")
assert code == 200, "Expected 200 when admin lists customers"

# Get own customer profile as Customer A -> should succeed (200)
code, res = make_request(f"http://localhost:8080/api/v1/customers/{cust_a_cust_id}", "GET", headers={"Authorization": f"Bearer {cust_a_token}"})
print(f"Get own customer profile: Status={code}, Response={res.strip()}")
assert code == 200, "Expected 200 when getting own customer profile"

# Get customer B profile as Customer A -> should fail (403 - IDOR protection)
code, res = make_request(f"http://localhost:8080/api/v1/customers/{cust_b_cust_id}", "GET", headers={"Authorization": f"Bearer {cust_a_token}"})
print(f"Get other customer profile (IDOR): Status={code}, Response={res.strip()}")
assert code == 403, "Expected 403 when trying to access other customer's profile"

# Update own customer profile as Customer A -> should succeed (200)
update_cust_body = {"name": "Customer A Customer Edited", "email": cust_a_email}
code, res = make_request(f"http://localhost:8080/api/v1/customers/{cust_a_cust_id}", "PUT", body=update_cust_body, headers={"Authorization": f"Bearer {cust_a_token}"})
print(f"Update own customer profile: Status={code}, Response={res.strip()}")
assert code == 200, "Expected 200 when updating own customer profile"

# Update customer B profile as Customer A -> should fail (403 - IDOR protection)
code, res = make_request(f"http://localhost:8080/api/v1/customers/{cust_b_cust_id}", "PUT", body=update_cust_body, headers={"Authorization": f"Bearer {cust_a_token}"})
print(f"Update other customer profile (IDOR): Status={code}, Response={res.strip()}")
assert code == 403, "Expected 403 when trying to update other customer profile"

# Delete customer as Customer -> should fail (403)
code, res = make_request(f"http://localhost:8080/api/v1/customers/{cust_b_cust_id}", "DELETE", headers={"Authorization": f"Bearer {cust_a_token}"})
print(f"Delete customer as Customer: Status={code}, Response={res.strip()}")
assert code == 403, "Expected 403 when customer tries to delete a customer"


# --- 4. CLEANUP (Admin deleting Customer B) ---
print("\n=== CLEANUP: DELETING USER AND CUSTOMER B AS ADMIN ===")
code, res = make_request(f"http://localhost:8080/api/v1/users/{cust_b_user_id}", "DELETE", headers={"Authorization": f"Bearer {admin_token}"})
print(f"Delete User B as Admin: Status={code}, Response={res.strip()}")
assert code == 200, "Expected 200 when admin deletes user"

code, res = make_request(f"http://localhost:8080/api/v1/customers/{cust_b_cust_id}", "DELETE", headers={"Authorization": f"Bearer {admin_token}"})
print(f"Delete Customer B as Admin: Status={code}, Response={res.strip()}")
assert code == 200, "Expected 200 when admin deletes customer"

print("\nALL CRUD AND IDOR TESTS PASSED SUCCESSFULLY!")
