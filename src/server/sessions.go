package server

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"net/http"
	"time"
)

// User
type User struct {
	ID       int
	Username string
}

// ==========================
// GET SESSION USER
// ==========================
func GetSessionUser(r *http.Request, dbConn *sql.DB) (*User, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return nil, errors.New("no session")
	}

	sessionID := cookie.Value

	var userID int
	var expiresAt time.Time

	err = dbConn.QueryRow(`
		SELECT user_id, expires_at
		FROM sessions
		WHERE id = $1
	`, sessionID).Scan(&userID, &expiresAt)

	if err != nil {
		return nil, errors.New("invalid session")
	}

	if time.Now().After(expiresAt) {
		_, _ = dbConn.Exec(`
			DELETE FROM sessions
			WHERE id = $1
		`, sessionID)

		return nil, errors.New("session expired")
	}

	var username string
	err = dbConn.QueryRow(`
		SELECT username
		FROM users
		WHERE id = $1
	`, userID).Scan(&username)

	if err != nil {
		return nil, err
	}

	return &User{
		ID:       userID,
		Username: username,
	}, nil
}

// ==========================
// CREATE SESSION
// ==========================
func CreateSession(userID int, dbConn *sql.DB) (string, error) {
	sessionIDBytes := make([]byte, 32)
	_, err := rand.Read(sessionIDBytes)
	if err != nil {
		return "", err
	}

	sessionID := hex.EncodeToString(sessionIDBytes)
	expiresAt := time.Now().Add(24 * time.Hour)

	_, err = dbConn.Exec(`
		INSERT INTO sessions (id, user_id, expires_at)
		VALUES ($1, $2, $3)
	`, sessionID, userID, expiresAt)

	return sessionID, err
}

// ==========================
// DESTROY SESSION
// ==========================
func DestroySession(sessionID string, dbConn *sql.DB) error {
	_, err := dbConn.Exec(`
		DELETE FROM sessions
		WHERE id = $1
	`, sessionID)

	return err
}
