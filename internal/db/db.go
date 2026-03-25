package db

import (
	"context"
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

// Open a new sqlite database connection to the hardcoded `./guestbook.db` file. Initializes and
// guarantees that the following migrations will have been run:
//   - migrations/001.sql
func InitDb() *sql.DB {
	// Attempt to open a database connection and return nil if impossible
	db, err := sql.Open("sqlite", "default.db")
	if err != nil {
		log.Panicln("Database file <./default.db> could not be opened.")
	}

	// Run migrations/001.sql
	sqlDBSetupBytes, err := os.ReadFile("migrations/001.sql")
	if err != nil {
		log.Panicln("Migration <./migrations/001.sql> could not be read/found.")
	}
	_, err = db.ExecContext(context.Background(), string(sqlDBSetupBytes))
	if err != nil {
		log.Panicln("Migration <./migrations/001.sql> failed to be commited to the database.")
	}

	log.Println("Database succesfully initialized.")
	return db
}
