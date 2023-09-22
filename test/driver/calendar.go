package driver

import (
	"app/calendar"
	"app/test/req"
	"net/http"
)

type CalendarDriver struct {
	url string
}

func NewCalendarDriver(url string) *CalendarDriver {
	return &CalendarDriver{url: url}
}

func (driver CalendarDriver) Create(resourceID string) (calendar.Calendar, error) {
	var newCalendar calendar.Calendar
	return newCalendar, makeJSONRequest(params{
		into:   &newCalendar,
		status: http.StatusOK,
		req: func() (*http.Response, error) {
			return req.Post(
				driver.url+"/calendar",
				map[string]string{"Content-Type": "application/json"},
				jsonBody(map[string]interface{}{"resource_id": resourceID}),
			)
		},
	})
}

func (driver CalendarDriver) FindByID(id string) (calendar.Calendar, error) {
	var found calendar.Calendar
	return found, makeJSONRequest(params{
		into:   &found,
		status: http.StatusOK,
		req: func() (*http.Response, error) {
			return req.Get(
				driver.url+"/calendar/"+id,
				map[string]string{"Content-Type": "application/json"},
			)
		},
	})
}
