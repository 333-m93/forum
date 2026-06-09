package database

import (
	"database/sql"
	"fmt"

	"forum.com/m/src/config"
	_ "github.com/lib/pq"
)

func Connect(cfg config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable", cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBName)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
