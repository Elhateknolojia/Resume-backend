package main

import (
    "log"
    "net/http"
    ghandlers "github.com/gorilla/handlers" // alias Gorilla handlers
    "github.com/gorilla/mux"
    "Backend/handlers"                      // your own handlers
    "Backend/middleware"
    "Backend/db"
    "os"

    "github.com/joho/godotenv"   // ✅ add this
    // "Backend/auth"
    // "fmt"
)


func main() {
    // password :="12345678"
    // password1 := "Jesus4Life"
    // hash := auth.HashPassword(password)
    // hash1 := auth.HashPassword(password1)
    // fmt.Println("Hashed password1:",hash1)
    // fmt.Println("Hashed password:",hash)

    // p1:=auth.HashToString(hash1)
    // p :=auth.HashToString(hash)
    // fmt.Println("Hash1", p1)
    // fmt.Println("Hash", p)

      err := godotenv.Load()
    if err != nil {
        log.Println("Warning: .env file not found, relying on system environment variables")
    }

    mongoURI := os.Getenv("MONGO_URI")
    if mongoURI == "" {
        log.Fatal("MONGO_URI not set in environment")
    }

    db.InitDB(mongoURI, "resumeDB", "users", "userinputs", "admin", "local")

    r := mux.NewRouter()
    // Public routes
    r.HandleFunc("/api/auth/signup", handlers.SignupHandler).Methods("POST")
    r.HandleFunc("/api/auth/login", handlers.LoginHandler).Methods("POST")
	r.HandleFunc("/api/admin/login", handlers.AdminLoginHandler).Methods("POST")

    r.HandleFunc("/api/admin/stats", handlers.AdminStatsHandler).Methods("GET")


	r.Handle("/api/resume/check-eligibility", middleware.AuthMiddleware(http.HandlerFunc(handlers.CheckEligibilityHandler))).Methods("POST")



	// AI routes
	r.Handle("/api/ai/process-pdf-text", middleware.AuthMiddleware(http.HandlerFunc(handlers.ProcessPdfTextHandler))).Methods("POST")
	r.Handle("/api/ai/generate-cover-letter", middleware.AuthMiddleware(http.HandlerFunc(handlers.GenerateCoverLetterHandler))).Methods("POST")
	r.Handle("/api/ai/save-blueprint", middleware.AuthMiddleware(http.HandlerFunc(handlers.SaveBlueprintHandler))).Methods("POST")

	r.Handle("/api/ai/polish-summary", middleware.AuthMiddleware(http.HandlerFunc(handlers.PolishSummaryHandler))).Methods("POST")


    // Protected routes
    r.Handle("/api/profile", middleware.AuthMiddleware(http.HandlerFunc(handlers.ProfileHandler))).Methods("GET", "PUT")
    r.Handle("/api/userinput", middleware.AuthMiddleware(middleware.RateLimitMiddleware(http.HandlerFunc(handlers.UserInputHandler)))).Methods("POST")
    r.Handle("/api/response", middleware.AuthMiddleware(middleware.RateLimitMiddleware(http.HandlerFunc(handlers.ResponseHandler)))).Methods("GET")

	r.Handle("/api/resume/download", middleware.AuthMiddleware(middleware.PremiumOnly(http.HandlerFunc(handlers.DownloadResumeHandler)))).Methods("GET")
	r.Handle("/api/coverletter/download", middleware.AuthMiddleware(middleware.PremiumOnly(http.HandlerFunc(handlers.DownloadCoverLetterHandler)))).Methods("GET")

    corsHandler := ghandlers.CORS(
        ghandlers.AllowedOrigins([]string{"http://localhost:3000/"}),
        ghandlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
        ghandlers.AllowedHeaders([]string{"Authorization", "Content-Type"}),
    )(r)


    log.Println("Server running on :8080")
    log.Fatal(http.ListenAndServe(":8080", corsHandler))
}
