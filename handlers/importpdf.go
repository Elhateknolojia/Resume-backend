package handlers
// import (
//     "encoding/json"
//     "net/http"
// 	"Backend/models"
//     "github.com/unidoc/unipdf/v3/extractor"
//     "github.com/unidoc/unipdf/v3/model"
//     "github.com/unidoc/unipdf/v3/creator"
//     "strings"
// )


// func ImportPdfHandler(w http.ResponseWriter, r *http.Request) {
//     file, _, err := r.FormFile("pdf")
//     if err != nil {
//         http.Error(w, "Upload error", http.StatusBadRequest)
//         return
//     }
//     defer file.Close()

//     pdfReader, _ := model.NewPdfReader(file)
//     numPages, _ := pdfReader.GetNumPages()

//     var pagesPayload []models.RawPageData

//     for i := 1; i <= numPages; i++ {
//         page, _ := pdfReader.GetPage(i)
//         ex, _ := extractor.New(page)
//         pageText, _, _, err := ex.ExtractPageText() // Gets coordinates and styles
//         if err != nil {
//             http.Error(w, "Failed to extract text", http.StatusInternalServerError)
//             return
//         }
//         pageData := models.RawPageData{PageNumber: i}

        
//             // Clean: Remove empty strings or excessive whitespace
//             cleanText := strings.TrimSpace(mark.Text)
//             if cleanText == "" {
//                 continue
//             }

//             pageData.TextBlocks = append(pageData.TextBlocks, models.TextBlock{
//                 Text:     cleanText,
//                 X:        mark.BBox.L,
//                 Y:        mark.BBox.T,
//                 FontSize: mark.FontSize,
//             })
        
//         pagesPayload = append(pagesPayload, pageData)
//     }

//     // HANDOFF POINT: 
//     // Usually, you would send 'pagesPayload' to your Python service via HTTP POST.
//     // For now, we return it as JSON to be consumed by your GPT-2 script.
//     w.Header().Set("Content-Type", "application/json")
//     json.NewEncoder(w).Encode(pagesPayload)
// }


// func ReconstructPdfHandler(w http.ResponseWriter, r *http.Request) {
//     var data models.ImportedPDF
//     if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
//         http.Error(w, "Invalid payload", http.StatusBadRequest)
//         return
//     }

//     c := creator.New()

//     if data.NameStyle != nil {
//         p := c.NewParagraph(data.Name)
//         p.SetFontSize(data.NameStyle.FontSize)
//         p.SetPos(data.NameStyle.X, data.NameStyle.Y)
//         c.Draw(p)
//     }

//     w.Header().Set("Content-Type", "application/pdf")
//     if err := c.Write(w); err != nil {
//         http.Error(w, "Failed to write PDF", http.StatusInternalServerError)
//     }
// }
