package comments

import (
	"context"

	"github.com/PegasusMKD/travel-dream-board/internal/db"
	"github.com/PegasusMKD/travel-dream-board/internal/utility"
)

type Repository interface {
	CreateComment(ctx context.Context, data *Comment) (*Comment, error)
	FindAllCommentsByRelatedEntity(ctx context.Context, _type db.CommentedOn, uuid string) ([]*Comment, error)
	UpdateCommentByUuid(ctx context.Context, uuid string, content string) error
	DeleteCommentByUuid(ctx context.Context, uuid string) error
}

type commentsRepositoryImpl struct {
	queries *db.Queries
}

func NewRepository(queries *db.Queries) Repository {
	return &commentsRepositoryImpl{
		queries: queries,
	}
}

func (repo *commentsRepositoryImpl) CreateComment(ctx context.Context, data *Comment) (*Comment, error) {
	commentedOnId, err := utility.UuidFromString(data.CommentedOnUuid)
	if err != nil {
		log.Error("Failed parsing UUID", "uuid", data.CommentedOnUuid, "error", err)
		return nil, err
	}

	userUuid, err := utility.UuidFromString(data.UserUuid)
	if err != nil {
		log.Error("Failed parsing UserUuid", "uuid", data.UserUuid, "error", err)
		return nil, err
	}

	params := db.CreateCommentParams{
		CreatedBy:       userUuid,
		Content:         data.Content,
		CommentedOn:     data.CommentedOn,
		CommentedOnUuid: commentedOnId,
	}

	ent, err := repo.queries.CreateComment(ctx, params)
	if err != nil {
		log.Error(
			"Failed creating entity",
			"created_by", data.UserUuid,
			"content", data.Content,
			"commented_on", data.CommentedOn,
			"commented_on_uuid", data.CommentedOnUuid,
			"error", err,
		)
		return nil, err
	}

	return FromEntity(ent), nil
}

func (repo *commentsRepositoryImpl) FindAllCommentsByRelatedEntity(ctx context.Context, _type db.CommentedOn, uuid string) ([]*Comment, error) {
	commentedOnId, err := utility.UuidFromString(uuid)
	if err != nil {
		log.Error("Failed parsing UUID", "uuid", uuid, "error", err)
		return nil, err
	}

	params := db.FindAllCommentsByCommentedOnUuidParams{
		CommentedOn:     _type,
		CommentedOnUuid: commentedOnId,
	}

	entities, err := repo.queries.FindAllCommentsByCommentedOnUuid(ctx, params)
	if err != nil {
		log.Error("Failed fetching comments", "type", _type, "commented_on_uuid", uuid, "error", err)
		return nil, err
	}

	data := []*Comment{}
	for _, val := range entities {
		data = append(data, FromEntity(val))
	}

	return data, nil
}

func (repo *commentsRepositoryImpl) UpdateCommentByUuid(ctx context.Context, uuid string, content string) error {
	id, err := utility.UuidFromString(uuid)
	if err != nil {
		log.Error("Failed parsing UUID", "uuid", uuid, "error", err)
		return err
	}

	params := db.UpdateCommentByUuidParams{
		Uuid:    id,
		Content: content,
	}

	return repo.queries.UpdateCommentByUuid(ctx, params)
}

func (repo *commentsRepositoryImpl) DeleteCommentByUuid(ctx context.Context, uuid string) error {
	id, err := utility.UuidFromString(uuid)
	if err != nil {
		log.Error("Failed parsing UUID", "uuid", uuid, "error", err)
		return err
	}

	return repo.queries.DeleteCommentByUuid(ctx, id)
}
