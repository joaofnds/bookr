package calendar

import "time"

type Event struct {
	ID       string              `json:"id"`
	Name     string              `json:"name"`
	Status   CalendarEventStatus `json:"status"`
	StartsAt time.Time           `json:"starts_at"`
	EndsAt   time.Time           `json:"ends_at"`
}
