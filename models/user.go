package models


type User struct {
    ID       string `json:"id" bson:"_id,omitempty"`
    Name     string `json:"name" bson:"name"`
    Email    string `json:"email" bson:"email"`
    Phone    string `json:"phone" bson:"phone"`
    Address  string `json:"address" bson:"address"`
    Password string `json:"password" bson:"password"` // store stringified hash
    IsAdmin  bool   `json:"isAdmin" bson:"isAdmin"`
    Tier     string `json:"tier" bson:"tier"`
	FreeDownloadsUsed int   `json:"freeDownloadsUsed" bson:"freeDownloadsUsed"`
}

type UserInput struct {
    UserID string `json:"userId" bson:"userId"`
    Text   string `json:"text" bson:"text"`
    Time   int64  `json:"time" bson:"time"`
}