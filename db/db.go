package db

import (
    "context"
    "time"
    "log"
    
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options" // ✅ this is required
)

var client *mongo.Client
var UserCollection *mongo.Collection
var InputCollection *mongo.Collection

func InitDB(uri, dbName string, collections ...string) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    clientOpts := options.Client().
        ApplyURI(uri).
        SetMaxPoolSize(100).        // ✅ allow up to 50 concurrent connections
        SetMinPoolSize(10).         // ✅ keep at least 5 warm connections
        SetMaxConnIdleTime(20 * time.Second) // ✅ recycle idle connections

    client, err:= mongo.Connect(ctx, clientOpts)
    if err != nil {
        log.Fatal(err)
    }

    db := client.Database(dbName)
    UserCollection = db.Collection("users")
    InputCollection = db.Collection("userinputs")

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
