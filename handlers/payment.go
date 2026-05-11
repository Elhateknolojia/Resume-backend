package handlers

import (
    "encoding/json"
    "net/http"
    "os"
    "bytes"
    "io/ioutil"
    "fmt"
    // "log"
)

type InitRequest struct {
    Email   string `json:"email"`
    TierId  string `json:"tierId"`
    Amount  int    `json:"amount"`
}

type PaystackResponse struct {
    Status  bool   `json:"status"`
    Message string `json:"message"`
    Data    struct {
        AuthorizationURL string `json:"authorization_url"`
        AccessCode       string `json:"access_code"`
        Reference        string `json:"reference"`
    } `json:"data"`
}

func InitiatePaymentHandler(w http.ResponseWriter, r *http.Request) {
    var req InitRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    // Paystack expects amount in kobo (KES cents)
    payload := map[string]interface{}{
    "email": req.Email,
    "amount": req.Amount * 100,
    "callback_url": fmt.Sprintf("http://localhost:3000/payment/callback?tierId=%s", req.TierId),
    "metadata": map[string]string{
        "tierId": req.TierId,
    },
}


    body, _ := json.Marshal(payload)

    paystackSecret := os.Getenv("PAYSTACK_SECRET_KEY")
    reqPaystack, _ := http.NewRequest("POST", "https://api.paystack.co/transaction/initialize", bytes.NewBuffer(body))
    reqPaystack.Header.Set("Authorization", "Bearer "+paystackSecret)
    reqPaystack.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(reqPaystack)
    if err != nil {
        http.Error(w, "Failed to reach Paystack", http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    respBody, _ := ioutil.ReadAll(resp.Body)

    var payResp PaystackResponse
    if err := json.Unmarshal(respBody, &payResp); err != nil {
        http.Error(w, "Failed to parse Paystack response", http.StatusInternalServerError)
        return
    }

    if !payResp.Status {
        http.Error(w, payResp.Message, http.StatusBadRequest)
        return
    }

    // Return authorization_url to frontend
    json.NewEncoder(w).Encode(map[string]interface{}{
        "success": true,
        "authorization_url": payResp.Data.AuthorizationURL,
        "reference": payResp.Data.Reference,
    })
}

func VerifyPaymentHandler(w http.ResponseWriter, r *http.Request) {
    reference := r.URL.Query().Get("reference")
    if reference == "" {
        http.Error(w, "Missing reference", http.StatusBadRequest)
        return
    }

    paystackSecret := os.Getenv("PAYSTACK_SECRET_KEY")
    url := "https://api.paystack.co/transaction/verify/" + reference

    reqPaystack, err := http.NewRequest("GET", url, nil)
    if err != nil {
        http.Error(w, "Failed to create request", http.StatusInternalServerError)
        return
    }
    reqPaystack.Header.Set("Authorization", "Bearer "+paystackSecret)

    client := &http.Client{}
    resp, err := client.Do(reqPaystack)
    if err != nil {
        http.Error(w, "Failed to reach Paystack", http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    respBody, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        http.Error(w, "Failed to read response", http.StatusInternalServerError)
        return
    }

    var payResp PaystackResponse
    if err := json.Unmarshal(respBody, &payResp); err != nil {
        http.Error(w, "Failed to parse Paystack response", http.StatusInternalServerError)
        return
    }

    if !payResp.Status {
        http.Error(w, payResp.Message, http.StatusBadRequest)
        return
    }

    // At this point, payment is verified
    json.NewEncoder(w).Encode(map[string]interface{}{
        "success": true,
        "reference": payResp.Data.Reference,
        "message": payResp.Message,
    })
}
