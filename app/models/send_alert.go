package models

import (
	"alert/app/parsers/alert"
	alert2 "alert/app/service_manager/alert"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/revel/revel"
)

type SendSMSResponse struct {
	Success string `json:"success"`
}

func SendAlertCommunication(request alert.SendAlertRequest, c *revel.Controller) (SendSMSResponse, error) {

	alertTeam := &AlertTeam{}
	team, err := alertTeam.GetTeamByName(request.Name)
	if err != nil {
		return SendSMSResponse{}, errors.New("Team not found")
	}

	if len(team.Developers) == 0 {
		return SendSMSResponse{}, errors.New("No Developer Found")
	}

	dev := team.Developers[0]

	return SendCommunication(dev, c)

}

func SendCommunication(dev Developers, c *revel.Controller) (resp SendSMSResponse, err error) {
	requestParams := map[string]interface{}{
		"phone_number": dev.PhoneNumber,
		"message":      "Too many 500",
	}

	seatAvlUrl := fmt.Sprintf("/v3/fd99c100-f88a-4d70-aaf7-393dbbd5d99f")
	serviceManager := alert2.AlertServiceManager{}.New(seatAvlUrl, requestParams)
	serviceManager.C = c
	seatAvlResp, err := serviceManager.Post()

	if err != nil {
		return resp, err
	}

	if unmarshaling_err := json.Unmarshal(seatAvlResp.Body, &resp); unmarshaling_err != nil {
		return resp, unmarshaling_err
	}

	return resp, nil

}
