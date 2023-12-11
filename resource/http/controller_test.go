package http_test

import (
	"context"
	"testing"
	"time"

	apphttp "app/adapter/http"
	"app/adapter/logger"
	"app/adapter/postgres"
	"app/adapter/validation"
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
			driver.DriverProvider,
			config.Module,
			apphttp.FiberModule,
			validation.Module,
			postgres.Module,
			resourcemodule.Module,
			test.TestIDProvider,
			test.TestClockProvider,
			fx.Invoke(func(app *fiber.App, controller *resourcehttp.Controller) {
				controller.Register(app)
			}),
			fx.Populate(&resourceService, &app, &idService, &clockService),
		).RequireStart()
	})

	BeforeEach(func() {
		idService.Reset()
		Must(resourceService.DeleteAll(context.Background()))
	})

	AfterAll(func() { fxApp.RequireStop() })

	It("creates a resource", func() {
		first := Must2(app.Resource.Create(resource.CreateResourceDTO{
			OwnID:   "first",
			Setup:   10 * time.Minute,
			Cleanup: 15 * time.Minute,
		}))

		Expect(first.ID).To(Equal(idService.Generated[0]))
		Expect(first.OwnID).To(Equal("first"))
		Expect(first.Setup).To(Equal(10 * time.Minute))
		Expect(first.Cleanup).To(Equal(15 * time.Minute))
		Expect(first.CreatedAt.String()).To(Equal(clockService.Now().String()))
		Expect(first.UpdatedAt.String()).To(Equal(clockService.Now().String()))
	})

	It("lists resources", func() {
		first := Must2(app.Resource.Create(resource.CreateResourceDTO{}))
		second := Must2(app.Resource.Create(resource.CreateResourceDTO{}))

		resources := Must2(app.Resource.All())
		Expect(resources).To(Equal([]resource.Resource{first, second}))
	})

	It("finds a resource by id", func() {
		first := Must2(app.Resource.Create(resource.CreateResourceDTO{}))
		second := Must2(app.Resource.Create(resource.CreateResourceDTO{}))

		Expect(Must2(app.Resource.FindByID(first.ID))).To(Equal(first))
		Expect(Must2(app.Resource.FindByID(second.ID))).To(Equal(second))
	})

	It("deletes a resource", func() {
		first := Must2(app.Resource.Create(resource.CreateResourceDTO{}))
		second := Must2(app.Resource.Create(resource.CreateResourceDTO{}))

		Must(app.Resource.Delete(first.ID))

		all := Must2(app.Resource.All())
		Expect(all).To(Equal([]resource.Resource{second}))
	})
})
