package jsonnewline

import "slices"

func NewLineToEscape(s string) string {
	inQuotes := false
	inBackslash := false
	stringArr := []rune(s)
	for i := 0; i < len(stringArr); i++ {
		s := stringArr[i]
		if inBackslash {
			inBackslash = false
			continue
		}
		switch s {
		case '\\':
			inBackslash = true
		case '"':
			inQuotes = !inQuotes
		case '\n':
			if inQuotes {
				stringArr[i] = '\\'
				i++
				stringArr = slices.Insert(stringArr, i, 'n')
			}
		}
	}
	return string(stringArr)
}
