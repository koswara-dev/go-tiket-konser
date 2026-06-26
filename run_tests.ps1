# Powershell Script to run all tests and generate result-testing.md

$apiKey = "juara-coding-super-secret"
$adminEmail = "adminkonser@gmail.com"
$adminPassword = "Indonesia"

$cust1Email = "customer1_$(Get-Random)@gmail.com"
$cust2Email = "customer2_$(Get-Random)@gmail.com"
$custPassword = "password123"

$results = @"
# Result Testing

Dokumen ini berisi hasil pengujian yang dilakukan untuk memvalidasi peran (RBAC), pembuatan booking tiket untuk 2 customer (termasuk kasus negatif), dan proteksi IDOR.

"@

# Helper to run a web request and return format
function Test-Request {
    param(
        [string]$Uri,
        [string]$Method,
        [string]$Body,
        [hashtable]$Headers
    )
    
    $reqHeaders = @{
        "x-api-key" = $apiKey
    }
    if ($Headers) {
        foreach ($key in $Headers.Keys) {
            $reqHeaders[$key] = $Headers[$key]
        }
    }

    try {
        if ($Body) {
            $res = Invoke-WebRequest -Uri $Uri -Method $Method -Headers $reqHeaders -Body $Body -ContentType "application/json" -UseBasicParsing
        } else {
            $res = Invoke-WebRequest -Uri $Uri -Method $Method -Headers $reqHeaders -UseBasicParsing
        }
        return [PSCustomObject]@{
            StatusCode = $res.StatusCode
            Content = $res.Content
        }
    } catch [System.Net.WebException] {
        $exceptionResponse = $_.Exception.Response
        $statusCode = [int]$exceptionResponse.StatusCode
        $stream = $exceptionResponse.GetResponseStream()
        $reader = New-Object System.IO.StreamReader($stream)
        $content = $reader.ReadToEnd()
        $reader.Close()
        $stream.Close()
        return [PSCustomObject]@{
            StatusCode = $statusCode
            Content = $content
        }
    } catch {
        return [PSCustomObject]@{
            StatusCode = 500
            Content = $_.Exception.Message
        }
    }
}

# --- STAGE 0: Get Tokens ---
Write-Host "Getting Admin Token..."
$adminLoginBody = @{ email = $adminEmail; password = $adminPassword } | ConvertTo-Json
$res = Test-Request -Uri "http://localhost:8080/api/v1/login" -Method Post -Body $adminLoginBody
$adminToken = (ConvertFrom-Json $res.Content).data.token

Write-Host "Registering Customer 1..."
$reg1Body = @{ email = $cust1Email; password = $custPassword; full_name = "Customer One" } | ConvertTo-Json
$res = Test-Request -Uri "http://localhost:8080/api/v1/register" -Method Post -Body $reg1Body
$cust1Id = (ConvertFrom-Json $res.Content).data.customer_id

Write-Host "Logging in Customer 1..."
$login1Body = @{ email = $cust1Email; password = $custPassword } | ConvertTo-Json
$res = Test-Request -Uri "http://localhost:8080/api/v1/login" -Method Post -Body $login1Body
$cust1Token = (ConvertFrom-Json $res.Content).data.token

Write-Host "Registering Customer 2..."
$reg2Body = @{ email = $cust2Email; password = $custPassword; full_name = "Customer Two" } | ConvertTo-Json
$res = Test-Request -Uri "http://localhost:8080/api/v1/register" -Method Post -Body $reg2Body
$cust2Id = (ConvertFrom-Json $res.Content).data.customer_id

Write-Host "Logging in Customer 2..."
$login2Body = @{ email = $cust2Email; password = $custPassword } | ConvertTo-Json
$res = Test-Request -Uri "http://localhost:8080/api/v1/login" -Method Post -Body $login2Body
$cust2Token = (ConvertFrom-Json $res.Content).data.token


# --- 1. Testing Role (RBAC) ---
$results += "`n## 1. Testing Role (Otorisasi RBAC)`n"

# Positive Test
Write-Host "Testing Role: Admin accessing Concert List..."
$resAdmin = Test-Request -Uri "http://localhost:8080/api/v1/concerts" -Method Get -Headers @{ "Authorization" = "Bearer $adminToken" }
$results += "### Positive Test: Admin Mengakses Daftar Konser
- **Request:** `GET /api/v1/concerts` dengan Token Admin
- **Status Code:** $($resAdmin.StatusCode)
- **Response:**
```json
$($resAdmin.Content | ConvertFrom-Json | ConvertTo-Json -Depth 5)
```
"

# Negative Test
Write-Host "Testing Role: Customer accessing Concert List..."
$resCust = Test-Request -Uri "http://localhost:8080/api/v1/concerts" -Method Get -Headers @{ "Authorization" = "Bearer $cust1Token" }
$results += "### Negative Test: Customer Mengakses Daftar Konser
- **Request:** `GET /api/v1/concerts` dengan Token Customer
- **Status Code:** $($resCust.StatusCode)
- **Response:**
```json
$($resCust.Content | ConvertFrom-Json | ConvertTo-Json -Depth 5)
```
"


# --- 2. Testing Post Data Booking Tiket (2 Customer) ---
$results += "`n## 2. Testing Post Data Booking Tiket`n"

# Customer 1 Positive Test
Write-Host "Booking ticket for Customer 1..."
$booking1Body = @{
    customer_id = $cust1Id
    booking_details = @(
        @{
            ticket_category_id = 1
            quantity = 1
        }
    )
} | ConvertTo-Json
$resBooking1 = Test-Request -Uri "http://localhost:8080/api/v1/bookings" -Method Post -Body $booking1Body -Headers @{ "Authorization" = "Bearer $cust1Token" }
$booking1Id = (ConvertFrom-Json $resBooking1.Content).data.id

$results += "### Positive Test: Customer 1 Booking Tiket
- **Customer ID:** $cust1Id
- **Request:** `POST /api/v1/bookings` (1 tiket kategori Gold)
- **Status Code:** $($resBooking1.StatusCode)
- **Response:**
```json
$($resBooking1.Content | ConvertFrom-Json | ConvertTo-Json -Depth 5)
```
"

# Customer 2 Positive Test
Write-Host "Booking ticket for Customer 2..."
$booking2Body = @{
    customer_id = $cust2Id
    booking_details = @(
        @{
            ticket_category_id = 2
            quantity = 2
        }
    )
} | ConvertTo-Json
$resBooking2 = Test-Request -Uri "http://localhost:8080/api/v1/bookings" -Method Post -Body $booking2Body -Headers @{ "Authorization" = "Bearer $cust2Token" }

$results += "### Positive Test: Customer 2 Booking Tiket
- **Customer ID:** $cust2Id
- **Request:** `POST /api/v1/bookings` (2 tiket kategori Silver)
- **Status Code:** $($resBooking2.StatusCode)
- **Response:**
```json
$($resBooking2.Content | ConvertFrom-Json | ConvertTo-Json -Depth 5)
```
"

# Negative Test: Quota Exceeded
Write-Host "Booking ticket negative test: Quota Exceeded..."
$bookingNegBody = @{
    customer_id = $cust1Id
    booking_details = @(
        @{
            ticket_category_id = 1
            quantity = 1000
        }
    )
} | ConvertTo-Json
$resBookingNeg = Test-Request -Uri "http://localhost:8080/api/v1/bookings" -Method Post -Body $bookingNegBody -Headers @{ "Authorization" = "Bearer $cust1Token" }

$results += "### Negative Test: Booking Melebihi Kuota Tersedia
- **Request:** `POST /api/v1/bookings` (1000 tiket kategori Gold)
- **Status Code:** $($resBookingNeg.StatusCode)
- **Response:**
```json
$($resBookingNeg.Content | ConvertFrom-Json | ConvertTo-Json -Depth 5)
```
"


# --- 3. Testing Proteksi IDOR ---
$results += "`n## 3. Testing Proteksi IDOR`n"

# Negative IDOR Test
Write-Host "Testing IDOR: Customer 2 reading Customer 1's booking..."
$resIDOR = Test-Request -Uri "http://localhost:8080/api/v1/bookings/$booking1Id" -Method Get -Headers @{ "Authorization" = "Bearer $cust2Token" }

$results += "### IDOR Test: Customer 2 Mengakses Booking Customer 1
- **Target Booking ID:** $booking1Id (Milik Customer 1)
- **Request:** `GET /api/v1/bookings/$booking1Id` dengan Token Customer 2
- **Status Code:** $($resIDOR.StatusCode)
- **Response:**
```json
$($resIDOR.Content | ConvertFrom-Json | ConvertTo-Json -Depth 5)
```
"

# Save results
$results | Out-File -FilePath "c:\Projects\JuaraCoding\golang-batch-1\go-tiket-konser\result-testing.md" -Encoding utf8
Write-Host "Saved results to result-testing.md"
