package validation

import (
	"net"
	"net/url"
	"strings"

	"page-insight-tool/internal/errors"
)

// URLValidator handles URL validation logic
type URLValidator struct{}

// NewURLValidator creates a new URL validator
func NewURLValidator() *URLValidator {
	return &URLValidator{}
}

// ValidateURL validates a URL for web scraping purposes
func (v *URLValidator) ValidateURL(rawURL string) error {
	if rawURL == "" {
		return errors.ErrMissingURL
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return errors.NewInvalidURLError(rawURL, err)
	}

	// Must have valid scheme (http or https)
	if u.Scheme != "http" && u.Scheme != "https" {
		return errors.NewInvalidURLError(rawURL, nil)
	}

	// Must have valid host
	if u.Host == "" {
		return errors.NewInvalidURLError(rawURL, nil)
	}

	// Extract hostname (remove port if present)
	hostname := u.Hostname()
	if hostname == "" {
		return errors.NewInvalidURLError(rawURL, nil)
	}

	// Check for valid hostname format
	if err := v.validateHostname(hostname); err != nil {
		return errors.NewInvalidURLError(rawURL, err)
	}

	return nil
}

// validateHostname checks if a hostname is valid according to RFC standards
func (v *URLValidator) validateHostname(hostname string) error {
	// Empty hostname is invalid
	if hostname == "" {
		return errors.NewInvalidInputError("hostname", hostname, "hostname cannot be empty")
	}

	// Hostname too long
	if len(hostname) > 253 {
		return errors.NewInvalidInputError("hostname", hostname, "hostname too long (max 253 characters)")
	}

	// Check for valid IP address (both IPv4 and IPv6)
	if net.ParseIP(hostname) != nil {
		return nil // IP addresses are valid
	}

	// For domain names, check basic format requirements
	hostname = strings.TrimSuffix(hostname, ".") // Remove trailing dot if present

	// Split into labels and validate each
	labels := strings.Split(hostname, ".")
	if len(labels) < 2 {
		// Single label hostnames like "localhost" are valid, but for web scraping
		// we want at least a TLD for practical purposes
		// However, "localhost" should be allowed for testing
		if hostname != "localhost" && !strings.HasSuffix(hostname, ".local") {
			return errors.NewInvalidInputError("hostname", hostname, "hostname must have at least a domain and TLD (e.g., example.com)")
		}
	}

	for i, label := range labels {
		if err := v.validateLabel(label, i); err != nil {
			return err
		}
	}

	return nil
}

// validateLabel checks if a DNS label is valid
func (v *URLValidator) validateLabel(label string, position int) error {
	if label == "" {
		return errors.NewInvalidInputError("hostname_label", label, "empty label in hostname")
	}

	if len(label) > 63 {
		return errors.NewInvalidInputError("hostname_label", label, "label too long (max 63 characters)")
	}

	// Must start and end with alphanumeric character
	if !isAlphaNumeric(label[0]) {
		return errors.NewInvalidInputError("hostname_label", label, "label must start with alphanumeric character")
	}

	if !isAlphaNumeric(label[len(label)-1]) {
		return errors.NewInvalidInputError("hostname_label", label, "label must end with alphanumeric character")
	}

	// Check all characters
	for _, r := range label {
		if !isAlphaNumeric(byte(r)) && r != '-' {
			return errors.NewInvalidInputError("hostname_label", label, "label contains invalid characters (only alphanumeric and hyphens allowed)")
		}
	}

	return nil
}

// isAlphaNumeric checks if a byte is alphanumeric
func isAlphaNumeric(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9')
}
