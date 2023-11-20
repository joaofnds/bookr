package driver

import (
	bookinghttp "app/booking/http"
	"net/http"
)

type BookingDriver struct {
	url string
}

func NewBookingDriver(baseURL string) *BookingDriver {
	return &BookingDriver{baseURL}
}

func (driver BookingDriver) Book(body bookinghttp.BookingRequestPayload) error {
	return makeJSONRequest(params{
		status: http.StatusCreated,
		req: func() (*http.Response, error) {
			return http.Post(driver.url+"/booking", "application/json", jsonBody(body))
		},
	})
}
