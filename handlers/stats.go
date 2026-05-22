package handlers

import (
    "encoding/json"
    "net/http"
    "context"
    "time"
	"log"
    "go.mongodb.org/mongo-driver/bson"
    "Backend/db" // ✅ import your db package
)

type Stats struct {
    ActiveUsers   int `json:"activeUsers"`
    JobPlacements int `json:"jobPlacements"`
}


func StatsHandler(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    


    if db.Client == nil {
        http.Error(w, "Database client not initialized", http.StatusInternalServerError)
        return
    }

    userCollection := db.Client.Database("resumeDB").Collection("users")
    if userCollection == nil {
        http.Error(w, "Users collection not found", http.StatusInternalServerError)
        return
    }

    // Count non-admin users
    userCount, err := userCollection.CountDocuments(ctx, bson.M{"isAdmin": false})
    if err != nil {
        log.Println("Error counting users:", err)
        http.Error(w, "Error counting users: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Count job placements
    jobCount, err := userCollection.CountDocuments(ctx, bson.M{
        "isAdmin": false,
        "resumeData.experience": bson.M{"$exists": true, "$ne": nil},
    })
    if err != nil {
        log.Println("Error counting job placements:", err)
        http.Error(w, "Error counting job placements: "+err.Error(), http.StatusInternalServerError)
        return
    }

    stats := Stats{
        ActiveUsers:   int(userCount) + 10000,
        JobPlacements: int(jobCount) + 1000,
    }

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(stats); err != nil {
        log.Println("Error encoding stats:", err)
        http.Error(w, "Error encoding stats", http.StatusInternalServerError)
    }
}
