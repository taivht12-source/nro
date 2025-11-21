package protocol

import (
	"encoding/binary"
	"io"
	"net"
	"sync"
)

// Session quản lý kết nối mạng với một Client.
// Nó đóng vai trò như một lớp trung gian giữa socket (net.Conn) và logic game.
//
// Chức năng chính:
// 1. Quản lý kết nối (Connect/Disconnect).
// 2. Đọc gói tin từ Client (Read Loop).
// 3. Gửi gói tin đến Client (Write Loop).
// 4. Mã hóa và Giải mã dữ liệu sử dụng thuật toán XOR với một khóa (Key).
//
// Ví dụ sử dụng:
//
//	conn, _ := listener.Accept()
//	session := protocol.NewSession(conn, 1)
//	session.Start() // Bắt đầu luồng đọc/ghi
//
//	// Gửi tin nhắn
//	msg := protocol.NewMessage(-7)
//	msg.WriteShort(100)
//	session.SendMessage(msg)
type Session struct {
	ID        int            // ID định danh của Session (thường là ID người chơi hoặc số thứ tự kết nối)
	Conn      net.Conn       // Đối tượng kết nối TCP thực tế
	Key       []byte         // Khóa mã hóa (Session Key). Dùng để XOR dữ liệu.
	CurR      int            // Con trỏ đọc: vị trí hiện tại trong mảng Key khi giải mã dữ liệu nhận được.
	CurW      int            // Con trỏ ghi: vị trí hiện tại trong mảng Key khi mã hóa dữ liệu gửi đi.
	SentKey   bool           // Cờ đánh dấu trạng thái Handshake. True nghĩa là đã trao đổi Key xong và bắt đầu mã hóa.
	Connected bool           // Trạng thái kết nối.
	SendChan  chan *Message  // Kênh (Channel) chứa các Message chờ gửi đi. Giúp việc gửi tin nhắn không chặn luồng chính.
	Player    interface{}    // Người chơi gắn với Session
	UserID    int            // ID của user đã đăng nhập (sau khi login thành công)
	mu        sync.Mutex     // Mutex để đảm bảo an toàn khi truy cập tài nguyên chia sẻ (như đóng kết nối).
	Handler   MessageHandler // Interface xử lý tin nhắn
}

type MessageHandler interface {
	OnMessage(session *Session, msg *Message)
}

// NewSession tạo một Session mới từ kết nối net.Conn.
// Khởi tạo các giá trị mặc định và channel gửi tin.
func NewSession(conn net.Conn, id int) *Session {
	return &Session{
		ID:        id,
		Conn:      conn,
		Key:       nil, // Ban đầu chưa có Key, sẽ được set sau khi Handshake thành công.
		Connected: true,
		SendChan:  make(chan *Message, 100), // Buffer 100 tin nhắn để tránh bị block nếu mạng chậm.
	}
}

// SetKey thiết lập khóa mã hóa cho Session.
// Được gọi khi Server quyết định Key (thường là mặc định hoặc random) và bắt đầu phiên mã hóa.
func (s *Session) SetKey(key []byte) {
	s.Key = key
	s.CurR = 0 // Reset con trỏ đọc về đầu mảng Key
	s.CurW = 0 // Reset con trỏ ghi về đầu mảng Key
}

// Start bắt đầu 2 goroutine riêng biệt cho việc Đọc và Gửi tin nhắn.
// Điều này giúp việc nhập/xuất dữ liệu không ảnh hưởng lẫn nhau.
func (s *Session) Start() {
	go s.writeLoop() // Goroutine chuyên gửi tin
	go s.readLoop()  // Goroutine chuyên đọc tin
}

// SendMessage đưa một gói tin vào hàng đợi gửi (SendChan).
// Hàm này không gửi ngay lập tức mà chỉ đẩy vào channel để writeLoop xử lý.
// Đây là cách xử lý bất đồng bộ (asynchronous) hiệu quả trong Go.
func (s *Session) SendMessage(msg *Message) {
	if s.Connected {
		s.SendChan <- msg
	}
}

// Close đóng kết nối và giải phóng tài nguyên.
// Đảm bảo chỉ đóng 1 lần duy nhất nhờ Mutex.
func (s *Session) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.Connected {
		s.Connected = false
		s.Conn.Close()    // Đóng socket
		close(s.SendChan) // Đóng channel để dừng writeLoop
	}
}

// --- Xử lý Mã Hóa (XOR Cipher) ---
// Cơ chế mã hóa của NRO là XOR từng byte dữ liệu với byte tương ứng trong chuỗi Key.
// Key được lặp lại vòng tròn (circular buffer).

// readKey giải mã một byte nhận được từ Client.
// Công thức: decoded_byte = (key[curR] & 0xFF) ^ (received_byte & 0xFF)
func (s *Session) readKey(b byte) byte {
	if len(s.Key) == 0 {
		return b // Nếu chưa có Key thì không giải mã
	}
	res := (s.Key[s.CurR] & 0xFF) ^ (b & 0xFF)
	s.CurR++
	if s.CurR >= len(s.Key) {
		s.CurR %= len(s.Key) // Quay vòng về đầu mảng Key nếu hết
	}
	return res
}

// writeKey mã hóa một byte trước khi gửi cho Client.
// Công thức: encoded_byte = (key[curW] & 0xFF) ^ (raw_byte & 0xFF)
func (s *Session) writeKey(b byte) byte {
	if len(s.Key) == 0 {
		return b // Nếu chưa có Key thì không mã hóa
	}
	res := (s.Key[s.CurW] & 0xFF) ^ (b & 0xFF)
	s.CurW++
	if s.CurW >= len(s.Key) {
		s.CurW %= len(s.Key) // Quay vòng về đầu mảng Key nếu hết
	}
	return res
}

// --- Vòng lặp Đọc/Ghi (Core Network Loop) ---

// readLoop liên tục đọc dữ liệu từ Socket, giải mã và tái tạo thành Message.
func (s *Session) readLoop() {
	defer s.Close() // Đảm bảo đóng kết nối khi vòng lặp kết thúc (do lỗi hoặc disconnect)
	for {
		// 1. Đọc Command (1 byte đầu tiên)
		// Command định danh loại gói tin (ví dụ: -7 là di chuyển).
		var cmdByte [1]byte
		_, err := io.ReadFull(s.Conn, cmdByte[:])
		if err != nil {
			return // Lỗi đọc (thường là client ngắt kết nối) -> thoát
		}

		cmd := cmdByte[0]
		// Nếu đã Handshake, byte Command này đã bị mã hóa -> cần giải mã.
		if s.SentKey {
			cmd = s.readKey(cmd)
		}

		// 2. Đọc Size (Kích thước dữ liệu đi kèm)
		// Cách đọc Size phụ thuộc vào trạng thái Handshake.
		var size int
		if s.SentKey {
			// Logic đọc size khi ĐÃ Handshake (theo giao thức NRO):
			// Size được gửi dưới dạng 2 byte, mỗi byte đều bị mã hóa.
			// Cần đọc 2 byte -> giải mã từng byte -> ghép lại thành số nguyên (Big Endian).
			var b1, b2 byte
			var buf [1]byte

			io.ReadFull(s.Conn, buf[:])
			b1 = buf[0]
			io.ReadFull(s.Conn, buf[:])
			b2 = buf[0]

			b1 = s.readKey(b1)
			b2 = s.readKey(b2)

			// Ghép 2 byte thành số int: (b1 << 8) | b2
			// Giải thích:
			// - b1 là byte cao (High Byte), b2 là byte thấp (Low Byte) theo quy tắc Big Endian.
			// - (b1 & 0xFF): Lấy giá trị unsigned của b1.
			// - << 8: Dịch trái 8 bit (tương đương nhân 256) để đưa b1 lên hàng "trăm" (trong hệ 16-bit).
			// - | b2: Cộng thêm giá trị của b2 vào.
			// Ví dụ: Size = 300 (00000001 00101100) -> b1=1, b2=44 -> (1<<8) | 44 = 256 + 44 = 300.
			size = (int(b1&0xFF) << 8) | int(b2&0xFF)
		} else {
			// Logic đọc size khi CHƯA Handshake:
			// Size được gửi dưới dạng Short (2 byte) chuẩn, không mã hóa.
			var sizeShort int16
			err := binary.Read(s.Conn, binary.BigEndian, &sizeShort)
			if err != nil {
				return
			}
			size = int(sizeShort)
		}

		// 3. Đọc Data (Dữ liệu chính của gói tin)
		data := make([]byte, size)
		if size > 0 {
			_, err := io.ReadFull(s.Conn, data)
			if err != nil {
				return
			}
			// Nếu đã Handshake, toàn bộ phần Data cũng bị mã hóa -> cần giải mã từng byte.
			if s.SentKey {
				for i := 0; i < len(data); i++ {
					data[i] = s.readKey(data[i])
				}
			}
		}

		// 4. Xử lý gói tin
		// Tại đây ta đã có một Message hoàn chỉnh (Command + Data sạch).
		if s.Handler != nil {
			s.Handler.OnMessage(s, NewMessageFromData(int8(cmd), data))
		} else {
			// Tạm thời in ra console để debug nếu chưa có Handler
			// fmt.Printf("Recv CMD: %d, Size: %d\n", int8(cmd), size)
		}

		// Xử lý Handshake đặc biệt (-27: GET_SESSION_ID):
		// Đây là gói tin đầu tiên Client gửi để yêu cầu kết nối.
		if int8(cmd) == -27 {
			s.sendSessionID()
		}
	}
}

// writeLoop liên tục lấy Message từ channel và gửi qua Socket.
func (s *Session) writeLoop() {
	for msg := range s.SendChan {
		s.doSendMessage(msg)
	}
}

// doSendMessage thực hiện việc mã hóa và ghi byte ra Socket.
func (s *Session) doSendMessage(msg *Message) error {
	data := msg.GetData()
	cmd := byte(msg.Command)

	// Buffer để gom dữ liệu ghi một lần (giảm syscall)
	// Cấu trúc: [Command] [Size] [Data]

	// 1. Ghi Command
	// Nếu đã Handshake, mã hóa byte Command trước khi gửi.
	if s.SentKey {
		cmd = s.writeKey(cmd)
	}
	s.Conn.Write([]byte{cmd})

	// 2. Ghi Size
	size := len(data)
	if s.SentKey {
		// Logic ghi size khi ĐÃ Handshake:
		// Tách size thành 2 byte (Big Endian) -> Mã hóa từng byte -> Gửi.
		b1 := byte(size >> 8)
		b2 := byte(size & 0xFF)

		b1 = s.writeKey(b1)
		b2 = s.writeKey(b2)

		s.Conn.Write([]byte{b1, b2})
	} else {
		// Logic ghi size khi CHƯA Handshake:
		// Gửi trực tiếp 2 byte Short (Big Endian).
		binary.Write(s.Conn, binary.BigEndian, int16(size))
	}

	// 3. Ghi Data
	if size > 0 {
		if s.SentKey {
			// Nếu đã Handshake, phải mã hóa từng byte Data.
			// Lưu ý: Phải copy data ra mảng mới để không làm hỏng dữ liệu gốc của Message
			// (vì Message có thể được dùng lại hoặc log).
			encodedData := make([]byte, size)
			copy(encodedData, data)
			for i := 0; i < size; i++ {
				encodedData[i] = s.writeKey(encodedData[i])
			}
			s.Conn.Write(encodedData)
		} else {
			// Chưa Handshake -> Gửi raw data.
			s.Conn.Write(data)
		}
	}

	return nil
}

// sendSessionID xử lý quy trình Handshake ban đầu.
func (s *Session) sendSessionID() {
	// Gửi gói tin -27 (GET_SESSION_ID) về Client.
	// Byte 1: Trạng thái kết nối (1 = Thành công).
	msg := NewMessage(-27)
	msg.WriteByte(1)
	s.SendMessage(msg)

	// Sau khi gửi gói tin này, Server thiết lập Key mặc định.
	// Từ gói tin tiếp theo trở đi, mọi dữ liệu sẽ được mã hóa.
	if s.Key == nil {
		s.SetKey([]byte("NguyenDucVuEntertainment")) // Key mặc định của NRO
		s.SentKey = true                             // Kích hoạt chế độ mã hóa
	}
}
