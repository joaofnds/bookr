package http_test

import (
	"context"
	"net/http"
	"testing"

	"app/adapter/featureflags"
	apphttp "app/adapter/http"
	"app/adapter/logger"
	"app/adapter/postgres"
	"app/adapter/validation"
	"app/config"
	"app/test"
	"app/test/driver"
	"app/test/matchers"
	"app/user"
	userhttp "app/user/http"
	usermodule "app/user/module"

	"github.com/gofiber/fiber/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestUserHTTP(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "user http suite")
}

var _ = Describe("/users", Ordered, func() {
	var (
		app         *driver.Driver
		fxApp       *fxtest.App
		userService *user.Service
	)

	BeforeAll(func() {
		var httpConfig apphttp.Config

		fxApp = fxtest.New(
			GinkgoT(),
			logger.NopLoggerProvider,
			test.Queue,
			test.Transaction,
			test.AvailablePortProvider,
			config.Module,
			featureflags.Module,
			apphttp.FiberModule,
			validation.Module,
			postgres.Module,
			usermodule.Module,
			driver.DriverProvider,
			fx.Invoke(func(fiberApp *fiber.App, controller *userhttp.Controller) {
				controller.Register(fiberApp)
			}),
			fx.Populate(&app, &httpConfig, &userService),
		).RequireStart()
	})

	BeforeEach(func() { matchers.Must(userService.DeleteAll(context.Background())) })

	AfterAll(func() { fxApp.RequireStop() })

	It("creates and gets user", func() {
		bob := app.User.MustCreateUser("bob")
		found := app.User.MustGetUser(bob.Name)

		Expect(found).To(Equal(bob))
	})

	It("lists users", func() {
		bob := app.User.MustCreateUser("bob")
		dave := app.User.MustCreateUser("dave")

		users := app.User.MustListUsers()
		Expect(users).To(Equal([]user.User{bob, dave}))
	})

	It("deletes users", func() {
		bob := app.User.MustCreateUser("bob")
		dave := app.User.MustCreateUser("dave")

		app.User.MustDeleteUser(dave.Name)

		_, err := app.User.GetUser(dave.Name)
		Expect(err).To(Equal(driver.RequestFailure{
			Status: http.StatusNotFound,
			Body:   "Not Found",
		}))

		users := app.User.MustListUsers()
		Expect(users).To(Equal([]user.User{bob}))
	})

	It("switches feature flag", func() {
		bob := app.User.MustCreateUser("bob")
		bobFeatures := app.User.MustGetFeature(bob.Name)
		Expect(bobFeatures).To(Equal(map[string]any{"cool-feature": "on"}))

		frank := app.User.MustCreateUser("frank")
		frankFeatures := app.User.MustGetFeature(frank.Name)
		Expect(frankFeatures).To(Equal(map[string]any{"cool-feature": "off"}))
	})
})
