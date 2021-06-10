package controllers

import (
	"encoding/json"
	"github.com/revel/revel"
)

type ApiController struct {
	ApplicationController
}

var SKIP_AUTHENTICATION_FOR_CONTROLLERS = []string{"AlertV1Controller"}
var SKIP_RULES_FOR_CONTROLLERS = []string{"AlertV1Controller"}
var SKIP_RULES_FOR_ACTIONS = []string{}
var API_STATUS = map[string]string{
	"failure": "failure",
	"success": "success",
}

type Error struct {
	Code     string      `json:"code"`
	Response interface{} `json:"response"`
	Error    string      `json:"error"`
	Event    interface{} `json:"event,omitempty"`
}

// Renders the response in json format
func (c *ApiController) RenderJsonResponse(o interface{}, status int) revel.Result {
	c.Response.Status = status

	return c.RenderJSON(o)
}

func (c *ApiController) RenderTextResponse(o string, status int) revel.Result {
	c.Response.Status = status
	return c.RenderText(o)

}

func (c *ApiController) MakeErrorResp(code string, error string, resp interface{}) interface{} {

	err := Error{
		Code:     code,
		Error:    error,
		Response: resp,
	}
	return err
}

func (c *ApiController) RenderErrorResponse(err error) revel.Result {
	code := err.Error()
	return c.RenderJsonResponse(c.MakeErrorResp("400", code, nil), 400)
}

// Bind params to given struct type
func (c *ApiController) BindParams(content []byte, request_struct interface{}) interface{} {
	err := json.Unmarshal(content, &request_struct)
	if err != nil {
		panic("SERVER_101")
	}
	return request_struct
}

// Authenticate all APIs call with basic auth
// It is included in interceptor to check every API call
func (c ApiController) Authenticate() revel.Result {
	return nil
}

func (c ApiController) SaveHeaders() revel.Result {

	return nil
}

func init() {
	revel.InterceptMethod(ApiController.SaveHeaders, revel.BEFORE)
}
