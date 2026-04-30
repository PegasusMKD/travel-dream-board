package transport

import (
	"time"

	"github.com/PegasusMKD/travel-dream-board/internal/comments"
	"github.com/PegasusMKD/travel-dream-board/internal/db"
	"github.com/PegasusMKD/travel-dream-board/internal/utility"
	"github.com/PegasusMKD/travel-dream-board/internal/votes"
)

type Transport struct {
	Uuid             string             `json:"uuid"`
	Url              string             `json:"url" binding:"required"`
	Title            string             `json:"title" binding:"required"`
	BoardUuid        string             `json:"board_uuid" binding:"required"`
	ImageUrl         *string            `json:"image_url"`
	Notes            *string            `json:"notes"`
	Status           db.TransportStatus `json:"status"`
	BookingReference *string            `json:"booking_reference"`
	Selected         bool               `json:"selected"`
	UserUuid         string             `json:"user_uuid"`
	AvgRating        float64            `json:"avg_rating"`
	RatingCount      int32              `json:"rating_count"`

	OutboundDepartingLocation *string    `json:"outbound_departing_location"`
	OutboundArrivingLocation  *string    `json:"outbound_arriving_location"`
	OutboundDepartingAt       *time.Time `json:"outbound_departing_at"`
	OutboundArrivingAt        *time.Time `json:"outbound_arriving_at"`
	InboundDepartingLocation  *string    `json:"inbound_departing_location"`
	InboundArrivingLocation   *string    `json:"inbound_arriving_location"`
	InboundDepartingAt        *time.Time `json:"inbound_departing_at"`
	InboundArrivingAt         *time.Time `json:"inbound_arriving_at"`

	OutboundDurationMinutes *int32 `json:"outbound_duration_minutes"`
	InboundDurationMinutes  *int32 `json:"inbound_duration_minutes"`
}

func FromEntity(entity db.Transport) *Transport {
	return &Transport{
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

		OutboundDepartingLocation: entity.OutboundDepartingLocation,
		OutboundArrivingLocation:  entity.OutboundArrivingLocation,
		OutboundDepartingAt:       utility.TimePtrFromTimestamptz(entity.OutboundDepartingAt),
		OutboundArrivingAt:        utility.TimePtrFromTimestamptz(entity.OutboundArrivingAt),
		InboundDepartingLocation:  entity.InboundDepartingLocation,
		InboundArrivingLocation:   entity.InboundArrivingLocation,
		InboundDepartingAt:        utility.TimePtrFromTimestamptz(entity.InboundDepartingAt),
		InboundArrivingAt:         utility.TimePtrFromTimestamptz(entity.InboundArrivingAt),

		OutboundDurationMinutes: entity.OutboundDurationMinutes,
		InboundDurationMinutes:  entity.InboundDurationMinutes,
	}
}

func FromGetTransportRow(entity db.GetTransportByUuidRow) *Transport {
	return &Transport{
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

		OutboundDepartingLocation: entity.OutboundDepartingLocation,
		OutboundArrivingLocation:  entity.OutboundArrivingLocation,
		OutboundDepartingAt:       utility.TimePtrFromTimestamptz(entity.OutboundDepartingAt),
		OutboundArrivingAt:        utility.TimePtrFromTimestamptz(entity.OutboundArrivingAt),
		InboundDepartingLocation:  entity.InboundDepartingLocation,
		InboundArrivingLocation:   entity.InboundArrivingLocation,
		InboundDepartingAt:        utility.TimePtrFromTimestamptz(entity.InboundDepartingAt),
		InboundArrivingAt:         utility.TimePtrFromTimestamptz(entity.InboundArrivingAt),

		OutboundDurationMinutes: entity.OutboundDurationMinutes,
		InboundDurationMinutes:  entity.InboundDurationMinutes,
	}
}

func FromFindTransportRow(entity db.FindAllTransportByBoardUuidRow) *Transport {
	return &Transport{
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

		OutboundDepartingLocation: entity.OutboundDepartingLocation,
		OutboundArrivingLocation:  entity.OutboundArrivingLocation,
		OutboundDepartingAt:       utility.TimePtrFromTimestamptz(entity.OutboundDepartingAt),
		OutboundArrivingAt:        utility.TimePtrFromTimestamptz(entity.OutboundArrivingAt),
		InboundDepartingLocation:  entity.InboundDepartingLocation,
		InboundArrivingLocation:   entity.InboundArrivingLocation,
		InboundDepartingAt:        utility.TimePtrFromTimestamptz(entity.InboundDepartingAt),
		InboundArrivingAt:         utility.TimePtrFromTimestamptz(entity.InboundArrivingAt),

		OutboundDurationMinutes: entity.OutboundDurationMinutes,
		InboundDurationMinutes:  entity.InboundDurationMinutes,
	}
}

type AggregatedTransport struct {
	Transport

	Comments []*comments.Comment `json:"comments"`
	Votes    []*votes.Vote       `json:"votes"`
}
