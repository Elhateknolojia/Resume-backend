package db

import (
    "context"
    

    "Backend/models"
    "go.mongodb.org/mongo-driver/bson"
    
)

func UpdateUser(email string, user *models.User) error {
    _, err := userCollection.UpdateOne(
        context.TODO(),
        bson.M{"email": email},
        bson.M{"$set": user},
    )
    return err
}
