package main

import (
	"context"
	"github.com/PegasusMKD/travel-dream-board/internal/db"
	"github.com/PegasusMKD/travel-dream-board/internal/utility"
	"github.com/jackc/pgx/v5/pgtype"
)

// testing logic
func getParam(userUuid *string) pgtype.UUID {
	var p pgtype.UUID
	if userUuid != nil && *userUuid != "" {
		p.Scan(*userUuid)
	}
	return p
}
