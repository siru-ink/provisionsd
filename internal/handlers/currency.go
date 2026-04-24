package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"git.siru.ink/siru/provisionsd/internal/templates"
)

func CreateCurrency(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the form into data representation
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Bad form data.", http.StatusBadRequest)
			return
		}

		// Get form data
		longname := r.FormValue("currency-longname")
		shortname := r.FormValue("currency-shortname")
		symbol := r.FormValue("currency-symbol")

		// Attempt to store data in database
		_, err = db.Exec("INSERT INTO currencies(longname, shortname, symbol) VALUES ('?', '?', '?')", longname, shortname, symbol)
		if err != nil {
			log.Printf("Inserting currency into database failed: %v\n", err)
			http.Error(w, "Failed storing data in database.", http.StatusInternalServerError)
		}

		// Redirect on success
		http.Redirect(w, r, "/currency/", http.StatusFound)
	}
}

// Show a pre-populated html form for creating a new currency in the database
func ShowCreateCurrencyForm(w http.ResponseWriter, r *http.Request) {
	templ := template.Must(template.ParseFS(templates.FS,
		"templates/base.html",
		"templates/css/main.css.html",
		"templates/form/currency/create.html",
	))

	type FormData struct {
		PostUrl string
	}

	formData := FormData{
		PostUrl: "/apiv1/form/currency/create/",
	}

	templ.Execute(w, formData)
}

func DestroyCurrency(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func ShowDestroyCurrencyForm(w http.ResponseWriter, r *http.Request) {
	//
}
