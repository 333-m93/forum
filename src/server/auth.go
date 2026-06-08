package server

import (
	"database/sql"
	"errors"
	"html"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func registerHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		renderRegisterForm(w, "")
	case http.MethodPost:
		username := strings.TrimSpace(r.FormValue("username"))
		password := r.FormValue("password")

		if username == "" || password == "" {
			renderRegisterForm(w, "Veuillez remplir tous les champs.")
			return
		}

		userID, err := createUser(username, password)
		if err != nil {
			if strings.Contains(err.Error(), "Duplicate entry") {
				renderRegisterForm(w, "Ce nom d'utilisateur existe déjà.")
				return
			}
			http.Error(w, "Erreur serveur : "+html.EscapeString(err.Error()), http.StatusInternalServerError)
			return
		}

		// Créer une session et rediriger vers l'accueil
		sessionID, err := CreateSession(userID, dbConn)
		if err != nil {
			http.Error(w, "Erreur serveur", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			Path:     "/",
			MaxAge:   86400,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		renderLoginForm(w, "")
	case http.MethodPost:
		username := strings.TrimSpace(r.FormValue("username"))
		password := r.FormValue("password")

		if username == "" || password == "" {
			renderLoginForm(w, "Veuillez remplir tous les champs.")
			return
		}

		userID, ok, err := authenticateUser(username, password)
		if err != nil {
			http.Error(w, "Erreur serveur : "+html.EscapeString(err.Error()), http.StatusInternalServerError)
			return
		}
		if !ok {
			renderLoginForm(w, "Nom d'utilisateur ou mot de passe incorrect.")
			return
		}

		// Créer une session
		sessionID, err := CreateSession(userID, dbConn)
		if err != nil {
			http.Error(w, "Erreur serveur", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			Path:     "/",
			MaxAge:   86400,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
	}
}

func renderRegisterForm(w http.ResponseWriter, message string) {
	data := struct{ Message string }{Message: message}
	if err := authTemplates.ExecuteTemplate(w, "register.html", data); err != nil {
		http.Error(w, "Erreur template : "+html.EscapeString(err.Error()), http.StatusInternalServerError)
	}
}

func renderLoginForm(w http.ResponseWriter, message string) {
	data := struct{ Message string }{Message: message}
	if err := authTemplates.ExecuteTemplate(w, "login.html", data); err != nil {
		http.Error(w, "Erreur template : "+html.EscapeString(err.Error()), http.StatusInternalServerError)
	}
}

func createUser(username, password string) (int, error) {
	if dbConn == nil {
		return 0, errors.New("connexion à la base non initialisée")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	result, err := dbConn.Exec(`INSERT INTO users (username, password_hash) VALUES (?, ?)`, username, string(hash))
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	return int(id), err
}

func authenticateUser(username, password string) (int, bool, error) {
	if dbConn == nil {
		return 0, false, errors.New("connexion à la base non initialisée")
	}

	var userID int
	var hash string
	err := dbConn.QueryRow(`SELECT id, password_hash FROM users WHERE username = ?`, username).Scan(&userID, &hash)
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

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == nil {
		DestroySession(cookie.Value, dbConn)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
