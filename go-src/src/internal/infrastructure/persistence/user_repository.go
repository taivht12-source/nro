package persistence

import (
	"database/sql"
	"errors"
	"nro-go/internal/core/domain"
	"nro-go/internal/core/ports"
)

type MySQLUserRepository struct {
	db *sql.DB
}

func NewMySQLUserRepository(db *sql.DB) ports.UserRepository {
	return &MySQLUserRepository{db: db}
}

func (r *MySQLUserRepository) GetByUsername(username string) (*domain.User, error) {
	query := "SELECT id, username, password, role, ban, active FROM account WHERE username = ?"
	row := r.db.QueryRow(query, username)

	user := &domain.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Role, &user.Ban, &user.Active)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Không tìm thấy
		}
		return nil, err
	}
	return user, nil
}

func (r *MySQLUserRepository) Create(user *domain.User) error {
	return errors.New("not implemented")
}
