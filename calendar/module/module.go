package module

import (
	"app/calendar"
	"app/calendar/adapter"
	"app/calendar/http"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"calendar",

	fx.Provide(calendar.NewService),

	fx.Provide(adapter.NewPostgresRepository),
	fx.Provide(func(repo *adapter.PostgresRepository) calendar.Repository { return repo }),

	fx.Provide(http.NewController),
)
