package middleware

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()

// InitLogger - funkcija za inicijalizaciju loggera
func InitLogger() {
	// Postavljamo format loga (JSON, tekstualni...); ovde je JSON format
	Logger.SetFormatter(&logrus.JSONFormatter{})

	// Logovanje u fajl (ako je potrebno), postaviće se na standardni izlaz ako dođe do greške
	file, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		logrus.Fatal("Failed to open log file:", err)
	}
	// Postavi fajl kao izlaz za logove
	logrus.SetOutput(file)
}

/*
// LogRequest - Middleware za logovanje HTTP zahteva
func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Započinjemo logovanje
		start := time.Now()

		// Nastavljamo sa obradom zahteva
		next.ServeHTTP(w, r)

		// Zapisujemo podatke u log
		Logger.WithFields(logrus.Fields{
			"method": r.Method,
			"url":    r.URL.Path,
			"ip":     r.RemoteAddr,
			"time":   time.Since(start),
		}).Info("HTTP request")
	})
}
*/
