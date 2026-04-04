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
	cookies   *sessions.CookieStore
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
	cookies = sessions.NewCookieStore(key)
	cookies.Options = &sessions.Options{
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

	r.NotFoundHandler = http.HandlerFunc(Http404NotFound)
	r.MethodNotAllowedHandler = http.HandlerFunc(Http405MethodNotAllowed)

	// Add static file serving i.e. css
	fs := http.FileServer(http.Dir("assets/"))
	r.PathPrefix("/static").Handler(http.StripPrefix("/static", fs))

	loginRouter := r.PathPrefix("/auth").Subrouter()
	loginRouter.HandleFunc("/", auth(authIndexRoute, cookies))
	loginRouter.HandleFunc("/login/", loginPost(defaultDB, cookies)).Methods("POST")
	loginRouter.HandleFunc("/login/", loginGet).Methods("GET")
	loginRouter.HandleFunc("/logout/", logoutRoute(cookies))

	log.Println("Starting server on port 11000")
	http.ListenAndServe(":11000", r)
}

func Http403Forbidden(w http.ResponseWriter, r *http.Request) {
	templ := template.Must(template.ParseFiles("templates/base.html",
		"templates/css/main.css.html", "templates/httpcodes/403.html"))
	templ.Execute(w, "")
}

func Http404NotFound(w http.ResponseWriter, r *http.Request) {
	templ := template.Must(template.ParseFiles("templates/base.html",
		"templates/css/main.css.html", "templates/httpcodes/404.html"))
	templ.Execute(w, "")
}

func Http405MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	templ := template.Must(template.ParseFiles("templates/base.html",
		"templates/css/main.css.html", "templates/httpcodes/405.html"))
	templ.Execute(w, "")
}

func auth(next http.HandlerFunc, cookies *sessions.CookieStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get authentication cookie from cookie store
		authCookie, err := cookies.Get(r, "auth")

		// Some kind of error in decoding an existing cookie
		if err != nil {
			log.Println("Failed decoding existing <auth> cookie: %w")
		}

		// If user is not authenticated (i.e. logged out) render 403 Forbidden
		if auth, isBool := authCookie.Values["status"].(bool); !isBool || !auth {
			Http403Forbidden(w, r)
			return
		}

		// User is authenticated, so continue loading route
		next(w, r)
	}
}

func loginGet(w http.ResponseWriter, r *http.Request) {
	templ := template.Must(template.ParseFiles("templates/base.html",
		"templates/css/main.css.html", "templates/login.html"))
	templ.Execute(w, "")
}

// Returns a HanlderFunc that processes login POST form information via the given database handle
func loginPost(db *sql.DB, cookies *sessions.CookieStore) http.HandlerFunc {
	// Actual http request handler with dependency on the `db` database reference
	return func(w http.ResponseWriter, r *http.Request) {
		// Request login form data
		uname := r.FormValue("uname")
		passwd := r.FormValue("passwd")

		// Check for a db user entry with given data
		var userid int
		err := db.
			QueryRow("SELECT id FROM users WHERE uname == ? AND passwd == ?", uname, passwd).
			Scan(&userid)

		// No matching user could be found, redirecting back to login form
		if err != nil {
			http.Redirect(w, r, "/auth/login/", http.StatusNotFound)
			return
		}

		// Store new authentication cookie. Error decoding existing cookie can be ignored since we
		// are storing a new cookie
		authCookie, _ := cookies.Get(r, "auth")

		// Save new auth:status:true cookie
		authCookie.Values["status"] = true
		err = authCookie.Save(r, w)

		// Log errors in storing auth cookies for later review
		if err != nil {
			log.Println("Saving <auth> cookie failed: %w", err)
		}

		// Redirect to auth index page to show logged in status
		http.Redirect(w, r, "/auth/", http.StatusFound)
	}
}

// Returns a handler func that will reset auth cookie used for storing logged-in information
func logoutRoute(cookies *sessions.CookieStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Retreive existing cookie or save a new one
		authCookie, _ := cookies.Get(r, "auth")

		// Set logged in status to false
		authCookie.Values["status"] = false
		err := authCookie.Save(r, w)

		// Log errors in storing auth cookies for later review
		if err != nil {
			log.Printf("Saving <auth> cookie failed: %v", err)
		}

		// Redirect to auth index page to show logged-out status
		http.Redirect(w, r, "/auth/", http.StatusFound)
	}
}

// Should always run after auth middleware
func authIndexRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello logged in user.")
}
