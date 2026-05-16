package handlers

import (
    // "encoding/json"
    // "net/http"
    // "os"
    // "log"
    "Backend/models"

    // "github.com/openai/openai-go" // ✅ official OpenAI Go SDK
    // "golang.org/x/net/context"

    "encoding/json"
    "net/http"
    "Backend/db"
    // "Backend/models"
    "context"
    "go.mongodb.org/mongo-driver/bson"
)


// helper to keep prompt logic clean
func buildCoverLetterPrompt(req models.CoverLetterRequest) string {
    return `
You are an expert career coach and professional writer.
Generate a professional, compelling, and ATS-friendly cover letter.

JOB DETAILS:
Title: ` + req.Job.Title + `
Company: ` + req.Job.Company + `
Job Description: ` + req.Job.Description + `

USER RESUME INFO:
Name: ` + req.Resume.FullName + `
Email: ` + req.Resume.Email + `
Phone: ` + req.Resume.Phone + `
Summary: ` + req.Resume.Summary + `
Experience: ` + req.Resume.Experience + `
Skills: ` + req.Resume.Skills + `

EMPHASIS: ` + string(req.Emphasis) + `

INSTRUCTIONS:
- Start with contact information for both the user and the company (placeholder for company address).
- Write a strong opening that captures attention and mentions the ` + req.Job.Title + ` role at ` + req.Job.Company + `.
- Align the user's skills and experience with the job description requirements.
- Use a professional yet persuasive tone.
- Keep it to approximately 300-400 words.
- Ensure it is structured in clear paragraphs.
- DO NOT add leading spaces or indentation to paragraphs.
- Return ONLY the content of the cover letter, ready to be printed.
`
}

// handlers/resume.go
// package handlers


func GetResumeHandler(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value("userID").(string)

    var user models.User
    err := db.UserCollection.FindOne(context.Background(), bson.M{"_id": userID}).Decode(&user)
    if err != nil {
        http.Error(w, "Resume not found", http.StatusNotFound)
        return
    }

    // Map ResumeData → UserResume
    resume := models.UserResume{
        FullName:   user.ResumeData.Name,
        Email:      user.ResumeData.Email,
        Phone:      user.ResumeData.Phone,
        Summary:    user.ResumeData.Summary,
        Experience: flattenExperience(user.ResumeData.Experience),
        Skills:     flattenSkills(user.ResumeData.Skills),
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resume)
}

func flattenExperience(exps []models.Experience) string {
    out := ""
    for _, e := range exps {
        out += e.Title + " at " + e.Company + " (" + e.StartDate + " - " + e.EndDate + ")\n" + e.Content + "\n"
    }
    return out
}

func flattenSkills(skills []models.Skill) string {
    out := ""
    for _, s := range skills {
        out += s.Name + ", "
    }
    return out
}
