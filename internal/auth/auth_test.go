package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sawitos/sawit/internal/auth"
)

func TestAuthenticateAdminUserSuccess(t *testing.T) {
	user, err := auth.AuthenticateUser("admin", "sawitos123")
	if err != nil {
		t.Fatalf("expected successful authentication for default admin: %v", err)
	}
	if user.Username != "admin" {
		t.Errorf("expected username 'admin', got %s", user.Username)
	}
}

func TestAuthenticateUserInvalidPassword(t *testing.T) {
	_, err := auth.AuthenticateUser("admin", "wrongpass")
	if err == nil {
		t.Error("expected authentication error for invalid password")
	}
}

func TestJWTTokenGenerationAndValidation(t *testing.T) {
	token, err := auth.GenerateToken("admin", "Administrator")
	if err != nil {
		t.Fatalf("expected clean JWT token generation: %v", err)
	}

	claims, err := auth.ValidateToken(token)
	if err != nil {
		t.Fatalf("expected valid token verification: %v", err)
	}
	if claims.Username != "admin" {
		t.Errorf("expected claim username 'admin', got %s", claims.Username)
	}
}

func TestAuthMiddlewareProtectedAccess(t *testing.T) {
	handler := auth.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))

	// Case 1: No header -> 401
	req1 := httptest.NewRequest("GET", "/api/v1/system", nil)
	rec1 := httptest.NewRecorder()
	handler.ServeHTTP(rec1, req1)
	if rec1.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 Unauthorized without token, got %d", rec1.Code)
	}

	// Case 2: Valid Bearer header -> 200
	token, _ := auth.GenerateToken("admin", "Administrator")
	req2 := httptest.NewRequest("GET", "/api/v1/system", nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	rec2 := httptest.NewRecorder()
	handler.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusOK {
		t.Errorf("expected 200 OK with valid token, got %d", rec2.Code)
	}
}
