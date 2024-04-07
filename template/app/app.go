package app

import (
	"github.com/weedien/notify-server/template/app/command"
	"github.com/weedien/notify-server/template/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateTemplate command.CreateTemplateHandler
	UpdateTemplate command.UpdateTemplateHandler
	DeleteTemplate command.DeleteTemplateHandler
}

type Queries struct {
	AllTemplates      query.AllTemplatesHandler
	TemplatesWithType query.TemplatesWithTypeHandler
	TemplateByID      query.TemplateIdHandler
}
