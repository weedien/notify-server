package command

import (
	"context"
	"errors"

	"github.com/sirupsen/logrus"
	"github.com/weedien/notify-server/common/decorator"
	"github.com/weedien/notify-server/template/domain/template"
)

type UpdateTemplateCommand struct {
	TemplateID int
	Topic      *string // 可选
	Content    *string // 可选
}

type UpdateTemplateHandler decorator.CommandHandler[UpdateTemplateCommand]

type updateTemplateHandler struct {
	templateRepo template.Repository
}

func NewUpdateTemplateContentHandler(
	templateRepo template.Repository,
	logger *logrus.Entry,
) UpdateTemplateHandler {
	if templateRepo == nil {
		panic("templateRepo is nil")
	}
	return decorator.ApplyCommandDecorators[UpdateTemplateCommand](
		updateTemplateHandler{templateRepo: templateRepo},
		logger,
	)
}

// Handle 更新模板
//
// topic为空串是不合法的，content为空串是合法的，所以需要通过nil来判断是否传入了content
func (h updateTemplateHandler) Handle(ctx context.Context, cmd UpdateTemplateCommand) error {
	if cmd.TemplateID < 1 {
		return errors.New("template id is invalid")
	}

	temp := template.NewEmptyEmailTemplate()

	if topic := cmd.Topic; topic != nil {
		if len(*topic) > 0 {
			temp.SetTopic(*topic)
		} else {
			return errors.New("topic should not be empty")
		}
	}
	if content := cmd.Content; content != nil {
		temp.SetContent(*content)
	}

	return h.templateRepo.Update(ctx, temp)
}
