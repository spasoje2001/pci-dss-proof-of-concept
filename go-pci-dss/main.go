package main

import (
	"log"
	"net/http"

	"go-pci-dss/internal/database"
	"go-pci-dss/internal/handlers"
	"go-pci-dss/internal/services"

	"github.com/gorilla/mux"
)

func main() {
	// 1. Povezivanje sa bazom podataka
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	defer db.Close()
	/*
		if err := database.ExecuteMigration(db, "000001_create_cardholders_table.up.sql"); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}*/

	// 2. Kreiranje servisa
	cardholderService := services.NewCardholderService(db)

	// 3. Kreiranje router-a
	r := mux.NewRouter()

	// 5. Definisanje ruta
	r.HandleFunc("/cardholders", handlers.GetCardholdersHandler(cardholderService)).Methods("GET")
	r.HandleFunc("/cardholders", handlers.CreateCardholderHandler(cardholderService)).Methods("POST")

	// 6. Pokretanje servera
	port := "8080"
	log.Printf("Server is running on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
