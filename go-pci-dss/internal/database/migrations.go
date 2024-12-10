package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
)

// ExecuteMigration pokreće SQL migracije (up ili down).
func ExecuteMigration(db *sql.DB, migrationName string) error {
	// Preuzimanje SQL fajla
	migrationPath := filepath.Join("migrations", migrationName)
	sqlQuery, err := ioutil.ReadFile(migrationPath)
	if err != nil {
		return fmt.Errorf("could not read migration file %s: %v", migrationName, err)
	}

	// Izvršavanje SQL upita
	_, err = db.Exec(string(sqlQuery))
	if err != nil {
		return fmt.Errorf("could not execute migration %s: %v", migrationName, err)
	}

	log.Printf("Migration %s applied successfully", migrationName)
	return nil
}
