package db

import (
    "context"
    "time"
    "go.mongodb.org/mongo-driver/bson"
    "Backend/models"
    "go.mongodb.org/mongo-driver/mongo/options"
    
)

func CreateUser(user models.User) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    _, err := userCollection.InsertOne(ctx, user)
    return err
}


func CreateOrUpdateUser(user *models.User) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    _, err := userCollection.UpdateOne(
        ctx,
        bson.M{"email": user.Email},
        bson.M{"$set": user},
        options.Update().SetUpsert(true),
    )
    return err
}
