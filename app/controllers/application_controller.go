package controllers

import (
	"github.com/revel/revel"
	"strings"
)

type ApplicationController struct {
	*revel.Controller
}

func (c *ApplicationController) RenderJsonResponse(o interface{}, status int) revel.Result {
	c.Response.Status = status
	return c.RenderJSON(o)
}

func (c *ApplicationController) GetCookiesFromHeaders() (ret map[string]string) {
	cookieString := c.Request.Header.Get("Cookie")
	ret = make(map[string]string)
	if cookieString == "" {
		return ret
	}
	stringArray := strings.Split(cookieString, ";")
	for _, element := range stringArray {
		cookieValArray := strings.Split(element, "=")
		if len(cookieValArray) != 2 {
			continue
		}
		ret[strings.Trim(cookieValArray[0], " ")] = strings.Trim(cookieValArray[1], " ")
	}

	return ret
}
