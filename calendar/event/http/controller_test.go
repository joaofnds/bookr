package http_test

import (
	apphttp "app/adapter/http"
	"app/adapter/logger"
	"app/adapter/postgres"
	"app/adapter/validation"
	"app/calendar"
	"app/calendar/event"
	eventhttp "app/calendar/event/http"
	eventmodule "app/calendar/event/module"
	calendarhttp "app/calendar/http"
	calendarmodule "app/calendar/module"
	"app/config"
	"app/resource"
	resourcehttp "app/resource/http"
	resourcemodule "app/resource/module"
	"app/test"
	"app/test/driver"
	"app/test/matchers"
	"context"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestEventHTTP(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "event http suite")
}

var _ = Describe("/events", Ordered, func() {
	var (
		app             *driver.Driver
		fxApp           *fxtest.App
		calendarService *calendar.Service
		resourceService *resource.Service
		idService       *test.TestIDService
		clockService    *test.TestClockService

		desk resource.Resource
		cal  calendar.Calendar
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

			resourcemodule.Module,
			calendarmodule.Module,
			eventmodule.Module,
			fx.Invoke(func(
				app *fiber.App,
				resourceController *resourcehttp.Controller,
				calController *calendarhttp.Controller,
				eventController *eventhttp.Controller,
			) {
				resourceController.Register(app)
				calController.Register(app)
				eventController.Register(app)
			}),

			fx.Populate(&calendarService, &resourceService, &app, &idService, &clockService),
		).RequireStart()
	})

	BeforeEach(func() {
		idService.Reset()
		matchers.Must(calendarService.DeleteAll(context.Background()))
		matchers.Must(resourceService.DeleteAll(context.Background()))

		desk = app.Resource.MustCreate(resource.CreateResourceDTO{})
		cal = app.Calendar.MustCreate(desk.ID)
	})

	AfterAll(func() { fxApp.RequireStop() })

	It("creates an event", func() {
		evt, err := app.Event.Create(cal.ID, eventhttp.CreateEventBody{
			Name:        "event",
			Description: "event description",
			Status:      event.Available,
			StartsAt:    clockService.Now().Add(1 * time.Hour),
			EndsAt:      clockService.Now().Add(2 * time.Hour),
		})

		Expect(err).ToNot(HaveOccurred())
		Expect(evt.ID).To(Equal(idService.Last()))
		Expect(evt.CalendarID).To(Equal(cal.ID))
		Expect(evt.Name).To(Equal("event"))
		Expect(evt.Description).To(Equal("event description"))
		Expect(evt.Status).To(Equal(event.Available))
		Expect(evt.StartsAt).To(Equal(clockService.Now().Add(1 * time.Hour)))
		Expect(evt.EndsAt).To(Equal(clockService.Now().Add(2 * time.Hour)))
		Expect(evt.CreatedAt).To(Equal(clockService.Now()))
		Expect(evt.UpdatedAt).To(Equal(clockService.Now()))
	})

	It("finds event by id", func() {
		evt := app.Event.MustCreate(cal.ID, eventhttp.CreateEventBody{
			Name:        "event",
			Description: "event description",
			Status:      event.Available,
			StartsAt:    clockService.Now().Add(1 * time.Hour),
			EndsAt:      clockService.Now().Add(2 * time.Hour),
		})

		found, err := app.Event.FindByID(evt.CalendarID, evt.ID)
		Expect(err).ToNot(HaveOccurred())
		Expect(found).To(Equal(evt))
	})

	It("finds event by calendar id", func() {
		evt := app.Event.MustCreate(cal.ID, eventhttp.CreateEventBody{
			Name:        "event",
			Description: "event description",
			Status:      event.Available,
			StartsAt:    clockService.Now().Add(1 * time.Hour),
			EndsAt:      clockService.Now().Add(2 * time.Hour),
		})

		calEvents, err := app.Event.FindByCalendarID(cal.ID)
		Expect(err).ToNot(HaveOccurred())
		Expect(calEvents).To(Equal([]event.Event{evt}))
	})

	It("deletes an event", func() {
		first := app.Event.MustCreate(cal.ID, eventhttp.CreateEventBody{
			Name:        "first event",
			Description: "first event description",
			Status:      event.Available,
			StartsAt:    clockService.Now().Add(1 * time.Hour),
			EndsAt:      clockService.Now().Add(2 * time.Hour),
		})
		second := app.Event.MustCreate(cal.ID, eventhttp.CreateEventBody{
			Name:        "second event",
			Description: "second event description",
			Status:      event.Available,
			StartsAt:    clockService.Now().Add(1 * time.Hour),
			EndsAt:      clockService.Now().Add(2 * time.Hour),
		})
		third := app.Event.MustCreate(cal.ID, eventhttp.CreateEventBody{
			Name:        "third event",
			Description: "third event description",
			Status:      event.Available,
			StartsAt:    clockService.Now().Add(1 * time.Hour),
			EndsAt:      clockService.Now().Add(2 * time.Hour),
		})

		Expect(app.Event.MustFindByCalendarID(cal.ID)).
			To(Equal([]event.Event{first, second, third}))

		err := app.Event.Delete(cal.ID, second.ID)
		Expect(err).ToNot(HaveOccurred())

		Expect(app.Event.MustFindByCalendarID(cal.ID)).
			To(Equal([]event.Event{first, third}))
	})
})
