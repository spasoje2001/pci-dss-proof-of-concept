package main

import (
	"log"
	"net/http"

	"go-pci-dss/internal/database"
	"go-pci-dss/internal/handlers"
	"go-pci-dss/internal/middleware"
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
	if err := database.ExecuteMigration(db, "000001_create_cardholders_table.up.sql"); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	cardholderService := services.NewCardholderService(db)
	userService := services.NewUserService(db)

	r := mux.NewRouter()

	r.Handle("/cardholders", middleware.AdminRoleMiddleware(handlers.GetCardholdersHandler(cardholderService))).Methods("GET")
	r.HandleFunc("/cardholders", handlers.CreateCardholderHandler(cardholderService)).Methods("POST")

	r.HandleFunc("/users", handlers.RegisterHandler(userService)).Methods("POST")
	r.HandleFunc("/users/login", handlers.LoginHandler(userService)).Methods("POST")

	port := "8080"
	log.Printf("Server is running on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
