package service_manager

// Imports
import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/revel/revel"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// Service Manager struct
type ServiceManager struct {
	Path     string
	Data     interface{}
	Headers  map[string]string
	Manager  ManagerInterface
	MetaData map[string]string // for sending auth keys or query params
	C        *revel.Controller
	IsMMT    string
	Name     string
}

// Instantiate new service manager
func (serviceManager *ServiceManager) New(minterface ManagerInterface, path string, data interface{}, headers map[string]string) {

	serviceManager.Manager = minterface
	serviceManager.Path = path
	serviceManager.Data = data
	serviceManager.Headers = headers
	serviceManager.Name = minterface.HystrixCommandName()
	serviceManager.MetaData = map[string]string{}
}

func getBytes(data interface{}) (b []byte, err error) {
	if data == nil {
		return nil, nil
	}
	var ok bool
	if b, ok = data.([]byte); ok {
		return b, nil
	}
	b, err = json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// Call the service
// This is the last layer which calls service
func (serviceManager *ServiceManager) CallNetwork(method string, url string, data interface{}) (serviceManagerResponse ServiceManagerResponse, err error) {

	data_bytes, d_err := getBytes(data)
	if d_err != nil {
		return ServiceManagerResponse{}.ErrorObject(), d_err
	}
	payload := bytes.NewReader(data_bytes)

	if strings.Contains(serviceManager.Headers["Content-Type"], "multipart/form-data") {
		p := data.(*bytes.Buffer)
		payload = bytes.NewReader(p.Bytes())
	}

	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return ServiceManagerResponse{}.ErrorObject(), err
	}

	for key := range serviceManager.Headers {
		req.Header.Set(key, serviceManager.Headers[key])
	}

	// For debugging request dump
	// if true {
	debug(httputil.DumpRequestOut(req, true))
	// }

	response, err := serviceManager.DoHTTPRequest(req)

	if err != nil {
		return ServiceManagerResponse{}.ErrorObject(), err
	}

	// Close response body once function returns
	defer func() {
		response.Body.Close()
	}()

	// For debugging response dump
	// if true {
	debug(httputil.DumpResponse(response, true))
	// }

	// body, err := ioutil.ReadAll(response.Body)
	serviceManagerResponse, err = ServiceManagerResponse{}.New(response)

	return serviceManagerResponse, err
}

// Create NR segment and do http request
func (serviceManager *ServiceManager) DoHTTPRequest(req *http.Request) (*http.Response, error) {
	response, err := httpClient.Do(req)
	return response, err
}

// Build url for the service.
// It adds query params based on request method type
func (serviceManager *ServiceManager) BuildUrl(method string) string {
	// Add auth params
	serviceManager.Manager.AddAuthParams(serviceManager)
	method = strings.ToUpper(method)
	baseUrl, err := url.Parse(serviceManager.Manager.BaseURL())
	if err != nil {
	}
	baseUrl.Path = serviceManager.Path
	baseQuery := baseUrl.Query()
	for key := range serviceManager.MetaData {
		if str, ok := serviceManager.MetaData[key]; ok {
			baseQuery.Set(key, str)
		}
	}

	if method == "GET" {
		var data_map map[string]interface{}
		if serviceManager.Data != nil {
			data_map = serviceManager.Data.(map[string]interface{})
		}
		for key := range data_map {
			if str, ok := data_map[key].(string); ok {
				baseQuery.Set(key, str)
				//Right now net url package does not provide any way of adding array query params. This is the only way
				//https://golang.org/pkg/net/url/
			} else if array_values, ok := data_map[key].([]string); ok {
				key = key + "[]"
				for _, v := range array_values {
					baseQuery.Add(key, v)
				}
			} else {

			}
		}
		// Flush data after adding get params
		serviceManager.Data = nil
	}

	if serviceManager.Headers["Content-Type"] == "multipart/form-data" {
		payload := &bytes.Buffer{}
		writer := multipart.NewWriter(payload)
		data := serviceManager.Data.(map[string]interface{})
		_ = writer.WriteField("mobile", data["mobile"].(string))
		_ = writer.WriteField("consumer", data["consumer"].(string))
		err := writer.Close()
		if err != nil {
			fmt.Println(err)
		}

		serviceManager.Headers["Content-Type"] = writer.FormDataContentType()
		serviceManager.Data = payload

	}

	baseUrl.RawQuery = baseQuery.Encode()

	return baseUrl.String()
}

func (serviceManager *ServiceManager) CallService(method string) (serviceManagerResponse ServiceManagerResponse, err error) {

	url := serviceManager.BuildUrl(method)
	// Adding hytrix to prevent API misbehaviour
	hystrixCommandName := serviceManager.Manager.HystrixCommandName()
	hystrix.Do(hystrixCommandName, func() error {
		defer func() {
			if r := recover(); r != nil {

			}
		}()

		serviceManagerResponse, err = serviceManager.CallNetwork(method, url, serviceManager.Data)
		if err != nil {

			return err
		}
		return nil
	}, func(err error) error {
		// FAllback

		return err
	})
	// Handle 5xx errors
	serviceManagerResponse, err = ValidateResponse(serviceManagerResponse, err)
	return serviceManagerResponse, err
}

// Do Get call with given params
func (serviceManager *ServiceManager) Get() (ServiceManagerResponse, error) {
	return serviceManager.CallService("GET")
}

// Do Post call with given params
func (serviceManager *ServiceManager) Post() (ServiceManagerResponse, error) {
	return serviceManager.CallService("POST")
}

func ValidateResponse(response ServiceManagerResponse, err error) (ServiceManagerResponse, error) {
	if response.StatusCode/100 == 5 {
		return response, fmt.Errorf("default")
	}
	return response, err
}

//func debug(data []byte, err error) {
func debug(data []byte, err error) {
	if err == nil {

	} else {

	}
}
