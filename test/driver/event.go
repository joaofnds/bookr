package driver

import (
	"app/calendar/event"
	"app/test/req"

	eventhttp "app/calendar/event/http"

	"net/http"
)

type EventDriver struct {
	url string
}

func NewEventDriver(baseURL string) *EventDriver {
	return &EventDriver{baseURL}
}

func (d *EventDriver) Create(calendarID string, dto eventhttp.CreateEventBody) (event.Event, error) {
	var e event.Event

	return e, makeJSONRequest(params{
		into:   &e,
		status: http.StatusCreated,
		req: func() (*http.Response, error) {
			return req.Post(
				d.url+"/calendar/"+calendarID+"/events",
				map[string]string{"Content-Type": "application/json"},
				jsonReader(dto),
			)
		},
	})
}

func (d *EventDriver) FindByID(calendarID string, eventID string) (event.Event, error) {
	var evt event.Event

	return evt, makeJSONRequest(params{
		into:   &evt,
		status: http.StatusOK,
		req: func() (*http.Response, error) {
			return req.Get(d.url+"/calendar/"+calendarID+"/events/"+eventID, nil)
		},
	})
}

func (d *EventDriver) FindByCalendarID(calendarID string) ([]event.Event, error) {
	var e []event.Event

	return e, makeJSONRequest(params{
		into:   &e,
		status: http.StatusOK,
		req: func() (*http.Response, error) {
			return req.Get(d.url+"/calendar/"+calendarID+"/events", nil)
		},
	})
}

func (d *EventDriver) Delete(calendarID string, eventID string) error {
	return makeJSONRequest(params{
		status: http.StatusNoContent,
		req: func() (*http.Response, error) {
			return req.Delete(d.url+"/calendar/"+calendarID+"/events/"+eventID, nil)
		},
	})
}
