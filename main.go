package main

import (
	"Backend/db"
	"Backend/handlers" // your own handlers
	"Backend/middleware"
	"log"
	"net/http"
	"os"

	ghandlers "github.com/gorilla/handlers" // alias Gorilla handlers
	"github.com/gorilla/mux"
	// "golang.org/x/oauth2"
	// "golang.org/x/oauth2/google"

	// "github.com/rs/cors"
	"github.com/joho/godotenv" // ✅ add this
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
    log.Printf("DEBUG MONGO_URI: [%s]", mongoURI)
    if mongoURI == "" {
        log.Fatal("MONGO_URI not set in environment")
    }
    db.InitDB(mongoURI, "resumeDB", "users", "userinputs", "jobs", "admin", "local")



    // in main.go
    // Replace this line:
// http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

// With this:
   



    r := mux.NewRouter()


    // Handle OPTIONS requests globally
    r.Methods(http.MethodOptions).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
    })
    // Public routes
    // r.HandleFunc("/api/pdf/import", handlers.ImportPdfHandler).Methods("POST")
    // PDF routes
    r.HandleFunc("/api/pdf/import", handlers.ImportPdfHandler).Methods("POST")
    r.HandleFunc("/api/pdf/save", handlers.SavePdfHandler).Methods("POST")


    r.HandleFunc("/api/docx/import", handlers.ImportDocxHandler).Methods("POST")
    r.HandleFunc("/api/docx/save", handlers.SaveDocxHandler).Methods("POST")


r.HandleFunc("/api/auth/google/url", handlers.GoogleLoginURLHandler).Methods("GET")
r.HandleFunc("/api/auth/google/callback", handlers.GoogleCallbackHandler).Methods("GET")

    // // PDF Reconstruction (Export)
    // r.HandleFunc("/api/pdf/reconstruct", handlers.ReconstructPdfHandler).Methods("POST")
     // Payment routes
    r.HandleFunc("/api/payment/initiate", handlers.InitiatePaymentHandler).Methods("POST")
    r.HandleFunc("/api/payment/verify", handlers.VerifyPaymentHandler).Methods("GET")


    //static file server for generated PDFs
     r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

    r.HandleFunc("/api/auth/signup", handlers.SignupHandler).Methods("POST")
    r.HandleFunc("/api/auth/login", handlers.LoginHandler).Methods("POST")
	r.HandleFunc("/api/admin/login", handlers.AdminLoginHandler).Methods("POST")

    r.HandleFunc("/api/auth/session", handlers.SessionHandler).Methods("GET")
    r.HandleFunc("/api/auth/logout", handlers.LogoutHandler).Methods("POST")


    r.HandleFunc("/api/admin/stats", handlers.AdminStatsHandler).Methods("GET")

    r.HandleFunc("/plans", handlers.PlansHandler).Methods("GET")

    r.HandleFunc("/api/stats", handlers.StatsHandler).Methods("GET")


	r.HandleFunc("/api/resume/check-eligibility", handlers.CheckEligibilityHandler).Methods("POST")
    r.Handle("/api/resume/export-html", middleware.AuthMiddleware(http.HandlerFunc(handlers.ExportHtmlHandler))).Methods("POST")

    r.Handle("/api/auth/refresh", middleware.AuthMiddleware(http.HandlerFunc(handlers.RefreshHandler))).Methods("POST")

    r.HandleFunc("/api/support", handlers.SupportHandler).Methods("POST")


	// AI routes
	r.Handle("/api/ai/process-pdf-text", middleware.AuthMiddleware(http.HandlerFunc(handlers.ProcessPdfTextHandler))).Methods("POST")
	r.Handle("/api/ai/generate-cover-letter", middleware.AuthMiddleware(http.HandlerFunc(handlers.GenerateCoverLetterHandler))).Methods("POST")
	r.Handle("/api/ai/save-blueprint", middleware.AuthMiddleware(http.HandlerFunc(handlers.SaveBlueprintHandler))).Methods("POST")

	r.Handle("/api/ai/polish-summary", middleware.AuthMiddleware(http.HandlerFunc(handlers.PolishSummaryHandler))).Methods("POST")
    r.Handle("/api/ai/coach", middleware.AuthMiddleware(http.HandlerFunc(handlers.CoachHandler))).Methods("POST")

    // Template routes
    r.HandleFunc("/api/templates/{id}", handlers.LoadTemplateHandler).Methods("GET")

    // OTP routes
    r.HandleFunc("/api/auth/send-otp", handlers.SendOTPHandler).Methods("POST")
    r.HandleFunc("/api/auth/verify-otp", handlers.VerifyOTPHandler).Methods("POST")
    
    r.HandleFunc("/api/payment/upgrade", handlers.UpgradeUserSubscriptionHandler).Methods("POST")

    // User routes
    r.Handle("/api/user", middleware.AuthMiddleware(http.HandlerFunc(handlers.GetUserHandler))).Methods("GET")
    r.Handle("/api/resume/me", middleware.AuthMiddleware(http.HandlerFunc(handlers.GetResumeHandler))).Methods("GET")


    r.HandleFunc("/coverletter/generate", handlers.GenerateCoverLetterHandler).Methods("POST")
    r.HandleFunc("/coverletter/suggestions", handlers.ImprovementSuggestionsHandler).Methods("POST")


// ✅ New Job routes with Gorilla Mux
r.HandleFunc("/api/jobs", handlers.GetJobs).Methods("GET")
r.HandleFunc("/api/jobs", handlers.CreateJob).Methods("POST")
r.HandleFunc("/api/jobs/{id}", handlers.DeleteJob).Methods("DELETE")


    // Protected routes
    r.Handle("/api/profile", middleware.AuthMiddleware(http.HandlerFunc(handlers.ProfileHandler))).Methods("GET", "PUT")
    r.Handle("/api/userinput", middleware.AuthMiddleware(middleware.RateLimitMiddleware(http.HandlerFunc(handlers.UserInputHandler)))).Methods("POST")
    r.Handle("/api/response", middleware.AuthMiddleware(middleware.RateLimitMiddleware(http.HandlerFunc(handlers.ResponseHandler)))).Methods("GET")
    
	r.Handle("/api/resume/generate-pdf", middleware.AuthMiddleware(http.HandlerFunc(handlers.GeneratePdfHandler))).Methods("POST")
    r.HandleFunc("/api/resume/save-pending", handlers.SavePendingResumeHandler).Methods("POST")

    
	// r.Handle("/api/coverletter/download", middleware.AuthMiddleware(middleware.PremiumOnly(http.HandlerFunc(handlers.DownloadCoverLetterHandler)))).Methods("GET")

corsHandler := ghandlers.CORS(
    ghandlers.AllowedOrigins([]string{
        // "http://localhost:4200",              // Angular dev
        // "http://localhost:4000",
        // "http://localhost:2000",
        // "http://localhost:3000",              // if you test on 3000
         "https://resume.elitesuites.top",
         "https://www.elitesuites.top",
         "https://elitesuites.top",
        "https://pdfeditor.elitesuites.top",
        "https://coverletter.elitesuites.top",
        "https://jobsearch.elitesuites.top",
        "https://resume-six-dun.vercel.app",
        "https://www.elitesuites.top/",
        "https://resumebuilder-pdfeditor.onrender.com",
        "https://coverletter-1-sbiz.onrender.com",
        "https://resume-backend-plmv.onrender.com",
        // "https://resume-backend-weld.vercel.app/",
        // "*",
    }),
    ghandlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
    ghandlers.AllowedHeaders([]string{"Authorization", "Content-Type"}),
    ghandlers.AllowCredentials(),
)(r)




    log.Println("Server running on :8080")
    log.Fatal(http.ListenAndServe(":8080", corsHandler))
}
