package template

import (
	"errors"
)

// EmailTemplate 邮件模板
type EmailTemplate struct {
	id           int
	templateType int // 为了避免和关键字冲突
	topic        string
	content      Content
	slots        []string
}

const (
	// IncompleteTemplate 不完整模板
	IncompleteTemplate = 1
	// CompleteTemplate 完整模板
	CompleteTemplate = 2
)

func NewEmailTemplate(topic string, content string) (*EmailTemplate, error) {
	if len(topic) > 255 {
		return nil, errors.New("topic is too long")
	}
	if len(content) > 65535 {
		return nil, errors.New("content is too long")
	}
	templateType := CompleteTemplate
	slots := extractPlaceholders(content)
	if len(slots) > 0 {
		// 不完整模板
		templateType = IncompleteTemplate
	}

	return &EmailTemplate{
		templateType: templateType,
		topic:        topic,
		content:      Content{String: content, Valid: true},
		slots:        slots,
	}, nil
}

func NewEmptyEmailTemplate() *EmailTemplate {
	return &EmailTemplate{}
}

func (t *EmailTemplate) SetTopic(topic string) {
	t.topic = topic
}

func (t *EmailTemplate) SetContent(content string) {
	t.content = Content{String: content, Valid: true}
	t.slots = t.content.Slots()
	if len(t.slots) > 0 {
		t.templateType = IncompleteTemplate
	} else {
		t.templateType = CompleteTemplate
	}
}

//func NewEmailTemplateWithID(id int, topic string, content string) (*EmailTemplate, error) {
//	if id <= 0 {
//		return nil, errors.New("id is invalid")
//	}
//	if len(topic) > 255 {
//		return nil, errors.New("topic is too long")
//	}
//	if len(content) > 65535 {
//		return nil, errors.New("content is too long")
//	}
//	templateType := 2
//	slots := extractPlaceholders(content)
//	if len(slots) > 0 {
//		// 不完整模板
//		templateType = 1
//	}
//	return &EmailTemplate{id: id, templateType: templateType, topic: topic, content: content, slots: slots}, nil
//}

func (t *EmailTemplate) ID() int {
	return t.id
}

func (t *EmailTemplate) Type() int {
	return t.templateType
}

func (t *EmailTemplate) Topic() string {
	return t.topic
}

func (t *EmailTemplate) Content() Content {
	return t.content
}

func (t *EmailTemplate) Slots() []string {
	return t.slots
}
