package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/unidoc/unipdf/v3/creator"
	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"
	// "os"
	// "fmt"
	// "github.com/jung-kurt/gofpdf"
	// "github.com/pdfcpu/pdfcpu/pkg/api"
	// "github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

type Section struct {
    ID       string `json:"id"`
    Content  string `json:"content"`
    FontName string `json:"fontName"`
    FontSize float64 `json:"fontSize"`
    Bold     bool   `json:"bold"`
    Italic   bool   `json:"italic"`
}

type ImportResponse struct {
    Sections []Section `json:"sections"`
}

// Import PDF and return editable sections// Import PDF and send to GPT-5.0 for styled JSON
func ImportPdfHandler(w http.ResponseWriter, r *http.Request) {
    file, _, err := r.FormFile("file")
    if err != nil {
        http.Error(w, "File upload error", http.StatusBadRequest)
        return
    }
    defer file.Close()

    pdfBytes, _ := io.ReadAll(file)

    // Build GPT request
	gptReq := map[string]interface{}{
		"model": "gpt-4.1", // or gpt-3.5-turbo depending on your access
		"messages": []map[string]string{
			{
				"role": "system",
				"content": "You are a PDF parser. Extract text sections with font info.",
			},
			{
				"role": "user",
				"content": fmt.Sprintf("Here is the PDF content (base64): %s", base64.StdEncoding.EncodeToString(pdfBytes)),
			},
		},
	}



    body, _ := json.Marshal(gptReq)

    // Call GPT-5.0 API
    apiKey := os.Getenv("OPENAI_API_KEY")
    reqGPT, _ := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(body))
    reqGPT.Header.Set("Content-Type", "application/json")
    reqGPT.Header.Set("Authorization", "Bearer "+apiKey)

    client := &http.Client{}
    respGPT, err := client.Do(reqGPT)
    if err != nil {
        http.Error(w, "AI service error", http.StatusInternalServerError)
        return
    }
    defer respGPT.Body.Close()

	rawBody, _ := io.ReadAll(respGPT.Body)
	log.Printf("Raw GPT response: %s", string(rawBody))

	// Reset the body reader so you can decode again
	respGPT.Body = io.NopCloser(bytes.NewBuffer(rawBody))


		var gptResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(respGPT.Body).Decode(&gptResp); err != nil {
		http.Error(w, "Failed to decode AI response", http.StatusInternalServerError)
		return
	}

	if len(gptResp.Choices) == 0 || gptResp.Choices[0].Message.Content == "" {
		log.Printf("GPT returned no usable choices, falling back to UniPDF")

		// Use UniPDF extraction directly on the uploaded file
		tmpFile := "static/tmp.pdf"
		os.WriteFile(tmpFile, pdfBytes, 0644)

		sections, err := ExtractStyledSections(tmpFile)
		if err != nil {
			http.Error(w, "Failed to extract PDF", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(ImportResponse{Sections: sections})
		return
	}


	var sections []Section
	if err := json.Unmarshal([]byte(gptResp.Choices[0].Message.Content), &sections); err != nil {
		log.Printf("Failed to parse GPT response: %v", err)
		http.Error(w, "Failed to parse AI response", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(ImportResponse{Sections: sections})


}

// Save edited sections back into PDF with styles
func SavePdfHandler(w http.ResponseWriter, r *http.Request) {
    var req struct {
        FileName string    `json:"fileName"`
        Sections []Section `json:"sections"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    c := creator.New()
    c.SetPageMargins(50, 50, 50, 50)

    // Load a base font (Helvetica is one of the built-in 14 fonts)
    baseFont, err := model.NewStandard14Font("Helvetica")
    if err != nil {
        http.Error(w, "Failed to load font", http.StatusInternalServerError)
        return
    }

    for _, s := range req.Sections {
        c.NewPage()

        p := c.NewParagraph(s.Content)

        // Apply styles
        font := baseFont
        if s.Bold && s.Italic {
            font, _ = model.NewStandard14Font("Helvetica-BoldOblique")
        } else if s.Bold {
            font, _ = model.NewStandard14Font("Helvetica-Bold")
        } else if s.Italic {
            font, _ = model.NewStandard14Font("Helvetica-Oblique")
        }

        p.SetFont(font)
        if s.FontSize > 0 {
            p.SetFontSize(s.FontSize)
        } else {
            p.SetFontSize(12) // default
        }

        c.Draw(p)
    }

    outPath := fmt.Sprintf("static/%s", req.FileName)
    if err := c.WriteToFile(outPath); err != nil {
        http.Error(w, "Failed to save PDF", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{
        "url": "/static/" + req.FileName,
    })
}

func ExtractStyledSections(filePath string) ([]Section, error) {
    f, err := os.Open(filePath)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    pdfReader, err := model.NewPdfReader(f)
    if err != nil {
        return nil, err
    }

    numPages, _ := pdfReader.GetNumPages()
    sections := []Section{}

    for i := 1; i <= numPages; i++ {
        page, _ := pdfReader.GetPage(i)
        ex, _ := extractor.New(page)

        pageText, _, _, err := ex.ExtractPageText()
        if err != nil {
            return nil, err
        }
		log.Printf("Error extracting page text: %v", err)


		marks := pageText.Marks()
		if marks != nil {
			elems := marks.Elements()
			for j, m := range elems {
				sections = append(sections, Section{
					ID:      fmt.Sprintf("%d-%d", i, j),
					Content: m.Text,
					// Font info not available in v3.69.0
					FontName: "",
					FontSize: 0,
					Bold:     false,
					Italic:   false,
				})
			}
		}

    }

    return sections, nil
}

func RebuildPdf(fileName string, sections []Section) error {
    c := creator.New()
    c.SetPageMargins(50, 50, 50, 50)

for _, s := range sections {
    c.NewPage()
    p := c.NewParagraph(s.Content)

    // Just use Helvetica, no style detection in v3.69.0
    font, err := model.NewStandard14Font("Helvetica")
    if err != nil {
        return err
    }
    p.SetFont(font)

    if s.FontSize > 0 {
        p.SetFontSize(s.FontSize)
    } else {
        p.SetFontSize(12)
    }
    p.SetMargins(10, 10, 10, 10)

    c.Draw(p)
}

    return c.WriteToFile("static/" + fileName)
}
