package id

import (
	"app/internal"

	"github.com/google/uuid"
	"go.uber.org/fx"
)

var _ internal.IDFactory = (*Service)(nil)

var Module = fx.Provide(func() internal.IDFactory { return New() })

func New() *Service {
	return &Service{}
}

type Service struct{}

func (s *Service) New() string {
	return uuid.NewString()
}
