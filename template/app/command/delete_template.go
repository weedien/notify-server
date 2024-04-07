package command

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/weedien/notify-server/common/decorator"
	"github.com/weedien/notify-server/template/domain/template"
)

type DeleteTemplateCommand struct {
	TemplateID int
}

type DeleteTemplateHandler decorator.CommandHandler[DeleteTemplateCommand]

type deleteTemplateHandler struct {
	templateRepo template.Repository
}

func NewDeleteTemplateHandler(
	templateRepo template.Repository,
	logger *logrus.Entry,
) DeleteTemplateHandler {
	if templateRepo == nil {
		panic("templateRepo is nil")
	}
	return decorator.ApplyCommandDecorators[DeleteTemplateCommand](
		deleteTemplateHandler{templateRepo: templateRepo},
		logger,
	)
}

func (h deleteTemplateHandler) Handle(ctx context.Context, cmd DeleteTemplateCommand) error {
	return h.templateRepo.Delete(ctx, cmd.TemplateID)
}
