package analyzer

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"page-insight-tool/internal/config"
	"page-insight-tool/internal/models"
	"page-insight-tool/internal/services/analyzer/extractors"

	"golang.org/x/net/html"
)

func TestAnalyzerService_BasicAnalysis(t *testing.T) {
	cfg := &config.Config{
		Analysis: config.AnalysisConfig{
			Timeout:     30,
			VerifySSL:   false,
			MaxBodySize: 10,
		},
		App: config.AppConfig{
			Name: "Test App",
		},
	}

	// Test with explicit extractors for basic analysis
	service, err := NewAnalyzerService(cfg,
		WithExtractors(
			&extractors.VersionExtractor{},
			&extractors.TitleExtractor{},
			&extractors.LoginFormExtractor{},
		),
	)
	if err != nil {
		t.Fatalf("Failed to create analyzer service: %v", err)
	}

	// Read test HTML file
	htmlPath := filepath.Join("testdata", "sample.html")
	htmlContent, err := os.ReadFile(htmlPath)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	testURL, _ := url.Parse("https://example.com")
	result, err := service.analyzeHTML(string(htmlContent), testURL)
	if err != nil {
		t.Fatalf("Analyze failed: %v", err)
	}

	// Verify basic analysis
	if result.HTMLVersion != "HTML5" {
		t.Errorf("HTMLVersion = %v, want HTML5", result.HTMLVersion)
	}
	if result.PageTitle != "Test Page" {
		t.Errorf("PageTitle = %v, want 'Test Page'", result.PageTitle)
	}

	if !result.HasLoginForm {
		t.Errorf("HasLoginForm = %v, want true", result.HasLoginForm)
	}
}

func TestAnalyzerService_CustomExtractors(t *testing.T) {
	cfg := &config.Config{
		Analysis: config.AnalysisConfig{
			Timeout:     30,
			VerifySSL:   false,
			MaxBodySize: 10,
		},
		App: config.AppConfig{
			Name: "Test App",
		},
	}

	// Test with custom extractors
	service, err := NewAnalyzerService(cfg,
		WithExtractors(
			&extractors.TitleExtractor{},
			&extractors.HeadingsExtractor{},
			// Only title and headings, no links or login forms
		),
	)
	if err != nil {
		t.Fatalf("Failed to create analyzer service: %v", err)
	}

	htmlPath := filepath.Join("testdata", "sample.html")
	htmlContent, err := os.ReadFile(htmlPath)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	testURL, _ := url.Parse("https://example.com")
	result, err := service.analyzeHTML(string(htmlContent), testURL)
	if err != nil {
		t.Fatalf("Analyze failed: %v", err)
	}

	// Should have title and headings
	if result.PageTitle == "" {
		t.Errorf("PageTitle should not be empty")
	}
	if result.Headings.H1 == 0 {
		t.Errorf("Should have H1 headings")
	}

	// Should NOT have links or login form (not configured)
	if result.Links.Internal != 0 || result.Links.External != 0 {
		t.Errorf("Should not have links when LinksExtractor is not configured")
	}
	if result.HasLoginForm {
		t.Errorf("Should not detect login form when LoginFormExtractor is not configured")
	}
}
func TestAnalyzerService_EmptyExtractors(t *testing.T) {
	cfg := &config.Config{
		Analysis: config.AnalysisConfig{
			Timeout:     30,
			VerifySSL:   false,
			MaxBodySize: 10,
		},
		App: config.AppConfig{
			Name: "Test App",
		},
	}

	// Test with no extractors - should fail fast
	_, err := NewAnalyzerService(cfg,
		WithExtractors(), // Empty extractors
	)
	if err == nil {
		t.Errorf("Expected error when creating service with no extractors, but got none")
	}
	if !strings.Contains(err.Error(), "no extractors configured") {
		t.Errorf("Expected specific error message about extractors, got: %v", err)
	}
}

// Test individual extractors
func TestTitleExtractor(t *testing.T) {
	extractor := &extractors.TitleExtractor{}
	result := &models.AnalysisResponse{}
	testURL, _ := url.Parse("https://example.com")

	// Test HTML with title
	htmlContent := `<html><head><title>Test Title</title></head><body></body></html>`
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		t.Fatalf("Failed to parse HTML: %v", err)
	}

	extractor.Extract(doc, testURL, result, htmlContent)
	if result.PageTitle != "Test Title" {
		t.Errorf("TitleExtractor failed: got %v, want 'Test Title'", result.PageTitle)
	}
}

func TestHeadingsExtractor(t *testing.T) {
	extractor := &extractors.HeadingsExtractor{}
	result := &models.AnalysisResponse{}
	testURL, _ := url.Parse("https://example.com")

	// Test HTML with headings
	htmlContent := `<html><body><h1>H1</h1><h2>H2</h2><h3>H3</h3></body></html>`
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		t.Fatalf("Failed to parse HTML: %v", err)
	}

	extractor.Extract(doc, testURL, result, htmlContent)
	expected := models.Headings{H1: 1, H2: 1, H3: 1}
	if result.Headings != expected {
		t.Errorf("HeadingsExtractor failed: got %+v, want %+v", result.Headings, expected)
	}
}

func TestLoginFormExtractor(t *testing.T) {
	extractor := &extractors.LoginFormExtractor{}
	result := &models.AnalysisResponse{}
	testURL, _ := url.Parse("https://example.com")

	// Test HTML with login form
	htmlContent := `<html><body><form><input type="password" name="pass"></form></body></html>`
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		t.Fatalf("Failed to parse HTML: %v", err)
	}

	extractor.Extract(doc, testURL, result, htmlContent)
	if !result.HasLoginForm {
		t.Errorf("LoginFormExtractor failed: should detect login form")
	}
}
