// handlers/auth.go
package handlers

import (
    "encoding/json"
    "net/http"
    "Resume-backend/models"
    "Resume-backend/Auth"
    "Resume-backend/Db"
)

func SignupHandler(w http.ResponseWriter, r *http.Request) {
    var user models.User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    hashed, err := auth.HashPassword(user.Password)
    if err != nil {
        http.Error(w, "Error hashing password", http.StatusInternalServerError)
        return
    }
    user.Password = hashed

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
    if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    user, err := db.GetUserByEmail(creds.Email)
    if err != nil || !auth.CheckPasswordHash(creds.Password, user.Password) {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    token, _ := auth.GenerateJWT(user.ID)
    json.NewEncoder(w).Encode(map[string]string{"token": token, "email": user.Email, "isAdmin": fmt.Sprintf("%v", user.IsAdmin)})
}
