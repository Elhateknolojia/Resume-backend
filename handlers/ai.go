package handlers

import (
    // "Backend/models"
    "log"
    "bytes"
    "encoding/json"
    "io"
    "net/http"
	
)

// Process PDF text (rewrite/shorten/expand/ats)
func ProcessPdfTextHandler(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Text  string `json:"text"`
        Action string `json:"action"`
        Model string `json:"model"`
    }
    json.NewDecoder(r.Body).Decode(&req)

    // Forward to Resume LLM service (port 8000)
    body, _ := json.Marshal(req)
    resp, err := http.Post("https://godtekllm.onrender.com/resume/rewrite", "application/json", bytes.NewBuffer(body))
    if err != nil {
        http.Error(w, "Resume AI service error", http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()
    io.Copy(w, resp.Body)
}

// Generate cover letter
func GenerateCoverLetterHandler(w http.ResponseWriter, r *http.Request) {
    var req struct {
        ResumeContent string `json:"resumeContent"`
        Model         string `json:"model"`
    }
    json.NewDecoder(r.Body).Decode(&req)

    body, _ := json.Marshal(req)
    resp, err := http.Post("https://godtekllm.onrender.com/coverletter/generate", "application/json", bytes.NewBuffer(body))
    if err != nil {
        http.Error(w, "Cover Letter AI service error", http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()
    io.Copy(w, resp.Body)
}

// Save blueprint
func SaveBlueprintHandler(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Sections []interface{} `json:"sections"`
    }
    json.NewDecoder(r.Body).Decode(&req)

    // Forward to Resume LLM service (port 8000)
    body, _ := json.Marshal(req)
    resp, err := http.Post("https://godtekllm.onrender.com/save-blueprint", "application/json", bytes.NewBuffer(body))
    if err != nil {
        http.Error(w, "Blueprint save error", http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()
    io.Copy(w, resp.Body)
}

func PolishSummaryHandler(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Summary string `json:"summary"`
        Model   string `json:"model"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    // Call your Python AI service (FastAPI on port 8000)
    body, _ := json.Marshal(req)
    resp, err := http.Post("https://godtekllm.onrender.com/resume/polish-summary", "application/json", bytes.NewBuffer(body))
    if err != nil {
        http.Error(w, "AI service unavailable", http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    var result struct{ Result string `json:"result"` }
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        http.Error(w, "Invalid AI response", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(result)
}



type CoachRequest struct {
    ResumeData json.RawMessage `json:"resumeData"`
    Model      string            `json:"model"`
}

type CoachResponse struct {
    AtsScore     int      `json:"atsScore"`
    Suggestions  []string `json:"suggestions"`
}

func CoachHandler(w http.ResponseWriter, r *http.Request) {
    var req CoachRequest
    bodyBytes, _ := io.ReadAll(r.Body)
    log.Println("Incoming coach request:", string(bodyBytes))

    if err := json.Unmarshal(bodyBytes, &req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }


    body, _ := json.Marshal(req)
    resp, err := http.Post("http://127.0.0.1:8000/resume/coach", "application/json", bytes.NewBuffer(body))
    if err != nil {
        http.Error(w, "Failed to reach AI service", http.StatusBadGateway)
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        http.Error(w, "AI service error", resp.StatusCode)
        return
    }

   data, _ := io.ReadAll(resp.Body)
    w.Header().Set("Content-Type", "application/json")
    w.Write(data)

}
