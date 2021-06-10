package service_manager

type ManagerInterface interface {
	BaseURL() string
	HystrixCommandName() string
	AddAuthParams(serviceManager *ServiceManager)
}
