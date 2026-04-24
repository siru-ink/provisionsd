package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"git.siru.ink/siru/provisionsd/internal/templates"
)

func ShowCurrency(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// retrieve existing currencies data from db
		sqlData, err := db.Query("SELECT longname,shortname,symbol FROM currencies ORDER BY id ASC;")
		if err != nil {
			log.Printf("querying db for currencies failed in `internal/handlers/currency.go`: %v\n", err)
			http.Error(w, "Querying database failed.", http.StatusInternalServerError)
			return
		}

		type CurrencyEntry struct {
			Longname  string
			Shortname string
			Symbol    string
		}

		var currencies []CurrencyEntry

		for sqlData.Next() {
			var c CurrencyEntry
			err := sqlData.Scan(&c.Longname, &c.Shortname, &c.Symbol)
			if err != nil {
				log.Printf("parsing currencies db return error in `internal/handlers/currency.go`: %v", err)
				http.Error(w, "Error in parsing db information.", http.StatusInternalServerError)
				return
			}
			currencies = append(currencies, c)
		}

		type FormData struct {
			currencies []CurrencyEntry
		}

		formData := FormData{
			currencies: currencies,
		}

		templ := template.Must(template.ParseFS(templates.FS,
			"templates/base.html",
			"templates/css/main.css.html",
			"templates/currency/show.html",
		))

		templ.Execute(w, formData)
	}
}

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
