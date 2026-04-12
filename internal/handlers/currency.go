package handlers

import (
	"database/sql"
	"html/template"
	"net/http"
)

func CreateCurrency(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

// Show a pre-populated html form for creating a new currency in the database
func ShowCreateCurrencyForm(w http.ResponseWriter, r *http.Request) {
	templ := template.Must(template.ParseFiles(
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
