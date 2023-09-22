package main

import (
	"app/adapter/clock"
	"app/adapter/event"
	"app/adapter/featureflags"
	"app/adapter/health"
	"app/adapter/http"
	"app/adapter/id"
	"app/adapter/logger"
	"app/adapter/metrics"
	"app/adapter/postgres"
	"app/adapter/queue"
	"app/adapter/redis"
	"app/adapter/tracing"
	"app/adapter/validation"
	calendar "app/calendar/module"
	"app/config"
	"app/kv"
	resource "app/resource/module"
	user "app/user/module"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		config.Module,
		logger.Module,
		metrics.Module,
		tracing.Module,
		health.Module,
		validation.Module,
		featureflags.Module,
		id.Module,
		clock.Module,

		event.Module,
		queue.ClientModule,
		http.Module,
		postgres.Module,
		redis.Module,

		user.Module,
		kv.Module,
		resource.Module,
		calendar.Module,
	).Run()
}
