package ports

import (
	"github.com/gofiber/fiber/v2"
	r "github.com/weedien/notify-server/common/result"
	"time"

	t "github.com/weedien/notify-server/template/app"
	"github.com/weedien/notify-server/template/app/command"
	"github.com/weedien/notify-server/template/app/query"
)

type HttpServer struct {
	app t.Application
}

func NewHttpServer(app t.Application) ServerInterface {
	return &HttpServer{app: app}
}

func (s *HttpServer) GetTemplates(c *fiber.Ctx, params GetTemplatesParams) error {
	templateModels, err := s.app.Queries.Templates.Handle(c.Context(), query.TemplatesQuery{
		Type:       params.Type,
		ContentLen: params.ContentLen,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(r.Fail(fiber.StatusInternalServerError, err.Error()))
	}

	templates := make([]Template, 0, len(templateModels))
	for _, templateModel := range templateModels {
		templates = append(templates, Template{
			Id:         templateModel.ID,
			Content:    templateModel.Content,
			Slots:      templateModel.Slots,
			Topic:      templateModel.Topic,
			Type:       templateModel.Type,
			UpdateTime: templateModel.UpdateTime.Format(time.DateTime),
		})
	}

	return c.Status(fiber.StatusOK).JSON(r.SuccessWithData(templates))
}

func (s *HttpServer) CreateTemplate(c *fiber.Ctx) error {
	var body CreateTemplateJSONBody
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if body.Content == nil || body.Topic == nil {
		return fiber.NewError(fiber.StatusBadRequest, "Content or Topic is nil")
	}

	err := s.app.Commands.CreateTemplate.Handle(c.Context(), command.CreateTemplateCommand{
		Content: body.Content,
		Topic:   body.Topic,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(r.CreateSuccess())
}

func (s *HttpServer) DeleteTemplateById(c *fiber.Ctx, tid int64) error {
	err := s.app.Commands.DeleteTemplate.Handle(c.Context(), command.DeleteTemplateCommand{TemplateID: tid})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusNoContent).JSON(r.DeleteSuccess())
}

func (s *HttpServer) GetTemplateById(c *fiber.Ctx, tid int64) error {
	templateModel, err := s.app.Queries.TemplateByID.Handle(c.Context(), query.TemplateIdQuery{ID: tid})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	template := Template{
		Id:         templateModel.ID,
		Content:    templateModel.Content,
		Slots:      templateModel.Slots,
		Topic:      templateModel.Topic,
		Type:       templateModel.Type,
		UpdateTime: templateModel.UpdateTime.Format(time.DateTime),
	}

	return c.Status(fiber.StatusOK).JSON(r.SuccessWithData(template))
}

func (s *HttpServer) UpdateTemplate(c *fiber.Ctx, tid int64) error {
	var body UpdateTemplateJSONBody
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	err := s.app.Commands.UpdateTemplate.Handle(c.Context(), command.UpdateTemplateCommand{
		Content:    body.Content,
		Topic:      body.Topic,
		TemplateID: tid,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(r.Success())
}
