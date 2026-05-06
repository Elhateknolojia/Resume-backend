package handlers

import (
    "crypto/rand"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
        "Backend/db"
    "gopkg.in/gomail.v2"
)

func sendEmailOTP(to string, otp string) error {
    m := gomail.NewMessage()
    m.SetHeader("From", "no-reply@yourdomain.com")
    m.SetHeader("To", to)
    m.SetHeader("Subject", "Your Verification Code")
    m.SetBody("text/plain", fmt.Sprintf("Your OTP code is: %s. It expires in 10 minutes.", otp))

    d := gomail.NewDialer("smtp.gmail.com", 587, "your-email@gmail.com", "your-app-password")
    return d.DialAndSend(m)
}


func SendOTPHandler(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Email string `json:"email"`
    }
    json.NewDecoder(r.Body).Decode(&req)

    user, err := db.GetUserByEmail(req.Email)
    if err != nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    otp := fmt.Sprintf("%06d", rand.Intn(1000000))
    expiry := time.Now().Add(10 * time.Minute).Unix()

    user.OTPCode = otp
    user.OTPExpiry = expiry
    db.UpdateUser(req.Email, user)

     if err := sendEmailOTP(req.Email, otp); err != nil {
        http.Error(w, "Failed to send OTP", http.StatusInternalServerError)
        return
    }
    
    json.NewEncoder(w).Encode(map[string]string{"message": "OTP sent"})
}

func VerifyOTPHandler(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Email string `json:"email"`
        Code  string `json:"code"`
    }
    json.NewDecoder(r.Body).Decode(&req)

    user, err := db.GetUserByEmail(req.Email)
    if err != nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    if user.OTPCode == req.Code && time.Now().Unix() <= user.OTPExpiry {
        json.NewEncoder(w).Encode(map[string]bool{"success": true})
    } else {
        json.NewEncoder(w).Encode(map[string]bool{"success": false})
    }
}

