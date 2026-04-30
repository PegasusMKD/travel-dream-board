package utility

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	durationHmRegex    = regexp.MustCompile(`^(?:(\d+)\s*h(?:ours?|rs?)?)?\s*(?:(\d+)\s*m(?:in(?:ute)?s?)?)?$`)
	durationColonRegex = regexp.MustCompile(`^(\d+):(\d+)$`)
)

// ParseDurationMinutes accepts whatever the LLM returns (e.g. "225", "3h 45m",
// "3:45") and best-effort converts to minutes. Returns nil if it can't make
// sense of the value.
func ParseDurationMinutes(s *string) *int32 {
	if s == nil {
		return nil
	}
	raw := strings.TrimSpace(*s)
	if raw == "" {
		return nil
	}

	if n, err := strconv.Atoi(raw); err == nil && n >= 0 {
		v := int32(n)
		return &v
	}

	lower := strings.ToLower(raw)
	var hours, mins int
	matched := false

	if m := durationHmRegex.FindStringSubmatch(lower); m != nil && (m[1] != "" || m[2] != "") {
		if m[1] != "" {
			hours, _ = strconv.Atoi(m[1])
		}
		if m[2] != "" {
			mins, _ = strconv.Atoi(m[2])
		}
		matched = true
	} else if m := durationColonRegex.FindStringSubmatch(lower); m != nil {
		hours, _ = strconv.Atoi(m[1])
		mins, _ = strconv.Atoi(m[2])
		matched = true
	}

	if !matched {
		return nil
	}
	v := int32(hours*60 + mins)
	return &v
}
