package command

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/weedien/notify-server/common/decorator"
	t "github.com/weedien/notify-server/template/domain/template"
)

type CreateTemplateCommand struct {
	Topic   *string
	Content *string
}

type CreateTemplateHandler decorator.CommandHandler[CreateTemplateCommand]

type createTemplateHandler struct {
	templateRepo t.Repository
}

func NewCreateTemplateHandler(
	templateRepo t.Repository,
	logger *logrus.Entry,
) CreateTemplateHandler {
	if templateRepo == nil {
		panic("templateRepo is nil")
	}
	return decorator.ApplyCommandDecorators[CreateTemplateCommand](
		createTemplateHandler{templateRepo: templateRepo},
		logger,
	)
}

func (h createTemplateHandler) Handle(ctx context.Context, cmd CreateTemplateCommand) error {
	template, err := t.NewEmailTemplate(cmd.Topic, cmd.Content)
	if err != nil {
		return err
	}
	return h.templateRepo.Create(ctx, template)
}
