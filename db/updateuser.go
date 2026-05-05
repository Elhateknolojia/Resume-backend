package db

import (
    "context"
    "time"

    "Backend/models"
    "go.mongodb.org/mongo-driver/bson"
    
)

func UpdateUser(id string, update models.User) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    _, err := userCollection.UpdateOne(
        ctx,
        bson.M{"_id": id},
        bson.M{"$set": update},
    )
    return err
}
