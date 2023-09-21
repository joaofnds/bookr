package resource

import "context"

type Repository interface {
	Create(ctx context.Context, resource Resource) (Resource, error)
	FindByID(ctx context.Context, id string) (Resource, error)
	Delete(ctx context.Context, id string) error

	All(ctx context.Context) ([]Resource, error)
	DeleteAll(ctx context.Context) error
}
