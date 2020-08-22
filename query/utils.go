package query

import "unicode"

// IsChar is char
func IsChar(c rune) bool {
	return !unicode.IsLetter(c) && !unicode.IsNumber(c) && c != '.' && c != '*'
}
