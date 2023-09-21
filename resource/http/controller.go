package http

import (
	"net/http"

	"app/resource"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func NewController(
	service *resource.Service,
	validator *validator.Validate,
	logger *zap.Logger,
) *Controller {
	return &Controller{
		service:   service,
		validator: validator,
		logger:    logger.Named("resource-controller"),
	}
}

type Controller struct {
	service   *resource.Service
	validator *validator.Validate
	logger    *zap.Logger
}

func (controller *Controller) Register(app *fiber.App) {
	resources := app.Group("/resources")
	resources.Get("/", controller.All)
	resources.Post("/", controller.Create)

	resourceID := resources.Group("/:id")
	resourceID.Get("/", controller.FindByID)
	resourceID.Delete("/", controller.Delete)
}

func (controller *Controller) All(ctx *fiber.Ctx) error {
	users, err := controller.service.All(ctx.UserContext())
	if err != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	return ctx.JSON(users)
}

func (controller *Controller) Create(ctx *fiber.Ctx) error {
	var body resource.CreateResourceDTO
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	if err := controller.validator.Struct(body); err != nil {
		return ctx.
			Status(http.StatusBadRequest).
			JSON(fiber.Map{"errors": validator.ValidationErrors(err.(validator.ValidationErrors))})
	}

	newResource, err := controller.service.Create(ctx.UserContext(), body)
	if err != nil {
		controller.logger.Error("failed to create resource", zap.Error(err))
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	return ctx.JSON(newResource)
}

func (controller *Controller) FindByID(ctx *fiber.Ctx) error {
	res, err := controller.service.FindByID(ctx.UserContext(), ctx.Params("id"))
	if err != nil {
		switch err {
		case resource.ErrNotFound:
			return ctx.SendStatus(http.StatusNotFound)
		default:
			controller.logger.Error("error finding resource", zap.Error(err))
			return ctx.SendStatus(http.StatusInternalServerError)
		}
	}

	return ctx.JSON(res)
}

func (controller *Controller) Delete(ctx *fiber.Ctx) error {
	err := controller.service.Delete(ctx.UserContext(), ctx.Params("id"))

	if err != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	return ctx.SendStatus(http.StatusOK)
}
