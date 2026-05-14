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
