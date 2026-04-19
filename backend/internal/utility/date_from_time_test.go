package utility

import (
	"testing"
	"time"
)

func TestDateFromTime(t *testing.T) {
	tm := time.Date(2023, 10, 25, 12, 0, 0, 0, time.UTC)
	date := DateFromTime(&tm)

	if !date.Valid {
		t.Error("expected valid date")
	}
	if date.Time != tm {
		t.Errorf("expected %v, got %v", tm, date.Time)
	}

	dateNil := DateFromTime(nil)
	if dateNil.Valid {
		t.Error("expected invalid date for nil time")
	}
}
