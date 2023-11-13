package adapter

import (
	"app/calendar/event"
	"context"
	"errors"

	"gorm.io/gorm"
)

var _ event.Repository = (*PostgresRepository)(nil)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) *PostgresRepository {
	return &PostgresRepository{db: db.Table("calendar_events")}
}

func (repository *PostgresRepository) Create(ctx context.Context, event event.Event) error {
	return gormErr(repository.db.WithContext(ctx).Create(&event))
}

func (repository *PostgresRepository) FindByID(ctx context.Context, id string) (event.Event, error) {
	var evt event.Event
	return evt, gormErr(repository.db.WithContext(ctx).First(&evt, "id = ?", id))
}

func (repository *PostgresRepository) FindByCalendarID(ctx context.Context, calendarID string) ([]event.Event, error) {
	var events []event.Event
	return events, gormErr(repository.db.WithContext(ctx).Find(&events, "calendar_id = ?", calendarID))
}

func (repository *PostgresRepository) Delete(ctx context.Context, id string) error {
	return gormErr(repository.db.WithContext(ctx).Delete(&event.Event{}, "id = ?", id))
}

func (repository *PostgresRepository) DeleteAll(ctx context.Context) error {
	return gormErr(repository.db.WithContext(ctx).Delete(&event.Event{}))
}

func gormErr(result *gorm.DB) error {
	switch {
	case result.Error == nil:
		return nil
	case errors.Is(result.Error, gorm.ErrRecordNotFound):
		return event.ErrNotFound
	default:
		return event.ErrRepository
	}
}
