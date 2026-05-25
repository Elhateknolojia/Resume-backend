package db

import (
    "context"
    "time"
    "log"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options" // ✅ this is required
)

var Client *mongo.Client
var UserCollection *mongo.Collection
var InputCollection *mongo.Collection
var JobCollection *mongo.Collection


func InitDB(uri, dbName string, collections ...string) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    clientOpts := options.Client().
        ApplyURI(uri).
        SetMaxPoolSize(100).
        SetMinPoolSize(10).
        SetMaxConnIdleTime(60 * time.Second)

    client, err := mongo.Connect(ctx, clientOpts)
    if err != nil {
        log.Fatal("MongoDB connect error:", err)
    }

    // ✅ Now ping after client is initialized
    if err := client.Ping(ctx, nil); err != nil {
        log.Fatalf("MongoDB connection failed: %v", err)
    }

    Client = client
    db := Client.Database(dbName)
    UserCollection = db.Collection("users")
    InputCollection = db.Collection("userinputs")
    JobCollection = db.Collection("jobs") // ✅ add this line
    log.Println("Connected to MongoDB:", dbName)
}



// Update OTP fields for a user
func UpdateUserOTP(email, otp string, expiry int64) error {
    _, err := UserCollection.UpdateOne(
        context.TODO(),
        bson.M{"email": email},
        bson.M{"$set": bson.M{
            "otpCode":   otp,
            "otpExpiry": expiry,
        }},
    )
    return err
}
