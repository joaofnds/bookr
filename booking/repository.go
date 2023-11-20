package booking

import "context"

type Repository interface {
	Book(context.Context, Request) error
}
