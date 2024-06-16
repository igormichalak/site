package view

import (
	"strings"
	"unicode/utf8"
)

func transformIgnoreCase(s string, from string, to func(match string) string, ignore []string) string {
	lowerS := []rune(strings.ToLower(s))
	lowerFrom := []rune(strings.ToLower(from))
	sRunes := []rune(s)

	var result strings.Builder

	i := 0
OuterLoop:
	for i < len(lowerS) {
		for _, is := range ignore {
			if strings.HasPrefix(string(lowerS[i:]), is) {
				result.WriteString(is)
				i += utf8.RuneCountInString(is)
				continue OuterLoop
			}
		}
		if strings.HasPrefix(string(lowerS[i:]), string(lowerFrom)) {
			match := sRunes[i : i+len(lowerFrom)]
			substitute := to(string(match))
			result.WriteString(substitute)
			i += len(match)
		} else {
			result.WriteRune(sRunes[i])
			i++
		}
	}

	return result.String()
}

func highlight(s string, words []string) string {
	for _, word := range words {
		s = transformIgnoreCase(s, word, func(match string) string {
			return `<span class="hl">` + match + `</span>`
		}, []string{`<span class="hl">`, `</span>`})
	}
	return s
}
