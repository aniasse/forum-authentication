package tools

import (
	"strings"
)

func IsInvalid(input string) bool {
	var rep bool

	if len([]rune(input)) > 1500 {
		rep = true
	} else if strings.TrimSpace(input) == "" {
		rep = true
	}

	return rep
}

func ValidExtension(s string) bool {
	return strings.HasSuffix(s, ".jpeg") || strings.HasSuffix(s, ".png") || strings.HasSuffix(s, ".avif") || strings.HasSuffix(s, ".jpg") || strings.HasSuffix(s, ".gif")
}
