// middleware/authmiddleware.go
package middleware

import (
    "context"
    "net/http"
    "strings"

    "github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("supersecretkey")

func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Missing token", http.StatusUnauthorized)
            return
        }

        tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
        claims := jwt.MapClaims{}
        token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
            return jwtKey, nil
        })
        if err != nil || !token.Valid {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        userID := claims["user_id"].(string)
        ctx := context.WithValue(r.Context(), "user_id", userID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
