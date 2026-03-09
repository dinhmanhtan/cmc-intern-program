# Homework Submission

**Họ tên:** Nguyễn Thị Huyền

## Hướng dẫn chạy

Do sử dụng database PostgreSQL lưu trữ online trên Neon, nên sử dụng key có sẵn trong file .env.example để kết nối database.

1. Clone repository về máy
2. cd homeworks/submissions/assets-api
3. go mod tidy
4. copy .env.example .env
5. go run cmd/server/main.go

## Các bài đã hoàn thành

- [x] Bài 1: Statistics APIs
      ![alt text](image.png)
      ![alt text](image-1.png)
      ![alt text](image-2.png)
      ![alt text](image-3.png)
- [x] Bài 2: Batch Create
      ![alt text](image-4.png)
      ![alt text](image-5.png)
- [x] Bài 3: Batch Delete
      ![alt text](image-6.png)

- [x] Bài 4: Connection Retry
      Khi thử viết sai mật khẩu database, API sẽ trả về lỗi kết nối. Sau khi sửa lại mật khẩu đúng, API sẽ tự động kết nối lại và hoạt động bình thường.
      ![alt text](image-8.png)
- [x] Bài 5: Health Check
      ![alt text](image-7.png)
- [x] Bài 6: Pagination (Bonus)
      ![alt text](image-9.png)
      ![alt text](image-10.png)
      ![alt text](image-11.png)
- [x] Bài 7: Search (Bonus)
      ![alt text](image-12.png)
      ![alt text](image-13.png)
      ![alt text](image-14.png)
