package ports

import "nro-go/internal/core/domain"

// UserRepository định nghĩa các phương thức tương tác với dữ liệu User.
type UserRepository interface {
	GetByUsername(username string) (*domain.User, error)
	Create(user *domain.User) error
}

// PlayerRepository định nghĩa các phương thức tương tác với dữ liệu Player.
type PlayerRepository interface {
	GetByUserID(userID int) ([]*domain.Player, error)
	GetByID(id int) (*domain.Player, error)
	Create(player *domain.Player) error
	Update(player *domain.Player) error
}
