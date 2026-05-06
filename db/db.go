package db

import (
    "context"
    "time"


    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options" // ✅ this is required
)

var client *mongo.Client
var userCollection *mongo.Collection
var inputCollection *mongo.Collection

func InitDB(uri, dbName string, collections ...string) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
    if err != nil {
        log.Fatal(err)
    }

    for _, col := range collections {
        log.Println("Initialized collection:", col)
    }
}

// Update OTP fields for a user
func UpdateUserOTP(email, otp string, expiry int64) error {
    _, err := userCollection.UpdateOne(
        context.TODO(),
        bson.M{"email": email},
        bson.M{"$set": bson.M{
            "otpCode":   otp,
            "otpExpiry": expiry,
        }},
    )
    return err
}
