package event

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

func (service *Service) Create(ctx context.Context, dto CreateEventDTO) (Event, error) {
	evt := Event{
		ID:        service.id.New(),
		CreatedAt: service.clock.Now(),
		UpdatedAt: service.clock.Now(),

		CalendarID:  dto.CalendarID,
		Name:        dto.Name,
		Description: dto.Description,
		Status:      dto.Status,
		StartsAt:    dto.StartsAt,
		EndsAt:      dto.EndsAt,
	}
	return evt, service.repository.Create(ctx, evt)
}

func (service *Service) FindByCalendarID(ctx context.Context, id string) ([]Event, error) {
	return service.repository.FindByCalendarID(ctx, id)
}

func (service *Service) Delete(ctx context.Context, id string) error {
	return service.repository.Delete(ctx, id)
}

func (service *Service) DeleteAll(ctx context.Context) error {
	return service.repository.DeleteAll(ctx)
}
