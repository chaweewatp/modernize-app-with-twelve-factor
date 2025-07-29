# modernize-app-with-twelve-factor

![modern-app](modern-app.png)
เป้าหมาย: เพื่อใช้ฝึกการพัฒนา app ตาม Twelve Factor App ในแต่ละขั้นตอน เริ่มจากโปรเจกต์ที่ไม่เป็น Twelve-Factor แล้วค่อยๆ refactor ทีละข้อ

## การ Setup ระบบ

ปกติขั้นตอนการ run ระบบขึ้นมา Admin ต้องพิมพ์คำสั่งตามขั้นตอนดังนี้นี้

### ขั้นตอนที่ 1 สร้าง database

```
docker run -d \
  --name my-postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=mydb \
  -p 5432:5432 \
  -v pgdata:/var/lib/postgresql/data \
  postgres:14
```

### ขั้นตอนที่ 2 การรัน Radis

```
docker run -d \
  --name my-redis \
  -p 6379:6379 \
  redis:7
```

### ขั้นตอนที่ 3 การรัน Service backend

```
cd backend
go mod init github.com/your-org/backend
go mod tidy
go build -o server
./server
```

### ขั้นตอนที่ 4 การรัน Service frontend

```
cd frontend
npm install
npm run dev
```

### ขั้นตอนที่ 5 ทดสอบการใช้งาน

```
http://localhost:3000
```

## modernize ทีละ step   ✅❌

### Step 1: ลองรีวิวโค้ดดูก่อนว่ามีส่วนไหนมั๊ยที่ไม่ตรงตาม Twelve Factor

##### 1.Codebase ✅

มีการทำ version control สำหรับ code แล้ว

##### 2.Dependencies ✅

- ใน backend มีการเก็บ depedency พร้อม version ไว้ใน go.mod เรียบร้อยแล้ว
- ใน frontend มีการเก็บ depedency พร้อม version ไว้ใน package.json เรียบร้อยแล้ว

##### 3.Config ❌

- ไม่มีการเรียกใช้ configuration ใน environment

##### 4.Backing services ✅

- มีการเรียกใช้ database, redis แยกออกจากตัวโปรแกรม อยู่แล้ว

##### 5.Build, Release and Run ❌

##### 6.Process ❌

##### 7.Port binding ✅

มี code อยู่แล้วใน backend

```
	fmt.Println("Backend running on port", port)
	r.Run(":" + port)
```

##### 8.Concurrency ❌

##### 9.Disposability ❌

##### 10.Dev/prod parity ❌

##### 11.Logs ❌

##### 12.Admin process ❌

- ปัญหา
  Admin ต้องทำงานหลายขั้นตอน หรือบางทีก็จำวิธีการเดิมไม่ได้ ส่งผลให้เกิดการทำงานผิดพลาดและล่าช้า

##### สรุป

| topic                     | pass |
| ------------------------- | :--: |
| 1. Codebase               |  ✅  |
| 2. Dependencies           |  ✅  |
| 3. Config                 |  ❌  |
| 4. Backing services       |  ✅  |
| 5. Build, Release and Run |  ❌  |
| 6. Process                |  ❌  |
| 7. Port binding           |  ✅  |
| 8. Concurrency            |  ❌  |
| 9. Disposability          |  ❌  |
| 10. Dev/prod parity       |  ❌  |
| 11. Logs                  |  ❌  |
| 12. Admin process         |  ❌  |

### Step 2: ค่อยๆ แก้ไขไปทีละจุด

#### เพิ่ม 3.config

มีอยู่ 2 รูปแบบคือจัดเก็บใน environment ของเครื่อง และใน .env file ซึ่งมีข้อดีข้อเสียแตกต่างกัน
ในที่นี้เราจะจัดเก็บใน environment ของเครื่อง

```
export APP_PORT=9000
export DATABASE_DSN="host=localhost user=postgres password=postgres dbname=mydb port=5432 sslmode=disable"
export REDIS_ADDR="localhost:6380"
```

ดูว่ามีค่า config ใน environment หรือยัง

```bash
env | grep port
env | grep dsn
env | grep redisAddr
```

แก้ไข code

```
// Load config
port := "8080"
dsn := "host=localhost user=postgres password=postgres dbname=mydb port=5432 sslmode=disable"
redisAddr := "localhost:6379"
```

เป็น

```
// Load config from environment variables
port := os.Getenv("APP_PORT")
if port == "" {
  port = "8080" // fallback default
}

dsn := os.Getenv("DATABASE_DSN")
if dsn == "" {
  dsn = "host=localhost user=postgres password=postgres dbname=mydb port=5432 sslmode=disable"
}

redisAddr := os.Getenv("REDIS_ADDR")
if redisAddr == "" {
  redisAddr = "localhost:6379"
}

fmt.Println("Port:", port)
fmt.Println("DSN:", dsn)
fmt.Println("Redis Addr:", redisAddr)
```

#### ปรับปรุง 5.Build, Release and Run

จากเดิมเราต้องมา run command เพื่อสร้าง databae, redis, backend, frontend ที่นี้เราสามารถมัดรวมกันทีเดียว โดยมีตัวจัดการเป็น docker-compose
โดยเราต้องสร้างไฟล์ docker-compose และ Dockerfile ใน backend และ frontend

```yml
### docker-compose.yml
version: "3.9"
services:
  frontend:
    build: ./frontend
    ports:
      - "3000:3000"

  backend:
    build: ./backend
    ports:
      - "8080:8080"
    depends_on:
      - db
      - cache

  db:
    image: postgres:14
    restart: always
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: mydb
    volumes:
      - postgres-data:/var/lib/postgresql/data

  cache:
    image: redis:7
    restart: always

volumes:
  postgres-data:
```

```bash
### Dockerfile ใน backend
FROM golang:1.21-alpine

WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o server

CMD ["./server"]
```

```bash
### Dockerfile ใน frontend
FROM node:18-alpine

WORKDIR /app
COPY . .
RUN npm install
RUN npm run build

EXPOSE 3000
CMD ["npm", "start"]
```

แล้วเราใช้คำสั่งเดียวในการรัน app

```bash
docker compose up --build -d
### output
```

#### ปรับปรุง 12.Admin Process

สร้าง makefile ขึ้นมา เพื่อเพิ่่ม process ในการ run, restart และ down ระบบขึ้นมา เพื่อให้ admin สามารถมาทำงานได้ในคำสั่งเดียว

```bash
###Makefile
run:
	docker compose up --build -d

restart:
	docker compose restart

down:
	docker compose down -v
```

admin สามารถรันคำสั่ง make ตามด้วย keyword ที่กำหนดไว้ได้

```bash
make run
make restart
make down
```
