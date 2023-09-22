package http

import (
	"app/calendar"
	"app/resource"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func NewController(
	service *calendar.Service,
	validator *validator.Validate,
	logger *zap.Logger,
) *Controller {
	return &Controller{
		service:   service,
		validator: validator,
		logger:    logger.Named("calendar-controller"),
	}
}

type Controller struct {
	service   *calendar.Service
	validator *validator.Validate
	logger    *zap.Logger
}

func (controller *Controller) Register(app *fiber.App) {
	calendars := app.Group("/calendar")
	calendars.Post("/", controller.Create)
	calendars.Get("/:id", controller.FindByID)
}

func (controller *Controller) Create(ctx *fiber.Ctx) error {
	var body ResourceIDPayload
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	if err := controller.validator.Struct(body); err != nil {
		return ctx.
			Status(http.StatusBadRequest).
			JSON(fiber.Map{"errors": validator.ValidationErrors(err.(validator.ValidationErrors))})
	}

	cal, err := controller.service.Create(ctx.Context(), body.ResourceID)
	if err != nil {
		controller.logger.Error("failed to create calendar", zap.Error(err))
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	return ctx.JSON(cal)
}

func (controller *Controller) FindByID(ctx *fiber.Ctx) error {
	res, err := controller.service.FindByID(ctx.UserContext(), ctx.Params("id"))
	switch err {
	case nil:
		return ctx.JSON(res)
	case resource.ErrNotFound:
		return ctx.SendStatus(http.StatusNotFound)
	default:
		controller.logger.Error("error finding calendar", zap.Error(err))
		return ctx.SendStatus(http.StatusInternalServerError)
	}
}

type ResourceIDPayload struct {
	ResourceID string `json:"resource_id" validate:"required"`
}
