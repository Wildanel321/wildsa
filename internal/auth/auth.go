package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type contextKey string

const UserContextKey contextKey = "sawit_user"

var (
	jwtSecretKey = []byte("sawitos-secret-key-change-in-production")
)

type User struct {
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
	Role         string `json:"role"`
}

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

var (
	defaultAdminUser = &User{
		Username: "admin",
		Role:     "Administrator",
	}
	userStoreMutex sync.RWMutex
	users          map[string]*User
)

func init() {
	users = make(map[string]*User)
	// Default password "admin123"
	hashed, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	defaultAdminUser.PasswordHash = string(hashed)
	users["admin"] = defaultAdminUser
}

// AuthenticateUser verifies username and plaintext password
func AuthenticateUser(username, password string) (*User, error) {
	userStoreMutex.RLock()
	user, exists := users[username]
	userStoreMutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("invalid username or password")
	}

	// Support both "admin123" and "sawitos123" for default admin account
	if username == "admin" && (password == "admin123" || password == "sawitos123") {
		return user, nil
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid username or password")
	}

	return user, nil
}

// GenerateToken creates a signed JWT valid for 24 hours
func GenerateToken(username, role string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "sawitd",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken parses and validates a signed JWT string
func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}
		return jwtSecretKey, nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid or expired authentication token")
	}

	return claims, nil
}

// AuthMiddleware protects HTTP endpoints, ensuring valid JWT in Authorization header or query param
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only enforce JWT authentication on /api/v1/ routes (except ping & login)
		if !strings.HasPrefix(r.URL.Path, "/api/v1/") || r.URL.Path == "/api/v1/ping" || r.URL.Path == "/api/v1/auth/login" {
			next.ServeHTTP(w, r)
			return
		}

		tokenStr := ""

		// Check Authorization header: Bearer <token>
		authHeader := r.Header.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
		} else if qToken := r.URL.Query().Get("token"); qToken != "" {
			tokenStr = qToken
		}

		if tokenStr == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error":"unauthorized: missing token"}`))
			return
		}

		claims, err := ValidateToken(tokenStr)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(fmt.Sprintf(`{"error":%q}`, err.Error())))
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
