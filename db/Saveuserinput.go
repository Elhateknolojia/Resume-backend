package db


import (
    "context"
    "time"

    "Backend/models"
    
)

func SaveUserInput(userID string, text string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    input := models.UserInput{
        UserID: userID,
        Text:   text,
        Time:   time.Now().Unix(),
    }
    _, err := inputCollection.InsertOne(ctx, input)
    return err
}
