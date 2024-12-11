package tests

import (
	"go-pci-dss/internal/middleware"
	"go-pci-dss/utils"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestAccessWithValidTokenButWrongRole(t *testing.T) {
	os.Setenv("JWT_SECRET_KEY", "yourSuperSecretJWTKey")
	token, err := utils.GenerateJWT(9, "Zela", "user")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	req, err := http.NewRequest("GET", "/cardholders", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()

	handler := middleware.AdminRoleMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Errorf("expected status 403, got %v", rr.Code)
	}
	actual := strings.TrimSpace(rr.Body.String())

	expected := "Forbidden"
	if actual != expected {
		t.Errorf("expected response body '%v', got '%v'", expected, actual)
	} else {
		t.Log("Test passed: response body matches expected.")
	}
}
