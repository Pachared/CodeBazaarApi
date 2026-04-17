# CodeBazaar API

Backend สำหรับโปรเจกต์ CodeBazaar เขียนด้วย `Go + Gin + GORM + PostgreSQL` และวาง contract ให้สัมพันธ์กับไฟล์ frontend ที่ใช้ `auth`, `products`, `checkout`, `seller`, `downloads`, `profile` และ `cookie consent`

## Tech Stack

- Go
- Gin
- GORM
- PostgreSQL
- Docker Compose สำหรับรันฐานข้อมูล local

## โครงสร้างโฟลเดอร์

```text
cmd/api
  main.go                 # entrypoint ของ API

internal/config
  config.go               # โหลด env และ config ทั้งระบบ

internal/database
  database.go             # เปิด connection และ auto-migrate

internal/models
  *.go                    # GORM models และ relation

internal/contracts
  api.go                  # request/response schema ที่ map กับ frontend

internal/repositories
  *.go                    # data access layer

internal/services
  *.go                    # business logic และ mapping response

internal/handlers
  *.go                    # Gin handlers

internal/middleware
  *.go                    # CORS และ mock current-user resolver

internal/routes
  router.go               # รวม route ทั้งหมด

internal/seed
  seed.go                 # seed data ให้ตรงกับ mock/flow ฝั่ง frontend

internal/httpx
  httpx.go                # helper สำหรับ error response
```

## Database Models

หลัก ๆ จะมีตารางเหล่านี้:

- `users`
  รองรับทั้ง buyer และ seller พร้อม field โปรไฟล์ที่ frontend ใช้ใน `AuthProvider`, `ProfilePage`, `CheckoutPage`
- `products`
  เก็บรายการสินค้าทั้งหมดที่ใช้ใน `HomePage`, `ProductsCatalogPage`, `ProductDetailPage`, `SellerStorePage`
- `orders`
  เก็บข้อมูลการ checkout
- `order_items`
  แยกสินค้าต่อรายการสั่งซื้อ เพื่อให้ดึง `seller/orders` ได้ง่าย
- `download_items`
  ใช้เป็นคลังดาวน์โหลดจริงแทน local storage ใน `DownloadsProvider`
- `cookie_consents`
  ใช้เก็บ cookie consent ได้ทั้งแบบผูก user หรือ session

## API ที่เตรียมไว้

รองรับทั้ง route ตรงและ route prefix `/api/v1`

- `POST /auth/google/start`
- `GET /products`
- `GET /products/featured`
- `GET /products/:productID`
- `GET /sellers`
- `GET /sellers/:sellerSlug`
- `GET /sellers/:sellerSlug/products`
- `POST /checkout/orders`
- `POST /seller/onboarding/google`
- `POST /seller/listings`
- `GET /seller/orders`
- `GET /me/profile`
- `PUT /me/profile`
- `GET /me/downloads`
- `POST /me/downloads/:libraryItemID/download`
- `GET /cookie-consent`
- `PUT /cookie-consent`
- `GET /me/cookie-consent`
- `PUT /me/cookie-consent`
- `GET /health`

## Auth/Context ชั่วคราวสำหรับ frontend

ตอนนี้ frontend ที่ให้มายังไม่ได้ส่ง token จริงมา จึงมี fallback สำหรับ local integration:

- ถ้าเรียก `auth/google/start` จะได้ buyer session ทดลองกลับไป
- ถ้าเรียก `seller/onboarding/google` จะได้ seller session ทดลองกลับไป
- ถ้าต้องการให้ route อย่าง `me/profile`, `me/downloads`, `seller/orders` ผูกกับ user เฉพาะราย ให้ส่ง header:
  - `X-User-ID`
  - หรือ `X-User-Email`
- ถ้าจะใช้ cookie consent แบบ anonymous ให้ส่ง `X-Session-Key`

## วิธีรัน

1. สตาร์ต PostgreSQL

```bash
docker compose up -d
```

2. สร้างไฟล์ env

```bash
cp .env.example .env
```

3. รัน API

```bash
go run ./cmd/api
```

เซิร์ฟเวอร์จะรันที่ `http://localhost:8080`

## Frontend Base URL ที่แนะนำ

ถ้าฝั่ง React เรียก path แบบ `/products/featured`, `/checkout/orders`, `/seller/orders`
ให้ตั้ง `apiBaseUrl` เป็น:

```text
http://localhost:8080
```

หรือถ้าต้องการใช้ namespace:

```text
http://localhost:8080/api/v1
```

## หมายเหตุ

- seed data ถูกออกแบบให้ใกล้กับ `mockProducts`, `mockSellerOrders` และ flow ใน frontend มากที่สุด
- `notification` ยังเป็นฝั่ง client UI queue ตาม `NotificationProvider` จึงไม่ได้แยกเป็นตารางเฉพาะ
- ถ้าจะต่อ OAuth จริงหรือ JWT ต่อจากนี้ สามารถเริ่มจาก middleware `CurrentUser` ได้ทันที
