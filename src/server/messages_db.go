package server

import (
	"database/sql"
)

func PostMessage(categoryID int, userID int, content string, dbConn *sql.DB) (*Message, error) {
	var msg Message

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

	var username string
	_ = dbConn.QueryRow(`
		SELECT username FROM users WHERE id = $1
	`, userID).Scan(&username)

	msg.Username = username

	return &msg, nil
}

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
