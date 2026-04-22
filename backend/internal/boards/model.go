package boards

import (
	"time"

	"github.com/PegasusMKD/travel-dream-board/internal/accomodations"
	"github.com/PegasusMKD/travel-dream-board/internal/activities"
	"github.com/PegasusMKD/travel-dream-board/internal/db"
	"github.com/PegasusMKD/travel-dream-board/internal/transport"
)

type Board struct {
	Uuid               string          `json:"uuid"`
	Name               string          `json:"name" binding:"required"`
	Details            *string         `json:"details"`
	LocationName       string          `json:"location_name" binding:"required"`
	StartsAt           *time.Time      `json:"starts_at"`
	LastsUntil         *time.Time      `json:"lasts_until"`
	Status             db.BoardsStatus `json:"status"`
	ThumbnailUrl       *string         `json:"thumbnail_url"`
	UserUuid           *string         `json:"user_uuid,omitempty"`
	AccomodationsCount int64           `json:"accomodations_count"`
	TransportCount     int64           `json:"transport_count"`
	ActivitiesCount    int64           `json:"activities_count"`
}

func FromEntity(ent *db.Board) *Board {
	var starts_at *time.Time
	if ent.StartsAt.Valid {
		starts_at = &ent.StartsAt.Time
	}

	var lasts_until *time.Time
	if ent.LastsUntil.Valid {
		lasts_until = &ent.LastsUntil.Time
	}

	var userUuid *string
	if ent.UserUuid.Valid {
		val := ent.UserUuid.String()
		userUuid = &val
	}

	return &Board{
		Uuid:         ent.Uuid.String(),
		Name:         ent.Name,
		Details:      ent.Details,
		LocationName: ent.LocationName,
		StartsAt:     starts_at,
		LastsUntil:   lasts_until,
		Status:       ent.Status,
		ThumbnailUrl: ent.ThumbnailUrl,
		UserUuid:     userUuid,
	}
}

func FromGetAllBoardsRow(row *db.GetAllBoardsRow) *Board {
	var starts_at *time.Time
	if row.StartsAt.Valid {
		starts_at = &row.StartsAt.Time
	}

	var lasts_until *time.Time
	if row.LastsUntil.Valid {
		lasts_until = &row.LastsUntil.Time
	}

	var userUuid *string
	if row.UserUuid.Valid {
		val := row.UserUuid.String()
		userUuid = &val
	}

	return &Board{
		Uuid:               row.Uuid.String(),
		Name:               row.Name,
		Details:            row.Details,
		LocationName:       row.LocationName,
		StartsAt:           starts_at,
		LastsUntil:         lasts_until,
		Status:             row.Status,
		ThumbnailUrl:       row.ThumbnailUrl,
		UserUuid:           userUuid,
		AccomodationsCount: row.AccomodationsCount,
		TransportCount:     row.TransportCount,
		ActivitiesCount:    row.ActivitiesCount,
	}
}

type AggregatedBoard struct {
	Board

	Accomodations []*accomodations.Accomodation `json:"accomodations"`
	Transport     []*transport.Transport        `json:"transport"`
	Activities    []*activities.Activity        `json:"activities"`
}
