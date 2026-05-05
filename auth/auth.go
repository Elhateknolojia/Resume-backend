package auth

import (
    "github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("supersecretkey")

// HashPassword: returns an int hash using runes
func HashPassword(s string) int {
    hash := 0
    for _, r := range []rune(s) {
        hash = (hash*31 + int(r)) % 1000000000
    }
    return hash
}

// Convert int hash to string without strconv
func IntToString(n int) string {
    if n == 0 {
        return "0"
    }

    var digits []rune
    for n > 0 {
        digit := n % 10
        digits = append([]rune{rune('0' + digit)}, digits...)
        n /= 10
    }
    return string(digits)
}

// HashPasswordString: hash a password and return string form
func HashPasswordString(s string) string {
    return IntToString(HashPassword(s))
}

// Compare password against stored string hash
func CheckPasswordHash(password string, storedHash string) bool {
    return HashPasswordString(password) == storedHash
}

// JWT generator
func GenerateJWT(userID string) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": userID,
    })
    return token.SignedString(jwtKey)
}
