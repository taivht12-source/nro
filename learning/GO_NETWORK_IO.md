# Đọc Dữ Liệu Mạng trong Go: `io.ReadFull` vs `binary.Read`

Trong lập trình mạng (Network Programming) với Go, việc đọc dữ liệu từ `net.Conn` (Socket) đòi hỏi sự chính xác vì dữ liệu truyền qua mạng là luồng (stream), có thể bị phân mảnh.

## 1. `io.ReadFull`

### Chức năng

Hàm `io.ReadFull(reader, buf)` đảm bảo **đọc đủ số lượng byte** để lấp đầy mảng `buf`.

- Nếu chưa đủ dữ liệu, nó sẽ chờ (block) cho đến khi đủ.
- Nếu kết nối bị đóng hoặc lỗi trước khi đọc đủ, nó sẽ trả về lỗi (ví dụ `io.EOF` hoặc `io.ErrUnexpectedEOF`).

### Tại sao cần dùng?

Khi bạn gọi `conn.Read(buf)`, nó có thể trả về số byte **ít hơn** kích thước của `buf` nếu gói tin bị chia nhỏ trên đường truyền. `io.ReadFull` giúp khắc phục điều này, đảm bảo bạn luôn có đủ dữ liệu mình cần.

### Ví dụ trong NRO Protocol

Đọc 1 byte Command hoặc đọc Header kích thước cố định.

```go
// Tạo buffer 1 byte để chứa Command
var cmdByte [1]byte

// Đọc chính xác 1 byte. Nếu mạng chậm, nó sẽ chờ.
// Nếu dùng conn.Read() thường, có rủi ro (dù nhỏ với 1 byte) nhưng io.ReadFull an toàn tuyệt đối.
_, err := io.ReadFull(s.Conn, cmdByte[:])
if err != nil {
    return // Lỗi hoặc mất kết nối
}
cmd := cmdByte[0]
```

## 2. `binary.Read`

### Chức năng

Hàm `binary.Read(reader, order, data)` dùng để đọc dữ liệu nhị phân có cấu trúc (như `int16`, `int32`, `float64`) từ luồng.

- **reader**: Nguồn đọc (ví dụ `s.Conn`).
- **order**: Quy tắc Endian (BigEndian hoặc LittleEndian). Giao thức mạng thường dùng **BigEndian**.
- **data**: Con trỏ đến biến cần lưu giá trị.

### Tại sao cần dùng?

Giúp code gọn gàng hơn khi bạn muốn đọc các kiểu số nguyên thay vì phải tự đọc mảng byte rồi dùng phép dịch bit (`<<`, `|`) để ghép lại.

### Ví dụ trong NRO Protocol

Đọc kích thước gói tin (Size) khi chưa Handshake (gửi dưới dạng `short` - 2 byte).

```go
var sizeShort int16

// Đọc 2 byte từ conn, tự động ghép lại thành số int16 theo chuẩn BigEndian
// Tương đương với: đọc 2 byte [b1, b2] -> size = (b1 << 8) | b2
err := binary.Read(s.Conn, binary.BigEndian, &sizeShort)
if err != nil {
    return
}
size := int(sizeShort)
```

## 3. So Sánh & Khi Nào Dùng

| Đặc điểm | `io.ReadFull` | `binary.Read` |
| :--- | :--- | :--- |
| **Mục đích** | Đọc mảng byte thô (Raw bytes) | Đọc số nguyên, struct (Structured data) |
| **Hiệu năng** | Rất nhanh, ít overhead | Chậm hơn chút do dùng Reflection |
| **Sử dụng khi** | Cần đọc dữ liệu đã mã hóa (cần giải mã từng byte), hoặc đọc chuỗi byte dài | Cần đọc nhanh các số Header chưa mã hóa, hoặc cấu trúc cố định |

### Áp dụng vào `session.go`

Trong `session.go`, chúng ta dùng kết hợp cả hai:

1. **Dùng `io.ReadFull` khi đã Handshake**:
    Vì dữ liệu bị mã hóa từng byte, ta không thể dùng `binary.Read` để đọc trực tiếp ra số đúng được. Ta phải đọc byte thô về, giải mã (XOR), rồi mới ghép lại.

    ```go
    // Đọc 2 byte size đã mã hóa
    var buf [1]byte
    io.ReadFull(s.Conn, buf[:]) // Đọc byte 1
    b1 := s.readKey(buf[0])     // Giải mã
    
    io.ReadFull(s.Conn, buf[:]) // Đọc byte 2
    b2 := s.readKey(buf[0])     // Giải mã
    
    size = (int(b1) << 8) | int(b2) // Ghép lại
    ```

2. **Dùng `binary.Read` khi chưa Handshake**:
    Lúc này dữ liệu chưa mã hóa, dùng `binary.Read` cho gọn code.

    ```go
    var sizeShort int16
    binary.Read(s.Conn, binary.BigEndian, &sizeShort)
    ```
