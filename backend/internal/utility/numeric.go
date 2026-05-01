package utility

import (
	"math/big"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
)

// NumericFromString converts a decimal string (e.g. "1234.56") to a pgtype.Numeric.
// Returns an invalid Numeric if s is nil or empty/unparseable so the column stores NULL.
func NumericFromString(s *string) pgtype.Numeric {
	if s == nil {
		return pgtype.Numeric{Valid: false}
	}
	raw := strings.TrimSpace(*s)
	if raw == "" {
		return pgtype.Numeric{Valid: false}
	}

	var n pgtype.Numeric
	if err := n.Scan(raw); err != nil {
		return pgtype.Numeric{Valid: false}
	}
	return n
}

// NumericToString formats a pgtype.Numeric back to a string with up to 2 fractional
// digits. Returns nil if the value is null/invalid.
func NumericToString(n pgtype.Numeric) *string {
	if !n.Valid || n.NaN {
		return nil
	}

	if n.Int == nil {
		zero := "0.00"
		return &zero
	}

	abs := new(big.Int).Abs(n.Int).String()
	exp := int(n.Exp)

	var s string
	if exp >= 0 {
		s = abs + strings.Repeat("0", exp) + ".00"
	} else {
		decimals := -exp
		if decimals >= len(abs) {
			abs = strings.Repeat("0", decimals-len(abs)+1) + abs
		}
		intPart := abs[:len(abs)-decimals]
		fracPart := abs[len(abs)-decimals:]
		// Pad/trim to exactly 2 decimal places
		if len(fracPart) < 2 {
			fracPart = fracPart + strings.Repeat("0", 2-len(fracPart))
		} else if len(fracPart) > 2 {
			// Round half-up to 2 decimals
			rounded, _ := strconv.Atoi(fracPart[:2])
			if fracPart[2] >= '5' {
				rounded++
				if rounded == 100 {
					rounded = 0
					i, _ := new(big.Int).SetString(intPart, 10)
					i.Add(i, big.NewInt(1))
					intPart = i.String()
				}
			}
			fracPart = strconv.Itoa(rounded)
			if len(fracPart) < 2 {
				fracPart = strings.Repeat("0", 2-len(fracPart)) + fracPart
			}
		}
		s = intPart + "." + fracPart
	}

	if n.Int.Sign() < 0 {
		s = "-" + s
	}
	return &s
}
