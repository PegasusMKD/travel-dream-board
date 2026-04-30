package utility

import "time"

// ParseWallClockTime parses an LLM-extracted datetime string and re-anchors
// it to UTC so the wall-clock components (Y/M/D H:M:S) survive any timezone
// conversion unchanged. Returns nil if the input is empty or unparseable.
func ParseWallClockTime(s *string) *time.Time {
	if s == nil || *s == "" {
		return nil
	}
	layouts := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05",
		"2006-01-02T15:04",
		"2006-01-02 15:04:05Z07:00",
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, *s); err == nil {
			wall := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, time.UTC)
			return &wall
		}
	}
	return nil
}
