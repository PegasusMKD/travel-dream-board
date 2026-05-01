package utility

import (
	"strings"

	"github.com/PegasusMKD/travel-dream-board/internal/db"
)

// ParseCurrencyCode normalizes a free-form currency string (PLN, "zł", "€",
// "ден", "USD", etc.) into the constrained db.CurrencyCode enum. Returns nil
// for nil/empty input, falls back to "unknown" when a price-bearing string is
// provided but the currency is not one of the supported codes.
func ParseCurrencyCode(s *string) *db.CurrencyCode {
	if s == nil {
		return nil
	}
	raw := strings.TrimSpace(*s)
	if raw == "" {
		return nil
	}

	lower := strings.ToLower(raw)

	switch lower {
	case "pln", "zł", "zl", "zloty", "złoty", "polish zloty", "polski złoty":
		c := db.CurrencyCodePLN
		return &c
	case "eur", "€", "euro":
		c := db.CurrencyCodeEUR
		return &c
	case "mkd", "ден", "denar", "macedonian denar":
		c := db.CurrencyCodeMKD
		return &c
	case "unknown":
		c := db.CurrencyCodeUnknown
		return &c
	}

	c := db.CurrencyCodeUnknown
	return &c
}

// NullCurrencyFromPtr converts an optional CurrencyCode pointer to a sqlc-style
// NullCurrencyCode for use in db params.
func NullCurrencyFromPtr(c *db.CurrencyCode) db.NullCurrencyCode {
	if c == nil {
		return db.NullCurrencyCode{Valid: false}
	}
	return db.NullCurrencyCode{Valid: true, CurrencyCode: *c}
}

// CurrencyPtrFromNull is the inverse of NullCurrencyFromPtr — extracts a *CurrencyCode
// from a sqlc-style NullCurrencyCode.
func CurrencyPtrFromNull(n db.NullCurrencyCode) *db.CurrencyCode {
	if !n.Valid {
		return nil
	}
	c := n.CurrencyCode
	return &c
}
