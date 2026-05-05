package handlers

import (
	"net/http"
	"encoding/json"
	"Backend/db"
)

// handlers/input.go
func UserInputHandler(w http.ResponseWriter, r *http.Request) {
    var input struct {
        Text string `json:"text"`
    }
    json.NewDecoder(r.Body).Decode(&input)

    userID := r.Context().Value("user_id").(string)
    db.SaveUserInput(userID, input.Text)

    json.NewEncoder(w).Encode(map[string]string{"message": "Input received"})
}
