package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"go-pci-dss/internal/models"
	"go-pci-dss/internal/services"
)

func RegisterHandler(s *services.UserService) http.HandlerFunc {
	log.Printf("U handleru")
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "invalid input", http.StatusBadRequest)
			return
		}

		if err := s.RegisterUser(user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("user registered"))
	}
}

func LoginHandler(s *services.UserService) http.HandlerFunc {
	log.Printf("U login handleru")
	return func(w http.ResponseWriter, r *http.Request) {
		// Dekodiramo podatke iz tela zahteva (postavljeni su user i password)
		var userInput models.LoginInput
		if err := json.NewDecoder(r.Body).Decode(&userInput); err != nil {
			http.Error(w, "invalid input", http.StatusBadRequest)
			return
		}

		// Proveravamo da li korisnik postoji i da li je lozinka tačna
		token, err := s.Login(userInput.Username, userInput.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Vraćamo token korisniku
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"token": token})
	}
}
