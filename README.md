# Go Backend Template (Skeleton Framework)

Đây là một template project (skeleton framework) viết bằng ngôn ngữ Go (Golang) được thiết kế theo hướng **Clean/Layered Architecture**. Template này cung cấp cấu trúc thư mục rõ ràng, hệ thống quản lý vòng đời ứng dụng (Provider & Registry), cơ chế Dependency Injection (DI) thủ công, xử lý lỗi tập trung, cấu hình môi trường và tích hợp sẵn cơ chế xác thực JWT cùng ORM GORM.

Project này được thiết kế để làm khung sườn (framework) vững chắc và dễ mở rộng cho các dịch vụ backend Go trong tương lai. Hai nghiệp vụ `User` và `Auth` được tích hợp sẵn làm các ví dụ (examples) tiêu chuẩn.

---

## 1. Cấu trúc thư mục (Directory Structure)

Dự án được phân chia thành các thư mục con trong `internal` để đảm bảo tính đóng gói và độc lập giữa các tầng:

```text
├── cmd
│   └── main.go                  # Điểm khởi đầu của ứng dụng (Entry point)
├── internal
│   ├── app
│   │   └── app.go               # Quản lý vòng đời ứng dụng (App lifecycle)
│   ├── config
│   │   └── config.go            # Khởi tạo và load cấu hình từ file .env
│   ├── database
│   │   └── postgres.go          # Thiết lập kết nối cơ sở dữ liệu PostgreSQL (GORM)
│   ├── di
│   │   ├── container.go         # Struct Container chứa toàn bộ dependencies toàn cục
│   │   └── providers.go         # Hàm BuildContainer thực hiện wire dependencies thủ công
│   ├── dto                      # Data Transfer Objects (Đối tượng trao đổi dữ liệu)
│   │   ├── request              # Chứa struct request từ Client (ví dụ: login, create user)
│   │   └── response             # Chứa struct response chuẩn hóa gửi về Client
│   ├── entity                   # Định nghĩa các Database Models (GORM structs)
│   ├── exception                # Định nghĩa hệ thống Custom Errors (AppError)
│   ├── handler                  # Tầng Handler/Controller (tiếp nhận HTTP request từ Gin)
│   ├── logger                   # Wrapper logger cho Zap Logger
│   ├── mapper                   # Chuyển đổi dữ liệu (Mapping) giữa Entity và DTO
│   ├── middleware               # Chứa các HTTP middleware (Error handling, Auth, v.v.)
│   ├── provider                 # Interface quản lý vòng đời của các dịch vụ chạy nền
│   │   ├── interface.go         # Định nghĩa Provider interface (Boot/Close)
│   │   ├── registry.go          # Registry để quản lý, boot và close nhiều Provider đồng thời
│   │   └── http                 # HTTP Provider (Gin Server)
│   │       ├── server.go        # Khởi tạo HTTP server và lắng nghe kết nối
│   │       ├── router.go        # Setup middleware toàn cục và đăng ký các sub-router
│   │       ├── user_router.go   # Router cho nghiệp vụ User
│   │       └── auth_router.go   # Router cho nghiệp vụ Auth
│   ├── repository               # Tầng Repository (giao tiếp trực tiếp với Database qua GORM)
│   ├── security                 # Chứa các thành phần bảo mật (Bcrypt hashing, JWT provider)
│   ├── service                  # Tầng Service chứa Business Logic chính của ứng dụng
│   └── validation               # Parse và định dạng các lỗi validation (Binding validation)
```

---

## 2. Các thành phần kiến trúc cốt lõi

### A. Quản lý vòng đời dịch vụ (Provider & Registry)
Dự án sử dụng mô hình **Provider** để trừu tượng hóa các dịch vụ chạy nền (như HTTP Server, gRPC Server, Worker, Cron Job).
- **Interface Provider** (`internal/provider/interface.go`): Định nghĩa 3 phương thức `Name()`, `Boot()`, và `Close()`.
- **Registry** (`internal/provider/registry.go`): Đăng ký các Provider và sử dụng `errgroup` để khởi chạy bất đồng bộ đồng thời tất cả các Provider. Khi ứng dụng nhận được tín hiệu tắt (SIGINT/SIGTERM), Registry sẽ đóng các Provider theo thứ tự ngược lại (LIFO) để đảm bảo shutdown an toàn (Graceful Shutdown).

### B. Dependency Injection thủ công (Manual DI Container)
Để kiểm soát luồng phụ thuộc rõ ràng mà không phụ thuộc vào các thư viện DI phản chiếu (reflection), dự án sử dụng Dependency Injection bằng tay:
- **`Container`** (`internal/di/container.go`): Một struct chứa các thành phần dùng chung hoặc cần thiết để ứng dụng vận hành (Config, Logger, Database, Registry).
- **`BuildContainer()`** (`internal/di/providers.go`): Hàm duy nhất chịu trách nhiệm khởi tạo kết nối database, các thành phần bảo mật, khởi tạo Repository -> Service -> Handler -> HTTP Provider và cuối cùng đăng ký vào Registry.

### C. Định dạng Response và Xử lý lỗi tập trung
- **Response chuẩn hóa** (`internal/dto/response/base.go`): Tất cả API đều trả về một format JSON đồng nhất:
  - Khi thành công: `{"code": 0, "result": <dữ liệu>}`
  - Khi có lỗi: `{"code": <mã lỗi>, "message": "<nội dung lỗi>"}`
- **AppError** (`internal/exception/error.go`): Struct Custom Error triển khai interface `error` chuẩn của Go, bao gồm `Code`, `Message`, `HttpStatus` và lỗi gốc `Err`.
- **Error Middleware** (`internal/middleware/error.go`): Bất cứ lỗi nào phát sinh ở Service hoặc Handler chỉ cần được gom lại và đẩy vào Gin Context thông qua `c.Error(err)`. Middleware sẽ chặn lại ở cuối pipeline, phân tích loại lỗi và tự động trả về JSON tương ứng với HTTP Status Code phù hợp.

---

## 3. Quy trình phát triển nghiệp vụ mới (Workflow)

Để thêm một nghiệp vụ mới vào hệ thống (ví dụ: `Product`), bạn cần thực hiện theo các bước chuẩn mực dưới đây:

### Bước 1: Định nghĩa Entity
Tạo file model GORM tại `internal/entity/product.go` để đại diện cho bảng dữ liệu trong DB.
```go
package entity

import "time"

type Product struct {
	ID        string    `gorm:"primaryKey;type:varchar(36)"`
	Name      string    `gorm:"type:varchar(255);not null"`
	Price     float64   `gorm:"type:decimal(10,2);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
```

### Bước 2: Định nghĩa Requests & Responses (DTO)
Tạo các struct để bind dữ liệu từ Client gửi lên hoặc định dạng dữ liệu trả về:
- `internal/dto/request/product.go`:
```go
package request

type CreateProductRequest struct {
	Name  string  `json:"name" binding:"required,min=3"`
	Price float64 `json:"price" binding:"required,gt=0"`
}
```
- `internal/dto/response/product.go`:
```go
package response

type ProductResponse struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}
```

### Bước 3: Tạo Repository
Tạo interface và struct triển khai để tương tác với GORM tại `internal/repository/product.go`.
```go
package repository

import (
	"context"
	"identify/internal/entity"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(ctx context.Context, product *entity.Product) error
	FindByID(ctx context.Context, id string) (*entity.Product, error)
}

type ProductRepositoryImp struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &ProductRepositoryImp{db: db}
}

func (r *ProductRepositoryImp) Create(ctx context.Context, p *entity.Product) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *ProductRepositoryImp) FindByID(ctx context.Context, id string) (*entity.Product, error) {
	var p entity.Product
	if err := r.db.WithContext(ctx).First(&p, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}
```

### Bước 4: Tạo Service (Business Logic)
Định nghĩa Interface và Implement xử lý nghiệp vụ tại `internal/service/product.go`.
```go
package service

import (
	"context"
	"identify/internal/dto/request"
	"identify/internal/dto/response"
	"identify/internal/entity"
	"identify/internal/repository"
	"github.com/google/uuid"
)

type ProductService interface {
	CreateProduct(ctx context.Context, req request.CreateProductRequest) (*response.ProductResponse, error)
}

type ProductServiceImp struct {
	productRepo repository.ProductRepository
}

func NewProductService(productRepo repository.ProductRepository) ProductService {
	return &ProductServiceImp{productRepo: productRepo}
}

func (s *ProductServiceImp) CreateProduct(ctx context.Context, req request.CreateProductRequest) (*response.ProductResponse, error) {
	product := &entity.Product{
		ID:    uuid.New().String(),
		Name:  req.Name,
		Price: req.Price,
	}

	if err := s.productRepo.Create(ctx, product); err != nil {
		return nil, err
	}

	return &response.ProductResponse{
		ID:    product.ID,
		Name:  product.Name,
		Price: product.Price,
	}, nil
}
```

### Bước 5: Tạo Handler (Controller)
Viết handler tiếp nhận request HTTP tại `internal/handler/product.go`. Ở đây, ta sử dụng `validation.ParseValidationError(err)` để parse lỗi validate và chuyển lỗi sang middleware qua `c.Error(err)`.
```go
package handler

import (
	"identify/internal/dto/request"
	"identify/internal/dto/response"
	"identify/internal/service"
	"identify/internal/validation"
	"net/http"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

func (h *ProductHandler) Create(c *gin.Context) {
	var body request.CreateProductRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.Error(validation.ParseValidationError(err))
		return
	}

	res, err := h.productService.CreateProduct(c.Request.Context(), body)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, response.NewSuccessResponse(res))
}
```

### Bước 6: Định cấu hình Route
Tạo file router nghiệp vụ tại `internal/provider/http/product_router.go`.
```go
package http

import (
	"identify/internal/handler"
	"github.com/gin-gonic/gin"
)

type ProductRouter struct {
	productHandler *handler.ProductHandler
}

func NewProductRouter(h *handler.ProductHandler) *ProductRouter {
	return &ProductRouter{productHandler: h}
}

func (r *ProductRouter) Register(engine *gin.Engine) {
	api := engine.Group("/api/v1/products")
	{
		api.POST("", r.productHandler.Create)
	}
}
```
Sau đó đăng ký Router mới này vào `SetupRoutes` tại `internal/provider/http/router.go`:
```go
// Register sub-routers
NewUserRouter(userHandler).Register(engine)
NewAuthRouter(authHandler).Register(engine)
NewProductRouter(productHandler).Register(engine) // Đăng ký thêm tại đây
```

### Bước 7: Cập nhật DI Container
Đăng ký các phụ thuộc mới trong hàm `BuildContainer()` tại `internal/di/providers.go`:
1. Thêm `ProductHandler` làm tham số khởi tạo cho `NewHTTPProvider`.
2. Khởi tạo Repo, Service, Handler của Product trong luồng DI:
```diff
 	// 5. Build repositories
 	userRepo := repository.NewUserRepository(db)
+	productRepo := repository.NewProductRepository(db)
 
 	// 5.5 Build transaction manager
 	txManager := repository.NewTransactionManager(db)
 
 	// 6. Build services
 	userService := service.NewUserService(userRepo, passwordEncoder, txManager)
 	authService := service.NewAuthService(userRepo, passwordEncoder, jwtProvider)
+	productService := service.NewProductService(productRepo)
 
 	// 7. Build handlers
 	userHandler := handler.NewUserHandler(userService)
 	authHandler := handler.NewAuthHandler(authService)
+	productHandler := handler.NewProductHandler(productService)
 
 	// 8. Build HTTP provider
-	httpProvider := http.NewHTTPProvider(cfg.App, userHandler, authHandler, log)
+	httpProvider := http.NewHTTPProvider(cfg.App, userHandler, authHandler, productHandler, log)
```

---

## 4. Cách chạy dự án (How to Run)

### 1. Chuẩn bị môi trường
Copy file `.env.example` thành `.env` và cập nhật thông tin cấu hình (Database Connection, JWT Secrets, App Port):
```bash
cp .env.example .env
```

### 2. Khởi chạy Database Postgres
Đảm bảo bạn đã bật Postgres và tạo Database khớp với cấu hình trong `.env`.

### 3. Run server
Chạy lệnh bên dưới tại thư mục gốc để khởi động ứng dụng:
```bash
go run cmd/main.go
```
Khi ứng dụng khởi động thành công, console sẽ hiển thị log thông báo:
`all providers booted successfully` và lắng nghe kết nối HTTP trên cổng được cấu hình (mặc định là `:8080`).
