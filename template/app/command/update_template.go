package command

import (
	"context"
	"errors"

	"github.com/sirupsen/logrus"
	"github.com/weedien/notify-server/common/decorator"
	t "github.com/weedien/notify-server/template/domain/template"
)

type UpdateTemplateCommand struct {
	TemplateID int64
	Topic      *string // 可选
	Content    *string // 可选
}

type UpdateTemplateHandler decorator.CommandHandler[UpdateTemplateCommand]

type updateTemplateHandler struct {
	templateRepo t.Repository
}

func NewUpdateTemplateHandler(
	templateRepo t.Repository,
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
		return errors.New("invalid template id")
	}

	temp, err := t.NewEmailTemplateForUpdate(cmd.TemplateID, cmd.Topic, cmd.Content)
	if err != nil {
		return err
	}

	return h.templateRepo.Update(ctx, temp)
}
