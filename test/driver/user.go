package driver

import (
	"app/test/matchers"
	"app/test/req"
	"app/user"

	"fmt"
	"net/http"
	"strings"
)

type UserDriver struct {
	url string
}

func NewUserDriver(baseURL string) UserDriver {
	return UserDriver{baseURL}
}

func (driver UserDriver) CreateUser(name string) (user.User, error) {
	var u user.User

	return u, makeJSONRequest(params{
		into:   &u,
		status: http.StatusCreated,
		req: func() (*http.Response, error) {
			return req.Post(
				driver.url+"/users",
				map[string]string{"Content-Type": "application/json"},
				strings.NewReader(fmt.Sprintf(`{"name":%q}`, name)),
			)
		},
	})
}

func (driver UserDriver) MustCreateUser(name string) user.User {
	return matchers.Must2(driver.CreateUser(name))
}

func (driver UserDriver) GetUser(name string) (user.User, error) {
	var u user.User

	return u, makeJSONRequest(params{
		into:   &u,
		status: http.StatusOK,
		req: func() (*http.Response, error) {
			return req.Get(
				driver.url+"/users/"+name,
				map[string]string{"Accept": "application/json"},
			)
		},
	})
}

func (driver UserDriver) MustGetUser(name string) user.User {
	return matchers.Must2(driver.GetUser(name))
}

func (driver UserDriver) ListUsers() ([]user.User, error) {
	var users []user.User
	return users, makeJSONRequest(params{
		into:   &users,
		status: http.StatusOK,
		req: func() (*http.Response, error) {
			return req.Get(
				driver.url+"/users",
				map[string]string{"Accept": "application/json"},
			)
		},
	})
}

func (driver UserDriver) MustListUsers() []user.User {
	return matchers.Must2(driver.ListUsers())
}

func (driver UserDriver) DeleteUser(name string) error {
	return makeJSONRequest(params{
		status: http.StatusOK,
		req: func() (*http.Response, error) {
			return req.Delete(
				driver.url+"/users/"+name,
				map[string]string{"Accept": "application/json"},
			)
		},
	})
}

func (driver UserDriver) MustDeleteUser(name string) {
	matchers.Must(driver.DeleteUser(name))
}

func (driver UserDriver) GetFeature(name string) (map[string]any, error) {
	var features map[string]any

	return features, makeJSONRequest(params{
		into:   &features,
		status: http.StatusOK,
		req: func() (*http.Response, error) {
			return req.Get(
				driver.url+"/users/"+name+"/feature",
				map[string]string{"Accept": "application/json"},
			)
		},
	})
}

func (driver UserDriver) MustGetFeature(name string) map[string]any {
	return matchers.Must2(driver.GetFeature(name))
}
