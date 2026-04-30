package utility

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func TimestamptzFromTime(t *time.Time) pgtype.Timestamptz {
	if t == nil {
		return pgtype.Timestamptz{Valid: false}
	}

	return pgtype.Timestamptz{
		Valid: true,
		Time:  *t,
	}
}

func TimePtrFromTimestamptz(ts pgtype.Timestamptz) *time.Time {
	if !ts.Valid {
		return nil
	}
	t := ts.Time
	return &t
}
