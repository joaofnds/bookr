package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
)

type Resource struct {
	ID      string        `json:"id"`
	Name    string        `json:"name"`
	Setup   time.Duration `json:"setup"`
	Cleanup time.Duration `json:"cleanup"`
}

type CalendarEventStatus int

const (
	Available CalendarEventStatus = iota
	Booked
	Cancelled
)

type CalendarEvent struct {
	ID       string              `json:"id"`
	Name     string              `json:"name"`
	Status   CalendarEventStatus `json:"status"`
	StartsAt time.Time           `json:"starts_at"`
	EndsAt   time.Time           `json:"ends_at"`
}

type Calendar struct {
	ID         string          `json:"id"`
	ResourceID string          `json:"resource_id"`
	Events     []CalendarEvent `json:"events"`
}

type BookingRequest struct {
	EventName       string    `json:"event_name"`
	ResourceID      string    `json:"resource_id"`
	CalendarID      string    `json:"calendar_id"`
	CalendarEventID string    `json:"calendar_event_id"`
	StartsAt        time.Time `json:"starts_at"`
	EndsAt          time.Time `json:"ends_at"`
}

var (
	desk = Resource{
		ID:      uuid.NewString(),
		Name:    "Desk",
		Setup:   5 * time.Minute,
		Cleanup: 10 * time.Minute,
	}

	meeting = CalendarEvent{
		ID:       uuid.NewString(),
		Name:     "Meeting 1",
		Status:   Booked,
		StartsAt: time.Now(),
		EndsAt:   time.Now().Add(30 * time.Minute),
	}

	available = CalendarEvent{
		ID:       uuid.NewString(),
		Name:     "Available",
		Status:   Available,
		StartsAt: time.Now().Add(1 * time.Hour),
		EndsAt:   time.Now().Add(2 * time.Hour),
	}

	deskCalendar = Calendar{
		ID:         uuid.NewString(),
		ResourceID: desk.ID,
		Events:     []CalendarEvent{meeting, available},
	}

	bookingRequest = BookingRequest{
		EventName:       "Meeting 2",
		ResourceID:      desk.ID,
		CalendarID:      deskCalendar.ID,
		CalendarEventID: available.ID,
		StartsAt:        available.StartsAt,
		EndsAt:          available.EndsAt,
	}
)

func findResource(id string) (Resource, error) {
	if desk.ID == id {
		return desk, nil
	}

	return Resource{}, errors.New("resource not found")
}

func findCalendar(id string) (Calendar, error) {
	if deskCalendar.ID == id {
		return deskCalendar, nil
	}

	return Calendar{}, errors.New("calendar not found")
}

func findCalendarEvent(calendarID string, eventId string) (CalendarEvent, error) {
	for _, event := range deskCalendar.Events {
		if event.ID == eventId {
			return event, nil
		}
	}
	return CalendarEvent{}, nil
}

func addEvent(cal Calendar, event CalendarEvent) Calendar {
	cal.Events = append(cal.Events, event)
	return cal
}

func removeEvent(calendar Calendar, event CalendarEvent) Calendar {
	for i, e := range calendar.Events {
		if e.ID == event.ID {
			calendar.Events = append(calendar.Events[:i], calendar.Events[i+1:]...)
			break
		}
	}
	return calendar
}

func validateBookingRequest(request BookingRequest) error {
	if request.ResourceID == "" {
		return errors.New("resource id is required")
	}

	if request.CalendarID == "" {
		return errors.New("calendar id is required")
	}

	if request.CalendarEventID == "" {
		return errors.New("calendar event id is required")
	}

	if request.StartsAt.IsZero() {
		return errors.New("starts at is required")
	}

	if request.EndsAt.IsZero() {
		return errors.New("ends at is required")
	}

	if request.StartsAt.After(request.EndsAt) {
		return errors.New("starts at must be before ends at")
	}

	if request.StartsAt.Equal(request.EndsAt) {
		return errors.New("starts at must be before ends at")
	}

	if request.StartsAt.Before(time.Now()) {
		return errors.New("starts at must be in the future")
	}

	return nil
}

func printCalendar(cal Calendar) {
	for _, event := range cal.Events {
		fmt.Printf(
			"%s %s - %s (%.0fm): %s\n",
			event.StartsAt.Format(time.DateOnly),
			event.StartsAt.Format(time.Kitchen),
			event.EndsAt.Format(time.Kitchen),
			event.EndsAt.Sub(event.StartsAt).Minutes(),
			event.Name,
		)
	}
}

func main() {
	printCalendar(deskCalendar)
	println()

	err := book(bookingRequest)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	printCalendar(deskCalendar)
}

func book(request BookingRequest) error {
	if err := validateBookingRequest(request); err != nil {
		return err
	}

	_, err := findResource(request.ResourceID)
	if err != nil {
		return err
	}

	cal, err := findCalendar(request.CalendarID)
	if err != nil {
		return err
	}

	event, err := findCalendarEvent(cal.ID, request.CalendarEventID)
	if err != nil {
		return err
	}

	if event.Status != Available {
		return errors.New("event is not available")
	}

	if request.StartsAt.Before(event.StartsAt) || request.EndsAt.After(event.EndsAt) {
		return errors.New("request is not within event time")
	}

	cal = removeEvent(cal, event)
	addEvent(cal, CalendarEvent{
		ID:       uuid.NewString(),
		Name:     request.EventName,
		Status:   Booked,
		StartsAt: request.StartsAt,
		EndsAt:   request.EndsAt,
	})

	return nil
}
