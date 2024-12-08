package handlers

import (
	"encoding/json"
	"net/http"

	"go-pci-dss/internal/models"
	"go-pci-dss/internal/services"
)

func GetCardholdersHandler(service *services.CardholderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cardholders, err := service.GetAllCardholders()
		if err != nil {
			http.Error(w, "Could not fetch cardholders", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(cardholders)
	}
}

func CreateCardholderHandler(service *services.CardholderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var cardholder models.Cardholder
		if err := json.NewDecoder(r.Body).Decode(&cardholder); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		if err := service.CreateCardholder(cardholder); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}
