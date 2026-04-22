package comments

import (
	"github.com/PegasusMKD/travel-dream-board/internal/db"
)

type Comment struct {
	Uuid            string         `json:"uuid"`
	UserUuid        string         `json:"user_uuid"`
	UserName        string         `json:"user_name,omitempty"`
	Content         string         `json:"content" binding:"required"`
	CommentedOn     db.CommentedOn `json:"commented_on" binding:"required"`
	CommentedOnUuid string         `json:"commented_on_uuid" binding:"required"`
}

func FromEntity(entity db.Comment) *Comment {
	return &Comment{
		Uuid:            entity.Uuid.String(),
		UserUuid:        entity.CreatedBy.String(),
		Content:         entity.Content,
		CommentedOn:     entity.CommentedOn,
		CommentedOnUuid: entity.CommentedOnUuid.String(),
	}
}

func FromRow(row db.FindAllCommentsByCommentedOnUuidRow) *Comment {
	return &Comment{
		Uuid:            row.Uuid.String(),
		UserUuid:        row.CreatedBy.String(),
		UserName:        row.UserName,
		Content:         row.Content,
		CommentedOn:     row.CommentedOn,
		CommentedOnUuid: row.CommentedOnUuid.String(),
	}
}
