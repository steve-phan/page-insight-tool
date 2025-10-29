package extractors

import (
	"net/url"
	"strings"
	"testing"

	"page-insight-tool/internal/models"

	"golang.org/x/net/html"
)

func TestHeadingsExtractor_Comprehensive(t *testing.T) {
	extractor := &HeadingsExtractor{}
	testURL, _ := url.Parse("https://example.com")

	tests := []struct {
		name     string
		html     string
		expected models.Headings
	}{
		{
			name:     "Empty HTML",
			html:     `<html><body></body></html>`,
			expected: models.Headings{},
		},
		{
			name:     "Single H1",
			html:     `<html><body><h1>Main Title</h1></body></html>`,
			expected: models.Headings{H1: 1},
		},
		{
			name:     "Multiple H1s",
			html:     `<html><body><h1>Title 1</h1><h1>Title 2</h1></body></html>`,
			expected: models.Headings{H1: 2},
		},
		{
			name:     "All heading levels",
			html:     `<html><body><h1>H1</h1><h2>H2</h2><h3>H3</h3><h4>H4</h4><h5>H5</h5><h6>H6</h6></body></html>`,
			expected: models.Headings{H1: 1, H2: 1, H3: 1, H4: 1, H5: 1, H6: 1},
		},
		{
			name:     "Mixed headings with content",
			html:     `<html><body><h1>Main</h1><p>Content</p><h2>Sub</h2><h3>Sub-sub</h3><div><h4>Nested</h4></div></body></html>`,
			expected: models.Headings{H1: 1, H2: 1, H3: 1, H4: 1},
		},
		{
			name:     "Headings with attributes",
			html:     `<html><body><h1 class="main">Title</h1><h2 id="sub">Subtitle</h2></body></html>`,
			expected: models.Headings{H1: 1, H2: 1},
		},
		{
			name:     "Case insensitive headings",
			html:     `<html><body><H1>Uppercase</H1><H2>Another</H2></body></html>`,
			expected: models.Headings{H1: 1, H2: 1},
		},
		{
			name:     "Headings with nested elements",
			html:     `<html><body><h1><span>Nested</span> <strong>Title</strong></h1><h2><em>Italic</em> Subtitle</h2></body></html>`,
			expected: models.Headings{H1: 1, H2: 1},
		},
		{
			name:     "Complex nested structure",
			html:     `<html><head><title>Test</title></head><body><div><section><h1>Section 1</h1><article><h2>Article 1</h2><h3>Sub Article</h3></article></section><section><h1>Section 2</h1><h2>Article 2</h2></section></div></body></html>`,
			expected: models.Headings{H1: 2, H2: 2, H3: 1},
		},
		{
			name:     "No headings",
			html:     `<html><body><p>Just paragraphs</p><div>And divs</div><span>And spans</span></body></html>`,
			expected: models.Headings{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := &models.AnalysisResponse{}
			doc, err := html.Parse(strings.NewReader(tt.html))
			if err != nil {
				t.Fatalf("Failed to parse HTML: %v", err)
			}

			extractor.Extract(doc, testURL, result, tt.html)

			if result.Headings != tt.expected {
				t.Errorf("HeadingsExtractor failed for %s: got %+v, want %+v",
					tt.name, result.Headings, tt.expected)
			}
		})
	}
}

func TestHeadingsExtractor_EdgeCases(t *testing.T) {
	extractor := &HeadingsExtractor{}
	testURL, _ := url.Parse("https://example.com")

	tests := []struct {
		name     string
		html     string
		expected models.Headings
	}{
		{
			name:     "Malformed HTML",
			html:     `<html><body><h1>Title</h1><h2>Subtitle</h2><h3>Broken</h3></body>`,
			expected: models.Headings{H1: 1, H2: 1, H3: 1},
		},
		{
			name:     "Self-closing headings (invalid but should be handled)",
			html:     `<html><body><h1/><h2/><h3/></body></html>`,
			expected: models.Headings{H1: 1, H2: 1, H3: 1},
		},
		{
			name:     "Headings with special characters",
			html:     `<html><body><h1>&lt;Special&gt; Characters</h1><h2>Unicode: 你好</h2></body></html>`,
			expected: models.Headings{H1: 1, H2: 1},
		},
		{
			name:     "Empty headings",
			html:     `<html><body><h1></h1><h2><span></span></h2><h3>   </h3></body></html>`,
			expected: models.Headings{H1: 1, H2: 1, H3: 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := &models.AnalysisResponse{}
			doc, err := html.Parse(strings.NewReader(tt.html))
			if err != nil {
				t.Fatalf("Failed to parse HTML: %v", err)
			}

			extractor.Extract(doc, testURL, result, tt.html)

			if result.Headings != tt.expected {
				t.Errorf("HeadingsExtractor failed for %s: got %+v, want %+v",
					tt.name, result.Headings, tt.expected)
			}
		})
	}
}

func TestHeadingsExtractor_Performance(t *testing.T) {
	extractor := &HeadingsExtractor{}
	testURL, _ := url.Parse("https://example.com")

	// Create HTML with many headings
	var htmlBuilder strings.Builder
	htmlBuilder.WriteString("<html><body>")
	for i := 0; i < 1000; i++ {
		htmlBuilder.WriteString("<h1>Heading ")
		htmlBuilder.WriteString(string(rune('0' + i%10)))
		htmlBuilder.WriteString("</h1>")
	}
	htmlBuilder.WriteString("</body></html>")

	htmlContent := htmlBuilder.String()
	result := &models.AnalysisResponse{}
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		t.Fatalf("Failed to parse HTML: %v", err)
	}

	extractor.Extract(doc, testURL, result, htmlContent)

	expected := models.Headings{H1: 1000}
	if result.Headings != expected {
		t.Errorf("HeadingsExtractor performance test failed: got %+v, want %+v",
			result.Headings, expected)
	}
}

func TestHeadingsExtractor_Idempotency(t *testing.T) {
	extractor := &HeadingsExtractor{}
	testURL, _ := url.Parse("https://example.com")
	htmlContent := `<html><body><h1>Title</h1><h2>Subtitle</h2></body></html>`
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		t.Fatalf("Failed to parse HTML: %v", err)
	}

	result := &models.AnalysisResponse{}

	// Run extractor multiple times - should accumulate counts
	for i := 0; i < 3; i++ {
		extractor.Extract(doc, testURL, result, htmlContent)
	}

	// After 3 runs, we should have 3x the counts
	expected := models.Headings{H1: 3, H2: 3}
	if result.Headings != expected {
		t.Errorf("HeadingsExtractor idempotency test failed: got %+v, want %+v",
			result.Headings, expected)
	}
}
