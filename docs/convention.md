# LokiForce Standard Coding & Architecture Conventions

Tài liệu này định nghĩa các tiêu chuẩn lập trình, quy tắc thiết kế kiến trúc và phong cách code áp dụng thống nhất trong toàn bộ dự án **LokiForce**.

---

## 1. Cấu Trúc Thư Mục (Clean Architecture)

Mỗi module trong thư mục `internal/` (ví dụ: `user`, `organization`, `project`, `team`) phải tuân thủ nghiêm ngặt mô hình kiến trúc Clean Architecture được chia làm 4 lớp chính:

```text
internal/module_name/
├── domain/            # Lớp Nghiệp vụ cốt lõi (Pure Domain)
│   ├── entity.go      # Khai báo cấu trúc dữ liệu Entity và Logic nghiệp vụ tự thân
│   ├── repository.go  # Định nghĩa Interface giao tiếp DB (không có mã thực thi)
│   └── error.go       # Định nghĩa các lỗi nghiệp vụ (Domain Errors)
│
├── application/       # Lớp Ứng dụng / Điều phối (Usecases)
│   ├── usecase.go     # Định nghĩa Input/Output DTOs và Interface Usecase công khai
│   ├── usecase_impl.go# Hiện thực hóa các usecase nghiệp vụ cụ thể
│   └── usecase_test.go# Kiểm thử đơn vị cho lớp Usecase
│
├── infrastructure/    # Lớp Hạ tầng (Adapters)
│   └── repository/
│       └── postgres.go# Hiện thực hóa repository bằng GORM / SQL cụ thể
│
└── delivery/          # Lớp Giao tiếp ngoài (HTTP Handlers / CLIs)
    └── http/
        ├── handler.go # Nhận request, validate dữ liệu, gọi Usecase, trả về JSON
        └── router.go  # Định nghĩa các tuyến đường API bằng Gin Gonic
```

### Nguyên tắc ranh giới phụ thuộc (Dependency Rule)
* **Domain** là lớp cốt lõi nhất. **Không được phép** import bất kỳ thư viện bên ngoài nào liên quan đến Database (GORM, SQL), Web Framework (Gin) hay lớp Application/Infrastructure.
* Mã nguồn ở các lớp ngoài được phép gọi vào lớp trong, lớp trong giao tiếp với lớp ngoài thông qua các **Interface** (Dependency Inversion).

---

## 2. Tiêu Chuẩn Viết Code Go (Go Coding Conventions)

### 2.1. Đặt tên (Naming Conventions)
* **Exported Identifiers:** Sử dụng chữ hoa đầu từ (`PascalCase`) cho struct, interface, function được phép truy cập từ bên ngoài package.
* **Internal/Local Identifiers:** Sử dụng chữ thường đầu từ (`camelCase`) cho biến cục bộ, struct private, hoặc các handler method nội bộ.
* **Interface:** Đặt tên ngắn gọn, rõ nghĩa (ví dụ: `UserUsecase`, `UserRepository`), tránh sử dụng tiền tố hoặc hậu tố thừa thãi nếu package name đã mô tả rõ.

### 2.2. Xử lý lỗi (Error Handling)
* Luôn khai báo các lỗi nghiệp vụ cụ thể trong file `domain/error.go` bằng `errors.New` thay vì dùng các chuỗi chữ thường rải rác.
* Lớp HTTP Handler sẽ dựa vào kiểu lỗi để trả về HTTP status code phù hợp (ví dụ: `domain.ErrUserNotFound` -> `404 Not Found`).

### 2.3. Mã hóa & Bảo mật (Security & Hashing)
* Tuyệt đối không lưu trữ mật khẩu dưới dạng văn bản thô (Plain Text) ở bất kỳ tầng nào.
* Việc kiểm tra mật khẩu hợp lệ và băm mật khẩu bằng `bcrypt` phải được thực hiện trực tiếp bên trong Domain Entity (`NewUser`) hoặc qua Domain Service trước khi gửi xuống database lưu trữ.
* Sử dụng phương thức mã hóa được đóng gói trong Entity để so khớp: `user.VerifyPassword(rawPassword)`.

---

## 3. Cấu hình (Configuration) & Dependency Injection

### 3.1. Quản lý cấu hình bằng Viper (YAML)
* Toàn bộ cấu hình hệ thống được lưu trữ trong file `config.yaml` và được tải lên thông qua **Viper**.
* Struct cấu hình phải ánh xạ cấu trúc phân cấp rõ ràng (Server, Database, Security):
  ```go
  type Config struct {
      Server   ServerConfig   `mapstructure:"server"`
      Database DatabaseConfig `mapstructure:"database"`
  }
  ```

### 3.2. Tiêm phụ thuộc (Dependency Injection) bằng Google Wire
* Mọi cấu phần dịch vụ phải được kết nối lỏng lẻo (loosely coupled) và tự động lắp ghép thông qua công cụ phát sinh mã nguồn tĩnh **Google Wire**.
* Khai báo liên kết interface với struct cụ thể trong file `cmd/api/wire.go`:
  ```go
  wire.Bind(new(domain.UserRepository), new(*repository.PostgresUserRepository))
  ```
* Chạy cập nhật DI bằng lệnh: `go run github.com/google/wire/cmd/wire gen` ở thư mục chứa tệp `wire.go`.

---

## 4. Chuẩn hóa Phản Hồi API (API Response Standard)

Tất cả các API trả về của hệ thống (Gin Gonic) phải sử dụng gói tiện ích dùng chung **`pkg/response`** để định dạng đầu ra thống nhất cho Client:

### Phản hồi thành công (Success Response)
Sử dụng `response.OK(c, data)` hoặc `response.Created(c, data)`. Định dạng JSON đầu ra:
```json
{
  "success": true,
  "data": {
    "user_id": "uuid-123"
  }
}
```

### Phản hồi lỗi (Failure Response)
Sử dụng `response.Fail(c, statusCode, errorMsg)`. Định dạng JSON đầu ra:
```json
{
  "success": false,
  "error": "invalid credentials"
}
```

---

## 5. Quy Tắc Viết Unit Test

* **Tên tệp:** Tệp test phải có hậu tố `_test.go` (ví dụ: `usecase_test.go`).
* **Độc lập dữ liệu:** Không thực hiện kết nối cơ sở dữ liệu thật trong unit test. Sử dụng cấu trúc lưu trữ tạm (In-memory) hoặc Mock object (ví dụ: `mocks.NewMockUserRepository()`).
* **Phạm vi kiểm thử:**
  * **Domain:** Test các điều kiện ràng buộc dữ liệu đầu vào (ví dụ: định dạng email, độ dài mật khẩu).
  * **Application/Usecase:** Test luồng đi của dữ liệu, tính chính xác khi xử lý thành công/thất bại.
  * **HTTP Delivery:** Sử dụng `httptest.NewRecorder()` và `gin.TestMode` để mô phỏng Client gọi API mà không cần mở cổng mạng thật.

---

## 6. Quy Tắc Truyền Context (Context Propagation)

* **Tham số đầu tiên:** Mọi method trong lớp Usecase (Application) và Repository (Domain/Infrastructure) phải nhận tham số đầu tiên là `ctx context.Context`.
* **GORM Context:** Khi truy vấn cơ sở dữ liệu qua GORM, bắt buộc phải sử dụng `.WithContext(ctx)` trước các phương thức truy vấn để hỗ trợ ngắt kết nối khi timeout hoặc hủy request:
  ```go
  r.db.WithContext(ctx).First(&model, "id = ?", id)
  ```
* **Luồng truyền:** Context bắt đầu từ `c.Request.Context()` của Gin ở tầng Delivery, được truyền xuyên suốt qua Usecase xuống Repository.

---

## 7. Graceful Shutdown (Tắt Máy Chủ An Toàn)

* **Tín hiệu hệ điều hành:** Khi ứng dụng tắt (ví dụ do Kubernetes xoá Pod hoặc restart), máy chủ HTTP phải bắt các tín hiệu `SIGINT` và `SIGTERM` để thực hiện shutdown an toàn.
* **Thời gian chờ (Timeout):** Thiết lập một khoảng thời gian chờ (ví dụ: 5 giây) để máy chủ hoàn thành các request đang xử lý trước khi dừng hoàn toàn ứng dụng, tránh làm mất dữ liệu của người dùng.
* **Ngắt kết nối DB:** Sau khi dừng nhận request mới, đóng các kết nối đến Database hoặc Message Queue một cách an toàn.

---

## 8. Ghi Log Cấu Trúc (Structured Logging)

* **Sử dụng slog:** Sử dụng gói thư viện chuẩn `log/slog` thay vì gói `log` truyền thống để xuất log dưới dạng cấu trúc khoá-giá trị (key-value) dễ phân tích và lưu trữ.
* **Mức độ Log (Log Levels):** Sử dụng các mức độ log phù hợp:
  * `slog.Info`: Ghi nhận các sự kiện thông thường của hệ thống (ví dụ: khởi chạy server thành công, xử lý request thành công).
  * `slog.Warn`: Cảnh báo các sự cố không nghiêm trọng (ví dụ: request bị từ chối do dữ liệu sai định dạng).
  * `slog.Error`: Ghi nhận các lỗi nghiêm trọng ảnh hưởng đến hệ thống (ví dụ: lỗi kết nối cơ sở dữ liệu, lỗi Usecase execution).

---

## 9. Thư Mục Templates Mẫu (Template Scaffold Folders)

* **Tránh viết cứng mã nguồn mẫu:** Để tránh việc viết cứng (hardcode) chuỗi nội dung file mã nguồn mẫu trong code Go, chúng ta sử dụng thư mục `templates/` ở gốc dự án.
* **Cấu trúc thư mục mẫu:**
  ```text
  templates/
  ├── golang/          # Chứa main.go, go.mod, Dockerfile mẫu
  └── nodejs/          # Chứa index.js, package.json, Dockerfile mẫu
  ```
* **Luồng hoạt động & Tích hợp Git:**
  * Thư mục mẫu được lưu cố định trong ứng dụng để tránh truy vấn DB không cần thiết.
  * Khi gọi API tạo Service, hệ thống tự động sinh mã nguồn khung vào một thư mục tạm từ thư mục mẫu tương ứng.
  * Hệ thống sử dụng dịch vụ `VersionControl` để tạo Repository từ xa (ví dụ: trên GitHub thông qua GitHub API sử dụng Personal Access Token).
  * Mã nguồn khung cục bộ được tự động khởi tạo Git (`git init`), commit và push (`git push`) lên Repository từ xa, sau đó xoá sạch thư mục tạm.
  * URL của Repository mới được lưu lại vào cơ sở dữ liệu làm thông tin tham chiếu cho Service.

