# Giao Thức Mạng NRO (Network Protocol)

Tài liệu này mô tả chi tiết cách Client và Server giao tiếp ở mức byte (byte-level).

## 1. Tổng Quan

- **Giao thức**: TCP/IP.
- **Cấu trúc**: Header + Body.
- **Mã hóa**: XOR đơn giản với khóa (Session Key).

## 2. Cấu Trúc Gói Tin (Packet Structure)

Mỗi gói tin gửi đi có dạng:
`[Command] [Size] [Data]`

### 2.1. Command (Lệnh)

- **Kích thước**: 1 byte.
- **Mô tả**: Định danh loại hành động (ví dụ: -7 là di chuyển).
- **Xử lý**:
  - Nếu chưa handshake (chưa gửi key): Gửi raw byte.
  - Nếu đã handshake: `byte_gui = command ^ key[index]`.

### 2.2. Size (Kích thước dữ liệu)

Quy định độ dài của phần Data. Cách ghi size phụ thuộc vào trạng thái kết nối và độ lớn của dữ liệu.

**Trường hợp 1: Chưa Handshake (`sentKey = false`)**

- Dùng 2 byte (Short) để ghi độ dài.
- `dos.writeShort(size)`

**Trường hợp 2: Đã Handshake (`sentKey = true`)**

- **Nếu là các lệnh đặc biệt** (-32, -66, -74, 11, -67, -87, 66):
  - Ghi 4 byte.
  - Byte 1: `writeKey(size)`
  - Byte 2: `size - 128`
  - Byte 3: `writeKey(size >> 8)`
  - Byte 4: `size >> 16` (có vẻ logic này trong code hơi lạ, cần check kỹ lại `MessageSendCollect.java` nếu muốn implement client riêng).
- **Nếu size > 0**:
  - Ghi 2 byte.
  - Byte 1: `writeKey(size >> 8)`
  - Byte 2: `writeKey(size & 0xFF)`

### 2.3. Data (Dữ liệu)

- Là mảng byte chứa nội dung chính.
- Nếu đã handshake, từng byte trong data sẽ được XOR với key:
  `data[i] = data[i] ^ key[current_index]`

## 3. Cơ Chế Mã Hóa (XOR Cipher)

Server và Client chia sẻ một chuỗi khóa (Key). Mặc định trong source này là:
`KEYS = "NguyenDucVuEntertainment".getBytes()`

Khi đọc/ghi:

1. Duy trì một con trỏ `curR` (khi đọc) và `curW` (khi ghi).
2. Byte thực tế = `Byte mạng ^ Key[cur % Key.length]`.
3. Tăng con trỏ sau mỗi lần xử lý 1 byte.

## 4. Quy Trình Kết Nối (Handshake)

1. **Client kết nối** đến Server (Socket connect).
2. **Server** (trong `Session`):
   - Khởi tạo `KEYS` mặc định.
   - Gửi `GET_SESSION_ID` (-27) kèm theo thông tin (nếu có).
3. **Client**:
   - Nhận gói tin -27.
   - Thiết lập Key (nếu server có gửi key mới, hoặc dùng key mặc định).
   - Đánh dấu `sentKey = true`.
4. Từ lúc này, mọi gói tin đều được mã hóa XOR.

## 5. Danh Sách Lệnh Quan Trọng (`consts.Cmd_message`)

| Command ID | Tên Constant | Mô tả |
| :--- | :--- | :--- |
| -7 | `PLAYER_MOVE` | Di chuyển nhân vật |
| -1 | `LOGIN` | Đăng nhập |
| -127 | `LUCKY_ROUND` | Vòng quay may mắn |
| -27 | `GET_SESSION_ID` | Handshake ban đầu |
| -26 | `DIALOG_MESSAGE` | Hiện thông báo hội thoại |
| -25 | `SERVER_MESSAGE` | Thông báo server |
| 44 | `CHAT_MAP` | Chat trong map |
| -71 | `CHAT_THEGIOI_CLIENT` | Chat thế giới |

## 6. Ví Dụ Đọc Gói Tin Di Chuyển (-7)

Khi Server nhận gói tin -7:

1. `Controller.onMessage` nhận được `Message` với `command = -7`.
2. Đọc dữ liệu:

   ```java
   try {
       byte b = msg.reader().readByte(); // Đọc 1 byte (thường là status)
       short x = msg.reader().readShort(); // Tọa độ X
       short y = msg.reader().readShort(); // Tọa độ Y
       // Xử lý logic di chuyển...
   } catch (IOException e) {
       e.printStackTrace();
   }
   ```

## 7. Lưu Ý Khi Debug

- Nếu bạn dùng Wireshark để bắt gói tin, bạn sẽ thấy dữ liệu bị mã hóa (trông lộn xộn) sau bước handshake.
- Để debug, hãy in log tại `MessageSendCollect.readMessage` trước khi XOR để thấy byte gốc.
