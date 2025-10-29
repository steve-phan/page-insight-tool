package models

import "time"

// HTTPError represents an HTTP error response structure
// This is purely for HTTP API responses, not business logic
type HTTPError struct {
	Message   string                 `json:"message" example:"Invalid URL"`
	Type      string                 `json:"type" example:"INVALID_URL"`
	Code      int                    `json:"code" example:"400"`
	Details   map[string]interface{} `json:"details,omitempty"`
	RequestID string                 `json:"request_id,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}

// Error implements the error interface for HTTPError
func (e *HTTPError) Error() string {
	return e.Message
}

// HTTPErrorResponse is a standardized error response wrapper
type HTTPErrorResponse struct {
	Error     *HTTPError `json:"error"`
	Success   bool       `json:"success"`
	Timestamp time.Time  `json:"timestamp"`
}

// NewHTTPErrorResponse creates a standardized error response
func NewHTTPErrorResponse(httpError *HTTPError) *HTTPErrorResponse {
	if httpError.Timestamp.IsZero() {
		httpError.Timestamp = time.Now()
	}

	return &HTTPErrorResponse{
		Error:     httpError,
		Success:   false,
		Timestamp: time.Now(),
	}
}
