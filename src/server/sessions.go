package server

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"net/http"
	"time"
)

// User représente un utilisateur
type User struct {
	ID       int
	Username string
}

// GetSessionUser récupère l'utilisateur depuis un cookie de session
func GetSessionUser(r *http.Request, dbConn *sql.DB) (*User, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return nil, errors.New("no session")
	}

	sessionID := cookie.Value
	var userID int
	var expiresAt time.Time

	err = dbConn.QueryRow(`
		SELECT user_id, expires_at FROM sessions WHERE id = ?
	`, sessionID).Scan(&userID, &expiresAt)

	if err != nil {
		return nil, errors.New("invalid session")
	}

	if time.Now().After(expiresAt) {
		dbConn.Exec(`DELETE FROM sessions WHERE id = ?`, sessionID)
		return nil, errors.New("session expired")
	}

	var username string
	err = dbConn.QueryRow(`SELECT username FROM users WHERE id = ?`, userID).Scan(&username)
	if err != nil {
		return nil, err
	}

	return &User{ID: userID, Username: username}, nil
}

// CreateSession crée une nouvelle session pour un utilisateur
func CreateSession(userID int, dbConn *sql.DB) (string, error) {
	sessionID := make([]byte, 32)
	_, err := rand.Read(sessionID)
	if err != nil {
		return "", err
	}
	sessionIDStr := hex.EncodeToString(sessionID)

	expiresAt := time.Now().Add(24 * time.Hour)

	_, err = dbConn.Exec(`
		INSERT INTO sessions (id, user_id, expires_at) VALUES (?, ?, ?)
	`, sessionIDStr, userID, expiresAt)

	return sessionIDStr, err
}

// DestroySession supprime une session
func DestroySession(sessionID string, dbConn *sql.DB) error {
	_, err := dbConn.Exec(`DELETE FROM sessions WHERE id = ?`, sessionID)
	return err
}
