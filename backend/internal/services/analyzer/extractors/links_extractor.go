package extractors

import (
	"net/url"
	"strings"

	"github.com/steve-phan/page-insight-tool/internal/models"

	"golang.org/x/net/html"
)

// LinksExtractor extracts link information from HTML documents
type LinksExtractor struct{}

// Name returns the extractor identifier
func (e *LinksExtractor) Name() string {
	return "links"
}

// Extract analyzes all anchor tags and categorizes them as internal, external, or inaccessible
func (e *LinksExtractor) Extract(doc *html.Node, base *url.URL, result *models.AnalysisResponse, rawHTML string) {
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			processLinkElement(n, base, &result.Links)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(doc)
}

// processLinkElement classifies a single <a> element
func processLinkElement(n *html.Node, base *url.URL, a *models.Links) {
	href, ok := getAttr(n, "href")
	if !ok || href == "" {
		a.Inaccessible++
		return
	}

	// Skip non-navigable protocols early
	if strings.HasPrefix(href, "#") ||
		strings.HasPrefix(href, "javascript:") ||
		strings.HasPrefix(href, "mailto:") ||
		strings.HasPrefix(href, "tel:") {
		return
	}

	parsed, err := url.Parse(href)
	if err != nil {
		a.Inaccessible++
		return
	}
	parsed = base.ResolveReference(parsed)

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		a.Inaccessible++
		return
	}

	if parsed.Host == "" { // relative, same host
		a.Internal++
		return
	}
	if strings.EqualFold(parsed.Hostname(), base.Hostname()) {
		a.Internal++
	} else {
		a.External++
	}
}

func getAttr(n *html.Node, key string) (string, bool) {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val, true
		}
	}
	return "", false
}
