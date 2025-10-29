package errors

import (
	"errors"
	"net"
	"net/http"
	"net/url"
	"strings"
	"syscall"
)

// Helper functions to check error types
func IsNetworkError(err error) bool {
	var domainErr *DomainError
	if errors.As(err, &domainErr) {
		switch domainErr.Type {
		case ErrorTypeNetworkTimeout, ErrorTypeNetworkConnection, ErrorTypeNetworkDNS, ErrorTypeNetworkSSL:
			return true
		}
	}
	return false
}

func IsInputValidationError(err error) bool {
	var domainErr *DomainError
	if errors.As(err, &domainErr) {
		switch domainErr.Type {
		case ErrorTypeInvalidInput, ErrorTypeInvalidURL:
			return true
		}
	}
	return false
}

func IsHTTPError(err error) bool {
	var domainErr *DomainError
	if errors.As(err, &domainErr) {
		switch domainErr.Type {
		case ErrorTypeHTTPForbidden, ErrorTypeHTTPNotFound, ErrorTypeHTTPNonOKStatus:
			return true
		}
	}
	return false
}

func IsContentProcessingError(err error) bool {
	var domainErr *DomainError
	if errors.As(err, &domainErr) {
		switch domainErr.Type {
		case ErrorTypeHTMLParse, ErrorTypeContentTooBig:
			return true
		}
	}
	return false
}

func IsInternalError(err error) bool {
	var domainErr *DomainError
	if errors.As(err, &domainErr) {
		return domainErr.Type == ErrorTypeInternal
	}
	return false
}

// ClassifyNetworkError converts network errors to domain-specific errors
// This function provides centralized error classification for network-related failures
func ClassifyNetworkError(targetURL string, err error) error {
	// Check for timeout errors
	if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
		return NewNetworkTimeoutError(targetURL, err)
	}

	// Check for connection refused
	if opErr, ok := err.(*net.OpError); ok {
		if syscallErr, ok := opErr.Err.(*syscall.Errno); ok {
			if *syscallErr == syscall.ECONNREFUSED {
				return NewNetworkConnectionError(targetURL, err)
			}
		}
	}

	// Check for DNS errors
	if dnsErr, ok := err.(*net.DNSError); ok {
		return NewNetworkDNSError(dnsErr.Name, err)
	}

	// Check for SSL/TLS errors
	errStr := err.Error()
	if strings.Contains(errStr, "certificate") ||
		strings.Contains(errStr, "tls") ||
		strings.Contains(errStr, "x509") {
		return NewNetworkSSLError(targetURL, err)
	}

	// Check for DNS lookup errors (fallback)
	if strings.Contains(errStr, "no such host") ||
		strings.Contains(errStr, "lookup") {
		// Extract hostname from URL for better error context
		if u, parseErr := url.Parse(targetURL); parseErr == nil {
			return NewNetworkDNSError(u.Hostname(), err)
		}
		return NewNetworkDNSError(targetURL, err)
	}

	// Default to generic network connection error
	return NewNetworkConnectionError(targetURL, err)
}

// ClassifyHTTPStatusError converts HTTP status codes to domain-specific errors
// This function provides centralized error classification for HTTP status code failures
func ClassifyHTTPStatusError(targetURL string, statusCode int) error {
	switch statusCode {
	case http.StatusForbidden:
		return NewHTTPForbiddenError(targetURL)
	case http.StatusNotFound:
		return NewHTTPNotFoundError(targetURL)
	default:
		return NewHTTPNonOKStatusError(targetURL, statusCode)
	}
}
