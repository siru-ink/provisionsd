package main

import (
	"crypto/rand"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var (
	key = func() []byte {
		k := make([]byte, 64)
		rand.Read(k)
		return k
	}()
	store = sessions.NewCookieStore(key)
)

func main() {
	initCookieStoreOptions()

	r := mux.NewRouter()

	loginRouter := r.PathPrefix("/auth").Subrouter()
	loginRouter.HandleFunc("/", authIndexRoute)
	loginRouter.HandleFunc("/login/", loginRoute)
	loginRouter.HandleFunc("/logout/", logoutRoute)

	fmt.Println("Starting server on port 11000")
	http.ListenAndServe(":11000", r)
}

func loginRoute(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "siru-cookie-name")
	if err != nil {
		fmt.Println("Store get failed with: %w", err)
	}
	session.Values["authenticated"] = true
	err = session.Save(r, w)
	if err != nil {
		fmt.Println("Session save failed with: %w", err)
	}
	fmt.Fprintf(w, "Currently on login page.")
	// http.Redirect(w, r, "/auth/", http.StatusSeeOther)
}
func logoutRoute(w http.ResponseWriter, r *http.Request) {}
func authIndexRoute(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "siru-cookie-name")

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		fmt.Fprintf(w, "Authenticated aka logged-in.")
		return
	}

	// Print secret message
	fmt.Fprintln(w, "Not authenticated aka logged-out.")
}

func initCookieStoreOptions() {
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		Secure:   false, // Set to false for local development without HTTPS
		SameSite: http.SameSiteLaxMode,
	}
}
