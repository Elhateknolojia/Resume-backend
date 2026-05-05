package main

import (
    "log"
    "net/http"

    "github.com/gorilla/mux"
    "Backend/handlers"
    "Backend/middleware"
    "Backend/db"
)

func main() {
    db.InitDB("mongodb://localhost:27017", "resumeDB")

    r := mux.NewRouter()

    // Public routes
    r.HandleFunc("/api/auth/signup", handlers.SignupHandler).Methods("POST")
    r.HandleFunc("/api/auth/login", handlers.LoginHandler).Methods("POST")
	r.HandleFunc("/api/admin/login", handlers.AdminLoginHandler).Methods("POST")



	r.Handle("/api/resume/check-eligibility", middleware.AuthMiddleware(http.HandlerFunc(handlers.CheckEligibilityHandler))).Methods("POST")



	// AI routes
	r.Handle("/api/ai/process-pdf-text", middleware.AuthMiddleware(http.HandlerFunc(handlers.ProcessPdfTextHandler))).Methods("POST")
	r.Handle("/api/ai/generate-cover-letter", middleware.AuthMiddleware(http.HandlerFunc(handlers.GenerateCoverLetterHandler))).Methods("POST")
	r.Handle("/api/ai/save-blueprint", middleware.AuthMiddleware(http.HandlerFunc(handlers.SaveBlueprintHandler))).Methods("POST")


    // Protected routes
    r.Handle("/api/profile", middleware.AuthMiddleware(http.HandlerFunc(handlers.ProfileHandler))).Methods("GET", "PUT")
    r.Handle("/api/userinput", middleware.AuthMiddleware(middleware.RateLimitMiddleware(http.HandlerFunc(handlers.UserInputHandler)))).Methods("POST")
    r.Handle("/api/response", middleware.AuthMiddleware(middleware.RateLimitMiddleware(http.HandlerFunc(handlers.ResponseHandler)))).Methods("GET")

	r.Handle("/api/resume/download", middleware.AuthMiddleware(middleware.PremiumOnly(http.HandlerFunc(handlers.DownloadResumeHandler)))).Methods("GET")
	r.Handle("/api/coverletter/download", middleware.AuthMiddleware(middleware.PremiumOnly(http.HandlerFunc(handlers.DownloadCoverLetterHandler)))).Methods("GET")


    log.Println("Server running on :8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}
