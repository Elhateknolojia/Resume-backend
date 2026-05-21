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

    userCollection := db.Client.Database("resumeDB").Collection("users")

    // Count non-admin users
    userCount, err := userCollection.CountDocuments(ctx, bson.M{"isAdmin": false})
    if err != nil {
        http.Error(w, "Error counting users: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Count job placements: users with at least one experience entry
    jobCount, err := userCollection.CountDocuments(ctx, bson.M{
        "isAdmin": false,
        "resumeData.experience": bson.M{"$ne": nil},
    })
    if err != nil {
		log.Println("StatsHandler error:", err)
        http.Error(w, "Error counting job placements: "+err.Error(), http.StatusInternalServerError)
        return
    }

    stats := Stats{
        ActiveUsers:   int(userCount) + 10000,
        JobPlacements: int(jobCount) + 1000,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(stats)
}
