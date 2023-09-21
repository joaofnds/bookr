package resource

import (
	"app/internal"
	"context"
)

type Service struct {
	id         internal.IDFactory
	clock      internal.ClockService
	repository Repository
}

func NewService(
	id internal.IDFactory,
	clock internal.ClockService,
	repository Repository,
) *Service {
	return &Service{id: id, clock: clock, repository: repository}
}

func (service *Service) Create(ctx context.Context, dto CreateResourceDTO) (Resource, error) {
	return service.repository.Create(ctx, Resource{
		ID:        service.id.New(),
		OwnID:     dto.OwnID,
		Setup:     dto.Setup,
		Cleanup:   dto.Cleanup,
		CreatedAt: service.clock.Now(),
		UpdatedAt: service.clock.Now(),
	})
}

func (service *Service) FindByID(ctx context.Context, id string) (Resource, error) {
	return service.repository.FindByID(ctx, id)
}

func (service *Service) Delete(ctx context.Context, id string) error {
	return service.repository.Delete(ctx, id)
}

func (service *Service) All(ctx context.Context) ([]Resource, error) {
	return service.repository.All(ctx)
}

func (service *Service) DeleteAll(ctx context.Context) error {
	return service.repository.DeleteAll(ctx)
}
