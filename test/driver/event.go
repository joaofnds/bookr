package driver

import (
	"app/calendar/event"
	"app/test/matchers"
	"app/test/req"

	eventhttp "app/calendar/event/http"

	"net/http"
)

type EventDriver struct {
	url string
}

func NewEventDriver(baseURL string) EventDriver {
	return EventDriver{baseURL}
}

func (driver EventDriver) Create(calendarID string, dto eventhttp.CreateEventBody) (event.Event, error) {
	var e event.Event

	return e, makeJSONRequest(params{
		into:   &e,
		status: http.StatusCreated,
		req: func() (*http.Response, error) {
			return req.Post(
				driver.url+"/calendar/"+calendarID+"/events",
				map[string]string{"Content-Type": "application/json"},
				jsonReader(dto),
			)
		},
	})
}

func (driver EventDriver) MustCreate(calendarID string, dto eventhttp.CreateEventBody) event.Event {
	return matchers.Must2(driver.Create(calendarID, dto))
}

func (driver EventDriver) FindByID(calendarID string, eventID string) (event.Event, error) {
	var evt event.Event

	return evt, makeJSONRequest(params{
		into:   &evt,
		status: http.StatusOK,
		req: func() (*http.Response, error) {
			return req.Get(driver.url+"/calendar/"+calendarID+"/events/"+eventID, nil)
		},
	})
}

func (driver EventDriver) MustFindByID(calendarID string, eventID string) event.Event {
	return matchers.Must2(driver.FindByID(calendarID, eventID))
}

func (driver EventDriver) FindByCalendarID(calendarID string) ([]event.Event, error) {
	var e []event.Event

	return e, makeJSONRequest(params{
		into:   &e,
		status: http.StatusOK,
		req: func() (*http.Response, error) {
			return req.Get(driver.url+"/calendar/"+calendarID+"/events", nil)
		},
	})
}

func (driver EventDriver) MustFindByCalendarID(calendarID string) []event.Event {
	return matchers.Must2(driver.FindByCalendarID(calendarID))
}

func (driver EventDriver) Delete(calendarID string, eventID string) error {
	return makeJSONRequest(params{
		status: http.StatusNoContent,
		req: func() (*http.Response, error) {
			return req.Delete(driver.url+"/calendar/"+calendarID+"/events/"+eventID, nil)
		},
	})
}

func (driver EventDriver) MustDelete(calendarID string, eventID string) {
	matchers.Must(driver.Delete(calendarID, eventID))
}
