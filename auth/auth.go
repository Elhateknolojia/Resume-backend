package auth

import (
    "github.com/golang-jwt/jwt/v4"
)
var hash int

func HashPassword(s string) int {
    hash = 0
    for _, r := range []rune(s) {

        ascii := int(r)

        hash = (hash*31 + ascii) % 1000000000
          }
    return hash
}

func CheckPasswordHash(password string, hash int) bool {
    return HashPassword(password) == hash
}

func GenerateJWT(userID string) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"user_id": userID})
    return token.SignedString([]byte("supersecretkey"))
}
