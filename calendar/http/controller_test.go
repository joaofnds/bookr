package http_test

import (
	"context"
	"testing"
	"time"

	apphttp "app/adapter/http"
	"app/adapter/logger"
	"app/adapter/postgres"
	"app/adapter/validation"
	"app/calendar"
	calendarhttp "app/calendar/http"
	calendarmodule "app/calendar/module"
	"app/config"
	"app/resource"
	resourcehttp "app/resource/http"
	resourcemodule "app/resource/module"
	"app/test"
	"app/test/driver"
	. "app/test/matchers"

	"github.com/gofiber/fiber/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestResourceHTTP(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "resource http suite")
}

var _ = Describe("/resources", Ordered, func() {
	var (
		app             *driver.Driver
		fxApp           *fxtest.App
		calendarService *calendar.Service
		resourceService *resource.Service
		idService       *test.TestIDService
		clockService    *test.TestClockService
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
			fx.Invoke(func(app *fiber.App, resourceController *resourcehttp.Controller, calcontroller *calendarhttp.Controller) {
				resourceController.Register(app)
				calcontroller.Register(app)
			}),
			fx.Populate(&calendarService, &resourceService, &app, &idService, &clockService),
		).RequireStart()
	})

	BeforeEach(func() {
		idService.Reset()
		Must(calendarService.DeleteAll(context.Background()))
		Must(resourceService.DeleteAll(context.Background()))
	})

	AfterAll(func() { fxApp.RequireStop() })

	It("creates a calendar", func() {
		res := Must2(app.Resource.Create(resource.CreateResourceDTO{
			OwnID:   "first",
			Setup:   10 * time.Minute,
			Cleanup: 15 * time.Minute,
		}))
		Expect(res.ID).To(Equal(idService.Generated[0]))

		cal := Must2(app.Calendar.Create(res.ID))

		Expect(cal.ID).To(Equal(idService.Generated[1]))
		Expect(cal.ResourceID).To(Equal(res.ID))
		Expect(cal.CreatedAt.String()).To(Equal(clockService.Now().String()))
		Expect(cal.UpdatedAt.String()).To(Equal(clockService.Now().String()))
	})

	It("finds a resource by id", func() {
		res := Must2(app.Resource.Create(resource.CreateResourceDTO{}))
		cal := Must2(app.Calendar.Create(res.ID))

		Expect(Must2(app.Calendar.FindByID(cal.ID))).To(Equal(cal))
	})
})
