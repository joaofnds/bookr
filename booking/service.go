package booking

import (
	"app/calendar"
	"app/calendar/event"
	"app/internal"
	"app/resource"
	"context"
)

type Service struct {
	clock    internal.ClockService
	resource *resource.Service
	calendar *calendar.Service
	event    *event.Service
}

func NewService(
	clock internal.ClockService,
	resource *resource.Service,
	calendar *calendar.Service,
	event *event.Service,
) *Service {
	return &Service{
		clock:    clock,
		resource: resource,
		calendar: calendar,
		event:    event,
	}
}

func (service *Service) Book(ctx context.Context, request Request) error {
	if err := service.validateBookingRequest(request); err != nil {
		return err
	}

	_, err := service.resource.FindByID(ctx, request.ResourceID)
	if err != nil {
		return err
	}

	cal, err := service.calendar.FindByID(ctx, request.CalendarID)
	if err != nil {
		return err
	}

	evt, err := service.event.FindByID(ctx, request.CalendarEventID)
	if err != nil {
		return err
	}

	if evt.Status != event.Available {
		return ErrEventNotAvailable
	}

	if request.StartsAt.Before(evt.StartsAt) || request.EndsAt.After(evt.EndsAt) {
		return ErrNotWithinEventTime
	}

	err = service.event.Delete(ctx, evt.ID)
	if err != nil {
		return err
	}

	_, err = service.event.Create(ctx, event.CreateEventDTO{
		Status:      event.Booked,
		CalendarID:  cal.ID,
		Name:        request.Name,
		Description: request.Description,
		StartsAt:    request.StartsAt,
		EndsAt:      request.EndsAt,
	})

	return err
}

func (service *Service) validateBookingRequest(request Request) error {
	if request.ResourceID == "" {
		return ErrMissingResourceID
	}

	if request.CalendarID == "" {
		return ErrMissingCalendarID
	}

	if request.CalendarEventID == "" {
		return ErrMissingCalendarEventID
	}

	if request.StartsAt.IsZero() {
		return ErrMissingStartsAt
	}

	if request.EndsAt.IsZero() {
		return ErrMissingEndsAt
	}

	if request.StartsAt.After(request.EndsAt) {
		return ErrStartAfterEnd
	}

	if request.StartsAt.Equal(request.EndsAt) {
		return ErrStartEqualEnd
	}

	if request.StartsAt.Before(service.clock.Now()) {
		return ErrStartAfterNow
	}

	return nil
}
