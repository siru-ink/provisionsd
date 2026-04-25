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
		// Step 1: Execute the SQL query to retrieve currency data
		sqlData, err := db.Query("SELECT longname,shortname,symbol FROM currencies ORDER BY id ASC;")
		if err != nil {
			log.Printf("querying db for currencies failed in `internal/handlers/currency.go`: %v\n", err)
			http.Error(w, "Querying database failed.", http.StatusInternalServerError)
			return
		}

		// Step 2: Define the structure to hold currency information
		type CurrencyEntry struct {
			Longname  string // Full name of the currency (e.g., "United States Dollar")
			Shortname string // Abbreviated form of the currency (e.g., "USD")
			Symbol    string // Currency symbol (e.g., "$")
		}

		var currencies []CurrencyEntry

		// Step 3: Populate the currencies slice with data from the database query
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

		// Step 4: Define the structure to be passed to the template
		type FormData struct {
			currencies []CurrencyEntry // Slice containing all currency information
		}

		formdata := FormData{
			currencies: currencies,
		}

		// Step 5: Load and execute the template
		templ := template.Must(template.ParseFS(templates.FS,
			"templates/base.html",          // Base template file
			"templates/css/main.css.html",  // CSS styles
			"templates/currency/show.html", // Currency-specific template
		))

		templ.Execute(w, formdata)
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

func ShowDestroyCurrencyForm(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Query database for currency list in ascending order by ID
		sqldata, err := db.Query("SELECT id, longname FROM currencies ORDER BY id ASC")
		if err != nil {
			log.Printf("error retrieving currencies from db in `internal/handlers/currency.go`: %v\n", err)
			http.Error(w, "Error retrieving information from database.", http.StatusInternalServerError)
			return
		}

		// Define structure to hold currency data
		type CurrencyEntry struct {
			Id       int
			Longname string
		}

		var currencies []CurrencyEntry

		// Read query results and populate currency list
		for sqldata.Next() {
			var c CurrencyEntry
			err := sqldata.Scan(&c.Id, &c.Longname)
			if err != nil {
				log.Printf("parsing of returned currencies db data failed: %v", err)
				http.Error(w, "Error parsing database information.", http.StatusInternalServerError)
				return
			}
			currencies = append(currencies, c)
		}

		// Define template data structure with currencies
		type TemplateData struct {
			Currencies []CurrencyEntry
		}

		templdata := TemplateData{
			Currencies: currencies,
		}

		// Parse and execute template to render response
		templ := template.Must(template.ParseFS(templates.FS,
			"templates/base.html",
			"templates/css/main.css.html",
			"templates/currency/destroy.html",
		))
		templ.Execute(w, templdata)
	}
}
