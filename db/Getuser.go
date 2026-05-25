package db

import (
    "context"
    "time"
    "log"
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

    if UserCollection == nil {
        log.Println("ERROR: UserCollection is nil")
        return 0
    }

    count, err := UserCollection.CountDocuments(ctx, bson.M{})
    if err != nil {
        log.Println("CountUsers error:", err)
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
    
func GetStats() map[string]interface{} {
    usersCount := CountUsers()
    adminsCount := CountAdmins()
    freeCount := CountByTier("free")
    premiumCount := CountByTier("premium")
    yearlyCount := CountByTier("1y")
    monthlyCount := CountByTier("1m")

    return map[string]interface{}{
        "users":   usersCount,
        "admins":  adminsCount,
        "free":    freeCount,
        "premium": premiumCount,
        "yearly":  yearlyCount,
        "monthly": monthlyCount,
    }
}
// GetAllUsers returns all users from the collection
func GetAllUsers() ([]models.User, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
    defer cancel()

    cursor, err := UserCollection.Find(ctx, bson.M{})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var users []models.User
    for cursor.Next(ctx) {
        var user models.User
        if err := cursor.Decode(&user); err != nil {
            return nil, err
        }
        users = append(users, user)
    }

    if err := cursor.Err(); err != nil {
        return nil, err
    }

    return users, nil
}
