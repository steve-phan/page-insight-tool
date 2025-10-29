package models

import "errors"

type AnalysisError struct {
	Message string    `json:"message" example:"Invalid URL"`
	Type    ErrorType `json:"type" example:"INVALID_URL"`
	Code    int       `json:"code" example:"400"`
}

type ErrorType string

const (
	ErrorTypeInvalidURL          ErrorType = "INVALID_URL"
	ErrorTypeForbidden           ErrorType = "FORBIDDEN"
	ErrorTypePageNotFound        ErrorType = "PAGE_NOT_FOUND"
	ErrorTypeTimeout             ErrorType = "TIMEOUT"
	ErrorTypeHTMLParse           ErrorType = "HTML_PARSE_ERROR"
	ErrorTypeNonOKStatus         ErrorType = "NON_OK_STATUS"
	ErrorTypeBodyTooLarge        ErrorType = "BODY_TOO_LARGE"
	ErrorTypeInternalServerError ErrorType = "INTERNAL_SERVER_ERROR"
)

// Predefined errors for common scenarios
var (
	ErrHTMLParse           = errors.New("failed to parse HTML")
	ErrNonOKStatus         = errors.New("non-200 status code")
	ErrBodyTooLarge        = errors.New("response body too large")
	ErrForbidden           = errors.New("access forbidden")
	ErrPageNotFound        = errors.New("page not found")
	ErrTimeout             = errors.New("request timeout")
	ErrInternalServerError = errors.New("internal server error")
	ErrInvalidURL          = errors.New("invalid URL")
)

func (e *AnalysisError) Error() string {
	return e.Message
}

// Helper functions to create structured errors
func NewInvalidURLError(message string) *AnalysisError {
	return &AnalysisError{
		Message: message,
		Type:    ErrorTypeInvalidURL,
		Code:    400,
	}
}

func NewForbiddenError(message string) *AnalysisError {
	return &AnalysisError{
		Message: message,
		Type:    ErrorTypeForbidden,
		Code:    403,
	}
}

func NewPageNotFoundError(message string) *AnalysisError {
	return &AnalysisError{
		Message: message,
		Type:    ErrorTypePageNotFound,
		Code:    404,
	}
}

func NewTimeoutError(message string) *AnalysisError {
	return &AnalysisError{
		Message: message,
		Type:    ErrorTypeTimeout,
		Code:    408,
	}
}

func NewHTMLParseError(message string) *AnalysisError {
	return &AnalysisError{
		Message: message,
		Type:    ErrorTypeHTMLParse,
		Code:    422,
	}
}

func NewInternalServerError(message string) *AnalysisError {
	return &AnalysisError{
		Message: message,
		Type:    ErrorTypeInternalServerError,
		Code:    500,
	}
}
