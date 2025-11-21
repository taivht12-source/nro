# Kiến Trúc Hệ Thống NRO Game Server

Tài liệu này mô tả chi tiết kiến trúc của mã nguồn gốc (Java) và kiến trúc mục tiêu (Go) cho dự án NRO_BLACK_GOKU.

## 1. Tổng Quan

Dự án là một Game Server cho trò chơi Ngọc Rồng Online (NRO), sử dụng giao thức TCP để giao tiếp với Client. Server chịu trách nhiệm quản lý kết nối, logic game, và lưu trữ dữ liệu.

## 2. Kiến Trúc Mã Nguồn Gốc (Java)

Mã nguồn Java được tổ chức theo mô hình Monolithic cổ điển, tập trung vào các package chính trong thư mục `src`.

### 2.1. Cấu Trúc Package (`src`)

* **`server`**: Chứa các thành phần cốt lõi của server.
  * `ServerManager.java`: **Entry Point**. Khởi tạo server, kết nối DB, và quản lý vòng đời ứng dụng.
  * `Controller.java`: **Dispatcher**. Nhận gói tin từ Network và điều hướng đến Service xử lý.
  * `Client.java`: Quản lý danh sách người chơi online (`players`).
  * `Manager.java`: Quản lý dữ liệu tĩnh (Cache) tải từ DB (Map, ItemTemplate, SkillTemplate).
* **`services`**: Chứa logic nghiệp vụ (Business Logic).
  * `Service.java`: Các dịch vụ chung (Gửi thông báo, update HP/MP).
  * `ItemService.java`, `SkillService.java`, `TaskService.java`: Logic chuyên biệt.
* **`player`**: Domain Model.
  * `Player.java`: Entity quan trọng nhất, đại diện cho người chơi và chứa mọi trạng thái (Inventory, Stats...).
* **`network`**: Lớp giao tiếp mạng.
  * `Session.java`: Quản lý kết nối TCP, mã hóa/giải mã.
  * `Message.java`: Cấu trúc gói tin.
* **`database`**: Data Access Layer (DAO) để truy cập MySQL.

### 2.2. Luồng Xử Lý (Flow)

1. **Khởi động**: `ServerManager.main()` -> `init()` (Load data) -> `activeServerSocket()` (Listen port).
2. **Giao tiếp Client-Server**:
    * **Nhận tin**: Socket -> `Collector` -> `QueueHandler` -> `Controller.onMessage()` -> `Service`.
    * **Gửi tin**: `Service` -> `Session.sendMessage()` -> `Sender` -> Socket.
3. **Cơ chế Chat (Ví dụ)**:
    * Client gửi lệnh chat -> `Controller` (case 44) -> `Command.chat()` -> Xử lý logic (cộng tiền, thông báo) -> `Service.sendThongBao()`.

### 2.3. Giao Thức Mạng (Protocol)

* **Transport**: TCP/IP.
* **Packet Structure**: `[Command (byte)] + [Data (variable)]`.
* **Security**: XOR Cipher với Key động hoặc tĩnh ("NguyenDucVuEntertainment").

## 3. Kiến Trúc Mục Tiêu (Go)

Dự án Go (`go-src`) đang được xây dựng lại theo kiến trúc **Hexagonal (Ports & Adapters)** kết hợp với **CQRS** để đảm bảo tính mở rộng và dễ bảo trì.

### 3.1. Cấu Trúc Thư Mục (`go-src/src`)

* **`cmd/`**: Entry points (ví dụ: `server/main.go`).
* **`internal/`**: Mã nguồn riêng tư của dự án.
  * **`core/`**:
    * **`domain/`**: Các Entity thuần túy (`User`, `Player`, `Map`, `Item`). Không phụ thuộc framework.
    * **`ports/`**: Các Interface định nghĩa Input/Output (`UserRepository`, `MapRepository`).
  * **`app/`**: Application Layer.
    * **`commands/`**: Xử lý logic thay đổi trạng thái (CQRS Command) như `LoginCommand`.
    * **`queries/`**: Xử lý logic đọc dữ liệu (CQRS Query).
  * **`infrastructure/`**: Triển khai các Interface (Adapters).
    * **`network/`**: `TCPServer`, `Controller`, `Session` (Go implementation).
    * **`persistence/`**: `MySQLUserRepository`, `MySQLMapRepository`.
    * **`config/`**: Đọc file cấu hình.
* **`pkg/`**: Các thư viện dùng chung (Public).
  * **`protocol/`**: Định nghĩa `Message`, `Session`, thuật toán XOR.

### 3.2. Điểm Cải Tiến So Với Java

1. **Hexagonal Architecture**: Tách biệt hoàn toàn logic game (Core) khỏi cơ sở hạ tầng (DB, Network). Dễ dàng thay thế MySQL bằng MongoDB hoặc TCP bằng WebSocket mà không sửa Core.
2. **CQRS**: Tách biệt luồng Ghi (Command) và Đọc (Query), giúp tối ưu hóa hiệu năng.
3. **Concurrency**: Sử dụng Goroutines và Channels của Go thay vì Thread/Lock của Java, giúp xử lý hàng nghìn kết nối nhẹ nhàng hơn.
4. **Microservices Ready**: Cấu trúc module hóa sẵn sàng để tách thành các service nhỏ (Auth Service, World Service, Chat Service) giao tiếp qua gRPC.

## 4. Roadmap Chuyển Đổi

* **Phase 1**: Foundation & Network (Đã xong).
* **Phase 2**: Core Domain & Auth (Đã xong).
* **Phase 3**: World & Movement (Đang thực hiện).
* **Phase 4**: Combat System.
* **Phase 5**: Item & Inventory.
* **Phase 6**: NPC & Quests.
* **Phase 7**: Optimization & Release.
* **Phase 8**: Microservices & gRPC.
