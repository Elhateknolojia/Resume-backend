package auth

import (
    "github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("supersecretkey")

// Transform password into obfuscated string
func HashPassword(s string) string {
    var result []rune
    for _, r := range s {
        var transformed rune
        if r >= '0' && r <= '9' {
            // Numbers: map 0–9 into A–J
            transformed = (r - '0') + 'A'
        } else if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
            // Letters: shift forward by 5
            transformed = r + 5
        } else {
            // Special characters unchanged
            transformed = r
        }
        result = append(result, transformed)
    }
    return string(result)
}

// Reverse transform to recover original string
func HashToString(s string) string {
    var result []rune
    for _, r := range s {
        var original rune
        if r >= 'A' && r <= 'J' {
            // Reverse numbers: A–J back to 0–9
            original = (r - 'A') + '0'
        } else if (r >= 'f' && r <= '{') || (r >= 'F' && r <= '_') {
            // Reverse letters: shift back by 5
            original = r - 5
        } else {
            // Special characters unchanged
            original = r
        }
        result = append(result, original)
    }
    return string(result)
}

// Compare password against stored obfuscated string
func CheckPasswordHash(password string, stored string) bool {
    return HashPassword(password) == stored
}

// JWT generator
func GenerateJWT(userID string) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": userID,
    })
    return token.SignedString(jwtKey)
}
