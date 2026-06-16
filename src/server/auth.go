package server

import (
	"html"
	"net/http"
	"strings"
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
			if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
				renderRegisterForm(w, "Ce nom d'utilisateur existe déjà.")
				return
			}
			http.Error(w, "Erreur serveur : "+html.EscapeString(err.Error()), http.StatusInternalServerError)
			return
		}

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
			renderLoginForm(w, "Identifiants incorrects.")
			return
		}

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
	}
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

func renderRegisterForm(w http.ResponseWriter, message string) {
	data := struct {
		Message string
	}{
		Message: message,
	}
	_ = authTemplates.ExecuteTemplate(w, "register.html", data)
}

func renderLoginForm(w http.ResponseWriter, message string) {
	data := struct {
		Message string
	}{
		Message: message,
	}
	_ = authTemplates.ExecuteTemplate(w, "login.html", data)
}
