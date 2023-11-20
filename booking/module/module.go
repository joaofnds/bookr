package module

import (
	"app/booking"
	"app/booking/http"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"booking",

	fx.Provide(booking.NewService),
	fx.Provide(http.NewController),
)
