// handlers/profile.go
package handlers

import (
	"net/http"
	"Backend/db"
	"encoding/json"
	"Backend/models"
    "context"
    "go.mongodb.org/mongo-driver/bson"
)

// handlers/profile.go

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
        json.NewDecoder(r.Body).Decode(&update)
        if err := db.UpdateUser(userID, &update); err != nil {
            http.Error(w, "Update failed", http.StatusInternalServerError)
            return
        }

        json.NewEncoder(w).Encode(map[string]string{"message": "Profile updated"})
    }
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
    email := r.URL.Query().Get("email")
    user, err := db.GetUserByEmail(email)
    if err != nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

func UpgradeUserSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Reference string `json:"reference"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    // Lookup transaction in Paystack (optional, double-check)
    // Then update user record in MongoDB:
    // e.g., set subscription = "premium", hasFreeDownloadLeft = false

    // Example pseudo-code:
    filter := bson.M{"reference": req.Reference}
    update := bson.M{"$set": bson.M{
        "subscription": "premium",
        "hasFreeDownloadLeft": false,
    }}
    _, err := db.Users.UpdateOne(context.TODO(), filter, update)
    if err != nil {
        http.Error(w, "Failed to update subscription", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(map[string]interface{}{
        "success": true,
        "message": "Subscription upgraded",
    })
}
