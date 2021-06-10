package controllers

import (
	"github.com/revel/revel"
)

type HomeController struct {
	ApplicationController
}

func (c HomeController) Index() revel.Result {
	if c.Params.Get("number") == "1" {
		panic("MyError")
	}
	return c.Render()
}

func (c HomeController) CreateHealthCheck() revel.Result {
	//healthcheck := &models.HealthCheck{}
	//healthcheck.CreateFirstRecord()
	return c.RenderText("Success")
}

func (c HomeController) HealthCheck() revel.Result {

	return c.RenderText("Success")
}

// Used on the irctc site
func (c HomeController) Mobile() revel.Result {
	return c.Redirect("https://goibibo.com/mobile")
}
