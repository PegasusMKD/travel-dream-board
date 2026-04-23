package scrapeprocess_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/PegasusMKD/travel-dream-board/internal/db"
	"github.com/PegasusMKD/travel-dream-board/internal/scrape_audit"
	scrapeprocess "github.com/PegasusMKD/travel-dream-board/internal/scrape_process"
	"github.com/PegasusMKD/travel-dream-board/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockRoundTripper struct {
	roundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.roundTripFunc(req)
}

func TestScrapeProcessService_Scrape_ClaudeFallback(t *testing.T) {
	htmlContent := `
	<!DOCTYPE html>
	<html>
	<head></head>
	<body>
		<h1>Some text for Claude to parse</h1>
	</body>
	</html>
	`

	// Mock http.DefaultTransport
	originalTransport := http.DefaultTransport
	defer func() { http.DefaultTransport = originalTransport }()

	http.DefaultTransport = &mockRoundTripper{
		roundTripFunc: func(req *http.Request) (*http.Response, error) {
			if req.URL.Host == "api.scrapingant.com" {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(htmlContent)),
					Request:    req,
				}, nil
			}
			if req.URL.Host == "openrouter.ai" {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(bytes.NewBufferString(`{
						"id": "chatcmpl-123",
						"object": "chat.completion",
						"created": 1677652288,
						"model": "anthropic/claude-3-haiku",
						"choices": [{
							"index": 0,
							"message": {
								"role": "assistant",
								"content": null,
								"tool_calls": [{
									"id": "call_1",
									"type": "function",
									"function": {
										"name": "extract_page_details",
										"arguments": "{\"title\": \"Claude Title\", \"description\": \"Claude Description\", \"image_url\": \"http://example.com/claude.png\"}"
									}
								}]
							},
							"finish_reason": "tool_calls"
						}]
					}`)),
					Header: http.Header{
						"Content-Type": []string{"application/json"},
					},
				}, nil
			}
			return originalTransport.RoundTrip(req)
		},
	}

	ctx := context.Background()
	mockAudit := new(mocks.MockscrapeauditService)
	svc := scrapeprocess.NewService("dummy-key", "dummy-ant-key", mockAudit)

	auditRecord := &scrapeaudit.ScrapeAudit{Uuid: "audit-3"}
	fakeURL := "https://booking.com/some-hotel"

	mockAudit.On("CreateScrapeAudit", ctx, fakeURL, mock.AnythingOfType("string")).Return(auditRecord, nil)
	mockAudit.On("UpdateScrapeAuditByUuid", ctx, auditRecord.Uuid, db.ScrapeStatusCompletedByAi, mock.AnythingOfType("*scrapeaudit.ScrapeResult")).Return(nil)

	result, err := svc.Scrape(ctx, fakeURL)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Claude Title", *result.Title)
	assert.Equal(t, "Claude Description", *result.Description)
	assert.Equal(t, "http://example.com/claude.png", *result.ImageUrl)

	mockAudit.AssertExpectations(t)
}

func TestScrapeProcessService_Scrape_SuccessWithMeta(t *testing.T) {
	htmlContent := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Test Page</title>
		<meta property="og:title" content="OG Title">
		<meta property="og:description" content="OG Description">
		<meta property="og:image" content="http://example.com/image.png">
		<meta property="og:site_name" content="OG Site">
	</head>
	<body>
		<h1>Hello</h1>
	</body>
	</html>
	`

	originalTransport := http.DefaultTransport
	defer func() { http.DefaultTransport = originalTransport }()

	http.DefaultTransport = &mockRoundTripper{
		roundTripFunc: func(req *http.Request) (*http.Response, error) {
			if req.URL.Host == "api.scrapingant.com" {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(htmlContent)),
					Request:    req,
				}, nil
			}
			return originalTransport.RoundTrip(req)
		},
	}

	ctx := context.Background()
	mockAudit := new(mocks.MockscrapeauditService)
	svc := scrapeprocess.NewService("dummy-key", "dummy-ant-key", mockAudit)

	auditRecord := &scrapeaudit.ScrapeAudit{Uuid: "audit-1"}
	fakeURL := "https://booking.com/some-hotel"

	// Expected calls
	mockAudit.On("CreateScrapeAudit", ctx, fakeURL, mock.AnythingOfType("string")).Return(auditRecord, nil)

	// Expect an update to happen (since we parse metas)
	mockAudit.On("UpdateScrapeAuditByUuid", ctx, auditRecord.Uuid, db.ScrapeStatusCompletedByOg, mock.AnythingOfType("*scrapeaudit.ScrapeResult")).Return(nil)

	result, err := svc.Scrape(ctx, fakeURL)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "OG Title", *result.Title)
	assert.Equal(t, "OG Description", *result.Description)
	assert.Equal(t, "http://example.com/image.png", *result.ImageUrl)
	assert.Equal(t, "OG Site", *result.SiteName)
	assert.Equal(t, fakeURL, result.InitialUrl)

	mockAudit.AssertExpectations(t)
}

func TestScrapeProcessService_Scrape_SuccessWithJSONLD(t *testing.T) {
	htmlContent := `
	<!DOCTYPE html>
	<html>
	<head>
		<script type="application/ld+json">
		{
			"@context": "https://schema.org",
			"@type": "Product",
			"name": "JSONLD Title",
			"description": "JSONLD Description",
			"image": "http://example.com/jsonld.png"
		}
		</script>
	</head>
	<body>
		<h1>Hello</h1>
	</body>
	</html>
	`

	originalTransport := http.DefaultTransport
	defer func() { http.DefaultTransport = originalTransport }()

	http.DefaultTransport = &mockRoundTripper{
		roundTripFunc: func(req *http.Request) (*http.Response, error) {
			if req.URL.Host == "api.scrapingant.com" {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(htmlContent)),
					Request:    req,
				}, nil
			}
			return originalTransport.RoundTrip(req)
		},
	}

	ctx := context.Background()
	mockAudit := new(mocks.MockscrapeauditService)
	svc := scrapeprocess.NewService("dummy-key", "dummy-ant-key", mockAudit)

	auditRecord := &scrapeaudit.ScrapeAudit{Uuid: "audit-2"}
	fakeURL := "https://booking.com/some-hotel"

	mockAudit.On("CreateScrapeAudit", ctx, fakeURL, mock.AnythingOfType("string")).Return(auditRecord, nil)
	// Expect it to update with JsonLd (for jsonld properties)
	mockAudit.On("UpdateScrapeAuditByUuid", ctx, auditRecord.Uuid, db.ScrapeStatusCompletedByJsonLd, mock.AnythingOfType("*scrapeaudit.ScrapeResult")).Return(nil)

	result, err := svc.Scrape(ctx, fakeURL)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "JSONLD Title", *result.Title)
	assert.Equal(t, "JSONLD Description", *result.Description)
	assert.Equal(t, "http://example.com/jsonld.png", *result.ImageUrl)

	mockAudit.AssertExpectations(t)
}

func TestScrapeProcessService_Scrape_NetworkError(t *testing.T) {
	originalTransport := http.DefaultTransport
	defer func() { http.DefaultTransport = originalTransport }()

	http.DefaultTransport = &mockRoundTripper{
		roundTripFunc: func(req *http.Request) (*http.Response, error) {
			if req.URL.Host == "api.scrapingant.com" {
				return nil, errors.New("network timeout")
			}
			return originalTransport.RoundTrip(req)
		},
	}

	ctx := context.Background()
	mockAudit := new(mocks.MockscrapeauditService)
	svc := scrapeprocess.NewService("dummy-key", "dummy-ant-key", mockAudit)

	auditRecord := &scrapeaudit.ScrapeAudit{Uuid: "audit-1"}
	url := "http://invalid-url.invalid-tld"

	mockAudit.On("CreateScrapeAudit", ctx, url, mock.AnythingOfType("string")).Return(auditRecord, nil)
	// It should mark it as failed
	mockAudit.On("UpdateScrapeAuditByUuid", ctx, auditRecord.Uuid, db.ScrapeStatusFailed, mock.AnythingOfType("*scrapeaudit.ScrapeResult")).Return(nil)

	result, err := svc.Scrape(ctx, url)

	assert.Error(t, err)
	assert.Nil(t, result)

	mockAudit.AssertExpectations(t)
}
