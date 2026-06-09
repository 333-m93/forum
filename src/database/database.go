package database

import (
	"database/sql"
	"fmt"

	"forum.com/m/src/config"
	_ "github.com/lib/pq"
)

func Connect(cfg config.Config) (*sql.DB, error) {

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s sslmode=require",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBName,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// test connexion réelle
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
