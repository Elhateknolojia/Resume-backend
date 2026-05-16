package handlers

import (
    "encoding/json"
    "net/http"
    "Backend/models"
    "Backend/db"
    "github.com/jung-kurt/gofpdf" // simple PDF generator
    "github.com/gorilla/mux"
    "fmt"
    // "strconv"
    // "io"
    "os"
    "os/exec"
    "log"
    "strings"
    "path/filepath"


)

type GeneratePdfResponse struct {
    PdfUrl string `json:"pdfUrl"`
}



type GeneratePdfRequest struct {
    ResumeData models.ResumeData `json:"resumeData"`
    Model      string            `json:"model"`
}

func GeneratePdfHandler(w http.ResponseWriter, r *http.Request) {
    var resumeData models.ResumeData
    if err := json.NewDecoder(r.Body).Decode(&resumeData); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    email := resumeData.Email
    if email == "" {
        http.Error(w, "Missing email", http.StatusBadRequest)
        return
    }

    user, err := db.GetUserByEmail(email)
    if err != nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    // Free tier enforcement
    if user.Tier == "free" {
        if user.FreeDownloadsUsed >= 1 {
            http.Error(w, "Free download limit reached", http.StatusForbidden)
            return
        }
        user.FreeDownloadsUsed++
        if err := db.UpdateUser(email, user); err != nil {
            http.Error(w, "Failed to update user", http.StatusInternalServerError)
            return
        }
    }

    // --- Build safe filename ---
    safeEmail := strings.ReplaceAll(email, "@", "_")
    safeEmail = strings.ReplaceAll(safeEmail, ".", "_")
    cwd, _ := os.Getwd()
    pdfPath := filepath.Join(cwd, "static", "resume_"+safeEmail+".pdf")
    htmlContent := `
            <!DOCTYPE html>
            <html>
            <head>
            <meta charset="utf-8">
            <title>Resume</title>
            <link href="https://cdn.jsdelivr.net/npm/tailwindcss@2.2.19/dist/tailwind.min.css" rel="stylesheet">
            <style>
                body { background: #f4f4f5; }
                .a4-paper { width: 794px; height: 1123px; margin: auto; background: white; box-shadow: 0 32px 64px -12px rgba(0,0,0,0.14); }
            </style>
            </head>
            <body>
            <div class="a4-paper p-8 font-sans">
                <h1 class="text-4xl font-black uppercase text-center mb-4">` + resumeData.Name + `</h1>
                <p class="text-center text-sm text-gray-500 mb-2">` + resumeData.Email + `</p>
                <p class="text-center text-sm text-gray-500 mb-6">` + resumeData.Phone + `</p>

                <h2 class="text-lg font-bold uppercase text-gray-400 border-b mb-2">Summary</h2>
                <p class="text-sm text-gray-600 italic mb-6">` + resumeData.Summary + `</p>

                <h2 class="text-lg font-bold uppercase text-gray-400 border-b mb-2">Experience</h2>
                <ul class="mb-6">`
            for _, exp := range resumeData.Experience {
                htmlContent += `<li class="mb-2">
                <span class="font-bold">` + exp.Title + `</span> at ` + exp.Company + ` (` + exp.StartDate + ` - ` + exp.EndDate + `)
                <p class="text-sm text-gray-600">` + exp.Content + `</p>
                </li>`
            }
            htmlContent += `</ul>

                <h2 class="text-lg font-bold uppercase text-gray-400 border-b mb-2">Education</h2>
                <ul class="mb-6">`
            for _, edu := range resumeData.Education {
                htmlContent += `<li class="mb-2">
                <span class="font-bold">` + edu.Degree + `</span> at ` + edu.School + ` (` + edu.StartDate + ` - ` + edu.EndDate + `)
                </li>`
            }
            htmlContent += `</ul>

                <h2 class="text-lg font-bold uppercase text-gray-400 border-b mb-2">Skills</h2>
                <ul class="flex flex-wrap gap-2">`
            for _, skill := range resumeData.Skills {
                htmlContent += `<li class="px-3 py-1 bg-gray-100 rounded text-xs font-bold uppercase">` + skill.Name + `</li>`
            }
            htmlContent += `</ul>

                <h2 class="text-lg font-bold uppercase text-gray-400 border-b mb-2">Referees</h2>
                <ul class="mb-6">`
            for _, ref := range resumeData.Referees {
                htmlContent += `<li class="mb-2">
                <span class="font-bold">` + ref.Name + `</span> — ` + ref.Email + `
                </li>`
            }
            htmlContent += `</ul>
            </div>
            </body>
            </html>`


    htmlPath := filepath.Join(cwd, "templates", "resume_"+safeEmail+".html")
    os.WriteFile(htmlPath, []byte(htmlContent), 0644)


    cmd := exec.Command("xvfb-run", "wkhtmltopdf", htmlPath, pdfPath)
    cmd.Stderr = os.Stderr // log errors

    if err := cmd.Run(); err != nil {
        log.Printf("wkhtmltopdf failed: %v", err)
        http.Error(w, "Failed to generate PDF", http.StatusInternalServerError)
        return
    }

    log.Printf("PDF generated at: %s", pdfPath)

    resp := GeneratePdfResponse{
        PdfUrl: "/static/resume_" + safeEmail + ".pdf",
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}


type SavePendingResponse struct {
    Message string `json:"message"`
    PdfUrl  string `json:"pdfUrl,omitempty"`
}

func init(){
    if _, err := os.Stat("static"); os.IsNotExist(err) {
        os.Mkdir("static", 0755)
    }
}

func SavePendingResumeHandler(w http.ResponseWriter, r *http.Request) {
    var resumeData models.ResumeData
    if err := json.NewDecoder(r.Body).Decode(&resumeData); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    email := resumeData.Email
    if email == "" {
        http.Error(w, "Missing email", http.StatusBadRequest)
        return
    }

    user, err := db.GetUserByEmail(email)
    if err != nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    // --- Generate PDF ---
    pdf := gofpdf.New("P", "mm", "A4", "")
    pdf.AddPage()
    pdf.SetFont("Arial", "B", 16)
    pdf.Cell(40, 10, "Resume")

    pdf.Ln(12)
    pdf.SetFont("Arial", "", 12)
    pdf.Cell(0, 10, "Name: "+resumeData.Name)
    pdf.Ln(8)
    pdf.Cell(0, 10, "Email: "+resumeData.Email)
    pdf.Ln(8)
    pdf.MultiCell(0, 10, "Summary: "+resumeData.Summary, "", "", false)

    
    safeEmail := strings.ReplaceAll(email, "@", "_")
    safeEmail = strings.ReplaceAll(safeEmail, ".", "_")
    filePath := "static/resume_" + safeEmail + ".pdf"
    
    if err := pdf.OutputFileAndClose(filePath); err != nil {
        log.Printf("PDF generation failed: %v", err)
        http.Error(w, "Failed to generate PDF", http.StatusInternalServerError)
        return
    }


    // --- Eligibility check ---
    if user.Tier == "premium" {
        // Premium users can download immediately
        resp := SavePendingResponse{
            Message: "Premium user — resume ready",
            PdfUrl:  "/static/resume_" + safeEmail + ".pdf",
        }
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(resp)
        return
    }

    // Free users → mark as pending
    user.PendingResume = filePath
    if err := db.UpdateUser(email, user); err != nil {
        http.Error(w, "Failed to update user", http.StatusInternalServerError)
        return
    }

    resp := SavePendingResponse{
        Message: "Resume saved as pending until subscription is active",
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}


type ExportHtmlRequest struct {
    Html string `json:"html"`
}

func ExportHtmlHandler(w http.ResponseWriter, r *http.Request) {
    var req ExportHtmlRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    safeEmail := "resume_export" // or derive from user session/email
    cwd, _ := os.Getwd()
    htmlPath := filepath.Join(cwd, "templates", safeEmail+".html")
    pdfPath := filepath.Join(cwd, "static", safeEmail+".pdf")

    // Write HTML file
    if err := os.WriteFile(htmlPath, []byte(req.Html), 0644); err != nil {
        http.Error(w, "Failed to write HTML", http.StatusInternalServerError)
        return
    }

    // Convert to PDF
    cmd := exec.Command("xvfb-run", "wkhtmltopdf", "--enable-local-file-access", htmlPath, pdfPath)
    cmd.Stderr = os.Stderr
    if err := cmd.Run(); err != nil {
        log.Printf("wkhtmltopdf failed: %v", err)
        http.Error(w, "Failed to generate PDF", http.StatusInternalServerError)
        return
    }

    resp := GeneratePdfResponse{PdfUrl: "/static/" + safeEmail + ".pdf"}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}





// Mock template store (later replace with DB or files)
// var templates = map[string]map[string]interface{}{
//     "modern": templateModern, // defined in templates/modern.json
//     "minimal": templateMinimal, // defined in templates/minimal.json
//     "ceo": templateCEO, // defined in templates/ceo.json
// }

func LoadTemplateHandler(w http.ResponseWriter, r *http.Request) {
    id := mux.Vars(r)["id"]
    path := fmt.Sprintf("templates/%s.json", id)

    data, err := os.ReadFile(path)

    if err != nil {
        http.Error(w, "Template not found", http.StatusNotFound)
        return
    }


    w.Header().Set("Content-Type", "application/json")
    w.Write(data)
}

