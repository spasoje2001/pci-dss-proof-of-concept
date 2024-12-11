package handlers

import (
	"encoding/json"
	"net/http"

	"go-pci-dss/internal/models"
	"go-pci-dss/internal/services"

	"github.com/skip2/go-qrcode"

	"github.com/sirupsen/logrus"
)

func RegisterHandler(s *services.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Logujemo dolazni zahtev
		logrus.WithFields(logrus.Fields{
			"method":   r.Method,
			"endpoint": r.URL.Path,
			"ip":       r.RemoteAddr,
		}).Info("Received registration request")

		var user models.User
		// Pokušaj dekodiranja korisničkog inputa
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			// Logujemo grešku pri dekodiranju
			logrus.WithFields(logrus.Fields{
				"ip":    r.RemoteAddr,
				"error": err.Error(),
			}).Error("Failed to decode user input")
			http.Error(w, "invalid input", http.StatusBadRequest)
			return
		}

		// Pokušaj registracije korisnika
		uri, err := s.RegisterUser(user)
		if err != nil {
			// Logujemo grešku prilikom registracije
			logrus.WithFields(logrus.Fields{
				"ip":    r.RemoteAddr,
				"error": err.Error(),
			}).Error("Failed to register user")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Logujemo uspešnu registraciju
		logrus.WithFields(logrus.Fields{
			"username": user.Username,
			"ip":       r.RemoteAddr,
		}).Info("User successfully registered")

		// Odgovor za uspešnu registraciju
		/*
			if uri, err := s.RegisterUser(user); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)

		*/

		qrCode, err := qrcode.Encode(uri, qrcode.Medium, 256) // Generišemo QR kod sa URI-jem
		if err != nil {
			http.Error(w, "Failed to generate QR code", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "image/png")
		w.WriteHeader(http.StatusCreated)
		if _, err := w.Write(qrCode); err != nil {
			logrus.WithFields(logrus.Fields{
				"username": user.Username,
				"ip":       r.RemoteAddr,
				"error":    err.Error(),
			}).Error("Failed to write QR code to response")
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
			return
		}
	}
}

func LoginHandler(s *services.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logrus.Infof("Received login request from IP: %s", r.RemoteAddr)

		var userInput models.LoginInput
		if err := json.NewDecoder(r.Body).Decode(&userInput); err != nil {
			logrus.WithFields(logrus.Fields{
				"username": userInput.Username,
				"ip":       r.RemoteAddr,
			}).Error("Invalid input data during login attempt")
			http.Error(w, "invalid input", http.StatusBadRequest)
			return
		}

		token, err := s.Login(userInput.Username, userInput.Password, userInput.TOTPSecret)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"username": userInput.Username,
				"ip":       r.RemoteAddr,
			}).Warn("Failed login attempt")
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		logrus.WithFields(logrus.Fields{
			"username": userInput.Username,
			"ip":       r.RemoteAddr,
			"time":     r.Header.Get("Date"),
		}).Info("User successfully logged in")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(map[string]string{"token": token}); err != nil {
			logrus.WithFields(logrus.Fields{
				"username": userInput.Username,
				"ip":       r.RemoteAddr,
				"error":    err.Error(),
			}).Error("Failed to encode login response")
			http.Error(w, "Failed to generate response", http.StatusInternalServerError)
			return
		}
	}
}
