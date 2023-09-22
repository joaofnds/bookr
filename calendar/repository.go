package calendar

import "context"

type Repository interface {
	Create(context.Context, Calendar) (Calendar, error)
	FindByID(context.Context, string) (Calendar, error)
	DeleteAll(context.Context) error
}
