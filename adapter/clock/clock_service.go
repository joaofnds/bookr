package clock

import (
	"app/internal"
	"time"

	"go.uber.org/fx"
)

var _ internal.ClockService = (*Service)(nil)

var Module = fx.Provide(func() internal.ClockService { return New() })

type Service struct{}

func New() *Service { return &Service{} }

func (s *Service) Now() time.Time { return time.Now().UTC() }
