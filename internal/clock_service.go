package internal

import "time"

type ClockService interface {
	Now() time.Time
}
