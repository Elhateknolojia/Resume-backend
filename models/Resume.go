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


type BlockStyle struct {
    X        int                    `json:"x"`
    Y        int                    `json:"y"`
    Width    int                    `json:"width"`
    Height   int                    `json:"height,omitempty"`
    Rotation *int                   `json:"rotation,omitempty"`
    IsLocked *bool                  `json:"isLocked,omitempty"`
    IsVisible *bool                 `json:"isVisible,omitempty"`
    Border   *string                `json:"border,omitempty"`
    Padding  *int                   `json:"padding,omitempty"`
    Style    map[string]interface{} `json:"style,omitempty"` // ✅ flexible
}


type ResumeData struct {
    Name             string          `json:"name"`
    Email            string          `json:"email"`
    Phone            string          `json:"phone"`
    PhoneCountryCode string          `json:"phoneCountryCode"`
    Location         string          `json:"location"`
    Summary          string          `json:"summary"`
    Sections         []ResumeSection `json:"sections"`
    Experience       []Experience    `json:"experience"`
    Education        []Education     `json:"education"`
    Referees         []Referee       `json:"referees"`
    Skills           []Skill         `json:"skills"`
    Hobbies          []Hobby         `json:"hobbies"`
    
    Website  *string `json:"website,omitempty"`
    QrCode   *string `json:"qrCode,omitempty"`
    SkillUrl *string `json:"skillUrl,omitempty"`

    
    Aesthetics       Aesthetics      `json:"aesthetics"`
    PageCount        int             `json:"pageCount"`

    MetadataStyle   *BlockStyle `json:"metadataStyle,omitempty"`
    ExperienceStyle *BlockStyle `json:"experienceStyle,omitempty"`
    EducationStyle  *BlockStyle `json:"educationStyle,omitempty"`
    SkillsStyle     *BlockStyle `json:"skillsStyle,omitempty"`
    RefereeStyle    *BlockStyle `json:"refereeStyle,omitempty"`
    QrStyle         *BlockStyle `json:"qrStyle,omitempty"`
    NameStyle       *BlockStyle `json:"nameStyle,omitempty"`
    EmailStyle      *BlockStyle `json:"emailStyle,omitempty"`
    PhoneStyle      *BlockStyle `json:"phoneStyle,omitempty"`
    SummaryStyle    *BlockStyle `json:"summaryStyle,omitempty"`

    
}


type Education struct {
    ID          string `json:"id"`
    School      string `json:"school"`
    Degree      string `json:"degree"`
    StartDate   string `json:"startDate"`
    EndDate     string `json:"endDate"`
    Description string `json:"description"`
}


type Hobby struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}