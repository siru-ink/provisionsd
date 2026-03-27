package main

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"git.siru.ink/siru/provisionsd/internal/db"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
)

var (
	store     *sessions.CookieStore
	defaultDB *sql.DB
)

func init() {
	// Load .env configuration variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: %w", err)
	}

	// Get secret key from .env file
	key, err := hex.DecodeString(os.Getenv("SECRET_KEY"))
	if err != nil {
		log.Fatal("Error retrieving secret key from .env file: %w", err)
	}

	// Set up cookie store
	store = sessions.NewCookieStore(key)
	store.Options = &sessions.Options{
		Path:     "/",
		Domain:   "dev.siru.ink",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		Secure:   false, // Set to false for local development without HTTPS
		SameSite: http.SameSiteStrictMode,
	}

	// Init default sqlite3 db
	defaultDB = db.InitDb()
}

func main() {
	r := mux.NewRouter()

	// Add static file serving i.e. css
	fs := http.FileServer(http.Dir("assets/"))
	r.PathPrefix("/static").Handler(http.StripPrefix("/static", fs))

	loginRouter := r.PathPrefix("/auth").Subrouter()
	loginRouter.HandleFunc("/", authIndexRoute)
	loginRouter.HandleFunc("/login/", authLoginPost).Methods("POST")
	loginRouter.HandleFunc("/login/", authLoginGet).Methods("GET")
	loginRouter.HandleFunc("/logout/", logoutRoute)

	log.Println("Starting server on port 11000")
	http.ListenAndServe(":11000", r)
}

func authLoginGet(w http.ResponseWriter, r *http.Request) {
	templ := template.Must(template.ParseFiles("templates/base.html", "templates/login.html"))
	templ.Execute(w, "")
}

func authLoginPost(w http.ResponseWriter, r *http.Request) {
	uname := r.FormValue("uname")
	passwd := r.FormValue("passwd")

	var userid int
	err := defaultDB.
		QueryRow("SELECT id FROM users WHERE uname == ? AND passwd == ?", uname, passwd).
		Scan(&userid)

	if err != nil {
		http.Redirect(w, r, "/auth/login/", http.StatusNotFound)
		return
	}

	authCookie, err := store.Get(r, "authCookie")
	if err != nil {
		log.Println("Cookie retrieval <authCookie> failed: %w", err)
	}
	authCookie.Values["logged-in"] = true
	err = authCookie.Save(r, w)
	if err != nil {
		log.Println("Cookie save <authCookie> failed: %w", err)
	}
	http.Redirect(w, r, "/auth/", http.StatusFound)
}

func logoutRoute(w http.ResponseWriter, r *http.Request) {}
func authIndexRoute(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "authCookie")
	if err != nil {
		log.Println("Failed retrieving is-authenticated cookie: %w")
	}

	// Check if user is authenticated
	if auth, isBool := session.Values["logged-in"].(bool); isBool && auth {
		fmt.Fprintf(w, "Authenticated aka logged-in.")
		return
	}

	fmt.Fprintln(w, "Not authenticated aka logged-out.")
}
