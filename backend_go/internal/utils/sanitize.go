package utils

import (
	"regexp"
	"strings"
)

var consecutiveSpaces = regexp.MustCompile(`\s+`)

func SanitizeString(input string) string {
	return consecutiveSpaces.ReplaceAllString(strings.TrimSpace(input), " ")
}
