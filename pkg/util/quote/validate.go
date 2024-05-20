package quote

import "fmt"

// ValidateBalancedQuotes quote balance detector used for Valdation
func ValidateBalancedQuotes(val any) error {
	s, ok := val.(string)
	if !ok {
		return fmt.Errorf("invalid type")
	}

	quoteCount := make(map[rune]int)

	for i, ch := range s {
		if i > 0 && s[i-1] == '\\' {
			continue
		}

		if ch == '\'' || ch == '"' {
			quoteCount[ch]++
		}
	}

	for _, count := range quoteCount {
		if count%2 != 0 {
			return fmt.Errorf("unbalanced quotes in %s", s)
		}
	}

	return nil
}
