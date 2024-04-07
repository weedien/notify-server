package query

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/weedien/notify-server/common/decorator"
)

type TemplatesWithTypeQuery struct {
	Type int
}

type TemplatesWithTypeHandler decorator.QueryHandler[TemplatesWithTypeQuery, []EmailTemplate]

type templatesWithTypeHandler struct {
	readModel TemplatesWithTypeReadModel
}

func NewTemplatesWithTypeHandler(
	readModel TemplatesWithTypeReadModel,
	logger *logrus.Entry,
) TemplatesWithTypeHandler {
	if readModel == nil {
		panic("readModel is nil")
	}
	return decorator.ApplyQueryDecorators[TemplatesWithTypeQuery, []EmailTemplate](
		templatesWithTypeHandler{readModel: readModel},
		logger,
	)
}

type TemplatesWithTypeReadModel interface {
	FindTemplatesByType(ctx context.Context, Type int) ([]EmailTemplate, error)
}

func (h templatesWithTypeHandler) Handle(ctx context.Context, cmd TemplatesWithTypeQuery) ([]EmailTemplate, error) {
	return h.readModel.FindTemplatesByType(ctx, cmd.Type)
}
