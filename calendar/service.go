package calendar

import (
	"app/internal"
	"context"
)

func NewService(
	id internal.IDFactory,
	clock internal.ClockService,
	repository Repository,
) *Service {
	return &Service{
		id:         id,
		clock:      clock,
		repository: repository,
	}
}

type Service struct {
	id         internal.IDFactory
	clock      internal.ClockService
	repository Repository
}

func (service *Service) Create(ctx context.Context, resourceID string) (Calendar, error) {
	return service.repository.Create(ctx, Calendar{
		ID:         service.id.New(),
		ResourceID: resourceID,
		CreatedAt:  service.clock.Now(),
		UpdatedAt:  service.clock.Now(),
	})
}

func (service *Service) FindByID(ctx context.Context, id string) (Calendar, error) {
	return service.repository.FindByID(ctx, id)
}

func (service *Service) DeleteAll(ctx context.Context) error {
	return service.repository.DeleteAll(ctx)
}
