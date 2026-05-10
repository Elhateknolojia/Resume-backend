package models

type MockTemplate struct {
    Name      string       `json:"name"`
    Email     string       `json:"email"`
    Phone     string       `json:"phone"`
    Location  string       `json:"location"`
    Summary   string       `json:"summary"`
    Experience []MockExperience `json:"mockExperience"`
    Education  []MockEducation  `json:"mockEducation"`
    Skills     []MockSkill      `json:"mockSkills"`
    Referees   []MockReferee    `json:"mockReferees"`
}

type MockExperience struct {
    Title     string `json:"title"`
    Company   string `json:"company"`
    StartDate string `json:"startDate"`
    EndDate   string `json:"endDate"`
    Content   string `json:"content"`
    Current   bool   `json:"current"`
}

type MockEducation struct {
    School    string `json:"school"`
    Degree    string `json:"degree"`
    StartDate string `json:"startDate"`
    EndDate   string `json:"endDate"`
}

type MockSkill struct {
    Name string `json:"name"`
}

type MockReferee struct {
    Name    string `json:"name"`
    Email   string `json:"email"`
    Phone   string `json:"phone"`
    Address string `json:"address"`
}
