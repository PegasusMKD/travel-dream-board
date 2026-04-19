package sharetokens_test

import (
	"context"
	"testing"

	"github.com/PegasusMKD/travel-dream-board/internal/boards"
	"github.com/PegasusMKD/travel-dream-board/internal/db"
	"github.com/PegasusMKD/travel-dream-board/internal/sharetokens"
	"github.com/PegasusMKD/travel-dream-board/internal/testutil"
	"github.com/stretchr/testify/require"
)

func TestShareTokensRepository(t *testing.T) {
	env := testutil.SetupTestDB(t)
	repo := sharetokens.NewRepository(env.Queries)
	ctx := context.Background()

	// 1. Create a user
	email := "shareuser@example.com"
	user, err := env.Queries.UpsertUser(ctx, db.UpsertUserParams{
		Email: &email,
		Name:  "Test Share User",
	})
	require.NoError(t, err)
	
	userUuidStrBytes, _ := user.Uuid.Value()
	userUuidStr := userUuidStrBytes.(string)

	// 2. Create a board
	boardRepo := boards.NewRepository(env.Queries)
	board, err := boardRepo.CreateBoard(ctx, &boards.Board{
		Name:         "Share Board",
		LocationName: "Share Location",
		UserUuid:     &userUuidStr,
	})
	require.NoError(t, err)

	t.Run("CreateShareToken", func(t *testing.T) {
		tokenStr := "some-random-token-xyz"
		
		createdToken, err := repo.CreateShareToken(ctx, tokenStr, board.Uuid)
		require.NoError(t, err)
		require.NotNil(t, createdToken)
		require.Equal(t, tokenStr, createdToken.Token)
		require.Equal(t, board.Uuid, createdToken.BoardUuid)
	})

	t.Run("CreateShareToken - Invalid Board UUID", func(t *testing.T) {
		_, err := repo.CreateShareToken(ctx, "another-token", "invalid-uuid")
		require.Error(t, err)
	})

	t.Run("GetShareToken", func(t *testing.T) {
		tokenStr := "token-for-get"
		_, err := repo.CreateShareToken(ctx, tokenStr, board.Uuid)
		require.NoError(t, err)

		fetchedToken, err := repo.GetShareToken(ctx, tokenStr)
		require.NoError(t, err)
		require.Equal(t, tokenStr, fetchedToken.Token)
		require.Equal(t, board.Uuid, fetchedToken.BoardUuid)
	})

	t.Run("GetShareTokensForBoard", func(t *testing.T) {
		tokenStr := "token-for-board-list"
		_, err := repo.CreateShareToken(ctx, tokenStr, board.Uuid)
		require.NoError(t, err)

		allTokens, err := repo.GetShareTokensForBoard(ctx, board.Uuid)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(allTokens), 1)
		
		var found bool
		for _, tok := range allTokens {
			if tok.Token == tokenStr {
				found = true
				break
			}
		}
		require.True(t, found, "Should find the newly created token")
	})

	t.Run("GetShareTokensForBoard - Invalid UUID", func(t *testing.T) {
		_, err := repo.GetShareTokensForBoard(ctx, "invalid-uuid")
		require.Error(t, err)
	})

	t.Run("DeleteShareToken", func(t *testing.T) {
		tokenStr := "token-for-delete"
		_, err := repo.CreateShareToken(ctx, tokenStr, board.Uuid)
		require.NoError(t, err)

		err = repo.DeleteShareToken(ctx, tokenStr, board.Uuid)
		require.NoError(t, err)

		// Try to fetch it again, should fail
		_, err = repo.GetShareToken(ctx, tokenStr)
		require.Error(t, err)
	})

	t.Run("DeleteShareToken - Invalid UUID", func(t *testing.T) {
		err := repo.DeleteShareToken(ctx, "token-for-delete", "invalid-uuid")
		require.Error(t, err)
	})
}
