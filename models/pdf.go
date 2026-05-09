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
