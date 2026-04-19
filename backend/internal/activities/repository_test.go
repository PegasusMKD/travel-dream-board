package activities_test

import (
	"context"
	"testing"

	"github.com/PegasusMKD/travel-dream-board/internal/activities"
	"github.com/PegasusMKD/travel-dream-board/internal/boards"
	"github.com/PegasusMKD/travel-dream-board/internal/db"
	"github.com/PegasusMKD/travel-dream-board/internal/testutil"
	"github.com/stretchr/testify/require"
)

func TestActivitiesRepository(t *testing.T) {
	env := testutil.SetupTestDB(t)
	repo := activities.NewRepository(env.Queries)
	ctx := context.Background()

	// 1. Create a user
	email := "activityuser@example.com"
	user, err := env.Queries.UpsertUser(ctx, db.UpsertUserParams{
		Email: &email,
		Name:  "Test Activity User",
	})
	require.NoError(t, err)

	userUuidStrBytes, _ := user.Uuid.Value()
	userUuidStr := userUuidStrBytes.(string)

	// 2. Create a board
	boardRepo := boards.NewRepository(env.Queries)
	board, err := boardRepo.CreateBoard(ctx, &boards.Board{
		Name:         "Activity Board",
		LocationName: "Activity Location",
		UserUuid:     &userUuidStr,
	})
	require.NoError(t, err)

	t.Run("CreateActivity", func(t *testing.T) {
		imgUrl := "http://example.com/image.jpg"
		actData := &activities.Activity{
			Url:       "http://example.com",
			Title:     "Visit Museum",
			BoardUuid: board.Uuid,
			ImageUrl:  &imgUrl,
			UserUuid:  userUuidStr,
		}

		createdAct, err := repo.CreateActivity(ctx, actData)
		require.NoError(t, err)
		require.NotNil(t, createdAct)
		require.NotEmpty(t, createdAct.Uuid)
		require.Equal(t, "Visit Museum", createdAct.Title)
		require.Equal(t, board.Uuid, createdAct.BoardUuid)
		require.Equal(t, userUuidStr, createdAct.UserUuid)
		require.Equal(t, imgUrl, *createdAct.ImageUrl)
	})

	t.Run("CreateActivity - Invalid Board UUID", func(t *testing.T) {
		actData := &activities.Activity{
			Url:       "http://example.com",
			Title:     "Visit Museum",
			BoardUuid: "invalid-uuid",
			UserUuid:  userUuidStr,
		}
		_, err := repo.CreateActivity(ctx, actData)
		require.Error(t, err)
	})

	t.Run("CreateActivity - Invalid User UUID", func(t *testing.T) {
		actData := &activities.Activity{
			Url:       "http://example.com",
			Title:     "Visit Museum",
			BoardUuid: board.Uuid,
			UserUuid:  "invalid-uuid",
		}
		_, err := repo.CreateActivity(ctx, actData)
		require.Error(t, err)
	})

	t.Run("GetActivityById", func(t *testing.T) {
		actData := &activities.Activity{
			Url:       "http://example.com/2",
			Title:     "Eat Pizza",
			BoardUuid: board.Uuid,
			UserUuid:  userUuidStr,
		}
		createdAct, err := repo.CreateActivity(ctx, actData)
		require.NoError(t, err)

		fetchedAct, err := repo.GetActivityById(ctx, createdAct.Uuid)
		require.NoError(t, err)
		require.Equal(t, createdAct.Uuid, fetchedAct.Uuid)
		require.Equal(t, "Eat Pizza", fetchedAct.Title)
	})

	t.Run("GetActivityById - Invalid UUID", func(t *testing.T) {
		_, err := repo.GetActivityById(ctx, "invalid-uuid")
		require.Error(t, err)
	})

	t.Run("GetAllActivitiesByBoardId", func(t *testing.T) {
		actData := &activities.Activity{
			Url:       "http://example.com/3",
			Title:     "Go Hiking",
			BoardUuid: board.Uuid,
			UserUuid:  userUuidStr,
		}
		_, err := repo.CreateActivity(ctx, actData)
		require.NoError(t, err)

		allActs, err := repo.GetAllActivitiesByBoardId(ctx, board.Uuid)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(allActs), 1)

		var found bool
		for _, a := range allActs {
			if a.Title == "Go Hiking" {
				found = true
				break
				t.Run("DeleteActivityById - Invalid UUID", func(t *testing.T) {
					err := repo.DeleteActivityById(ctx, "invalid-uuid")
					require.Error(t, err)
				})
			}
		}
		require.True(t, found, "Should find the newly created activity")
	})

	t.Run("GetAllActivitiesByBoardId - Invalid UUID", func(t *testing.T) {
		_, err := repo.GetAllActivitiesByBoardId(ctx, "invalid-uuid")
		require.Error(t, err)
	})

	t.Run("UpdateActivityById", func(t *testing.T) {
		actData := &activities.Activity{
			Url:       "http://example.com/4",
			Title:     "Go Surfing",
			BoardUuid: board.Uuid,
			UserUuid:  userUuidStr,
		}
		createdAct, err := repo.CreateActivity(ctx, actData)
		require.NoError(t, err)

		notes := "Bring surfboard"
		createdAct.Notes = &notes
		createdAct.Title = "Go Surfing Updated"
		createdAct.Selected = true

		err = repo.UpdateActivityById(ctx, createdAct.Uuid, createdAct)
		require.NoError(t, err)

		updatedAct, err := repo.GetActivityById(ctx, createdAct.Uuid)
		require.NoError(t, err)
		require.Equal(t, "Go Surfing Updated", updatedAct.Title)
		require.Equal(t, notes, *updatedAct.Notes)
		require.True(t, updatedAct.Selected)
	})

	t.Run("UpdateActivityById - Invalid UUID", func(t *testing.T) {
		err := repo.UpdateActivityById(ctx, "invalid-uuid", &activities.Activity{})
		require.Error(t, err)
	})

	t.Run("DeleteActivityById", func(t *testing.T) {
		actData := &activities.Activity{
			Url:       "http://example.com/5",
			Title:     "Delete Me",
			BoardUuid: board.Uuid,
			UserUuid:  userUuidStr,
		}
		createdAct, err := repo.CreateActivity(ctx, actData)
		require.NoError(t, err)

		err = repo.DeleteActivityById(ctx, createdAct.Uuid)
		require.NoError(t, err)

		// Try to fetch it again, should fail
		_, err = repo.GetActivityById(ctx, createdAct.Uuid)
		require.Error(t, err)
	})
}
