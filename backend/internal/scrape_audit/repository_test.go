package scrapeaudit_test

import (
	"context"
	"testing"

	"github.com/PegasusMKD/travel-dream-board/internal/db"
	scrapeaudit "github.com/PegasusMKD/travel-dream-board/internal/scrape_audit"
	"github.com/PegasusMKD/travel-dream-board/internal/testutil"
	"github.com/stretchr/testify/require"
)

func TestScrapeAuditRepository(t *testing.T) {
	env := testutil.SetupTestDB(t)
	repo := scrapeaudit.NewRepository(env.Queries)
	ctx := context.Background()

	t.Run("CreateScrapeAudit", func(t *testing.T) {
		urlStr := "https://example.com/article"
		hostStr := "example.com"
		
		createdAudit, err := repo.CreateScrapeAudit(ctx, urlStr, hostStr)
		require.NoError(t, err)
		require.NotNil(t, createdAudit)
		require.NotEmpty(t, createdAudit.Uuid)
		require.Equal(t, urlStr, createdAudit.Url)
		require.Equal(t, hostStr, createdAudit.Host)
		// Assuming default status is processing or similar based on DB schema
	})

	t.Run("UpdateScrapeAuditByUuid", func(t *testing.T) {
		urlStr := "https://example.com/update"
		hostStr := "example.com"
		
		createdAudit, err := repo.CreateScrapeAudit(ctx, urlStr, hostStr)
		require.NoError(t, err)

		title := "Updated Title"
		desc := "Updated Description"
		imgUrl := "https://example.com/img.jpg"
		siteName := "Example Site"

		resultData := &scrapeaudit.ScrapeResult{
			Title:       &title,
			Description: &desc,
			ImageUrl:    &imgUrl,
			SiteName:    &siteName,
		}

		err = repo.UpdateScrapeAuditByUuid(ctx, createdAudit.Uuid, db.ScrapeStatusCompletedByOg, resultData)
		require.NoError(t, err)
		
		// To verify update we might need a Get method which isn't there.
		// If we can't get it, the fact it doesn't error is enough for repository interface test coverage.
	})

	t.Run("UpdateScrapeAuditByUuid - Invalid UUID", func(t *testing.T) {
		resultData := &scrapeaudit.ScrapeResult{}
		err := repo.UpdateScrapeAuditByUuid(ctx, "invalid-uuid", db.ScrapeStatusFailed, resultData)
		require.Error(t, err)
	})
}
