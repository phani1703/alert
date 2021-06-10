package mongodb

var MaxPool int
var PATH string
var DBNAME string
var Method string
var DialTimeout int

func CheckAndInitServiceConnection() {
	if service.baseSession == nil {
		service.URL = PATH
		err := service.New()
		if err != nil {
			panic(err)
		}
	}
}
