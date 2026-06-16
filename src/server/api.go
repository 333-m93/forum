package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// =====================
// RESPONSE STRUCT
// =====================
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// =====================
// REQUEST BODY (JSON)
// =====================
type MessageBody struct {
	Category string `json:"category"`
	Content  string `json:"content"`
}

// =====================
// MAIN HANDLER
// =====================
func messagesAPIHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	switch r.Method {

	case http.MethodGet:
		handleGetMessages(w, r)

	case http.MethodPost:
		handlePostMessage(w, r)

	default:
		writeJSON(w, http.StatusMethodNotAllowed, APIResponse{
			Success: false,
			Message: "Méthode non autorisée",
		})
	}
}

// =====================
// GET MESSAGES
// =====================
func handleGetMessages(w http.ResponseWriter, r *http.Request) {

	categoryName := strings.TrimSpace(r.URL.Query().Get("category"))

	if categoryName == "" {
		writeJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "category manquant",
		})
		return
	}

	cat, err := GetCategoryByName(categoryName, dbConn)
	if err != nil {
		log.Printf("⚠️ Category not found: %q — DB error: %v", categoryName, err)
		writeJSON(w, http.StatusNotFound, APIResponse{
			Success: false,
			Message: "Catégorie introuvable",
		})
		return
	}

	messages, err := GetMessages(cat.ID, dbConn)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Erreur récupération messages",
		})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    messages,
	})
}

// =====================
// POST MESSAGE
// =====================
func handlePostMessage(w http.ResponseWriter, r *http.Request) {

	user, err := GetSessionUser(r, dbConn)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, APIResponse{
			Success: false,
			Message: "Non authentifié",
		})
		return
	}

	var categoryName, content string

	contentType := r.Header.Get("Content-Type")

	// =====================
	// JSON MODE
	// =====================
	if strings.Contains(contentType, "application/json") {

		var body MessageBody

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeJSON(w, http.StatusBadRequest, APIResponse{
				Success: false,
				Message: "JSON invalide",
			})
			return
		}

		categoryName = body.Category
		content = body.Content

	} else {
		// =====================
		// FORM / MULTIPART MODE (ton chat.js actuel)
		// =====================
		if err := r.ParseForm(); err != nil {
			writeJSON(w, http.StatusBadRequest, APIResponse{
				Success: false,
				Message: "Formulaire invalide",
			})
			return
		}

		categoryName = r.FormValue("category")
		content = r.FormValue("content")
	}

	categoryName = strings.TrimSpace(categoryName)
	content = strings.TrimSpace(content)

	// =====================
	// VALIDATION
	// =====================
	if categoryName == "" || content == "" {
		writeJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Paramètres manquants",
		})
		return
	}

	if len(content) > 5000 {
		content = content[:5000]
	}

	cat, err := GetCategoryByName(categoryName, dbConn)
	if err != nil {
		log.Printf("⚠️ Category not found: %q — DB error: %v", categoryName, err)
		writeJSON(w, http.StatusNotFound, APIResponse{
			Success: false,
			Message: "Catégorie introuvable",
		})
		return
	}

	msg, err := PostMessage(cat.ID, user.ID, content, dbConn)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Erreur ajout message",
		})
		return
	}

	writeJSON(w, http.StatusCreated, APIResponse{
		Success: true,
		Data:    msg,
	})
}

// =====================
// HELPERS
// =====================
func writeJSON(w http.ResponseWriter, status int, data APIResponse) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
