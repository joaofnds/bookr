package booking

import "time"

type Request struct {
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	ResourceID      string    `json:"resource_id"`
	CalendarID      string    `json:"calendar_id"`
	CalendarEventID string    `json:"calendar_event_id"`
	StartsAt        time.Time `json:"starts_at"`
	EndsAt          time.Time `json:"ends_at"`
}
