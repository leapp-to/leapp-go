package api

import (
	"encoding/json"
	"net/http"

	"github.com/leapp-to/leapp-go/pkg/executor"
)

type checkTargetParams struct {
	TargetHost string `json:"target_host"`
	TargetUser string `json:"target_user"`
	Status     bool   `json:"check_target_service_status"`
	Debug      bool   `json:"debug"`
}

func buildCheckTargetInput(p *checkTargetParams) (string, error) {
	data := map[string]interface{}{
		"target_host":                 ObjValue{p.TargetHost},
		"check_target_service_status": ObjValue{p.Status},
		"target_user_name":            ObjValue{p.TargetUser},
	}

	j, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(j), nil
}

func checkTargetHandler(request *http.Request) (*executor.Command, error) {
	var params checkTargetParams

	if err := json.NewDecoder(request.Body).Decode(&params); err != nil {
		return nil, err
	}

	actorInput, err := buildCheckTargetInput(&params)
	if err != nil {
		return nil, err
	}

	c := executor.New("remote-target-check-group", actorInput)
	return c, nil
}
