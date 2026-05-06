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
var userCollection *mongo.Collection
var inputCollection *mongo.Collection

func InitDB(uri, dbName string, collections ...string) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var err error
    client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
    if err != nil {
        log.Fatal(err)
    }

    db := client.Database(dbName)

    // initialize collections explicitly
    userCollection = db.Collection("users")
    inputCollection = db.Collection("userinputs")
    // add more if needed: adminCollection := db.Collection("admin")

    log.Println("Connected to MongoDB:", dbName)
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
