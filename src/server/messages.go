package server

import (
	"database/sql"
	"time"
)

type Message struct {
	ID         int       `json:"id"`
	CategoryID int       `json:"category_id"`
	UserID     int       `json:"user_id"`
	Username   string    `json:"username"`
	AvatarURL  string    `json:"avatar_url"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

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

func GetCategoryByName(name string, dbConn *sql.DB) (*Category, error) {
	var c Category

	err := dbConn.QueryRow(`
		SELECT id, name, description
		FROM categories
		WHERE LOWER(TRIM(name)) = LOWER(TRIM($1))
		LIMIT 1
	`, name).Scan(&c.ID, &c.Name, &c.Description)

	if err == nil {
		return &c, nil
	}

	if err != sql.ErrNoRows {
		return nil, err
	}

	err = dbConn.QueryRow(`
		INSERT INTO categories (name, description)
		VALUES ($1, $2)
		ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
		RETURNING id, name, description
	`, name, "Catégorie "+name).Scan(&c.ID, &c.Name, &c.Description)

	if err != nil {
		return nil, err
	}

	return &c, nil
}
