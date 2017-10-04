package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/leapp-to/leapp-go/pkg/executor"
)

type destroyContainerParams struct {
	ContainerName string `json:"container_name"`
	TargetHost    string `json:"target_host"`
	TargetUser    string `json:"target_user"`
}

func destroyContainerHandler(request *http.Request) (interface{}, error) {
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
		return "", err
	}

	c := executor.New("destroy-container", string(actorInput))
	r := c.Execute()

	log.Println(r.Stderr)

	var out interface{}
	err = json.Unmarshal([]byte(r.Stdout), &out)
	return out, err
}
