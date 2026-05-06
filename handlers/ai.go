package handlers

import (
    "bytes"
    "encoding/json"
    "io"
    "net/http"
	"fmt"
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
    resp, err := http.Post("http://localhost:8000/rewrite", "application/json", bytes.NewBuffer(body))
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
    resp, err := http.Post("http://localhost:8001/generate", "application/json", bytes.NewBuffer(body))
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
    resp, err := http.Post("http://localhost:8000/save-blueprint", "application/json", bytes.NewBuffer(body))
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
    resp, err := http.Post("http://localhost:8000/polish-summary", "application/json",
        bytes.NewBuffer([]byte(fmt.Sprintf(`{"summary":"%s","model":"%s"}`, req.Summary, req.Model))))
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

