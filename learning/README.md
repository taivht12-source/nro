# Hướng Dẫn Học Code Game NRO (Dragon Ball Online)

Chào mừng bạn đến với mã nguồn game NRO (Ngọc Rồng Online). Tài liệu này sẽ giúp bạn hiểu cấu trúc dự án và cách bắt đầu lập trình.

## 1. Yêu Cầu Hệ Thống

Để chạy và phát triển dự án này, bạn cần:

- **Java Development Kit (JDK)**: Phiên bản 8 hoặc cao hơn (khuyên dùng JDK 8 hoặc 11 vì game cũ thường ổn định với bản này).
- **IDE**: NetBeans (khuyên dùng vì có folder `nbproject`), IntelliJ IDEA, hoặc Eclipse.
- **Database**: MySQL (XAMPP hoặc MySQL Server riêng).
- **Công cụ quản lý Database**: Navicat, phpMyAdmin, hoặc MySQL Workbench.

## 2. Cài Đặt & Chạy Server

### Bước 1: Cấu hình Database

1. Tạo database tên `nro` (hoặc tên khác tùy ý).
2. Import file `nro.sql` vào database vừa tạo.
3. Mở file `Config.properties` để cấu hình kết nối:

   ```properties
   database.host=localhost
   database.port=3306
   database.name=nro
   database.user=root
   database.pass=
   ```

### Bước 2: Cấu hình Server

Cũng trong `Config.properties`, bạn có thể chỉnh:

- `server.port`: Cổng server (mặc định 14445).
- `server.ip`: IP server (để local hoặc public).

### Bước 3: Chạy Server

- **Cách 1 (Dùng IDE)**: Mở project, tìm class `server.ServerManager`, chuột phải chọn **Run File** (hoặc chạy hàm `main`).
- **Cách 2 (Dùng file bat)**: Chạy file `run.bat` (lưu ý cần build ra file jar trong thư mục `dist` trước).

## 3. Cấu Trúc Dự Án (Source Code)

Source code nằm trong thư mục `src`. Dưới đây là các package quan trọng:

### `server` - Core của Server

- **`ServerManager.java`**: **Điểm bắt đầu (Entry Point)**. Chứa hàm `main`, khởi tạo server, kết nối database, và start các luồng xử lý (Boss, Shenron...).
- **`Controller.java`**: **Bộ xử lý trung tâm**. Nhận tất cả các gói tin (packet) từ Client gửi lên và điều hướng đến các chức năng tương ứng.
- **`Client.java`**: Quản lý danh sách người chơi đang online (`players`, `players_id`...).
- **`Manager.java`**: Quản lý dữ liệu tĩnh tải từ database (Map, Item Template, Skill Template...).

### `services` - Các dịch vụ (Logic game)

Chứa các class xử lý logic cụ thể và gửi gói tin về Client:

- **`Service.java`**: Các hàm gửi thông báo, cập nhật HP/MP, thông tin nhân vật cơ bản.
- **`ItemService.java`**: Xử lý liên quan đến vật phẩm.
- **`SkillService.java`**: Xử lý kỹ năng, chiêu thức.
- **`TaskService.java`**: Xử lý nhiệm vụ.
- **`ChatGlobalService.java`**: Xử lý chat thế giới.

### `player` - Đối tượng người chơi

- **`Player.java`**: Class đại diện cho một nhân vật. Chứa mọi thông tin: chỉ số, hành trang, đệ tử, nhiệm vụ...

### `database` - Kết nối dữ liệu

- Chứa các class DAO (Data Access Object) để đọc/ghi dữ liệu vào MySQL.

## 4. Luồng Hoạt Động Cơ Bản

### Khi Server Start

1. `ServerManager.main()` chạy.
2. `ServerManager.init()`: Load dữ liệu từ DB vào RAM (thông qua `Manager.gI()`).
3. `ServerManager.activeServerSocket()`: Mở cổng mạng lắng nghe kết nối.

### Khi Người Chơi Gửi Lệnh (Packet)

1. Client gửi gói tin (Message).
2. `Network` nhận và chuyển cho `Controller.onMessage(session, msg)`.
3. `Controller` đọc `msg.command` (mã lệnh) để biết người chơi muốn làm gì.
   - Ví dụ: `cmd = -7` là di chuyển.
4. `Controller` gọi đến Service tương ứng để xử lý logic (ví dụ `PlayerService.playerMove`).
5. Service gửi lại gói tin phản hồi cho Client (nếu cần).

## 5. Hướng Dẫn Code: Thêm Một Lệnh Chat Đơn Giản

Giả sử bạn muốn thêm lệnh chat "test" để nhận 1000 vàng.

**Bước 1**: Mở `server.Controller.java` hoặc nơi xử lý chat (thường là `ChatGlobalService` hoặc `Service`).
Tuy nhiên, logic chat thường nằm trong `Controller` case `-71` (Chat Global) hoặc `44` (Chat Command).

Tìm đến `case 44` trong `Controller.java`:

```java
case 44:
    if (player != null) {
        Command.gI().chat(player, _msg.reader().readUTF());
    }
    break;
```

Nó gọi sang `server.Command.java`. Hãy mở file đó.

**Bước 2**: Trong `Command.java`, thêm logic:

```java
public void chat(Player player, String text) {
    if (text.equals("test")) {
        player.inventory.gold += 1000;
        Service.gI().sendMoney(player); // Cập nhật tiền về client
        Service.gI().sendThongBao(player, "Bạn đã nhận được 1000 vàng!");
    }
}
```

## 6. Lời Khuyên Khi Học

1. **Đọc Log**: Khi chạy server, chú ý cửa sổ Console/Log. Nó sẽ báo lỗi hoặc thông tin quan trọng.
2. **Debug**: Dùng tính năng Debug của IDE để đặt Breakpoint và xem code chạy từng dòng.
3. **Backup**: Luôn backup database và source code trước khi sửa những thay đổi lớn.

## 7. Giao Tiếp Client - Server (Network)

Game sử dụng giao thức TCP thông qua `java.net.Socket`. Dữ liệu được đóng gói thành các **Message** (Gói tin).

### Cấu Trúc Gói Tin (`network.Message`)

Mỗi gói tin gửi đi hoặc nhận về đều là một đối tượng `Message`, bao gồm:

- **Command (`byte`)**: Mã lệnh (ví dụ: `-7` là di chuyển, `-1` là đăng nhập).
- **Data**: Dữ liệu đi kèm (đọc/ghi bằng `DataInputStream`/`DataOutputStream`).

### Quy Trình Gửi/Nhận

1. **Gửi tin (Server -> Client)**:
   - Tạo `Message msg = new Message(command)`.
   - Ghi dữ liệu: `msg.writer().writeInt(...)`, `msg.writer().writeUTF(...)`.
   - Gửi đi: `session.sendMessage(msg)`.
   - Code nằm trong `network.Session.sendMessage()` và `network.Sender`.

2. **Nhận tin (Client -> Server)**:
   - `network.Collector` liên tục đọc dữ liệu từ Socket.
   - Khi có dữ liệu, nó tạo `Message` và đẩy vào hàng đợi.
   - `network.QueueHandler` lấy message ra và gọi `Controller.onMessage()` để xử lý.

### Ví Dụ Code Gửi Tin

Xem trong `services.Service.java`, ví dụ hàm gửi thông báo:

```java
public void sendThongBao(Player player, String thongBao) {
    Message msg = null;
    try {
        msg = new Message(-26); // -26 là lệnh hiện thông báo
        msg.writer().writeUTF(thongBao);
        player.sendMessage(msg); // Gửi về client của player đó
    } catch (Exception e) {
        // Xử lý lỗi
    } finally {
        if (msg != null) {
            msg.cleanup();
            msg.dispose();
        }
    }
}
```
