package calendar

import "errors"

var (
	ErrNotFound   = errors.New("resource not found")
	ErrRepository = errors.New("repository error")
)
