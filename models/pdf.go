package models

type RawPageData struct {
    PageNumber int           `json:"page_number"`
    TextBlocks []TextBlock   `json:"text_blocks"`
}

type TextBlock struct {
    Text     string  `json:"t"` // Raw text content
    X        float64 `json:"x"` // Horizontal position
    Y        float64 `json:"y"` // Vertical position
    FontSize float64 `json:"s"` // Font size (helps LLM identify headers)
}


// models/importedpdf.go

type ImportedPDF struct {
    Name       string     `json:"name"`
    NameStyle  *BlockStyle `json:"nameStyle"`
    // add Email, Phone, Summary, Experience, etc.
}

type importedBlockStyle struct {
    X        float64 `json:"x"`
    Y        float64 `json:"y"`
    FontSize float64 `json:"fontSize"`
}
