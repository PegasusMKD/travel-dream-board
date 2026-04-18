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
	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/shared/constant"
)

const USER_AGENT = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
const ACCEPT_LANGUAGE = "en-US,en;q=0.9"

type Service interface {
	Scrape(ctx context.Context, url string) (*scrapeaudit.ScrapeResult, error)
}

type scrapeProcessServiceImpl struct {
	client *http.Client

	anthropicKey string

	scrapeAuditService scrapeaudit.Service
}

func NewService(anthropicKey string, scrapeAuditService scrapeaudit.Service) Service {
	transport := &http.Transport{
		MaxIdleConns:          10,
		IdleConnTimeout:       30 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
	}

	return &scrapeProcessServiceImpl{
		anthropicKey:       anthropicKey,
		scrapeAuditService: scrapeAuditService,
		client: &http.Client{
			Timeout:   15 * time.Second,
			Transport: transport,
			// Optional: Custom redirect policy
			// CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// 	 if len(via) >= 5 {
			// 		 return errors.New("stopped after 5 redirects")
			// 	 }
			// 	 return nil
			// },
		},
	}
}

func (svc *scrapeProcessServiceImpl) Scrape(ctx context.Context, url string) (*scrapeaudit.ScrapeResult, error) {
	out := scrapeaudit.ScrapeResult{
		InitialUrl: url,
	}

	host := getHostFromURL(url)
	auditRecord, err := svc.scrapeAuditService.CreateScrapeAudit(ctx, url, host)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Error("Failed creating request for scraping", "url", url, "error", err)
		svc.scrapeAuditService.UpdateScrapeAuditByUuid(ctx, auditRecord.Uuid, db.ScrapeStatusFailed, &out)
		return nil, err
	}

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

	out.ActualUrl = resp.Request.URL.String()

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

	client := anthropic.NewClient(option.WithAPIKey(s.anthropicKey))

	// 2. Define the Tool using the correct Union and Schema structures
	toolSchema := anthropic.ToolUnionParam{
		OfTool: &anthropic.ToolParam{
			Name:        "extract_page_details",
			Description: anthropic.String("Extracts the main title, description, and primary image URL."),
			InputSchema: anthropic.ToolInputSchemaParam{
				Type: constant.Object("object"),
				Properties: map[string]interface{}{
					"title": map[string]string{
						"type":        "string",
						"description": "The name of the accommodation, activity, or main product.",
					},
					"description": map[string]string{
						"type":        "string",
						"description": "A short 1-2 sentence description.",
					},
					"image_url": map[string]string{
						"type":        "string",
						"description": "The absolute URL to the main image.",
					},
				},
				Required: []string{"title", "description", "image_url"},
			},
		},
	}

	// 3. Force Claude to use the tool via ToolChoiceUnionParam
	toolChoice := anthropic.ToolChoiceUnionParam{
		OfTool: &anthropic.ToolChoiceToolParam{
			Type: constant.Tool("tool"),
			Name: "extract_page_details",
		},
	}

	// 4. Send the Request
	message, err := client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:      anthropic.ModelClaudeHaiku4_5,
		MaxTokens:  1024,
		Tools:      []anthropic.ToolUnionParam{toolSchema},
		ToolChoice: toolChoice,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock("Extract the details from this text: " + text)),
		},
	})
	if err != nil {
		log.Error("Failed with Haiku 4.5 fallback", "error", err)
		return
	}

	commitedUpdate := false

	// 5. Parse the Response
	for _, block := range message.Content {
		if block.Type == "tool_use" {
			var extracted AIExtraction

			inputBytes, _ := json.Marshal(block.Input)

			if err := json.Unmarshal(inputBytes, &extracted); err == nil {
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
				return
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
