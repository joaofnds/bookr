package http_test

import (
	apphttp "app/adapter/http"
	"app/adapter/logger"
	"app/adapter/postgres"
	"app/adapter/validation"
	bookinghttp "app/booking/http"
	bookingmodule "app/booking/module"
	"app/calendar/event"
	eventhttp "app/calendar/event/http"
	eventmodule "app/calendar/event/module"
	calendarhttp "app/calendar/http"
	calendarmodule "app/calendar/module"
	"app/config"
	"app/internal"
	"app/resource"
	resourcehttp "app/resource/http"
	resourcemodule "app/resource/module"
	"app/test"
	"app/test/driver"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestBookingHTTP(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "booking http suite")
}

var _ = Describe("/booking", Ordered, func() {
	var (
		app   *driver.Driver
		fxApp *fxtest.App
		clock internal.ClockService
	)

	BeforeAll(func() {
		fxApp = fxtest.New(
			GinkgoT(),
			logger.NopLoggerProvider,
			test.Transaction,
			test.AvailablePortProvider,
			test.TestIDProvider,
			test.TestClockProvider,

			driver.DriverProvider,

			config.Module,
			apphttp.FiberModule,
			validation.Module,
			postgres.Module,

			bookingmodule.Module,
			calendarmodule.Module,
			eventmodule.Module,
			resourcemodule.Module,
			fx.Invoke(func(
				app *fiber.App,
				bookingcontroller *bookinghttp.Controller,
				calcontroller *calendarhttp.Controller,
				eventController *eventhttp.Controller,
				resourceController *resourcehttp.Controller,
			) {
				bookingcontroller.Register(app)
				calcontroller.Register(app)
				eventController.Register(app)
				resourceController.Register(app)
			}),
			fx.Populate(&app, &clock),
		).RequireStart()
	})

	AfterAll(func() { fxApp.RequireStop() })

	It("books on top of an available event", func() {
		res := app.Resource.MustCreate(resource.CreateResourceDTO{})
		cal := app.Calendar.MustCreate(res.ID)
		evt := app.Event.MustCreate(cal.ID, eventhttp.CreateEventBody{
			Name:        "event",
			Description: "event description",
			Status:      event.Available,
			StartsAt:    clock.Now().Add(1 * time.Hour),
			EndsAt:      clock.Now().Add(2 * time.Hour),
		})

		err := app.Booking.Book(bookinghttp.BookingRequestPayload{
			Name:            "test-event-name",
			ResourceID:      res.ID,
			CalendarID:      cal.ID,
			CalendarEventID: evt.ID,
			StartsAt:        evt.StartsAt,
			EndsAt:          evt.EndsAt,
		})

		Expect(err).To(BeNil())
	})

	When("booking on top of a booked event", func() {
		It("returns an error", func() {
			res := app.Resource.MustCreate(resource.CreateResourceDTO{})
			cal := app.Calendar.MustCreate(res.ID)
			evt := app.Event.MustCreate(cal.ID, eventhttp.CreateEventBody{
				Name:        "event",
				Description: "event description",
				Status:      event.Booked,
				StartsAt:    clock.Now().Add(1 * time.Hour),
				EndsAt:      clock.Now().Add(2 * time.Hour),
			})

			err := app.Booking.Book(bookinghttp.BookingRequestPayload{
				Name:            "test-event-name",
				ResourceID:      res.ID,
				CalendarID:      cal.ID,
				CalendarEventID: evt.ID,
				StartsAt:        evt.StartsAt,
				EndsAt:          evt.EndsAt,
			})

			Expect(err).To(MatchError(driver.RequestFailure{
				Status: http.StatusBadRequest,
				Body:   `{"error":"booking error: event is not available"}`,
			}))
		})
	})

	When("start is after end", func() {
		It("returns an error", func() {
			res := app.Resource.MustCreate(resource.CreateResourceDTO{})
			cal := app.Calendar.MustCreate(res.ID)
			evt := app.Event.MustCreate(cal.ID, eventhttp.CreateEventBody{
				Name:        "event",
				Description: "event description",
				Status:      event.Available,
				StartsAt:    clock.Now().Add(1 * time.Hour),
				EndsAt:      clock.Now().Add(2 * time.Hour),
			})

			err := app.Booking.Book(bookinghttp.BookingRequestPayload{
				Name:            "test-event-name",
				ResourceID:      res.ID,
				CalendarID:      cal.ID,
				CalendarEventID: evt.ID,
				StartsAt:        evt.EndsAt,
				EndsAt:          evt.StartsAt,
			})

			Expect(err).To(MatchError(driver.RequestFailure{
				Status: http.StatusBadRequest,
				Body:   `{"error":"booking error: starts at must be before ends at"}`,
			}))
		})
	})

	When("start is equal to end", func() {
		It("returns an error", func() {
			res := app.Resource.MustCreate(resource.CreateResourceDTO{})
			cal := app.Calendar.MustCreate(res.ID)
			evt := app.Event.MustCreate(cal.ID, eventhttp.CreateEventBody{
				Name:        "event",
				Description: "event description",
				Status:      event.Available,
				StartsAt:    clock.Now().Add(1 * time.Hour),
				EndsAt:      clock.Now().Add(2 * time.Hour),
			})

			err := app.Booking.Book(bookinghttp.BookingRequestPayload{
				Name:            "test-event-name",
				ResourceID:      res.ID,
				CalendarID:      cal.ID,
				CalendarEventID: evt.ID,
				StartsAt:        evt.StartsAt,
				EndsAt:          evt.StartsAt,
			})

			Expect(err).To(MatchError(driver.RequestFailure{
				Status: http.StatusBadRequest,
				Body:   `{"error":"booking error: starts at must be before ends at"}`,
			}))
		})
	})

	When("start is in the past", func() {
		It("returns an error", func() {
			res := app.Resource.MustCreate(resource.CreateResourceDTO{})
			cal := app.Calendar.MustCreate(res.ID)
			evt := app.Event.MustCreate(cal.ID, eventhttp.CreateEventBody{
				Name:        "event",
				Description: "event description",
				Status:      event.Available,
				StartsAt:    clock.Now().Add(1 * time.Hour),
				EndsAt:      clock.Now().Add(2 * time.Hour),
			})

			err := app.Booking.Book(bookinghttp.BookingRequestPayload{
				Name:            "test-event-name",
				ResourceID:      res.ID,
				CalendarID:      cal.ID,
				CalendarEventID: evt.ID,
				StartsAt:        evt.StartsAt.Add(-2 * time.Hour),
				EndsAt:          evt.EndsAt,
			})

			Expect(err).To(MatchError(driver.RequestFailure{
				Status: http.StatusBadRequest,
				Body:   `{"error":"booking error: starts at must be in the future"}`,
			}))
		})
	})

	When("start is before event start", func() {
		It("returns an error", func() {
			res := app.Resource.MustCreate(resource.CreateResourceDTO{})
			cal := app.Calendar.MustCreate(res.ID)
			evt := app.Event.MustCreate(cal.ID, eventhttp.CreateEventBody{
				Name:        "event",
				Description: "event description",
				Status:      event.Available,
				StartsAt:    clock.Now().Add(1 * time.Hour),
				EndsAt:      clock.Now().Add(2 * time.Hour),
			})

			err := app.Booking.Book(bookinghttp.BookingRequestPayload{
				Name:            "test-event-name",
				ResourceID:      res.ID,
				CalendarID:      cal.ID,
				CalendarEventID: evt.ID,
				StartsAt:        evt.StartsAt.Add(-5 * time.Minute),
				EndsAt:          evt.EndsAt,
			})

			Expect(err).To(MatchError(driver.RequestFailure{
				Status: http.StatusBadRequest,
				Body:   `{"error":"booking error: request is not within event time"}`,
			}))
		})
	})

	When("end is after event end", func() {
		It("returns an error", func() {
			res := app.Resource.MustCreate(resource.CreateResourceDTO{})
			cal := app.Calendar.MustCreate(res.ID)
			evt := app.Event.MustCreate(cal.ID, eventhttp.CreateEventBody{
				Name:        "event",
				Description: "event description",
				Status:      event.Available,
				StartsAt:    clock.Now().Add(1 * time.Hour),
				EndsAt:      clock.Now().Add(2 * time.Hour),
			})

			err := app.Booking.Book(bookinghttp.BookingRequestPayload{
				Name:            "test-event-name",
				ResourceID:      res.ID,
				CalendarID:      cal.ID,
				CalendarEventID: evt.ID,
				StartsAt:        evt.StartsAt,
				EndsAt:          evt.EndsAt.Add(5 * time.Minute),
			})

			Expect(err).To(MatchError(driver.RequestFailure{
				Status: http.StatusBadRequest,
				Body:   `{"error":"booking error: request is not within event time"}`,
			}))
		})
	})

	It("updates the event", func() {
		res := app.Resource.MustCreate(resource.CreateResourceDTO{})
		cal := app.Calendar.MustCreate(res.ID)
		evt := app.Event.MustCreate(cal.ID, eventhttp.CreateEventBody{
			Name:        "event",
			Description: "event description",
			Status:      event.Available,
			StartsAt:    clock.Now().Add(1 * time.Hour),
			EndsAt:      clock.Now().Add(2 * time.Hour),
		})

		bookingRequest := bookinghttp.BookingRequestPayload{
			Name:            "test-event-name",
			Description:     "test event description",
			ResourceID:      res.ID,
			CalendarID:      cal.ID,
			CalendarEventID: evt.ID,
			StartsAt:        evt.StartsAt,
			EndsAt:          evt.EndsAt,
		}
		err := app.Booking.Book(bookingRequest)

		Expect(err).To(BeNil())

		events := app.Event.MustFindByCalendarID(cal.ID)
		Expect(events).To(HaveLen(1))
		Expect(events[0].Status).To(Equal(event.Booked))
		Expect(events[0].Name).To(Equal(bookingRequest.Name))
		Expect(events[0].Description).To(Equal(bookingRequest.Description))
		Expect(events[0].StartsAt).To(Equal(evt.StartsAt))
		Expect(events[0].EndsAt).To(Equal(evt.EndsAt))
	})
})
