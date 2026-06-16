package server

import (
	"net/http"
	"strconv"
)

func adminHandler(w http.ResponseWriter, r *http.Request) {
	user, err := GetSessionUser(r, dbConn)
	if err != nil || !user.IsAdmin {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	stats, _ := GetAdminStats()
	users, _ := GetAllUsers()
	messages, _ := GetAllMessages()

	data := struct {
		User     *User
		Stats    AdminStats
		Users    []AdminUser
		Messages []AdminMessage
	}{
		User:     user,
		Stats:    stats,
		Users:    users,
		Messages: messages,
	}

	_ = authTemplates.ExecuteTemplate(w, "admin.html", data)
}

func adminBanHandler(w http.ResponseWriter, r *http.Request) {
	user, err := GetSessionUser(r, dbConn)
	if err != nil || !user.IsAdmin {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	userID, _ := strconv.Atoi(r.FormValue("user_id"))
	if userID > 0 {
		_ = BanUser(userID)
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func adminUnbanHandler(w http.ResponseWriter, r *http.Request) {
	user, err := GetSessionUser(r, dbConn)
	if err != nil || !user.IsAdmin {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	userID, _ := strconv.Atoi(r.FormValue("user_id"))
	if userID > 0 {
		_ = UnbanUser(userID)
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func adminDeleteMessageHandler(w http.ResponseWriter, r *http.Request) {
	user, err := GetSessionUser(r, dbConn)
	if err != nil || !user.IsAdmin {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	messageID, _ := strconv.Atoi(r.FormValue("message_id"))
	if messageID > 0 {
		_ = DeleteMessageByID(messageID)
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
