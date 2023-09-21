package module

import (
	"app/resource"
	"app/resource/adapter"
	"app/resource/http"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"resource",

	fx.Provide(resource.NewService),

	fx.Provide(adapter.NewPostgresRepository),
	fx.Provide(func(repo *adapter.PostgresRepository) resource.Repository { return repo }),

	fx.Provide(http.NewController),
)
