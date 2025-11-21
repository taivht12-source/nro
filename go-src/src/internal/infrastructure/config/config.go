package config

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// Config chứa toàn bộ cấu hình của Server.
type Config struct {
	ServerPort int
	DBHost     string
	DBPort     int
	DBUser     string
	DBPass     string
	DBName     string
	RequireDB  bool // Nếu false, server sẽ chạy mà không cần DB (cho testing)
}

var globalConfig *Config

// Get trả về cấu hình toàn cục.
func Get() *Config {
	if globalConfig == nil {
		// Fallback nếu chưa load (hoặc load mặc định)
		return &Config{
			ServerPort: 14445,
			DBHost:     "localhost",
			DBPort:     3306,
			DBUser:     "root",
			DBPass:     "",
			DBName:     "nro",
			RequireDB:  false, // Default: không yêu cầu DB
		}
	}
	return globalConfig
}

// Load đọc file cấu hình (thường là Config.properties hoặc .env).
// Ở đây ta sẽ parse file Config.properties của server cũ để tương thích.
func Load(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	cfg := &Config{
		ServerPort: 14445, // Default
		DBHost:     "localhost",
		DBPort:     3306,
		DBUser:     "root",
		DBPass:     "",
		DBName:     "nro",
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		switch key {
		case "server.port":
			if p, err := strconv.Atoi(val); err == nil {
				cfg.ServerPort = p
			}
		case "database.host":
			cfg.DBHost = val
		case "database.port":
			if p, err := strconv.Atoi(val); err == nil {
				cfg.DBPort = p
			}
		case "database.user":
			cfg.DBUser = val
		case "database.pass":
			cfg.DBPass = val
		case "database.name":
			cfg.DBName = val
		case "database.required":
			cfg.RequireDB = (val == "true" || val == "1")
		}
	}

	globalConfig = cfg
	return scanner.Err()
}
