package query

import "time"

type EmailTemplate struct {
	ID         int64     `json:"id"`
	Type       int       `json:"type"`
	Topic      string    `json:"topic"`
	Content    string    `json:"content"`
	Slots      []string  `json:"slots"`
	UpdateTime time.Time `json:"updateTime"`
}
