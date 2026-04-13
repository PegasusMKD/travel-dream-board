package utility

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func DateFromTime(t *time.Time) pgtype.Date {
	if t == nil {
		return pgtype.Date{
			Valid: false,
		}
	}

	return pgtype.Date{
		Valid: true,
		Time:  *t,
	}
}
