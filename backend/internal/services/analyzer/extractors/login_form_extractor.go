package extractors

import (
	"net/url"

	"page-insight-tool/internal/models"

	"golang.org/x/net/html"
)

// LoginFormExtractor detects login forms in HTML documents
type LoginFormExtractor struct{}

// Name returns the extractor identifier
func (e *LoginFormExtractor) Name() string {
	return "login_form"
}

// Extract searches for forms containing password input fields to detect login forms
func (e *LoginFormExtractor) Extract(doc *html.Node, base *url.URL, result *models.AnalysisResponse, rawHTML string) {
	if result.HasLoginForm {
		return // Already found
	}

	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "form" {
			if formHasPasswordField(n) {
				result.HasLoginForm = true
				return
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
			if result.HasLoginForm {
				return
			}
		}
	}
	walk(doc)
}

// formHasPasswordField checks if a form contains a password field
func formHasPasswordField(form *html.Node) bool {
	var hasPassword bool
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "input" {
			inputType, _ := getAttr(n, "type")
			if inputType == "password" {
				hasPassword = true
				return
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
			if hasPassword {
				return
			}
		}
	}
	walk(form)
	return hasPassword
}
