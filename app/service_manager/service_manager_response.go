package service_manager

import (
	"io/ioutil"
	"net/http"
)

type ServiceManagerResponse struct {
	Body       []byte
	Status     string // e.g. "200 OK"
	StatusCode int    // e.g. 200
}

func (serviceManagerResponse ServiceManagerResponse) New(response *http.Response) (ServiceManagerResponse, error) {
	body, err := ioutil.ReadAll(response.Body)

	serviceManagerResponse.Body = body
	serviceManagerResponse.Status = response.Status
	serviceManagerResponse.StatusCode = response.StatusCode

	return serviceManagerResponse, err
}

func (serviceManagerResponse ServiceManagerResponse) ErrorObject() ServiceManagerResponse {
	serviceManagerResponse.Body = nil
	serviceManagerResponse.Status = "0"
	serviceManagerResponse.StatusCode = 0

	return serviceManagerResponse
}
