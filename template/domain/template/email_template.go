package template

import (
	"errors"
)

// EmailTemplate 邮件模板
type EmailTemplate struct {
	id           int64
	templateType int // 为了避免和关键字冲突
	topic        string
	content      *Content
}

const (
	// IncompleteTemplate 不完整模板
	IncompleteTemplate = 1
	// CompleteTemplate 完整模板
	CompleteTemplate = 2
)

// NewEmailTemplate 创建一个新的邮件模板
//
// 不允许 topic 和 content 为 nil
func NewEmailTemplate(topic *string, content *string) (*EmailTemplate, error) {
	if topic == nil || content == nil {
		return nil, errors.New("topic and content should not be nil")
	}
	if len(*topic) < 1 {
		return nil, errors.New("topic should not be empty")
	}
	if len(*topic) > 255 {
		return nil, errors.New("topic is too long")
	}
	if len(*content) > 65535 {
		return nil, errors.New("content is too long")
	}
	c := NewContent(*content)

	templateType := CompleteTemplate
	if len(c.Slots()) > 0 {
		templateType = IncompleteTemplate
	}

	return &EmailTemplate{
		templateType: templateType,
		topic:        *topic,
		content:      c,
	}, nil
}

// NewEmailTemplateForUpdate 创建一个新的邮件模板，用于更新
//
// topic 和 content 可以为 nil
func NewEmailTemplateForUpdate(id int64, topic *string, content *string) (*EmailTemplate, error) {
	temp := &EmailTemplate{id: id, templateType: CompleteTemplate}
	if topic != nil {
		if len(*topic) < 1 {
			return nil, errors.New("topic should not be empty")
		}
		if len(*topic) > 255 {
			return nil, errors.New("topic is too long")
		}
		temp.topic = *topic
	}
	if content != nil {
		if len(*content) > 65535 {
			return nil, errors.New("content is too long")
		}
		temp.content = NewContent(*content)
		if len(temp.content.Slots()) > 0 {
			temp.templateType = IncompleteTemplate
		}
	}
	return temp, nil
}

func (t *EmailTemplate) ID() int64 {
	return t.id
}

func (t *EmailTemplate) Type() int {
	return t.templateType
}

func (t *EmailTemplate) Topic() string {
	return t.topic
}

func (t *EmailTemplate) Content() *Content {
	return t.content
}

//func validate(topic string, content string) error {
//	if len(topic) < 1 {
//		return errors.New("topic should not be empty")
//	} else if len(topic) > 255 {
//		return errors.New("topic is too long")
//	}
//
//	if len(content) > 65535 {
//		return errors.New("content is too long")
//	}
//	return nil
//}
