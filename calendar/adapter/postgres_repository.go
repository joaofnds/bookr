package adapter

import (
	"app/calendar"
	"context"
	"errors"

	"gorm.io/gorm"
)

var _ calendar.Repository = (*PostgresRepository)(nil)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) *PostgresRepository {
	return &PostgresRepository{db}
}

func (repository *PostgresRepository) Create(ctx context.Context, cal calendar.Calendar) (calendar.Calendar, error) {
	return cal, gormErr(repository.db.WithContext(ctx).Create(&cal))
}

func (repository *PostgresRepository) FindByID(ctx context.Context, id string) (calendar.Calendar, error) {
	var cal calendar.Calendar
	return cal, gormErr(repository.db.WithContext(ctx).First(&cal, "id = ?", id))
}

func (repository *PostgresRepository) DeleteAll(ctx context.Context) error {
	return gormErr(repository.db.WithContext(ctx).Delete(&calendar.Calendar{}))
}

func gormErr(result *gorm.DB) error {
	switch {
	case result.Error == nil:
		return nil
	case errors.Is(result.Error, gorm.ErrRecordNotFound):
		return calendar.ErrNotFound
	default:
		return calendar.ErrRepository
	}
}
