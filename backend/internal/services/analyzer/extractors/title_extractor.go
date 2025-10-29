package extractors

import (
	"net/url"
	"strings"

	"github.com/steve-phan/page-insight-tool/internal/models"

	"golang.org/x/net/html"
)

// TitleExtractor extracts page title from HTML documents
type TitleExtractor struct{}

// Name returns the extractor identifier
func (e *TitleExtractor) Name() string {
	return "title"
}

// Extract finds and extracts the page title from the HTML document
func (e *TitleExtractor) Extract(doc *html.Node, base *url.URL, result *models.AnalysisResponse, rawHTML string) {
	if result.PageTitle != "" {
		return // Already found
	}

	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" && n.FirstChild != nil {
			result.PageTitle = normalizeText(n.FirstChild.Data)
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
			if result.PageTitle != "" {
				return
			}
		}
	}
	walk(doc)
}

// normalizeText cleans and normalizes text content
func normalizeText(text string) string {
	// Remove extra whitespace and normalize
	text = strings.TrimSpace(text)
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ReplaceAll(text, "\t", " ")
	// Replace multiple spaces with single space
	for strings.Contains(text, "  ") {
		text = strings.ReplaceAll(text, "  ", " ")
	}
	return text
}
