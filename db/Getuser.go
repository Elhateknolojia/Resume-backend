package db

import (
    "context"
    "time"

    "Backend/models"
    "go.mongodb.org/mongo-driver/bson"
    
)

func GetUserByEmail(email string) (models.User, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    var user models.User
    err := userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
    return user, err
}

func GetUserByID(id string) (models.User, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    var user models.User
    err := userCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
    return user, err
}
