package auth

import (
    "github.com/golang-jwt/jwt/v4"
    "time"
    "log"
)

var jwtKey = []byte("supersecretkey")

// Transform password into obfuscated string

func HashPassword(s string) string {
    timestart := time.Now()
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
    log.Printf("Password hashing in %s", time.Since(timestart))
    return string(result)
}

// Reverse transform to recover original string
func HashToString(s string) string {
    timestart := time.Now()
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
    log.Printf("Password unhashing in %s", time.Since(timestart))
    return string(result)
}

// Compare password against stored obfuscated string
func CheckPasswordHash(password string, stored string) bool {
    return HashPassword(password) == stored
}

// JWT claims
type Claims struct {
    UserID   string `json:"user_id"`
    Email    string `json:"email"`
    Role     string `json:"role"`
    Tier     string `json:"tier"`
    Verified bool   `json:"verified"`
    Name     string `json:"name"`   // ✅ add this
    jwt.RegisteredClaims
}


func GenerateJWT(userID, email, role, tier string, verified bool, name string) (string, error) {
    claims := &Claims{
        UserID:   userID,
        Email:    email,
        Role:     role,
        Tier:     tier,
        Verified: verified,
        Name:     name,   // ✅ include name
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtKey)
}

