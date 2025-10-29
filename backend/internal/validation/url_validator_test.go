package validation

import (
	"strings"
	"testing"
)

func TestURLValidator_ValidateURL(t *testing.T) {
	validator := NewURLValidator()

	tests := []struct {
		name      string
		url       string
		wantError bool
		errorType string
	}{
		{
			name:      "valid domain",
			url:       "https://google.com",
			wantError: false,
		},
		{
			name:      "valid domain with path",
			url:       "https://example.com/path?query=1",
			wantError: false,
		},
		{
			name:      "valid localhost",
			url:       "https://localhost",
			wantError: false,
		},
		{
			name:      "valid localhost with port",
			url:       "http://localhost:3000",
			wantError: false,
		},
		{
			name:      "valid IPv4",
			url:       "https://127.0.0.1",
			wantError: false,
		},
		{
			name:      "valid IPv6",
			url:       "https://[::1]",
			wantError: false,
		},
		{
			name:      "valid local domain",
			url:       "https://test.local",
			wantError: false,
		},
		{
			name:      "invalid single letter domain",
			url:       "https://p",
			wantError: true,
			errorType: "INVALID_URL",
		},
		{
			name:      "invalid single word domain",
			url:       "https://test",
			wantError: true,
			errorType: "INVALID_URL",
		},
		{
			name:      "invalid no scheme",
			url:       "google.com",
			wantError: true,
			errorType: "INVALID_URL",
		},
		{
			name:      "invalid ftp scheme",
			url:       "ftp://example.com",
			wantError: true,
			errorType: "INVALID_URL",
		},
		{
			name:      "empty URL",
			url:       "",
			wantError: true,
		},
		{
			name:      "invalid URL format",
			url:       "not-a-url",
			wantError: true,
			errorType: "INVALID_URL",
		},
		{
			name:      "URL with spaces",
			url:       "https://example .com",
			wantError: true,
			errorType: "INVALID_URL",
		},
		{
			name:      "URL with umlauts",
			url:       "https://exämple.com",
			wantError: true,
			errorType: "INVALID_URL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateURL(tt.url)

			if tt.wantError {
				if err == nil {
					t.Errorf("ValidateURL() expected error for %s, got nil", tt.url)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateURL() unexpected error for %s: %v", tt.url, err)
				}
			}
		})
	}
}

func TestURLValidator_ValidateHostname(t *testing.T) {
	validator := NewURLValidator()

	tests := []struct {
		name      string
		hostname  string
		wantError bool
	}{
		{
			name:      "valid domain",
			hostname:  "example.com",
			wantError: false,
		},
		{
			name:      "valid subdomain",
			hostname:  "www.example.com",
			wantError: false,
		},
		{
			name:      "valid localhost",
			hostname:  "localhost",
			wantError: false,
		},
		{
			name:      "valid local domain",
			hostname:  "test.local",
			wantError: false,
		},
		{
			name:      "valid IPv4",
			hostname:  "192.168.1.1",
			wantError: false,
		},
		{
			name:      "invalid single letter",
			hostname:  "p",
			wantError: true,
		},
		{
			name:      "invalid single word",
			hostname:  "test",
			wantError: true,
		},
		{
			name:      "empty hostname",
			hostname:  "",
			wantError: true,
		},
		{
			name:      "hostname too long",
			hostname:  "a" + strings.Repeat(".very-long-label-name-that-exceeds-limits", 10) + ".com",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.validateHostname(tt.hostname)

			if tt.wantError {
				if err == nil {
					t.Errorf("validateHostname() expected error for %s, got nil", tt.hostname)
				}
			} else {
				if err != nil {
					t.Errorf("validateHostname() unexpected error for %s: %v", tt.hostname, err)
				}
			}
		})
	}
}
