package accomodations_test

import (
	"context"
	"testing"

	"github.com/PegasusMKD/travel-dream-board/internal/accomodations"
	"github.com/PegasusMKD/travel-dream-board/internal/boards"
	"github.com/PegasusMKD/travel-dream-board/internal/db"
	"github.com/PegasusMKD/travel-dream-board/internal/testutil"
	"github.com/stretchr/testify/require"
)

func TestAccomodationsRepository(t *testing.T) {
	env := testutil.SetupTestDB(t)
	repo := accomodations.NewRepository(env.Queries)
	ctx := context.Background()

	// 1. Create a user
	email := "accuser@example.com"
	user, err := env.Queries.UpsertUser(ctx, db.UpsertUserParams{
		Email: &email,
		Name:  "Test Acc User",
	})
	require.NoError(t, err)

	userUuidStrBytes, _ := user.Uuid.Value()
	userUuidStr := userUuidStrBytes.(string)

	// 2. Create a board
	boardRepo := boards.NewRepository(env.Queries)
	board, err := boardRepo.CreateBoard(ctx, &boards.Board{
		Name:         "Acc Board",
		LocationName: "Acc Location",
		UserUuid:     &userUuidStr,
	})
	require.NoError(t, err)

	t.Run("CreateAccomodation", func(t *testing.T) {
		imgUrl := "http://example.com/hotel.jpg"
		accData := &accomodations.Accomodation{
			Url:       "http://example.com/hotel",
			Title:     "Book Hotel",
			BoardUuid: board.Uuid,
			ImageUrl:  &imgUrl,
			UserUuid:  userUuidStr,
		}

		createdAcc, err := repo.CreateAccomodation(ctx, accData)
		require.NoError(t, err)
		require.NotNil(t, createdAcc)
		require.NotEmpty(t, createdAcc.Uuid)
		require.Equal(t, "Book Hotel", createdAcc.Title)
		require.Equal(t, board.Uuid, createdAcc.BoardUuid)
		require.Equal(t, userUuidStr, createdAcc.UserUuid)
		require.Equal(t, imgUrl, *createdAcc.ImageUrl)
	})

	t.Run("CreateAccomodation - Invalid Board UUID", func(t *testing.T) {
		accData := &accomodations.Accomodation{
			Url:       "http://example.com/hotel",
			Title:     "Book Hotel",
			BoardUuid: "invalid-uuid",
			UserUuid:  userUuidStr,
		}
		_, err := repo.CreateAccomodation(ctx, accData)
		require.Error(t, err)
	})

	t.Run("CreateAccomodation - Invalid User UUID", func(t *testing.T) {
		accData := &accomodations.Accomodation{
			Url:       "http://example.com/hotel",
			Title:     "Book Hotel",
			BoardUuid: board.Uuid,
			UserUuid:  "invalid-uuid",
		}
		_, err := repo.CreateAccomodation(ctx, accData)
		require.Error(t, err)
	})

	t.Run("GetAccomodationById", func(t *testing.T) {
		accData := &accomodations.Accomodation{
			Url:       "http://example.com/hostel",
			Title:     "Book Hostel",
			BoardUuid: board.Uuid,
			UserUuid:  userUuidStr,
		}
		createdAcc, err := repo.CreateAccomodation(ctx, accData)
		require.NoError(t, err)

		fetchedAcc, err := repo.GetAccomodationById(ctx, createdAcc.Uuid)
		require.NoError(t, err)
		require.Equal(t, createdAcc.Uuid, fetchedAcc.Uuid)
		require.Equal(t, "Book Hostel", fetchedAcc.Title)
	})

	t.Run("GetAccomodationById - Invalid UUID", func(t *testing.T) {
		_, err := repo.GetAccomodationById(ctx, "invalid-uuid")
		require.Error(t, err)
	})

	t.Run("GetAllAccomodationsByBoardId", func(t *testing.T) {
		accData := &accomodations.Accomodation{
			Url:       "http://example.com/airbnb",
			Title:     "Book Airbnb",
			BoardUuid: board.Uuid,
			UserUuid:  userUuidStr,
		}
		_, err := repo.CreateAccomodation(ctx, accData)
		require.NoError(t, err)

		allAccs, err := repo.GetAllAccomodationsByBoardId(ctx, board.Uuid)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(allAccs), 1)

		var found bool
		for _, a := range allAccs {
			if a.Title == "Book Airbnb" {
				found = true
				break
				t.Run("DeleteAccomodationById - Invalid UUID", func(t *testing.T) {
					err := repo.DeleteAccomodationById(ctx, "invalid-uuid")
					require.Error(t, err)
				})
			}
		}
		require.True(t, found, "Should find the newly created accomodation")
	})

	t.Run("GetAllAccomodationsByBoardId - Invalid UUID", func(t *testing.T) {
		_, err := repo.GetAllAccomodationsByBoardId(ctx, "invalid-uuid")
		require.Error(t, err)
	})

	t.Run("UpdateAccomodationById", func(t *testing.T) {
		accData := &accomodations.Accomodation{
			Url:       "http://example.com/tent",
			Title:     "Book Tent",
			BoardUuid: board.Uuid,
			UserUuid:  userUuidStr,
		}
		createdAcc, err := repo.CreateAccomodation(ctx, accData)
		require.NoError(t, err)

		notes := "Bring sleeping bag"
		createdAcc.Notes = &notes
		createdAcc.Title = "Book Tent Updated"
		createdAcc.Selected = true

		err = repo.UpdateAccomodationById(ctx, createdAcc.Uuid, createdAcc)
		require.NoError(t, err)

		updatedAcc, err := repo.GetAccomodationById(ctx, createdAcc.Uuid)
		require.NoError(t, err)
		require.Equal(t, "Book Tent Updated", updatedAcc.Title)
		require.Equal(t, notes, *updatedAcc.Notes)
		require.True(t, updatedAcc.Selected)
	})

	t.Run("UpdateAccomodationById - Invalid UUID", func(t *testing.T) {
		err := repo.UpdateAccomodationById(ctx, "invalid-uuid", &accomodations.Accomodation{})
		require.Error(t, err)
	})

	t.Run("DeleteAccomodationById", func(t *testing.T) {
		accData := &accomodations.Accomodation{
			Url:       "http://example.com/delete",
			Title:     "Delete Me",
			BoardUuid: board.Uuid,
			UserUuid:  userUuidStr,
		}
		createdAcc, err := repo.CreateAccomodation(ctx, accData)
		require.NoError(t, err)

		err = repo.DeleteAccomodationById(ctx, createdAcc.Uuid)
		require.NoError(t, err)

		// Try to fetch it again, should fail
		_, err = repo.GetAccomodationById(ctx, createdAcc.Uuid)
		require.Error(t, err)
	})
}
