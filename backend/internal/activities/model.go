package activities

import (
	"github.com/PegasusMKD/travel-dream-board/internal/comments"
	"github.com/PegasusMKD/travel-dream-board/internal/db"
	"github.com/PegasusMKD/travel-dream-board/internal/votes"
)

type Activity struct {
	Uuid             string              `json:"uuid"`
	Url              string              `json:"url" binding:"required"`
	Title            string              `json:"title" binding:"required"`
	BoardUuid        string              `json:"board_uuid" binding:"required"`
	ImageUrl         *string             `json:"image_url"`
	Notes            *string             `json:"notes"`
	Status           db.ActivitiesStatus `json:"status"`
	BookingReference *string             `json:"booking_reference"`
	Selected         bool                `json:"selected"`
	UserUuid         string              `json:"user_uuid"`
	AvgRating        float64             `json:"avg_rating"`
	RatingCount      int32               `json:"rating_count"`
}

func FromEntity(entity db.Activity) *Activity {
	return &Activity{
		Uuid:             entity.Uuid.String(),
		BoardUuid:        entity.BoardUuid.String(),
		Url:              entity.Url,
		Title:            entity.Title,
		ImageUrl:         entity.ImageUrl,
		Notes:            entity.Notes,
		Status:           entity.Status,
		BookingReference: entity.BookingReference,
		Selected:         entity.Selected,
		UserUuid:         entity.UserUuid.String(),
	}
}

func FromGetActivityRow(entity db.GetActivityByUuidRow) *Activity {
	return &Activity{
		Uuid:             entity.Uuid.String(),
		BoardUuid:        entity.BoardUuid.String(),
		Url:              entity.Url,
		Title:            entity.Title,
		ImageUrl:         entity.ImageUrl,
		Notes:            entity.Notes,
		Status:           entity.Status,
		BookingReference: entity.BookingReference,
		Selected:         entity.Selected,
		UserUuid:         entity.UserUuid.String(),
		AvgRating:        entity.AvgRating,
		RatingCount:      entity.RatingCount,
	}
}

func FromFindActivitiesRow(entity db.FindAllActivitiesByBoardUuidRow) *Activity {
	return &Activity{
		Uuid:             entity.Uuid.String(),
		BoardUuid:        entity.BoardUuid.String(),
		Url:              entity.Url,
		Title:            entity.Title,
		ImageUrl:         entity.ImageUrl,
		Notes:            entity.Notes,
		Status:           entity.Status,
		BookingReference: entity.BookingReference,
		Selected:         entity.Selected,
		UserUuid:         entity.UserUuid.String(),
		AvgRating:        entity.AvgRating,
		RatingCount:      entity.RatingCount,
	}
}

type AggregatedActivity struct {
	Activity

	Comments []*comments.Comment `json:"comments"`
	Votes    []*votes.Vote       `json:"votes"`
}
