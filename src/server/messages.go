package server

import (
	"database/sql"
	"time"
)

// =====================
// STRUCTS
// =====================

type Message struct {
	ID         int       `json:"id"`
	CategoryID int       `json:"category_id"`
	UserID     int       `json:"user_id"`
	Username   string    `json:"username"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// =====================
// CATEGORIES
// =====================

func GetCategories(dbConn *sql.DB) ([]Category, error) {

	rows, err := dbConn.Query(`
		SELECT id, name, description
		FROM categories
		ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category

	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Description); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	return categories, rows.Err()
}

// =====================
// CATEGORY BY NAME (auto-create if missing)
// =====================

func GetCategoryByName(name string, dbConn *sql.DB) (*Category, error) {

	var c Category

	err := dbConn.QueryRow(`
		SELECT id, name, description
		FROM categories
		WHERE LOWER(TRIM(name)) = LOWER(TRIM($1))
		LIMIT 1
	`, name).Scan(
		&c.ID,
		&c.Name,
		&c.Description,
	)

	if err == nil {
		return &c, nil
	}

	if err != sql.ErrNoRows {
		return nil, err
	}

	// Category not found — auto-create it
	err = dbConn.QueryRow(`
		INSERT INTO categories (name, description)
		VALUES ($1, $2)
		ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
		RETURNING id, name, description
	`, name, "Catégorie "+name).Scan(
		&c.ID,
		&c.Name,
		&c.Description,
	)

	if err != nil {
		return nil, err
	}

	return &c, nil
}

// =====================
// POST MESSAGE (FIX POSTGRES)
// =====================

func PostMessage(categoryID int, userID int, content string, dbConn *sql.DB) (*Message, error) {

	var msg Message

	// 🔥 INSERT + RETURNING (IMPORTANT POSTGRES)
	err := dbConn.QueryRow(`
		INSERT INTO messages (category_id, user_id, content)
		VALUES ($1, $2, $3)
		RETURNING id, category_id, user_id, content, created_at
	`, categoryID, userID, content).Scan(
		&msg.ID,
		&msg.CategoryID,
		&msg.UserID,
		&msg.Content,
		&msg.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	// 🔥 récupérer username
	var username string
	_ = dbConn.QueryRow(`
		SELECT username FROM users WHERE id = $1
	`, userID).Scan(&username)

	msg.Username = username

	return &msg, nil
}

// =====================
// GET MESSAGES
// =====================

func GetMessages(categoryID int, dbConn *sql.DB) ([]Message, error) {

	rows, err := dbConn.Query(`
		SELECT 
			m.id,
			m.category_id,
			m.user_id,
			COALESCE(u.username, 'unknown') AS username,
			m.content,
			m.created_at
		FROM messages m
		LEFT JOIN users u ON u.id = m.user_id
		WHERE m.category_id = $1
		ORDER BY m.created_at ASC
		LIMIT 50
	`, categoryID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message

	for rows.Next() {
		var msg Message

		err := rows.Scan(
			&msg.ID,
			&msg.CategoryID,
			&msg.UserID,
			&msg.Username,
			&msg.Content,
			&msg.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		messages = append(messages, msg)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}
