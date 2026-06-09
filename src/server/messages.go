package server

import (
	"database/sql"
	"time"
)

// Message représente un message du forum
type Message struct {
	ID         int       `json:"id"`
	CategoryID int       `json:"category_id"`
	UserID     int       `json:"user_id"`
	Username   string    `json:"username"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}

// Category représente une catégorie de chat
type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// GetCategories
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
		var cat Category
		if err := rows.Scan(&cat.ID, &cat.Name, &cat.Description); err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}

	return categories, rows.Err()
}

// GetCategoryByName (POSTGRES FIX)
func GetCategoryByName(name string, dbConn *sql.DB) (*Category, error) {
	var cat Category

	err := dbConn.QueryRow(`
		SELECT id, name, description
		FROM categories
		WHERE name = $1
	`, name).Scan(&cat.ID, &cat.Name, &cat.Description)

	if err != nil {
		return nil, err
	}

	return &cat, nil
}

// PostMessage (POSTGRES FIX)
func PostMessage(categoryID int, userID int, content string, dbConn *sql.DB) (*Message, error) {
	result, err := dbConn.Exec(`
		INSERT INTO messages (category_id, user_id, content)
		VALUES ($1, $2, $3)
	`, categoryID, userID, content)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	var username string
	_ = dbConn.QueryRow(`
		SELECT username FROM users WHERE id = $1
	`, userID).Scan(&username)

	return &Message{
		ID:         int(id),
		CategoryID: categoryID,
		UserID:     userID,
		Username:   username,
		Content:    content,
		CreatedAt:  time.Now(),
	}, nil
}

// GetMessages (POSTGRES FIX)
func GetMessages(categoryID int, dbConn *sql.DB) ([]Message, error) {
	rows, err := dbConn.Query(`
		SELECT m.id, m.category_id, m.user_id, u.username, m.content, m.created_at
		FROM messages m
		JOIN users u ON m.user_id = u.id
		WHERE m.category_id = $1
		ORDER BY m.created_at DESC
		LIMIT 50
	`, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(
			&msg.ID,
			&msg.CategoryID,
			&msg.UserID,
			&msg.Username,
			&msg.Content,
			&msg.CreatedAt,
		); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	// reverse order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, rows.Err()
}
