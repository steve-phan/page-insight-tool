package extractors

import (
	"net/url"

	"page-insight-tool/internal/models"

	"golang.org/x/net/html"
)

// HeadingsExtractor extracts heading counts from HTML documents
type HeadingsExtractor struct{}

// Name returns the extractor identifier
func (e *HeadingsExtractor) Name() string {
	return "headings"
}

// Extract counts all heading elements (h1-h6) in the HTML document
func (e *HeadingsExtractor) Extract(doc *html.Node, base *url.URL, result *models.AnalysisResponse, rawHTML string) {
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && isHeadingTag(n.Data) {
			incrementHeadingCount(n.Data, &result.Headings)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(doc)
}

// isHeadingTag checks if a tag is a heading tag
func isHeadingTag(tag string) bool {
	return tag == "h1" || tag == "h2" || tag == "h3" || tag == "h4" || tag == "h5" || tag == "h6"
}

func incrementHeadingCount(tag string, hs *models.Headings) {
	switch tag {
	case "h1":
		hs.H1++
	case "h2":
		hs.H2++
	case "h3":
		hs.H3++
	case "h4":
		hs.H4++
	case "h5":
		hs.H5++
	case "h6":
		hs.H6++
	}
}
