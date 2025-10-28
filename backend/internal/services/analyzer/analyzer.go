package services

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"page-insight-tool/internal/config"
	"page-insight-tool/internal/models"

	"golang.org/x/net/html"
)

// ---------- Errors ----------
var (
	ErrBodyTooLarge  = errors.New("response body exceeds limit")
	ErrInvalidURL    = errors.New("invalid URL")
	ErrHTMLParse     = errors.New("failed to parse HTML")
	ErrNonOKStatus   = errors.New("non-OK HTTP status")
	ErrRedirectLimit = errors.New("too many redirects")
)

type AnalyzerService struct {
	cfg        *config.Config
	httpClient *http.Client
	userAgent  string
}

func NewAnalyzerService(cfg *config.Config) *AnalyzerService {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		Proxy:                 http.ProxyFromEnvironment,
		MaxConnsPerHost:       20,
		IdleConnTimeout:       30 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	client := &http.Client{
		Timeout:       cfg.Analysis.Timeout * time.Second,
		Transport:     transport,
		CheckRedirect: redirectPolicy(5),
	}

	return &AnalyzerService{
		cfg:        cfg,
		httpClient: client,
		userAgent:  cfg.App.Name,
	}
}

func redirectPolicy(max int) func(*http.Request, []*http.Request) error {
	return func(req *http.Request, via []*http.Request) error {
		if len(via) >= max {
			return ErrRedirectLimit
		}
		return nil
	}
}

// ---------- Public API ----------
func (s *AnalyzerService) Analyze(ctx context.Context, rawURL string) (models.AnalysisResponse, error) {
	u, err := normalizeURL(rawURL)
	if err != nil {
		return models.AnalysisResponse{}, err
	}

	htmlContent, err := s.fetchHTML(ctx, u)
	if err != nil {
		return models.AnalysisResponse{}, err
	}

	return s.analyzeHTML(htmlContent, u)
}

func (s *AnalyzerService) fetchHTML(ctx context.Context, u *url.URL) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("User-Agent", s.userAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%w: %s", ErrNonOKStatus, resp.Status)
	}

	// Respect context while reading body
	bodyReader := io.LimitReader(resp.Body, int64(s.cfg.Analysis.MaxBodySize)*1024*1024+1)
	data, err := io.ReadAll(bodyReader)
	if err != nil {
		return "", fmt.Errorf("read body: %w", err)
	}
	if len(data) > int(s.cfg.Analysis.MaxBodySize*1024*1024) {
		return "", ErrBodyTooLarge
	}
	return string(data), nil
}

func (s *AnalyzerService) analyzeHTML(raw string, base *url.URL) (models.AnalysisResponse, error) {
	doc, err := html.Parse(strings.NewReader(raw))
	if err != nil {
		return models.AnalysisResponse{}, ErrHTMLParse
	}

	a := models.AnalysisResponse{
		HTMLVersion: detectHTMLVersion(doc, raw),
		Links:       models.Links{},
		Headings:    models.Headings{},
	}

	// Single tree walk
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type != html.ElementNode {
			goto next
		}

		switch n.Data {
		case "title":
			if a.PageTitle == "" && n.FirstChild != nil {
				a.PageTitle = normalizeText(n.FirstChild.Data)
			}
		case "a":
			s.processLink(n, base, &a)
		case "h1", "h2", "h3", "h4", "h5", "h6":
			s.incrementHeading(n.Data, &a.Headings)
		case "form":
			if !a.HasLoginForm {
				a.HasLoginForm = formHasPassword(n)
			}
		}

	next:
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(doc)

	return a, nil
}

func normalizeURL(raw string) (*url.URL, error) {
	if raw == "" {
		return nil, ErrInvalidURL
	}
	if !strings.Contains(raw, "://") {
		raw = "https://" + raw
	}
	u, err := url.Parse(raw)
	if err != nil || u.Host == "" {
		return nil, ErrInvalidURL
	}
	return u, nil
}

func normalizeText(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Join(strings.Fields(s), " ")
	return s
}

// detectHTMLVersion inspects the doctype and optional version attribute.
func detectHTMLVersion(doc *html.Node, raw string) string {
	// 1. Look for <!DOCTYPE …>
	if strings.Contains(strings.ToLower(raw[:min(200, len(raw))]), "<!doctype html>") {
		return "HTML5"
	}
	if strings.Contains(strings.ToLower(raw), "html 4.01") {
		return "HTML 4.01"
	}
	if strings.Contains(strings.ToLower(raw), "xhtml 1.0") {
		return "XHTML 1.0"
	}
	if strings.Contains(strings.ToLower(raw), "xhtml 1.1") {
		return "XHTML 1.1"
	}

	// 2. Fallback to <html version="…">
	if doc.Data == "html" {
		if v, ok := getAttr(doc, "version"); ok && v != "" {
			return v
		}
	}
	return "Unknown"
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// processLink classifies a single <a> element.
func (s *AnalyzerService) processLink(n *html.Node, base *url.URL, a *models.AnalysisResponse) {
	href, ok := getAttr(n, "href")
	if !ok || href == "" {
		a.Links.Inaccessible++
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
		a.Links.Inaccessible++
		return
	}
	parsed = base.ResolveReference(parsed)

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		a.Links.Inaccessible++
		return
	}

	if parsed.Host == "" { // relative, same host
		a.Links.Internal++
		return
	}
	if strings.EqualFold(parsed.Hostname(), base.Hostname()) {
		a.Links.Internal++
	} else {
		a.Links.External++
	}
}

// incrementHeading safely increments the correct counter.
func (s *AnalyzerService) incrementHeading(tag string, hs *models.Headings) {
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

// formHasPassword stops at the first password-type input.
func formHasPassword(form *html.Node) bool {
	var found bool
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if found {
			return
		}
		if n.Type == html.ElementNode && n.Data == "input" {
			typ, _ := getAttr(n, "type")
			name, _ := getAttr(n, "name")
			typ = strings.ToLower(typ)
			name = strings.ToLower(name)

			if typ == "password" ||
				strings.Contains(name, "pass") ||
				strings.Contains(name, "pwd") {
				found = true
				return
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(form)
	return found
}

func getAttr(n *html.Node, key string) (string, bool) {
	for _, a := range n.Attr {
		if a.Key == key {
			return a.Val, true
		}
	}
	return "", false
}
