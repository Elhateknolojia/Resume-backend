package models

type Admin struct {
    ID       string `json:"id" bson:"_id,omitempty"`
    Email    string `json:"email" bson:"email"`
    Password string `json:"password" bson:"password"`
    Role     string `json:"role" bson:"role"` // "admin" or "user"
    Tier     string `json:"tier" bson:"tier"` // "free" or "premium"
}
