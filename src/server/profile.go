package server

import (
	"html/template"
	"strings"
	"time"
)

type ProfilePageData struct {
	User           *User
	JoinedAt       string
	MessageCount   int
	RecentMessages []RecentMessage
	Success        string
	Error          string
}

type RecentMessage struct {
	Content      string
	CategoryName string
	CreatedAt    time.Time
}

func buildProfileData(user *User, success, errMsg string) ProfilePageData {
	var messageCount int
	_ = dbConn.QueryRow(`SELECT COUNT(*) FROM messages WHERE user_id = $1`, user.ID).Scan(&messageCount)

	var joinedAt time.Time
	_ = dbConn.QueryRow(`SELECT created_at FROM users WHERE id = $1`, user.ID).Scan(&joinedAt)

	rows, err := dbConn.Query(`
		SELECT m.content, c.name, m.created_at
		FROM messages m
		JOIN categories c ON c.id = m.category_id
		WHERE m.user_id = $1
		ORDER BY m.created_at DESC
		LIMIT 5
	`, user.ID)
	if err != nil {
		rows = nil
	}

	var recent []RecentMessage
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var m RecentMessage
			if rows.Scan(&m.Content, &m.CategoryName, &m.CreatedAt) == nil {
				recent = append(recent, m)
			}
		}
	}

	return ProfilePageData{
		User:           user,
		JoinedAt:       joinedAt.Format("2 January 2006"),
		MessageCount:   messageCount,
		RecentMessages: recent,
		Success:        success,
		Error:          errMsg,
	}
}

func templateFuncs() template.FuncMap {
	return template.FuncMap{
		"uppercase": strings.ToUpper,
	}
}
