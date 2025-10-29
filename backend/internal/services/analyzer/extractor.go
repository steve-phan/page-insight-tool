package analyzer

import (
	"net/url"

	"page-insight-tool/internal/models"

	"golang.org/x/net/html"
)

// Extractor defines the interface for HTML analysis extractors
// Each extractor is responsible for extracting specific information from HTML documents
type Extractor interface {
	// Name returns a unique identifier for this extractor
	Name() string

	// Extract processes the HTML document and populates the result with extracted data
	// The extractor should be idempotent and safe to call multiple times
	Extract(doc *html.Node, base *url.URL, result *models.AnalysisResponse, rawHTML string)
}
