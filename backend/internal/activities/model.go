package activities

import (
	"time"

	"github.com/PegasusMKD/travel-dream-board/internal/comments"
	"github.com/PegasusMKD/travel-dream-board/internal/db"
	"github.com/PegasusMKD/travel-dream-board/internal/utility"
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

	StartAt *time.Time `json:"start_at"`
	EndAt   *time.Time `json:"end_at"`
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
		StartAt:          utility.TimePtrFromTimestamptz(entity.StartAt),
		EndAt:            utility.TimePtrFromTimestamptz(entity.EndAt),
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
		StartAt:          utility.TimePtrFromTimestamptz(entity.StartAt),
		EndAt:            utility.TimePtrFromTimestamptz(entity.EndAt),
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
		StartAt:          utility.TimePtrFromTimestamptz(entity.StartAt),
		EndAt:            utility.TimePtrFromTimestamptz(entity.EndAt),
	}
}

type AggregatedActivity struct {
	Activity

	Comments []*comments.Comment `json:"comments"`
	Votes    []*votes.Vote       `json:"votes"`
}
