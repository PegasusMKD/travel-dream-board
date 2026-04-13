package comments

import (
	"github.com/PegasusMKD/travel-dream-board/internal/db"
)

type Comment struct {
	Uuid            string         `json:"uuid"`
	CreatedBy       string         `json:"created_by" binding:"required"`
	Content         string         `json:"content" binding:"required"`
	CommentedOn     db.CommentedOn `json:"commented_on" binding:"required"`
	CommentedOnUuid string         `json:"commented_on_uuid" binding:"required"`
}
