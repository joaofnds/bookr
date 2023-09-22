package calendar

type CalendarEventStatus int

const (
	Available CalendarEventStatus = iota
	Booked
	Cancelled
)
