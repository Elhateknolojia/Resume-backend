package handlers

import (
	"Backend/auth"
	"Backend/db"
	"Backend/middleware"
	"Backend/models"
    "time"
	"encoding/json"
	// "go/token"
	// "hash"
	"log"
	"net/http"

	// "time"
	"github.com/golang-jwt/jwt/v4"
)

// SignupHandler: create new user with transformed password
func SignupHandler(w http.ResponseWriter, r *http.Request) {
    start := time.Now()
    var creds struct {
        Name     string `json:"name"`
        Email    string `json:"email"`
        Phone    string `json:"phone"`
        Address  string `json:"address"`
        Password string `json:"password"`
    }
    log.Printf("Decode signup request in %s", time.Since(start))
    if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    // Check if user already exists
    existingUserstart := time.Now()
    existingUser, _ := db.GetUserByEmail(creds.Email)
    log.Printf("Existing user lookup in %s", time.Since(existingUserstart))


    if existingUser != nil {
        w.WriteHeader(http.StatusConflict)
        json.NewEncoder(w).Encode(map[string]string{
            "error":   "User already exists, try forget password",
            "message": "The email is already registered",
        })
        return
    }

    creatuserstart := time.Now()
    user := models.User{
        Name:     creds.Name,
        Email:    creds.Email,
        Phone:    creds.Phone,
        Address:  creds.Address,
        // Store transformed password string
        Password: auth.HashPassword(creds.Password),
        IsAdmin:  false,
        Tier:     "free",
    }
    log.Printf("User  creation in %s", time.Since(creatuserstart))

    otp := generateOTP()
    expiry := time.Now().Add(10 * time.Minute).Unix()

    user.OTPCode = otp
    user.OTPExpiry = expiry

    if err := db.CreateUser(user); err != nil {
        http.Error(w, "Error saving user", http.StatusInternalServerError)
        return
    }

    if err := sendEmailOTP(user.Email, otp); err != nil {
        log.Printf("SendGrid error: %v", err)
        http.Error(w, "Failed to send OTP", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
json.NewEncoder(w).Encode(map[string]string{"message": "OTP sent"})
}

// LoginHandler: verify password by re-transforming input// LoginHandler: verify password by re-transforming input
func LoginHandler(w http.ResponseWriter, r *http.Request) {
    start := time.Now()
    var creds struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    // ⏱ Decode request
    if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }
    log.Printf("[Login] Decode request in %s", time.Since(start))

    // ⏱ DB lookup
    dbStart := time.Now()
    user, err := db.GetUserByEmail(creds.Email)
    log.Printf("[Login] DB lookup in %s", time.Since(dbStart))

    if err != nil || user == nil {
        w.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(w).Encode(map[string]string{
            "error":   "Email mismatch",
            "message": "No account found with that email",
        })
        return
    }

    // ⏱ Password check
    hashStart := time.Now()
    if !auth.CheckPasswordHash(creds.Password, user.Password) {
        log.Printf("[Login] Password verification in %s", time.Since(hashStart))
        w.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(w).Encode(map[string]string{
            "error":   "Password mismatch",
            "message": "The password you entered is incorrect",
        })
        return
    }
    log.Printf("[Login] Password verification in %s", time.Since(hashStart))

    // ⏱ Token generation
    tokenStart := time.Now()
    token, err := middleware.GenerateToken(user.ID)
    if err != nil {
        http.Error(w, "Could not generate token", http.StatusInternalServerError)
        return
    }
    log.Printf("[Login] Token generation in %s", time.Since(tokenStart))

    // ✅ Success response
    resp := map[string]interface{}{
        "token":          token,
        "userId":         user.ID,
        "email":          user.Email,
        "isAdmin":        user.IsAdmin,
        "tier":           user.Tier,
        "passwordStatus": "Password is a match",
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}



func RefreshHandler(w http.ResponseWriter, r *http.Request) {
    authHeader := r.Header.Get("Authorization")
    if authHeader == "" {
        http.Error(w, "Missing token", http.StatusUnauthorized)
        return
    }

    tokenStr := authHeader[len("Bearer "):]
    claims := jwt.MapClaims{}
    token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
        return middleware.JwtKey(), nil
    })
    if err != nil || !token.Valid {
        http.Error(w, "Invalid token", http.StatusUnauthorized)
        return
    }

    // Issue new token with 10‑minute expiry
    userID := claims["user_id"].(string)
    newToken, err := middleware.GenerateToken(userID)
    if err != nil {
        http.Error(w, "Failed to refresh token", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{"token": newToken})
}
