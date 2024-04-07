package template

import (
	"regexp"
)

type Content struct {
	String string
	Valid  bool
}

func (c *Content) Slots() []string {
	return extractPlaceholders(c.String)
}

func extractPlaceholders(text string) []string {
	re := regexp.MustCompile(`{{([^}]+)}}`)
	matches := re.FindAllStringSubmatch(text, -1)
	placeholders := make([]string, 0, len(matches))
	for _, match := range matches {
		placeholders = append(placeholders, match[1])
	}
	return placeholders
}
