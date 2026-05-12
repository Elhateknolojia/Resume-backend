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

    if err := db.CreateUser(user); err != nil {
        http.Error(w, "Error saving user", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"message": "User created"})
}

// LoginHandler: verify password by re-transforming input
func LoginHandler(w http.ResponseWriter, r *http.Request) {
    start := time.Now()
    var creds struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    log.Printf("Decode login request in %s", time.Since(start))
    if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }


    dbSart := time.Now()
    user, err := db.GetUserByEmail(creds.Email)
    log.Printf("DB lookup in %s", time.Since(dbSart))

    hashStart := time.Now()
    if err != nil || !auth.CheckPasswordHash(creds.Password, user.Password) {
        w.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(w).Encode(map[string]string{
            "error":   "Invalid email or password",
            "message": "Please check your credentials and try again",
        })
        return
    }
    log.Printf("Password verification in %s", time.Since(hashStart))


    tokenStart := time.Now()
    token, err := middleware.GenerateToken(user.ID)
    if err != nil {
        http.Error(w, "Could not generate token", http.StatusInternalServerError)
        return
    }
    log.Printf("Token generation in %s", time.Since(tokenStart))


            resp := models.LoginResponse{
                Token:   token,
                UserID:  user.ID,
                Email:   user.Email,
                IsAdmin: user.IsAdmin,
                Tier:    user.Tier,
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
