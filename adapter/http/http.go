package http

import (
	"app/adapter/health"
	bookinghttp "app/booking/http"
	eventhttp "app/calendar/event/http"
	"app/kv"
	resourcehttp "app/resource/http"
	userhttp "app/user/http"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"http",
	FiberModule,
	fx.Invoke(func(app *fiber.App, healthController *health.Controller) { healthController.Register(app) }),
	fx.Invoke(func(app *fiber.App, userController *userhttp.Controller) { userController.Register(app) }),
	fx.Invoke(func(app *fiber.App, kvController *kv.Controller) { kvController.Register(app) }),
	fx.Invoke(func(app *fiber.App, resourceController *resourcehttp.Controller) { resourceController.Register(app) }),
	fx.Invoke(func(app *fiber.App, eventController *eventhttp.Controller) { eventController.Register(app) }),
	fx.Invoke(func(app *fiber.App, bookingController *bookinghttp.Controller) { bookingController.Register(app) }),
)
