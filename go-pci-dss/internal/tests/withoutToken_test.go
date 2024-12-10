package tests

import (
	"go-pci-dss/internal/middleware"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAccessWithoutToken(t *testing.T) {
	req, err := http.NewRequest("GET", "/cardholders", nil)
	if err != nil {
		t.Fatal(err) // Zaustavi test ako je došlo do greške pri kreiranju zahteva
	}

	rr := httptest.NewRecorder() // Kreiranje recordera za hvatanje odgovora

	// Omotavanje handler-a u middleware
	handler := middleware.AdminRoleMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Pozivanje handler-a sa testnim zahtevom
	handler.ServeHTTP(rr, req)

	// Provera status koda
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %v", rr.Code)
	}

	expected := "Missing token"
	actual := strings.TrimSpace(rr.Body.String())

	t.Logf("Response body: '%s'", actual)

	if actual != expected {
		t.Errorf("expected response body '%v', got '%v'", expected, actual)
	} else {
		t.Log("Test passed: response body matches expected.")
	}
}
