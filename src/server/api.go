package server

import (
	"encoding/json"
	"net/http"
	"strings"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type MessageBody struct {
	Category string `json:"category"`
	Content  string `json:"content"`
}

func messagesAPIHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		handleGetMessages(w, r)
	case http.MethodPost:
		handlePostMessage(w, r)
	default:
		json.NewEncoder(w).Encode(APIResponse{
			Success: false,
			Message: "method not allowed",
		})
	}
}

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
		writeJSON(w, http.StatusNotFound, APIResponse{
			Success: false,
			Message: "catégorie introuvable",
		})
		return
	}

	messages, err := GetMessages(cat.ID, dbConn)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "erreur messages",
		})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    messages,
	})
}

func handlePostMessage(w http.ResponseWriter, r *http.Request) {
	user, err := GetSessionUser(r, dbConn)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, APIResponse{
			Success: false,
			Message: "Non authentifié",
		})
		return
	}

	var body MessageBody

	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "json invalide",
		})
		return
	}

	categoryName := strings.TrimSpace(body.Category)
	content := strings.TrimSpace(body.Content)

	if categoryName == "" || content == "" {
		writeJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "paramètres manquants",
		})
		return
	}

	cat, err := GetCategoryByName(categoryName, dbConn)
	if err != nil {
		writeJSON(w, http.StatusNotFound, APIResponse{
			Success: false,
			Message: "catégorie introuvable",
		})
		return
	}

	msg, err := PostMessage(cat.ID, user.ID, content, dbConn)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "erreur message",
		})
		return
	}

	writeJSON(w, http.StatusCreated, APIResponse{
		Success: true,
		Data:    msg,
	})
}

func writeJSON(w http.ResponseWriter, status int, data APIResponse) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
