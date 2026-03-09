# Homework Submission

**Họ tên:** Nguyễn Nhật Minh

## Các bài đã hoàn thành

- [x] Bài 1: Statistics APIs
- [x] Bài 2: Batch Create
- [x] Bài 3: Batch Delete
- [x] Bài 4: Connection Retry
- [x] Bài 5: Health Check
- [x] Bài 6: Pagination
- [x] Bài 7: Search

## Hướng dẫn chạy project

```powershell
# 1. Khởi động Docker Desktop

# 2. Chạy PostgreSQL
cd homeworks/Day2/mini-asm
docker-compose up -d

# 3. Chạy server
go run cmd/server/main.go
```

---

## Bài 1: Statistics APIs

### GET /assets/stats
![Bài 1 - Stats](screenshots/bai1-stats.png)

### GET /assets/count?type=domain&status=active
![Bài 1 - Count](screenshots/bai1-count.png)

---

## Bài 2: Batch Create

### POST /assets/batch (Success - 201 Created)
![Bài 2 - Success](screenshots/bai2-success.png)

### POST /assets/batch (Error - 400 Rollback)
![Bài 2 - Error](screenshots/bai2-error.png)

---

## Bài 3: Batch Delete

### DELETE /assets/batch?ids=...
![Bài 3 - Delete](screenshots/bai3-delete.png)

---

## Bài 4: Connection Retry

### Server logs khi DB tắt → retry → DB bật lại → connected
![Bài 4 - Retry](screenshots/bai4-retry.png)

---

## Bài 5: Health Check

### GET /health (DB connected - 200 OK)
![Bài 5 - Health](screenshots/bai5-health.png)

---

## Bài 6: Pagination

### GET /assets?page=1&limit=2
![Bài 6 - Page 1](screenshots/bai6-page1.png)

### GET /assets?page=2&limit=2
![Bài 6 - Page 2](screenshots/bai6-page2.png)

---

## Bài 7: Search

### GET /assets/search?q=.com
![Bài 7 - Search](screenshots/bai7-search.png)

### GET /assets/search?q=google
![Bài 7 - Search](screenshots/bai7-search2.png)
