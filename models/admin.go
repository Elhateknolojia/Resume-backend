package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Job struct {
    ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Title    string             `bson:"title" json:"title"`
    Location string             `bson:"location" json:"location"`
    Type     string             `bson:"type" json:"type"`
    Tag      string             `bson:"tag" json:"tag"`
}

type Admin struct {
    ID       string `json:"id" bson:"_id,omitempty"`
    Email    string `json:"email" bson:"email"`
    Name     string  `json:"name" bson:"name"`
    Password string `json:"password" bson:"password"`
    Role     string `json:"role" bson:"role"` // "admin" or "user"
    Tier     string `json:"tier" bson:"tier"` // "free" or "premium"
}
