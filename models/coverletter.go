package models

type CoverLetterData struct {
    JobDescription   string `json:"jobDescription"`
    InstitutionName  string `json:"institutionName"`
    PositionTitle    string `json:"positionTitle"`
    Requirements     string `json:"requirements"`
    GeneratedLetter  string `json:"generatedLetter"`
}
