package query

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/weedien/notify-server/common/decorator"
)

type TemplateIdQuery struct {
	ID int
}

type TemplateIdHandler decorator.QueryHandler[TemplateIdQuery, EmailTemplate]

type templateIdHandler struct {
	readModel TemplateIdReadModel
}

func NewTemplateIdHandler(
	readModel TemplateIdReadModel,
	logger *logrus.Entry,
) TemplateIdHandler {
	if readModel == nil {
		panic("readModel is nil")
	}
	return decorator.ApplyQueryDecorators[TemplateIdQuery, EmailTemplate](
		templateIdHandler{readModel: readModel},
		logger,
	)
}

type TemplateIdReadModel interface {
	FindTemplateByID(ctx context.Context, id int) (EmailTemplate, error)
}

func (h templateIdHandler) Handle(ctx context.Context, q TemplateIdQuery) (EmailTemplate, error) {
	return h.readModel.FindTemplateByID(ctx, q.ID)
}
