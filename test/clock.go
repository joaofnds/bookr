package test

import (
	"app/internal"
	"time"

	"go.uber.org/fx"
)

var _ internal.ClockService = (*TestClockService)(nil)

var TestClockProvider = fx.Options(
	fx.Provide(func() *TestClockService {
		return &TestClockService{now: time.Now().UTC().Truncate(time.Millisecond)}
	}),
	fx.Provide(func(service *TestClockService) internal.ClockService { return service }),
)

type TestClockService struct {
	now time.Time
}

func (service *TestClockService) Now() time.Time { return service.now }
