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
        Email string `json:"email"`
        Tier  string `json:"tier"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    filter := bson.M{"email": req.Email}
    update := bson.M{"$set": bson.M{
        "tier": req.Tier,             // change free → premium/2w/1m/1y
        "freeDownloadsUsed": 0,       // reset free counter
    }}

    _, err := db.UserCollection.UpdateOne(context.TODO(), filter, update)
    if err != nil {
        http.Error(w, "Failed to update subscription", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(map[string]interface{}{
        "success": true,
        "message": "Tier upgraded successfully",
    })
}
