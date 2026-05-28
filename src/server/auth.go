package server

import (
    "database/sql"
    "errors"
    "fmt"
    "html"
    "html/template"
    "net/http"
    "strings"

    "golang.org/x/crypto/bcrypt"
)

var authTemplates = template.Must(template.ParseGlob("templates/*.html"))

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

        err := createUser(username, password)
        if err != nil {
            if strings.Contains(err.Error(), "Duplicate entry") {
                renderRegisterForm(w, "Ce nom d'utilisateur existe déjà.")
                return
            }
            http.Error(w, "Erreur serveur : "+html.EscapeString(err.Error()), http.StatusInternalServerError)
            return
        }

        fmt.Fprintf(w, `<h1>Inscription réussie</h1><p>Bienvenue %s !</p><p><a href="/login">Se connecter</a></p>`, html.EscapeString(username))
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

        ok, err := authenticateUser(username, password)
        if err != nil {
            http.Error(w, "Erreur serveur : "+html.EscapeString(err.Error()), http.StatusInternalServerError)
            return
        }
        if !ok {
            renderLoginForm(w, "Nom d'utilisateur ou mot de passe incorrect.")
            return
        }

        fmt.Fprintf(w, `<h1>Connexion réussie</h1><p>Bienvenue %s !</p><p><a href="/">Retour à l'accueil</a></p>`, html.EscapeString(username))
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

func createUser(username, password string) error {
    if dbConn == nil {
        return errors.New("connexion à la base non initialisée")
    }

    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    _, err = dbConn.Exec(`INSERT INTO users (username, password_hash) VALUES (?, ?)`, username, string(hash))
    return err
}

func authenticateUser(username, password string) (bool, error) {
    if dbConn == nil {
        return false, errors.New("connexion à la base non initialisée")
    }

    var hash string
    err := dbConn.QueryRow(`SELECT password_hash FROM users WHERE username = ?`, username).Scan(&hash)
    if err != nil {
        if err == sql.ErrNoRows {
            return false, nil
        }
        return false, err
    }

    err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    if err != nil {
        return false, nil
    }

    return true, nil
}
