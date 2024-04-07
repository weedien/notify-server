package query

type EmailTemplate struct {
	ID      int      `json:"id"`
	Type    int      `json:"type"`
	Topic   string   `json:"topic"`
	Content string   `json:"content"`
	Slots   []string `json:"slots"`
}
