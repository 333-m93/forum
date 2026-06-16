package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {

	dsn := os.Getenv("DATABASE_URL")

	if dsn == "" {
		log.Fatal("❌ DATABASE_URL manquant")
	}

	log.Println("🔌 Connecting DB...")

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("✅ DB connected")

	return db, nil
}

func Migrate(db *sql.DB) error {

	stmts := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(255) UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS sessions (
			id TEXT PRIMARY KEY,
			user_id INT NOT NULL,
			expires_at TIMESTAMP NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE INDEX IF NOT EXISTS idx_sessions_expires ON sessions (expires_at)`,
		`CREATE TABLE IF NOT EXISTS categories (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) UNIQUE NOT NULL,
			description TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS messages (
			id SERIAL PRIMARY KEY,
			category_id INT NOT NULL,
			user_id INT NOT NULL,
			content TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE INDEX IF NOT EXISTS idx_messages_category_time ON messages (category_id, created_at DESC)`,
		`INSERT INTO categories (name, description) VALUES
			('Chat général', 'Discussions générales sur tous sujets'),
			('MMA', 'Discussions et ressources sur MMA'),
			('Boxe Anglaise', 'Discussions et ressources sur Boxe Anglaise'),
			('Muay Thai', 'Discussions et ressources sur Muay Thai'),
			('Jujitsu Brésilien', 'Discussions et ressources sur Jujitsu Brésilien'),
			('Grappling', 'Discussions et ressources sur Grappling'),
			('Autres sports de combat', 'Discussions sur autres sports de combat')
			ON CONFLICT (name) DO NOTHING`,
	}

	for i, s := range stmts {
		if _, err := db.Exec(s); err != nil {
			log.Printf("⚠️ Migration step %d failed: %v", i+1, err)
			return err
		}
		log.Printf("  ✅ Migration step %d OK", i+1)
	}

	// Verify categories exist
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM categories").Scan(&count); err != nil {
		log.Printf("⚠️ Cannot count categories: %v", err)
		return err
	}
	log.Printf("✅ Migrations done — %d categories in DB", count)

	return nil
}
