# Queue Ticket System (interview-question-005)

โปรเจกต์สอบสำหรับ **example.com** — ระบบรับบัตรคิว IT 05 (หน้า 05-1 / 05-2 / 05-3)

ออกบัตรคิวตั้งแต่ `A0` ถึง `Z9` แสดงหมายเลขคิวที่ออกล่าสุด และล้างคิวกลับเป็น `00` พร้อมป้องกันการกดรับคิวพร้อมกันไม่ให้เลขซ้ำ

## ความต้องการตามโจทย์

| หน้า | Route | การทำงาน |
|------|-------|----------|
| **IT 05-1** | `/` | กด **รับบัตรคิว** → ไปหน้าแสดงเลขคิว (05-2) · กด **ล้างคิว** → ไปหน้า 05-3 |
| **IT 05-2** | `/queue` | แสดงหมายเลขคิว · กด **กลับไปหน้ารับบัตรคิว** → 05-1 |
| **IT 05-3** | `/reset` | กด **ล้างคิว** → แสดง `00` · กด **กลับไปหน้ารับบัตรคิว** → 05-1 |

**กฎรันคิว (ตัวอย่าง):** คิวปัจจุบัน `A9` → กดรับบัตร → ได้ `B0` · หลัง `Z9` ระบบไม่วนอัตโนมัติ ต้องล้างคิวก่อน

**ฐานข้อมูล:** ออกแบบ `queue_state` (สถานะคิวปัจจุบัน 1 แถว) + `queue_history` (ประวัติบัตรที่ออก) · ใช้ transaction + `SELECT … FOR UPDATE` กันคิวซ้ำเมื่อกดพร้อมกัน

## Tech Stack

- **Backend:** Go 1.24, Fiber, GORM (`example.com/interview-question-005/backend`)
- **Frontend:** Vue 3, Vite, Vue Router, TypeScript
- **Database:** PostgreSQL 16 (Docker Compose สำหรับ local dev)

## สิ่งที่ต้องมีก่อนรัน

- [Docker](https://www.docker.com/) หรือ [OrbStack](https://orbstack.dev/) (สำหรับ PostgreSQL)
- [Go](https://go.dev/) 1.24+
- [Node.js](https://nodejs.org/) 20+ และ npm

> โจทย์ไม่บังคับ Docker/SQLite — ใช้ Postgres เพราะรองรับ row lock และตรงกับ test concurrency

## Quick Start (สำหรับผู้ตรวจ)

### 1. เริ่ม PostgreSQL

จาก root โปรเจกต์:

```bash
docker compose up -d
```

รอจน container healthy (`docker compose ps`)

Migration ใน `backend/migrations/` จะรันตอนสร้าง volume ครั้งแรก  
ถ้าเคยรันแล้ว schema เก่า ให้ reset volume:

```bash
docker compose down -v
docker compose up -d
```

### 2. รัน Backend

```bash
cd backend
cp .env.example .env
go mod download
go run ./cmd/api
```

API: http://localhost:8080

### 3. รัน Frontend

เปิด terminal ใหม่:

```bash
cd frontend
cp .env.example .env
npm install
npm run dev
```

เว็บ: http://localhost:5173

### 4. ทดสอบ flow ตามโจทย์ (มือ)

1. เปิด http://localhost:5173 (IT 05-1) → กด **รับบัตรคิว** → เห็นเลขคิวที่ IT 05-2  
2. กด **กลับไปหน้ารับบัตรคิว**  
3. กด **ล้างคิว** → ไป IT 05-3 → กด **ล้างคิว** → เห็น `00`  
4. กลับ 05-1 → รับบัตรอีกครั้ง → ได้ `A0`

### 5. รัน Automated Tests

Unit tests (ไม่ต้องมี DB):

```bash
cd backend
go test ./internal/service
```

Concurrency test (ต้องมี Postgres ตาม `docker compose`):

```bash
cd backend
TEST_DATABASE_URL="postgres://queue_user:queue_password@localhost:5432/queue_db?sslmode=disable" go test ./tests
```

รัน backend tests ทั้งหมด:

```bash
cd backend
go test ./...
```

## Environment Variables

คัดลอกจาก `.env.example` — **อย่า commit ไฟล์ `.env`**

**Backend** (`backend/.env`):

```env
APP_PORT=8080
DATABASE_URL=postgres://queue_user:queue_password@localhost:5432/queue_db?sslmode=disable
CORS_ALLOWED_ORIGINS=http://localhost:5173,http://127.0.0.1:5173
```

**Frontend** (`frontend/.env`):

```env
VITE_API_BASE_URL=http://localhost:8080
```

## โครงสร้างโปรเจกต์

```text
interview-question-005/
  backend/
    cmd/api/main.go          # HTTP server
    internal/
      config/                # env config
      database/              # GORM connection
      handler/               # REST handlers
      model/                 # entities
      repository/            # DB + FOR UPDATE
      service/               # queue logic A0–Z9
    migrations/001_create_queue_tables.sql
    tests/concurrency_test.go
  frontend/
    src/
      api/queue.ts
      views/                 # IT 05-1, 05-2, 05-3
      router/
  docker-compose.yml         # PostgreSQL only
  README.md
```

## API

| Method | Path | คำอธิบาย |
|--------|------|----------|
| GET | `/health` | health check |
| GET | `/api/queue/current` | คิวปัจจุบัน (`00` หรือ `A0`–`Z9`) |
| POST | `/api/queue/next` | ออกบัตรคิวถัดไป |
| POST | `/api/queue/reset` | ล้างคิวเป็น `00` |

**ตัวอย่าง**

`GET /api/queue/current`

```json
{ "current_queue": "A5" }
```

`POST /api/queue/next`

```json
{ "queue_number": "A6" }
```

เมื่อคิวถึง `Z9` แล้ว:

```json
{ "error": "Queue limit reached" }
```

`POST /api/queue/reset`

```json
{ "queue_number": "00" }
```

## Database Schema

**`queue_state`** — แถวเดียว `id = 1`

| คอลัมน์ | ความหมาย |
|---------|----------|
| `current_queue` | ค่าที่แสดง (`00` หรือ `A0`–`Z9`) |
| `current_letter` | `A`–`Z` หรือ NULL ตอน reset |
| `current_number` | `0`–`9` หรือ NULL ตอน reset |
| `updated_at` | เวลาอัปเดตล่าสุด |

**`queue_history`** — เก็บเฉพาะบัตรที่ออก (ไม่เก็บ `00` จากการล้างคิว)

## กฎรันคิว

- สถานะหลังล้างคิว: `00`
- บัตรแรกหลังล้าง: `A0`
- `A9` → `B0`, `B9` → `C0`, … `Y9` → `Z0`, `Z8` → `Z9`
- หลัง `Z9` ไม่วนกลับ — ต้องล้างคิว

## การป้องกันกดรับคิวพร้อมกัน

`POST /api/queue/next` ทำงานใน transaction เดียว:

```sql
BEGIN;
SELECT id, current_letter, current_number, current_queue
FROM queue_state
WHERE id = 1
FOR UPDATE;
-- คำนวณคิวถัดไป
UPDATE queue_state ...
INSERT INTO queue_history ...
COMMIT;
```

`FOR UPDATE` ล็อกแถว `queue_state` ทำให้ request พร้อมกันรอคิวกัน — ไม่คำนวณเลขซ้ำ

## หมายเหตุการส่งงาน

- Repository:
  - [https://github.com/Eursukkul/interview-question-005]
  - [https://gitlab.com/chalermphan.eur/interview-question-005]
- Package/module ใช้ domain **`example.com`** ตามข้อกำหนดสอบ
- ไฟล์ที่ไม่ commit: `.env`, `node_modules/`, `frontend/dist/`, `frontend/.vite/`

## Troubleshooting

| อาการ | แนวทางแก้ |
|-------|-----------|
| Backend ต่อ DB ไม่ได้ | ตรวจ `docker compose ps` ว่า postgres healthy · ตรง `DATABASE_URL` |
| Frontend เรียก API ไม่ได้ | ตรวจ backend รันที่ `:8080` · `VITE_API_BASE_URL` ใน `frontend/.env` |
| คิวผิดปกติ / schema เก่า | `docker compose down -v && docker compose up -d` |
| `go test ./tests` skip | ตั้ง `TEST_DATABASE_URL` ตามตัวอย่างในหัวข้อ Tests |
