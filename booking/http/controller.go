package http

import (
	"app/booking"
	"errors"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Controller struct {
	logger    *zap.Logger
	validator *validator.Validate
	service   *booking.Service
}

func NewController(
	logger *zap.Logger,
	validator *validator.Validate,
	service *booking.Service,
) *Controller {
	return &Controller{
		logger:    logger.Named("booking-controller"),
		validator: validator,
		service:   service,
	}
}

func (controller *Controller) Register(app *fiber.App) {
	bookings := app.Group("/booking")
	bookings.Post("/", controller.Create)
}

func (controller *Controller) Create(ctx *fiber.Ctx) error {
	var body BookingRequestPayload
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	if err := controller.validator.Struct(body); err != nil {
		return ctx.
			Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"errors": validator.ValidationErrors(err.(validator.ValidationErrors))})
	}

	request := booking.Request{
		Name:            body.Name,
		Description:     body.Description,
		ResourceID:      body.ResourceID,
		CalendarID:      body.CalendarID,
		CalendarEventID: body.CalendarEventID,
		StartsAt:        body.StartsAt.UTC(),
		EndsAt:          body.EndsAt.UTC(),
	}

	err := controller.service.Book(ctx.UserContext(), request)
	switch {
	case err == nil:
		return ctx.SendStatus(http.StatusCreated)
	case errors.Is(err, booking.ErrBooking):
		return ctx.
			Status(http.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	default:
		controller.logger.Error("failed to book", zap.Error(err))
		return ctx.SendStatus(http.StatusInternalServerError)
	}
}

type BookingRequestPayload struct {
	ResourceID      string    `json:"resource_id" validate:"required,uuid4"`
	CalendarID      string    `json:"calendar_id" validate:"required,uuid4"`
	CalendarEventID string    `json:"calendar_event_id" validate:"required,uuid4"`
	Name            string    `json:"name" validate:"required,min=1"`
	Description     string    `json:"description"`
	StartsAt        time.Time `json:"starts_at" validate:"required"`
	EndsAt          time.Time `json:"ends_at" validate:"required"`
}
