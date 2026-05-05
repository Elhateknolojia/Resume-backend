package db

func CreateUser(user models.User) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    _, err := userCollection.InsertOne(ctx, user)
    return err
}
