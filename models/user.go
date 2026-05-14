package models


type User struct {
    ID       string `json:"id" bson:"_id,omitempty"`
    Name     string `json:"name" bson:"name"`
    Email    string `json:"email" bson:"email"`
    Phone    string `json:"phone" bson:"phone"`
    Address  string `json:"address" bson:"address"`
    Password string `json:"password" bson:"password"` // hashed
    IsAdmin  bool   `json:"isAdmin" bson:"isAdmin"`
    Tier     string `json:"tier" bson:"tier"`

    FreeDownloadsUsed int    `json:"freeDownloadsUsed" bson:"freeDownloadsUsed"`
    PendingResume     string `json:"pendingResume,omitempty" bson:"pendingResume,omitempty"`
    OTPCode           string `json:"otpCode" bson:"otpCode"`
    OTPExpiry         int64  `json:"otpExpiry" bson:"otpExpiry"`

    // 🔹 New fields for resume ecosystem
    ResumeData       ResumeData   `json:"resumeData" bson:"resumeData"`
    ImportedPdfData  []byte       `json:"importedPdfData,omitempty" bson:"importedPdfData,omitempty"` // raw PDF bytes or reference
    CoverLetterData  string       `json:"coverLetterData,omitempty" bson:"coverLetterData,omitempty"`
}


type UserInput struct {
    UserID string `json:"userId" bson:"userId"`
    Text   string `json:"text" bson:"text"`
    Time   int64  `json:"time" bson:"time"`
}

type LoginResponse struct {
    Token   string `json:"token"`
    UserID  string `json:"user_id"`
    Email   string `json:"email"`
    IsAdmin bool   `json:"isAdmin"`
    Tier    string `json:"tier"`
}
