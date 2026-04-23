package scrapeprocess

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PegasusMKD/travel-dream-board/internal/db"
	scrapeaudit "github.com/PegasusMKD/travel-dream-board/internal/scrape_audit"
	"github.com/PuerkitoBio/goquery"
	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

const USER_AGENT = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
const ACCEPT_LANGUAGE = "en-US,en;q=0.9"

type Service interface {
	Scrape(ctx context.Context, url string) (*scrapeaudit.ScrapeResult, error)
}

type scrapeProcessServiceImpl struct {
	client *http.Client

	openrouterKey  string
	scrapingAntKey string

	scrapeAuditService scrapeaudit.Service
}

func NewService(openrouterKey string, scrapingAntKey string, scrapeAuditService scrapeaudit.Service) Service {
	return &scrapeProcessServiceImpl{
		openrouterKey:      openrouterKey,
		scrapingAntKey:     scrapingAntKey,
		scrapeAuditService: scrapeAuditService,
		client: &http.Client{
			Timeout:   60 * time.Second,
			Transport: http.DefaultTransport,
		},
	}
}

func (svc *scrapeProcessServiceImpl) Scrape(ctx context.Context, url string) (*scrapeaudit.ScrapeResult, error) {
	out := scrapeaudit.ScrapeResult{
		InitialUrl: url,
	}

	host := getHostFromURL(url)
	auditRecord, err := svc.scrapeAuditService.CreateScrapeAudit(ctx, url, host)

	reqUrl := "https://api.scrapingant.com/v2/general"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqUrl, nil)
	if err != nil {
		log.Error("Failed creating request for scraping", "url", url, "error", err)
		svc.scrapeAuditService.UpdateScrapeAuditByUuid(ctx, auditRecord.Uuid, db.ScrapeStatusFailed, &out)
		return nil, err
	}

	q := req.URL.Query()
	q.Add("url", url)
	q.Add("x-api-key", svc.scrapingAntKey)
	q.Add("browser", "true")
	req.URL.RawQuery = q.Encode()

	req.Header.Set("User-Agent", USER_AGENT)
	req.Header.Set("Accept-Language", ACCEPT_LANGUAGE)

	resp, err := svc.client.Do(req)
	if err != nil {
		log.Error("Failed fetching HTML page", "url", url, "error", err)
		svc.scrapeAuditService.UpdateScrapeAuditByUuid(ctx, auditRecord.Uuid, db.ScrapeStatusFailed, &out)
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Error("Failed processing HTML page", "url", url, "error", err)
		svc.scrapeAuditService.UpdateScrapeAuditByUuid(ctx, auditRecord.Uuid, db.ScrapeStatusFailed, &out)
		return nil, err
	}

	out.ActualUrl = url

	svc.parseMetaTags(ctx, auditRecord.Uuid, doc, &out)

	svc.parseJSONLD(ctx, auditRecord.Uuid, doc, &out)

	if out.Title == nil || *out.Title == "" || out.ImageUrl == nil || *out.ImageUrl == "" || out.Description == nil || *out.Description == "" {
		svc.fallbackToClaude(ctx, auditRecord.Uuid, doc, &out)
	}

	return &out, nil
}

func (svc *scrapeProcessServiceImpl) parseMetaTags(ctx context.Context, uuid string, doc *goquery.Document, out *scrapeaudit.ScrapeResult) {
	commitedUpdate := false

	// Fallback standard title
	if title := doc.Find("title").First().Text(); title != "" {
		commitedUpdate = true
		title := strings.TrimSpace(title)
		out.Title = &title
	}

	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		property, _ := s.Attr("property")
		name, _ := s.Attr("name")
		content, _ := s.Attr("content")
		content = strings.TrimSpace(content)
		// OG Tags
		if property == "og:title" && content != "" {
			commitedUpdate = true
			out.Title = &content
		} else if property == "og:description" && content != "" {
			commitedUpdate = true
			out.Description = &content
		} else if property == "og:image" && content != "" {
			commitedUpdate = true
			out.ImageUrl = &content
		} else if property == "og:site_name" && content != "" {
			commitedUpdate = true
			out.SiteName = &content
		}
		// Standard Meta Tags Backup
		if name == "description" && (out.Description == nil || *out.Description == "") {
			commitedUpdate = true
			out.Description = &content
		}
	})

	if commitedUpdate {
		svc.scrapeAuditService.UpdateScrapeAuditByUuid(ctx, uuid, db.ScrapeStatusCompletedByOg, out)
	}

}

func (svc *scrapeProcessServiceImpl) parseJSONLD(ctx context.Context, uuid string, doc *goquery.Document, out *scrapeaudit.ScrapeResult) {
	commitedUpdate := false

	doc.Find(`script[type="application/ld+json"]`).Each(func(i int, s *goquery.Selection) {
		var data map[string]any

		if err := json.Unmarshal([]byte(s.Text()), &data); err != nil {
			return
		}

		if name, ok := data["name"].(string); ok && (out.Title == nil || *out.Title == "") {
			commitedUpdate = true
			out.Title = &name
		}

		if description, ok := data["description"].(string); ok && (out.Description == nil || *out.Description == "") {
			commitedUpdate = true
			out.Description = &description
		}

		if img, ok := data["image"].(string); ok && (out.ImageUrl == nil || *out.ImageUrl == "") {
			commitedUpdate = true
			out.ImageUrl = &img
		}
	})

	if commitedUpdate {
		svc.scrapeAuditService.UpdateScrapeAuditByUuid(ctx, uuid, db.ScrapeStatusCompletedByJsonLd, out)
	}
}

func (s *scrapeProcessServiceImpl) fallbackToClaude(ctx context.Context, uuid string, doc *goquery.Document, out *scrapeaudit.ScrapeResult) {
	// 1. Remove noise (scripts, styles, svgs) to save tokens
	doc.Find("script, style, svg, nav, footer").Remove()

	text := doc.Text()
	// Clean up excess whitespace and truncate to approx ~15,000 characters
	text = strings.Join(strings.Fields(text), " ")
	if len(text) > 15000 {
		text = text[:15000]
	}

	config := openai.DefaultConfig(s.openrouterKey)
	config.BaseURL = "https://openrouter.ai/api/v1"
	client := openai.NewClientWithConfig(config)

	params := openai.ChatCompletionRequest{
		Model:     "anthropic/claude-3-haiku",
		MaxTokens: 1024,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: "Extract the details from this text: " + text,
			},
		},
		Tools: []openai.Tool{
			{
				Type: openai.ToolTypeFunction,
				Function: &openai.FunctionDefinition{
					Name:        "extract_page_details",
					Description: "Extracts the main title, description, and primary image URL.",
					Parameters: jsonschema.Definition{
						Type: jsonschema.Object,
						Properties: map[string]jsonschema.Definition{
							"title": {
								Type:        jsonschema.String,
								Description: "The name of the accommodation, activity, or main product.",
							},
							"description": {
								Type:        jsonschema.String,
								Description: "A short 1-2 sentence description.",
							},
							"image_url": {
								Type:        jsonschema.String,
								Description: "The absolute URL to the main image.",
							},
						},
						Required: []string{"title", "description", "image_url"},
					},
				},
			},
		},
		ToolChoice: map[string]interface{}{
			"type": "function",
			"function": map[string]string{
				"name": "extract_page_details",
			},
		},
	}

	resp, err := client.CreateChatCompletion(ctx, params)
	if err != nil {
		log.Error("Failed with Haiku fallback", "error", err)
		return
	}

	commitedUpdate := false

	if len(resp.Choices) > 0 && len(resp.Choices[0].Message.ToolCalls) > 0 {
		toolCall := resp.Choices[0].Message.ToolCalls[0]
		if toolCall.Function.Name == "extract_page_details" {
			var extracted AIExtraction
			if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &extracted); err == nil {
				if (out.Title == nil || *out.Title == "") && extracted.Title != "" {
					commitedUpdate = true
					out.Title = &extracted.Title
				}
				if (out.Description == nil || *out.Description == "") && extracted.Description != "" {
					commitedUpdate = true
					out.Description = &extracted.Description
				}
				if (out.ImageUrl == nil || *out.ImageUrl == "") && extracted.ImageURL != "" {
					commitedUpdate = true
					out.ImageUrl = &extracted.ImageURL
				}
			}
		}
	}

	if commitedUpdate {
		s.scrapeAuditService.UpdateScrapeAuditByUuid(ctx, uuid, db.ScrapeStatusCompletedByAi, out)
	}
}

func getHostFromURL(URL string) string {
	parsedURL, err := url.Parse(URL)
	if err != nil {
		return ""
	}
	host := parsedURL.Hostname()
	host = strings.TrimPrefix(host, "www.")
	return host
}
