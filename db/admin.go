package db

import (
    "context"
    "time"

    "Backend/models"
    "go.mongodb.org/mongo-driver/bson"
)

func GetAdminByEmail(email string) (models.Admin, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    var admin models.Admin
    err := UserCollection.FindOne(ctx, bson.M{"email": email, "role": "admin"}).Decode(&admin)
    return admin, err
}


func GetStats() map[string]interface{} {
    usersCount := CountUsers()
    adminsCount := CountAdmins()
    premiumCount := CountByTier("premium")
    freeCount := CountByTier("free")

    return map[string]interface{}{
        "users":   usersCount,
        "admins":  adminsCount,
        "premium": premiumCount,
        "free":    freeCount,
    }
}