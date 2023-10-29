package event

type Status int

const (
	Available Status = iota
	Booked
	Cancelled
)
