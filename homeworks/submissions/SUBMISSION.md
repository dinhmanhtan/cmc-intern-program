# Homework Submission

**Họ tên:** Trần Hữu Đức

## Các bài đã hoàn thành

- [x] Bài 1: Statistics APIs
- [x] Bài 2: Batch Create
- [x] Bài 3: Batch Delete
- [x] Bài 4: Concurrent-safe Create
- [x] Bài 5: In-memory Health Check
- [x] Bài 6: Pagination (Bonus)
- [x] Bài 7: Search (Bonus)

---

## Bài 1: Statistics APIs (20 điểm)

### 1.1 GET /assets/stats

**Data**
```bash 
❯ curl -s -X POST http://localhost:8080/assets \
   -H "Content-Type: application/json" \
   -d '{"name":"example.com","type":"domain"}' | jq

{
  "id": "e69aee2f-f35d-4164-a4bb-df544d779b8e",
  "name": "example.com",
  "type": "domain",
  "status": "active",
  "created_at": "2026-06-14T10:45:27.50583883+07:00",
  "updated_at": "2026-06-14T10:45:27.505838863+07:00"
}
❯ curl -s -X POST http://localhost:8080/assets \
   -H "Content-Type: application/json" \
   -d '{"name":"192.168.1.1","type":"ip"}' | jq

{
  "id": "e15c3536-4478-4356-979f-9343dc92e82b",
  "name": "192.168.1.1",
  "type": "ip",
  "status": "active",
  "created_at": "2026-06-14T10:45:40.584160964+07:00",
  "updated_at": "2026-06-14T10:45:40.584161021+07:00"
}
❯ curl -s -X POST http://localhost:8080/assets \
   -H "Content-Type: application/json" \
   -d '{"name":"web-service","type":"service","status":"inactive"}' | jq

{
  "id": "b0d49bff-f94c-4dd2-a3bf-aa4a7bfc1b69",
  "name": "web-service",
  "type": "service",
  "status": "active",
  "created_at": "2026-06-14T10:45:53.248061676+07:00",
  "updated_at": "2026-06-14T10:45:53.248061705+07:00"
}
```

**Request:**
```bash
curl http://localhost:8080/assets/stats
```

**Response:**
```json
{"total":3,"by_type":{"domain":1,"ip":1,"service":1},"by_status":{"active":3}}
```

### 1.2 GET /assets/count

**Request (all):**
```bash
curl http://localhost:8080/assets/count
```

**Response:**
```json
{"count": 3, "filters": {}}
```

**Request (filter by type):**
```bash
curl "http://localhost:8080/assets/count?type=domain"
```

**Response:**
```json
{"count": 1, "filters": {"type": "domain"}}
```

**Request (filter by type + status):**
```bash
curl "http://localhost:8080/assets/count?type=domain&status=active"
```

**Response:**
```json
{"count": 1, "filters": {"type": "domain", "status": "active"}}
```

---

## Bài 2: Batch Create Assets (25 điểm)

### Success case

**Data:**
```bash
❯ curl -s -X POST http://localhost:8080/assets/batch \
   -H "Content-Type: application/json" \
   -d '{
   "assets": [
   {"name":"batch1.com","type":"domain"},
   {"name":"batch2.com","type":"domain"},
   {"name":"10.0.0.1","type":"ip"}
   ]
   }' | jq

{
  "created": 3,
  "ids": [
    "e6e6f281-cdfc-4293-a3e4-3bd836e5241a",
    "3bde55bd-9b27-4054-be08-4c29cd3c0466",
    "d859c7a5-1e84-402a-8904-2c1bf831905c"
  ]
}
```

### Error case (invalid type - all or nothing)

**Request:**
```bash
❯ curl -s -X POST http://localhost:8080/assets/batch \
   -H "Content-Type: application/json" \
   -d '{
   "assets": [
   {"name":"good.com","type":"domain"},
   {"name":"bad.com","type":"invalid_type"}
   ]
   }' | jq
```

**Response (400 Bad Request):**
```json
{"error": "invalid asset type: must be domain, ip, or service"}
```

**Verify - count vẫn = 6 (không có asset nào được tạo):**
```bash
curl http://localhost:8080/assets/count
```
```json
{"count":6,"filters":{}}
```

---

## Bài 3: Batch Delete Assets (20 điểm)

### Tạo 2 assets để xóa

**Request:**
```bash
❯ ID1=$(curl -s -X POST http://localhost:8080/assets \
   -H "Content-Type: application/json" \
   -d '{"name":"delete-me1.com","type":"domain"}' | jq -r '.id')

  ID2=$(curl -s -X POST http://localhost:8080/assets \
   -H "Content-Type: application/json" \
   -d '{"name":"delete-me2.com","type":"domain"}' | jq -r '.id')

  echo "ID1=$ID1 ID2=$ID2"

ID1=7dff78b3-84d7-4d0a-ac38-e9c487fd0ab0 ID2=49c809f5-2c55-4ab8-9090-4bfa31f34205
```

### Batch delete (2 real IDs + 1 fake ID)

**Request:**
```bash
curl -X DELETE "http://localhost:8080/assets/batch?ids=$ID1,$ID2,fake-uuid-123"
```

**Response (200 OK):**
```json
{"deleted": 2, "not_found": 1}
```

---

## Bài 4: Concurrent-safe Create (25 điểm)

### Bắn 20 request song song

**Request:**
```bash
for i in $(seq 1 20); do
  curl -s -X POST http://localhost:8080/assets \
    -H "Content-Type: application/json" \
    -d "{\"name\":\"concurrent-$i.com\",\"type\":\"domain\"}" > /dev/null &
done
wait
```

### Verify - tất cả 20 assets được tạo thành công

**Request:**
```bash
curl http://localhost:8080/assets/count
```

**Response:**
```json
{"count": 26, "filters": {}}
```

- Server không crash, không data corruption
- 26 = 6 (existing) + 2 (batch create) + 20 (concurrent create) — chính xác

---

## Bài 5: In-memory Health Check (15 điểm)

### GET /health

**Request:**
```bash
curl http://localhost:8080/health
```

**Response (200 OK):**
```json
{"status":"ok","storage":{"type":"in-memory","asset_count":26},"uptime_seconds":10195.698410865,"timestamp":"2026-06-14T13:27:47+07:00"}
```

---

## Bài 6: Pagination & Filtering (15 điểm) - BONUS

### Page 1, limit 5

**Request:**
```bash
curl "http://localhost:8080/assets?page=1&limit=5"
```

**Response (200 OK):**
```json
{"data":[{"id":"5056a5b2-0735-4935-afec-5f8034590fe8","name":"concurrent-6.com","type":"domain","status":"active","created_at":"2026-06-14T13:27:01.829253125+07:00","updated_at":"2026-06-14T13:27:01.829253199+07:00"},{"id":"cffc8480-6160-48e8-8985-f316b4c56d2c","name":"concurrent-9.com","type":"domain","status":"active","created_at":"2026-06-14T13:27:01.829185062+07:00","updated_at":"2026-06-14T13:27:01.829185146+07:00"},{"id":"8270cd01-1abd-4e21-bef4-e7f1ca6397d5","name":"concurrent-11.com","type":"domain","status":"active","created_at":"2026-06-14T13:27:01.829183926+07:00","updated_at":"2026-06-14T13:27:01.829183978+07:00"},{"id":"4632ad1d-15aa-41c7-93d6-8eb9b2275eb0","name":"concurrent-12.com","type":"domain","status":"active","created_at":"2026-06-14T13:27:01.829181481+07:00","updated_at":"2026-06-14T13:27:01.829181554+07:00"},{"id":"f398d995-2563-4369-86f1-f833921c5fcf","name":"concurrent-2.com","type":"domain","status":"active","created_at":"2026-06-14T13:27:01.829174432+07:00","updated_at":"2026-06-14T13:27:01.829174525+07:00"}],"pagination":{"page":1,"limit":5,"total":26,"total_pages":6}}
```

---

## Bài 7: Search by Name (10 điểm) - BONUS

### Search "example"

**Request:**
```bash
curl "http://localhost:8080/assets/search?q=example"
```

**Response:** 
```bash
[{"id":"e69aee2f-f35d-4164-a4bb-df544d779b8e","name":"example.com","type":"domain","status":"active","created_at":"2026-06-14T10:45:27.50583883+07:00","updated_at":"2026-06-14T10:45:27.505838863+07:00"}]
```

### Search ".com"

**Request:**
```bash
curl "http://localhost:8080/assets/search?q=.com"
```

**Response:** 
```bash
[{"id":"5056a5b2-0735-4935-afec-5f8034590fe8","name":"concurrent-6.com","type":"domain","status":"active","created_at":"2026-06-14T13:27:01.829253125+07:00","updated_at":"2026-06-14T13:27:01.829253199+07:00"},{"id":"cffc8480-6160-48e8-8985-f316b4c56d2c","name":"concurrent-9.com","type":"domain","status":"active","created_at":"2026-06-14T13:27:01.829185062+07:00","updated_at":"2026-06-14T13:27:01.829185146+07:00"},{"id":"8270cd01-1abd-4e21-bef4-e7f1ca6397d5","name":"concurrent-11.com","type":"domain","status":"active","created_at":"2026-06-14T13:27:01.829183926+07:00","updated_at":"2026-06-14T13:27:01.829183978+07:00"},{"id":"4632ad1d-15aa-41c7-93d6-8eb9b2275eb0","name":"concurrent-12.com","type":"domain","status":"active","created_at":"2026-06-14T13:27:01.829181481+07:00","updated_at":"2026-06-14T13:27:01.829181554+07:00"},{"id":"f398d995-2563-4369-86f1-f833921c5fcf","name":"concurrent-2.com","type":"domain","status":"active","created_at":"2026-06-14T13:27:01.829174432+07:00","updated_at":"2026-06-14T13:27:01.829174525+07:00"},{"id":"3cffe8d3-5945-4c9c-89bf-d797b3e62631","name":"concurrent-8.com","type":"domain","status":"active","created_at":"2026-06-14T13:27:01.829169614+07:00","updated_at":"2026-06-14T13:27:01.82916971+07:00"},{"id":"72864e54-be11-4f41-906f-8ae771b08bdb","name":"concurrent-18.com","type":"domain","status":"active","created_at":"2026-06-14T13:27:01.829165824+07:00","updated_at":"2026-06-14T13:27:01.829165879+07:00"},{"id":"2e404c20-2da2-40da-b871-1e63eea22604","name":"concurrent-7.com","type":"domain","status":"active","created_at":"2026-06-14T13:27:01.829158516+07:00","updated_at":"2026-06-14T13:27:01.829158559+07:00"},{"id":"7779c7ba-3042-4c69-aab7-599cf17608fb","name":"concurrent-17.com","type":"domain","status":"active","created_at":"2026-06-14T13:27:01.82915762+07:00","updated_at":"2026-06-14T13:27:01.829157682+07:00"},{"id":"f7aeb49b-2135-4d87-bef1-20806ef3c13b","name":"concurrent-1.com","type":"domain","status":"active","created_at":"2026-06-14T13:27:01.829155365+07:00","updated_at":"2026-06-14T13:27:01.82915543+07:00"},{"id":"7a325e0e-7e39-4057-a8be-306344e6261d","name":"concurrent-20.com","type":"domain","status":"active","created_at":"2026-06-14T13:27:01.829154965+07:00","updated_at":"2026-06-14T13:27:01.82915509+07:00"},{"id":"378207b3-b4e7-4b9c-b809-0229bc2fb592","name":"concurrent-3.com","type":"domain","status":"active","created_at":"2026-06-14T13:27:01.829147597+07:00","updated_at":"2026-06-14T13:27:01.829147658+07:00"},{"id":"c1ef6ba2-700c-4813-b5d4-e08c8483f117","name":"concurrent-15.com","type":"domain","status":"active","created_at":"2026-06-14T13:27:01.829147023+07:00","updated_at":"2026-06-14T13:27:01.829147106+07:00"},{"id":"8e58374a-11c1-4660-baaa-a1c48cccdcf3","name":"concurrent-5.com","type":"domain","status":"active","created_at":"2026-06-14T13:27:01.82914647+07:00","updated_at":"2026-06-14T13:27:01.82914653+07:00"},{"id":"5384d44a-4fc9-4bc2-80fc-62d291a7cc47","name":"concurrent-4.com","type":"domain","status":"active","created_at":"2026-06-14T13:27:01.829144713+07:00","updated_at":"2026-06-14T13:27:01.829144781+07:00"},{"id":"ce833b92-94b7-42e1-a6e1-7d320fd78ce8","name":"concurrent-14.com","type":"domain","status":"active","created_at":"2026-06-14T13:27:01.82914449+07:00","updated_at":"2026-06-14T13:27:01.829144529+07:00"},{"id":"88110d33-3fd2-4c93-8236-489cb738a3cc","name":"concurrent-13.com","type":"domain","status":"active","created_at":"2026-06-14T13:27:01.829139881+07:00","updated_at":"2026-06-14T13:27:01.82913995+07:00"},{"id":"5448676a-aeea-49d0-a2de-f9aae0ead4e0","name":"concurrent-19.com","type":"domain","status":"active","created_at":"2026-06-14T13:27:01.829139251+07:00","updated_at":"2026-06-14T13:27:01.829139393+07:00"},{"id":"9781f8b3-7c0f-41d6-b41a-b7f60bc60c37","name":"concurrent-16.com","type":"domain","status":"active","created_at":"2026-06-14T13:27:01.829138993+07:00","updated_at":"2026-06-14T13:27:01.829139059+07:00"},{"id":"ce62e2a4-0d8c-45fd-afe6-61762a5d4630","name":"concurrent-10.com","type":"domain","status":"active","created_at":"2026-06-14T13:27:01.829138875+07:00","updated_at":"2026-06-14T13:27:01.829138951+07:00"},{"id":"3bde55bd-9b27-4054-be08-4c29cd3c0466","name":"batch2.com","type":"domain","status":"active","created_at":"2026-06-14T10:49:02.935555881+07:00","updated_at":"2026-06-14T10:49:02.935555881+07:00"},{"id":"e6e6f281-cdfc-4293-a3e4-3bd836e5241a","name":"batch1.com","type":"domain","status":"active","created_at":"2026-06-14T10:49:02.935555881+07:00","updated_at":"2026-06-14T10:49:02.935555881+07:00"},{"id":"e69aee2f-f35d-4164-a4bb-df544d779b8e","name":"example.com","type":"domain","status":"active","created_at":"2026-06-14T10:45:27.50583883+07:00","updated_at":"2026-06-14T10:45:27.505838863+07:00"}]
```

### Search "CONCURRENT" (case-insensitive)

**Request:**
```bash
curl "http://localhost:8080/assets/search?q=CONCURRENT"
```

**Response:** 
```bash
[{"id":"5056a5b2-0735-4935-afec-5f8034590fe8","name":"concurrent-6.com","type":"domain","status":"active","created_at":"2026-06-14T13:27:01.829253125+07:00","updated_at":"2026-06-14T13:27:01.829253199+07:00"},{"id":"cffc8480-6160-48e8-8985-f316b4c56d2c","name":"concurrent-9.com","type":"domain","status":"active","created_at":"2026-06-14T13:27:01.829185062+07:00","updated_at":"2026-06-14T13:27:01.829185146+07:00"},{"id":"8270cd01-1abd-4e21-bef4-e7f1ca6397d5","name":"concurrent-11.com","type":"domain","status":"active","created_at":"2026-06-14T13:27:01.829183926+07:00","updated_at":"2026-06-14T13:27:01.829183978+07:00"},{"id":"4632ad1d-15aa-41c7-93d6-8eb9b2275eb0","name":"concurrent-12.com","type":"domain","status":"active","created_at":"2026-06-14T13:27:01.829181481+07:00","updated_at":"2026-06-14T13:27:01.829181554+07:00"},{"id":"f398d995-2563-4369-86f1-f833921c5fcf","name":"concurrent-2.com","type":"domain","status":"active","created_at":"2026-06-14T13:27:01.829174432+07:00","updated_at":"2026-06-14T13:27:01.829174525+07:00"},{"id":"3cffe8d3-5945-4c9c-89bf-d797b3e62631","name":"concurrent-8.com","type":"domain","status":"active","created_at":"2026-06-14T13:27:01.829169614+07:00","updated_at":"2026-06-14T13:27:01.82916971+07:00"},{"id":"72864e54-be11-4f41-906f-8ae771b08bdb","name":"concurrent-18.com","type":"domain","status":"active","created_at":"2026-06-14T13:27:01.829165824+07:00","updated_at":"2026-06-14T13:27:01.829165879+07:00"},{"id":"2e404c20-2da2-40da-b871-1e63eea22604","name":"concurrent-7.com","type":"domain","status":"active","created_at":"2026-06-14T13:27:01.829158516+07:00","updated_at":"2026-06-14T13:27:01.829158559+07:00"},{"id":"7779c7ba-3042-4c69-aab7-599cf17608fb","name":"concurrent-17.com","type":"domain","status":"active","created_at":"2026-06-14T13:27:01.82915762+07:00","updated_at":"2026-06-14T13:27:01.829157682+07:00"},{"id":"f7aeb49b-2135-4d87-bef1-20806ef3c13b","name":"concurrent-1.com","type":"domain","status":"active","created_at":"2026-06-14T13:27:01.829155365+07:00","updated_at":"2026-06-14T13:27:01.82915543+07:00"},{"id":"7a325e0e-7e39-4057-a8be-306344e6261d","name":"concurrent-20.com","type":"domain","status":"active","created_at":"2026-06-14T13:27:01.829154965+07:00","updated_at":"2026-06-14T13:27:01.82915509+07:00"},{"id":"378207b3-b4e7-4b9c-b809-0229bc2fb592","name":"concurrent-3.com","type":"domain","status":"active","created_at":"2026-06-14T13:27:01.829147597+07:00","updated_at":"2026-06-14T13:27:01.829147658+07:00"},{"id":"c1ef6ba2-700c-4813-b5d4-e08c8483f117","name":"concurrent-15.com","type":"domain","status":"active","created_at":"2026-06-14T13:27:01.829147023+07:00","updated_at":"2026-06-14T13:27:01.829147106+07:00"},{"id":"8e58374a-11c1-4660-baaa-a1c48cccdcf3","name":"concurrent-5.com","type":"domain","status":"active","created_at":"2026-06-14T13:27:01.82914647+07:00","updated_at":"2026-06-14T13:27:01.82914653+07:00"},{"id":"5384d44a-4fc9-4bc2-80fc-62d291a7cc47","name":"concurrent-4.com","type":"domain","status":"active","created_at":"2026-06-14T13:27:01.829144713+07:00","updated_at":"2026-06-14T13:27:01.829144781+07:00"},{"id":"ce833b92-94b7-42e1-a6e1-7d320fd78ce8","name":"concurrent-14.com","type":"domain","status":"active","created_at":"2026-06-14T13:27:01.82914449+07:00","updated_at":"2026-06-14T13:27:01.829144529+07:00"},{"id":"88110d33-3fd2-4c93-8236-489cb738a3cc","name":"concurrent-13.com","type":"domain","status":"active","created_at":"2026-06-14T13:27:01.829139881+07:00","updated_at":"2026-06-14T13:27:01.82913995+07:00"},{"id":"5448676a-aeea-49d0-a2de-f9aae0ead4e0","name":"concurrent-19.com","type":"domain","status":"active","created_at":"2026-06-14T13:27:01.829139251+07:00","updated_at":"2026-06-14T13:27:01.829139393+07:00"},{"id":"9781f8b3-7c0f-41d6-b41a-b7f60bc60c37","name":"concurrent-16.com","type":"domain","status":"active","created_at":"2026-06-14T13:27:01.829138993+07:00","updated_at":"2026-06-14T13:27:01.829139059+07:00"},{"id":"ce62e2a4-0d8c-45fd-afe6-61762a5d4630","name":"concurrent-10.com","type":"domain","status":"active","created_at":"2026-06-14T13:27:01.829138875+07:00","updated_at":"2026-06-14T13:27:01.829138951+07:00"}]
```

---

## Kiến trúc triển khai

Tuân thủ Clean Architecture (handler → service → storage):

```
internal/
├── model/
│   └── asset.go          ← Thêm Statistics, Pagination, BatchCreate/Response structs
├── storage/
│   ├── storage.go        ← Thêm interface: GetStatistics, Count, BatchCreate, BatchDelete
│   └── memory/
│       └── memory.go     ← Implement: sync.RWMutex (Bài 4), loop-based counting, all-or-nothing batch
├── service/
│   └── asset_service.go  ← Validate-before-write (Bài 2), pagination logic, storage info
└── handler/
    ├── asset_handler.go  ← GET /assets/stats, /count, /search; POST /batch; DELETE /batch; pagination
    └── health_handler.go ← Nhận *AssetService, trả storage.type + asset_count + uptime
```

### Cách chạy project

```bash
cd app/session2-basic-api
go run cmd/server/main.go
```

Server sẽ lắng nghe trên `http://localhost:8080`.
