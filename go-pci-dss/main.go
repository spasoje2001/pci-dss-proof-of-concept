package main

import (
	"go-pci-dss/handlers"
	"go-pci-dss/repositories"
	"go-pci-dss/services"
	"log"
	"net/http"
)

func main() {
	// Kreiraj skladište
	repo := repositories.NewCardRepository()

	// Kreiraj servis
	service := services.NewCardService(repo)

	// Kreiraj handler
	handler := handlers.NewCardHandler(service)

	// Postavi rute
	http.HandleFunc("/card/save", handler.SaveCard)
	http.HandleFunc("/card/get", handler.GetCard)

	// Pokreni server
	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
