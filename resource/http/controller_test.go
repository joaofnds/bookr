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
	"app/test/matchers"

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
		matchers.Must(resourceService.DeleteAll(context.Background()))
	})

	AfterAll(func() { fxApp.RequireStop() })

	It("creates a resource", func() {
		first := app.Resource.MustCreate(resource.CreateResourceDTO{
			OwnID:   "first",
			Setup:   10 * time.Minute,
			Cleanup: 15 * time.Minute,
		})

		Expect(first.ID).To(Equal(idService.Generated[0]))
		Expect(first.OwnID).To(Equal("first"))
		Expect(first.Setup).To(Equal(10 * time.Minute))
		Expect(first.Cleanup).To(Equal(15 * time.Minute))
		Expect(first.CreatedAt.String()).To(Equal(clockService.Now().String()))
		Expect(first.UpdatedAt.String()).To(Equal(clockService.Now().String()))
	})

	It("lists resources", func() {
		first := app.Resource.MustCreate(resource.CreateResourceDTO{})
		second := app.Resource.MustCreate(resource.CreateResourceDTO{})

		resources := app.Resource.MustAll()
		Expect(resources).To(Equal([]resource.Resource{first, second}))
	})

	It("finds a resource by id", func() {
		first := app.Resource.MustCreate(resource.CreateResourceDTO{})
		second := app.Resource.MustCreate(resource.CreateResourceDTO{})

		Expect(app.Resource.MustFindByID(first.ID)).To(Equal(first))
		Expect(app.Resource.MustFindByID(second.ID)).To(Equal(second))
	})

	It("deletes a resource", func() {
		first := app.Resource.MustCreate(resource.CreateResourceDTO{})
		second := app.Resource.MustCreate(resource.CreateResourceDTO{})

		app.Resource.MustDelete(first.ID)

		all := app.Resource.MustAll()
		Expect(all).To(Equal([]resource.Resource{second}))
	})
})
