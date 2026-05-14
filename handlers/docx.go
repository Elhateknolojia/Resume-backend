package handlers

import (
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "io"
    "github.com/unidoc/unioffice/document"
    "github.com/unidoc/unioffice/schema/soo/wml"
)

type DocxBlock struct {
    ID      string `json:"id"`
    Type    string `json:"type"`    // "paragraph" or "heading"
    Content string `json:"content"` // inner text or HTML
    Bold    bool   `json:"bold"`
    Italic  bool   `json:"italic"`
    Underline bool `json:"underline"`
}


type ImportDocxResponse struct {
    Blocks []DocxBlock `json:"blocks"`
}

func ImportDocxHandler(w http.ResponseWriter, r *http.Request) {
    file, _, err := r.FormFile("file")
    if err != nil {
        http.Error(w, "File upload error", http.StatusBadRequest)
        return
    }
    defer file.Close()

    tmpPath := "static/tmp.docx"
    out, _ := os.Create(tmpPath)
    defer out.Close()
    _, _ = io.Copy(out, file)

    doc, err := document.Open(tmpPath)
    if err != nil {
        http.Error(w, "Failed to open DOCX", http.StatusInternalServerError)
        return
    }

    blocks := []DocxBlock{}
    idCounter := 1

for _, para := range doc.Paragraphs() {
    text := ""
    bold := false
    italic := false
    underline := false

    for _, run := range para.Runs() {
        text += run.Text()

    props := run.Properties()
    if props.IsBold() {
        bold = true
    }
    if props.IsItalic() {
        italic = true
    }
    if props.Underline() != wml.ST_UnderlineNone {
        underline = true
    }
    }

    blockType := "paragraph"
    if para.Style() == "Heading1" || para.Style() == "Heading2" {
        blockType = "heading"
    }

    blocks = append(blocks, DocxBlock{
        ID:        fmt.Sprintf("%d", idCounter),
        Type:      blockType,
        Content:   text,
        Bold:      bold,
        Italic:    italic,
        Underline: underline,
    })
    idCounter++
}


    json.NewEncoder(w).Encode(ImportDocxResponse{Blocks: blocks})
}


func SaveDocxHandler(w http.ResponseWriter, r *http.Request) {
    var req struct {
        FileName string     `json:"fileName"`
        Blocks   []DocxBlock `json:"blocks"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    doc := document.New()
    for _, b := range req.Blocks {
        para := doc.AddParagraph()
        if b.Type == "heading" {
            para.SetStyle("Heading1")
        }
        run := para.AddRun()
        run.AddText(b.Content)
    }

    outPath := fmt.Sprintf("static/%s", req.FileName)
    if err := doc.SaveToFile(outPath); err != nil {
        http.Error(w, "Failed to save DOCX", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{
        "url": "/static/" + req.FileName,
    })
}
