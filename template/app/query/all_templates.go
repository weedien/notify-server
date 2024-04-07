package query

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/weedien/notify-server/common/decorator"
)

type AllTemplatesQuery struct{}

type AllTemplatesHandler decorator.QueryHandler[AllTemplatesQuery, []EmailTemplate]

type allTemplatesHandler struct {
	readModel AllTemplatesReadModel
}

func NewAllTemplatesHandler(
	readModel AllTemplatesReadModel,
	logger *logrus.Entry,
) AllTemplatesHandler {
	if readModel == nil {
		panic("readModel is nil")
	}
	return decorator.ApplyQueryDecorators[AllTemplatesQuery, []EmailTemplate](
		allTemplatesHandler{readModel: readModel},
		logger,
	)
}

type AllTemplatesReadModel interface {
	AllTemplates(ctx context.Context) ([]EmailTemplate, error)
}

func (h allTemplatesHandler) Handle(ctx context.Context, _ AllTemplatesQuery) ([]EmailTemplate, error) {
	return h.readModel.AllTemplates(ctx)
}
