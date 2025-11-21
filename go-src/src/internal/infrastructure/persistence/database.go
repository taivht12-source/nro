package persistence

import (
	"database/sql"
	"fmt"
	"nro/src/internal/infrastructure/config"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// ConnectDB khởi tạo kết nối đến MySQL.
func ConnectDB() error {
	cfg := config.Get()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
		cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName)

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	// Cấu hình Connection Pool
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Hour)

	// Kiểm tra kết nối
	if err := db.Ping(); err != nil {
		return err
	}

	fmt.Println("Connected to Database successfully!")
	return nil
}

// GetDB trả về đối tượng DB connection.
func GetDB() *sql.DB {
	return db
}

// CloseDB đóng kết nối.
func CloseDB() {
	if db != nil {
		db.Close()
	}
}
