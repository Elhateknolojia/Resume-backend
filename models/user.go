package models

type User struct {
    ID       string `json:"id"`
    Name     string `json:"name"`
    Email    string `json:"email"`
    Phone    string `json:"phone"`
    Address  string `json:"address"`
    Password string `json:"-"` // stored hashed, not exposed
    IsAdmin  bool   `json:"isAdmin"`
}

type UserInput struct {
    UserID string `json:"userId" bson:"userId"`
    Text   string `json:"text" bson:"text"`
    Time   int64  `json:"time" bson:"time"`
}