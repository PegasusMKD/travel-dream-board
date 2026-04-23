package comments_test

import (
	"context"
	"testing"

	"github.com/PegasusMKD/travel-dream-board/internal/accomodations"
	"github.com/PegasusMKD/travel-dream-board/internal/boards"
	"github.com/PegasusMKD/travel-dream-board/internal/comments"
	"github.com/PegasusMKD/travel-dream-board/internal/db"
	"github.com/PegasusMKD/travel-dream-board/internal/testutil"
	"github.com/stretchr/testify/require"
)

func TestCommentsRepository(t *testing.T) {
	env := testutil.SetupTestDB(t)
	repo := comments.NewRepository(env.Queries)
	ctx := context.Background()

	// 1. Create a user
	email := "commentuser@example.com"
	user, err := env.Queries.UpsertUser(ctx, db.UpsertUserParams{
		Email: &email,
		Name:  "Test Comment User",
	})
	require.NoError(t, err)

	userUuidStrBytes, _ := user.Uuid.Value()
	userUuidStr := userUuidStrBytes.(string)

	// 2. Create a board
	boardRepo := boards.NewRepository(env.Queries)
	board, err := boardRepo.CreateBoard(ctx, &boards.Board{
		Name:         "Comment Board",
		LocationName: "Comment Location",
		UserUuid:     &userUuidStr,
	})
	require.NoError(t, err)

	// 3. Create an accomodation
	accRepo := accomodations.NewRepository(env.Queries)
	acc, err := accRepo.CreateAccomodation(ctx, &accomodations.Accomodation{
		Url:       "http://example.com/hostel",
		Title:     "Hostel",
		BoardUuid: board.Uuid,
		UserUuid:  userUuidStr,
	})
	require.NoError(t, err)

	t.Run("CreateComment", func(t *testing.T) {
		commentData := &comments.Comment{
			Content:         "Great hostel!",
			CommentedOn:     db.CommentedOnAccomodation,
			CommentedOnUuid: acc.Uuid,
			UserUuid:        userUuidStr,
		}

		createdComment, err := repo.CreateComment(ctx, commentData)
		require.NoError(t, err)
		require.NotNil(t, createdComment)
		require.NotEmpty(t, createdComment.Uuid)
		require.Equal(t, "Great hostel!", createdComment.Content)
		require.Equal(t, db.CommentedOnAccomodation, createdComment.CommentedOn)
		require.Equal(t, acc.Uuid, createdComment.CommentedOnUuid)
		require.Equal(t, userUuidStr, createdComment.UserUuid)
	})

	t.Run("CreateComment - Invalid CommentedOn UUID", func(t *testing.T) {
		commentData := &comments.Comment{
			Content:         "Great hostel!",
			CommentedOn:     db.CommentedOnAccomodation,
			CommentedOnUuid: "invalid-uuid",
			UserUuid:        userUuidStr,
		}
		_, err := repo.CreateComment(ctx, commentData)
		require.Error(t, err)
	})

	t.Run("CreateComment - Invalid User UUID", func(t *testing.T) {
		commentData := &comments.Comment{
			Content:         "Great hostel!",
			CommentedOn:     db.CommentedOnAccomodation,
			CommentedOnUuid: acc.Uuid,
			UserUuid:        "invalid-uuid",
		}
		_, err := repo.CreateComment(ctx, commentData)
		require.Error(t, err)
	})

	t.Run("FindAllCommentsByRelatedEntity", func(t *testing.T) {
		commentData := &comments.Comment{
			Content:         "Another comment",
			CommentedOn:     db.CommentedOnAccomodation,
			CommentedOnUuid: acc.Uuid,
			UserUuid:        userUuidStr,
		}
		_, err := repo.CreateComment(ctx, commentData)
		require.NoError(t, err)

		allComments, err := repo.FindAllCommentsByRelatedEntity(ctx, db.CommentedOnAccomodation, acc.Uuid)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(allComments), 1)

		var found bool
		for _, c := range allComments {
			if c.Content == "Another comment" {
				found = true
				break
			}
		}
		require.True(t, found, "Should find the newly created comment")
	})

	t.Run("FindAllCommentsByRelatedEntity - Invalid UUID", func(t *testing.T) {
		_, err := repo.FindAllCommentsByRelatedEntity(ctx, db.CommentedOnAccomodation, "invalid-uuid")
		require.Error(t, err)
	})

	t.Run("UpdateCommentByUuid", func(t *testing.T) {
		commentData := &comments.Comment{
			Content:         "To update",
			CommentedOn:     db.CommentedOnAccomodation,
			CommentedOnUuid: acc.Uuid,
			UserUuid:        userUuidStr,
		}
		createdComment, err := repo.CreateComment(ctx, commentData)
		require.NoError(t, err)

		err = repo.UpdateCommentByUuid(ctx, createdComment.Uuid, "Updated content")
		require.NoError(t, err)

		updatedComments, err := repo.FindAllCommentsByRelatedEntity(ctx, db.CommentedOnAccomodation, acc.Uuid)
		require.NoError(t, err)

		var found bool
		for _, c := range updatedComments {
			if c.Uuid == createdComment.Uuid && c.Content == "Updated content" {
				found = true
				break
			}
		}
		require.True(t, found, "Should find the updated comment content")
	})

	t.Run("UpdateCommentByUuid - Invalid UUID", func(t *testing.T) {
		err := repo.UpdateCommentByUuid(ctx, "invalid-uuid", "Updated content")
		require.Error(t, err)
	})

	t.Run("DeleteCommentByUuid", func(t *testing.T) {
		commentData := &comments.Comment{
			Content:         "To delete",
			CommentedOn:     db.CommentedOnAccomodation,
			CommentedOnUuid: acc.Uuid,
			UserUuid:        userUuidStr,
		}
		createdComment, err := repo.CreateComment(ctx, commentData)
		require.NoError(t, err)

		err = repo.DeleteCommentByUuid(ctx, createdComment.Uuid)
		require.NoError(t, err)
	})

	t.Run("DeleteCommentByUuid - Invalid UUID", func(t *testing.T) {
		err := repo.DeleteCommentByUuid(ctx, "invalid-uuid")
		require.Error(t, err)
	})
}
