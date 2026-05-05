// handlers/auth.go
package handlers

import (
    "encoding/json"
    "net/http"
    "Backend/models"
    "Backend/auth"
    "Backend/db"
)

// handlers/auth.go
func SignupHandler(w http.ResponseWriter, r *http.Request) {
    var creds struct {
        Name     string `json:"name"`
        Email    string `json:"email"`
        Phone    string `json:"phone"`
        Address  string `json:"address"`
        Password string `json:"password"`
    }
    if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

	// Check if user already exists
    existingUser, _ := db.GetUserByEmail(creds.Email)
    if existingUser != nil {
        w.WriteHeader(http.StatusConflict) // 409 Conflict
        json.NewEncoder(w).Encode(map[string]string{
            "error":   "User already exists, try forget password",
            "message": "The email is already registered",
        })
        return
    }

    user := models.User{
        Name:     creds.Name,
        Email:    creds.Email,
        Phone:    creds.Phone,
        Address:  creds.Address,
        Password: auth.HashPasswordString(creds.Password), // store stringified hash
        IsAdmin:  false,
        Tier:     "free",
    }

    if err := db.CreateUser(user); err != nil {
        http.Error(w, "Error saving user", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"message": "User created"})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    var creds struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    json.NewDecoder(r.Body).Decode(&creds)

    user, err := db.GetUserByEmail(creds.Email)
    if err != nil || !auth.CheckPasswordHash(creds.Password, user.Password) {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    token, _ := auth.GenerateJWT(user.ID)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "success": true,
        "email":   user.Email,
        "isAdmin": user.IsAdmin,
        "tier":    user.Tier,
        "token":   token,
    })
}
