package handlers

import (
    "encoding/json"
    "net/http"
    "Backend/db"
)

func CheckEligibilityHandler(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Email string `json:"email"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    user, err := db.GetUserByEmail(req.Email)
    if err != nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    // Default flags
    canDownload := false
    isPremium := false
    hasFreeDownloadLeft := false

    // Logic based on tier
    switch user.Tier {
    case "premium":
        canDownload = true
        isPremium = true
        hasFreeDownloadLeft = false // premium users don’t need free downloads
    case "free":
        // Example: allow one free download if not used yet
        if user.FreeDownloadsUsed < 1 {
            canDownload = true
            hasFreeDownloadLeft = true
        } else {
            canDownload = false
            hasFreeDownloadLeft = false
        }
    }

    json.NewEncoder(w).Encode(map[string]interface{}{
        "canDownload":        canDownload,
        "isPremium":          isPremium,
        "hasFreeDownloadLeft": hasFreeDownloadLeft,
    })
}
