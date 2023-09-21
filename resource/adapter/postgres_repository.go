package adapter

import (
	"app/resource"
	"context"
	"errors"

	"gorm.io/gorm"
)

var _ resource.Repository = (*PostgresRepository)(nil)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) *PostgresRepository {
	return &PostgresRepository{db}
}

func (repository *PostgresRepository) Create(ctx context.Context, resource resource.Resource) (resource.Resource, error) {
	return resource, gormErr(repository.db.WithContext(ctx).Create(&resource))
}

func (repository *PostgresRepository) FindByID(ctx context.Context, id string) (resource.Resource, error) {
	var resource resource.Resource
	return resource, gormErr(repository.db.WithContext(ctx).First(&resource, "id = ?", id))
}

func (repository *PostgresRepository) Delete(ctx context.Context, id string) error {
	return gormErr(repository.db.WithContext(ctx).Exec("DELETE FROM resources WHERE id = ?", id))
}

func (repository *PostgresRepository) All(ctx context.Context) ([]resource.Resource, error) {
	var resources []resource.Resource
	return resources, gormErr(repository.db.WithContext(ctx).Find(&resources))
}

func (repository *PostgresRepository) DeleteAll(ctx context.Context) error {
	return gormErr(repository.db.WithContext(ctx).Exec("DELETE FROM resources"))
}

func gormErr(result *gorm.DB) error {
	switch {
	case result.Error == nil:
		return nil
	case errors.Is(result.Error, gorm.ErrRecordNotFound):
		return resource.ErrNotFound
	default:
		return resource.ErrRepository
	}
}
