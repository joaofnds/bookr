package event

import "context"

type Repository interface {
	Create(context.Context, Event) error
	FindByID(context.Context, string) (Event, error)
	FindByCalendarID(context.Context, string) ([]Event, error)
	Delete(context.Context, string) error
	DeleteAll(context.Context) error
}
