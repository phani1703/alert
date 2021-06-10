package alert

import (
	service_manager2 "alert/app/service_manager"
	"github.com/revel/revel"
)

type AlertServiceManager struct {
	service_manager2.ServiceManager
}

func (alertServiceManager AlertServiceManager) New(path string, data interface{}) AlertServiceManager {
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	alertServiceManager.ServiceManager.New(alertServiceManager, path, data, headers)
	return alertServiceManager
}

func (alertServiceManager AlertServiceManager) BaseURL() string {
	return revel.Config.StringDefault("alert.base_url", "")
	//return "http://demo8553856.mockable.io"
}

func (alertServiceManager AlertServiceManager) HystrixCommandName() string {
	irctcHystrixCommand := revel.Config.StringDefault("alert.hystrix.command", "")
	return irctcHystrixCommand
}

func (alertServiceManager AlertServiceManager) AddAuthParams(serviceManager *service_manager2.ServiceManager) {

}

// helper function
// copied from https://golang.org/src/net/http/client.go?h=basicAuth#L347

