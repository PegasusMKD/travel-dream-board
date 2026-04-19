package boards_test

import (
	"context"
	"testing"
	"time"

	"github.com/PegasusMKD/travel-dream-board/internal/boards"
	"github.com/PegasusMKD/travel-dream-board/internal/db"
	"github.com/PegasusMKD/travel-dream-board/internal/testutil"
	"github.com/stretchr/testify/require"
)

func TestBoardsRepository(t *testing.T) {
	env := testutil.SetupTestDB(t)
	repo := boards.NewRepository(env.Queries)
	ctx := context.Background()

	email := "test@example.com"
	user, err := env.Queries.UpsertUser(ctx, db.UpsertUserParams{
		Email: &email,
		Name:  "Test User",
	})
	require.NoError(t, err)

	// Ensure string is correctly extracted from UUID
	userUuidStr := user.Uuid.String()
	// userUuidStr will likely be something like "[uuid string format]" or we need to extract string
	userUuidStrBytes, _ := user.Uuid.Value()
	userUuidStr = userUuidStrBytes.(string)

	t.Run("CreateBoard", func(t *testing.T) {
		boardData := &boards.Board{
			Name:         "Test Trip",
			LocationName: "Paris",
			UserUuid:     &userUuidStr,
		}

		createdBoard, err := repo.CreateBoard(ctx, boardData)
		require.NoError(t, err)
		require.NotNil(t, createdBoard)
		require.NotEmpty(t, createdBoard.Uuid)
		require.Equal(t, "Test Trip", createdBoard.Name)
		require.Equal(t, "Paris", createdBoard.LocationName)
		require.Equal(t, userUuidStr, *createdBoard.UserUuid)
	})

	t.Run("CreateBoard - Invalid User UUID", func(t *testing.T) {
		invalidUuid := "invalid-uuid"
		boardData := &boards.Board{
			Name:         "Invalid User Trip",
			LocationName: "Paris",
			UserUuid:     &invalidUuid,
		}

		_, err := repo.CreateBoard(ctx, boardData)
		require.Error(t, err)
	})

	t.Run("GetBoardById", func(t *testing.T) {
		boardData := &boards.Board{
			Name:         "Rome Trip",
			LocationName: "Rome",
			UserUuid:     &userUuidStr,
		}
		createdBoard, err := repo.CreateBoard(ctx, boardData)
		require.NoError(t, err)

		fetchedBoard, err := repo.GetBoardById(ctx, createdBoard.Uuid)
		require.NoError(t, err)
		require.Equal(t, createdBoard.Uuid, fetchedBoard.Uuid)
		require.Equal(t, "Rome Trip", fetchedBoard.Name)
	})

	t.Run("GetBoardById - Invalid UUID", func(t *testing.T) {
		_, err := repo.GetBoardById(ctx, "invalid-uuid")
		require.Error(t, err)
	})

	t.Run("GetAllBoards", func(t *testing.T) {
		boardData := &boards.Board{
			Name:         "London Trip",
			LocationName: "London",
			UserUuid:     &userUuidStr,
		}
		_, err := repo.CreateBoard(ctx, boardData)
		require.NoError(t, err)

		allBoards, err := repo.GetAllBoards(ctx, userUuidStr)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(allBoards), 1)

		var found bool
		for _, b := range allBoards {
			if b.Name == "London Trip" {
				found = true
				break
				t.Run("DeleteBoardById - Invalid UUID", func(t *testing.T) {
					err := repo.DeleteBoardById(ctx, "invalid-uuid")
					require.Error(t, err)
				})
			}
		}
		require.True(t, found, "Should find the newly created London Trip")
	})

	t.Run("GetAllBoards - Invalid UUID", func(t *testing.T) {
		_, err := repo.GetAllBoards(ctx, "invalid-uuid")
		require.Error(t, err)
	})

	t.Run("UpdateBoardById", func(t *testing.T) {
		boardData := &boards.Board{
			Name:         "Madrid Trip",
			LocationName: "Madrid",
			UserUuid:     &userUuidStr,
		}
		createdBoard, err := repo.CreateBoard(ctx, boardData)
		require.NoError(t, err)

		details := "Detailed info about Madrid"
		createdBoard.Details = &details
		createdBoard.Name = "Madrid Updated Trip"

		startsAt := time.Now().Add(24 * time.Hour).Truncate(24 * time.Hour).UTC() // since it's a date usually
		createdBoard.StartsAt = &startsAt

		err = repo.UpdateBoardById(ctx, createdBoard.Uuid, createdBoard)
		require.NoError(t, err)

		updatedBoard, err := repo.GetBoardById(ctx, createdBoard.Uuid)
		require.NoError(t, err)
		require.Equal(t, "Madrid Updated Trip", updatedBoard.Name)
		require.Equal(t, details, *updatedBoard.Details)
		// Date matching might be tricky due to pgtype.Date mapping but let's check year/month/day
		require.Equal(t, startsAt.Year(), updatedBoard.StartsAt.Year())
		require.Equal(t, startsAt.Month(), updatedBoard.StartsAt.Month())
		require.Equal(t, startsAt.Day(), updatedBoard.StartsAt.Day())
	})

	t.Run("UpdateBoardById - Invalid UUID", func(t *testing.T) {
		err := repo.UpdateBoardById(ctx, "invalid-uuid", &boards.Board{})
		require.Error(t, err)
	})

	t.Run("DeleteBoardById", func(t *testing.T) {
		boardData := &boards.Board{
			Name:         "Delete Trip",
			LocationName: "Delete City",
			UserUuid:     &userUuidStr,
		}
		createdBoard, err := repo.CreateBoard(ctx, boardData)
		require.NoError(t, err)

		err = repo.DeleteBoardById(ctx, createdBoard.Uuid)
		require.NoError(t, err)

		// Try to fetch it again, should fail
		_, err = repo.GetBoardById(ctx, createdBoard.Uuid)
		require.Error(t, err)
	})
}
