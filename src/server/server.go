package server

import (
	"database/sql"
	"log"
	"net/http"
	"strings"
	"sync"

	"forum.com/m/src/config"
	"forum.com/m/src/database"
)

type Category struct {
	Slug        string
	Title       string
	Description string
}

type Message struct {
	Author  string
	Content string
	When    string
}

var dbConn *sql.DB

var categories = []Category{
	{Slug: "general", Title: "Chat général", Description: "Un espace pour parler de tout, des entraînements aux compétitions."},
	{Slug: "mma", Title: "MMA", Description: "Techniques mixtes, entraînement et actualité du MMA."},
	{Slug: "boxe-anglaise", Title: "Boxe Anglaise", Description: "Discussion sur la boxe pieds-poings, les combats et les gants."},
	{Slug: "muay-thai", Title: "Muay Thai", Description: "Échanges sur le muay thai, les clinschs et les pads."},
	{Slug: "jujitsu-bresilien", Title: "Jujitsu Brésilien", Description: "Soumissions, garde et sparring entre amis combattants."},
	{Slug: "grappling", Title: "Grappling", Description: "Lutte, prise de soumission et progression au sol."},
	{Slug: "autres-sports", Title: "Autres sports de combat", Description: "Arts martiaux et sports de combat variés pour tous."},
}

var discussionMessages = map[string][]Message{
	"general": {
		{Author: "Rayan", Content: "Bienvenue sur le chat général du Dojo !", When: "Il y a 2 minutes"},
		{Author: "Mila", Content: "Quel est votre meilleur conseil pour rester motivé après un entraînement difficile ?", When: "Il y a 1 minute"},
	},
	"mma": {
		{Author: "Léo", Content: "Qui a suivi le dernier combat de la saison ?", When: "Aujourd'hui"},
	},
}

var messagesMu sync.RWMutex

func helloHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Name       string
		Categories []Category
	}{
		Name:       "Le Dojo",
		Categories: categories,
	}
	if err := authTemplates.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, "Erreur template : "+err.Error(), http.StatusInternalServerError)
	}
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Name       string
		Categories []Category
	}{
		Name:       "Le Dojo",
		Categories: categories,
	}
	if err := authTemplates.ExecuteTemplate(w, "chat.html", data); err != nil {
		http.Error(w, "Erreur template : "+err.Error(), http.StatusInternalServerError)
	}
}

func discussionHandler(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/discussion/")
	if slug == "" {
		http.Redirect(w, r, "/chat", http.StatusSeeOther)
		return
	}

	category, ok := findCategory(slug)
	if !ok {
		http.NotFound(w, r)
		return
	}

	data := struct {
		Name       string
		Categories []Category
		Category   Category
		Messages   []Message
		Error      string
		Success    bool
	}{
		Name:       "Le Dojo",
		Categories: categories,
		Category:   category,
	}

	if r.Method == http.MethodPost {
		author := strings.TrimSpace(r.FormValue("author"))
		content := strings.TrimSpace(r.FormValue("content"))
		if author == "" || content == "" {
			data.Error = "Merci de renseigner un pseudo et un message."
		} else {
			messagesMu.Lock()
			discussionMessages[slug] = append(discussionMessages[slug], Message{Author: author, Content: content, When: "À l'instant"})
			messagesMu.Unlock()
			data.Success = true
		}
	}

	messagesMu.RLock()
	data.Messages = discussionMessages[slug]
	messagesMu.RUnlock()

	if err := authTemplates.ExecuteTemplate(w, "discussion.html", data); err != nil {
		http.Error(w, "Erreur template : "+err.Error(), http.StatusInternalServerError)
	}
}

func findCategory(slug string) (Category, bool) {
	for _, category := range categories {
		if category.Slug == slug {
			return category, true
		}
	}
	return Category{}, false
}

func StartServer() {
	cfg := config.Load()
	var err error
	dbConn, err = database.Connect(cfg)
	if err != nil {
		log.Fatalf("Erreur de connexion à la DB : %v", err)
	}
	defer dbConn.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/", helloHandler)
	mux.HandleFunc("/chat", chatHandler)
	mux.HandleFunc("/discussion/", discussionHandler)
	mux.HandleFunc("/register", registerHandler)
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/profile", profileHandler)
	mux.HandleFunc("/logout", logoutHandler)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Printf("Serveur démarré sur %s\n", cfg.Port)
	if err := http.ListenAndServe(cfg.Port, mux); err != nil {
		log.Fatalf("Erreur du serveur : %v", err)
	}
}
