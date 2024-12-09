package middleware

import (
	"context"
	"go-pci-dss/utils"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

type contextKey string

const userContextKey contextKey = "user"

func AdminRoleMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			logrus.WithFields(logrus.Fields{
				"ip":       r.RemoteAddr,
				"endpoint": r.URL.Path,
			}).Error("Unauthorized request - Missing token")
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// Validiramo token
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"ip":       r.RemoteAddr,
				"endpoint": r.URL.Path,
				"error":    err.Error(),
			}).Error("User attempted unauthorized access or action - Invalid token")
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if claims.Role != "admin" {
			logrus.WithFields(logrus.Fields{
				"ip":       r.RemoteAddr,
				"endpoint": r.URL.Path,
				"user":     claims.Username,
			}).Warn("Forbidden: User does not have admin role")
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		logrus.WithFields(logrus.Fields{
			"ip":       r.RemoteAddr,
			"endpoint": r.URL.Path,
			"user":     claims.Username,
		}).Info("Admin access granted")

		ctx := context.WithValue(r.Context(), userContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserRoleMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			logrus.WithFields(logrus.Fields{
				"ip":       r.RemoteAddr,
				"endpoint": r.URL.Path,
			}).Error("Unauthorized request - Missing token")
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// Validiramo token
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"ip":       r.RemoteAddr,
				"endpoint": r.URL.Path,
				"error":    err.Error(),
			}).Error("Invalid token")
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Proveravamo da li je korisnik sa ulogom "user"
		if claims.Role != "user" {
			logrus.WithFields(logrus.Fields{
				"ip":       r.RemoteAddr,
				"endpoint": r.URL.Path,
				"user":     claims.Username,
			}).Warn("Forbidden: User does not have user role")
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		logrus.WithFields(logrus.Fields{
			"ip":       r.RemoteAddr,
			"endpoint": r.URL.Path,
			"user":     claims.Username,
		}).Info("User access granted")

		// Dodajemo claims u context
		ctx := context.WithValue(r.Context(), userContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
