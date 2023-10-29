package module

import (
	"app/calendar/event"
	"app/calendar/event/adapter"
	"app/calendar/event/http"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"event",

	fx.Provide(adapter.NewPostgresRepository),
	fx.Provide(func(repo *adapter.PostgresRepository) event.Repository { return repo }),

	fx.Provide(event.NewService),
	fx.Provide(http.NewController),
)
