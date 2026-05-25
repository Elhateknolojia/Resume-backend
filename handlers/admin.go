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

  token, err := auth.GenerateJWT(admin.ID, admin.Email, "admin", admin.Tier, true, admin.Name)
   if err != nil {
        http.Error(w, "Token generation failed", http.StatusInternalServerError)
        return
    }

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

func AdminUsersHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    users, err := db.GetAllUsers()
    if err != nil {
        http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}
