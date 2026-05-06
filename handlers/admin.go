package handlers

import (
    "encoding/json"
    "net/http"

    "Backend/auth"
    "Backend/db"
)

func AdminLoginHandler(w http.ResponseWriter, r *http.Request) {
    var creds struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    json.NewDecoder(r.Body).Decode(&creds)

    admin, err := db.GetAdminByEmail(creds.Email)
    if err != nil || !auth.CheckPasswordHash(creds.Password, admin.Password) {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    token, _ := auth.GenerateJWT(admin.ID)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "success": true,
        "email":   admin.Email,
        "role":    admin.Role,
        "tier":    admin.Tier,
        "token":   token,
    })
}

// AdminStatsHandler returns basic system stats for the admin dashboard
func AdminStatsHandler(w http.ResponseWriter, r *http.Request) {
    stats := db.GetStats() // implement this in db package
    json.NewEncoder(w).Encode(stats)
}