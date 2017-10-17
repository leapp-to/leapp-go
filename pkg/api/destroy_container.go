package api

import (
	"encoding/json"
	"net/http"

	"github.com/leapp-to/leapp-go/pkg/executor"
)

type destroyContainerParams struct {
	TargetHost    string `json:"target_host"`
	TargetUser    string `json:"target_user"`
	ContainerName string `json:"container_name"`
	Debug         bool   `json:"debug"`
}

func destroyContainerHandler(request *http.Request) (*executor.Command, error) {
	var params destroyContainerParams

	if err := json.NewDecoder(request.Body).Decode(&params); err != nil {
		return nil, err
	}

	d := map[string]interface{}{
		"container_name":   ObjValue{params.ContainerName},
		"target_host":      ObjValue{params.TargetHost},
		"target_user_name": ObjValue{params.TargetUser},
	}

	actorInput, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}

	c := executor.New("destroy-container", string(actorInput))
	return c, nil
}
