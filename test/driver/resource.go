package driver

import (
	"app/resource"
	"app/test/matchers"
	"app/test/req"
	"net/http"
)

type ResourceDriver struct {
	url string
}

func NewResourceDriver(url string) ResourceDriver {
	return ResourceDriver{url: url}
}

func (driver ResourceDriver) All() ([]resource.Resource, error) {
	var resources []resource.Resource
	return resources, makeJSONRequest(params{
		into:   &resources,
		status: http.StatusOK,
		req: func() (*http.Response, error) {
			return req.Get(
				driver.url+"/resources",
				map[string]string{"Content-Type": "application/json"},
			)
		},
	})
}

func (driver ResourceDriver) MustAll() []resource.Resource {
	return matchers.Must2(driver.All())
}

func (driver ResourceDriver) Create(dto resource.CreateResourceDTO) (resource.Resource, error) {
	var newResource resource.Resource
	return newResource, makeJSONRequest(params{
		into:   &newResource,
		status: http.StatusOK,
		req: func() (*http.Response, error) {
			return req.Post(
				driver.url+"/resources",
				map[string]string{"Content-Type": "application/json"},
				jsonBody(map[string]interface{}{
					"own_id":  dto.OwnID,
					"setup":   dto.Setup,
					"cleanup": dto.Cleanup,
				}),
			)
		},
	})
}

func (driver ResourceDriver) MustCreate(dto resource.CreateResourceDTO) resource.Resource {
	return matchers.Must2(driver.Create(dto))
}

func (driver ResourceDriver) FindByID(id string) (resource.Resource, error) {
	var found resource.Resource
	return found, makeJSONRequest(params{
		into:   &found,
		status: http.StatusOK,
		req: func() (*http.Response, error) {
			return req.Get(
				driver.url+"/resources/"+id,
				map[string]string{"Content-Type": "application/json"},
			)
		},
	})
}

func (driver ResourceDriver) MustFindByID(id string) resource.Resource {
	return matchers.Must2(driver.FindByID(id))
}

func (driver ResourceDriver) Delete(id string) error {
	return makeJSONRequest(params{
		status: http.StatusOK,
		req: func() (*http.Response, error) {
			return req.Delete(driver.url+"/resources/"+id, nil)
		},
	})
}

func (driver ResourceDriver) MustDelete(id string) {
	matchers.Must(driver.Delete(id))
}
