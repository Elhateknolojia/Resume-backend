package handlers


// handlers/input.go
func UserInputHandler(w http.ResponseWriter, r *http.Request) {
    var input struct {
        Text string `json:"text"`
    }
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    // Save input to DB or forward to AI service
    db.SaveUserInput(r.Context().Value("user_id").(string), input.Text)

    json.NewEncoder(w).Encode(map[string]string{"message": "Input received"})
}
