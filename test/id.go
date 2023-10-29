package test

import (
	"app/internal"

	"github.com/google/uuid"
	"go.uber.org/fx"
)

var _ internal.IDFactory = (*TestIDService)(nil)

var TestIDProvider = fx.Options(
	fx.Provide(func() *TestIDService { return &TestIDService{} }),
	fx.Provide(func(service *TestIDService) internal.IDFactory { return service }),
)

type TestIDService struct {
	Generated []string
}

func (service *TestIDService) New() string {
	id := uuid.NewString()
	service.Generated = append(service.Generated, id)
	return id
}

func (service *TestIDService) Last() string {
	return service.Generated[len(service.Generated)-1]
}

func (service *TestIDService) Reset() {
	clear(service.Generated)
}
