package commands

import (
	"errors"
	"fmt"
	"nro-go/internal/core/domain"
	"nro-go/internal/core/ports"
)

// LoginCommand chứa thông tin yêu cầu đăng nhập.
type LoginCommand struct {
	Username string
	Password string
}

// LoginHandler xử lý logic đăng nhập.
type LoginHandler struct {
	userRepo ports.UserRepository
}

func NewLoginHandler(userRepo ports.UserRepository) *LoginHandler {
	return &LoginHandler{userRepo: userRepo}
}

// Handle thực thi lệnh đăng nhập.
func (h *LoginHandler) Handle(cmd LoginCommand) (*domain.User, error) {
	// 1. Tìm user trong DB
	user, err := h.userRepo.GetByUsername(cmd.Username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("tài khoản không tồn tại")
	}

	// 2. Kiểm tra mật khẩu (Lưu ý: Server cũ có thể dùng MD5 hoặc Plaintext, cần check lại logic cũ)
	// Tạm thời so sánh trực tiếp (Plaintext) cho giống logic cơ bản, sau này sẽ thêm Hash.
	if user.Password != cmd.Password {
		return nil, errors.New("mật khẩu không chính xác")
	}

	// 3. Kiểm tra Ban
	if user.Ban > 0 {
		return nil, errors.New("tài khoản đã bị khóa")
	}

	fmt.Printf("User %s logged in successfully\n", user.Username)
	return user, nil
}
