package adapter

import (
	"app/booking"
	"context"
	"errors"

	"gorm.io/gorm"
)

var _ booking.Repository = (*PostgresRepository)(nil)

type PostgresRepository struct {
}

func NewPostgresRepository() *PostgresRepository {
	return &PostgresRepository{}
}

func (repository *PostgresRepository) Book(ctx context.Context, request booking.Request) error {
	return gormErr(nil)
}

func gormErr(result *gorm.DB) error {
	switch {
	case result.Error == nil:
		return nil
	case errors.Is(result.Error, gorm.ErrRecordNotFound):
		return booking.ErrNotFound
	default:
		return booking.ErrRepository
	}
}
