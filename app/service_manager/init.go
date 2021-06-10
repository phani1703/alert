package service_manager

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/revel/revel"
	"net/http"
)

const (
	MILLISECONDS_IN_SEC = 1000
)

// Package variables
var (
	httpClient *http.Client
)

func DoInit() {
	ignoreSsl := revel.Config.BoolDefault("service_manager.skip_ssl", false)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: ignoreSsl},
	}
	httpClient = &http.Client{Transport: tr}
	ConfigureHystrix()
}

func ConfigureHystrix() {
	var hystrixCommandClass []string

	hytrixCommandString := revel.Config.StringDefault("service_manager.hystrix.command", "")

	if hytrixCommandString == "" {
		hystrixCommandClass = []string{}
	} else {
		hytrixCommandByte := []byte(hytrixCommandString)
		err := json.Unmarshal(hytrixCommandByte, &hystrixCommandClass)
		if err != nil {
			panic(err)
		}
	}

	for _, v := range hystrixCommandClass {
		hystrixCommandName := revel.Config.StringDefault(v, "")

		if hystrixCommandName == "" {
			err_string := fmt.Sprintf("CommandName not found in config for key %s", v)
			panic(err_string)
		}

		hystrix.ConfigureCommand(hystrixCommandName, hystrix.CommandConfig{
			Timeout:               revel.Config.IntDefault("hystrix.timeout_sec", 0) * MILLISECONDS_IN_SEC,
			MaxConcurrentRequests: revel.Config.IntDefault("hystrix.max_conc_requests", 0),
			ErrorPercentThreshold: revel.Config.IntDefault("hystrix.err_percent+threshold", 0),
		})
	}
}
