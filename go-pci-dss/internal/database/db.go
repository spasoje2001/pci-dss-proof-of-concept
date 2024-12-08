package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // PostgreSQL driver
)

// Connect funkcija koja se koristi za povezivanje sa bazom podataka
func Connect() (*sql.DB, error) {
	// Učitaj vrednosti iz .env fajla
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Uzimamo konekcioni string iz okruženja
	connStr := os.Getenv("DB_CONNECTION_STRING")

	// Povezivanje sa bazom
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("could not open connection: %v", err)
	}

	// Proveravamo da li je konekcija validna
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("could not connect to the database: %v", err)
	}

	return db, nil
}
