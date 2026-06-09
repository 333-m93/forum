package server

import (
	"encoding/json"
	"net/http"
	"strings"
)

// APIResponse structure standard
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// Body JSON fallback
type MessageBody struct {
	Category string `json:"category"`
	Content  string `json:"content"`
}

// Handler principal
func messagesAPIHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	switch r.Method {
	case http.MethodGet:
		handleGetMessages(w, r)
	case http.MethodPost:
		handlePostMessage(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{
			Success: false,
			Message: "Méthode non autorisée",
		})
	}
}

// -------------------- GET --------------------
func handleGetMessages(w http.ResponseWriter, r *http.Request) {
	categoryName := strings.TrimSpace(r.URL.Query().Get("category"))
	if categoryName == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{
			Success: false,
			Message: "category manquant",
		})
		return
	}

	cat, err := GetCategoryByName(categoryName, dbConn)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(APIResponse{
			Success: false,
			Message: "Catégorie introuvable",
		})
		return
	}

	messages, err := GetMessages(cat.ID, dbConn)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{
			Success: false,
			Message: "Erreur récupération messages",
		})
		return
	}

	json.NewEncoder(w).Encode(APIResponse{
		Success: true,
		Data:    messages,
	})
}

// -------------------- POST --------------------
func handlePostMessage(w http.ResponseWriter, r *http.Request) {
	user, err := GetSessionUser(r, dbConn)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(APIResponse{
			Success: false,
			Message: "Non authentifié",
		})
		return
	}

	var categoryName, content string

	contentType := r.Header.Get("Content-Type")

	// JSON
	if strings.Contains(contentType, "application/json") {
		var body MessageBody
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(APIResponse{
				Success: false,
				Message: "JSON invalide",
			})
			return
		}

		categoryName = body.Category
		content = body.Content
	} else {
		_ = r.ParseForm()
		categoryName = r.FormValue("category")
		content = r.FormValue("content")
	}

	categoryName = strings.TrimSpace(categoryName)
	content = strings.TrimSpace(content)

	if categoryName == "" || content == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{
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
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(APIResponse{
			Success: false,
			Message: "Catégorie introuvable",
		})
		return
	}

	// ❌ IMPORTANT: ne PAS escape HTML ici
	msg, err := PostMessage(cat.ID, user.ID, content, dbConn)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{
			Success: false,
			Message: "Erreur ajout message",
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(APIResponse{
		Success: true,
		Data:    msg,
	})
}
