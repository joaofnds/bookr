package http

import (
	"app/calendar/event"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Controller struct {
	service   *event.Service
	validator *validator.Validate
	logger    *zap.Logger
}

func NewController(
	service *event.Service,
	validator *validator.Validate,
	logger *zap.Logger,
) *Controller {
	return &Controller{
		service:   service,
		validator: validator,
		logger:    logger.Named("event-controller"),
	}
}

func (controller *Controller) Register(app *fiber.App) {
	events := app.Group("calendar/:calendarID/events")
	events.Get("/", controller.findByCalendarID)
	events.Post("/", controller.create)

	event := events.Group("/:eventID")
	event.Get("/", controller.findByID)
	event.Delete("/", controller.delete)
}

func (controller *Controller) create(ctx *fiber.Ctx) error {
	var body CreateEventBody
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	if err := controller.validator.Struct(body); err != nil {
		return ctx.
			Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"errors": validator.ValidationErrors(err.(validator.ValidationErrors))})
	}

	cal, err := controller.service.Create(ctx.Context(), event.CreateEventDTO{
		CalendarID:  ctx.Params("calendarID"),
		Name:        body.Name,
		Description: body.Description,
		Status:      body.Status,
		StartsAt:    body.StartsAt,
		EndsAt:      body.EndsAt,
	})
	if err != nil {
		controller.logger.Error("failed to create event", zap.Error(err))
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return ctx.Status(http.StatusCreated).JSON(cal)
}

func (controller *Controller) findByCalendarID(ctx *fiber.Ctx) error {
	events, err := controller.service.FindByCalendarID(ctx.Context(), ctx.Params("calendarID"))
	if err != nil {
		controller.logger.Error("failed to find events by calendar id", zap.Error(err))
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return ctx.JSON(events)
}

func (controller *Controller) findByID(ctx *fiber.Ctx) error {
	event, err := controller.service.FindByID(ctx.Context(), ctx.Params("eventID"))
	if err != nil {
		controller.logger.Error("failed to find event by id", zap.Error(err))
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return ctx.JSON(event)
}

func (controller *Controller) delete(ctx *fiber.Ctx) error {
	err := controller.service.Delete(ctx.Context(), ctx.Params("eventID"))

	if err != nil {
		controller.logger.Error("failed to delete event", zap.Error(err))
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

type CreateEventBody struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Status      event.Status `json:"status"`
	StartsAt    time.Time    `json:"starts_at"`
	EndsAt      time.Time    `json:"ends_at"`
}
