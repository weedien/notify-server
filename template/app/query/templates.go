package query

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/weedien/notify-server/common/decorator"
)

type TemplatesQuery struct {
	Type       *int
	ContentLen *int
}

type TemplatesHandler decorator.QueryHandler[TemplatesQuery, []EmailTemplate]

type templatesHandler struct {
	readModel TemplatesReadModel
}

func NewTemplatesHandler(
	readModel TemplatesReadModel,
	logger *logrus.Entry,
) TemplatesHandler {
	if readModel == nil {
		panic("readModel is nil")
	}
	return decorator.ApplyQueryDecorators[TemplatesQuery, []EmailTemplate](
		templatesHandler{readModel: readModel},
		logger,
	)
}

type TemplatesReadModel interface {
	Templates(ctx context.Context, query TemplatesQuery) ([]EmailTemplate, error)
}

func (h templatesHandler) Handle(ctx context.Context, query TemplatesQuery) ([]EmailTemplate, error) {
	return h.readModel.Templates(ctx, query)
}
