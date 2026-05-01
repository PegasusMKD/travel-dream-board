package scrapeprocess

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PegasusMKD/travel-dream-board/internal/db"
	scrapeaudit "github.com/PegasusMKD/travel-dream-board/internal/scrape_audit"
	"github.com/PegasusMKD/travel-dream-board/internal/utility"
	"github.com/PuerkitoBio/goquery"
	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

const USER_AGENT = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
const ACCEPT_LANGUAGE = "en-US,en;q=0.9"

type Service interface {
	Scrape(ctx context.Context, url string) (*scrapeaudit.ScrapeResult, error)
	ExtractFromImage(ctx context.Context, url string, imageBytes []byte) (*AIExtraction, error)
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

	isScrapingAntError := resp.StatusCode >= 400 || strings.Contains(doc.Text(), "Free subscription plan")

	if isScrapingAntError {
		// Fallback to direct HTTP get
		directReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err == nil {
			directReq.Header.Set("User-Agent", USER_AGENT)
			directReq.Header.Set("Accept-Language", ACCEPT_LANGUAGE)
			if directResp, err := svc.client.Do(directReq); err == nil {
				defer directResp.Body.Close()
				if directDoc, err := goquery.NewDocumentFromReader(directResp.Body); err == nil {
					doc = directDoc
				}
			}
		}
	}

	svc.parseMetaTags(ctx, auditRecord.Uuid, doc, &out)

	svc.parseJSONLD(ctx, auditRecord.Uuid, doc, &out)

	if out.Title == nil || *out.Title == "" || out.ImageUrl == nil || *out.ImageUrl == "" || out.Description == nil || *out.Description == "" {
		svc.fallbackToClaude(ctx, auditRecord.Uuid, doc, url, &out)
	}

	// Hard fallback for missing or unknown title
	if out.Title == nil || *out.Title == "" || *out.Title == "<UNKNOWN>" || *out.Title == "UNKNOWN" {
		fallbackTitle := strings.Title(host)
		if fallbackTitle == "" {
			fallbackTitle = "Unknown"
		}
		out.Title = &fallbackTitle
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

func (s *scrapeProcessServiceImpl) fallbackToClaude(ctx context.Context, uuid string, doc *goquery.Document, actualUrl string, out *scrapeaudit.ScrapeResult) {
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
				Content: "Extract the details from this text. The text is from the URL: " + actualUrl + ". If the text looks like an error message (like 'Free subscription plan' or bot protection), ignore the text and infer the brand name (e.g., 'Google Flights', 'Austrian Airlines') from the URL. Never output <UNKNOWN>. If this page describes transport (flight, train, bus) with a round-trip booking, populate both the outbound and inbound legs (locations, ISO 8601 datetimes, and total duration in minutes). For one-way trips, leave all inbound_* fields null. If this page describes an activity or event (tour, concert, museum visit, etc.), populate start_at, end_at, and location (venue name + city, or full address if shown) if available. If a total/displayed price is visible on the page, populate price (numeric string with up to 2 decimals) and currency (PLN, EUR, MKD, or 'unknown' if a different/unrecognized currency). " + text,
			},
		},
		Tools: []openai.Tool{
			{
				Type: openai.ToolTypeFunction,
				Function: &openai.FunctionDefinition{
					Name:        "extract_page_details",
					Description: "Extracts the main title, description, primary image URL, and optional transport leg details.",
					Parameters:  extractionSchema(),
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
			log.Info("Claude text extraction raw output", "args", toolCall.Function.Arguments)
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

				out.OutboundDepartingLocation = extracted.OutboundDepartingLocation
				out.OutboundArrivingLocation = extracted.OutboundArrivingLocation
				out.OutboundDepartingAt = parseLegTime(extracted.OutboundDepartingAt)
				out.OutboundArrivingAt = parseLegTime(extracted.OutboundArrivingAt)
				out.InboundDepartingLocation = extracted.InboundDepartingLocation
				out.InboundArrivingLocation = extracted.InboundArrivingLocation
				out.InboundDepartingAt = parseLegTime(extracted.InboundDepartingAt)
				out.InboundArrivingAt = parseLegTime(extracted.InboundArrivingAt)

				out.StartAt = parseLegTime(extracted.StartAt)
				out.EndAt = parseLegTime(extracted.EndAt)
				out.Location = extracted.Location

				out.OutboundDurationMinutes = utility.ParseDurationMinutes(extracted.OutboundDurationMinutes)
				out.InboundDurationMinutes = utility.ParseDurationMinutes(extracted.InboundDurationMinutes)

				out.Price = extracted.Price
				out.Currency = utility.ParseCurrencyCode(extracted.Currency)
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

func (s *scrapeProcessServiceImpl) ExtractFromImage(ctx context.Context, urlStr string, imageBytes []byte) (*AIExtraction, error) {
	host := getHostFromURL(urlStr)
	if host == "" {
		host = "uploaded_image"
	}
	auditRecord, err := s.scrapeAuditService.CreateScrapeAudit(ctx, urlStr, host)

	out := scrapeaudit.ScrapeResult{
		InitialUrl: urlStr,
		ActualUrl:  urlStr,
	}

	config := openai.DefaultConfig(s.openrouterKey)
	config.BaseURL = "https://openrouter.ai/api/v1"
	client := openai.NewClientWithConfig(config)

	b64Image := base64.StdEncoding.EncodeToString(imageBytes)
	dataUrl := "data:image/jpeg;base64," + b64Image

	params := s.imageExtractionClaudeConfig(dataUrl)

	resp, err := client.CreateChatCompletion(ctx, params)
	if err != nil {
		log.Error("Failed with Haiku image extraction", "error", err)
		if auditRecord != nil {
			s.scrapeAuditService.UpdateScrapeAuditByUuid(ctx, auditRecord.Uuid, db.ScrapeStatusFailed, &out)
		}
		return nil, err
	}

	if len(resp.Choices) > 0 && len(resp.Choices[0].Message.ToolCalls) > 0 {
		toolCall := resp.Choices[0].Message.ToolCalls[0]
		if toolCall.Function.Name == "extract_page_details" {
			log.Info("Claude image extraction raw output", "args", toolCall.Function.Arguments)
			var extracted AIExtraction
			if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &extracted); err == nil {
				if auditRecord != nil {
					out.Title = &extracted.Title
					out.Description = &extracted.Description
					out.ImageUrl = &extracted.ImageURL
					s.scrapeAuditService.UpdateScrapeAuditByUuid(ctx, auditRecord.Uuid, db.ScrapeStatusCompletedByImageExtraction, &out)
				}
				return &extracted, nil
			}
		}
	}

	if auditRecord != nil {
		s.scrapeAuditService.UpdateScrapeAuditByUuid(ctx, auditRecord.Uuid, db.ScrapeStatusFailed, &out)
	}
	return nil, nil
}

func (s *scrapeProcessServiceImpl) imageExtractionClaudeConfig(dataUrl string) openai.ChatCompletionRequest {
	return openai.ChatCompletionRequest{
		Model:     "anthropic/claude-3-haiku",
		MaxTokens: 1024,
		Messages: []openai.ChatCompletionMessage{
			{
				Role: openai.ChatMessageRoleUser,
				MultiContent: []openai.ChatMessagePart{
					{
						Type: openai.ChatMessagePartTypeText,
						Text: "Extract the details from this image. If it depicts transport (flight, train, bus), populate the outbound and (if round-trip) inbound leg locations, ISO 8601 datetimes, and total duration in minutes. For one-way trips, leave all inbound_* fields null. If it depicts an activity or event (tour, concert, museum visit, etc.), populate start_at, end_at, and location (venue name + city, or full address if shown). If a price is shown, populate price (numeric string with up to 2 decimals) and currency (PLN, EUR, MKD, or 'unknown' if a different/unrecognized currency).",
					},
					{
						Type: openai.ChatMessagePartTypeImageURL,
						ImageURL: &openai.ChatMessageImageURL{
							URL: dataUrl,
						},
					},
				},
			},
		},
		Tools: []openai.Tool{
			{
				Type: openai.ToolTypeFunction,
				Function: &openai.FunctionDefinition{
					Name:        "extract_page_details",
					Description: "Extracts the main title, description, primary image URL, and optional transport leg details.",
					Parameters:  extractionSchema(),
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
}

func extractionSchema() jsonschema.Definition {
	return jsonschema.Definition{
		Type: jsonschema.Object,
		Properties: map[string]jsonschema.Definition{
			"title": {
				Type:        jsonschema.String,
				Description: "The name of the accommodation, activity, transport, or main product.",
			},
			"description": {
				Type:        jsonschema.String,
				Description: "A short 1-2 sentence description.",
			},
			"image_url": {
				Type:        jsonschema.String,
				Description: "The absolute URL to the main image.",
			},
			"outbound_departing_location": {
				Type:        jsonschema.String,
				Description: "Origin (city or airport code) for the outbound/first leg of a transport booking. Null if not shown or not a transport page.",
			},
			"outbound_arriving_location": {
				Type:        jsonschema.String,
				Description: "Destination for the outbound/first leg of a transport booking. Null if not shown or not a transport page.",
			},
			"outbound_departing_at": {
				Type:        jsonschema.String,
				Description: "Outbound departure date and time as shown on the ticket/page. Use ISO 8601 local format YYYY-MM-DDTHH:MM:SS with NO timezone offset — just the wall-clock time at the departure location. Null if not shown.",
			},
			"outbound_arriving_at": {
				Type:        jsonschema.String,
				Description: "Outbound arrival date and time as shown on the ticket/page. Use ISO 8601 local format YYYY-MM-DDTHH:MM:SS with NO timezone offset — wall-clock time at the arrival location. Null if not shown.",
			},
			"inbound_departing_location": {
				Type:        jsonschema.String,
				Description: "Origin for the inbound/return leg. Null if one-way or not shown.",
			},
			"inbound_arriving_location": {
				Type:        jsonschema.String,
				Description: "Destination for the inbound/return leg. Null if one-way or not shown.",
			},
			"inbound_departing_at": {
				Type:        jsonschema.String,
				Description: "Return/inbound departure date and time (YYYY-MM-DDTHH:MM:SS, no timezone offset — wall-clock at the departure location). Null if one-way or not shown.",
			},
			"inbound_arriving_at": {
				Type:        jsonschema.String,
				Description: "Return/inbound arrival date and time (YYYY-MM-DDTHH:MM:SS, no timezone offset — wall-clock at the arrival location). Null if one-way or not shown.",
			},
			"start_at": {
				Type:        jsonschema.String,
				Description: "Activity/event start date and time as shown on the ticket/page. Use ISO 8601 local format YYYY-MM-DDTHH:MM:SS with NO timezone offset — wall-clock time at the venue. Null if not an activity or not shown.",
			},
			"end_at": {
				Type:        jsonschema.String,
				Description: "Activity/event end date and time (YYYY-MM-DDTHH:MM:SS, no timezone offset — wall-clock at the venue). Null if not an activity or not shown.",
			},
			"location": {
				Type:        jsonschema.String,
				Description: "Activity/event venue or location — venue name with city (e.g. 'Stadion Narodowy, Warsaw') or a full street address if shown. Null if not an activity or not shown.",
			},
			"outbound_duration_minutes": {
				Type:        jsonschema.String,
				Description: "Total duration of the outbound transport leg expressed as integer minutes in a numeric string (e.g. '225' for 3h 45m). Includes layovers if multi-stop. Null if not a transport page or not shown.",
			},
			"inbound_duration_minutes": {
				Type:        jsonschema.String,
				Description: "Total duration of the inbound/return transport leg as integer minutes in a numeric string (e.g. '225'). Null if one-way or not shown.",
			},
			"price": {
				Type:        jsonschema.String,
				Description: "Numeric price as a string with up to 2 decimals (e.g. '1234.56'). Use the total/displayed price for the booking. Null if no price is shown.",
			},
			"currency": {
				Type:        jsonschema.String,
				Description: "ISO 4217 currency code if recognizable: 'PLN', 'EUR', or 'MKD'. If a price is shown but the currency is unclear or different (USD, GBP, etc.), return 'unknown'. Null if no price.",
			},
		},
		Required: []string{"title", "description", "image_url"},
	}
}

func parseLegTime(s *string) *time.Time {
	if s == nil || *s == "" {
		return nil
	}
	layouts := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05",
		"2006-01-02T15:04",
		"2006-01-02 15:04:05Z07:00",
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, *s); err == nil {
			// Wall-clock: discard offset, re-anchor components as UTC so
			// the displayed time matches what's printed on the ticket.
			wall := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, time.UTC)
			return &wall
		}
	}
	log.Warn("Failed parsing leg time from LLM extraction", "value", *s)
	return nil
}
