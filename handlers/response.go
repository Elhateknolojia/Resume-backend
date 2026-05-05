package handlers

// handlers/response.go
func ResponseHandler(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value("user_id").(string)

    // Example: call FastAPI resume service
    resp, err := http.Post("http://127.0.0.1:8000/rewrite", "application/json",
        bytes.NewBuffer([]byte(`{"input":"sample text"}`)))
    if err != nil {
        http.Error(w, "AI service error", http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    var result map[string]string
    json.NewDecoder(resp.Body).Decode(&result)

    json.NewEncoder(w).Encode(map[string]interface{}{
        "user":     userID,
        "response": result["output"],
    })
}
