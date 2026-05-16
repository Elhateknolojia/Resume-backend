package models

// JobDetails represents the job posting info
type JobDetails struct {
    Title       string `json:"title"`
    Company     string `json:"company"`
    Description string `json:"description"`
}

// UserResume represents the key resume fields needed for cover letter generation
type UserResume struct {
    FullName   string `json:"fullName"`
    Email      string `json:"email"`
    Phone      string `json:"phone"`
    Summary    string `json:"summary"`
    Experience string `json:"experience"`
    Skills     string `json:"skills"`
}

// GenerationEmphasis defines the writing focus
// e.g. "skills-focused", "experience-focused", "achievements-focused", "balanced"
type GenerationEmphasis string

const (
    SkillsFocused      GenerationEmphasis = "skills-focused"
    ExperienceFocused  GenerationEmphasis = "experience-focused"
    AchievementsFocused GenerationEmphasis = "achievements-focused"
    Balanced           GenerationEmphasis = "balanced"
)

// CoverLetterRequest is the payload sent from Angular to Go backend
type CoverLetterRequest struct {
    Job      JobDetails        `json:"job"`
    Resume   UserResume        `json:"resume"`
    Emphasis GenerationEmphasis `json:"emphasis"`
}

// CoverLetterResponse is what the backend returns to Angular
type CoverLetterResponse struct {
    CoverLetter string `json:"coverLetter"`
    Error       string `json:"error,omitempty"`
}

// ImprovementSuggestionsRequest is for polishing an existing draft
type ImprovementSuggestionsRequest struct {
    Content string    `json:"content"`
    Job     JobDetails `json:"job"`
}

// ImprovementSuggestionsResponse returns AI feedback
type ImprovementSuggestionsResponse struct {
    Suggestions []string `json:"suggestions"`
    Error       string   `json:"error,omitempty"`
}
