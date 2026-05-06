package handlers

import (
    "net/http"
    "Backend/db"
    
)

func DownloadResumeHandler(w http.ResponseWriter, r *http.Request) {
    email := r.URL.Query().Get("email")
    if email == "" {
        http.Error(w, "Missing email", http.StatusBadRequest)
        return
    }

    user, err := db.GetUserByEmail(email)
    if err != nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    // If free tier, enforce one free download
    if user.Tier == "free" {
        if user.FreeDownloadsUsed >= 1 {
            http.Error(w, "Free download limit reached", http.StatusForbidden)
            return
        }
        // Increment usage
        user.FreeDownloadsUsed++
        if err := db.UpdateUser(email, user); err != nil {
            http.Error(w, "Failed to update user", http.StatusInternalServerError)
            return
        }
    }

    // TODO: generate or fetch the actual PDF file
    w.Header().Set("Content-Type", "application/pdf")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("%PDF-1.4\n%...mock resume pdf...")) // replace with real PDF bytes
}
