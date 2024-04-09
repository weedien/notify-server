package service

import (
	"github.com/sirupsen/logrus"
	"github.com/weedien/notify-server/template/adapters"
	t "github.com/weedien/notify-server/template/app"
	"github.com/weedien/notify-server/template/app/command"
	"github.com/weedien/notify-server/template/app/query"
)

func NewTemplateApp() t.Application {
	db := adapters.DB()

	templateRepo := adapters.NewEmailTemplateRepository(db)

	logger := logrus.NewEntry(logrus.StandardLogger())

	return t.Application{
		Commands: t.Commands{
			CreateTemplate: command.NewCreateTemplateHandler(templateRepo, logger),
			UpdateTemplate: command.NewUpdateTemplateHandler(templateRepo, logger),
			DeleteTemplate: command.NewDeleteTemplateHandler(templateRepo, logger),
		},
		Queries: t.Queries{
			Templates:    query.NewTemplatesHandler(templateRepo, logger),
			TemplateByID: query.NewTemplateIdHandler(templateRepo, logger),
		},
	}
}
