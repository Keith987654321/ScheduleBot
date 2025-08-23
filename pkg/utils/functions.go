package utils

import (
	"unicode"
)

func CapitalizeFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}

	runeArr := []rune(s)
	runeArr[0] = unicode.ToUpper(runeArr[0])
	return string(runeArr)
}
