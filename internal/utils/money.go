package utils

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseAmount converts a decimal string like "10.50" to paise (int64 = 1050).
// This avoids all floating-point rounding issues with money.
func ParseAmount(s string) (int64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("amount is required")
	}

	parts := strings.Split(s, ".")
	if len(parts) > 2 {
		return 0, fmt.Errorf("invalid amount format")
	}

	rupees, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil || rupees < 0 {
		return 0, fmt.Errorf("invalid amount")
	}

	var paise int64
	if len(parts) == 2 {
		// Normalise to exactly 2 decimal places
		dec := parts[1]
		if len(dec) > 2 {
			return 0, fmt.Errorf("amount can have at most 2 decimal places")
		}
		if len(dec) == 1 {
			dec += "0"
		}
		p, err := strconv.ParseInt(dec, 10, 64)
		if err != nil || p < 0 {
			return 0, fmt.Errorf("invalid decimal part")
		}
		paise = p
	}

	total := rupees*100 + paise
	if total <= 0 {
		return 0, fmt.Errorf("amount must be greater than zero")
	}
	return total, nil
}

// FormatAmount converts paise back to "₹10.50" display string.
func FormatAmount(paise int64) string {
	rupees := paise / 100
	cents := paise % 100
	return fmt.Sprintf("₹%d.%02d", rupees, cents)
}
