package tests

import (
	"go-pci-dss/internal/middleware"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAccessWithInvalidToken(t *testing.T) {
	req, err := http.NewRequest("GET", "/cardholders", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer invalid_token")

	rr := httptest.NewRecorder()

	handler := middleware.AdminRoleMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %v", rr.Code)
	}
	actual := strings.TrimSpace(rr.Body.String())

	expected := "Invalid token"

	if actual != expected {
		t.Errorf("expected response body '%v', got '%v'", expected, actual)
	} else {
		t.Log("Test passed: response body matches expected.")
	}
}
