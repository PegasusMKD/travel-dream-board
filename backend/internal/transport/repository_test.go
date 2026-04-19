package transport_test

import (
	"context"
	"testing"

	"github.com/PegasusMKD/travel-dream-board/internal/boards"
	"github.com/PegasusMKD/travel-dream-board/internal/db"
	"github.com/PegasusMKD/travel-dream-board/internal/testutil"
	"github.com/PegasusMKD/travel-dream-board/internal/transport"
	"github.com/stretchr/testify/require"
)

func TestTransportRepository(t *testing.T) {
	env := testutil.SetupTestDB(t)
	repo := transport.NewRepository(env.Queries)
	ctx := context.Background()

	// 1. Create a user
	email := "transuser@example.com"
	user, err := env.Queries.UpsertUser(ctx, db.UpsertUserParams{
		Email: &email,
		Name:  "Test Transport User",
	})
	require.NoError(t, err)

	userUuidStrBytes, _ := user.Uuid.Value()
	userUuidStr := userUuidStrBytes.(string)

	// 2. Create a board
	boardRepo := boards.NewRepository(env.Queries)
	board, err := boardRepo.CreateBoard(ctx, &boards.Board{
		Name:         "Transport Board",
		LocationName: "Transport Location",
		UserUuid:     &userUuidStr,
	})
	require.NoError(t, err)

	t.Run("CreateTransport", func(t *testing.T) {
		imgUrl := "http://example.com/bus.jpg"
		transData := &transport.Transport{
			Url:       "http://example.com/bus",
			Title:     "Book Bus",
			BoardUuid: board.Uuid,
			ImageUrl:  &imgUrl,
			UserUuid:  userUuidStr,
		}

		createdTrans, err := repo.CreateTransport(ctx, transData)
		require.NoError(t, err)
		require.NotNil(t, createdTrans)
		require.NotEmpty(t, createdTrans.Uuid)
		require.Equal(t, "Book Bus", createdTrans.Title)
		require.Equal(t, board.Uuid, createdTrans.BoardUuid)
		require.Equal(t, userUuidStr, createdTrans.UserUuid)
		require.Equal(t, imgUrl, *createdTrans.ImageUrl)
	})

	t.Run("CreateTransport - Invalid Board UUID", func(t *testing.T) {
		transData := &transport.Transport{
			Url:       "http://example.com/bus",
			Title:     "Book Bus",
			BoardUuid: "invalid-uuid",
			UserUuid:  userUuidStr,
		}
		_, err := repo.CreateTransport(ctx, transData)
		require.Error(t, err)
	})

	t.Run("CreateTransport - Invalid User UUID", func(t *testing.T) {
		transData := &transport.Transport{
			Url:       "http://example.com/bus",
			Title:     "Book Bus",
			BoardUuid: board.Uuid,
			UserUuid:  "invalid-uuid",
		}
		_, err := repo.CreateTransport(ctx, transData)
		require.Error(t, err)
	})

	t.Run("GetTransportById", func(t *testing.T) {
		transData := &transport.Transport{
			Url:       "http://example.com/train",
			Title:     "Book Train",
			BoardUuid: board.Uuid,
			UserUuid:  userUuidStr,
		}
		createdTrans, err := repo.CreateTransport(ctx, transData)
		require.NoError(t, err)

		fetchedTrans, err := repo.GetTransportById(ctx, createdTrans.Uuid)
		require.NoError(t, err)
		require.Equal(t, createdTrans.Uuid, fetchedTrans.Uuid)
		require.Equal(t, "Book Train", fetchedTrans.Title)
	})

	t.Run("GetTransportById - Invalid UUID", func(t *testing.T) {
		_, err := repo.GetTransportById(ctx, "invalid-uuid")
		require.Error(t, err)
	})

	t.Run("GetAllTransportsByBoardId", func(t *testing.T) {
		transData := &transport.Transport{
			Url:       "http://example.com/flight",
			Title:     "Book Flight",
			BoardUuid: board.Uuid,
			UserUuid:  userUuidStr,
		}
		_, err := repo.CreateTransport(ctx, transData)
		require.NoError(t, err)

		allTrans, err := repo.GetAllTransportsByBoardId(ctx, board.Uuid)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(allTrans), 1)

		var found bool
		for _, t := range allTrans {
			if t.Title == "Book Flight" {
				found = true
				break
			}
		}
		require.True(t, found, "Should find the newly created transport")
	})

	t.Run("DeleteTransportById - Invalid UUID", func(t *testing.T) {
		err := repo.DeleteTransportById(ctx, "invalid-uuid")
		require.Error(t, err)
	})

	t.Run("GetAllTransportsByBoardId - Invalid UUID", func(t *testing.T) {
		_, err := repo.GetAllTransportsByBoardId(ctx, "invalid-uuid")
		require.Error(t, err)
	})

	t.Run("UpdateTransportById", func(t *testing.T) {
		transData := &transport.Transport{
			Url:       "http://example.com/car",
			Title:     "Rent Car",
			BoardUuid: board.Uuid,
			UserUuid:  userUuidStr,
		}
		createdTrans, err := repo.CreateTransport(ctx, transData)
		require.NoError(t, err)

		notes := "Need automatic"
		createdTrans.Notes = &notes
		createdTrans.Title = "Rent Car Updated"
		createdTrans.Selected = true

		err = repo.UpdateTransportById(ctx, createdTrans.Uuid, createdTrans)
		require.NoError(t, err)

		updatedTrans, err := repo.GetTransportById(ctx, createdTrans.Uuid)
		require.NoError(t, err)
		require.Equal(t, "Rent Car Updated", updatedTrans.Title)
		require.Equal(t, notes, *updatedTrans.Notes)
		require.True(t, updatedTrans.Selected)
	})

	t.Run("UpdateTransportById - Invalid UUID", func(t *testing.T) {
		err := repo.UpdateTransportById(ctx, "invalid-uuid", &transport.Transport{})
		require.Error(t, err)
	})

	t.Run("DeleteTransportById", func(t *testing.T) {
		transData := &transport.Transport{
			Url:       "http://example.com/delete",
			Title:     "Delete Me",
			BoardUuid: board.Uuid,
			UserUuid:  userUuidStr,
		}
		createdTrans, err := repo.CreateTransport(ctx, transData)
		require.NoError(t, err)

		err = repo.DeleteTransportById(ctx, createdTrans.Uuid)
		require.NoError(t, err)

		// Try to fetch it again, should fail
		_, err = repo.GetTransportById(ctx, createdTrans.Uuid)
		require.Error(t, err)
	})
}
