package api

import (
	"encoding/json"
	"net/http"

	"github.com/leapp-to/leapp-go/pkg/executor"
)

type destroyContainerParams struct {
	ContainerName string `json:"container_name"`
	TargetHost    string `json:"target_host"`
	TargetUser    string `json:"target_user"`
}

func destroyContainerHandler(request *http.Request) (*executor.Command, error) {
	var params destroyContainerParams

	if err := json.NewDecoder(request.Body).Decode(&params); err != nil {
		return nil, err
	}

	d := map[string][]interface{}{
		"container_name":   ChannelData(ObjValue{params.ContainerName}),
		"target_host":      ChannelData(ObjValue{params.TargetHost}),
		"target_user_name": ChannelData(ObjValue{params.TargetUser}),
	}

	actorInput, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}

	c := executor.New("destroy-container", string(actorInput))
	return c, nil
}
