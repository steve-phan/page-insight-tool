package models

type AnalysisRequest struct {
	URL string `json:"url" validate:"required,url" example:"https://www.google.com"`
}

type AnalysisResponse struct {
	HTMLVersion  string   `json:"html_version" example:"HTML5"`
	PageTitle    string   `json:"page_title" example:"Google"`
	Headings     Headings `json:"headings"`
	Links        Links    `json:"links"`
	HasLoginForm bool     `json:"has_login_form" example:"true"`
	AnalysisTime int64    `json:"analysis_time_ms" example:"150"`
}

type Headings struct {
	H1 int `json:"h1" example:"10"`
	H2 int `json:"h2" example:"20"`
	H3 int `json:"h3" example:"30"`
	H4 int `json:"h4" example:"40"`
	H5 int `json:"h5" example:"50"`
	H6 int `json:"h6" example:"60"`
}

type Links struct {
	Internal     int `json:"internal" example:"10"`
	External     int `json:"external" example:"20"`
	Inaccessible int `json:"inaccessible" example:"30"`
}
