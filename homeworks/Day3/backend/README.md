# Mini EASM API - Day 3 Homework

## Overview

Day 3 mở rộng hệ thống EASM với **IP Scan**, **Port Scan**, **CORS**, **CI/CD**, **Docker**, và **Export Reports**.

## Quick Start

```bash
# Start database
cd homeworks/Day3
docker compose up -d db

# Start backend
cd homeworks/Day3/backend
go run ./cmd/server

# Start frontend (new terminal)
cd homeworks/Day3/frontend
npm install && npm run dev
# → http://localhost:5173
```

## Docker Compose (Full Stack)

```bash
cd homeworks/Day3
docker compose up -d

# Services:
# - DB:       localhost:7432 (PostgreSQL)
# - Backend:  http://localhost:8080
# - Frontend: http://localhost:3000
```

## API

```bash
# Health check
curl http://localhost:8080/health

# Create IP asset
curl -s -X POST http://localhost:8080/assets \
  -H "Content-Type: application/json" \
  -d '{"name":"127.0.0.1","type":"ip"}' | jq

# IP Scan (geolocation + ASN)
curl -s -X POST http://localhost:8080/assets/{ID}/scan \
  -H "Content-Type: application/json" \
  -d '{"scan_type":"ip"}' | jq

# Port Scan (localhost/private IPs only!)
curl -s -X POST http://localhost:8080/assets/{ID}/scan \
  -H "Content-Type: application/json" \
  -d '{"scan_type":"port"}' | jq

# Export assets as CSV
curl -o assets.csv "http://localhost:8080/assets/export?format=csv"

# Export assets as JSON
curl -o assets.json "http://localhost:8080/assets/export?format=json"

# Export scan results for 1 asset
curl -o results.json "http://localhost:8080/assets/{ID}/results/export"
```

## Unit Tests

```bash
cd homeworks/Day3/backend

# Run all unit tests
go test ./internal/model/... ./internal/scanner/... ./internal/handler/... -v

# Run with coverage
go test -cover ./...

# Generate HTML coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## What's New in Day 3

| Feature | File | Mô tả |
|---------|------|-------|
| IP Scan | `scanner/ip_scanner.go` | Geolocation + ASN (ip-api.com) |
| Port Scan | `scanner/port_scanner.go` | TCP scan, private IPs only |
| CORS Middleware | `middleware/cors.go` | Cho phép frontend kết nối |
| Export Reports | `handler/export_handler.go` | CSV/JSON download |
| Unit Tests | `scanner/*_test.go` | 13 test functions mới |
| CI/CD | `.github/workflows/ci.yml` | Gosec, Trivy, Gitleaks |

## Project Structure

```
homeworks/Day3/backend/
├── cmd/server/main.go              ← Entry point (CORS integrated)
├── internal/
│   ├── model/
│   │   ├── scan.go                 ← ScanTypeIP, IPScanResult, PortScanResult
│   │   └── scan_test.go            ← Extended tests (27 pass)
│   ├── scanner/
│   │   ├── ip_scanner.go           ← NEW: IP Geolocation & ASN
│   │   ├── ip_scanner_test.go      ← NEW: 6 tests
│   │   ├── port_scanner.go         ← NEW: TCP Port Scanner
│   │   └── port_scanner_test.go    ← NEW: 7 tests
│   ├── middleware/
│   │   └── cors.go                 ← NEW: CORS Middleware
│   ├── handler/
│   │   ├── export_handler.go       ← NEW: CSV/JSON Export
│   │   └── export_handler_test.go  ← NEW: Export tests
│   └── service/
│       └── scan_service.go         ← Updated: IP/Port dispatch
├── api.yml                         ← Updated: v7.0.0, ip scan type, export paths
└── Dockerfile                      ← Docker build
```

## Security Notes

- **Port Scanner**: Chỉ cho phép scan localhost và private IP ranges (RFC 1918)
  - ✅ 127.x.x.x, 10.x.x.x, 172.16-31.x.x, 192.168.x.x
  - ❌ Public IPs bị từ chối với lỗi `unauthorized: only private/localhost IPs allowed`
- **CORS**: Cấu hình `Access-Control-Allow-Origin: *` (production nên giới hạn origin)
