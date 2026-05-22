package handlers

import (
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "strings"

    "github.com/go-resty/resty/v2"
)

type SupportTicket struct {
    Name    string `json:"name"`
    Email   string `json:"email"`
    Type    string `json:"type"`
    Message string `json:"message"`
    IsGuest bool   `json:"isGuest"`
}

func SupportHandler(w http.ResponseWriter, r *http.Request) {
    var ticket SupportTicket
    if err := json.NewDecoder(r.Body).Decode(&ticket); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    subject := fmt.Sprintf("[Support Request] %s", ticket.Type)
    body := fmt.Sprintf("From: %s (%s)\n\nType: %s\n\nMessage:\n%s",
        ticket.Name, ticket.Email, ticket.Type, ticket.Message)

    if err := sendMailgunEmail(subject, body); err != nil {
        http.Error(w, "Failed to send support request", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "status":  "success",
        "message": "Support request submitted",
    })
}

func sendMailgunEmail(subject, body string) error {
    apiKey := os.Getenv("MAILGUN_API_KEY")
    domain := os.Getenv("MAILGUN_DOMAIN")
    baseURL := os.Getenv("MAILGUN_BASE_URL")
    supportEmail := os.Getenv("SUPPORT_EMAIL") // e.g. support@elitesuites.com
    adminEmail := os.Getenv("ADMIN_EMAIL")     // optional

    client := resty.New()
    toList := []string{supportEmail}
    if adminEmail != "" {
        toList = append(toList, adminEmail)
    }

    _, err := client.R().
        SetBasicAuth("api", apiKey).
        SetFormData(map[string]string{
            "from":    "Elitesuites Support <no-reply@" + domain + ">",
            "to":      strings.Join(toList, ","),
            "subject": subject,
            "text":    body,
        }).
        Post(fmt.Sprintf("%s/v3/%s/messages", baseURL, domain))

    return err
}
