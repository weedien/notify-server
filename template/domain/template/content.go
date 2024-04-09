package template

import (
	"regexp"
)

type Content struct {
	inner string
	slots []string
}

func NewContent(content string) *Content {
	return &Content{
		inner: content,
		slots: extractPlaceholders(content),
	}
}

func (c *Content) String() string {
	return c.inner
}

func (c *Content) Slots() []string {
	return c.slots
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

func (c *Content) ContainsAllSlots(slots []string) bool {
	for _, slot := range slots {
		if !c.containsSlot(slot) {
			return false
		}
	}
	return true
}

func (c *Content) containsSlot(slot string) bool {
	for _, s := range c.slots {
		if s == slot {
			return true
		}
	}
	return false
}
