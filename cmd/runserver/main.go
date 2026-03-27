package main

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
)

var (
	store *sessions.CookieStore
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
		Domain:   "localhost",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		Secure:   false, // Set to false for local development without HTTPS
		SameSite: http.SameSiteLaxMode,
	}
}

func main() {
	r := mux.NewRouter()

	loginRouter := r.PathPrefix("/auth").Subrouter()
	loginRouter.HandleFunc("/", authIndexRoute)
	loginRouter.HandleFunc("/login/", loginRoute)
	loginRouter.HandleFunc("/logout/", logoutRoute)

	log.Println("Starting server on port 11000")
	http.ListenAndServe(":11000", r)
}

func loginRoute(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "is-authenticated")
	if err != nil {
		log.Println("Store get failed with: %w", err)
	}
	session.Values["state"] = true
	err = session.Save(r, w)
	if err != nil {
		log.Println("Session save failed with: %w", err)
	}
	fmt.Fprintf(w, "Currently on login page.")
	// http.Redirect(w, r, "/auth/", http.StatusSeeOther)
	db.
}
func logoutRoute(w http.ResponseWriter, r *http.Request) {}
func authIndexRoute(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "is-authenticated")
	if err != nil {
		log.Println("Failed retrieving is-authenticated cookie: %w")
	}

	// Check if user is authenticated
	if auth, isBool := session.Values["state"].(bool); isBool && auth {
		fmt.Fprintf(w, "Authenticated aka logged-in.")
		return
	}

	fmt.Fprintln(w, "Not authenticated aka logged-out.")
}
