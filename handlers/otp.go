package handlers

import (
    "crypto/rand"
    "encoding/json"
    "fmt"
    "math/big"
    "net/http"
    "os"
    "time"
    "Backend/db"
    "log"
    "Backend/auth"
    "Backend/models"
    "gopkg.in/gomail.v2"
)

// generateOTP creates a cryptographically secure 6-digit code
func generateOTP() string {
    n, _ := rand.Int(rand.Reader, big.NewInt(1000000))
    return fmt.Sprintf("%06d", n.Int64())
}

func sendEmailOTP(to string, otp string) error {
    m := gomail.NewMessage()
    m.SetHeader("From", os.Getenv("SMTP_USER"))
    m.SetHeader("To", to)
    m.SetHeader("Subject", "Your Verification Code")
    m.SetBody("text/plain", fmt.Sprintf("Your OTP code is: %s. It expires in 10 minutes.", otp))

    d := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("SMTP_USER"), os.Getenv("SMTP_PASS"))
    return d.DialAndSend(m)
}

func SendOTPHandler(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Name     string `json:"name"`
        Email    string `json:"email"`
        Phone    string `json:"phone"`
        Address  string `json:"address"`
        Role     string `json:"role"`
        Password string `json:"password"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    otp := generateOTP()
    expiry := time.Now().Add(10 * time.Minute).Unix()

    hashedPassword := auth.HashPassword(req.Password)

    user := &models.User{
        Name:      req.Name,
        Email:     req.Email,
        Phone:     req.Phone,
        Address:   req.Address,
        Tier:      "free",
        OTPCode:   otp,
        OTPExpiry: expiry,
        IsAdmin:   false,
        Password:  hashedPassword,
    }

    if err := db.CreateOrUpdateUser(user); err != nil {
        log.Printf("DB save error: %v", err)
        http.Error(w, "Failed to save user", http.StatusInternalServerError)
        return
    }

    if err := sendEmailOTP(req.Email, otp); err != nil {
        log.Printf("email send error: %v", err)
        http.Error(w, "Failed to send OTP", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"message": "OTP sent"})
}

func VerifyOTPHandler(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Email string `json:"email"`
        Code  string `json:"code"`
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

    success := user.OTPCode == req.Code && time.Now().Unix() <= user.OTPExpiry
    if success {
    user.OTPCode = "VERIFIED"
    user.OTPExpiry = 0
    user.Tier = "free"
    db.CreateOrUpdateUser(user)
}


    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]bool{"success": success})
}


