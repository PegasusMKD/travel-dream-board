package votes_test

import (
	"context"
	"testing"

	"github.com/PegasusMKD/travel-dream-board/internal/accomodations"
	"github.com/PegasusMKD/travel-dream-board/internal/boards"
	"github.com/PegasusMKD/travel-dream-board/internal/db"
	"github.com/PegasusMKD/travel-dream-board/internal/testutil"
	"github.com/PegasusMKD/travel-dream-board/internal/votes"
	"github.com/stretchr/testify/require"
)

func TestVotesRepository(t *testing.T) {
	env := testutil.SetupTestDB(t)
	repo := votes.NewRepository(env.Queries)
	ctx := context.Background()

	// 1. Create a user
	email := "voteuser@example.com"
	user, err := env.Queries.UpsertUser(ctx, db.UpsertUserParams{
		Email: &email,
		Name:  "Test Vote User",
	})
	require.NoError(t, err)
	
	userUuidStrBytes, _ := user.Uuid.Value()
	userUuidStr := userUuidStrBytes.(string)

	// 2. Create a board
	boardRepo := boards.NewRepository(env.Queries)
	board, err := boardRepo.CreateBoard(ctx, &boards.Board{
		Name:         "Vote Board",
		LocationName: "Vote Location",
		UserUuid:     &userUuidStr,
	})
	require.NoError(t, err)

	// 3. Create an accomodation
	accRepo := accomodations.NewRepository(env.Queries)
	acc, err := accRepo.CreateAccomodation(ctx, &accomodations.Accomodation{
		Url:       "http://example.com/votehostel",
		Title:     "Vote Hostel",
		BoardUuid: board.Uuid,
		UserUuid:  userUuidStr,
	})
	require.NoError(t, err)

	t.Run("CreateVote", func(t *testing.T) {
		voteData := &votes.Vote{
			Rank:        1,
			VotedOn:     db.VotedOnAccomodation,
			VotedOnUuid: acc.Uuid,
			UserUuid:    userUuidStr,
		}
		
		createdVote, err := repo.CreateVote(ctx, voteData)
		require.NoError(t, err)
		require.NotNil(t, createdVote)
		require.NotEmpty(t, createdVote.Uuid)
		require.Equal(t, int32(1), createdVote.Rank)
		require.Equal(t, db.VotedOnAccomodation, createdVote.VotedOn)
		require.Equal(t, acc.Uuid, createdVote.VotedOnUuid)
		require.Equal(t, userUuidStr, createdVote.UserUuid)
	})

	t.Run("CreateVote - Invalid VotedOn UUID", func(t *testing.T) {
		voteData := &votes.Vote{
			Rank:        1,
			VotedOn:     db.VotedOnAccomodation,
			VotedOnUuid: "invalid-uuid",
			UserUuid:    userUuidStr,
		}
		_, err := repo.CreateVote(ctx, voteData)
		require.Error(t, err)
	})

	t.Run("CreateVote - Invalid User UUID", func(t *testing.T) {
		voteData := &votes.Vote{
			Rank:        1,
			VotedOn:     db.VotedOnAccomodation,
			VotedOnUuid: acc.Uuid,
			UserUuid:    "invalid-uuid",
		}
		_, err := repo.CreateVote(ctx, voteData)
		require.Error(t, err)
	})

	t.Run("FindAllVotesByRelatedEntity", func(t *testing.T) {
		voteData := &votes.Vote{
			Rank:        2,
			VotedOn:     db.VotedOnAccomodation,
			VotedOnUuid: acc.Uuid,
			UserUuid:    userUuidStr,
		}
		_, err := repo.CreateVote(ctx, voteData)
		require.NoError(t, err)

		allVotes, err := repo.FindAllVotesByRelatedEntity(ctx, db.VotedOnAccomodation, acc.Uuid)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(allVotes), 1)
		
		var found bool
		for _, v := range allVotes {
			if v.Rank == 2 {
				found = true
				break
			}
		}
		require.True(t, found, "Should find the newly created vote")
	})

	t.Run("FindAllVotesByRelatedEntity - Invalid UUID", func(t *testing.T) {
		_, err := repo.FindAllVotesByRelatedEntity(ctx, db.VotedOnAccomodation, "invalid-uuid")
		require.Error(t, err)
	})

	t.Run("UpdateVoteByUuid", func(t *testing.T) {
		voteData := &votes.Vote{
			Rank:        3,
			VotedOn:     db.VotedOnAccomodation,
			VotedOnUuid: acc.Uuid,
			UserUuid:    userUuidStr,
		}
		createdVote, err := repo.CreateVote(ctx, voteData)
		require.NoError(t, err)

		err = repo.UpdateVoteByUuid(ctx, createdVote.Uuid, 5)
		require.NoError(t, err)

		updatedVotes, err := repo.FindAllVotesByRelatedEntity(ctx, db.VotedOnAccomodation, acc.Uuid)
		require.NoError(t, err)
		
		var found bool
		for _, v := range updatedVotes {
			if v.Uuid == createdVote.Uuid && v.Rank == 5 {
				found = true
				break
			}
		}
		require.True(t, found, "Should find the updated vote rank")
	})

	t.Run("UpdateVoteByUuid - Invalid UUID", func(t *testing.T) {
		err := repo.UpdateVoteByUuid(ctx, "invalid-uuid", 5)
		require.Error(t, err)
	})

	t.Run("DeleteVoteByUuid", func(t *testing.T) {
		voteData := &votes.Vote{
			Rank:        4,
			VotedOn:     db.VotedOnAccomodation,
			VotedOnUuid: acc.Uuid,
			UserUuid:    userUuidStr,
		}
		createdVote, err := repo.CreateVote(ctx, voteData)
		require.NoError(t, err)

		err = repo.DeleteVoteByUuid(ctx, createdVote.Uuid)
		require.NoError(t, err)
	})

	t.Run("DeleteVoteByUuid - Invalid UUID", func(t *testing.T) {
		err := repo.DeleteVoteByUuid(ctx, "invalid-uuid")
		require.Error(t, err)
	})
}
