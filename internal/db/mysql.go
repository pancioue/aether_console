package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	Host string
	Port string
	User string
	Pass string
	Name string
}

func LoadConfig() Config {
	return Config{
		Host: getenv("DB_HOST", "127.0.0.1"),
		Port: getenv("DB_PORT", "3306"),
		User: getenv("DB_USER", "root"),
		Pass: getenv("DB_PASS", ""),
		Name: getenv("DB_NAME", "test"),
	}
}

func Open(cfg Config) (*sql.DB, error) {
	// parseTime=true: 讓 DATETIME/TIMESTAMP 能掃進 time.Time
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Local",
		cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.Name,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// 先給保守的 pool 值，之後再依 Cloud Run concurrency 調整
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}