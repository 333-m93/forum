package server

import (
	"time"
)

type AdminUser struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	IsAdmin   bool      `json:"is_admin"`
	IsBanned  bool      `json:"is_banned"`
	MessageCount int    `json:"message_count"`
	CreatedAt time.Time `json:"created_at"`
}

type AdminMessage struct {
	ID           int       `json:"id"`
	Content      string    `json:"content"`
	Username     string    `json:"username"`
	CategoryName string    `json:"category_name"`
	CreatedAt    time.Time `json:"created_at"`
}

type AdminStats struct {
	UserCount    int `json:"user_count"`
	MessageCount int `json:"message_count"`
	CategoryCount int `json:"category_count"`
}

func GetAdminStats() (AdminStats, error) {
	var s AdminStats
	err := dbConn.QueryRow(`
		SELECT
			(SELECT COUNT(*) FROM users),
			(SELECT COUNT(*) FROM messages),
			(SELECT COUNT(*) FROM categories)
	`).Scan(&s.UserCount, &s.MessageCount, &s.CategoryCount)
	return s, err
}

func GetAllUsers() ([]AdminUser, error) {
	rows, err := dbConn.Query(`
		SELECT u.id, u.username, COALESCE(u.is_admin, FALSE), COALESCE(u.is_banned, FALSE),
		       COUNT(m.id) as msg_count, u.created_at
		FROM users u
		LEFT JOIN messages m ON m.user_id = u.id
		GROUP BY u.id, u.username, u.is_admin, u.is_banned, u.created_at
		ORDER BY u.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []AdminUser
	for rows.Next() {
		var u AdminUser
		if err := rows.Scan(&u.ID, &u.Username, &u.IsAdmin, &u.IsBanned, &u.MessageCount, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func GetAllMessages() ([]AdminMessage, error) {
	rows, err := dbConn.Query(`
		SELECT m.id, m.content, COALESCE(u.username, 'unknown'), c.name, m.created_at
		FROM messages m
		LEFT JOIN users u ON u.id = m.user_id
		JOIN categories c ON c.id = m.category_id
		ORDER BY m.created_at DESC
		LIMIT 200
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []AdminMessage
	for rows.Next() {
		var m AdminMessage
		if err := rows.Scan(&m.ID, &m.Content, &m.Username, &m.CategoryName, &m.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, rows.Err()
}

func BanUser(userID int) error {
	_, err := dbConn.Exec(`UPDATE users SET is_banned = TRUE WHERE id = $1`, userID)
	return err
}

func UnbanUser(userID int) error {
	_, err := dbConn.Exec(`UPDATE users SET is_banned = FALSE WHERE id = $1`, userID)
	return err
}

func DeleteMessageByID(messageID int) error {
	_, err := dbConn.Exec(`DELETE FROM messages WHERE id = $1`, messageID)
	return err
}
