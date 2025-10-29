package extractors

import (
	"net/url"
	"strings"

	"github.com/steve-phan/page-insight-tool/internal/models"

	"golang.org/x/net/html"
)

// VersionExtractor extracts HTML version from DOCTYPE
type VersionExtractor struct{}

// Name returns the extractor identifier
func (e *VersionExtractor) Name() string {
	return "version"
}

// Extract detects the HTML version from DOCTYPE
func (e *VersionExtractor) Extract(doc *html.Node, base *url.URL, result *models.AnalysisResponse, rawHTML string) {
	if result.HTMLVersion != "" {
		return // Already found
	}

	result.HTMLVersion = DetectHTMLVersion(rawHTML)
}

// DetectHTMLVersion detects the HTML version from DOCTYPE in raw HTML
func DetectHTMLVersion(raw string) string {
	// Check for DOCTYPE in raw HTML
	rawLower := strings.ToLower(raw)

	if strings.Contains(rawLower, "<!doctype html>") {
		return "HTML5"
	}
	if strings.Contains(rawLower, "<!doctype html 4.01") {
		return "HTML 4.01"
	}
	if strings.Contains(rawLower, "<!doctype html 4.0") {
		return "HTML 4.0"
	}
	if strings.Contains(rawLower, "<!doctype html 3.2") {
		return "HTML 3.2"
	}
	if strings.Contains(rawLower, "<!doctype html 2.0") {
		return "HTML 2.0"
	}
	if strings.Contains(rawLower, "<!doctype") {
		return "HTML 4.01" // Default for older DOCTYPEs
	}

	return "Unknown"
}
