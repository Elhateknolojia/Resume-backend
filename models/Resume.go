package models

type Referee struct {
    ID      string `json:"id"`
    Name    string `json:"name"`
    Email   string `json:"email"`
    Phone   string `json:"phone"`
    Address string `json:"address"`
}

type Experience struct {
    ID        string `json:"id"`
    Title     string `json:"title"`
    Company   string `json:"company"`
    StartDate string `json:"startDate"`
    EndDate   string `json:"endDate"`
    Current   bool   `json:"current"`
    Content   string `json:"content"`
}

type ResumeSection struct {
    ID      string `json:"id"`
    Title   string `json:"title"`
    Content string `json:"content"`
}

type ResumeElement struct {
    ID        string  `json:"id"`
    Type      string  `json:"type"` // image, line, box, text
    Content   *string `json:"content,omitempty"`
    URL       *string `json:"url,omitempty"`
    X         int     `json:"x"`
    Y         int     `json:"y"`
    Width     int     `json:"width"`
    Height    int     `json:"height"`
    Rotation  *int    `json:"rotation,omitempty"`
    IsLocked  *bool   `json:"isLocked,omitempty"`
    IsVisible *bool   `json:"isVisible,omitempty"`
}

type Aesthetics struct {
    FontFamily     string          `json:"fontFamily"`
    PrimaryColor   string          `json:"primaryColor"`
    BackgroundColor string         `json:"backgroundColor"`
    FontSize       int             `json:"fontSize"`
    Elements       []ResumeElement `json:"elements"`
}

type Skill struct {
    Name        string `json:"name"`
    Level       int    `json:"level"`
    DisplayMode string `json:"displayMode,omitempty"`
}

type ResumeData struct {
    Name            string          `json:"name"`
    Email           string          `json:"email"`
    Phone           string          `json:"phone"`
    PhoneCountryCode string         `json:"phoneCountryCode"`
    Location        string          `json:"location"`
    Summary         string          `json:"summary"`
    Sections        []ResumeSection `json:"sections"`
    Experience      []Experience    `json:"experience"`
    Education       []string        `json:"education"`
    Referees        []Referee       `json:"referees"`
    Skills          []Skill         `json:"skills"`
    Hobbies         []string        `json:"hobbies"`
    Website         *string         `json:"website,omitempty"`
    Aesthetics      Aesthetics      `json:"aesthetics"`
    PageCount       int             `json:"pageCount"`
}
