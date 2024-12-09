package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"go-pci-dss/internal/models"
	"go-pci-dss/internal/services"
	"go-pci-dss/utils"

	"github.com/sirupsen/logrus"
)

/*
type contextKey string

const userContextKey contextKey = "user"*/

func GetCardholdersHandler(service *services.CardholderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if tokenString == "" {
			logrus.Warn("Missing token in Authorization header")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err.Error(),
				"ip":    r.RemoteAddr,
			}).Warn("Invalid token")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		logrus.WithFields(logrus.Fields{
			"username": claims.Username,
			"role":     claims.Role,
			"ip":       r.RemoteAddr,
		}).Info("User attempting to fetch cardholders")

		cardholders, err := service.GetAllCardholders()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"username": claims.Username,
				"ip":       r.RemoteAddr,
				"error":    err.Error(),
			}).Error("Failed to fetch cardholders")
			http.Error(w, "Could not fetch cardholders", http.StatusInternalServerError)
			return
		}

		logrus.WithFields(logrus.Fields{
			"username":  claims.Username,
			"ip":        r.RemoteAddr,
			"num_cards": len(cardholders),
		}).Info("Successfully fetched cardholders")

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cardholders)
	}
}

func CreateCardholderHandler(service *services.CardholderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if tokenString == "" {
			logrus.Warn("Missing token in Authorization header")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err.Error(),
				"ip":    r.RemoteAddr,
			}).Warn("Invalid token")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		logrus.WithFields(logrus.Fields{
			"username": claims.Username,
			"role":     claims.Role,
			"ip":       r.RemoteAddr,
		}).Info("User attempting to add a card")

		var cardholder models.Cardholder
		if err := json.NewDecoder(r.Body).Decode(&cardholder); err != nil {
			logrus.WithFields(logrus.Fields{
				"username": claims.Username,
				"ip":       r.RemoteAddr,
				"error":    err.Error(),
			}).Error("Failed to decode cardholder input")
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		if err := service.CreateCardholder(cardholder); err != nil {
			logrus.WithFields(logrus.Fields{
				"username": claims.Username,
				"ip":       r.RemoteAddr,
				"error":    err.Error(),
			}).Error("Failed to create cardholder")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		maskedCardNumber := maskCardNumber(cardholder.CardNumber)

		logrus.WithFields(logrus.Fields{
			"username": claims.Username,
			"ip":       r.RemoteAddr,
			"card":     maskedCardNumber,
		}).Info("Successfully added a card")

		w.WriteHeader(http.StatusCreated)
	}
}
func maskCardNumber(cardNumber string) string {
	if len(cardNumber) < 4 {
		return cardNumber // Ako broj kartice nije duži od 4 cifre, ne maskiramo ga
	}
	return "**** **** **** " + cardNumber[len(cardNumber)-4:] // Maskiraj sve osim poslednje četiri cifre
}
