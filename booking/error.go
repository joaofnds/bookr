package booking

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound   = errors.New("resource not found")
	ErrRepository = errors.New("repository error")

	ErrBooking                = errors.New("booking error")
	ErrMissingResourceID      = fmt.Errorf("%w: missing resource id", ErrBooking)
	ErrMissingCalendarID      = fmt.Errorf("%w: missing calendar id", ErrBooking)
	ErrMissingCalendarEventID = fmt.Errorf("%w: missing calendar event id", ErrBooking)
	ErrMissingStartsAt        = fmt.Errorf("%w: missing starts at", ErrBooking)
	ErrMissingEndsAt          = fmt.Errorf("%w: missing ends at", ErrBooking)
	ErrStartAfterEnd          = fmt.Errorf("%w: starts at must be before ends at", ErrBooking)
	ErrStartEqualEnd          = fmt.Errorf("%w: starts at must be before ends at", ErrBooking)
	ErrStartAfterNow          = fmt.Errorf("%w: starts at must be in the future", ErrBooking)

	ErrEventNotAvailable  = fmt.Errorf("%w: event is not available", ErrBooking)
	ErrNotWithinEventTime = fmt.Errorf("%w: request is not within event time", ErrBooking)
)
