package handlers

import (
	"encoding/json"
	"go-pci-dss/models"
	"go-pci-dss/services"
	"net/http"
)

type CardHandler struct {
	service *services.CardService
}

func NewCardHandler(service *services.CardService) *CardHandler {
	return &CardHandler{service: service}
}

func (h *CardHandler) SaveCard(w http.ResponseWriter, r *http.Request) {
	var card models.Card
	if err := json.NewDecoder(r.Body).Decode(&card); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	if err := h.service.SaveCard(card); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *CardHandler) GetCard(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	card, err := h.service.GetCard(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(card)
}
