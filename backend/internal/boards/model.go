package boards

import (
	"time"

	"github.com/PegasusMKD/travel-dream-board/internal/accomodations"
	"github.com/PegasusMKD/travel-dream-board/internal/activities"
	"github.com/PegasusMKD/travel-dream-board/internal/db"
	"github.com/PegasusMKD/travel-dream-board/internal/transport"
)

type Board struct {
	Uuid         string          `json:"uuid"`
	Name         string          `json:"name" binding:"required"`
	Details      *string         `json:"details"`
	LocationName string          `json:"location_name" binding:"required"`
	StartsAt     *time.Time      `json:"starts_at"`
	LastsUntil   *time.Time      `json:"lasts_until"`
	Status       db.BoardsStatus `json:"status"`
	ThumbnailUrl *string         `json:"thumbnail_url"`
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

	return &Board{
		Uuid:         ent.Uuid.String(),
		Name:         ent.Name,
		Details:      ent.Details,
		LocationName: ent.LocationName,
		StartsAt:     starts_at,
		LastsUntil:   lasts_until,
		Status:       ent.Status,
		ThumbnailUrl: ent.ThumbnailUrl,
	}
}

type AggregatedBoard struct {
	Board

	Accomodations []*accomodations.Accomodation `json:"accomodations"`
	Transport     []*transport.Transport        `json:"transport"`
	Activities    []*activities.Activity        `json:"activities"`
}
