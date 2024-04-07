package service

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/weedien/notify-server/template/adapters"
	t "github.com/weedien/notify-server/template/app"
	"github.com/weedien/notify-server/template/app/command"
	"github.com/weedien/notify-server/template/app/query"
)

func NewTemplateApp(ctx context.Context) t.Application {
	db := adapters.GetDB()

	templateRepo := adapters.NewEmailTemplateRepository(db)

	logger := logrus.NewEntry(logrus.StandardLogger())

	return t.Application{
		Commands: t.Commands{
			CreateTemplate: command.NewCreateTemplateHandler(templateRepo, logger),
			UpdateTemplate: command.NewUpdateTemplateHandler(templateRepo, logger),
			DeleteTemplate: command.NewDeleteTemplateHandler(templateRepo, logger),
		},
		Queries: t.Queries{
			AllTemplates:      query.NewAllTemplatesHandler(templateRepo, logger),
			TemplatesWithType: query.NewTemplatesWithTypeHandler(templateRepo, logger),
			TemplateByID:      query.NewTemplateIdHandler(templateRepo, logger),
		},
	}
}
