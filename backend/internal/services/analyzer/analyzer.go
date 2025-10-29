package analyzer

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/steve-phan/page-insight-tool/internal/config"
	domainerrors "github.com/steve-phan/page-insight-tool/internal/errors"
	"github.com/steve-phan/page-insight-tool/internal/models"

	"golang.org/x/net/html"
)

// AnalysisOption defines a function that configures an analyzer
type AnalysisOption func(*AnalyzerConfig)

// AnalyzerConfig holds configuration for the analyzer
type AnalyzerConfig struct {
	extractors []Extractor
}

// AnalyzerService uses functional options for extensible analysis
type AnalyzerService struct {
	cfg        *config.Config
	httpClient *http.Client
	userAgent  string
	config     *AnalyzerConfig
}

// NewAnalyzerService creates a new analyzer service
func NewAnalyzerService(cfg *config.Config, options ...AnalysisOption) (*AnalyzerService, error) {
	// Default configuration
	analyzerConfig := &AnalyzerConfig{
		extractors: []Extractor{},
	}

	// Apply options
	for _, option := range options {
		option(analyzerConfig)
	}

	if len(analyzerConfig.extractors) == 0 {
		return nil, fmt.Errorf("no extractors configured: at least one extractor must be provided")
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: !cfg.Analysis.VerifySSL,
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
		config:     analyzerConfig,
	}, nil
}

// Functional Options

// WithExtractors configures custom extractors
func WithExtractors(extractors ...Extractor) AnalysisOption {
	return func(config *AnalyzerConfig) {
		config.extractors = extractors
	}
}

// Analyze performs HTML analysis using configured extractors
func (s *AnalyzerService) Analyze(ctx context.Context, rawURL string) (models.AnalysisResponse, error) {
	start := time.Now()

	u, err := normalizeURL(rawURL)
	if err != nil {
		return models.AnalysisResponse{}, err
	}

	htmlContent, err := s.fetchHTML(ctx, u)
	if err != nil {
		return models.AnalysisResponse{}, err
	}

	result, err := s.analyzeHTML(htmlContent, u)
	if err != nil {
		return models.AnalysisResponse{}, err
	}

	result.AnalysisTime = int64(time.Since(start) / time.Millisecond)
	return result, nil
}

// analyzeHTML performs analysis using configured extractors
func (s *AnalyzerService) analyzeHTML(raw string, base *url.URL) (models.AnalysisResponse, error) {
	doc, err := html.Parse(strings.NewReader(raw))
	if err != nil {
		return models.AnalysisResponse{}, domainerrors.NewHTMLParseError(base.String(), err)
	}

	result := models.AnalysisResponse{
		Links:    models.Links{},
		Headings: models.Headings{},
	}

	// Run all configured extractors
	for _, extractor := range s.config.extractors {
		extractor.Extract(doc, base, &result, raw)
	}

	return result, nil
}

// fetchHTML fetches HTML content (shared implementation)
func (s *AnalyzerService) fetchHTML(ctx context.Context, u *url.URL) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return "", domainerrors.ClassifyNetworkError(u.String(), err)
	}
	req.Header.Set("User-Agent", s.userAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", domainerrors.ClassifyNetworkError(u.String(), err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", domainerrors.ClassifyHTTPStatusError(u.String(), resp.StatusCode)
	}

	bodyReader := io.LimitReader(resp.Body, int64(s.cfg.Analysis.MaxBodySize)*1024*1024+1)
	data, err := io.ReadAll(bodyReader)
	if err != nil {
		return "", domainerrors.ClassifyNetworkError(u.String(), err)
	}
	if len(data) > int(s.cfg.Analysis.MaxBodySize*1024*1024) {
		return "", domainerrors.NewContentTooBigError(u.String(), int64(len(data)), int64(s.cfg.Analysis.MaxBodySize*1024*1024))
	}
	return string(data), nil
}

// Helper functions for analyzer

// normalizeURL normalizes and validates a URL
func normalizeURL(rawURL string) (*url.URL, error) {
	if rawURL == "" {
		return nil, domainerrors.NewInvalidURLError(rawURL, nil)
	}

	// Add https:// if no scheme is provided
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		rawURL = "https://" + rawURL
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, domainerrors.NewInvalidURLError(rawURL, err)
	}

	if u.Scheme == "" || u.Host == "" {
		return nil, domainerrors.NewInvalidURLError(rawURL, nil)
	}

	return u, nil
}

// redirectPolicy creates a redirect policy function
func redirectPolicy(maxRedirects int) func(req *http.Request, via []*http.Request) error {
	return func(req *http.Request, via []*http.Request) error {
		if len(via) >= maxRedirects {
			return domainerrors.NewInternalError(fmt.Sprintf("stopped after %d redirects", maxRedirects), nil)
		}
		return nil
	}
}
