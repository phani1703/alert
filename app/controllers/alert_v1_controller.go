package controllers

import (
	"alert/app/models"
	alert2 "alert/app/parsers/alert"
	"github.com/pkg/errors"
	"github.com/revel/revel"
	"net/http"
)

type AlertV1Controller struct {
	ApiController
}

func (c AlertV1Controller) CreateTeam() revel.Result {
	// get the request body
	var request_body alert2.CreateTeamRequest
	c.BindParams(c.Params.JSON, &request_body)

	// validate the request
	err := ValidateCreateTeamRequest(request_body)
	if err != nil {
		return c.RenderErrorResponse(err)
	}

	// set the team name to db
	err = models.SetTeamToDb(request_body)
	if err != nil {
		return c.RenderErrorResponse(err)
	}

	return c.RenderTextResponse("success", http.StatusOK)

}

func ValidateCreateTeamRequest(request alert2.CreateTeamRequest) error {
	if request.Team.Name == "" {
		return errors.New("Team Name cannot be empty")
	}

	if len(request.Developers) == 0 {
		return errors.New("Developers cannot be empty")
	}

	return nil
}

func (c AlertV1Controller) SendAlert() revel.Result {
	var request_body alert2.SendAlertRequest
	c.BindParams(c.Params.JSON, &request_body)

	err := ValidateSendAlertRequest(request_body)
	if err != nil {
		return c.RenderErrorResponse(err)
	}

	resp, err := models.SendAlertCommunication(request_body, c.Controller)
	if err != nil {
		return c.RenderErrorResponse(err)
	}

	return c.RenderJsonResponse(resp, http.StatusOK)

}

func ValidateSendAlertRequest(request alert2.SendAlertRequest) error {
	if request.Name == "" {
		return errors.New("Team Name cannot be empty")
	}

	return nil
}
