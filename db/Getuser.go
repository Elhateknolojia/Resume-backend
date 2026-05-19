package db

import (
    "context"
    "time"

    "Backend/models"
    "go.mongodb.org/mongo-driver/bson"
)

// Get user by email with timeout
func GetUserByEmail(email string) (*models.User, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
    defer cancel()

    var user models.User
    err := UserCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
    if err != nil {
        return nil, err
    }
    return &user, nil
}

// Get user by ID with timeout
func GetUserByID(id string) (*models.User, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
    defer cancel()

    var user models.User
    err := UserCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
    if err != nil {
        return nil, err
    }
    return &user, nil
}

// Count all users
func CountUsers() int64 {
    ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
    defer cancel()

    count, err := UserCollection.CountDocuments(ctx, bson.M{})
    if err != nil {
        return 0
    }
    return count
}

// Count all admins
func CountAdmins() int64 {
    ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
    defer cancel()

    count, err := UserCollection.CountDocuments(ctx, bson.M{"isAdmin": true})
    if err != nil {
        return 0
    }
    return count
}

// Count users by tier (e.g. "free", "premium")
func CountByTier(tier string) int64 {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    count, err := UserCollection.CountDocuments(ctx, bson.M{"tier": tier})
    if err != nil {
        return 0
    }
    return count
}
    