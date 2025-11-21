package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

// Message đại diện cho một gói tin trong giao thức NRO.
// Bao gồm mã lệnh (Command) và dữ liệu (Data).
type Message struct {
	Command int8         // Mã lệnh (ví dụ: -7 là di chuyển)
	Data    []byte       // Dữ liệu thô của gói tin
	reader  *bytes.Reader // Dùng để đọc dữ liệu từ Data
	writer  *bytes.Buffer // Dùng để ghi dữ liệu vào Data
}

// NewMessage tạo một Message mới với mã lệnh chỉ định.
func NewMessage(command int8) *Message {
	return &Message{
		Command: command,
		writer:  new(bytes.Buffer),
	}
}

// NewMessageFromData tạo một Message từ dữ liệu đã có (thường dùng khi nhận từ mạng).
func NewMessageFromData(command int8, data []byte) *Message {
	return &Message{
		Command: command,
		Data:    data,
		reader:  bytes.NewReader(data),
	}
}

// Cleanup giải phóng tài nguyên (trong Go có GC nên thường không cần, nhưng giữ lại cho giống logic cũ nếu cần pool).
func (m *Message) Cleanup() {
	m.Data = nil
	m.reader = nil
	m.writer = nil
}

// --- Các hàm ĐỌC dữ liệu (Reader) ---

// ReadByte đọc 1 byte từ gói tin.
func (m *Message) ReadByte() (int8, error) {
	if m.reader == nil {
		return 0, errors.New("reader is nil")
	}
	b, err := m.reader.ReadByte()
	return int8(b), err
}

// ReadShort đọc 2 byte (int16) từ gói tin.
func (m *Message) ReadShort() (int16, error) {
	if m.reader == nil {
		return 0, errors.New("reader is nil")
	}
	var v int16
	err := binary.Read(m.reader, binary.BigEndian, &v)
	return v, err
}

// ReadInt đọc 4 byte (int32) từ gói tin.
func (m *Message) ReadInt() (int32, error) {
	if m.reader == nil {
		return 0, errors.New("reader is nil")
	}
	var v int32
	err := binary.Read(m.reader, binary.BigEndian, &v)
	return v, err
}

// ReadUTF đọc chuỗi String (UTF-8) từ gói tin.
// Cấu trúc: [2 byte độ dài] + [bytes chuỗi].
func (m *Message) ReadUTF() (string, error) {
	if m.reader == nil {
		return "", errors.New("reader is nil")
	}
	// Đọc độ dài chuỗi (short)
	lenStr, err := m.ReadShort()
	if err != nil {
		return "", err
	}
	// Đọc các byte dữ liệu
	buf := make([]byte, lenStr)
	_, err = io.ReadFull(m.reader, buf)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

// --- Các hàm GHI dữ liệu (Writer) ---

// WriteByte ghi 1 byte vào gói tin.
func (m *Message) WriteByte(v int8) {
	m.writer.WriteByte(byte(v))
}

// WriteShort ghi 2 byte (int16) vào gói tin.
func (m *Message) WriteShort(v int16) {
	binary.Write(m.writer, binary.BigEndian, v)
}

// WriteInt ghi 4 byte (int32) vào gói tin.
func (m *Message) WriteInt(v int32) {
	binary.Write(m.writer, binary.BigEndian, v)
}

// WriteUTF ghi chuỗi String (UTF-8) vào gói tin.
func (m *Message) WriteUTF(v string) {
	lenStr := int16(len(v))
	m.WriteShort(lenStr)
	m.writer.WriteString(v)
}

// GetData trả về mảng byte dữ liệu đã ghi.
func (m *Message) GetData() []byte {
	if m.writer != nil {
		return m.writer.Bytes()
	}
	return m.Data
}
