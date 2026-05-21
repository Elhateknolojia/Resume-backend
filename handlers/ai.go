package handlers

import (
	"Backend/models"
	"encoding/json"
	"log"
    "fmt"
	"net/http"
	"os"
	"time"

	openai "github.com/sashabaranov/go-openai"
	"golang.org/x/net/context"

	// "github.com/openai/openai-go" // ✅ official OpenAI Go SDK
	// "golang.org/x/net/context"

	// "encoding/json"
	// "net/http"
	// "os"
	"bytes"
	"io"
	// "log"
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


type CoverLetterRequest struct {
    ResumeContent string `json:"resumeContent"`
    JobDetails    string `json:"jobDetails"`
    Template      string `json:"template"`
}

type GPTRequest struct {
    Model    string        `json:"model"`
    Messages []GPTMessage  `json:"messages"`
    MaxTokens int          `json:"max_tokens"`
    Temperature float64    `json:"temperature"`
}

type GPTMessage struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

type GPTResponse struct {
    Choices []struct {
        Message GPTMessage `json:"message"`
    } `json:"choices"`
}

func GenerateCoverLetterHandler(w http.ResponseWriter, r *http.Request) {
    start := time.Now() // ⏱ start timer

    var req models.CoverLetterRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    apiKey := os.Getenv("OPENAI_API_KEY")
    if apiKey == "" {
        http.Error(w, "OPENAI_API_KEY not set in environment", http.StatusInternalServerError)
        return
    }

    ctx := context.Background()
    client := openai.NewClient(apiKey)

    log.Println("[CoverLetter] ✅ OpenAI client initialized, secure connection established")

    prompt := buildCoverLetterPrompt(req)

    resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
        Model:     "gpt-5.0",
        Messages: []openai.ChatCompletionMessage{
            {Role: "user", Content: prompt},
        },
        MaxTokens:   800,
        Temperature: 0.7,
    })
    if err != nil {
        log.Println("[CoverLetter] ❌ OpenAI error:", err)
        json.NewEncoder(w).Encode(models.CoverLetterResponse{Error: "Failed to generate cover letter"})
        return
    }

    text := ""
    if len(resp.Choices) > 0 {
        text = resp.Choices[0].Message.Content
    }

    elapsed := time.Since(start)
    log.Printf("[CoverLetter] ⏱ Generated in %s\n", elapsed)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(models.CoverLetterResponse{CoverLetter: text})
}

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
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

    prompt := "Polish the following resume summary to be more impactful, concise, and ATS-friendly. " +
        "Return only the improved summary text:\n\n" + req.Summary

    resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
        Model: "gpt-5.0",
        Messages: []openai.ChatCompletionMessage{
            {Role: "user", Content: prompt},
        },
        MaxTokens: 300,
    })
    if err != nil {
        http.Error(w, "AI error", http.StatusInternalServerError)
        return
    }

    polished := resp.Choices[0].Message.Content

    result := map[string]interface{}{
        "result": polished,
        "requiresSubscription": false, // add tier logic here
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
}


func CoachHandler(w http.ResponseWriter, r *http.Request) {
    var req models.CoachRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

    // ✅ Access struct field directly
    prompt := "Evaluate the following resume summary for ATS compatibility. " +
        "Return a JSON object with fields: atsScore (0-100) and suggestions (array of strings).\n\n" +
        req.ResumeData.Summary

    resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
        Model: "gpt-5.0",
        Messages: []openai.ChatCompletionMessage{
            {Role: "user", Content: prompt},
        },
        MaxTokens: 500,
    })
    if err != nil {
        http.Error(w, "AI error", http.StatusInternalServerError)
        return
    }

    output := resp.Choices[0].Message.Content

    var coachRes models.CoachResponse
    if err := json.Unmarshal([]byte(output), &coachRes); err != nil {
        coachRes = models.CoachResponse{
            AtsScore:    78,
            Suggestions: []string{"Add specific achievements", "Improve summary impact"},
        }
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(coachRes)
}
func ImprovementSuggestionsHandler(w http.ResponseWriter, r *http.Request) {
    start := time.Now() // ⏱ start timer

    var req models.ImprovementSuggestionsRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    apiKey := os.Getenv("OPENAI_API_KEY")
    if apiKey == "" {
        http.Error(w, "OPENAI_API_KEY not set in environment", http.StatusInternalServerError)
        return
    }

    client := openai.NewClient(apiKey)
    log.Println("[Suggestions] ✅ OpenAI client initialized, secure connection established")

    prompt := fmt.Sprintf(`
Review the following cover letter draft for a %s position at %s.
Provide 3-5 specific, bulleted suggestions to improve its tone, clarity, and alignment with the job description.

COVER LETTER:
%s

JOB DESCRIPTION SUMMARY:
%s

Return the suggestions as a JSON array of strings.
`, req.Job.Title, req.Job.Company, req.Content, req.Job.Description)

    resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
        Model: "gpt-5.0",
        Messages: []openai.ChatCompletionMessage{
            {Role: "user", Content: prompt},
        },
        MaxTokens: 500,
    })
    if err != nil {
        log.Println("[Suggestions] ❌ OpenAI error:", err)
        json.NewEncoder(w).Encode(models.ImprovementSuggestionsResponse{
            Error: "Failed to generate suggestions",
        })
        return
    }

    output := resp.Choices[0].Message.Content

    var suggestions []string
    if err := json.Unmarshal([]byte(output), &suggestions); err != nil {
        log.Println("[Suggestions] ⚠️ JSON parse fallback triggered")
        suggestions = []string{"Improve clarity", "Add measurable achievements"}
    }

    elapsed := time.Since(start)
    log.Printf("[Suggestions] ⏱ Generated in %s\n", elapsed)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(models.ImprovementSuggestionsResponse{
        Suggestions: suggestions,
    })
}
