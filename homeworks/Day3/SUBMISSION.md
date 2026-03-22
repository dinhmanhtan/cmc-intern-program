# Homework Submission - Day 3

**Họ tên:** Nguyễn Nhật Minh

## Các bài đã hoàn thành

- [x] Bài 1: Mở rộng Scan API
- [x] Bài 2: Viết Unit Tests
- [x] Bài 3: Tích hợp Frontend
- [x] Bài 4: CI/CD với GitHub Actions
- [x] Bài 5: Deploy với Docker Compose
- [x] Bài 6: Tính năng EASM mới - Export Reports CSV/JSON (Bonus)
- [ ] Bài 7: Deploy lên Cloud VM (Bonus)
- [ ] Bài 8: Domain & TLS/HTTPS (Bonus)
- [ ] Bài 9: Auto Deploy on Merge (Bonus)

---

## Bài 1: Mở rộng Scan API (25 điểm)

Thêm 2 loại scan mới vào hệ thống:

**IP Scan (`scan_type: "ip"`)** - Passive scan dùng `ip-api.com` lấy geolocation + ASN  
**Port Scan (`scan_type: "port"`)** - Active scan TCP, chỉ cho phép private IPs (safety check)

Files thêm mới:
- `internal/scanner/ip_scanner.go`
- `internal/scanner/port_scanner.go`
- `internal/model/scan.go` → thêm `ScanTypeIP`, `IPScanResult`, `PortScanResult`
- `internal/service/scan_service.go` → thêm dispatch + in-memory cache

### Screenshots

**[Ảnh 1] IP Scan result:**
![IP Scan Result](screenshots/bai1-ip-scan-result.png)

**[Ảnh 2] Port Scan result:**
![Port Scan Result](screenshots/bai1-port-scan-result.png)

**[Ảnh 3] Safety check - Port scan bị từ chối với public IP:**
![Port Scan Safety Check](screenshots/bai1-safety-check.png)

---

## Bài 2: Viết Unit Tests (25 điểm)

| File | Tests |
|------|-------|
| `model/scan_test.go` | Mở rộng + ScanTypeIP (27 tests PASS) |
| `scanner/ip_scanner_test.go` | 6 test functions |
| `scanner/port_scanner_test.go` | 7 test functions |
| `handler/export_handler_test.go` | 5 test functions |

### Screenshots

**[Ảnh 4a] Unit Tests - Model & Scanner:**
![Unit Tests Model](screenshots/bai2-unit-tests-model.png)

**[Ảnh 4b] Unit Tests - All PASS:**
![Unit Tests Pass](screenshots/bai2-unit-tests-pass.png)

---

## Bài 3: Tích hợp Frontend (20 điểm)

- `internal/middleware/cors.go` → CORS middleware mới
- `cmd/server/main.go` → wrap mux với `CORSMiddleware`
- `frontend/.env` → `VITE_API_URL=http://localhost:8080`

### Screenshots

**[Ảnh 5] Danh sách assets trên Frontend:**
![Frontend Asset List](screenshots/bai3-frontend-assets.png)

**[Ảnh 6] Kết quả scan trên Frontend:**
![Frontend Scan Results](screenshots/bai3-frontend-scan-results.png)

---

## Bài 4: CI/CD với GitHub Actions (25 điểm)

File: `.github/workflows/ci.yml`

6 jobs: Backend Build/Test → Frontend Build → Gosec → Trivy → Gitleaks → Docker Build  
Trigger: push vào branch `homework`, paths `homeworks/Day3/**`

### Screenshots

**[Ảnh 7] GitHub Actions - All jobs green:**
![CI All Green](screenshots/bai4-ci-all-green.png)

**[Ảnh 8] Security scan output:**
![Security Scan Output](screenshots/bai4-security-scan.png)

---

## Bài 5: Deploy với Docker Compose (15 điểm)

File: `homeworks/Day3/docker-compose.yml`  
Stack: PostgreSQL 15 (7432) + Go Backend (8080) + Nginx Frontend (3000)

### Screenshots

**[Ảnh 9] docker compose ps - All services running:**
![Docker Compose PS](screenshots/bai5-docker-compose-ps.png)

**[Ảnh 10] Backend health check:**
![Backend Health Check](screenshots/bai5-health-check.png)

**[Ảnh 11] Frontend ở port 3000 (từ Docker):**
![Frontend Docker](screenshots/bai5-frontend-docker.png)

---

## Bài 6: Export Reports - Bonus (15 điểm)

File: `internal/handler/export_handler.go`

| Endpoint | Mô tả |
|----------|-------|
| `GET /assets/export?format=csv` | Download CSV (Excel-ready) |
| `GET /assets/export?format=json` | Download JSON |
| `GET /assets/{id}/results/export` | Download scan results JSON |

### Screenshots

**[Ảnh 12] Export CSV:**
![Export CSV](screenshots/bai6-export-csv.png)

**[Ảnh 13] Export JSON:**
![Export JSON](screenshots/bai6-export-json.png)
