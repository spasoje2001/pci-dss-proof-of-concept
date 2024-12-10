package tests

import (
	"go-pci-dss/internal/middleware"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestAccessWithValidTokenButWrongRole(t *testing.T) {
	os.Setenv("JWT_SECRET_KEY", "yourSuperSecretJWTKey")
	validToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo5LCJ1c2VybmFtZSI6IlplbGEiLCJyb2xlIjoidXNlciIsImV4cCI6MTczMzg3NjIyNCwiaXNzIjoiZ28tcGNpLWRzcyJ9.cT9QOlCvCc4NJM_MNLkH8Dn0eSViAd0EX9P9yluLOGc"
	req, err := http.NewRequest("GET", "/cardholders", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+validToken)

	rr := httptest.NewRecorder()

	handler := middleware.AdminRoleMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Errorf("expected status 403, got %v", rr.Code)
	}

	expected := "Forbidden"
	if rr.Body.String() != expected {
		t.Errorf("expected response body %v, got %v", expected, rr.Body.String())
	}
}
