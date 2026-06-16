package server

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func profileAvatarHandler(w http.ResponseWriter, r *http.Request) {
	user, err := GetSessionUser(r, dbConn)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, APIResponse{Success: false, Message: "Non authentifié"})
		return
	}

	if err := r.ParseMultipartForm(2 << 20); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Message: "Fichier trop volumineux (max 2Mo)"})
		return
	}

	file, header, err := r.FormFile("avatar")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Message: "Fichier manquant"})
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Message: "Le fichier doit être une image"})
		return
	}

	ext := filepath.Ext(header.Filename)
	if ext == "" {
		switch contentType {
		case "image/jpeg":
			ext = ".jpg"
		case "image/png":
			ext = ".png"
		case "image/gif":
			ext = ".gif"
		case "image/webp":
			ext = ".webp"
		default:
			ext = ".jpg"
		}
	}

	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
	if !allowedExts[strings.ToLower(ext)] {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Message: "Format non supporté (jpg, png, gif, webp)"})
		return
	}

	if user.AvatarURL != "" {
		oldPath := "." + user.AvatarURL
		os.Remove(oldPath)
	}

	os.MkdirAll("static/avatars", 0755)

	filename := fmt.Sprintf("%d%s", user.ID, ext)
	savePath := filepath.Join("static", "avatars", filename)

	dst, err := os.Create(savePath)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Message: "Erreur serveur"})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Message: "Erreur serveur"})
		return
	}

	avatarURL := "/static/avatars/" + filename
	_, err = dbConn.Exec(`UPDATE users SET avatar_url = $1 WHERE id = $2`, avatarURL, user.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Message: "Erreur serveur"})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: map[string]string{"avatar_url": avatarURL}})
}
