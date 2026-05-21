package handlers

import (
    "encoding/json"
    "net/http"
    "context"
    "time"

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

    // Count all non-admin users
    userCount, err := userCollection.CountDocuments(ctx, bson.M{"isAdmin": false})
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Example: job placements derived from resumes with experience entries
    resumeCollection := db.Client.Database("resumeDB").Collection("resumes")
    jobCount, err := resumeCollection.CountDocuments(ctx, bson.M{"experience.0": bson.M{"$exists": true}})
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    stats := Stats{
        ActiveUsers:   int(userCount) + 10000, // offset
        JobPlacements: int(jobCount) + 1000,   // offset
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(stats)
}

