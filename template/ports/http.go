package ports

import (
	"github.com/gofiber/fiber/v2"
	"strconv"

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
	templateModels, err := s.app.Queries.TemplatesWithType.Handle(c.Context(), query.TemplatesWithTypeQuery{Type: *params.Type})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	templates := make([]Template, 0, len(templateModels))
	for _, templateModel := range templateModels {
		templates = append(templates, Template{
			Content: templateModel.Content,
			Slots:   templateModel.Slots,
			Topic:   templateModel.Topic,
			Type:    templateModel.Type,
		})
	}

	return c.Status(fiber.StatusOK).JSON(templates)
}

func (s *HttpServer) CreateTemplate(c *fiber.Ctx) error {
	var template Template
	if err := c.BodyParser(&template); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	err := s.app.Commands.CreateTemplate.Handle(c.Context(), command.CreateTemplateCommand{
		Content: template.Content,
		Topic:   template.Topic,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(fiber.StatusCreated).SendString("Template created")
}

func (s *HttpServer) DeleteTemplateById(c *fiber.Ctx, tid string) error {
	id, err := strconv.Atoi(tid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid id")
	}

	err = s.app.Commands.DeleteTemplate.Handle(c.Context(), command.DeleteTemplateCommand{TemplateID: id})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(fiber.StatusOK).SendString("Template deleted")
}

func (s *HttpServer) GetTemplateById(c *fiber.Ctx, tid string) error {
	id, err := strconv.Atoi(tid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid id")
	}

	templateModel, err := s.app.Queries.TemplateByID.Handle(c.Context(), query.TemplateIdQuery{ID: id})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	template := Template{
		Content: templateModel.Content,
		Slots:   templateModel.Slots,
		Topic:   templateModel.Topic,
		Type:    templateModel.Type,
	}

	return c.Status(fiber.StatusOK).JSON(template)
}

func (s *HttpServer) UpdateTemplate(c *fiber.Ctx, tid string) error {
	var body UpdateTemplateJSONBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	id, err := strconv.Atoi(tid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid id")
	}

	err = s.app.Commands.UpdateTemplate.Handle(c.Context(), command.UpdateTemplateCommand{
		Content:    body.Content,
		Topic:      body.Topic,
		TemplateID: id,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(fiber.StatusOK).SendString("Template updated")
}
