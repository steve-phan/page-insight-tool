package models

type AnalysisError struct {
	Message string    `json:"message" example:"Invalid URL"`
	Type    ErrorType `json:"type" example:"INVALID_URL"`
	Code    int       `json:"code" example:"400"`
}

type ErrorType string

const (
	ErrorTypeInvalidURL          ErrorType = "INVALID_URL"
	ErrorTypeInternalServerError ErrorType = "INTERNAL_SERVER_ERROR"
)

func (e *AnalysisError) Error() string {
	return e.Message
}
