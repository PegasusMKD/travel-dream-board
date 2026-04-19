package scrapeaudit_test

import (
	"context"
	"testing"

	"github.com/PegasusMKD/travel-dream-board/internal/db"
	"github.com/PegasusMKD/travel-dream-board/internal/scrape_audit"
	mock_scrape_audit "github.com/PegasusMKD/travel-dream-board/internal/scrape_audit/mocks"
	"github.com/stretchr/testify/assert"
)

func TestScrapeAuditService_CreateScrapeAudit(t *testing.T) {
	mockRepo := new(mock_scrape_audit.Repository)
	svc := scrapeaudit.NewService(mockRepo)

	ctx := context.Background()
	url := "http://example.com"
	host := "example.com"
	expectedAudit := &scrapeaudit.ScrapeAudit{Uuid: "audit-1", Url: url, Host: host}

	mockRepo.On("CreateScrapeAudit", ctx, url, host).Return(expectedAudit, nil)

	result, err := svc.CreateScrapeAudit(ctx, url, host)
	assert.NoError(t, err)
	assert.Equal(t, expectedAudit, result)
	mockRepo.AssertExpectations(t)
}

func TestScrapeAuditService_UpdateScrapeAuditByUuid(t *testing.T) {
	mockRepo := new(mock_scrape_audit.Repository)
	svc := scrapeaudit.NewService(mockRepo)

	ctx := context.Background()
	uuid := "audit-1"
	status := db.ScrapeStatusCompletedByAi
	title := "Test Title"
	data := &scrapeaudit.ScrapeResult{Title: &title}

	mockRepo.On("UpdateScrapeAuditByUuid", ctx, uuid, status, data).Return(nil)

	err := svc.UpdateScrapeAuditByUuid(ctx, uuid, status, data)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
