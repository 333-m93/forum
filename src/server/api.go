package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// RESPONSE
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// BODY
type MessageBody struct {
	CategoryID int    `json:"category_id"`
	Content    string `json:"content"`
}

// MAIN
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
			Message: "Method not allowed",
		})
	}
}

// GET
func handleGetMessages(w http.ResponseWriter, r *http.Request) {

	idStr := r.URL.Query().Get("category_id")
	if idStr == "" {
		json.NewEncoder(w).Encode(APIResponse{
			Success: false,
			Message: "category_id manquant",
		})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		json.NewEncoder(w).Encode(APIResponse{
			Success: false,
			Message: "category_id invalide",
		})
		return
	}

	messages, err := GetMessages(id, dbConn)
	if err != nil {
		json.NewEncoder(w).Encode(APIResponse{
			Success: false,
			Message: "erreur messages",
		})
		return
	}

	json.NewEncoder(w).Encode(APIResponse{
		Success: true,
		Data:    messages,
	})
}

// POST
func handlePostMessage(w http.ResponseWriter, r *http.Request) {

	user, err := GetSessionUser(r, dbConn)
	if err != nil {
		json.NewEncoder(w).Encode(APIResponse{
			Success: false,
			Message: "non authentifié",
		})
		return
	}

	var body MessageBody

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		json.NewEncoder(w).Encode(APIResponse{
			Success: false,
			Message: "json invalide",
		})
		return
	}

	body.Content = strings.TrimSpace(body.Content)

	if body.CategoryID == 0 || body.Content == "" {
		json.NewEncoder(w).Encode(APIResponse{
			Success: false,
			Message: "paramètres manquants",
		})
		return
	}

	msg, err := PostMessage(body.CategoryID, user.ID, body.Content, dbConn)
	if err != nil {
		log.Println(err)
		json.NewEncoder(w).Encode(APIResponse{
			Success: false,
			Message: "erreur insertion message",
		})
		return
	}

	json.NewEncoder(w).Encode(APIResponse{
		Success: true,
		Data:    msg,
	})
}
