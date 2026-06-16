package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type ReactionBody struct {
	MessageID int    `json:"message_id"`
	Emoji     string `json:"emoji"`
}

func reactionsAPIHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, err := GetSessionUser(r, dbConn)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, APIResponse{
			Success: false,
			Message: "Non authentifié",
		})
		return
	}

	switch r.Method {
	case http.MethodPost:
		handleAddReaction(w, r, user)
	case http.MethodDelete:
		handleRemoveReaction(w, r, user)
	default:
		json.NewEncoder(w).Encode(APIResponse{
			Success: false,
			Message: "method not allowed",
		})
	}
}

func handleAddReaction(w http.ResponseWriter, r *http.Request, user *User) {
	var body ReactionBody

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "json invalide",
		})
		return
	}

	emoji := strings.TrimSpace(body.Emoji)
	if emoji == "" || body.MessageID == 0 {
		writeJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "paramètres manquants",
		})
		return
	}

	if err := AddReaction(body.MessageID, user.ID, emoji); err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "erreur réaction",
		})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Success: true,
	})
}

func handleRemoveReaction(w http.ResponseWriter, r *http.Request, user *User) {
	messageIDStr := r.URL.Query().Get("message_id")
	emoji := strings.TrimSpace(r.URL.Query().Get("emoji"))

	messageID, err := strconv.Atoi(messageIDStr)
	if err != nil || emoji == "" {
		writeJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "paramètres manquants",
		})
		return
	}

	if err := RemoveReaction(messageID, user.ID, emoji); err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "erreur réaction",
		})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Success: true,
	})
}
