// handlers/careers.go

package handlers

import (
	"Backend/db"
	"Backend/models"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetJobs fetches all jobs
func GetJobs(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    cursor, err := db.JobCollection.Find(ctx, bson.M{})
    if err != nil {
        http.Error(w, "Failed to fetch jobs", http.StatusInternalServerError)
        return
    }
    defer cursor.Close(ctx)

    var jobs []models.Job
    if err = cursor.All(ctx, &jobs); err != nil {
        http.Error(w, "Failed to decode jobs", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(jobs)
}

// CreateJob inserts a new job
func CreateJob(w http.ResponseWriter, r *http.Request) {
    var job models.Job
    if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    job.ID = primitive.NewObjectID()

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    _, err := db.JobCollection.InsertOne(ctx, job)
    if err != nil {
        http.Error(w, "Failed to insert job", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(job)
}

// DeleteJob removes a job by ID
func DeleteJob(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        http.Error(w, "Invalid job ID", http.StatusBadRequest)
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    res, err := db.JobCollection.DeleteOne(ctx, bson.M{"_id": objID})
    if err != nil {
        http.Error(w, "Failed to delete job", http.StatusInternalServerError)
        return
    }
    if res.DeletedCount == 0 {
        http.Error(w, "Job not found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"message": "Job deleted successfully"})
}
