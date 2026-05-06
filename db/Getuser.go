package db

import (
    "context"

    "Backend/models"
    "go.mongodb.org/mongo-driver/bson"
    
)

func GetUserByEmail(email string) (*models.User, error) {
    var user models.User
    err := userCollection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
    if err != nil {
        return nil, err
    }
    return &user, nil
}


func GetUserByID(id string) (*models.User, error) {
    var user models.User
    err := userCollection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&user)
    if err != nil {
        return nil, err
    }
    return &user, nil
}


// Count all users
func CountUsers() int64 {
    count, err := userCollection.CountDocuments(context.TODO(), bson.M{})
    if err != nil {
        return 0
    }
    return count
}

// Count all admins
func CountAdmins() int64 {
    count, err := userCollection.CountDocuments(context.TODO(), bson.M{"isAdmin": true})
    if err != nil {
        return 0
    }
    return count
}

// Count users by tier (e.g. "free", "premium")
func CountByTier(tier string) int64 {
    count, err := userCollection.CountDocuments(context.TODO(), bson.M{"tier": tier})
    if err != nil {
        return 0
    }
    return count
}