package server

import (
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func profileHandler(w http.ResponseWriter, r *http.Request) {
	user, err := GetSessionUser(r, dbConn)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		updateProfile(w, r, user)
		return
	}

	data := buildProfileData(user, "", "")
	_ = authTemplates.ExecuteTemplate(w, "profile.html", data)
}

func updateProfile(w http.ResponseWriter, r *http.Request, user *User) {
	username := strings.TrimSpace(r.FormValue("username"))
	bio := strings.TrimSpace(r.FormValue("bio"))

	if username == "" {
		data := buildProfileData(user, "", "Le nom d'utilisateur ne peut pas être vide.")
		_ = authTemplates.ExecuteTemplate(w, "profile.html", data)
		return
	}

	if username != user.Username {
		var exists int
		_ = dbConn.QueryRow(`SELECT COUNT(*) FROM users WHERE username = $1 AND id != $2`, username, user.ID).Scan(&exists)
		if exists > 0 {
			data := buildProfileData(user, "", "Ce nom d'utilisateur est déjà pris.")
			_ = authTemplates.ExecuteTemplate(w, "profile.html", data)
			return
		}
	}

	_, err := dbConn.Exec(`UPDATE users SET username = $1, bio = $2 WHERE id = $3`, username, bio, user.ID)
	if err != nil {
		data := buildProfileData(user, "", "Erreur lors de la mise à jour.")
		_ = authTemplates.ExecuteTemplate(w, "profile.html", data)
		return
	}

	user.Username = username
	user.Bio = bio

	data := buildProfileData(user, "Profil mis à jour.", "")
	_ = authTemplates.ExecuteTemplate(w, "profile.html", data)
}

func profilePasswordHandler(w http.ResponseWriter, r *http.Request) {
	user, err := GetSessionUser(r, dbConn)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	currentPassword := r.FormValue("current_password")
	newPassword := r.FormValue("new_password")

	if currentPassword == "" || newPassword == "" {
		data := buildProfileData(user, "", "Remplis tous les champs.")
		_ = authTemplates.ExecuteTemplate(w, "profile.html", data)
		return
	}

	var hash string
	_ = dbConn.QueryRow(`SELECT password_hash FROM users WHERE id = $1`, user.ID).Scan(&hash)

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(currentPassword)); err != nil {
		data := buildProfileData(user, "", "Mot de passe actuel incorrect.")
		_ = authTemplates.ExecuteTemplate(w, "profile.html", data)
		return
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		data := buildProfileData(user, "", "Erreur serveur.")
		_ = authTemplates.ExecuteTemplate(w, "profile.html", data)
		return
	}

	_, err = dbConn.Exec(`UPDATE users SET password_hash = $1 WHERE id = $2`, string(newHash), user.ID)
	if err != nil {
		data := buildProfileData(user, "", "Erreur serveur.")
		_ = authTemplates.ExecuteTemplate(w, "profile.html", data)
		return
	}

	data := buildProfileData(user, "Mot de passe changé.", "")
	_ = authTemplates.ExecuteTemplate(w, "profile.html", data)
}

func profileDeleteHandler(w http.ResponseWriter, r *http.Request) {
	user, err := GetSessionUser(r, dbConn)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, APIResponse{Success: false, Message: "Non authentifié"})
		return
	}

	if user.AvatarURL != "" {
		os.Remove("." + user.AvatarURL)
	}

	_, err = dbConn.Exec(`DELETE FROM users WHERE id = $1`, user.ID)
	if err != nil {
		log.Printf("Delete account error: %v", err)
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Message: "Erreur serveur"})
		return
	}

	cookie, _ := r.Cookie("session_id")
	if cookie != nil {
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

	writeJSON(w, http.StatusOK, APIResponse{Success: true})
}
