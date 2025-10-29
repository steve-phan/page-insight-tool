package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// Domain error types following Go best practices
// These errors represent business logic failures, not HTTP concerns

// ErrorType represents the category of domain error
type ErrorType string

const (
	// Input validation errors
	ErrorTypeInvalidInput ErrorType = "INVALID_INPUT"
	ErrorTypeInvalidURL   ErrorType = "INVALID_URL"

	// Network/external service errors
	ErrorTypeNetworkTimeout    ErrorType = "NETWORK_TIMEOUT"
	ErrorTypeNetworkConnection ErrorType = "NETWORK_CONNECTION"
	ErrorTypeNetworkDNS        ErrorType = "NETWORK_DNS"
	ErrorTypeNetworkSSL        ErrorType = "NETWORK_SSL"

	// HTTP-specific errors
	ErrorTypeHTTPForbidden   ErrorType = "HTTP_FORBIDDEN"
	ErrorTypeHTTPNotFound    ErrorType = "HTTP_NOT_FOUND"
	ErrorTypeHTTPNonOKStatus ErrorType = "HTTP_NON_OK_STATUS"

	// Content processing errors
	ErrorTypeHTMLParse     ErrorType = "HTML_PARSE"
	ErrorTypeContentTooBig ErrorType = "CONTENT_TOO_BIG"

	// Internal/system errors
	ErrorTypeInternal ErrorType = "INTERNAL"
)

// DomainError represents a business domain error with context
type DomainError struct {
	Type       ErrorType
	Message    string
	Cause      error
	StatusCode int // HTTP status code hint for mapping
	Details    map[string]interface{}
}

func (e *DomainError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

func (e *DomainError) Unwrap() error {
	return e.Cause
}

// Sentinel errors for common cases (following Go conventions)
var (
	ErrMissingURL = errors.New("URL parameter is required")
)

// Domain error constructors following Go conventions

// Input Validation Errors
func NewInvalidURLError(url string, cause error) *DomainError {
	return &DomainError{
		Type:       ErrorTypeInvalidURL,
		Message:    fmt.Sprintf("invalid URL: %s", url),
		Cause:      cause,
		StatusCode: http.StatusBadRequest,
		Details:    map[string]interface{}{"url": url},
	}
}

func NewInvalidInputError(field string, value interface{}, reason string) *DomainError {
	return &DomainError{
		Type:       ErrorTypeInvalidInput,
		Message:    fmt.Sprintf("invalid %s: %s", field, reason),
		StatusCode: http.StatusBadRequest,
		Details: map[string]interface{}{
			"field":  field,
			"value":  value,
			"reason": reason,
		},
	}
}

// Network Errors
func NewNetworkTimeoutError(url string, cause error) *DomainError {
	return &DomainError{
		Type:       ErrorTypeNetworkTimeout,
		Message:    fmt.Sprintf("network timeout for URL: %s", url),
		Cause:      cause,
		StatusCode: http.StatusRequestTimeout,
		Details:    map[string]interface{}{"url": url},
	}
}

func NewNetworkConnectionError(url string, cause error) *DomainError {
	return &DomainError{
		Type:       ErrorTypeNetworkConnection,
		Message:    fmt.Sprintf("connection failed for URL: %s", url),
		Cause:      cause,
		StatusCode: http.StatusBadGateway,
		Details:    map[string]interface{}{"url": url},
	}
}

func NewNetworkDNSError(hostname string, cause error) *DomainError {
	return &DomainError{
		Type:       ErrorTypeNetworkDNS,
		Message:    fmt.Sprintf("DNS resolution failed for host: %s", hostname),
		Cause:      cause,
		StatusCode: http.StatusBadGateway,
		Details:    map[string]interface{}{"hostname": hostname},
	}
}

func NewNetworkSSLError(url string, cause error) *DomainError {
	return &DomainError{
		Type:       ErrorTypeNetworkSSL,
		Message:    fmt.Sprintf("SSL/TLS error for URL: %s", url),
		Cause:      cause,
		StatusCode: http.StatusBadGateway,
		Details:    map[string]interface{}{"url": url},
	}
}

// HTTP Errors
func NewHTTPForbiddenError(url string) *DomainError {
	return &DomainError{
		Type:       ErrorTypeHTTPForbidden,
		Message:    fmt.Sprintf("access forbidden to URL: %s", url),
		StatusCode: http.StatusForbidden,
		Details:    map[string]interface{}{"url": url},
	}
}

func NewHTTPNotFoundError(url string) *DomainError {
	return &DomainError{
		Type:       ErrorTypeHTTPNotFound,
		Message:    fmt.Sprintf("page not found: %s", url),
		StatusCode: http.StatusNotFound,
		Details:    map[string]interface{}{"url": url},
	}
}

func NewHTTPNonOKStatusError(url string, statusCode int) *DomainError {
	return &DomainError{
		Type:       ErrorTypeHTTPNonOKStatus,
		Message:    fmt.Sprintf("HTTP %d error for URL: %s", statusCode, url),
		StatusCode: http.StatusBadGateway,
		Details: map[string]interface{}{
			"url":             url,
			"response_status": statusCode,
		},
	}
}

// Content Processing Errors
func NewHTMLParseError(url string, cause error) *DomainError {
	return &DomainError{
		Type:       ErrorTypeHTMLParse,
		Message:    fmt.Sprintf("failed to parse HTML from URL: %s", url),
		Cause:      cause,
		StatusCode: http.StatusUnprocessableEntity,
		Details:    map[string]interface{}{"url": url},
	}
}

func NewContentTooBigError(url string, size int64, maxSize int64) *DomainError {
	return &DomainError{
		Type:       ErrorTypeContentTooBig,
		Message:    fmt.Sprintf("content too large: %d bytes (max: %d)", size, maxSize),
		StatusCode: http.StatusRequestEntityTooLarge,
		Details: map[string]interface{}{
			"url":      url,
			"size":     size,
			"max_size": maxSize,
		},
	}
}

// Internal Errors
func NewInternalError(message string, cause error) *DomainError {
	return &DomainError{
		Type:       ErrorTypeInternal,
		Message:    message,
		Cause:      cause,
		StatusCode: http.StatusInternalServerError,
	}
}
