package server

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

func createUser(username, password string) (int, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	var id int

	err = dbConn.QueryRow(`
		INSERT INTO users (username, password_hash)
		VALUES ($1, $2)
		RETURNING id
	`, username, string(hash)).Scan(&id)

	return id, err
}

func authenticateUser(username, password string) (int, bool, error) {
	var userID int
	var hash string

	err := dbConn.QueryRow(`
		SELECT id, password_hash
		FROM users
		WHERE username = $1
	`, username).Scan(&userID, &hash)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, false, nil
		}
		return 0, false, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return 0, false, nil
	}

	return userID, true, nil
}
