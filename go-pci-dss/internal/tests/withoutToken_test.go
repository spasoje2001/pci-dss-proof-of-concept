package tests

import (
	"go-pci-dss/internal/middleware"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAccessWithoutToken(t *testing.T) {
	req, err := http.NewRequest("GET", "/cardholders", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := middleware.AdminRoleMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %v", rr.Code)
	}

	expected := "Missing token"
	if rr.Body.String() != expected {
		t.Errorf("expected response body %v, got %v", expected, rr.Body.String())
	}
}
