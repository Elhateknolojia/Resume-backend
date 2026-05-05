package handlers

import (
    "net/http"
    "fmt"
)

func DownloadResumeHandler(w http.ResponseWriter, r *http.Request) {
    tier := r.Context().Value("tier").(string)
    if tier != "premium" {
        http.Error(w, "Free users cannot download resumes", http.StatusForbidden)
        return
    }
    // Example: serve file or generate PDF
    fmt.Fprintln(w, "Resume PDF download started...")
}

func DownloadCoverLetterHandler(w http.ResponseWriter, r *http.Request) {
    tier := r.Context().Value("tier").(string)
    if tier != "premium" {
        http.Error(w, "Free users cannot download cover letters", http.StatusForbidden)
        return
    }
    fmt.Fprintln(w, "Cover Letter PDF download started...")
}
