// handlers/profile.go
package handlers

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value("user_id").(string)

    if r.Method == http.MethodGet {
        user, err := db.GetUserByID(userID)
        if err != nil {
            http.Error(w, "User not found", http.StatusNotFound)
            return
        }
        json.NewEncoder(w).Encode(user)
    }

    if r.Method == http.MethodPut {
        var update models.User
        if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
            http.Error(w, "Invalid request", http.StatusBadRequest)
            return
        }
        if err := db.UpdateUser(userID, update); err != nil {
            http.Error(w, "Update failed", http.StatusInternalServerError)
            return
        }
        json.NewEncoder(w).Encode(map[string]string{"message": "Profile updated"})
    }
}
